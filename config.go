// config.go
// used to access config file and parse it
package main

import (
	"errors"
	"os"
)

var (
	configData ConfigData

	configFile = "config.json"

	errChannelNotRegistered = errors.New("Channel not yet registered")
	errChannelRegistered    = errors.New("Channel already registered")
)

type ConfigData struct {
	DiscordKey      string       `json:"discord-key"`
	SourceDir       string       `json:"source-dir"`
	Prefix          string       `json:"bot-prefix"`
	ModRoleName     string       `json:"mod-role-name"`
	ChannelsLogging []string     `json:"channels-logging"`
	EnabledRoles    RolesEnabled `json:"roles-enabled"`
	MutedRole       string       `json:"muted-role"`
}

type RolesEnabled struct {
	UserRoles  []string `json:"user"`
	AdminRoles []string `json:"admin"`
}

func getDebugMode() bool {
	if _, err := os.Stat("./production"); err == nil {
		return false
	}
	return true
}

func loadConfig() error {
	return loadJSON(configFile, &configData)
}

func writeConfig() error {
	return writeJSON(configFile, configData)
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
