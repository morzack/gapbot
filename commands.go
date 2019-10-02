package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// process commands that can only be run in dms
func DMCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	switch command {
	case "help":
		DMHelp(s, m)
	case "ping":
		Ping(s, m)
	case "register":
		Register(s, m)
	default:
		DefaultHelp(s, m)
	}
}

// process commands as normal user
func UserCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	switch command {
	case "help":
		ServerHelp(s, m)
	case "user":
		UserInfo(s, m)
	case "server":
		ServerInfo(s, m)
	default:
		DMCommand(s, m, command)
	}
}

// process commands as admin
func AdminCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	switch command {
	case "purge":
		Purge(s, m)
		logCommand(s, m)
	case "addlog":
		AddLoggingChannelCommand(s, m)
		logCommand(s, m)
	case "removelog":
		RemoveLoggingChannelCommand(s, m)
		logCommand(s, m)
	case "kick":
		KickUser(s, m)
		logCommand(s, m)
	case "ban":
		BanUser(s, m)
		logCommand(s, m)
	case "help":
		AdminHelp(s, m)
	case "register":
		TempMassRegister(s, m)
	default:
		UserCommand(s, m, command)
	}
}

func DMHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, getDMHelpEmbed())
}

func ServerHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, getServerHelpEmbed())
}

func AdminHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, getAdminHelpEmbed())
}

func DefaultHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s isn't a valid command. Use %shelp to learn more", strings.TrimPrefix(m.Content, configData.Prefix), configData.Prefix))
}

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

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
	err = s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
	if err != nil {
		fmt.Printf("Error deleting messages in channel %s: %s", m.ChannelID, err)
		s.ChannelMessageSend(m.ChannelID, "Unable to delete messages. Please check permissions and try again")
		return
	}
}

//UserInfo embed
func UserInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting guild: %s", err)
		return
	}
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

	prevMessage, err := s.ChannelMessages(m.ChannelID, 1, "", "", "")
	if err != nil {
		fmt.Printf("Error retrieving previous message: %s", err)
	}
	err = s.ChannelMessageDelete(m.ChannelID, prevMessage[0].ID)
	if err != nil {
		fmt.Printf("Error deleting previous message: %s", err)
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s was %sed because of reason: %s", user.Mention(), method, reason))
}

func KickUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	RemoveUser(s, m, false)
}

func BanUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	RemoveUser(s, m, true)
}

//ServerInfo embed
func ServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting guild: %s", err)
		return
	}
	guildOwner, err := s.User(g.OwnerID)
	if err != nil {
		fmt.Printf("Error getting guild owner: %s", err)
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, getServerEmbed(s, g, guildOwner))
}

func AddLoggingChannelCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := addLoggingChannel(m.ChannelID)
	if err == nil {
		s.ChannelMessageSend(m.ChannelID, "This channel will now be used for logging")
	} else if err == errChannelRegistered {
		s.ChannelMessageSend(m.ChannelID, "This channel is already set up for logging")
	} else {
		s.ChannelMessageSend(m.ChannelID, "There was an error while setting up this channel for logging")
		fmt.Printf("Error while configuring %s for logging: %s", m.ChannelID, err)
	}
}

func RemoveLoggingChannelCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := removeLoggingChannel(m.ChannelID)
	if err == nil {
		s.ChannelMessageSend(m.ChannelID, "This channel will no longer be used for logging")
	} else if err == errChannelNotRegistered {
		s.ChannelMessageSend(m.ChannelID, "This channel has not yet been configured for logging")
	} else {
		s.ChannelMessageSend(m.ChannelID, "There was an error while removing this channel's logging status")
		fmt.Printf("Error while removing %s from logging: %s", m.ChannelID, err)
	}
}

func TempMassRegister(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting guild: %s", err)
	}
	for _, mem := range guild.Members {
		if !mem.User.Bot {
			c, err := s.UserChannelCreate(mem.User.ID)
			if err != nil {
				fmt.Printf("Error creating channel: %s", err)
			}
			if configData.Users[mem.User.ID] == "" {
				s.ChannelMessageSend(c.ID, fmt.Sprintf("Please send me '%sregister {your first and last name} {grade as a number}' or ask for '%s help'", configData.Prefix, configData.Prefix))
			}
		}
	}
}
