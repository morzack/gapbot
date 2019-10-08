package main

import (
	"fmt"
	"sort"
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

func Role(s *discordgo.Session, m *discordgo.MessageCreate, removing bool) {
	// if an admin is modifying another persons' role
	if getAdmin(s, m) && len(m.Mentions) > 0 {
		role := strings.SplitN(strings.TrimPrefix(m.Content, configData.Prefix), " ", 3)[2]
		AddDelRole(m, m.Mentions[0], s, role, removing)
		return
	}
	role := strings.SplitN(strings.TrimPrefix(m.Content, configData.Prefix), " ", 2)[1]
	AddDelRole(m, m.Author, s, role, removing)
}

func AddDelRole(m *discordgo.MessageCreate, u *discordgo.User, s *discordgo.Session, roleName string, removing bool) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting guild: %s", err)
		return
	}
	rs, err := s.GuildRoles(m.GuildID)
	if err != nil {
		fmt.Printf("Error getting roles: %s", err)
	}
	for _, role := range rs {
		if strings.ToLower(role.Name) == roleName {
			if removing {
				err := s.GuildMemberRoleRemove(g.ID, u.ID, role.ID)
				if err != nil {
					fmt.Printf("Error removing role %s to %s", roleName, m.Author.Username)
					return
				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role ``%s`` succesfully removed.", roleName))
				return
			}
			err := s.GuildMemberRoleAdd(g.ID, u.ID, role.ID)
			if err != nil {
				fmt.Printf("Error adding role %s to %s", roleName, m.Author.Username)
				return
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role ``%s`` succesfully added.", roleName))
			return
		}
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role ``%s`` not found. Please try again or contact an admin.", roleName))
}

func getGuildRoles(s *discordgo.Session, m *discordgo.MessageCreate) ([]*discordgo.Role, error) {
	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})
	return roles, nil
}

func getAvailableRoles(s *discordgo.Session, m *discordgo.MessageCreate, u *discordgo.User) ([]*discordgo.Role, error) {
	admin := getAdmin(s, m)
	roles := []*discordgo.Role{}
	guildRoles, err := getGuildRoles(s, m)
	if err != nil {
		return nil, err
	}
	for _, role := range guildRoles {
		if itemInSlice(role.Name, configData.EnabledRoles.UserRoles) {
			roles = append(roles, role)
		}
		if admin && itemInSlice(role.Name, configData.EnabledRoles.AdminRoles) {
			roles = append(roles, role)
		}
	}
	return roles, nil
}
