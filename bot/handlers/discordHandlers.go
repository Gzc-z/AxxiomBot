package handlers

import (
	"axiom/bot/interactions"

	"github.com/bwmarrin/discordgo"
)

type funcHandler func(*discordgo.Session, *discordgo.InteractionCreate)

var (
	commandInteractions = map[string]funcHandler{
		"axiom-test": interactions.AxiomTest,
		"pts":        interactions.PtsCommandResponse,
	}
	messageComponentInteractions = map[string]funcHandler{
		"newGroupTag": interactions.PtsNewGroupTag,
	}
	// submitModalInteraction
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		if handler, ok := commandInteractions[data.Name]; ok {
			handler(s, i)
		}
	case discordgo.InteractionMessageComponent:
		data := i.MessageComponentData()
		if handler, ok := messageComponentInteractions[data.CustomID]; ok {
			handler(s, i)
		}
	case discordgo.InteractionModalSubmit:
		// data := i.ModalSubmitData()
	}
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "hello" {
		s.ChannelMessageSendReply(m.ChannelID, "world!", m.Reference())
	}
}
