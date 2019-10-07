package main

import (
	"fmt"
	"strings"

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

func getAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return getRole(s, m.Author, m.Message, configData.ModRoleName)
}

func Role(s *discordgo.Session, m *discordgo.MessageCreate, rem bool) {
	content := strings.Fields(strings.TrimPrefix(m.Content, configData.Prefix))
	role := content[1]
	//Command for admins
	if getAdmin(s, m) {
		//If they mention someone
		if len(m.Mentions) > 0 {
			AddDelRole(m, m.Mentions[0], s, strings.ToLower(role), rem)
		} else { // if they don't mention anyone
			AddDelRole(m, m.Author, s, strings.ToLower(role), rem)
		}
	} else { //If not admin
		AddDelRole(m, m.Author, s, strings.ToLower(role), rem)
	}
}

// need role name
func AddDelRole(m *discordgo.MessageCreate, u *discordgo.User, s *discordgo.Session, r string, rem bool) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting guild: %s", err)
	}
	rs, err := s.GuildRoles(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting roles: %s", err)
	}
	for _, role := range rs {
		if strings.ToLower(role.Name) == r {
			if rem {
				err := s.GuildMemberRoleRemove(g.ID, u.ID, role.ID)
				if err != nil {
					fmt.Printf("Error removing role %s to %s", r, m.Author.Username)
				} else {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role ``%s`` succesfully removed.", r))
					return
				}
			} else {
				err := s.GuildMemberRoleAdd(g.ID, u.ID, role.ID)
				if err != nil {
					fmt.Printf("Error adding role %s to %s", r, m.Author.Username)
				} else {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role ``%s`` succesfully added.", r))
					return
				}
			}
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role ``%s`` not found. Please try again or contact an admin.", r))
	}
}
