package handlers

import (
	"fmt"

	"axiom/src/interactions/grouptags"

	"github.com/bwmarrin/discordgo"
)

type funcICHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

var (
	// for on strings assigning to a func type
	commandInteractions = map[string]funcICHandler{
		"pts": grouptags.PtsCommandResponse,
	}
	messageComponentInteractions = map[string]funcICHandler{
		"selectGroupTag": grouptags.PtsGroupTagResponse,
		"newGroupTag":    grouptags.PtsGroupTagResponse,
		"delGroupTag":    grouptags.PtsGroupTagResponse,
		"createTag":      grouptags.PtsGroupTagResponse,
		"delTag":         grouptags.PtsGroupTagResponse,
		"ptsReturn":      grouptags.PtsCommandResponse,
		"leftPage":       grouptags.IncrementPage,
		"rightPage":      grouptags.DecrementPage,
	}
	submitModalInteractions = map[string]funcICHandler{
		"submitNewGroupTag": grouptags.SubmitNewGrouptag,
		"submitDelGroupTag": grouptags.DelGroupTag,
		"submitNewTag":      grouptags.SubmitNewTag,
		"submitDelTag":      grouptags.SubmitDelTag,
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
