package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ConfigData struct {
	DiscordToken string
	Twitter      Twitter
	Geocode      string
	TimeZoneDB   string
}

type Twitter struct {
	Token       string
	TokenSecret string
	Key         string
	KeySecret   string
}

var Config ConfigData

func ConfigInit() {
	readin, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatalln("Error: " + err.Error())
	} //put some error handling here later
	_ = json.Unmarshal(readin, &Config)
}
