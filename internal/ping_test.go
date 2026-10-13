package bot

import (
	"fmt"
	"os"
	"testing"

	"axxiom/internal/config"
)

// Integration test ping message
func TestPing(t *testing.T) {
	cfg := config.Load()
	bot := NewBot(cfg.Discord)
	config.SetDeps(cfg) // improve injection of dependencies... pls

	fmt.Println("logged")
	s := bot.Session

	channel := os.Getenv("TEST_CHANNELID")
	_, err := s.ChannelMessageSend(channel, "ping")
	t.Log(err)
}
