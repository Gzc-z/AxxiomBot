// TODO: find a way to sort commands with servers

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
	// public commands
	Global []*discordgo.ApplicationCommand

	// private commands
	Guild []*discordgo.ApplicationCommand
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

// Commands get application commands with infos and whether is enabled
func Commands() scoped {
	var comm scoped
	for _, c := range commands {
		if c.Public {
			comm.Global = append(comm.Global, c.Cmd)
		} else {
			comm.Guild = append(comm.Guild, c.Cmd)
		}
	}
	return comm
}
