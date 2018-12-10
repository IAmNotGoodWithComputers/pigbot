package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AdviceCommand struct {
	BotCommandBase
}

type Advice struct {
	Slip struct{
		Advice string `json:"advice"`
		SlipId string `json:"slip_id"`
	} `json:"slip"`
}

func (a *AdviceCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!advice")
}

func (a *AdviceCommand) Exec(context *MessageContext) {
	resp, _ := http.Get("https://api.adviceslip.com/advice")
	advice := Advice{}
	bytes, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(bytes, &advice)

	if err != nil {
		fmt.Println(err.Error())
	}
	context.Session.ChannelMessageSend(context.Message.ChannelID, advice.Slip.Advice)
}

func (a *AdviceCommand) Info() string {
	return `**!advice**
gives you some solid, high quality life advice`
}