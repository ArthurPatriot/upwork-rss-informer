package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	rs = NewRssSended()
)

type (
	RssFeed struct {
		Channel RssChannel `xml:"channel"`
	}

	RssChannel struct {
		Items []RssItem `xml:"item"`
	}

	RssItem struct {
		Title   string `xml:"title"`
		Link    string `xml:"link"`
		PubDate string `xml:"pubDate"`
		Guid    string `xml:"guid"`
	}

	RssSended struct {
		sync.Mutex
		rsp map[string]bool
	}
)

func NewRssSended() *RssSended {
	return &RssSended{
		rsp: make(map[string]bool),
	}
}

func (rs *RssSended) Sended(item *RssItem) {
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()

	rs.rsp[item.Guid] = true
}

func (rs *RssSended) IsSended(item *RssItem) bool {

	if rs.rsp[item.Guid] {
		return true
	}

	return false

}

func monitorRssFeed(rssConf RssFeedConfig) {

	log.Println("Start Monitoring...", rssConf.Name)

	ticker := time.NewTicker(time.Duration(rand.Intn(AppConfig.IntervalTo-AppConfig.IntervalFrom)+AppConfig.IntervalFrom) * time.Second)
	stopTickers := make(chan bool)

	for {
		select {
		case <-stopTickers:
			log.Println("Stop Monitor Feed...", rssConf.Name)
			return
		case <-ticker.C:
			if AppConfig.Debug {
				log.Println("Handle Feed...", rssConf.Name)
			}
			handleRssFeed(&rssConf)
		}
	}
}

func handleRssFeed(rssConf *RssFeedConfig) {

	f := RssFeed{}

	fd, err := getRssFeed(rssConf.Link)
	if err != nil {
		return
	}

	err = xml.Unmarshal(fd, &f)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range f.Channel.Items {

		if AppConfig.Debug {
			log.Println("Checking Exists for...", item.Title)
		}

		if !rssConf.Started {
			rs.Sended(&item)
			continue
		}

		if rs.IsSended(&item) {
			continue
		}

		if AppConfig.Debug {
			log.Println("Need Send...", item.Title)
		}

		go sendRssFeed(item, rssConf)

	}

	if !rssConf.Started {
		rssConf.Started = true
	}

}

func getRssFeed(rssLink string) ([]byte, error) {

	r, err := http.Get(rssLink)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return body, nil

}
