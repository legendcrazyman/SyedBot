package messageHandler

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

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
		droppedchars := regexp.MustCompile(`[^a-z0-9 _-]`)
		clipped = droppedchars.ReplaceAllString(clipped, "")
		spaces := regexp.MustCompile(` `)
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
		message := fmt.Sprintf("%s%s (%d%%)", clipped, wholesomestat, wholesomeamt)
		s.ChannelMessageSend(m.ChannelID, message)
	}

	if strings.HasPrefix(m.Content, "?whitecatify ") {
		clipped := strings.Replace(m.Content, "?whitecatify ", "", 1)
		s.ChannelMessageSend(m.ChannelID, "holy shit guys, "+clipped)
	}

	imsearch := regexp.MustCompile(`^((.|\n)*?)( |^)(([iI]'?[mM])|[iI] [aA][mM]) `)
	if imsearch.MatchString(m.Content) {
		if rand.Intn(6) ==  1 { // should probably make this a changeable setting
			s.ChannelMessageSend(m.ChannelID, "hi "+imsearch.ReplaceAllString(m.Content, ""))
		}
	}
	

	if strings.HasPrefix(m.Content, "?tweet ") {
		if !CountVotes(s, m, 2) {
			return
		}	
		clipped, _ := m.ContentWithMoreMentionsReplaced(s)
		clipped = strings.Replace(clipped, "?tweet ", "", 1)
		go commands.Tweet(s, m, clipped)

	}

	if strings.HasPrefix(m.Content, "?retweet ") {
		if !CountVotes(s, m, 2) {
			return
		}	
		clipped := strings.Replace(m.Content, "?retweet ", "", 1)
		go commands.Retweet(s, m, clipped)

	}

	if strings.HasPrefix(m.Content, "?reply ") {
		if !CountVotes(s, m, 2) {
			return
		}	
		clipped, _ := m.ContentWithMoreMentionsReplaced(s)
		clipped = strings.Replace(clipped, "?reply ", "", 1)
		go commands.Reply(s, m, clipped)
	}
	/*
	if strings.HasPrefix(m.Content, "?quote ") {
		clipped := strings.Replace(m.Content, "?quote ", "", 1)
		go commands.Quote(s, m, clipped)
	}*/

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

	// todo: move this somewhere else
	// also, it would probably look better if the command was checked for formatting before the vote counter appears
	if strings.HasPrefix(m.Content, "?rename ") {
		if !CountVotes(s, m, 4) {
			return
		}	
		clipped := strings.Replace(m.Content, "?rename ", "", 1)
		idregex := regexp.MustCompile(`<@!*\d+>`)
		id := idregex.FindString(m.Content)
		clipped = idregex.ReplaceAllString(clipped, "")
		name := strings.TrimSpace(clipped)
		if id != "" {
			idregex = regexp.MustCompile(`\d+`)
			id = idregex.FindString(id)
			err := s.GuildMemberNickname(m.GuildID, id, name)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Rename failed!")
				log.Println(err)
			}

		}
	}

	if strings.HasPrefix(m.Content, "?setavatar ") {
		if !CountVotes(s, m, 4) {
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
	if strings.HasPrefix(m.Content, "?play ") {
		urlregex := regexp.MustCompile(`((e\/)|(v=))[A-Za-z0-9]+`) //cba to make a better match
		video := urlregex.FindString(m.Content)
		clipped := video[2:]

		go commands.PlayVideo(s, m, clipped)
	
	}
}

func CountVotes(s *discordgo.Session, m *discordgo.MessageCreate, amount int) bool {
	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	s.MessageReactionAdd(m.ChannelID, m.ID, "🖕")
	time.Sleep(10 * time.Second)
	reactionMessage, _ := s.ChannelMessage(m.ChannelID, m.ID)

	upvote := 0
	downvote := 0
	for _, x := range reactionMessage.Reactions {
		if x.Emoji.Name == "✅" {
			upvote = x.Count
		} else if x.Emoji.Name == "🖕" {
			downvote = x.Count
		}
	} 
	if upvote > amount && upvote - downvote > 1 {
		return true
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Not enough upvotes! (need at least %d)", amount))
		return false
	}
}
