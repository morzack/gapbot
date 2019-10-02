// config.go
// used to access config file and parse it
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
	configData ConfigData
	debugMode  = true

	errChannelNotRegistered = errors.New("Channel not yet registered")
	errChannelRegistered    = errors.New("Channel already registered")
	errUserRegistered       = errors.New("User is already registered")
)

type ConfigData struct {
	DiscordKey      string            `json:"discord-key"`
	SourceDir       string            `json:"source-dir"`
	Prefix          string            `json:"bot-prefix"`
	ModRoleName     string            `json:"mod-role-name"`
	MemberRole      string            `json:"member-role"`
	ChannelsLogging []string          `json:"channels-logging"`
	Users           map[string]string `json:"users"`
	NameChannel     string            `json:"names-channel"`
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

func loadConfig() error {
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
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return err
	}

	return nil
}

func writeConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	marshalledJSON, err := json.Marshal(configData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, marshalledJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func removeLoggingChannel(channel string) error {
	for i, channelID := range configData.ChannelsLogging {
		if channelID == channel {
			configData.ChannelsLogging = append(configData.ChannelsLogging[:i], configData.ChannelsLogging[i+1:]...) // remove the channel, go doesn't have an easy removal function
			return writeConfig()
		}
	}
	return errChannelNotRegistered
}

func addLoggingChannel(channel string) error {
	for _, channelID := range configData.ChannelsLogging {
		if channelID == channel {
			return errChannelRegistered
		}
	}
	configData.ChannelsLogging = append(configData.ChannelsLogging, channel)
	return writeConfig()
}

func Register(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Fields(strings.TrimPrefix(m.Content, configData.Prefix))
	if configData.Users[m.Author.ID] == "" {
		configData.Users[m.Author.ID] = content[1] + " " + content[2]
		s.ChannelMessageSend(configData.NameChannel, fmt.Sprintf("%s: %s, %sth grade", m.Author.Username, configData.Users[m.Author.ID], content[3]))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are already registered as: %s", configData.Users[m.Author.ID]))
	}
	return writeConfig()
}

func Deregister(s *discordgo.Session, u *discordgo.User) error {
	if configData.Users[u.ID] == "" {
		delete(configData.Users, u.ID)
	} else {
		fmt.Printf("That user is not registered")
	}
	return writeConfig()
}
