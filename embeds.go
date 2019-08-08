// embeds.go
// creates various embeds

package main

import (
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

func getHelpEmbed() *discordgo.MessageEmbed {
	embed := getBaseEmbed()
	embed.Title = "Gapbot help"
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name: "Basic",
			Value: "`ping` - ping the bot\n" +
				"`help` - it seems like you figured this one out\n" +
				"`avatar` `@user` - display an image of your avatar or up to 4 others\n" +
				"`user` `@user` - display information about a user",
			Inline: false,
		},
	}
	return embed
}

func getAdminHelpEmbed() *discordgo.MessageEmbed {
	embed := getBaseEmbed()
	embed.Title = "Admin Commands"
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "Admin",
			Value:  "`purge #` - purge a given number of messages (up to 99) from a channel",
			Inline: false,
		},
	}
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
