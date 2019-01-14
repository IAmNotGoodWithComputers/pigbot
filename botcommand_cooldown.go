package main

import (
	"fmt"
	"strconv"
	"strings"
)

type SetThresholdCommand struct {
	BotCommandBase
}

func (c *SetThresholdCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!cooldown") &&
		UserIsAdmin(context.Message.Author.ID)
}

func (c *SetThresholdCommand) Exec(context *MessageContext) {
	parts := strings.Split(context.Message.Content, " ")
	if len(parts) != 3 {
		context.Session.ChannelMessageSend(context.Message.ChannelID, "invalid parameters. Usage: !cooldown [!command] [seconds]")
	}

	threshold := strings.Trim(parts[2], " \t\n!")
	command := strings.Trim(parts[1], " \t\n!")
	intVal, err := strconv.Atoi(threshold)

	if err != nil {
		context.Session.ChannelMessageSend(context.Message.ChannelID, err.Error())
		return
	}

	SetCommandPolicy(command, context.Message.GuildID, intVal)

	context.Session.ChannelMessageSend(context.Message.ChannelID,
		fmt.Sprintf("command cooldown for %s set to %s seconds", command, threshold))
}

func (c *SetThresholdCommand) CommandCategory() int {
	return COMMAND_CATEGORY_PRODUCTIVE
}