package main

import (
	"fmt"
	"strings"
)

var shitUsers []string = make([]string, 0)

type DisallowCommand struct {
	BotCommandBase
}

func (c *DisallowCommand) Satisfies(context *MessageContext) bool {
	if !strings.HasPrefix(context.Message.Content, "!disallow") {
		return false
	}

	author := context.Message.Author.ID
	guild, guildErr := context.Session.Guild(context.Message.GuildID)

	if guildErr != nil {
		return false
	}

	if (context.Message.Author.ID == "157636475117240320") {
		return true
	}

	for _, member := range guild.Members {
		if member.User.ID == author {
			for _, role := range member.Roles {
				if role == "209260592911876096" {
					return true
				}
			}
		}
	}

	return false
}

func (c *DisallowCommand) Exec(context *MessageContext) {
	userId := strings.Replace(context.Message.Content, "!disallow <@", "", -1)
	userId = strings.Replace(userId, ">", "", -1)

	if userId == "157636475117240320" {
		context.Session.ChannelMessageSend(
			context.Message.ChannelID,
			fmt.Sprintf("I am not going to disallow daddy, instead, " +
				"<@%s> has been disallowed from using the bot :^)", context.Message.Author.ID))
		return
	}

	guild, guildErr := context.Session.Guild(context.Message.GuildID)

	if guildErr != nil {
		return
	}

	if UserIsBlacklisted(userId) {
		if ToggleUserBlacklist(userId) {
			context.Session.ChannelMessageSend(context.Message.ChannelID,
				"user has been removed from the blacklist")
			return
		} else {
			context.Session.ChannelMessageSend(context.Message.ChannelID,
				"could not change user blacklist status")
		}
	}

	for _, member := range guild.Members {
		if member.User.ID == userId {
			context.Session.ChannelMessageSend(context.Message.ChannelID,
				fmt.Sprintf("<@%s> has been revoked the permission to use the bot", userId))
			ToggleUserBlacklist(userId)
			return
		}
	}

	context.Session.ChannelMessageSend(context.Message.ChannelID, fmt.Sprintf("user %s not found", userId))
}

func (c *DisallowCommand) CommandCategory() int {
	return COMMAND_CATEGORY_PRODUCTIVE
}