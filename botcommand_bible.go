package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"jaytaylor.com/html2text"
	"net/http"
	"strconv"
	"strings"
)

type BibleCommand struct {
	BotCommandBase
}

func (b *BibleCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!bible") &&
		context.Message.GuildID != CSG_SERVERID
}

func (b *BibleCommand) Exec(context *MessageContext) {
	msgLen := len(strings.Split(context.Message.Content, " "))
	if msgLen > 2 {
		b.getBibleVerse(context)
	} else if msgLen < 2 {
		b.randomBibleVerse(context)
	}
}

func fmtBook(msg []string) string {
	if len(msg) == 3 {
		return string(msg[1])
	} else if len(msg) == 4 {
		return fmt.Sprintf("%s-%s", msg[1], msg[2])
	}
	return ""
}

func fmtChpVrs(msg []string) []string {
	if len(msg) == 3 {
		return strings.Split(string(msg[2]), ":")
	} else if len(msg) == 4 {
		return strings.Split(string(msg[3]), ":")
	}
	return []string{}
}

func (b *BibleCommand) getBibleVerse(context *MessageContext) {
	msg := strings.Split(context.Message.Content, " ")
	msgBook := fmtBook(msg)
	if msgBook != "" {
		chpVrs := fmtChpVrs(msg)
		if len(chpVrs) == 2 {
			msgChpt := chpVrs[0]
			msgVrs := chpVrs[1]
			chpt, _ := strconv.ParseInt(msgChpt, 10, 32)
			vrs, _ := strconv.ParseInt(msgVrs, 10, 32)

			if chpt > 0 && vrs > 0 {
				url := fmt.Sprintf("https://dailyverses.net/%s/%d/%d", msgBook, chpt, vrs)
				resp, _ := http.Get(url)
				if resp.StatusCode == 200 {
					doc, _ := goquery.NewDocumentFromReader(resp.Body)
					htmlVerse, _ := doc.Find("div.bibleVerse").First().Html()
					txt, _ := html2text.FromString(htmlVerse)
					verse := strings.Split(txt, "(")[0]
					context.Session.ChannelMessageSend(context.Message.ChannelID, verse)
				}
			}
		}
	}
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
