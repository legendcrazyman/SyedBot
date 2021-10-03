package commands

import (
	"SyedBot/config"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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

		// image tweeting (not well tested)
		urlregex, _ := regexp.Compile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`) // stolen
		if urlregex.MatchString(arg) {
			srcurl := urlregex.FindStringSubmatch(arg)[0]
			method := "GET"

			client := &http.Client {
			}
			req, err := http.NewRequest(method, srcurl, nil)

			if err != nil {
				fmt.Println(err)
				return
			}
			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			TwitterSession := anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret)
			media, err := TwitterSession.UploadMedia(base64.StdEncoding.EncodeToString(body))
			if err != nil {
				TwitterSession.Close() //could just pass the same session but don't care lol
				fmt.Println(err)
				TweetText(s, m, arg)
				return	
			}
			vals := url.Values{}
			vals.Set("media_ids", strconv.FormatInt(media.MediaID, 10))
			arg = strings.ReplaceAll(arg, srcurl, "")
			tweet, err := TwitterSession.PostTweet(arg, vals)
			if err != nil {
				log.Println("Tweet post failed" + err.Error())
			} else {
				tweeturl := "https://twitter.com/BotSyed/status/" + strconv.Itoa(int(tweet.Id))
				s.ChannelMessageSend(m.ChannelID, tweeturl)
			}
		} else {
			TweetText(s, m, arg)
		}


		
	} else {
		s.ChannelMessageSend(m.ChannelID, "Not enough upvotes! (need at least 3)")
	}	
}

func TweetText (s *discordgo.Session, m *discordgo.MessageCreate, text string) {
	TwitterSession := anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret)
	tweet, err := TwitterSession.PostTweet(text, url.Values{})
	if err != nil {
		log.Println("Tweet post failed" + err.Error())
	} else {
		tweeturl := "https://twitter.com/BotSyed/status/" + strconv.Itoa(int(tweet.Id))
		s.ChannelMessageSend(m.ChannelID, tweeturl)
	}
	TwitterSession.Close()	
}