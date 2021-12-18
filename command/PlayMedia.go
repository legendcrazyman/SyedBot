package commands

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
	"layeh.com/gopus"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

func PlayVideo(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	vidregex := regexp.MustCompile(`((e\/)|(v=))[A-Za-z0-9\-\_]+`) //cba to make a better match
	video := vidregex.FindString(arg)
	videoID := video[2:]
	client := youtube.Client{}

	audio, err := client.GetVideo(videoID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid Video URL!")
		return
	}

	format := audio.Formats.FindByItag(140) // only get videos with audio
	stream, _, err := client.GetStream(audio, format)
	if err != nil {
		log.Println(err)
	}
	var channel string
	preCon := false
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			channel = vs.ChannelID
		}
		if vs.UserID == s.State.User.ID {
			preCon = true
		}
	}
	if preCon {
		s.ChannelMessageSend(m.ChannelID, "Bot is already connected to another channel in this server!")
		return
	}
	vc, err := s.ChannelVoiceJoin(g.ID, channel, false, false)
	if err != nil {
		log.Println("Join error:", err)
		return
	}

	run := exec.Command("ffmpeg", "-f", "m4a", "-i", "pipe:", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1" )
	run.Stdin = stream

	ffmpegout, err := run.StdoutPipe()
	if err != nil {
		log.Println("piperror:", err)
		return
	}

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	err = run.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer run.Process.Kill()

	err = vc.Speaking(true)
	if err != nil {
		log.Println("Send error:", err)
	}

	send := make(chan []int16, 2)

	toClose := make(chan bool)
	go func() {
		SendPCM(vc, send)
		toClose <- true
	}()

	for {
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		select {
		case send <- audiobuf:
		case <-toClose:
			break
		}
	}
	time.Sleep(250 * time.Millisecond)
	err = vc.Speaking(false)
	if err != nil {
		log.Println("Send error:", err)
	}
	defer vc.Disconnect()
	defer close(send)
}
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		log.Println("PCM failed")
		return
	}

	opusEncoder, err := gopus.NewEncoder(frameRate, 2, gopus.Audio)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		recv, ok := <-pcm
		if !ok {
			return
		}

		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			log.Println(err)
			return
		}
		
		if v.Ready == false || v.OpusSend == nil {
			return
		}
		v.OpusSend <- opus
	}
}