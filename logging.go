package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func logCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := fmt.Sprintf("%s ran: %s in <#%s>", m.Author.Mention(), m.Content, m.ChannelID)
	postToLogs(s, m.GuildID, message)
}

func logServerInvite(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := fmt.Sprintf("Deleted server invite in message: `%s` by %s", m.ContentWithMentionsReplaced(), m.Author.Mention())
	postToLogs(s, m.GuildID, message)
}

func postToLogs(s *discordgo.Session, guildID string, message string) {
	for _, channelID := range loadedConfigData.ChannelsLogging {
		channel, err := s.Channel(channelID)
		if err != nil {
			fmt.Printf("Error confirming channel guild id for %s: %s", channelID, err)
		} else if channel.GuildID == guildID {
			_, err := s.ChannelMessageSend(channelID, message)
			if err != nil {
				fmt.Printf("Error logging message %s in channel %s", message, channelID)
			}
		}
	}
}
