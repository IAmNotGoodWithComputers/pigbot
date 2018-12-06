package main

import (
	"github.com/bwmarrin/discordgo"
	"jaytaylor.com/html2text"
	"strconv"
	"strings"
)

type PostCommand struct {
}

func (h *PostCommand) Satisfies(context *MessageContext) bool {
	return context.Message.ChannelID == "247482564070080532" && // only allow messages from moderator channel
		strings.HasPrefix(context.Message.Content, "!post")
}

func (h *PostCommand) Exec(context *MessageContext) {
	postId, _ := strconv.Atoi(strings.Trim(strings.Replace(context.Message.Content, "!post ", "", -1), " "))
	post := context.BotRegistry.ThreadHandler.GetPost(postId)

	postText, _ := html2text.FromString(post.Content)

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00, // Green
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Post",
				Value:  postText,
				Inline: false,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: post.FileUrl,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
		Title:     strconv.Itoa(post.PostId),
	}

	// post into #csg-reviews
	context.Session.ChannelMessageSendEmbed("491172628917256192", embed)
}
