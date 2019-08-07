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
	s.ChannelMessageSendEmbed(m.ChannelID, getAvatarEmbed(m.Author))
}

//Purge command
func Purge(s *discordgo.Session, m *discordgo.MessageCreate, n int) {
	messages, err := s.ChannelMessages(m.ChannelID, n, m.ID, "", "")
	if err != nil {
		fmt.Printf("Error getting messages: %s", err)
		s.ChannelMessageSend(m.ChannelID, "You can't delete that many messages!")
		return
	}
	for _, element := range messages {
		s.ChannelMessageDelete(m.ChannelID, element.ID)
	}
}
