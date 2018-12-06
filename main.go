package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
)

type BotConfig struct {
	ApiKey          string
	WolframAlphaKey string
}

func main() {
	botConfig := parseBotConfig()

	auth := "Bot " + botConfig.ApiKey

	discord, err := discordgo.New(auth)
	dieOnError("Could not create discord session", err)

	user, err := discord.User("@me")
	dieOnError("Bot has no account", err)

	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "June 4th Tiananmen Square Massacre")
	})

	fmt.Printf("Started bot %s: (%s) \n", user.Username, user.ID)

	BotInit(discord, botConfig)

	// prevent process from exiting
	<-make(chan struct{})
}

func parseBotConfig() *BotConfig {
	botConfig := BotConfig{}

	fcontent, err := ioutil.ReadFile("config.json")
	dieOnError("Could not read config file", err)

	err = json.Unmarshal(fcontent, &botConfig)
	dieOnError("Could not parse config file", err)

	return &botConfig
}

func dieOnError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
