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
}

func (c *CatCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!cat")
}

func (c *CatCommand) Exec(context *MessageContext) {
	rnd := strconv.FormatInt(time.Now().UnixNano(), 10)
	fname := rnd + ".jpg"
	url := "https://cataas.com/cat?" + rnd
	img := fetchCatImg(fname, url)

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
	//context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed)
}

func fetchCatImg(file string, url string) string {
	out, _ := os.Create("cats/" + file)
	defer out.Close()

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	_, err := io.Copy(out, resp.Body)
	if err == nil {
		println(err)
	}
	//cat, _ := ioutil.ReadFile("cats/" + file)
	return "cats/" + file
}

func (c *CatCommand) Info() string {
	return `**!cat**
Fetch a random cat image (everyone on the server will see a different cat image)`
}
