// config.go
// used to access config file and parse it
package main

import (
	"errors"
	"os"
)

var (
	loadedConfigData configDataStruct

	configFile = "config.json"

	errChannelNotRegistered = errors.New("Channel not yet registered")
	errChannelRegistered    = errors.New("Channel already registered")
)

type configDataStruct struct {
	DiscordKey      string                       `json:"discord-key"`
	SourceDir       string                       `json:"source-dir"`
	Prefix          string                       `json:"bot-prefix"`
	ModRoleName     string                       `json:"mod-role-name"`
	ChannelsLogging []string                     `json:"channels-logging"`
	EnabledRoles    enabledPermissionRolesStruct `json:"roles-enabled"`
	MutedRole       string                       `json:"muted-role"`
}

type enabledPermissionRolesStruct struct {
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
	return loadJSON(configFile, &loadedConfigData)
}

func writeConfig() error {
	return writeJSON(configFile, loadedConfigData)
}

func removeLoggingChannel(channel string) error {
	for i, channelID := range loadedConfigData.ChannelsLogging {
		if channelID == channel {
			loadedConfigData.ChannelsLogging = append(loadedConfigData.ChannelsLogging[:i], loadedConfigData.ChannelsLogging[i+1:]...) // remove the channel, go doesn't have an easy removal function
			return writeConfig()
		}
	}
	return errChannelNotRegistered
}

func addLoggingChannel(channel string) error {
	for _, channelID := range loadedConfigData.ChannelsLogging {
		if channelID == channel {
			return errChannelRegistered
		}
	}
	loadedConfigData.ChannelsLogging = append(loadedConfigData.ChannelsLogging, channel)
	return writeConfig()
}
