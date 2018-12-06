package main

import (
	"fmt"
	"strings"
)

type SelfbanCommand struct {
}

func (s *SelfbanCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!selfban") ||
		strings.HasPrefix(context.Message.Content, "!sudoku")
}

func (s *SelfbanCommand) Exec(context *MessageContext) {
	context.Session.GuildBanCreate(context.Message.GuildID, context.Message.Author.ID, 1)

	context.Session.ChannelMessageSend(context.Message.ChannelID,
		fmt.Sprintf("<@%s> has been banned from the server", context.Message.Author.ID))
}

func (s *SelfbanCommand) Info() string {
	return `**!selfban** [Alias: !sudoku]
starts a nice game of sudoku`
}
