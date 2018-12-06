package main

import (
	"github.com/PuerkitoBio/goquery"
	"jaytaylor.com/html2text"
	"net/http"
	"strings"
)

type BibleCommand struct {

}

func (b *BibleCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!bible")
}

func (b *BibleCommand) Exec(context *MessageContext) {
	resp, _ :=  http.Get("https://dailyverses.net/random-bible-verse")
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	htmlVerse, _ := doc.Find("div.bibleVerse").First().Html()

	txt , _ := html2text.FromString(htmlVerse)

	verse := strings.Split(txt, "(")[0]

	context.Session.ChannelMessageSend(context.Message.ChannelID, verse)
}

func (b *BibleCommand) Info() string {
	return `**!bible** 
returns a random bible verse`
}
