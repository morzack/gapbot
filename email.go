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
				reg := IsUserRegistered(member.User)
				if !reg {
					log.Printf("not registered!")
					channel, err := s.UserChannelCreate(member.User.ID)
					if err != nil {
						log.Printf("Error creating channel: %s", err)
					}
					s.ChannelMessageSend(channel.ID, "Please enter your first and last name, as they appear in your email address")
					s.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.MessageCreate) {
						c, _ := s.Channel(m.ChannelID)
						if !m.Author.Bot && c.GuildID == "" {
							SendEmail(m)
							s.ChannelMessageSend(channel.ID, "Please enter your first and last name, as they appear in your email address")
						}
					})
					s.AddHandlerOnce(func(s *discordgo.Session, m *discordgo.MessageCreate) {
						c, _ := s.Channel(m.ChannelID)
						if !m.Author.Bot && c.GuildID == "" {
							content := strings.Fields(m.Content)
							log.Printf("code entered: " + content[0])
							CheckCode(m.Author, content[0])
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

func CheckCode(u *discordgo.User, c string) {
	id, _ := strconv.Atoi(u.ID)
	code := (configData.Values[0] * int(math.Pow(float64(id), 2))) + (configData.Values[1] * id) + configData.Values[2]
	if c == strconv.Itoa(code) {
		log.Printf("Registering user!")
		err := RegisterUser(u)
		if err != nil {
			log.Printf("Error adding user: %s", err)
		}
	}
}
