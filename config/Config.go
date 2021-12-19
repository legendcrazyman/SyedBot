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
	Prefix 		 string
	MongoDB		 MongoDB
}

type Twitter struct {
	Token       string
	TokenSecret string
	Key         string
	KeySecret   string
}

type MongoDB struct {
	Hostname	string
	Username	string
	Password	string
}

var Config ConfigData

func ConfigInit() {
	readin, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatalln("Error: " + err.Error())
	} //put some error handling here later
	_ = json.Unmarshal(readin, &Config)
}
