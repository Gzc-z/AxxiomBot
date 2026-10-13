// Package config
package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
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
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Failed to get caller")
	}

	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, "../../.env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Println(err)
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

func SetDeps(deps Config) {
	cfg = &deps
}

func MinecraftServer() Minecraft {
	return cfg.Minecraft
}

func MinecraftAddress() string {
	return net.JoinHostPort(cfg.Minecraft.IP, cfg.Minecraft.PORT)
}

func GetGuildID() string {
	return cfg.Discord.GuildID
}

func GetAppID() string {
	return cfg.Discord.AppID
}

func GetAIToken() string {
	return cfg.AI.AI_TOKEN
}
