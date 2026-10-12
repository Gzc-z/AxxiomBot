// Package config
package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	cfg  *Config
	once sync.Once
)

type Minecraft struct {
	IP    string
	PORT  string
	Forge string
}
type Discord struct {
	GuildID string
	AppID   string
	Token   string
}
type AI struct {
	AI_TOKEN string
}

type Config struct {
	Discord   Discord
	Minecraft Minecraft
	AI        AI
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(os.Getwd())
		log.Fatal("Error loading .env file\n", err)
	}

	cfg := Config{
		Discord: Discord{
			GuildID: getEnv("DISCORD_GUILD_ID"),
			AppID:   getEnv("APPLICATION_ID"),
			Token:   getEnv("DISCORD_BOT_TOKEN"),
		},
		Minecraft: Minecraft{
			PORT: getEnv("MINECRAFT_SERVER_PORT"),
			IP:   getEnv("MINECRAFT_SERVER_IP"),
		},
		AI: AI{
			AI_TOKEN: getEnv("AI_TOKEN"),
		},
	}

	return cfg
}

func getEnv(envVar string) string {
	key, exist := os.LookupEnv(envVar)
	if !exist {
		return fmt.Sprintf("Error: %s does'n exist\n", envVar)
	}
	return key
}

func GetGuildID() string {
	return cfg.Discord.GuildID
}

func GetMinecraft() Minecraft {
	return cfg.Minecraft
}

func GetAppID() string {
	return cfg.Discord.AppID
}

func GetAIKey() string {
	return cfg.AI.AI_TOKEN
}
