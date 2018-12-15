package main

import (
	"fmt"
	"strconv"
	"strings"
)

var sudokuVictimCounter int = 0

type SelfbanCommand struct {
	BotCommandBase
}

func (s *SelfbanCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!selfban") ||
		strings.HasPrefix(context.Message.Content, "!sudoku")
}

func (s *SelfbanCommand) Exec(context *MessageContext) {
	sudokuVictimCounter ++;

	context.Session.GuildMemberDeleteWithReason(context.Message.GuildID, context.Message.Author.ID,
		fmt.Sprintf("Sudoku victim %s", strconv.Itoa(sudokuVictimCounter)))

	context.Session.ChannelMessageSend(context.Message.ChannelID,
		fmt.Sprintf("<@%s> has been banned from the server", context.Message.Author.ID))
}

func (s *SelfbanCommand) Info() string {
	return `**!selfban** [Alias: !sudoku]
starts a nice game of sudoku`
}
