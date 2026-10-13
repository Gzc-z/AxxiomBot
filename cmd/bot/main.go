package main

import (
	bot "axxiom/internal"
	"axxiom/internal/config"
)

func main() {
	cfg := config.Load()
	discordBot := bot.NewBot(cfg.Discord)

	discordBot.Start(cfg)
	defer discordBot.Close()
}
