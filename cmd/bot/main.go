package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	bot "axxiom/internal"
	"axxiom/internal/config"
)

func main() {
	cfg := config.Load()

	discordBot := bot.NewBot(cfg)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go discordBot.Start()

	<-ctx.Done()
	stop() // i realy don't know why this is needed

	discordBot.Close()
}
