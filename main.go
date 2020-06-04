package main

import (
	"fmt"
)

func main() {

	for _, rssConf := range AppConfig.RssFeeds {
		go monitorRssFeed(rssConf)
	}

	var input string
	_, _ = fmt.Scanln(&input)

}
