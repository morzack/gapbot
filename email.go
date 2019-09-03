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
	var h func()
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
					RegisteringUser(s, channel, h)
				}
			}
		}
	}
}

func RegisteringUser(s *discordgo.Session, channel *discordgo.Channel, h func()) {
	s.ChannelMessageSend(channel.ID, "Please enter your cary academy email!")
	h = s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		c, _ := s.Channel(m.ChannelID)
		if !m.Author.Bot && c.GuildID == "" {
			s.ChannelMessageSend()
			h()
		}
	})
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

func CheckCode(u *discordgo.User, c string) bool {
	id, _ := strconv.Atoi(u.ID)
	code := (configData.Values[0] * int(math.Pow(float64(id), 2))) + (configData.Values[1] * id) + configData.Values[2]
	if c == strconv.Itoa(code) {
		err := RegisterUser(u)
		if err != nil {
			log.Printf("Error adding user: %s", err)
			return false
		}
	}
	return true
}
