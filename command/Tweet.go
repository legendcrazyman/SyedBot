package commands

import (
	"SyedBot/config"
	"SyedBot/utilities"
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

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
)

/* TODO: make something so that we don't have to make a new twitter api session each time we use a command lol
also, I'm pretty sure all of these commands can be combined into a single one
*/
var twit *anaconda.TwitterApi


func Tweet(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if !utilities.CountVotes(s, m, 2) {
		return
	}
	twit = anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret) //why does this need to go here
	urlregex := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`) // stolen
	text := arg
	vals := url.Values{}
	if urlregex.MatchString(text) {
		srcurl := urlregex.FindStringSubmatch(text)[0]
		text_nourl := strings.ReplaceAll(text, srcurl, "")
		res, err := http.Head(srcurl)
		if err != nil {
			fmt.Println(err)
			return
		}

		mediatype := res.Header.Get("Content-Type")
		log.Println(mediatype)
		if (strings.HasPrefix(mediatype, "image")) {
			err = AppendImg(s, m, srcurl, &vals)
		} else if strings.HasPrefix(mediatype, "video") {
			err = AppendVid(s, m, srcurl, &vals)
		} 
		if err != nil {
			log.Println(err)
		} else {
			text = text_nourl
		}
	}
	tweet, err := twit.PostTweet(text, vals)
	if err != nil {
		log.Println("Tweet post failed" + err.Error())
		s.ChannelMessageSend(m.ChannelID, "Tweet post failed")
	} else {
		tweeturl := "https://twitter.com/BotSyed/status/" + tweet.IdStr
		s.ChannelMessageSend(m.ChannelID, tweeturl)
	}
}

func Retweet (s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if !utilities.CountVotes(s, m, 2) {
		return
	}
	twit = anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret) //why does this need to go here
	id, err := URLtoID(arg)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Retweet failed") 
	} else {
		_, err := twit.Retweet(id, true)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Retweet failed") 
		} else {
			s.ChannelMessageSend(m.ChannelID, "done lol") 
		}
	} 
}

func Reply (s *discordgo.Session, m *discordgo.MessageCreate, arg string) {	
	if !utilities.CountVotes(s, m, 2) {
		return
	}
	urlregex := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	if urlregex.MatchString(arg) {
		twit = anaconda.NewTwitterApiWithCredentials(config.Config.Twitter.Token, config.Config.Twitter.TokenSecret, config.Config.Twitter.Key, config.Config.Twitter.KeySecret) //why does this need to go here
		srcurl := urlregex.FindStringSubmatch(arg)[0]
		text := strings.ReplaceAll(arg, srcurl, "") 
		id, err := URLtoID(srcurl)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Reply failed") 
		} else {
			tweet, err := twit.GetTweet(id, url.Values{})
			if err != nil {
				log.Println("Invalid Tweet ID")
				return
			}
			vals := url.Values{}
			vals.Set("in_reply_to_status_id", tweet.IdStr)
			if urlregex.MatchString(text) {
				srcurl := urlregex.FindStringSubmatch(text)[0]
				text_nourl := strings.ReplaceAll(text, srcurl, "")
				res, err := http.Head(srcurl)
				if err != nil {
					fmt.Println(err)
					return
				}

				mediatype := res.Header.Get("Content-Type")
				log.Println(mediatype)
				if (strings.HasPrefix(mediatype, "image")) {
					err = AppendImg(s, m, srcurl, &vals)
				} else if strings.HasPrefix(mediatype, "video") {
					err = AppendVid(s, m, srcurl, &vals)
				} 
				if err != nil {
					log.Println(err)
				} else {
					text = text_nourl
				}
			} 
			status := fmt.Sprintf("@%s %s", tweet.User.ScreenName, text)
			reply, err := twit.PostTweet(status, vals)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Reply failed") 
			} else {
				tweeturl := "https://twitter.com/BotSyed/status/" + reply.IdStr
				s.ChannelMessageSend(m.ChannelID, tweeturl) 
			}
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please include a Tweet to reply to")
	}
	
}
/*
func Quote (s *discordgo.Session, m*discordgo.MessageCreate, arg string) {
	urlregex := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	if urlregex.MatchString(arg) {
		if CountVotes(s, m) {
			srcurl := urlregex.FindStringSubmatch(arg)[0]
			text := strings.ReplaceAll(arg, srcurl, "") 
			id, err := URLtoID(srcurl)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Quote Tweet failed") 
			} else {
				tweet, err := twit.GetTweet(id, url.Values{})
				if err != nil {
					log.Println("Invalid Tweet ID")
					return
				}
				vals := url.Values{}
				vals.Set("quoted_status_id", tweet.IdStr)
				vals.Set("is_quote_status", "true")
				if urlregex.MatchString(text) {
					srcurl := urlregex.FindStringSubmatch(text)[0]
					text_nourl := strings.ReplaceAll(text, srcurl, "")
					req, err := http.NewRequest("HEAD", srcurl, nil)
					if err != nil {
						fmt.Println(err)
						return
					}
					res, err := client.Do(req)
					if err != nil {
						fmt.Println(err)
						return
					}

					mediatype := res.Header.Get("Content-Type")
					log.Println(mediatype)
					if (strings.HasPrefix(mediatype, "image")) {
						err = AppendImg(s, m, srcurl, &vals)
					} else if strings.HasPrefix(mediatype, "video") {
						err = AppendVid(s, m, srcurl, &vals)
					} 
					if err != nil {
						log.Println(err)
					} else {
						text = text_nourl
					}
				} 
				quote, err := twit.PostTweet(text, vals)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Quote failed") 
				} else {
					tweeturl := "https://twitter.com/BotSyed/status/" + quote.IdStr
					s.ChannelMessageSend(m.ChannelID, tweeturl) 
				}
	
			}
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please include a Tweet to quote")
	}
}
*/
func AppendVid (s *discordgo.Session, m *discordgo.MessageCreate, srcurl string, vals *url.Values) error {
	res, err := http.Get(srcurl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// currently, anaconda does not have a way to set the media_category field
	// so for now, video uploads are limited to 30 seconds in length
	media, err := twit.UploadVideoInit(len(body), res.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	chunk := 0
	// 5mb chunks (api doc says use 1mb chunks? but this works ? ? idk)
	for i := 0; i < len(body); i += 5242879 { 
		err = twit.UploadVideoAppend(media.MediaIDString, chunk,
			base64.StdEncoding.EncodeToString(
				body[i:int(math.Min(float64(i) + 5242879, float64(len(body))))], // this is disease
			),
		)
		if err != nil {
			return err
		}
		chunk++
	}
	videoMedia, err := twit.UploadVideoFinalize(media.MediaIDString)
	if err != nil {
		return err
	}
	vals.Set("media_ids", videoMedia.MediaIDString)
	vals.Set("possibly_sensitive", "true")
	return nil
}

func AppendImg (s *discordgo.Session, m *discordgo.MessageCreate, srcurl string, vals *url.Values) error {
	res, err := http.Get(srcurl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	media, err := twit.UploadMedia(base64.StdEncoding.EncodeToString(body))
	if err != nil {
		return err	
	}
	vals.Set("media_ids", media.MediaIDString)
	vals.Set("possibly_sensitive", "true")
	return nil
}

func URLtoID (url string) (int64, error) {
	urlclip := regexp.MustCompile(`(^https:\/\/twitter.com\/.*\/status\/)|(\?.+)`)
	id := urlclip.ReplaceAllString(url, "")
	idint, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, errors.New("ID conversion failed")
	}
	return idint, nil

}
