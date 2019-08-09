// commands.go
// this file contains all of the commands, but not underlying logic

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// process commands as normal user
func UserCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	switch command {
	case "help":
		Help(s, m)
	case "ping":
		Ping(s, m)
	case "avatar":
		Avatar(s, m)
	case "user":
		UserInfo(s, m)
	default:
		DefaultHelp(s, m)
	}
}

// process commands as admin
func AdminCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	switch command {
	case "purge":
		Purge(s, m)
	case "help":
		AdminHelp(s, m)
	case "kick":
		KickUser(s, m)
	case "ban":
		BanUser(s, m)
	default:
		UserCommand(s, m, command)
	}
}

//Help command
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, getHelpEmbed())
}

func AdminHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, getAdminHelpEmbed())
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
func Purge(s *discordgo.Session, m *discordgo.MessageCreate) {
	fields := strings.Fields(m.Content)
	n, err := strconv.Atoi(fields[len(fields)-1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Make sure a number of messages to delete is specified at the end of the command")
		return
	}
	if n > 99 || n < 1 {
		s.ChannelMessageSend(m.ChannelID, "Please enter a number between 1 and 99 (inclusive)")
		return
	}
	var messageIDs []string
	messages, err := s.ChannelMessages(m.ChannelID, n+1, "", "", "")
	if err != nil {
		fmt.Printf("Error getting messages: %s", err)
		return
	}
	for _, element := range messages {
		messageIDs = append(messageIDs, element.ID)
	}
	s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
}

//UserInfo embed command
func UserInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	g, _ := s.Guild(m.GuildID)
	if len(m.Mentions) > 0 {
		if len(m.Mentions) > 4 {
			s.ChannelMessageSend(m.ChannelID, "Make sure to mention less than 5 users")
			return
		}
		for _, u := range m.Mentions {
			s.ChannelMessageSendEmbed(m.ChannelID, getUserEmbed(u, s, g))
		}
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, getUserEmbed(m.Author, s, g))
}

// RemoveUser -- handle the kick/ban command based on status of param ban
func RemoveUser(s *discordgo.Session, m *discordgo.MessageCreate, ban bool) {
	method := "kick"
	if ban {
		method = "ban"
	}

	fields := strings.SplitN(strings.TrimPrefix(m.Content, configData.Prefix), " ", 3)
	if len(fields) < 3 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please make sure to specify a user and reason when %sing", method))
		return
	}
	reason := fields[2]
	if len(m.Mentions) < 1 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Make sure to specify a user to %s", method))
		return
	}
	user := m.Mentions[0]

	var err error
	if ban {
		err = s.GuildBanCreateWithReason(m.GuildID, user.ID, reason, 1)
	} else {
		err = s.GuildMemberDeleteWithReason(m.GuildID, user.ID, reason)
	}
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to %s %s for reason %s", method, user.Mention(), reason))
		fmt.Printf("Error when %sing user %s, %s", method, user.Mention(), err)
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s was %sed because of reason: %s", user.Mention(), method, reason))
}

func KickUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	RemoveUser(s, m, false)
}

func BanUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	RemoveUser(s, m, true)
}
