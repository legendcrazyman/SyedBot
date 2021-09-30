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

type AniData struct {
	Media struct {
		ID    int `json:"id"`
		Title struct {
			English string `json:"english"`
			Romaji  string `json:"romaji"`
		} `json:"title"`
		Type       string   `json:"type"`
		Genres     []string `json:"genres"`
		CoverImage struct {
			Large string `json:"large"`
			Color string `json:"color"`
		} `json:"coverImage"`
		Status			  string 	  `json:"status"`
		Season            string      `json:"season"`
		SeasonYear        int         `json:"seasonYear"`
		Episodes          int 		  `json:"episodes"`
		AverageScore      int         `json:"averageScore"`
		MeanScore         int         `json:"meanScore"`
		Description       string      `json:"description"`
		NextAiringEpisode struct {
			AiringAt int `json:"airingAt"`
			Episode  int `json:"episode"`
		} `json:"nextAiringEpisode"`
	} `json:"Media"`
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
				Media(search: "` + clipped + `", type: ANIME, sort: SEARCH_MATCH	) {
					id
					title {
						romaji
						english
					}
					type
					genres
					coverImage {
						large
						color
					}
					status
					season
					seasonYear
					episodes
					averageScore
					meanScore
					description (asHtml: false)
					nextAiringEpisode {
						airingAt
						episode
					}
				}
				
			}
		`)
		var graphqlResponse AniData
		
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
			log.Println(err.Error())
			s.ChannelMessageSend(m.ChannelID, "Anime not found!")
			return
		}
		color := 0xFFFFFF
		if graphqlResponse.Media.CoverImage.Color != "" {
			colorhexstring := strings.Replace(graphqlResponse.Media.CoverImage.Color, "#", "", 1)
			colorvalue, _ := strconv.ParseInt(colorhexstring, 16, 64)
			color = int(colorvalue)
		} 

		var title string
		var subtitle string
		if graphqlResponse.Media.Title.English != "" {
			title = graphqlResponse.Media.Title.English 
			subtitle = "**" + graphqlResponse.Media.Title.Romaji + "**\n\n"
		} else {
			title = graphqlResponse.Media.Title.Romaji
			subtitle = ""
		}
		var genres string
		for i, s := range graphqlResponse.Media.Genres {
			if i == 0 {
				genres += s
			} else {
				genres += ", " + s
			}
		}
		var airingTime string
		if graphqlResponse.Media.Status == "RELEASING" {
			convtime := time.Unix(int64(graphqlResponse.Media.NextAiringEpisode.AiringAt), 0)
			airingTime = "\n**Next Airing: **Episode " + strconv.Itoa(graphqlResponse.Media.NextAiringEpisode.Episode)  + " on " + convtime.Month().String() + " " + strconv.Itoa(convtime.Day()) + " " + strconv.Itoa(convtime.Year())
		} else {
			airingTime = ""
		}
		var episodes string
		if graphqlResponse.Media.Episodes != 0 {
			episodes = "\n**Episodes:  **" + strconv.Itoa(graphqlResponse.Media.Episodes)
		} else {
			episodes = "\n**Not Yet Aired**"
		}
		description := strings.Split(graphqlResponse.Media.Description, "<br>")[0] // only use everything before the first linebreak returned by description

		season := "\n\n**Season:  **" + strings.Title(strings.ToLower(graphqlResponse.Media.Season) + " " + strconv.Itoa(graphqlResponse.Media.SeasonYear))
		
		averageScore := "\n**Average Score:  **" + strconv.Itoa(graphqlResponse.Media.AverageScore) + "%"
		embed := &discordgo.MessageEmbed{
			Author:      	&discordgo.MessageEmbedAuthor{},
			Color:      	color,
			Description: 	subtitle + description + season + episodes  + averageScore + airingTime,
			URL:			"https://anilist.co/anime/" + strconv.Itoa(graphqlResponse.Media.ID),
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Genres",
					Value:  genres,
					Inline: false,
				},
			},
			
			Image: &discordgo.MessageEmbedImage{
				URL: graphqlResponse.Media.CoverImage.Large,
			},
			Title:     title,
		}
		
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
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
