package slash

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"axxiom/internal/config"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/retry"
	"github.com/bwmarrin/discordgo"
)

var AIinstance *openrouter.OpenRouter

func getAIInstance() *openrouter.OpenRouter {
	retryConfig := retry.Config{
		Strategy: "backoff",
		Backoff: &retry.BackoffStrategy{
			InitialInterval: 1,
			MaxInterval:     30,
			Exponent:        1.1,
			MaxElapsedTime:  80,
		},
		RetryConnectionErrors: true,
	}

	AIinstance = openrouter.New(
		openrouter.WithRetryConfig(retryConfig),
		openrouter.WithSecurity(config.GetAIToken()),
		openrouter.WithTimeout(time.Second*20))

	return AIinstance
}

type Message struct {
	Content string
	Author  string
}

var (
	promptFile  = "internal/interactions/private/developer-prompt.txt"
	historyMsgs = make([]Message, 0, 2)
	once        sync.Once
)

func Axx(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	userMessage := i.ApplicationCommandData().Options[0].StringValue()
	if len(userMessage) > 200 {
		s.ChannelMessageSend(i.ChannelID, "mensagem muito grande")
		return fmt.Errorf("message too long")
	}
	s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("`%s`", userMessage))
	s.ChannelTyping(i.ChannelID)

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*10),
	)
	defer cancel()

	devPrompt, err := getDevPrompt()
	if err != nil {
		s.ChannelMessageSend(i.ChannelID, "Sem prompt de desenvolvedor")
	}
	messages := []components.ChatMessages{
		components.CreateChatMessagesDeveloper(
			components.ChatDeveloperMessage{
				Role:    components.ChatDeveloperMessageRoleDeveloper,
				Content: components.CreateChatDeveloperMessageContentStr(devPrompt),
			},
		),
		components.CreateChatMessagesUser(
			components.ChatUserMessage{
				Role:    components.ChatUserMessageRoleUser,
				Content: components.CreateChatUserMessageContentStr(userMessage),
			},
		),
	}

	if AIinstance == nil {
		AIinstance = getAIInstance()
	}
	res, err := AIinstance.Chat.Send(
		ctx,
		components.ChatRequest{
			Model:    openrouter.Pointer("@preset/default"),
			Messages: messages,
		},
		components.MetadataLevelEnabled.ToPointer(),
	)
	if err != nil {
		s.ChannelMessageSend(i.ChannelID, "erro ao gerar resposta: pode ser picos requisições")
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

func (h *Message) addToHistory(m components.ChatMessages) {
	historyMsgs = append(historyMsgs, *h)
}

func getDevPrompt() (string, error) {
	var loadErr error
	var devPrompt string
	once.Do(func() {
		f, err := os.ReadFile(promptFile)
		if err != nil {
			loadErr = err
			return
		}
		if string(f) == "" {
			loadErr = fmt.Errorf("empty prompt")
			return
		}
		devPrompt = string(f)
	})
	if loadErr != nil {
		return "", loadErr
	}
	return devPrompt, nil
}
