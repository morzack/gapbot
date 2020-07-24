// config.go
// used to access config file and parse it
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	loadedConfigData configDataStruct

	configFile = "config.json"

	errChannelNotRegistered = errors.New("Channel not yet registered")
	errChannelRegistered    = errors.New("Channel already registered")
)

type configDataStruct struct {
	DiscordKey      string                       `json:"discord-key"`
	LastFMKey       string                       `json:"last-fm-key"`
	LastFMSecret    string                       `json:"last-fm-secret"`
	SourceDir       string                       `json:"source-dir"`
	Prefix          string                       `json:"bot-prefix"`
	ModRoleName     string                       `json:"mod-role-name"`
	ChannelsLogging []string                     `json:"channels-logging"`
	EnabledRoles    enabledPermissionRolesStruct `json:"roles-enabled"`
	MutedRole       string                       `json:"muted-role"`
	OpsUsers        []string                     `json:"ops-users"`
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

func startBotConfig(s *discordgo.Session) {
	s.UpdateStatus(0, fmt.Sprintf("Type %shelp", loadedConfigData.Prefix))
}
