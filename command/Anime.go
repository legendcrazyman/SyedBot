package commands

import (
	structs "SyedBot/struct"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)




func Anime(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	query :=  "search: \\\"" + arg + "\\\", type: ANIME, sort: SEARCH_MATCH"
	AnimeMedia(s, m, query, 1)
}
 
func AniRand(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {

	// randomizing anime by choosing a random sort order
	// this list of orders should probably go somewhere else
	orders := [18]string {
		"ID",
		"TITLE_ROMAJI",
		"TITLE_ENGLISH",
		"TITLE_NATIVE",
		"TYPE",
		"FORMAT",
		"START_DATE",
		"END_DATE",
		"SCORE",
		"POPULARITY",
		"TRENDING",
		"EPISODES",
		"DURATION",
		"STATUS",
		"CHAPTERS",
		"VOLUMES",
		"UPDATED_AT",
		"FAVOURITES",
	}
	order := orders[rand.Intn(18)]
	if rand.Intn(2) == 1 {
		order += "_DESC"
	}
	query := "type: ANIME, sort: " + order 
	yearRangeParse, _ := regexp.Compile(`[yY]:[0-9]{4}-[0-9]{4}`)
	yearParse, _ := regexp.Compile(`[yY]:[0-9]{4}`)
	scoreRangeParse, _ := regexp.Compile(`[sS]:[0-9]{1,2}-[0-9]{1,2}`) // this regex does not handle values greater than 99 (highest score on anilist is 92)
	scoreParse, _ := regexp.Compile(`[sS]:[0-9]{1,2}`)
	genreParse, _ := regexp.Compile(`[gG]:.+?(([yY]:)|([sS]:)|$)`)
	if yearRangeParse.MatchString(arg) {
		yearRange := yearRangeParse.FindString(arg)
		yearRange = strings.Replace(strings.ToLower(yearRange), "y:", "", 1)
		years := strings.Split(yearRange, "-")
		query += ", startDate_greater: " + years[0] + "0000"
		query += ", startDate_lesser: " + years[1] + "0000"
	} else if yearParse.MatchString(arg) {
		year := yearParse.FindString(arg)
		year = strings.Replace(strings.ToLower(year), "y:", "", 1)
		query += ", startDate_greater: " + year + "0000"
	}
	if scoreRangeParse.MatchString(arg) {
		scoreRange := scoreRangeParse.FindString(arg)
		scoreRange = strings.Replace(strings.ToLower(scoreRange), "s:", "", 1)
		years := strings.Split(scoreRange, "-")
		query += ", averageScore_greater: " + years[0]
		query += ", averageScore_lesser: " + years[1]
	} else if scoreParse.MatchString(arg) {
		score := scoreParse.FindString(arg)
		score = strings.Replace(strings.ToLower(score), "s:", "", 1)
		query += ", averageScore_greater: " + score
	}
	if genreParse.MatchString(arg) {
		genreArg := genreParse.FindString(arg)
		genreArg = strings.Replace(strings.ToLower(genreArg), "g:", "", 1)
		genreArg = strings.Replace(genreArg, "s:", "", 1)
		genreArg = strings.Replace(genreArg, "y:", "", 1)
		var divider string
		if strings.Contains(genreArg, ", ") {
			divider = ", "
		} else if strings.Contains(genreArg, ",") {
			divider = ","
		} else {
			divider = " "
		}
		genres := strings.Split(genreArg, divider)
		for i, s := range genres {
			if i == 0 {
				query += ", genre_in: ["
			}
			if s != "" {
				query += "\\\"" + s + "\\\","
			}
			if i == len(genres) - 1 {
				query += "]"
			}
		}
	}
	log.Println(query)
	AnimeMedia(s, m, query, 3)
}


func AnimeMedia(s *discordgo.Session, m *discordgo.MessageCreate, arg string, results int) {
	url := "https://graphql.anilist.co"
	method := "POST"

	payload := strings.NewReader("{\"query\":\"query { Page(page: 1, perPage: " + strconv.Itoa(results) + ") { media(" + arg + ") { id title { romaji english } type genres coverImage { large color } status season seasonYear episodes averageScore meanScore format description (asHtml: false) nextAiringEpisode { airingAt episode } characters(sort:ROLE, page: 1, perPage: 4) { edges { node { id name { full } } voiceActors(language: JAPANESE) { id name { full } } } } } } }\",\"variables\":{}}")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Anime not found!")
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Anime not found!")
		return
	}
	var response structs.AniData

	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Anime not found!")
		return
	}
	graphqlResponse := response.Data

	var index int
	if len(graphqlResponse.Page.Media) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Anime not found!")
		return
	}
	if results == 1 {
		index = 0
	} else {
		max := results
		if len(graphqlResponse.Page.Media) < results {
			max = len(graphqlResponse.Page.Media)
		}
		index = rand.Intn(max)
	}
	media := graphqlResponse.Page.Media[index]

	color := 0xFFFFFF
	if media.CoverImage.Color != "" {
		colorhexstring := strings.Replace(media.CoverImage.Color, "#", "", 1)
		colorvalue, _ := strconv.ParseInt(colorhexstring, 16, 64)
		color = int(colorvalue)
	}

	var title string
	var subtitle string
	if media.Title.English != "" {
		title = media.Title.English
		if media.Title.English != media.Title.Romaji {
			subtitle = "**" + media.Title.Romaji + "**\n"
		}
	} else {
		title = media.Title.Romaji
	}
	var genres string
	for i, s := range media.Genres {
		if i == 0 {
			genres += s
		} else {
			genres += ", " + s
		}
	}
	var airingTime string
	if media.Status == "RELEASING" {
		convtime := time.Unix(int64(media.NextAiringEpisode.AiringAt), 0)
		airingTime = "\n**Next Airing: **Episode " + strconv.Itoa(media.NextAiringEpisode.Episode) + " on " + convtime.Month().String() + " " + strconv.Itoa(convtime.Day()) + " " + strconv.Itoa(convtime.Year())
	}
	var episodes string
	if media.Format != "MOVIE" {
		if media.Episodes != 0 {
			episodes = "\n**Episodes:  **" + strconv.Itoa(media.Episodes)
		} else if media.Status != "RELEASING" {
			episodes = "\n**Not Yet Aired**"
		}
	}
	
	description := strings.Split(media.Description, "<br>")[0] + "\n\n" // only use everything before the first linebreak returned by description

	var characters string
	for _, s := range media.Characters.Edges {	
		characters += "[" + s.Node.Name.Full + "](http://anilist.co/character/" + strconv.Itoa(s.Node.ID) + ") "
		if len(s.VoiceActors) > 0 {
			characters += "[(" + s.VoiceActors[0].Name.Full + ")](http://anilist.co/staff/" + strconv.Itoa(s.VoiceActors[0].ID) + ")\n"
		}
	}
	re, err := regexp.Compile(`(?:<[\/a-z]*>)`)
	if err != nil {
		log.Println(err.Error())
	}
	description = re.ReplaceAllString(description, "")
	var format string
	switch media.Format {
	case "TV":
		format = "*TV Series*\n\n"
	case "TV_SHORT":
		format = "*TV Short*\n\n"
	case "MOVIE":
		format = "*Movie*\n\n"
	case "SPECIAL":
		format = "*Special*\n\n"
	case "MUSIC":
		format = "*Music*\n\n"
	default:
		format = "*" + media.Format + "*\n\n"
	}
	var season string
	if media.Season != "" {
		season = "**Season:  **" + strings.Title(strings.ToLower(media.Season)+" "+strconv.Itoa(media.SeasonYear))
	}
	var averageScore string
	if media.AverageScore != 0 {
		averageScore = "\n**Average Score:  **" + strconv.Itoa(media.AverageScore) + "%"
	} else {
		averageScore = "\n**Mean Score:  **" + strconv.Itoa(media.MeanScore) + "%"
	}
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       color,
		Description: subtitle + format + description + season + episodes + averageScore + airingTime,
		URL:         "https://anilist.co/anime/" + strconv.Itoa(media.ID),

		Image: &discordgo.MessageEmbedImage{
			URL: media.CoverImage.Large,
		},
		Title: title,
	}
	if genres != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "\n\nGenres",
			Value: genres,
		})
	}
	if characters != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "\nCharacters",
			Value: characters,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)	
}

func AniStaff(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	url := "https://graphql.anilist.co"
	method := "POST"

	payload := strings.NewReader("{\"query\":\"query { Staff(search: \\\"" + arg + "\\\", sort: SEARCH_MATCH ) { id gender age primaryOccupations dateOfBirth { year month day } name { full } image { large } characters(sort: FAVOURITES_DESC, page: 1, perPage: 3 ) { nodes { id name { full } media(sort: POPULARITY_DESC) { nodes { id title { romaji english } } } } } } }\",\"variables\":{}}")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Person not found!")
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Person not found!")
		return
	}
	var response structs.AniStaffData

	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Person not found!")
		return
	}
	graphqlResponse := response.Data

	var occupations string
	for i, s := range graphqlResponse.Staff.PrimaryOccupations {
		if i == 0 {
			occupations += s
		} else {
			occupations += ", " + s
		}
	}
	if occupations != "" {
		occupations = "*" + occupations + "*\n"
	}

	var birth string
	if graphqlResponse.Staff.DateOfBirth.Day != 0 {
		birth = "\n**Birth: **" + strconv.Itoa(graphqlResponse.Staff.DateOfBirth.Day) + " " + time.Month(graphqlResponse.Staff.DateOfBirth.Month).String()
		if graphqlResponse.Staff.DateOfBirth.Year != 0 {
			birth += " " + strconv.Itoa(graphqlResponse.Staff.DateOfBirth.Year)
		}
	}

	var age string
	if graphqlResponse.Staff.Age != 0 {
		age = "\n**Age: **" + strconv.Itoa(graphqlResponse.Staff.Age)
	}

	var gender string
	if graphqlResponse.Staff.Gender != "" {
		gender = "\n**Gender: **" + graphqlResponse.Staff.Gender
	}

	var roles string
	for _, s := range graphqlResponse.Staff.Characters.Nodes {
		roles += "[" + s.Name.Full + "](https://anilist.co/character/" + strconv.Itoa(s.ID) + ") "
		if s.Media.Nodes[0].Title.English != "" {
			roles += "[(" + s.Media.Nodes[0].Title.English
		} else {
			roles += "[(" + s.Media.Nodes[0].Title.Romaji
		}
		roles += ")](https://anilist.co/anime/" + strconv.Itoa(s.Media.Nodes[0].ID) + ")\n"
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xFFFFFF,
		Description: occupations + birth + age + gender,
		URL:         "https://anilist.co/staff/" + strconv.Itoa(graphqlResponse.Staff.ID),
		Image: &discordgo.MessageEmbedImage{
			URL: graphqlResponse.Staff.Image.Large,
		},
		Title: graphqlResponse.Staff.Name.Full,
	}
	if roles != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "\n\nCharacter Roles",
			Value: roles,
		})
	}


	s.ChannelMessageSendEmbed(m.ChannelID, embed)	
}

func AniChar(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	url := "https://graphql.anilist.co"
	method := "POST"

	payload := strings.NewReader("{\"query\":\" query { Character(search: \\\"" + arg + "\\\", sort: SEARCH_MATCH){ id gender age name { full } dateOfBirth { year month day } image { large } media(sort: POPULARITY_DESC, page: 1, perPage: 3){ nodes{ id title { english romaji } } edges { node { id } voiceActors (language: JAPANESE sort: FAVOURITES_DESC){ id name { full } } } } } }\",\"variables\":{}}")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Character not found!")
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Character not found!")
		return
	}
	var response structs.AniCharData

	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Character not found!")
		return
	}
	graphqlResponse := response.Data

	var birth string
	if graphqlResponse.Character.DateOfBirth.Day != 0 {
		birth = "\n\n**Birth: **" + strconv.Itoa(graphqlResponse.Character.DateOfBirth.Day) + " " + time.Month(graphqlResponse.Character.DateOfBirth.Month).String()
		if graphqlResponse.Character.DateOfBirth.Year != 0 {
			birth += " " + strconv.Itoa(graphqlResponse.Character.DateOfBirth.Year)
		}
	}

	var age string
	if graphqlResponse.Character.Age != "" {
		age = "\n**Age: **" + graphqlResponse.Character.Age
	}

	var gender string
	if graphqlResponse.Character.Gender != "" {
		gender = "\n**Gender: **" + graphqlResponse.Character.Gender
	}

	// currently, this command displays the most favorited japanese voice actresses in the most popular media appearance
	// afaik, there's no way to get all voice actresses across all appearances besides procedurally compiling a list

	var portrayal string
	if len(graphqlResponse.Character.Media.Edges) > 0 {
		for _, s := range graphqlResponse.Character.Media.Edges[0].VoiceActors {
			portrayal += "[" + s.Name.Full + "](https://anilist.co/staff/" + strconv.Itoa(graphqlResponse.Character.Media.Edges[0].VoiceActors[0].ID) + ")\n"
		} 
	}
	


	var appearances string
	var series string
	for i, s := range graphqlResponse.Character.Media.Nodes {
		if s.Title.English != "" {
			appearances += "[" + s.Title.English
		} else {
			appearances += "[" + s.Title.Romaji
		}
		if i == 0 {
			series = "*" + appearances + "](https://anilist.co/anime/" + strconv.Itoa(s.ID) + ")*"
		}
		appearances += "](https://anilist.co/anime/" + strconv.Itoa(s.ID) + ")\n"
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xFFFFFF,
		Description: series + birth + age + gender,
		URL:         "https://anilist.co/character/" + strconv.Itoa(graphqlResponse.Character.ID),
		Image: &discordgo.MessageEmbedImage{
			URL: graphqlResponse.Character.Image.Large,
		},
		Title: graphqlResponse.Character.Name.Full,
	}
	if appearances != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "\n\nAppearances",
			Value: appearances,
		})
	}
	if portrayal != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "\nPortrayed By",
			Value: portrayal,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)	
}