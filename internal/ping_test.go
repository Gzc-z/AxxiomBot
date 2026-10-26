//go:build integration

package bot

import (
	"flag"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"axxiom/internal/config"
)

var discordMaxLen = 1000

func loadText() []string {
	args := flag.Args()
	if len(args) == 0 {
		return nil
	}
	if n := utf8.RuneCountInString(strings.Join(args, " ")); n > discordMaxLen {
		return nil
	}
	return args
}

// Integration test ping message
func TestPing(t *testing.T) {
	t.Parallel()

	cfg := config.Load()
	bot := NewBot(cfg)
	s := bot.Session

	channel := os.Getenv("TEST_CHANNELID")
	if channel == "" {
		t.Fail()
		return
	}

	text := loadText()
	if text == nil {
		text = []string{"ping"}
	}
	if _, err := s.ChannelMessageSend(channel, strings.Join(text, " ")); err != nil {
		t.Fatalf("failed to send message: %v", err)
	}
}
