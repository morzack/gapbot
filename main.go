package main

import (
	"fmt"
	"log"
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

	dg, err := discordgo.New("Bot " + configData.DiscordKey)
	if err != nil {
		fmt.Printf("Error creating discordgo session: %s", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(ready)
	dg.AddHandler(GuildCreate)
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

func GuildCreate(s *discordgo.Session, g *discordgo.Guild) {
	if len(g.ID) != 0 {
		for _, member := range g.Members {
			for _, user := range configData.Users {
				if !member.User.Bot {
					if member.User.ID == user {
						continue
					} else {
						channel, err := s.UserChannelCreate(member.User.ID)
						if err != nil {
							log.Printf("Error creating channel: %s", err)
						}
						s.ChannelMessageSend(channel.ID, "Please enter the code that has been emailed to you.")
					}
				}
			}
		}
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
		} else if getBotmod(s, m) {
			AdminCommand(s, m, content[0])
		} else {
			UserCommand(s, m, content[0])
		}
	}
}
