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

func NewBot(cfg config.Discord) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal("something's wrong, can't create bot")
		panic(1)
	}

	app, err := session.Application("@me")
	if err != nil {
		slog.Error("session application @me", "Error", err.Error())
		return nil, err
	}
	return &Bot{
		Session: session,
		GuildID: cfg.GuildID,
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
		ds.AddHandler(handler)
	}
	go ds.AddHandler(bot.commandCreate)
}

func (bot Bot) Close() {
	ds := bot.Session
	defer ds.Close()

	pg, _ := bot.listGuilds(ds)
	for _, guild := range pg {
		bot.commandDelete(guild.ID)
	}
	bot.commandDelete("") // delete global commands
}

func (bot Bot) commandCreate(s *discordgo.Session, r *discordgo.Ready) {
	commands := interactions.Commands()

	pg, prg := bot.listGuilds(s)
	for _, guild := range pg {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guild.ID, commands.Global)
		if err != nil {
			slog.Error("sync commands", "guild_id", guild.ID, "err", err)
			continue
		}
	}
	for _, guild := range prg {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guild.ID, commands.Private)
		if err != nil {
			slog.Error("sync commands", "guild_id", guild.ID, "err", err)
			continue
		}
	}
}

// for _, g := range bot.listGuilds(s) {
// 	if _, err := s.ApplicationCommandBulkOverwrite(appID, g.ID, interactions.Commands); err != nil {
// 	}
// 	slog.Info("commands synced", "guild_id", g.ID, "count", len(interactions.Commands))
// }

//
// for _, guild := range session.State.Guilds {
//     guildCommands := make([]*discordgo.ApplicationCommand, 0)
//
//     guildCommands = append(guildCommands, publicCommands...)
//
//     if isAuthorizedGuild(guild.ID) {
//         guildCommands = append(guildCommands, privateCommands...)
//     }
//
//     _, err := session.ApplicationCommandBulkOverwrite(
//         appID,
//         guild.ID,
//         guildCommands,
//     )
//
//     if err != nil {
//         log.Printf("Erro no servidor %s: %v", guild.ID, err)
//         continue
//     }
// }

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

// list guilds returns public and private guilds
func (bot Bot) listGuilds(s *discordgo.Session) ([]*discordgo.Guild, []*discordgo.Guild) {
	guilds := s.State.Guilds

	var allowedGuilds []*discordgo.Guild
	roles := strings.Split(bot.GuildID, ",")

	for _, guild := range guilds {
		if slices.Contains(roles, guild.ID) {
			allowedGuilds = append(allowedGuilds, guild)
		}
	}
	return guilds, allowedGuilds
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
