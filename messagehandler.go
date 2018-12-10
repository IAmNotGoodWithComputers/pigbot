package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var CSG_FUNCOMMAND_THROTTLE int64 = 30
const CSG_SERVERID string = "189466684938125312"

const COMMAND_CATEGORY_FUN = 1
const COMMAND_CATEGORY_PRODUCTIVE = 2

var csgLastFunCommand int64 = 0

type BotCommandBase struct {}

func (b *BotCommandBase) CommandCategory() int {
	return COMMAND_CATEGORY_FUN
}

func (b *BotCommandBase) Info() string {
	return ""
}

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
	CommandCategory() int
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
	context := &MessageContext{Session: session, Message: message, BotRegistry: m.BotRegistry}
	for _, handler := range m.Receivers {
		if handler.Satisfies(context) {
			if (handler.CommandCategory() == COMMAND_CATEGORY_FUN && message.GuildID == CSG_SERVERID) {
				deadLine := time.Now().Unix() - CSG_FUNCOMMAND_THROTTLE
				if (deadLine <= csgLastFunCommand) {
					return
				}

				csgLastFunCommand = time.Now().Unix()
			}

			go handler.Exec(context)
		}
	}
}

