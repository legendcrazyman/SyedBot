package main
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strings"
	"math/rand"
	"net/url"
	"strconv"
	//"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
)

type Config struct { 
	DiscordToken	string
	Twitter			Twitter

}

type Twitter struct {
	Token		string
	TokenSecret	string
	Key			string
	KeySecret	string
}

var config Config

func init() {
	readin, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Error: " + err.Error())
	} //put some error handling here later
	_ = json.Unmarshal(readin, &config)
}



func main() {

	
	DiscordToken := config.DiscordToken

	DiscordSession, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalln("Error creating Discord session" + err.Error())
	}
	
	rand.Seed(time.Now().UnixNano())
	DiscordSession.AddHandler(MessageHandler)
	DiscordSession.AddHandler(ReactHandler)
	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages
	
	err = DiscordSession.Open()
	if err != nil {
		log.Fatalln("Error opening Discord connection" + err.Error())
	}	
	
	log.Println("Bot started")


	//Run until term signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc
	

	//Close the bot
	DiscordSession.Close()
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	log.Println(m.Content)

	if m.Content == "piss" {
		s.ChannelMessageSend(m.ChannelID, "shid")
	}

	if m.Content == "syed" {
		s.ChannelMessageSend(m.ChannelID, "ji?")
	}

	if m.Content == "salam" {
		s.ChannelMessageSend(m.ChannelID, "salam")
	}

	if m.Content == "dsd" {
		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
		time.Sleep(2 * time.Second)
		reactionMessage, _ := s.ChannelMessage(m.ChannelID, m.ID)

		for _, x := range reactionMessage.Reactions {
		log.Println(reactionMessage.Reactions[0].Emoji)
			if x.Emoji.Name == "✅" && x.Count > 1 {
				s.ChannelMessageSend(m.ChannelID, "yeaaah")
			}
			log.Println(x.Emoji.Name)
		}
	}	

	if strings.HasPrefix(m.Content, "?tweet ") {
		clipped := strings.Replace(m.Content, "?tweet ", "", 1)
		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
		time.Sleep(10 * time.Second)
		reactionMessage, _ := s.ChannelMessage(m.ChannelID, m.ID)

		for _, x := range reactionMessage.Reactions {
		log.Println(reactionMessage.Reactions[0].Emoji)
			if x.Emoji.Name == "✅" {
				if x.Count > 2 {
					TwitterSession := anaconda.NewTwitterApiWithCredentials(config.Twitter.Token, config.Twitter.TokenSecret, config.Twitter.Key, config.Twitter.KeySecret)
					tweet, err := TwitterSession.PostTweet(clipped, url.Values{})
					if err != nil {
						log.Fatalf("Tweet post failed" + err.Error())
					} else {
						tweeturl := "https://twitter.com/BotSyed/status/"+ strconv.Itoa(int(tweet.Id))
						s.ChannelMessageSend(m.ChannelID, tweeturl)
					}
					TwitterSession.Close()
				} else {
					s.ChannelMessageSend(m.ChannelID, "Not enough votes! (need at least 2)")
				}
				
			}
		}
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


}


//Handler doesn't actually detect reactions, not sure why

func ReactHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	log.Println("react")

	if r.UserID == s.State.User.ID {
		return
	}
	log.Println(r.Emoji.Name)
	if r.Emoji.Name == "📌" {
		
		reactionMessage, _ := s.ChannelMessage(r.ChannelID, r.MessageID)
		for _, x := range reactionMessage.Reactions {
			if x.Emoji.Name == "📌" && x.Count > 1 {
				s.ChannelMessagePin(r.ChannelID, r.MessageID)

			}
		}

	}
}
