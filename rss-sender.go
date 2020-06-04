package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type (
	TgMessage struct {
		ChatId    string `json:"chat_id"`
		ParseMode string `json:"parse_mode"`
		Text      string `json:"text"`
	}
)

func sendRssFeed(item RssItem, rssConf *RssFeedConfig) {

	tgmB := new(bytes.Buffer)
	tgmTB := strings.Builder{}
	tgm := TgMessage{
		ChatId:    AppConfig.ChatID,
		ParseMode: "html",
	}

	tgmTB.WriteString(item.Title + "\n" + item.PubDate)

	//TODO: Future tool
	//if rssConf.Detailed {
	//	ia, ok := item.GetDetails(item.Link)
	//	if ok && len(ia) > 0 {
	//
	//		tgmTB.WriteString("\n")
	//
	//		for dk, dv := range ia {
	//			tgmTB.WriteString("\n" + strings.Title(dk) + ": " + dv)
	//		}
	//
	//	}
	//}

	tgmTB.WriteString("\n\n" + item.Link + "\n\n#" + rssConf.Name)

	tgm.Text = tgmTB.String()

	_ = json.NewEncoder(tgmB).Encode(tgm)

	if AppConfig.Debug {
		log.Println("Send ->", item.Title)
	}

	rs.Sended(&item)

	_, _ = http.Post("https://api.telegram.org/bot"+AppConfig.BotToken+"/sendMessage", "application/json; charset=utf-8", tgmB)

}
