package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	loadedUserData userDataStruct
	userFile       = "users.json"

	errUserInputInvalid  = errors.New("User input invalid")
	errUserRegistered    = errors.New("User is already registered")
	errUserNotRegistered = errors.New("User not yet registered")
)

type userDataStruct struct {
	Users       map[string]userStruct `json:"users"`
	NameChannel string                `json:"names-channel"`
}

type userStruct struct {
	DiscordID string `json:"id"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Grade     int    `json:"grade"`
}

func loadUsers() error {
	err := loadJSON(userFile, &loadedUserData)
	if err != nil {
		return err
	}

	// if user map is empty it needs to not be
	// i'm going to be honest -- i think this is redundant
	// but it doesn't hurt so it'll stay in for now
	if len(loadedUserData.Users) == 0 {
		loadedUserData.Users = make(map[string]userStruct)
	}

	return nil
}

func writeUsers() error {
	return writeJSON(userFile, loadedUserData)
}

func registerUserCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Fields(strings.TrimPrefix(m.Content, loadedConfigData.Prefix))
	r := regexp.MustCompile(`^(?P<first>\w+) (?P<last>\w{2,}) (?P<grade>[6-9]|1[0-2]|a)$`)
	subMatch := r.FindStringSubmatch(strings.Join(content[1:], " "))

	if _, present := loadedUserData.Users[m.Author.ID]; !present {
		if subMatch == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You entered something wrong.  Don't forget your grade should be 6-12 (or 'a' if you're an alumni)."))
			return errUserInputInvalid
		}
		var grade int
		var err error
		if subMatch[3] == "a" {
			grade = 13
		} else {
			grade, err = strconv.Atoi(subMatch[3])
			if err != nil {
				fmt.Printf("Error parsing grade (input was %s) to int: %s", subMatch[3], err)
				return errUserInputInvalid
			}
		}
		loadedUserData.Users[m.Author.ID] = userStruct{
			FirstName: strings.Title(subMatch[1]),
			LastName:  strings.Title(subMatch[2]),
			DiscordID: m.Author.ID,
			Grade:     grade,
		}
		pushNewUser(m.Author, s)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have registered as: %s %s, %s", loadedUserData.Users[m.Author.ID].FirstName, loadedUserData.Users[m.Author.ID].LastName, getGradeString(loadedUserData.Users[m.Author.ID].Grade)))
		return writeUsers()
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are already registered as: %s %s", loadedUserData.Users[m.Author.ID].FirstName, loadedUserData.Users[m.Author.ID].LastName))
	return errUserRegistered
}

func pushNewUser(user *discordgo.User, s *discordgo.Session) {
	s.ChannelMessageSend(loadedUserData.NameChannel, fmt.Sprintf("%s: %s %s, %s", user.Mention(), loadedUserData.Users[user.ID].FirstName, loadedUserData.Users[user.ID].LastName, getGradeString(loadedUserData.Users[user.ID].Grade)))
}

func deregisterUserCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if len(m.Mentions) <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Unable to deregister a user without a mention.")
		return errUserNotRegistered
	}
	u := m.Mentions[0]
	if _, present := loadedUserData.Users[u.ID]; present {
		delete(loadedUserData.Users, u.ID)
		s.ChannelMessageSend(loadedUserData.NameChannel, fmt.Sprintf("%s was removed as a member", u.Username))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s was unable to be deregistered -- not registered in the first place", u.Username))
		return errUserNotRegistered
	}
	return writeUsers()
}

func getGradeString(grade int) string {
	// if grade is 13 then they're an alumn
	if grade > 12 {
		return "Alumni"
	}
	return fmt.Sprintf("%dth grade", grade)
}
