package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func getRole(s *discordgo.Session, u *discordgo.User, m *discordgo.Message, roleName string) bool {
	// see what the user permission level is
	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting roles: %s", err)
		return false
	}
	mem, err := s.GuildMember(m.GuildID, u.ID)
	if err != nil {
		fmt.Printf("Error getting members: %s", err)
		return false
	}
	for _, role := range roles {
		if role.Name == roleName {
			// check and see if user has role
			for _, v := range mem.Roles {
				if role.ID == v {
					return true
				}
			}
		}
	}
	return false
}

func getBotmod(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return getRole(s, m.Author, m.Message, configData.ModRoleName)
}

func getOwner(s *discordgo.Session, u *discordgo.User, m *discordgo.Message) (bool, error) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return false, err
	}
	return guild.OwnerID == u.ID, nil
}
