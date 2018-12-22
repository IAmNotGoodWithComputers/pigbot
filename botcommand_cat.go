package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"os"
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
	img, fname := fetchCatImg()

	embed := &discordgo.MessageEmbed{
		Author:    &discordgo.MessageEmbedAuthor{},
		Color:     0x00ff00, // Green
		Fields:    []*discordgo.MessageEmbedField{},
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
		Title:     "Random Cat",
	}

	message := &discordgo.MessageSend{
		Embed: embed,
	}

	f, _ := os.Open(img)
	embed.Image = &discordgo.MessageEmbedImage{
		URL: fmt.Sprintf("attachment://%s", fname),
	}
	message.Files = []*discordgo.File{
		&discordgo.File{
			Name:   fname,
			Reader: f,
		},
	}

	context.Session.ChannelMessageSendComplex(context.Message.ChannelID, message)

	f.Close()
	os.Remove(img)
}

func fetchCatImg() (string, string) {
	rnd := strconv.FormatInt(time.Now().UnixNano(), 10)
	file := rnd + ".jpg"
	url := "https://cataas.com/cat?" + rnd

	out, _ := os.Create("cats/" + file)
	defer out.Close()

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	io.Copy(out, resp.Body)
	return "cats/" + file, file
}

func (c *CatCommand) Info() string {
	return `**!cat**
Fetch a random cat image (everyone on the server will see the same cat image)`
}
