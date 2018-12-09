package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var userMessages map[string]int64 = make(map[string]int64)
const fifteen_seconds int64 = 15


type BotRegistry struct {
	ThreadHandler *ThreadHandler
	BotConfig     *BotConfig
	InfoMessages  []string
}

type MessageHandler struct {
	Receivers   []MessageReceiver
	BotRegistry *BotRegistry
}

type MessageContext struct {
	Session     *discordgo.Session
	Message     *discordgo.MessageCreate
	BotRegistry *BotRegistry
}

type MessageReceiver interface {
	Satisfies(*MessageContext) bool
	Exec(*MessageContext)
	Info() string
}

func CreateMessageHandler() *MessageHandler {
	h := MessageHandler{}
	h.Receivers = make([]MessageReceiver, 0)
	return &h
}

func (m *MessageHandler) RegisterReceiver(handler MessageReceiver) {
	m.Receivers = append(m.Receivers, handler)
	m.BotRegistry.InfoMessages = append(m.BotRegistry.InfoMessages, handler.Info())
}

func (m *MessageHandler) OnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Quick and dirty to see if it will work fine: rate limit bot to 1 message / 15 seconds in csg
	// except for commands in #mods
	// @todo make clean if it proves to be an effective counter measure to bot abuse
	if message.GuildID == "189466684938125312" && message.ChannelID != "247482564070080532" {
		canSend := false

		if timestamp, ok := userMessages[message.Author.ID]; ok {
			deadLine := time.Now().Unix() - fifteen_seconds
			if (deadLine > timestamp) {
				canSend = true;
				userMessages[message.Author.ID] = time.Now().Unix()
			}
		} else {
			canSend = true;
			userMessages[message.Author.ID] = time.Now().Unix()
		}

		if (!canSend) {
			return
		}
	}

	context := &MessageContext{Session: session, Message: message, BotRegistry: m.BotRegistry}
	for _, handler := range m.Receivers {
		if handler.Satisfies(context) {
			go handler.Exec(context)
		}
	}
}
