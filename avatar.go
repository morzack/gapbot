// avatar.go
// creates the embed for the avatar command

package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	AvatarEmbed = &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  #7FFFD4,
		Title:  "Avatar",
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Contact @Valis#7360 or @Patchkat#9990 if something is broken",
			IconURL: "https://cdn.discordapp.com/attachments/384528590160658434/603418681405341706/pfp.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "The Avatar",
				Value: &discordgo.messageAuthor.AvatarURL,
				Inline: false,
			},
		},
	}
)
