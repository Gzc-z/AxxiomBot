package handlers

import (
	"fmt"
	"os"

	"axiom/bot/interactions"

	"github.com/bwmarrin/discordgo"
)

type funcICHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

var (
	commandInteractions = map[string]funcICHandler{
		"pts": interactions.PtsCommandResponse,
	}
	messageComponentInteractions = map[string]funcICHandler{
		"newGroupTag":    interactions.PtsGroupTagResponse,
		"selectGroupTag": interactions.PtsGroupTagResponse,
	}
	submitModalInteractions = map[string]funcICHandler{
		"submitNewGroupTag": interactions.SubmitNewGrouptag,
	}
)

func interactionCreateErrors(handler funcICHandler, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := handler(s, i)
	if err != nil {
		fmt.Println('\n', err)
		os.Exit(1)
	}
}

// TODO: implement fallback response
func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		if handler, ok := commandInteractions[data.Name]; ok {
			go interactionCreateErrors(handler, s, i)
		}
	case discordgo.InteractionMessageComponent:
		data := i.MessageComponentData()
		if handler, ok := messageComponentInteractions[data.CustomID]; ok {
			go interactionCreateErrors(handler, s, i)
		}
	case discordgo.InteractionModalSubmit:
		data := i.ModalSubmitData()
		if handler, ok := submitModalInteractions[data.CustomID]; ok {
			go interactionCreateErrors(handler, s, i)
		}
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
