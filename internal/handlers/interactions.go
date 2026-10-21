// Package handlers
//
// handler interactions like create messages
// and redirect to your own functionality
package handlers

import (
	"log/slog"
	"runtime/debug"

	"axxiom/internal/interactions/slash"

	"github.com/bwmarrin/discordgo"
)

type icHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

var (
	commandInteractions = map[string]icHandler{
		"pts": slash.PtsCommandResponse,
		"axx": slash.Axx,
	}
	messageComponentInteractions = map[string]icHandler{
		"ptsReturn":      slash.PtsCommandResponse,
		"selectGroupTag": slash.PtsGroupTag,
		"newGroupTag":    slash.PtsGroupTag,
		"delGroupTag":    slash.PtsGroupTag,
		"createTag":      slash.PtsGroupTag,
		"delTag":         slash.PtsGroupTag,
		"rightPage":      slash.PtsGroupTag,
		"leftPage":       slash.PtsGroupTag,
	}
	submitModalInteractions = map[string]icHandler{
		"submitNewGroupTag": slash.SubmitNewGrouptag,
		"submitDelGroupTag": slash.DelGroupTag,
		"submitNewTag":      slash.SubmitNewTag,
		"submitDelTag":      slash.SubmitDelTag,
	}
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var h icHandler
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h = commandInteractions[i.ApplicationCommandData().Name]
	case discordgo.InteractionMessageComponent:
		d := i.MessageComponentData()
		key := d.CustomID
		if key == "groupOptions" || len(d.Values) > 0 {
			key = d.Values[0]
		}
		h = messageComponentInteractions[key]
	case discordgo.InteractionModalSubmit:
		h = submitModalInteractions[i.ModalSubmitData().CustomID]
	}
	if h != nil {
		runWithRecover(h, s, i)
	}
}

func runWithRecover(h icHandler, s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("handler panicked", "interaction", i.ID, "type", i.Type, "panic", err, "stack", string(debug.Stack()))
			reply(s, i.Interaction)
		}
	}()
	if err := h(s, i); err != nil {
		slog.Error("handler failed", "interaction", i.ID, "type", i.Type, "err", err)
		reply(s, i.Interaction)
	}
}

func reply(s *discordgo.Session, in *discordgo.Interaction) {
	const msg = "Algo deu errado. Por favor, tente novamente mais tarde."
	err := s.InteractionRespond(in, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg, Flags: discordgo.MessageFlagsEphemeral},
	})
	if err == nil {
		return
	}
	// Followup message failed
	_, ferr := s.FollowupMessageCreate(in, true, &discordgo.WebhookParams{
		Content: msg,
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	if ferr != nil {
		slog.Error("failed to notify user", "interaction", in.ID, "err", err, "followup_err", ferr)
	}
}
