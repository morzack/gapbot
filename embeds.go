package main

import (
	"strconv"

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
				"`register name grade` - enter your [name] and [grade] (as a number) to register",
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
			"`server` - displays info about the server",
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
			"`removelog` - stop using a channel for logging",
		Inline: false,
	})
	return embed
}

func getUserEmbed(user *discordgo.User, s *discordgo.Session, g *discordgo.Guild) *discordgo.MessageEmbed {
	m, _ := s.GuildMember(g.ID, user.ID)
	b := "No"
	if user.Bot {
		b = "Yes"
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
			Name:   "ID",
			Value:  user.ID,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Bot?",
			Value:  b,
			Inline: true,
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
