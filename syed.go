package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "SyedBot/config"
	handlers "SyedBot/handler"

	"github.com/bwmarrin/discordgo"
)
func main() {

	config.ConfigInit()
	DiscordToken := config.Config.DiscordToken
	DiscordSession, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalln("Error creating Discord session" + err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	DiscordSession.AddHandler(handlers.MessageHandler)
	DiscordSession.AddHandler(handlers.ReactHandler)
	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages
	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessageReactions

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

