// config.go
// used to access config file and parse it
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	configData ConfigData
	debugMode  = true

	errChannelNotRegistered = errors.New("Channel not yet registered")
	errChannelRegistered    = errors.New("Channel already registered")
	errUserRegistered       = errors.New("User is already registered")
)

type ConfigData struct {
	DiscordKey      string         `json:"discord-key"`
	SourceDir       string         `json:"source-dir"`
	Prefix          string         `json:"bot-prefix"`
	ModRoleName     string         `json:"mod-role-name"`
	ChannelsLogging []string       `json:"channels-logging"`
	Users           map[int]string `json:"users"`
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
