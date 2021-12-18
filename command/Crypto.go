package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Crypto(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	id := strings.ToLower(arg)
	droppedchars := regexp.MustCompile(`[^a-z0-9 _-]`)
	id = droppedchars.ReplaceAllString(id, "")
	id = strings.ReplaceAll(id, " ", "-")
	url := "https://api.coingecko.com/api/v3/simple/price?ids=" + id + "&vs_currencies=usd"
	method := "GET"
  
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err)
	  return
	}
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	defer res.Body.Close()
  
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	var result map[string]map[string]float64
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Invalid Name")
		return
	}
	price := result[id]["usd"]
	if price == 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid Name")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "$" + strconv.FormatFloat(price, 'f', -1, 64) + " (USD)")
}