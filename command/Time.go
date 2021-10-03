package commands

import (
	"SyedBot/config"
	structs "SyedBot/struct"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)



func Time(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	zone, err := time.LoadLocation(arg)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Invalid timezone (currently this command only works with TimeZoneDB names (https://timezonedb.com/)")
		return
	}
	currenttime := time.Now().In(zone)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("It is %01d", currenttime.Hour())+":"+fmt.Sprintf("%02d", currenttime.Minute())+" "+zone.String())
	//today I learned you can use Sprintf to format stuff into strings without printing
	//rats
}

func TimeUntil (s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	var hours string
	var minutes string
	if strings.Contains(arg, ":") {
		if string(arg[1]) == ":" {
			hours += string(arg[0])
			minutes += string(arg[2]) + string(arg[3])
		} else if string(arg[2]) == ":" {
			hours += string(arg[0]) + string(arg[1])
			minutes += string(arg[3]) + string(arg[4])
		}
	} else {
		if len(arg) == 1 {
			hours += string(arg[0])
			minutes = "0"
		} else if len(arg) == 2 {
			hours += string(arg[0]) + string(arg[1])
			minutes = "0"
		} else if len(arg) == 3 {
			hours += string(arg[0])
			minutes += string(arg[1]) + string(arg[2])
		} else if len(arg) == 4 {
			hours += string(arg[0]) + string(arg[1])
			minutes += string(arg[2]) + string(arg[3])
		}
	}
	hoursint, err := strconv.Atoi(hours)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Invalid Format")
		return
	} else if hoursint > 23 {
		s.ChannelMessageSend(m.ChannelID, "Invalid Time")
		return
	}
	minutesint, err := strconv.Atoi(minutes)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Invalid Format")
		return
	} else if minutesint > 59 {
		s.ChannelMessageSend(m.ChannelID, "Invalid Time")
		return
	}
	currenttime := time.Now().UTC()
	target := time.Date(currenttime.Year(), currenttime.Month(), currenttime.Day(), hoursint, minutesint, 0, 0, currenttime.Location())
	timeuntil := int(time.Until(target).Minutes())
	if timeuntil < 0 {
		timeuntil += 1440
	}
	var output string
	if int(timeuntil/60) != 0 {
		output = strconv.Itoa(timeuntil/60) + " hours, " + strconv.Itoa(timeuntil%60) + " minutes"
	} else if int(timeuntil%60) != 0 {
		output = strconv.Itoa(timeuntil%60) + " minutes"
	} else {
		output = "Right freaking now"
	}
	s.ChannelMessageSend(m.ChannelID, output)
}

func TimeIn (s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	url := "https://geocode.xyz"
	method := "POST"

	payload_loc := strings.NewReader("locate=" + arg + "&json=1&Key=" + config.Config.Geocode)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload_loc)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "__cflb=0H28vTE11mXeuU6nLEGMumyL4X6iAPif9KvtBGzZSfF; geocode.xyz=686140097; xyzh=xyzh")

	geores, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer geores.Body.Close()

	geobody, err := ioutil.ReadAll(geores.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var geoResponse structs.GeoData

	if err := json.Unmarshal(geobody, &geoResponse); err != nil {
		log.Println(err.Error())
		return

	}

	if geoResponse.Longt == "0.00000" {
		s.ChannelMessageSend(m.ChannelID, "Invalid City")
		return
	}

	url = "http://api.timezonedb.com/v2.1/get-time-zone?key=" + config.Config.TimeZoneDB + "&by=position&format=json&lat=" + geoResponse.Latt + "&lng=" + geoResponse.Longt
	method = "GET"

	payload_time := strings.NewReader("")
	req, err = http.NewRequest(method, url, payload_time)

	if err != nil {
		fmt.Println(err)
		return
	}
	timeres, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer timeres.Body.Close()

	timebody, err := ioutil.ReadAll(timeres.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var timeResponse structs.GeoTimeData
	if err := json.Unmarshal(timebody, &timeResponse); err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Invalid City")
		return
	}
	currenttime := time.Now().UTC()
	convtime := currenttime.Add(1000000000 * time.Duration(timeResponse.GmtOffset)) // rofl
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("It is %01d", convtime.Hour())+":"+fmt.Sprintf("%02d", convtime.Minute())+" in "+geoResponse.Standard.City+", "+geoResponse.Standard.Countryname)
}