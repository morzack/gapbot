package main

import (
	"fmt"
	"log"
	"math"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func RegisterUsers(s *discordgo.Session, g *discordgo.Guild) {
	if len(g.ID) != 0 {
		for _, member := range g.Members {
			log.Printf("Username: " + member.User.Username)
			if !member.User.Bot {
				log.Printf("not a bot!")
				err := registerUser(member.User)
				if err == nil {
					log.Printf("not registered!")
					channel, err := s.UserChannelCreate(member.User.ID)
					if err != nil {
						log.Printf("Error creating channel: %s", err)
					}
					s.ChannelMessageSend(channel.ID, "Please enter your first and last name, as they appear in your email address")
					s.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.MessageCreate) {
						if len(m.GuildID) == 0 {
							if !m.Author.Bot {
								log.Printf("Got second message!" + strconv.Itoa(len(m.GuildID)))
								content := strings.Fields(m.Content)
								log.Printf("Name entered: " + content[0] + content[1])
								SendEmail(m)
							}
						}
					})
				}
			}
		}
	}
}

func SendEmail(m *discordgo.MessageCreate) {
	user := configData.Username
	pass := configData.Password
	id, _ := strconv.Atoi(m.Author.ID)
	code := (configData.Values[0] * int(math.Pow(float64(id), 2))) + (configData.Values[1] * id) + configData.Values[2]
	content := strings.Fields(m.Content)
	to := fmt.Sprintf("%v_%v@caryacademy.org", content[0], content[1])
	log.Printf(strconv.Itoa(code))

	msg := "From: CA Discord Bot" + "\n" +
		"To: " + to + "\n" +
		"Subject: Send this code to the bot\n\n" +
		strconv.Itoa(code)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", user, pass, "smtp.gmail.com"),
		user, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}

func TestEmail(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Printf("Error creating channel: %s", err)
	}
	s.ChannelMessageSend(channel.ID, "Please enter the code that has been emailed to you.")
}
