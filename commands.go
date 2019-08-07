// commands.go
// this file contains all of the commands, but not underlying logic

package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//Help command
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, getHelpEmbed())
}

//DefaultHelp command
func DefaultHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s isn't a valid command. Use %shelp to learn more", strings.TrimPrefix(m.Content, configData.Prefix), configData.Prefix))
}

//Ping command
func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

//Avatar command
func Avatar(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Mentions) > 0 {
		if len(m.Mentions) > 4 {
			s.ChannelMessageSend(m.ChannelID, "Make sure to mention less than 5 users")
			return
		}
		for _, u := range m.Mentions {
			s.ChannelMessageSendEmbed(m.ChannelID, getAvatarEmbed(u))
		}
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, getAvatarEmbed(m.Author))
}

//Purge command
func Purge(s *discordgo.Session, m *discordgo.MessageCreate, n int) {
  var messageIDs [] string
	messages, err := s.ChannelMessages(m.ChannelID, n+1, "", "", "")
	if err != nil || n > 100 {
		fmt.Printf("Error getting messages: %s", err)
		s.ChannelMessageSend(m.ChannelID, "Please enter a number between 1 and 100")
		return
	}
	for _, element := range messages {
    messageIDs = append(messageIDs, element.ID)
	}
  s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
}
