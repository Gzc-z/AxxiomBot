// Package config
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token   string
	GuildID string
}

func GetGuildID() string {
	cfg := Load()
	return cfg.GuildID
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		Token:   getEnv("DISCORD_BOT_TOKEN"),
		GuildID: getEnv("DISCORD_GUILD_ID"),
	}
}

func getEnv(envVar string) string {
	key, exist := os.LookupEnv(envVar)
	if !exist {
		return fmt.Sprintln(envVar, "key: dt exist")
	}
	return key
}
