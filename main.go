package main

import (
	"fmt"
	"os"
	"os/signal"
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

	// get discordgo instance
	dg, err := discordgo.New("Bot " + loadedConfigData.DiscordKey)
	if err != nil {
		fmt.Printf("Error creating discordgo session: %s", err)
		return
	}

	// get lastfm instance
	createLastFMAPIInstance()

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
	startBotConfig(s)
}

// called when a new user enters the server
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if _, present := loadedUserData.Users[m.User.ID]; !present {
		c, err := s.UserChannelCreate(m.User.ID)
		if err != nil {
			fmt.Printf("Error creating channel: %s", err)
			return
		}
		s.ChannelMessageSend(c.ID, fmt.Sprintf("Please send me `%sregister {first name} {last name} {grade # (or 'a' if you're an alumni)}` (e.g. `%sregister Jono Jenkens 12`)", loadedConfigData.Prefix, loadedConfigData.Prefix))
	}
}

// called when a message is created on a channel this has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all of this bot's messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	isDM := m.GuildID == ""

	if strings.HasPrefix(m.Content, loadedConfigData.Prefix) {
		content := strings.Fields(strings.TrimPrefix(m.Content, loadedConfigData.Prefix))
		if isDM {
			dmCommand(s, m, content[0])
		} else {
			if _, present := loadedUserData.Users[m.Author.ID]; present {
				if getOps(s, m) {
					opsCommand(s, m, content[0])
				} else if getAdmin(s, m) {
					adminCommand(s, m, content[0])
				} else {
					userCommand(s, m, content[0])
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "You need to register before using the bot in the server!")
				s.ChannelMessageSend(m.ChannelID, "Message an admin if you need help.")
			}
		}
	}
}
