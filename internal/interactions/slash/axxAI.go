package slash

import (
	"context"
	"fmt"
	"time"

	"axxiom/internal/config"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/retry"
	"github.com/bwmarrin/discordgo"
)

// func preLoadAxx() {
// }
var AIinstance = openrouter.New(
	openrouter.WithRetryConfig(
		retry.Config{
			Strategy: "backoff",
			Backoff: &retry.BackoffStrategy{
				InitialInterval: 1,
				MaxInterval:     30,
				Exponent:        1.1,
				MaxElapsedTime:  100,
			},
			RetryConnectionErrors: true,
		}),
	openrouter.WithSecurity(config.GetAIKey()),
	openrouter.WithTimeout(time.Second*20),
)

func Axx(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	message := data.Options[0].StringValue() // option 0 from options of command
	if len(message) > 200 {
		s.ChannelMessageSend(i.ChannelID, "mensagem muito grande")
		return fmt.Errorf("message too long")
	}
	s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("`%s`", message))

	s.ChannelTyping(i.ChannelID)
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*10),
	)
	defer cancel()

	userMessage := &components.ChatUserMessage{
		Role: "user",
		Content: components.ChatUserMessageContent{
			Str: openrouter.Pointer(message),
		},
	}

	res, err := AIinstance.Chat.Send(ctx, components.ChatRequest{
		Model: openrouter.Pointer("@preset/default"),
		Messages: []components.ChatMessages{{
			ChatUserMessage: userMessage,
		}},
	}, components.MetadataLevelEnabled.ToPointer())
	if err != nil {
		s.ChannelMessageSend(i.ChannelID, "erro ao gerar resposta\npode ser picos requisições")
		return err
	}

	resp, ok := res.ChatResult.Choices[0].Message.Content.Get()
	if !ok {
		s.ChannelMessageSend(i.ChannelID, "erro ao gerar resposta\nerro: 500")
		return fmt.Errorf("error getting response")
	}
	s.ChannelMessageSend(i.ChannelID, *resp.Str)
	return nil
}
