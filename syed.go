package main
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	session.AddHandler(messageHandler)
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

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "piss" {
		s.ChannelMessageSend(m.ChannelID, "shid")
	}

	if m.Content == "syed" {
		s.ChannelMessageSend(m.ChannelID, "gi?")
	}

	log.Println("lol")
}
