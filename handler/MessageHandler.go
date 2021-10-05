package messageHandler

import (
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	commands "SyedBot/command"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "piss" {
		s.ChannelMessageSend(m.ChannelID, "shid")
	}

	if m.Content == "salam" {
		s.ChannelMessageSend(m.ChannelID, "salam")
	}

	if m.Content == "?github" {
		s.ChannelMessageSend(m.ChannelID, "https://github.com/Monko2k/SyedBot")
	}

	if strings.HasPrefix(m.Content, "?time") {
		clipped := strings.Replace(m.Content, "?time", "", 1)
		clipped = strings.Replace(clipped, " ", "", 1)
		if strings.HasPrefix(clipped, "in ") {
			clipped = strings.Replace(clipped, "in ", "", 1)
			clipped = strings.ReplaceAll(clipped, " ", "%20")
			go commands.TimeIn(s, m, clipped)
		} else if strings.HasPrefix(clipped, "until ") {
			clipped = strings.Replace(clipped, "until ", "", 1)
			go commands.TimeUntil(s, m, clipped)
		} else {
			go commands.Time(s, m, clipped)
		}

	}

	if strings.HasPrefix(m.Content, "?stock ") {
		clipped := strings.Replace(m.Content, "?stock ", "", 1)
		go commands.Stock(s, m, clipped)
	}

	if strings.HasPrefix(m.Content, "?crypto ") {
		clipped := strings.Replace(m.Content, "?crypto ", "", 1)
		clipped = strings.ToLower(clipped)
		droppedchars, _ := regexp.Compile(`[^a-z0-9 _-]`)
		clipped = droppedchars.ReplaceAllString(clipped, "")
		spaces, _ := regexp.Compile(` `)
		clipped = spaces.ReplaceAllString(clipped, "-")
		go commands.Crypto(s, m, clipped)
	}

	if strings.HasPrefix(m.Content, "?wholesome ") {
		clipped := strings.Replace(m.Content, "?wholesome ", "", 1)
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
		message := clipped + " is " + strconv.Itoa(wholesomeamt) + "% wholesome\n" + clipped + wholesomestat

		s.ChannelMessageSend(m.ChannelID, message)
	}

	if strings.HasPrefix(m.Content, "?whitecatify ") {
		clipped := strings.Replace(m.Content, "?whitecatify ", "", 1)
		s.ChannelMessageSend(m.ChannelID, "holy shit guys, "+clipped)
	}

	imsearch, err := regexp.Compile(`^((.|\n)*?)( |^)(([iI]'?[mM])|[iI] [aA][mM]) `)
	if err != nil {
		log.Println(err.Error())
	} else {
		if imsearch.MatchString(m.Content) {
			if rand.Intn(101) < 20 { // should probably make this a changeable setting
				s.ChannelMessageSend(m.ChannelID, "hi "+imsearch.ReplaceAllString(m.Content, ""))
			}
		}
	}

	if strings.HasPrefix(m.Content, "?tweet ") {
		clipped := strings.Replace(m.Content, "?tweet ", "", 1)
		go commands.Tweet(s, m, clipped)

	}

	if strings.HasPrefix(m.Content, "?choose ") {
		clipped := strings.Replace(m.Content, "?choose ", "", 1)
		var divider string
		if strings.Contains(m.Content, ", ") {
			divider = ", "
		} else {
			divider = " "
		}
		options := strings.Split(clipped, divider)
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

	if strings.HasPrefix(m.Content, "?anime ") {
		clipped := strings.Replace(m.Content, "?anime ", "", 1)
		go commands.Anime(s, m, clipped)
	}

	if strings.HasPrefix(m.Content, "?anirand") {
		clipped := strings.Replace(m.Content, "?anirand", "", 1)
		clipped = strings.Replace(clipped, " ", "", 1)
		go commands.AniRand(s, m, clipped)
	}

	if strings.HasPrefix(m.Content, "?anistaff ") {
		clipped := strings.Replace(m.Content, "?anistaff ", "", 1)
		go commands.AniStaff(s, m, clipped)
	}

	if strings.HasPrefix(m.Content, "?anichar ") {
		clipped := strings.Replace(m.Content, "?anichar ", "", 1)
		go commands.AniChar(s, m, clipped)
	}

}