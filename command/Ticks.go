package commands

import (
	"SyedBot/utilities"
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Tick(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if len(m.Mentions) < 1 {
		s.ChannelMessageSend(m.ChannelID, "You must tag a user to tick")
		return
	}
	// move this to utilty function
	idregex := regexp.MustCompile(`<@!*\d+>`)
	quote := idregex.ReplaceAllString(arg, "")
	quote = strings.TrimSpace(quote)
	if quote == "" {
		s.ChannelMessageSend(m.ChannelID, "You must provide a quote")
		return
	}
	id := m.Mentions[0].ID
	var nickname string
	collection := utilities.Database.Database("SyedBot").Collection("nicknames")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() == mongo.ErrNoDocuments {
		nickname = m.Mentions[0].String()
	} else if result.Err() != nil {
		log.Println(result.Err())
		return
	} else {
		data := struct {
			Nickname string `bson:"nickname"`
		}{}
		err := result.Decode(&data)
		if err != nil {
			nickname = m.Mentions[0].String()
		} else {
			nickname = data.Nickname
		}
	}
	collection = utilities.Database.Database("SyedBot").Collection("quotes")
	_, err := collection.InsertOne(ctx, 
		bson.D{
			{"user", id}, 
			{"quote", quote}, 
			{"nickname", nickname}, 
			{"time", primitive.NewDateTimeFromTime(time.Now())},
		},
	)

	if err != nil {
		log.Println(err)
	}
}

func SetNick(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	if len(m.Mentions) < 1 {
		s.ChannelMessageSend(m.ChannelID, "You must tag a user to set a nickname")
		return
	}
	// move this to utilty function
	idregex := regexp.MustCompile(`<@!*\d+>`)
	nick := idregex.ReplaceAllString(arg, "")
	nick = strings.TrimSpace(nick)
	if nick == "" {
		s.ChannelMessageSend(m.ChannelID, "You must provide a nickname")
		return
	}
	id := m.Mentions[0].ID
	collection := utilities.Database.Database("SyedBot").Collection("nicknames")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, 
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"nickname", nick}}},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Println(err)
	}
	collection = utilities.Database.Database("SyedBot").Collection("quotes")
	_, err = collection.UpdateMany(ctx, 
		bson.M{"user": id},
		bson.D{
			{"$set", bson.D{{"nickname", nick}}},
		},
	)

	if err != nil {
		log.Println(err)
	}

}


