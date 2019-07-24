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
)

type ConfigData struct {
	DiscordKey string `json:"discord-key"`
	SourceDir  string `json:"source-dir"`
	Prefix     string `json:"bot-prefix"`
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
