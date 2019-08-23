package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"strconv"
	"math"

	"github.com/bwmarrin/discordgo"
)

func SendEmail(m *discordgo.MessageCreate) {
	user := configData.Username
	pass := configData.Password
	id, _ := strconv.Atoi(m.Author.ID)
	code := (configData.Values[0] * int(math.Pow(float64(id), 2))) + (configData.Values[1] * id) + configData.Values[2]
	content := strings.Fields(strings.TrimPrefix(m.Content, configData.Prefix))
	to := fmt.Sprintf("%v_%v@caryacademy.org", content[1], content[2])
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
