// Package bot
package bot

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"axxiom/internal/config"
	"axxiom/internal/interactions"
	"axxiom/internal/interactions/handlers"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	GuildID string
	AppID   string
}

type environments interface {
	getEnvironments(cfg config.Config)
}

func NewBot(cfg config.Config) *Bot {
	session, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Fatal("something's wrong, can't create bot")
		panic(1)
	}

	config.SetDeps(cfg)

	return &Bot{
		Session: session,
		GuildID: cfg.Discord.GuildID,
		AppID:   cfg.Discord.AppID,
	}
}

func (bot *Bot) Start() {
	ds := bot.Session
	ds.Identify.Intents |= discordgo.IntentMessageContent
	ds.Identify.Intents |= discordgo.IntentGuilds
	ds.Identify.Intents |= discordgo.IntentGuildMembers
	bot.AddSessionEvents()

	if err := ds.Open(); err != nil {
		panic(err)
	}

	// ds.AddHandler(onGuildCreate)

	userBot, _ := ds.User("@me")
	log.Printf("Logged in as: %v#%v", userBot.Username, userBot.Discriminator)
}

// TODO: scheduled events
func (bot Bot) AddSessionEvents() {
	ds := bot.Session
	handlers := []any{
		handlers.MessageCreate,
		handlers.InteractionCreate,
	}
	for _, handler := range handlers {
		go ds.AddHandler(handler)
	}
	go ds.AddHandler(bot.commandCreate)
}

func (bot Bot) Close() {
	ds := bot.Session
	defer ds.Close()

	guilds := bot.listGuilds(ds)
	for _, guild := range guilds {
		go bot.commandDelete(guild)
	}

	// delete global
	cmds, _ := ds.ApplicationCommands(bot.AppID, "")
	for _, v := range cmds {
		err := ds.ApplicationCommandDelete(bot.AppID, "", v.ID)
		if err != nil {
			log.Println(err)
		}
	}
}

func (bot Bot) commandCreate(s *discordgo.Session, r *discordgo.Ready) {
	guilds := bot.listGuilds(s)
	for _, guild := range guilds {
		for _, v := range interactions.Commands {
			go s.ApplicationCommandCreate(s.State.User.ID, guild.ID, v)
			fmt.Printf("/%s created\n", v.Name)
		}
	}
}

func (bot Bot) commandDelete(guild *discordgo.UserGuild) {
	s := bot.Session
	cmds, err := s.ApplicationCommands(bot.AppID, guild.ID)
	if err != nil {
		log.Println(err)
	}
	for _, v := range cmds {
		err := s.ApplicationCommandDelete(bot.AppID, guild.ID, v.ID)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("/%s deleted\n", v.Name)
	}
}

func (bot Bot) listGuilds(s *discordgo.Session) []*discordgo.UserGuild {
	guilds, err := s.UserGuilds(100, "", "", false)
	if err != nil {
		log.Panicf("Cannot get guilds: %v", err)
	}

	var allowedGuilds []*discordgo.UserGuild
	roles := strings.Split(bot.GuildID, ",")

	for _, guild := range guilds {
		if slices.Contains(roles, guild.ID) {
			allowedGuilds = append(allowedGuilds, guild)
		}
	}
	return allowedGuilds
}

func onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		return
	}

	for _, channel := range g.Channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			return
		}
		_, err := s.ChannelMessageSend(channel.ID,
			"Olá! Serei seu bot de utilidades\n`digite .help`",
		)
		if err == nil {
			break
		}
	}
}
