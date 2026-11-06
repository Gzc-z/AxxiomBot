// Package bot
package bot

import (
	"fmt"
	"log"
	"log/slog"
	"slices"
	"strings"

	"axxiom/internal/config"
	"axxiom/internal/handlers"
	"axxiom/internal/interactions"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	GuildID string
	AppID   string
}

func NewBot(cfg config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Fatal("something's wrong, can't create bot")
		panic(1)
	}

	config.SetDeps(cfg)

	app, err := session.Application("@me")
	if err != nil {
		slog.Error("session application @me", "Error", err.Error())
		return nil, err
	}
	return &Bot{
		Session: session,
		GuildID: cfg.Discord.GuildID,
		AppID:   app.ID,
	}, nil
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

	// TODO: verify which guilds have already been sent
	// ds.AddHandler(onGuildCreate)

	userBot, _ := ds.User("@me")
	slog.Info(fmt.Sprintf("Logged in as: %v#%v", userBot.Username, userBot.Discriminator))
}

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
		bot.commandDelete(guild.ID)
	}
	bot.commandDelete("") // delete global commands
}

func (bot Bot) commandCreate(s *discordgo.Session, r *discordgo.Ready) {
	guilds := bot.listGuilds(s)
	for _, guild := range guilds {
		for _, v := range interactions.Commands {
			go s.ApplicationCommandCreate(s.State.User.ID, guild.ID, v)
			slog.Info(fmt.Sprintf("%s command created: /%s\n", guild.ID, v.Name))
		}
	}
}

func (bot Bot) commandDelete(guildID string) {
	s := bot.Session
	cmds, err := s.ApplicationCommands(bot.AppID, guildID)
	if err != nil {
		slog.Error("Cannot delete command", "error", err.Error())
	}
	for _, v := range cmds {
		err := s.ApplicationCommandDelete(bot.AppID, guildID, v.ID)
		if err != nil {
			slog.Error("Cannot delete server command", "error", err.Error())
		}
		fmt.Printf("/%s deleted\n", v.Name)
	}
}

func (bot Bot) listGuilds(s *discordgo.Session) []*discordgo.UserGuild {
	guilds, err := s.UserGuilds(100, "", "", false)
	if err != nil {
		slog.Warn("Cannot get guilds", "error", err.Error())
		return nil
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
		_, err := s.ChannelMessageSend(
			channel.ID,
			"Olá! Serei seu bot de utilidades\n`digite .help`",
		)
		if err == nil {
			break
		}
	}
}
