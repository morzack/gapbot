package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// load config
	err := loadConfig()
	if err != nil {
		fmt.Printf("Error getting config: %s", err)
		return
	}

	// load user config
	err = loadUsers()
	if err != nil {
		fmt.Printf("Error getting users: %s", err)
		return
	}

	dg, err := discordgo.New("Bot " + configData.DiscordKey)
	if err != nil {
		fmt.Printf("Error creating discordgo session: %s", err)
		return
	}

	dg.AddHandler(guildMemberAdd)
	dg.AddHandler(messageCreate)
	dg.AddHandler(ready)

	err = dg.Open()
	if err != nil {
		fmt.Printf("Error opening connection: %s", err)
		return
	}

	fmt.Printf("Bot launched. Send interrupt to exit")
	fmt.Printf("invite link: https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8\n", dg.State.User.ID)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// close the session
	dg.Close()
}

// called when discord sends the ready state
func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, fmt.Sprintf("Type %shelp", configData.Prefix))
}

// called when a new user enters the server
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if _, present := userData.Users[m.User.ID]; !present {
		c, err := s.UserChannelCreate(m.User.ID)
		if err != nil {
			fmt.Printf("Error creating channel: %s", err)
			return
		}
		s.ChannelMessageSend(c.ID, fmt.Sprintf("Please send me `%sregister {first name} {last name} {grade #}` (e.g. `%sregister Jono Jenkens 12`)", configData.Prefix, configData.Prefix))
	}
}

// called when a message is created on a channel this has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all of this bot's messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	isDM := m.GuildID == ""

	r, _ := regexp.Compile("https*:\\/\\/discord.gg\\/(invite\\/)*[a-zA-Z0-9]{6}")
	if !isDM && r.MatchString(m.Content) {
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		logServerInvite(s, m)
		if err != nil {
			fmt.Printf("Failed to delete invite %s: %s", m.Content, err)
			return
		}
	}
	if strings.HasPrefix(m.Content, configData.Prefix) {
		content := strings.Fields(strings.TrimPrefix(m.Content, configData.Prefix))
		if isDM {
			DMCommand(s, m, content[0])
		} else {
			if _, present := userData.Users[m.Author.ID]; present {
				if getAdmin(s, m) {
					AdminCommand(s, m, content[0])
				} else {
					UserCommand(s, m, content[0])
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "You need to register before using the bot in the server!")
				s.ChannelMessageSend(m.ChannelID, "Message an admin if you need help.")
			}
		}
	}
}
