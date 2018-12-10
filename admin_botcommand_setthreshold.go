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
		context.Message.ChannelID == "247482564070080532"
}

func (c *SetThresholdCommand) Exec(context *MessageContext) {
	threshold := strings.Replace(context.Message.Content, "!cooldown ", "", -1)
	threshold = strings.Trim(threshold, " \t\n!")

	intVal, err := strconv.Atoi(threshold)

	if err != nil {
		context.Session.ChannelMessageSend(context.Message.ChannelID, err.Error())
		return
	}

	CSG_FUNCOMMAND_THROTTLE = int64(intVal)

	context.Session.ChannelMessageSend(context.Message.ChannelID, fmt.Sprintf("fun command cooldown set to %s seconds",
		threshold))
}

func (c *SetThresholdCommand) CommandCategory() int {
	return COMMAND_CATEGORY_PRODUCTIVE
}