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
	//"fmt"
	
	"github.com/bwmarrin/discordgo"
)

type Config struct { 
	DiscordToken	string
}



func main() {
	var config Config
	readin, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Error: " + err.Error())
	} //put some error handling here later
	_ = json.Unmarshal(readin, &config)
	Token := config.DiscordToken
	
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln("Error creating Discord session" + err.Error())
	}

	session.AddHandler(MessageHandler)
	session.AddHandler(ReactHandler)
	session.Identify.Intents = discordgo.IntentsGuildMessages
	
	err = session.Open()
	if err != nil {
		log.Fatalln("Error opening Discord connection" + err.Error())
	}	
	
	log.Println("Bot started")
	
	//Run until term signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc
	

	//Close the bot
	session.Close()
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

	if m.Content == "test" {
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

	if strings.HasPrefix(m.Content, "?choose ") {
		clipped := strings.Replace(m.Content, "?choose ", "", 1)
		options := strings.Split(clipped, ", ")
		if len(options) == 0 {
			return
		} else if len(options) == 1 {
			s.ChannelMessageSend(m.ChannelID, options[0])
		} else {
			selection := rand.Intn(len(options) - 1)
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
