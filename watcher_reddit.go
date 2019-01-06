package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
	"strings"
	"time"
)

var postedUrls []string = make([]string, 0)

var whiteList []string = []string{"mechanical", "keyboard", "mouse"}
var blackList []string = []string{"laptop", "prebuilt", "corsair", "steelseries", "logitech", "evga", "hyperx", "razer", "mousepad", "ssd"}

type RedditWatcher struct {

}

func (r *RedditWatcher) StartRedditWatcher(session *discordgo.Session) {
	go r.Watch(session)
}

func (r *RedditWatcher) Watch(session *discordgo.Session) {
	r.CheckReddit(session)
	time.Sleep(5 * time.Minute)
	r.Watch(session)
}

func (r *RedditWatcher) CheckReddit(session *discordgo.Session){
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.reddit.com/r/buildapcsales/.rss")

	for _, item := range feed.Items {
		title := item.Title
		url := item.Link

		if containsAnyWhitelist(title) &&
			!containsAnyBlacklist(title) &&
			!hasBeenPosted(url) {
			postedUrls = append(postedUrls, url)

			session.ChannelMessageSend("287939846662651906", fmt.Sprintf("**%s**\n<%s>", title, url))

			time.Sleep(5 * time.Second)
		}
	}
}

func containsAnyWhitelist(inp string) bool {
	inp = strings.ToLower(inp)
	for _, term := range whiteList {
		if strings.Contains(inp, term) {
			return true
		}
	}

	return false
}

func containsAnyBlacklist(inp string) bool {
	inp = strings.ToLower(inp)
	for _, term := range blackList {
		if strings.Contains(inp, term) {
			return true
		}
	}

	return false
}

func hasBeenPosted(inp string) bool {
	inp = strings.ToLower(inp)
	for _, url := range postedUrls {
		if url == inp {
			return true
		}
	}

	return false
}