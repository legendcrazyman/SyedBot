package messageHandler

import (
	"SyedBot/config"
	"math/rand"
	"regexp"
	"strings"

	commands "SyedBot/command"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	imsearch := regexp.MustCompile(`^((.|\n)*?)( |^)(([iI]'?[mM])|[iI] [aA][mM]) `)
	if imsearch.MatchString(m.Content) {
		if rand.Intn(6) ==  1 { // should probably make this a changeable setting
			s.ChannelMessageSend(m.ChannelID, "hi "+imsearch.ReplaceAllString(m.Content, ""))
		}
	}

	if m.Content == "piss" {
		s.ChannelMessageSend(m.ChannelID, "shid")
	}

	if m.Content == "salam" {
		s.ChannelMessageSend(m.ChannelID, "salam")
	}
	
	if !strings.HasPrefix(m.Content, config.Config.Prefix) { 
		return
	}
	message, _ := m.ContentWithMoreMentionsReplaced(s)
	split := strings.SplitN(message[1:], " ", 2)
	splitWithMentions := strings.SplitN(m.Content[1:], " ", 2)
	cmd := split[0]
	var arg string
	var argWithMentions string
	if len(split) > 1 {
		arg = split[1]
		argWithMentions = splitWithMentions[1]
	} else {
		arg = ""
		argWithMentions = ""
	}
	

	switch(cmd) {
		case "github":
			s.ChannelMessageSend(m.ChannelID, "https://github.com/Monko2k/SyedBot")
		case "time": 
			go commands.Time(s, m, arg)
		case "timein":
			go commands.TimeIn(s, m, arg)
		case "timeuntil": 
			go commands.TimeUntil(s, m, arg)
		case "stock":
			go commands.Stock(s, m, arg)
		case "crypto":
			go commands.Crypto(s, m, arg)
		case "wholesome":
			go commands.Wholesome(s, m, arg)
		case "tweet":
			go commands.Tweet(s, m, arg)
		case "retweet":
			go commands.Retweet(s, m, arg)
		case "reply":
			go commands.Reply(s, m, arg)
		case "choose": 
			go commands.Choose(s, m, arg)
		case "anime":
			go commands.Anime(s, m, arg)
		case "anirand":
			go commands.AniRand(s, m, arg)
		case "anistaff":
			go commands.AniStaff(s, m, arg)
		case "anichar": 
			go commands.AniChar(s, m, arg)
		case "rename":
			go commands.Rename(s, m, argWithMentions)
		case "setavatar":
			go commands.SetAvatar(s, m, arg)
		case "play": 
			go commands.PlayVideo(s, m, arg)
	} 
}

