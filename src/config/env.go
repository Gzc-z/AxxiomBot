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
	Discord   Discord
	Minecraft Minecraft
}

func load() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		cfg = &Config{
			Discord: Discord{
				GuildID: getEnv("DISCORD_GUILD_ID"),
				AppID:   getEnv("APPLICATION_ID"),
			},
			Minecraft: Minecraft{
				PORT: getEnv("MINECRAFT_SERVER_PORT"),
				IP:   getEnv("MINECRAFT_SERVER_IP"),
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
	return cfg.Discord.GuildID
}

func GetMinecraft() Minecraft {
	cfg := load()
	return cfg.Minecraft
}

func GetAppID() string {
	cfg := load()
	return cfg.Discord.AppID
}
