package main

import (
	"bufio"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"jaytaylor.com/html2text"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PostCommand struct {
	BotCommandBase
}

func (h *PostCommand) Satisfies(context *MessageContext) bool {
	return UserIsAdmin(context.Message.Author.ID) &&
		strings.HasPrefix(context.Message.Content, "!post")
}

func (h *PostCommand) Exec(context *MessageContext) {
	postId, _ := strconv.Atoi(strings.Trim(strings.Replace(context.Message.Content, "!post ", "", -1), " "))
	post := context.BotRegistry.ThreadHandler.GetPost(postId)

	if post == nil {
		context.Session.ChannelMessageSend(context.Message.ChannelID,
			"this post is not in my registry")

		return
	}

	fileName := ""
	// download image
	if post.FileUrl != "" {
		fileName = strconv.Itoa(post.PostId) + ".jpg"

		resp, _ := http.Get(post.FileUrl)
		defer resp.Body.Close()

		fbytes, _ := ioutil.ReadAll(resp.Body)
		fpointer, _ := os.Create(fileName)
		writer := bufio.NewWriter(fpointer)
		writer.Write(fbytes)
		writer.Flush()
	}

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
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
		Title:     strconv.Itoa(post.PostId),
	}

	message := &discordgo.MessageSend{
		Embed: embed,
	}

	if fileName != "" {
		f, _ := os.Open(fileName)
		embed.Image = &discordgo.MessageEmbedImage{
			URL: fmt.Sprintf("attachment://%s", fileName),
		}
		message.Files = []*discordgo.File{
			&discordgo.File{
				Name:   fileName,
				Reader: f,
			},
		}
	}

	// post into #csg-reviews
	context.Session.ChannelMessageSendComplex("491172628917256192", message)

	if fileName != "" {
		os.Remove(fileName)
	}
}

func (h *PostCommand) Info() string {
	return ""
}

func (h *PostCommand) CommandCategory() int {
	return COMMAND_CATEGORY_PRODUCTIVE
}