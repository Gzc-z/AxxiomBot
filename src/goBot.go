// Package bot
package bot

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"axxiom/src/config"
	"axxiom/src/interactions"
	"axxiom/src/interactions/handlers"

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
	go bot.run()
}

func (bot *Bot) run() {
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

func (bot *Bot) Close() {
	ds := bot.Session
	defer ds.Close()

	guilds := listGuilds(ds)
	for _, guild := range guilds {
		bot.commandDelete(guild)
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
	guilds := listGuilds(s)
	for _, guild := range guilds {
		for _, v := range interactions.Commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, v)
			if err != nil {
				log.Println("Cannot create '%v' err: %v\n", v.Name, err)
			}
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

func listGuilds(s *discordgo.Session) []*discordgo.UserGuild {
	guilds, err := s.UserGuilds(100, "", "", false)
	if err != nil {
		log.Panicf("Cannot get guilds: %v", err)
	}

	var allowedGuilds []*discordgo.UserGuild
	roles := strings.Split(config.GetGuildID(), ",")

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
		if channel.Type == discordgo.ChannelTypeGuildText {
			_, err := s.ChannelMessageSend(channel.ID,
				"Olá! Serei seu bot de utilidades\n`.help`",
			)
			if err == nil {
				break
			}
		}
	}
}
