package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "SyedBot/config"
	handlers "SyedBot/handler"
	utilities "SyedBot/utilities"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds

	err = DiscordSession.Open()
	if err != nil {
		log.Fatalln("Error opening Discord connection" + err.Error())
	}

	log.Println("Bot started")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoAuth := options.Credential{
		Username: config.Config.MongoDB.Username,
		Password: config.Config.MongoDB.Password,
	}
	clientOpts := options.Client().ApplyURI(config.Config.MongoDB.Hostname).SetAuth(mongoAuth)
	utilities.Database, err = mongo.Connect(ctx, clientOpts)
	log.Println("Database connected")
	if err != nil {
		log.Println("Database failed to connect:", err)
	}

	defer func() {
		if err = utilities.Database.Disconnect(ctx); err != nil {
			log.Fatalln("whoops")
			//make it reconnect? for now just restart the bot
		}
	}()


	//Run until term signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc

	//Close the bot
	DiscordSession.Close()
}

