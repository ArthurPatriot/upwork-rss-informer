package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	AppConfig Config
)

type Config struct {
	BotToken     string          `json:"bot_token"`
	ChatID       string          `json:"chat_id"`
	IntervalFrom int             `json:"interval_from"`
	IntervalTo   int             `json:"interval_to"`
	Debug        bool            `json:"debug"`
	RssFeeds     []RssFeedConfig `json:"rss_feeds"`
}

type RssFeedConfig struct {
	Name    string `json:"name"`
	Link    string `json:"link"`
	Started bool   `json:"started"`
}

func init() {

	log.Println("Start Loading Config & App")

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Config File Not Loaded: ", err)
	}
	defer configFile.Close()

	configByte, _ := ioutil.ReadAll(configFile)

	if err := json.Unmarshal(configByte, &AppConfig); err != nil {
		log.Fatal("Config Parsing Error: ", err)
	}

	log.Println("App Config Loaded Successfully")

}
