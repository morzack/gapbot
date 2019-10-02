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
	Users map[string]string `json:"users"`
	Roles [][]string        `json:"roles"`
}

func getConfigPath() (string, error) {
	if !debugMode {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			return "", errors.New("Use Linux and set your $HOME variable you filthy casual")
		}
		return fmt.Sprintf("%s/.config/gapbot/config.json", homePath), nil
	}
	return fmt.Sprintf("./config.json"), nil
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
	configPath, err := getConfigPath()
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
	configPath, err := getConfigPath()
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
	content := strings.Fields(strings.TrimPrefix(m.Content, userData.Prefix))
	if userData.Users[m.Author.ID] == "" {
		if len(content) == 4 {
			userData.Users[m.Author.ID] = strings.Title(content[1]) + " " + strings.Title(content[2])
			s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s: %s, %sth grade", m.Author.Username, userData.Users[m.Author.ID], content[3]))
			err := s.GuildMemberRoleAdd(userData.Roles[4][1], m.Author.ID, userData.Roles[4][0])
			if err != nil {
				fmt.Printf("Here da error: %s", err)
			}
			fmt.Printf("GuildID: %s", userData.Roles[4][1])
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are missing something"))
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are already registered as: %s", userData.Users[m.Author.ID]))
	}
	return writeConfig()
}

func Deregister(s *discordgo.Session, u *discordgo.User) error {
	if userData.Users[u.ID] != "" {
		delete(userData.Users, u.ID)
		s.GuildMemberRoleRemove(userData.Roles[4][1], u.ID, userData.Roles[4][0])
		s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s was removed as a member", u.Username))
	} else {
		fmt.Printf("That user is not registered")
	}
	return writeConfig()
}
