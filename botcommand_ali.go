package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"net/url"
	"strings"
)

type AliCommand struct {
}

func (a *AliCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!ali")
}

func (a *AliCommand) Exec(context *MessageContext) {
	searchUrl := "https://www.aliexpress.com/wholesale?catId=0&SearchText=%s"
	term := url.QueryEscape(strings.Replace(context.Message.Content, "!ali ", "", 1))

	searchUrl = fmt.Sprintf(searchUrl, term)

	resp, _ := http.Get(searchUrl)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Author:    &discordgo.MessageEmbedAuthor{},
		Color:     0x00ff00, // Green
		Fields:    make([]*discordgo.MessageEmbedField, 0),
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
		Title:     strings.Replace(context.Message.Content, "!ali ", "", 1),
		URL:       searchUrl,
	}

	doc.Find("a.product").Each(func(i int, current *goquery.Selection) {
		if len(embed.Fields) > 3 {
			return
		}

		link := strings.Trim(strings.Split(current.AttrOr("href", ""), ".html")[0]+".html",
			" \t\n")

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   current.Text(),
			Value:  "[Open on Ali Express](https:" + link + ")",
			Inline: false,
		})
	})

	if len(embed.Fields) == 0 {
		context.Session.ChannelMessageSend(context.Message.ChannelID,
			"No results :(")

		return
	}

	context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed)
}

func (a *AliCommand) Info() string {
	return `**!ali [searchterm]**
Search Ali Express for a specific search term`
}
