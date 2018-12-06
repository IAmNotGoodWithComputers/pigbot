package main

import (
	"strings"
)

type AboutCommand struct {
}

func (c *AboutCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!about") ||
		strings.HasPrefix(context.Message.Content, "!help")
}

func (c *AboutCommand) Exec(context *MessageContext) {
	message := ""

	for _, info := range context.BotRegistry.InfoMessages {
		if info == "" {
			continue
		}
		message += info + "\n\n"
	}

	context.Session.ChannelMessageSend(context.Message.ChannelID, message)
}

func (c *AboutCommand) Info() string {
	return ""
}
