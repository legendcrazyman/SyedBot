package commands

import (
	"SyedBot/config"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
)

func Tweet(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
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

	if upvote > 3 && upvote-downvote > 2 {
		TwitterSession := anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret)
		tweet, err := TwitterSession.PostTweet(arg, url.Values{})
		if err != nil {
			log.Println("Tweet post failed" + err.Error())
		} else {
			tweeturl := "https://twitter.com/BotSyed/status/" + strconv.Itoa(int(tweet.Id))
			s.ChannelMessageSend(m.ChannelID, tweeturl)
		}
		TwitterSession.Close()
	} else {
		s.ChannelMessageSend(m.ChannelID, "Not enough upvotes! (need at least 3)")
	}	
}