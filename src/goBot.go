// Package bot
package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"axiom/src/config"
	"axiom/src/handlers"
	"axiom/src/interactions"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Bot struct {
	Session *discordgo.Session
	GuildID string
	AppID   string
}

func NewBot() *Bot {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token, exist := os.LookupEnv("DISCORD_BOT_TOKEN")
	if !exist {
		log.Fatal("Error loading .env file")
	}

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("something's wrong, can't create discord bot")
		os.Exit(1)
	}

	return &Bot{
		Session: sess,
		GuildID: config.GetGuildID(),
		AppID:   config.GetAppID(),
	}
}

func (bot *Bot) Start() {
	ds := bot.Session
	ds.Identify.Intents |= discordgo.IntentMessageContent
	ds.Identify.Intents |= discordgo.IntentGuilds
	ds.Identify.Intents |= discordgo.IntentGuildMembers

	bot.SessionEvents()

	if err := ds.Open(); err != nil {
		panic(err)
	}
	defer ds.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	userBot, _ := ds.User("@me")
	log.Printf("Logged in as: %v#%v", userBot.Username, userBot.Discriminator)
	<-sig
}

func (bot *Bot) Close() error {
	ds := bot.Session

	cmds, err := ds.ApplicationCommands(bot.AppID, bot.GuildID)
	if err != nil {
		log.Println(err)
	}
	if len(cmds) != 0 {
		for _, v := range cmds {
			err := ds.ApplicationCommandDelete(bot.AppID, bot.GuildID, v.ID)
			if err != nil {
				return err
			}
			fmt.Printf("\ncommand /%s deleted", v.Name)
		}
	}
	return nil
}

func (bot Bot) SessionEvents() {
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
	for _, v := range interactions.Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, bot.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		fmt.Printf("/%s created\n", v.Name)
	}
}
