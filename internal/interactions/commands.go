package interactions

import (
	"github.com/bwmarrin/discordgo"
)

type command struct {
	Cmd *discordgo.ApplicationCommand

	// enabled for public access
	Public bool
	// TODO: make use of
	// Run    func(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type scoped struct {
	Global *discordgo.ApplicationCommand
	Guild  *discordgo.ApplicationCommand
}

var commands = []command{
	{
		Public: false,
		Cmd: &discordgo.ApplicationCommand{
			Name:        "pts",
			Description: "sistema de pontuação",
			Type:        discordgo.ChatApplicationCommand,
			Version:     "v0.8",
		},
	},
	{
		Public: true,
		Cmd: &discordgo.ApplicationCommand{
			Name:        "axx",
			Description: "AI",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "mensagem",
					Description: "sua mensagem",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
			Version: "v0.1",
		},
	},
}

// Commands get application commands with infos and wheter is enabled
func Commands() struct{} {
	var pc []*discordgo.ApplicationCommand
	var prc []*discordgo.ApplicationCommand

	for _, command := range commands {
		if command.Public {
			pc = append(pc, command.Cmd)
		} else {
			prc = append(prc, command.Cmd)
		}
	}
	return commands
}
