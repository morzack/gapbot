package main

import (
	"fmt"

	"github.com/ParrotCaws/insults"
	"github.com/bwmarrin/discordgo"
)

var shakespearean = insults.InsultList{
	Column1: []string{"artless", "bawdy", "beslubbering", "bootless", "churlish", "cockered", "clouted", "craven", "currish", "darkish", "dissembling",
		"droning", "errant", "fawning", "fobbing", "froward", "frothy", "gleeking", "goatish", "gorbellied", "impertinent", "infectious", "jarring",
		"loggerheaded", "lumpish", "mammering", "mangled", "mewling", "paunchy", "pribbling", "puking", "puny", "qualling", "rank", "reeky", "roguish",
		"ruttish", "saucy", "spleeny", "spongy", "surly", "tottering", "unmuzzled", "vain", "venomed", "villainous", "warped", "wayward", "weedy", "yeasty"},
	Column2: []string{"base-court", "bat-fowling", "beef-witted", "beetle-headed", "boil-brained", "clapper-clawed", "clay-brained", "common-kissing", "crook-plated",
		"dismal-dreaming", "dizzy-eyed", "doghearted", "dread-boiled", "earth-vexing", "elf-skinned", "fat-kidneyed", "fen-sucked", "flap-mouthed", "fly-bitten", "folly-fallen",
		"fool-born", "full-gorged", "guts-gripping", "half-faced", "hasty-witted", "hedge-born", "hell-hated", "idle-headed", "ill-breeding", "ill-nurtured", "knotty-pated",
		"milk-livered", "motley-minded", "onion-eyed", "plume-plucked", "pottle-deep", "pox-marked", "reeling-ripe", "rough-hewn", "rude-growing", "rump-fed", "shard-borne",
		"sheep-biting", "spur-galled", "swag-bellied", "tardy-gaited", "tickle-brained", "toad-spotted", "urchin-snouted", "weather-bitten"},
	Column3: []string{"apple-john", "baggage", "barnacle", "bladder", "boar-pig", "bugbear", "bum-bailey", "canker-blossom", "clack-dish", "clotpole", "coxcomb", "codpiece",
		"death-token", "dewberry", "flap-dragon", "flax-wench", "flirt-gill", "foot-licker", "fustilarian", "giglet", "gudgeon", "haggard", "harpy", "hedge-pig", "horn-beast",
		"hugger-mugger", "joithead", "lewdster", "lout", "maggot-pie", "malt-worm", "mammet", "meals", "minnow", "miscreant", "moldwarp", "mumble-news", "nut-hook", "pigeon-egg",
		"pignut", "puttock", "pumpion", "ratsbane", "scut", "skainsmate", "strumpet", "varlot", "vassal", "whey-face", "wagtail"},
	You: "Thou",
}

func Insult(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>: %s", m.Mentions[0].ID, shakespearean.RandInsult()))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>: You didn't mention anyone. %s", m.Author.ID, shakespearean.RandInsult()))
}
