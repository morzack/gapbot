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
	User        *discordgo.User `json:"user"`
	Name        string          `json:"name"`
	Grade       string          `json:"grade"`
	Currency    int             `json:"funds"`
	Infractions int             `json:"infractions"`
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

	var student Student

	if _, present := userData.Users[m.Author.ID]; !present {
		if subMatch == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You entered something wrong.  Don't forget your grade should be 6-12."))
			return errUserInputInvalid
		}
		student.User = m.Author
		student.Name = fmt.Sprintf("%s %s", strings.Title(subMatch[1]), strings.Title(subMatch[2]))
		student.Grade = subMatch[3]
		userData.Users[m.Author.ID] = student
		s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s: %s, %sth grade", m.Author.Username, userData.Users[m.Author.ID].Name, subMatch[3]))
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have registered as: %s, %sth grade", userData.Users[m.Author.ID].Name, subMatch[3]))
		return writeUsers()
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are already registered as: %s", userData.Users[m.Author.ID].Name))
	return errUserRegistered
}

func Deregister(s *discordgo.Session, m *discordgo.MessageCreate) error {
	u := m.Mentions[0]
	if _, present := userData.Users[u.ID]; present {
		delete(userData.Users, u.ID)
		s.ChannelMessageSend(userData.NameChannel, fmt.Sprintf("%s was removed as a member", u.Username))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s was unable to be deregistered -- not registered in the first place", u.Username))
		return errUserNotRegistered
	}
	return writeUsers()
}
