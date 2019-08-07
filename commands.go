// commands.go
// this file contains all of the commands, but not underlying logic

package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, getHelpEmbed())
}

func DefaultHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s isn't a valid command. Use %shelp to learn more", strings.TrimPrefix(m.Content, configData.Prefix), configData.Prefix))
}

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func Avatar(s *discordgo.Session, m *discordgo.MessageCreate, author *discordgo.User) {
	s.ChannelMessageSendEmbed(m.ChannelID, getAvatarEmbed(author))
}
