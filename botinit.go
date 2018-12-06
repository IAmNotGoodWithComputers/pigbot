package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func BotInit(discord *discordgo.Session, config *BotConfig) {
	// prepare
	threadHandler := CreateThreadHandler()
	threadHandler.StartWatcher()

	registry := &BotRegistry{ThreadHandler: threadHandler, BotConfig: config}
	messageHandler := CreateMessageHandler()
	messageHandler.BotRegistry = registry

	messageHandler.RegisterReceiver(new(PostCommand))
	messageHandler.RegisterReceiver(new(JokeCommand))
	messageHandler.RegisterReceiver(new(CatCommand))
	messageHandler.RegisterReceiver(new(AnswerCommand))
	messageHandler.RegisterReceiver(new(ThreadCommand))
	messageHandler.RegisterReceiver(new(SelfbanCommand))

	discord.AddHandler(messageHandler.OnMessage)
	discord.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate) {
		guild, _ := session.State.Guild(message.GuildID)
		channel, _ := session.State.Channel(message.ChannelID)
		go fmt.Println(fmt.Sprintf("[%s]-#%s (%s): %s", guild.Name, channel.Name, message.Author.Username, message.Content))
	})

	discord.Open()
}
