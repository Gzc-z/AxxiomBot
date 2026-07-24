package handlers

import (
	"fmt"

	"axiom/bot/interactions"

	"github.com/bwmarrin/discordgo"
)

type funcICHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

var (
	// for on strings assigning to a func type
	commandInteractions = map[string]funcICHandler{
		"pts": interactions.PtsCommandResponse,
	}
	messageComponentInteractions = map[string]funcICHandler{
		"selectGroupTag": interactions.PtsGroupTagResponse,
		"newGroupTag":    interactions.PtsGroupTagResponse,
		"delGroupTag":    interactions.PtsGroupTagResponse,
		"createTag":      interactions.PtsGroupTagResponse,
		"delTag":         interactions.PtsGroupTagResponse,
		"ptsReturn":      interactions.PtsCommandResponse,
		"leftPage":       interactions.IncrementPage,
		"rightPage":      interactions.DecrementPage,
	}
	submitModalInteractions = map[string]funcICHandler{
		"submitNewGroupTag": interactions.SubmitNewGrouptag,
		"submitDelGroupTag": interactions.DelGroupTag,
		"submitNewTag":      interactions.SubmitNewTag,
		"submitDelTag":      interactions.SubmitDelTag,
	}
)

// TODO: maybe use generics
func interactionCreateErrors(handler funcICHandler, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := handler(s, i)
	// TODO: write errors in logs
	if err != nil {
		fmt.Println('\n', err)
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
		if data.CustomID == "groupOptions" {
			if handler, ok := messageComponentInteractions[data.Values[0]]; ok {
				go interactionCreateErrors(handler, s, i)
			}
			break
		}
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
