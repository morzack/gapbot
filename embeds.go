package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getBaseEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x993399, // purple for yukari, duh
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Contact @Valis#7360 or @Patchkat#9990 for more info",
			IconURL: "https://cdn.discordapp.com/attachments/384528590160658434/603418681405341706/pfp.png",
		},
	}
}

func getAvatarEmbed(user *discordgo.User) *discordgo.MessageEmbed {
	embed := getBaseEmbed()
	embed.Title = "Avatar"
	embed.Image = &discordgo.MessageEmbedImage{
		URL: user.AvatarURL(""),
	}
	return embed
}

func getDMHelpEmbed() *discordgo.MessageEmbed {
	embed := getBaseEmbed()
	embed.Title = "Gapbot Commands"
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name: "Basic",
			Value: "`ping` - ping the bot\n" +
				"`help` - it seems like you figured this one out\n" +
				"`bigletters` - make big text",
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name: "Last.fm",
			Value: "`lastregister [username]` - link your discord account to last.fm\n" +
				"`lastplayed` - get your most recently scrobbled song off last.fm\n" +
				"`lastloved` - get your most recently loved song off last.fm",
			Inline: false,
		},
	}
	return embed
}

func getServerHelpEmbed() *discordgo.MessageEmbed {
	embed := getDMHelpEmbed()
	embed.Title = "Gapbot Commands"
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "Server",
		Value: "`user @` - display [@]'s user info\n" +
			"`server` - displays info about the server\n" +
			"`addrole` - give yourself a role\n" +
			"`delrole` - remove a role from yourself\n" +
			"`roles` - list roles available to you\n" +
			"`myroles` - list your roles\n" +
			"`ownerhist` - see owner history of this server",
		Inline: false,
	})
	return embed
}

func getAdminHelpEmbed() *discordgo.MessageEmbed {
	embed := getServerHelpEmbed()
	embed.Title = "Gapbot Commands"
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "Admin",
		Value: "`purge #` - purge [#] messages (up to 99)\n" +
			"`ban @ reason` - ban [@] for a [reason]\n" +
			"`kick @ reason` - kick [@] for a [reason]\n" +
			"`addlog` - start using a channel for logging\n" +
			"`removelog` - stop using a channel for logging\n" +
			"`deregister @` - deregister [@]\n" +
			"`addrole @ role` - give [@] a [role]\n" +
			"`delrole @ role` - remove a [role] from [@]\n" +
			"`mute @` - give [@] the configured muted role\n" +
			"`unmute @` - remove the muted role from [@]",
		Inline: false,
	})
	return embed
}

func getRolesEmbed(roles []*discordgo.Role, title string) *discordgo.MessageEmbed {
	embed := getBaseEmbed()
	// embed.Title = "Available Roles"

	roleNames := []string{}
	for _, role := range roles {
		if role.Name != "@everyone" {
			roleNames = append(roleNames, role.Name)
		}
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   title,
		Value:  "- " + strings.Join(roleNames, "\n- "),
		Inline: false,
	})

	return embed
}

func getUserEmbed(user *discordgo.User, s *discordgo.Session, g *discordgo.Guild) *discordgo.MessageEmbed {
	m, _ := s.GuildMember(g.ID, user.ID)
	b := ""
	uid := ""
	idTitle := ""
	if user.Bot {
		b = "Yes"
		uid = user.ID
		idTitle = "ID"
	} else {
		b = "No"
		ok := false
		var tempUserData userStruct
		if tempUserData, ok = loadedUserData.Users[user.ID]; !ok {
			idTitle = "ID"
			uid = tempUserData.DiscordID
		} else {
			idTitle = "Name"
			uid = fmt.Sprintf("%s %s, %s", tempUserData.FirstName, tempUserData.LastName, getGradeString(tempUserData.Grade))
		}
	}
	i, _ := m.JoinedAt.Parse()
	t := i.Format("02/01/2006 15:04:05 EST")
	embed := getBaseEmbed()
	embed.Title = user.String()
	embed.Image = &discordgo.MessageEmbedImage{
		URL: user.AvatarURL(""),
	}
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   idTitle,
			Value:  uid,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Bot?",
			Value:  b,
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "Joined the server",
			Value:  t,
			Inline: false,
		},
	}
	return embed
}

func getServerEmbed(s *discordgo.Session, g *discordgo.Guild, u *discordgo.User) *discordgo.MessageEmbed {
	url := discordgo.EndpointGuildIcon(g.ID, g.Icon)
	embed := getBaseEmbed()
	embed.Title = g.Name
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "Members",
			Value:  strconv.Itoa(g.MemberCount),
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "Owner",
			Value:  u.String(),
			Inline: false,
		},
	}
	return embed
}

func getOwnerHistEmbed(s *discordgo.Session) *discordgo.MessageEmbed {

	embed := getBaseEmbed()
	embed.Title = "Owner History"
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name: "Patchkat",
			Value: "December 2018 (Server Creation) - June 2020",
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name: "Astroturtle",
			Value: "June 2020 - Present",
			Inline: false,
		},
	}
	return embed
	
}

func getLastFMTrackEmbed(mostRecentTrack lastFMSong) *discordgo.MessageEmbed {
	embed := getBaseEmbed()
	embed.Title = mostRecentTrack.Name
	embed.URL = mostRecentTrack.URL
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: mostRecentTrack.ImageThumbnail,
	}
	embed.Fields = []*discordgo.MessageEmbedField{}
	if mostRecentTrack.Artist != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Artist",
			Value:  mostRecentTrack.Artist,
			Inline: false,
		})
	}
	if mostRecentTrack.Album != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Album",
			Value:  mostRecentTrack.Album,
			Inline: false,
		})
	}
	return embed
}
