package commands

import (
	utilities "SyedBot/utilities"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Wholesome(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	wholesomeamt := rand.Intn(101)
	var wholesomestat string
	if wholesomeamt == 0 {
		wholesomestat = " is the least wholesome of them all."
	} else if wholesomeamt < 25 {
		wholesomestat = " is definitively unwholesome."
	} else if wholesomeamt < 50 {
		wholesomestat = " is pretty unwholesome."
	} else if wholesomeamt < 75 {
		wholesomestat = " is pretty wholesome!"
	} else if wholesomeamt < 100 {
		wholesomestat = " is incredibly wholesome!"
	} else {
		wholesomestat = " is super freaking wholesome!"
	}
	message := fmt.Sprintf("%s%s (%d%%)", arg, wholesomestat, wholesomeamt)
	s.ChannelMessageSend(m.ChannelID, message)
}

func Choose(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	var divider string
	if strings.Contains(m.Content, ", ") {
		divider = ", "
	} else {
		divider = " "
	}
	options := strings.Split(arg, divider)
	log.Println(len(options))
	for _, x := range options {
		log.Println(x)
	}
	if len(options) == 0 {
		return
	} else if len(options) == 1 {
		s.ChannelMessageSend(m.ChannelID, options[0])
	} else {
		selection := rand.Intn(len(options))
		s.ChannelMessageSend(m.ChannelID, options[selection])
	}
}

func Rename(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if !utilities.CountVotes(s, m, 4) {
		return
	}
	idregex := regexp.MustCompile(`<@!*\d+>`)
	newname := idregex.ReplaceAllString(arg, "")
	if len(m.Mentions) < 1 {
		s.ChannelMessageSend(m.ChannelID, "You must tag a user to rename!")
	} else {
		err := s.GuildMemberNickname(m.GuildID, m.Mentions[0].ID, newname)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Rename failed!")
			log.Println(err)
		}

	}
}

func SetAvatar(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if !utilities.CountVotes(s, m, 4) {
		return
	}
	urlregex := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	url := urlregex.FindString(m.Content)
	if url != "" {
		head, err := http.Head(url)
		contentType := head.Header.Get("Content-Type")
		if err != nil || !strings.HasPrefix(contentType, "image") {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Invalid URL Content")
			return
		}
		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer res.Body.Close()

		img, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return
		}
		base64img :=  base64.StdEncoding.EncodeToString(img)
		avatar := fmt.Sprintf("data:%s;base64,%s", contentType, base64img)
		_, err = s.UserUpdate("", "", "", avatar, "")
		if err != nil {
			log.Println(err)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "You must include an image URL")
	}
}

func PingVoice(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if len(arg) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Min of 3 characters.")
		return
	}

	vsearch := regexp.MustCompile(`.*` + regexp.QuoteMeta(arg) + `.*`)

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	for _, c := range guild.Channels {
		if c.Type == discordgo.ChannelTypeGuildVoice && vsearch.MatchString(c.Name) {
			users := ""
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == c.ID {
					users += "<@" + vs.UserID + "> "
				}
			}
			if len(users) != 0 {
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> "+users)
				return
			}
		}
	}
}
