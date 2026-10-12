package bot

import (
	"fmt"
	"strings"
	"testing"

	"axxiom/internal/config"
)

func TestPing(t *testing.T) {
	cfg := config.Load()
	bot := NewBot(cfg)

	fmt.Println("logged")
	s := bot.Session

	guilds := strings.Split(config.GetGuildID(), ",")
	s.ChannelMessageSend(guilds[0], "ping")
}
