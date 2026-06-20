package handlers

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "pts",
		Description: "points system",
		Type:        discordgo.ChatApplicationCommand,
		// Version:     "v1.0", not yet
	},
}
