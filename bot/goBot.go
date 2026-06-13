// Package bot
package bot

import (
	"fmt"
	"log"
	"os"

	"axiom/bot/config"
	"axiom/bot/handlers"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	GuildID string
}

func NewBot(cfg config.Config) *Bot {
	bot, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal("something's wrong, can't create discord bot")
		os.Exit(1)
	}

	return &Bot{
		Session: bot,
		GuildID: cfg.GuildID,
	}
}

func (bot *Bot) SessionEvents() {
	ds := bot.Session
	handlers := []any{
		handlers.MessageCreate,
		handlers.InteractionCreate,
	}
	for _, handler := range handlers {
		go ds.AddHandler(handler)
	}
	go ds.AddHandler(bot.applicationCommandCreate)
}

func (bot Bot) applicationCommandCreate(s *discordgo.Session, r *discordgo.Ready) {
	for _, v := range handlers.Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, bot.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		fmt.Printf("/%s created\n", v.Name)
	}
}
