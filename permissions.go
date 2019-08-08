package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func getBotmod(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// see what the user permission level is
	author, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		fmt.Printf("Unable to get author data: %s", err)
		return false
	}
	modroles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		fmt.Printf("Error querying mod roles: %s", err)
		return false
	}
	for _, role := range modroles {
		if role.Name == configData.ModRoleName {
			// check and see if user has role
			for _, v := range author.Roles {
				if role.ID == v {
					return true
				}
			}
		}
	}
	return false
}
