package commands

import (
	"SyedBot/config"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
)

/* TODO: make something so that we don't have to make a new twitter api session each time we use a command lol
also, I'm pretty sure all of these commands can be combined into a single one
*/


func countVotes(s *discordgo.Session, m *discordgo.MessageCreate) bool{
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
	if upvote > 3 && upvote - downvote > 2 {
		return true
	} else {
		s.ChannelMessageSend(m.ChannelID, "Not enough upvotes! (need at least 3)")
		return false
	}
}

func Tweet(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	

	if countVotes(s, m) {
		twit := anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret)
		// image tweeting (not well tested)
		urlregex, _ := regexp.Compile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`) // stolen
		if urlregex.MatchString(arg) {

			
			srcurl := urlregex.FindStringSubmatch(arg)[0]
			text := strings.ReplaceAll(arg, srcurl, "")
			client := &http.Client {
			}
			req, err := http.NewRequest("GET", srcurl, nil)

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

			mediatype := res.Header.Get("Content-Type")
			log.Println(mediatype)
			if (strings.HasPrefix(mediatype, "image")) {
				TweetImg(s, m, twit, body, text)
			} else if strings.HasPrefix(mediatype, "video") {
				TweetVid(s, m, twit, body, mediatype, text)
			} else {
				TweetText(s, m, twit, arg)
			}

		} else {
			TweetText(s, m, twit, arg)
		}
		twit.Close()
	}
}

func TweetText (s *discordgo.Session, m *discordgo.MessageCreate, t *anaconda.TwitterApi, text string) {
	tweet, err := t.PostTweet(text, url.Values{})
	if err != nil {
		log.Println("Tweet post failed" + err.Error())
		s.ChannelMessageSend(m.ChannelID, "Tweet post failed")
	} else {
		tweeturl := "https://twitter.com/BotSyed/status/" + tweet.IdStr
		s.ChannelMessageSend(m.ChannelID, tweeturl)
	}
}

func Retweet (s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if countVotes(s, m) {
		id, err := URLtoID(arg)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Retweet failed") 
		} else {
			TwitterSession := anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret)
			_, err := TwitterSession.Retweet(id, true)
			TwitterSession.Close()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Retweet failed") 
			} else {
				s.ChannelMessageSend(m.ChannelID, "done lol") 
			}
		}
	} 
}

func Reply (s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	urlregex, _ := regexp.Compile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	if urlregex.MatchString(arg) {
		if countVotes(s, m) {
			srcurl := urlregex.FindStringSubmatch(arg)[0]
			text := strings.ReplaceAll(arg, srcurl, "") 
			id, err := URLtoID(srcurl)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Reply failed") 
			} else {
				TwitterSession := anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret)
				tweet, err := TwitterSession.GetTweet(id, url.Values{})
				if err != nil {
					log.Println("Invalid Tweet ID")
					return
				}
				vals := url.Values{}
				vals.Set("in_reply_to_status_id", tweet.IdStr)
				status := fmt.Sprintf("@%s %s", tweet.User.ScreenName, text)
				reply, err := TwitterSession.PostTweet(status, vals)
				TwitterSession.Close()
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Reply failed") 
				} else {
					tweeturl := "https://twitter.com/BotSyed/status/" + reply.IdStr
					s.ChannelMessageSend(m.ChannelID, tweeturl) 
				}
	
			}
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please include a Tweet to reply to")
	}
	
}

func TweetVid (s *discordgo.Session, m *discordgo.MessageCreate, t *anaconda.TwitterApi, body []byte, mediatype string, arg string) {
	log.Println(len(body))
	media, err := t.UploadVideoInit(len(body), mediatype)
	if err != nil {
		fmt.Println(err)
		TweetText(s, m, t, arg)
		return	
	}

	chunk := 0
	// 5mb chunks
	for i := 0; i < len(body); i += 5242879 { 
		err = t.UploadVideoAppend(media.MediaIDString, chunk,
			base64.StdEncoding.EncodeToString(
				body[i:int(math.Min(float64(i) + 5242879, float64(len(body))))], // this is disease
			),
		)
		if err != nil {
			log.Println(err.Error())
			return
		}
		chunk++
	}

	videoMedia, err := t.UploadVideoFinalize(media.MediaIDString)
	if err != nil {
		log.Println(err.Error())
		return
	}
	vals := url.Values{}
	vals.Set("media_ids", strconv.FormatInt(videoMedia.MediaID, 10))
	vals.Set("possibly_sensitive", "true")
	tweet, err := t.PostTweet(arg, vals)
	if err != nil {
		log.Println("Tweet post failed" + err.Error())
		s.ChannelMessageSend(m.ChannelID, "Tweet post failed")
	} else {
		tweeturl := "https://twitter.com/BotSyed/status/" + strconv.Itoa(int(tweet.Id))
		s.ChannelMessageSend(m.ChannelID, tweeturl)
	}
}

func TweetImg (s *discordgo.Session, m *discordgo.MessageCreate, t *anaconda.TwitterApi, body []byte, arg string) {
	media, err := t.UploadMedia(base64.StdEncoding.EncodeToString(body))
	if err != nil {
		fmt.Println(err)
		TweetText(s, m, t, arg)
		return	
	}
	vals := url.Values{}
	vals.Set("media_ids", strconv.FormatInt(media.MediaID, 10))
	vals.Set("possibly_sensitive", "true")
	tweet, err := t.PostTweet(arg, vals)
	if err != nil {
		log.Println("Tweet post failed" + err.Error())
		s.ChannelMessageSend(m.ChannelID, "Tweet post failed")
	} else {
		tweeturl := "https://twitter.com/BotSyed/status/" + strconv.Itoa(int(tweet.Id))
		s.ChannelMessageSend(m.ChannelID, tweeturl)
	}
}

func URLtoID (url string) (int64, error) {
	urlclip, _ := regexp.Compile(`(^https:\/\/twitter.com\/.*\/status\/)|(\?.+)`)
	id := urlclip.ReplaceAllString(url, "")
	idint, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, errors.New("ID conversion failed")
	}
	return idint, nil

}