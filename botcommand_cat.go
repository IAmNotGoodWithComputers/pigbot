package main

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
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
	uplImg(img)
	println(img)

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00, // Green
		Fields: []*discordgo.MessageEmbedField{},
		Image: &discordgo.MessageEmbedImage{
			URL: url,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
		Title:     "Random Cat",
	}

	context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed)
}

func uplImg(img []byte) {
	resp, _ := http.Post("https://api.imgur.com/3/upload", "text/plain; charset=UTF-8", bytes.NewBuffer(img))
	println(resp)
}

func fetchCatImg(file string, url string) []byte {
	out, _ := os.Create("cats/" + file)
	defer out.Close()

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	_, err := io.Copy(out, resp.Body)
	if err == nil {
		println(err)
	}
	cat, _ := ioutil.ReadFile("cats/" + file)
	return cat
}

func (c *CatCommand) Info() string {
	return `**!cat**
Fetch a random cat image (everyone on the server will see a different cat image)`
}
