package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	userData UserData
	student  Student
	userFile = "users.json"

	errUserInputInvalid  = errors.New("User input invalid")
	errUserRegistered    = errors.New("User is already registered")
	errUserNotRegistered = errors.New("User not yet registered")
)

type UserData struct {
	Users       map[string]Student `json:"users"`
	NameChannel string             `json:"names-channel"`
}

type Student struct {
	ID    string
	Name  string
	Grade string
	Money int
}

func loadUsers() error {
	err := loadJSON(userFile, &userData)
	if err != nil {
		return err
	}
	if len(userData.Users) == 0 {
		userData.Users = make(map[string]Student)
	}

	return nil
}

func writeUsers() error {
	return writeJSON(userFile, userData)
}

func Register(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Fields(strings.TrimPrefix(m.Content, configData.Prefix))
	r := regexp.MustCompile(`^(?P<first>\w+) (?P<last>\w+) (?P<grade>[6-9]|1[0-2])$`)
	subMatch := r.FindStringSubmatch(strings.Join(content[1:], " "))

	if userData.Users[m.Author.ID] == "" {
		if subMatch == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You entered something wrong.  Don't forget your grade should be 6-12."))
			return errUserInputInvalid
		} else {
			userData.Users[m.Author.ID] = fmt.Sprintf("%s %s", strings.Title(subMatch[1]), strings.Title(subMatch[2]))
			s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s: %s, %sth grade", m.Author.Username, userData.Users[m.Author.ID], subMatch[3]))
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are already registered as: %s", userData.Users[m.Author.ID]))
		return errUserRegistered
	}
	return writeUsers()
}

func Deregister(s *discordgo.Session, m *discordgo.MessageCreate) error {
	u := m.Mentions[0]
	if userData.Users[u.ID] != "" {
		delete(userData.Users, u.ID)
		s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s was removed as a member", u.Username))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s was unable to be deregistered -- not registered in the first place", u.Username))
		return errUserNotRegistered
	}
	return writeUsers()
}
