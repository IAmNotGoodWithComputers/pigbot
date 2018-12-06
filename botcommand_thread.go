package main

import (
	"strings"
)

type ThreadCommand struct {
}

func (c *ThreadCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!thread")
}

func (c *ThreadCommand) Exec(context *MessageContext) {
	if !context.BotRegistry.ThreadHandler.ThreadExists {
		context.Session.ChannelMessage(context.Message.ChannelID,
			"there seems to be no /csg/ thread up right now")
		return
	}

	context.Session.ChannelMessage(context.Message.ChannelID,
		context.BotRegistry.ThreadHandler.CurrentThread)
}
