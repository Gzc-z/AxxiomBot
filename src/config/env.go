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
	IP   string
	PORT string
}
type Discord struct {
	GuildID string
	AppID   string
}

type Config struct {
	discord   Discord
	minecraft Minecraft
}

func load() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		cfg = &Config{
			discord: Discord{
				GuildID: getEnv("DISCORD_GUILD_ID"),
				AppID:   getEnv("APPLICATION_ID"),
			},
			minecraft: Minecraft{
				IP:   getEnv("MINECRAFT_SERVER_IP"),
				PORT: getEnv("MINECRAFT_SERVER_PORT"),
			},
		}
	})
	return cfg
}

func getEnv(envVar string) string {
	key, exist := os.LookupEnv(envVar)
	if !exist {
		return fmt.Sprintf("Error: %s", envVar)
	}
	return key
}

func GetGuildID() string {
	cfg := load()
	return cfg.discord.GuildID
}

func GetMinecraft() Minecraft {
	cfg := load()
	return cfg.minecraft
}

func GetAppID() string {
	cfg := load()
	return cfg.discord.AppID
}
