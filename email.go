package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SendEmail(m *discordgo.MessageCreate) {
	user := configData.Username
	pass := configData.Password
	content := strings.Fields(strings.TrimPrefix(m.Content, configData.Prefix))
	to := fmt.Sprintf("%v_%v@caryacademy.org", content[1], content[2])

	msg := "From: CA Discord Bot" + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		fmt.Sprintf("%v", content[3:])

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", user, pass, "smtp.gmail.com"),
		user, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}
