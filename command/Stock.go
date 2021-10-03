package commands

import (
	structs "SyedBot/struct"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Stock(s *discordgo.Session, m *discordgo.MessageCreate, arg string) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/" + arg + "?region=US&lang=en-US&includePrePost=false&interval=2m&useYfid=true&range=1d&corsDomain=finance.yahoo.com&.tsrc=finance"
	method := "GET"

	payload := strings.NewReader("locate=tallinn&json=1&key=68600968168611176251x101318")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Referer", "https://finance.yahoo.com/quote/"+arg+"/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
	var stockResponse structs.StockData
	if err := json.Unmarshal(body, &stockResponse); err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Invalid Symbol")
		return
	}
	if len(stockResponse.Chart.Result) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid Symbol")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Market Price: "+fmt.Sprintf("%f", stockResponse.Chart.Result[0].Meta.RegularMarketPrice)+"\nPrevious Close: "+fmt.Sprintf("%f", stockResponse.Chart.Result[0].Meta.PreviousClose))
	//I'll make this a nice embed later
}
