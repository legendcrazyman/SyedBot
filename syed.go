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
	"context"
	//"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
	"github.com/machinebox/graphql"
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
			TwitterSession := anaconda.NewTwitterApiWithCredentials(config.Twitter.Token, config.Twitter.TokenSecret, config.Twitter.Key, config.Twitter.KeySecret)
			tweet, err := TwitterSession.PostTweet(clipped, url.Values{})
			if err != nil {
				log.Println("Tweet post failed" + err.Error())
			} else {
				tweeturl := "https://twitter.com/BotSyed/status/"+ strconv.Itoa(int(tweet.Id))
				s.ChannelMessageSend(m.ChannelID, tweeturl)
			}				
			TwitterSession.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, "Not enough upvotes! (need at least 3)")
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

	if strings.HasPrefix(m.Content, "?anime ") {
		clipped := strings.Replace(m.Content, "?anime ", "", 1)
		graphqlClient := graphql.NewClient("https://graphql.anilist.co")
		graphqlRequest := graphql.NewRequest(`
			{
					Media(search: "` + clipped + `", type: ANIME, sort: FAVOURITES_DESC) {
						id
						title {
							romaji
							english
							native
						}
						type
						genres
					}
				
			}
		`)
		var graphqlResponse interface{}
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
			log.Fatalln(err.Error())
		}
		log.Println(graphqlResponse)
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
