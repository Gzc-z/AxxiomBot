package interactions

import (
	"github.com/bwmarrin/discordgo"
)

type cmds struct {
	// enabled for public access
	Enabled bool
	Cmd     *discordgo.ApplicationCommand
}

var commands = []cmds{
	{
		Enabled: false,
		Cmd: &discordgo.ApplicationCommand{
			Name:        "pts",
			Description: "sistema de pontuação",
			Type:        discordgo.ChatApplicationCommand,
			Version:     "v0.8",
		},
	},
	{
		Enabled: true,
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

// GetCommands get application commands with infos and wheter is enabled
func GetCommands() []cmds {
	cmd := make([]cmds, 0)
	for _, c := range commands {
		if c.Enabled {
			cmd = append(cmd, c)
		}
	}
	return cmd
}
