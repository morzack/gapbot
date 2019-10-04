package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	userData UserData
)

type UserData struct {
	Users       map[string]string `json:"users"`
	NameChannel string            `json:"names-channel"`
}

func getUsersPath() (string, error) {
	if !debugMode {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			return "", errors.New("Use Linux and set your $HOME variable you filthy casual")
		}
		return fmt.Sprintf("%s/.config/gapbot/users.json", homePath), nil
	}
	return fmt.Sprintf("./users.json"), nil
}

func loadUsers() error {
	configPath, err := getUsersPath()
	if err != nil {
		return err
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &userData)
	if err != nil {
		return err
	}

	return nil
}

func writeUsers() error {
	configPath, err := getUsersPath()
	if err != nil {
		return err
	}
	marshalledJSON, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, marshalledJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Register(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Fields(strings.TrimPrefix(m.Content, configData.Prefix))
	if userData.Users[m.Author.ID] == "" {
		if len(content) == 4 {
			userData.Users[m.Author.ID] = strings.Title(content[1]) + " " + strings.Title(content[2])
			s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s: %s, %sth grade", m.Author.Username, userData.Users[m.Author.ID], content[3]))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are missing something"))
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are already registered as: %s", userData.Users[m.Author.ID]))
	}
	return writeUsers()
}

func Deregister(s *discordgo.Session, u *discordgo.User) error {
	if userData.Users[u.ID] != "" {
		delete(userData.Users, u.ID)
		s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s was removed as a member", u.Username))
	} else {
		fmt.Printf("That user is not registered")
	}
	return writeUsers()
}
