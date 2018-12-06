package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type JokeCommand struct {
}

type Joke struct {
	Id     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func (h *JokeCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!joke")
}

func (h *JokeCommand) Exec(context *MessageContext) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	req.Header.Set("Accept", "application/json")
	res, _ := client.Do(req)
	bytes, _ := ioutil.ReadAll(res.Body)

	joke := Joke{}
	err := json.Unmarshal(bytes, &joke)
	if err != nil {
		context.Session.ChannelMessageSend(context.Message.ChannelID,
			fmt.Sprintf("Could not fetch joke: %s", err.Error()))
		return
	}

	context.Session.ChannelMessageSend(context.Message.ChannelID, joke.Joke)
}

func (h *JokeCommand) Info() string {
	return `**!joke**
Gets a random joke from icanhazdadjoke`
}
