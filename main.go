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

	botInit(discord, botConfig)

	// prevent process from exiting
	<-make(chan struct{})
}

func botInit(discord *discordgo.Session, config *BotConfig) {
	// prepare
	threadHandler := CreateThreadHandler()
	threadHandler.StartWatcher()

	registry := &BotRegistry{ThreadHandler: threadHandler, BotConfig: config}
	messageHandler := CreateMessageHandler()
	messageHandler.BotRegistry = registry

	messageHandler.RegisterReceiver(new(AboutCommand))

	messageHandler.RegisterReceiver(new(PostCommand))
	messageHandler.RegisterReceiver(new(JokeCommand))
	messageHandler.RegisterReceiver(new(CatCommand))
	messageHandler.RegisterReceiver(new(AnswerCommand))
	messageHandler.RegisterReceiver(new(ThreadCommand))
	messageHandler.RegisterReceiver(new(SelfbanCommand))
	messageHandler.RegisterReceiver(new(AliCommand))

	discord.AddHandler(messageHandler.OnMessage)
	discord.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate) {
		guild, _ := session.State.Guild(message.GuildID)
		channel, _ := session.State.Channel(message.ChannelID)
		go fmt.Println(fmt.Sprintf("[%s]-#%s (%s): %s", guild.Name, channel.Name, message.Author.Username, message.Content))
	})

	discord.Open()
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
