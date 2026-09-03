package interactions

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "pts",
		Description: "sistema de pontuação",
		Type:        discordgo.ChatApplicationCommand,
		Version:     "v0.8",
	},
	{
		Name:        "axx",
		Description: "AI",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "mensagem",
				Description: "usa mensagem",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
		Version: "v0.1",
	},
}
