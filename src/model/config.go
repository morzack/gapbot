package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	DiscordKey   string `json:"discord-key"`
	LastFMKey    string `json:"last-fm-key"`
	LastFMSecret string `json:"last-fm-secret"`

	Prefix string `json:"bot-prefix"`

	ModRoleName string   `json:"mod-role-name"`
	OpsUsers    []string `json:"ops-users"`
	MutedRole   string   `json:"muted-role"`

	ChannelsLogging []string `json:"channels-logging"`

	UserRoles  []string `json:"user-roles"`
	AdminRoles []string `json:"admin-roles"`
}

type ConfigInt interface {
	Load() error
	Get() Config
	Save() error
	Update(loggingChannels, userRoles, adminRoles []string) error
}

type configIntImpl struct {
	path   string
	config Config
}

func NewConfigInt(path string) (ConfigInt, error) {
	config := configIntImpl{
		path:   path,
		config: Config{},
	}
	if err := config.Load(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *configIntImpl) Load() error {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &c.config)
	if err != nil {
		return err
	}
	return nil
}

func (c *configIntImpl) Get() Config {
	return c.config
}

func (c *configIntImpl) Save() error {
	marshalled, err := json.Marshal(c.config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.path, marshalled, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *configIntImpl) Update(loggingChannels, userRoles, adminRoles []string) error {
	oldConfig := c.config

	if loggingChannels != nil {
		c.config.ChannelsLogging = loggingChannels
	}
	if userRoles != nil {
		c.config.UserRoles = userRoles
	}
	if adminRoles != nil {
		c.config.AdminRoles = adminRoles
	}

	if err := c.Save(); err != nil {
		c.config = oldConfig
		return err
	}
	return nil
}
