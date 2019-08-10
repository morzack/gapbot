package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func logCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := fmt.Sprintf("%s ran: `%s` in <#%s>", m.Author.Mention(), m.Content, m.ChannelID)
	postToLogs(s, message)
}

func postToLogs(s *discordgo.Session, message string) {
	for _, channelID := range configData.ChannelsLogging {
		_, err := s.ChannelMessageSend(channelID, message)
		if err != nil {
			fmt.Printf("Error logging message %s in channel %s", message, channelID)
		}
	}
}
