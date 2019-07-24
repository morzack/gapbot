// help.go
// used to create embeds for help function

package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	HelpEmbed = &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x993399, // purple for yukari, duh
		Title:  "Gapbot help",
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Contact @Valis#7360 or @Patchkat#9990 for more info",
			IconURL: "https://cdn.discordapp.com/attachments/384528590160658434/603418681405341706/pfp.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Basic",
				Value: "`ping` - ping the bot\n" +
					"`help` - it seems like you figured this one out\n" +
					"`purge` - purge x amount of messages from the channel\n" +
					"`avatar` - display an image of your avatar",
				Inline: false,
			},
		},
	}
)
