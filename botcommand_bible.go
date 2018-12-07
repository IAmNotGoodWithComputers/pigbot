package main

import (
	"fmt"
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
	msgLen := len(strings.Split(context.Message.Content, " "))
	if msgLen > 1 {
		b.getBibleVerse(context)
	} else {
		b.randomBibleVerse(context)
	}
}

func (b *BibleCommand) getBibleVerse(context *MessageContext) {
	msgBook := strings.Split(context.Message.Content, " ")[1]
	msgChpt := strings.Split(strings.Split(context.Message.Content, " ")[2], ":")[0]
	msgVrs := strings.Split(strings.Split(context.Message.Content, " ")[2], ":")[1]

	url := fmt.Sprintf("https://dailyverses.net/%s/%s/%s", msgBook, msgChpt, msgVrs)
	resp, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	htmlVerse, _ := doc.Find("div.bibleVerse").First().Html()
	txt, _ := html2text.FromString(htmlVerse)
	verse := strings.Split(txt, "(")[0]
	context.Session.ChannelMessageSend(context.Message.ChannelID, verse)
}

func (b *BibleCommand) randomBibleVerse(context *MessageContext) {
	resp, _ := http.Get("https://dailyverses.net/random-bible-verse")
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	htmlVerse, _ := doc.Find("div.bibleVerse").First().Html()
	txt, _ := html2text.FromString(htmlVerse)
	verse := strings.Split(txt, "(")[0]
	context.Session.ChannelMessageSend(context.Message.ChannelID, verse)
}

func (b *BibleCommand) Info() string {
	return `**!bible** 
returns a random bible verse 
**!bible _book_ _chapter:verse_**
returns a specific bible verse
refer to <https://dailyverses.net/bible-books> to see all available books and verses`
}
