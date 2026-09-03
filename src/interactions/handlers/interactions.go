package handlers

import (
	"fmt"

	"axxiom/src/interactions/slash"

	"github.com/bwmarrin/discordgo"
)

type funcICHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

var (
	// for on strings assigning to a func type
	commandInteractions = map[string]funcICHandler{
		"pts": slash.PtsCommandResponse,
		"axx": slash.Axx,
	}
	messageComponentInteractions = map[string]funcICHandler{
		"ptsReturn":      slash.PtsCommandResponse,
		"selectGroupTag": slash.PtsGroupTag,
		"newGroupTag":    slash.PtsGroupTag,
		"delGroupTag":    slash.PtsGroupTag,
		"createTag":      slash.PtsGroupTag,
		"delTag":         slash.PtsGroupTag,
		"rightPage":      slash.PtsGroupTag,
		"leftPage":       slash.PtsGroupTag,
	}
	submitModalInteractions = map[string]funcICHandler{
		"submitNewGroupTag": slash.SubmitNewGrouptag,
		"submitDelGroupTag": slash.DelGroupTag,
		"submitNewTag":      slash.SubmitNewTag,
		"submitDelTag":      slash.SubmitDelTag,
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
