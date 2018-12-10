package main

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

type CatCommand struct {
	BotCommandBase
}

func (c *CatCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!cat")
}

func (c *CatCommand) Exec(context *MessageContext) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00, // Green
		Fields: []*discordgo.MessageEmbedField{},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cataas.com/cat?" + strconv.FormatInt(time.Now().UnixNano(), 10),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
		Title:     "Random Cat",
	}

	context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed)
}

func (c *CatCommand) Info() string {
	return `**!cat**
Fetch a random cat image (everyone on the server will see a different cat image)`
}
