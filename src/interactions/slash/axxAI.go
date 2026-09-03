package slash

import (
	"context"
	"fmt"

	"axxiom/src/config"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/genai"
)

func Axx(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	message := data.Options[0].StringValue()

	query := fmt.Sprintf("`%s`\nAxx:", message)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: query,
		},
	})
	s.ChannelTyping(i.ChannelID)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.GetAI(),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}

	genconfig := &genai.GenerateContentConfig{
		Temperature:     genai.Ptr(float32(0.7)),
		MaxOutputTokens: 400,
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: `
				Você se chama Axxiom.
				Você é um assistente especialista jogos, minecraft, configurações, e muito mais.
				Você responde como um bot de ajuda, e não como um humano.
				Seu objetivo é ajudar o usuário a resolver suas dúvidas de forma mais eficiente, rápida, precisa e com poucas palavras.
				Sempre responda de forma clara, direta
				`},
			},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3.6-flash",
		genai.Text(message),
		genconfig,
	)
	if err != nil {
		return err
	}

	if result.UsageMetadata != nil {
		fmt.Printf("Prompt tokens: %d\n", result.UsageMetadata.PromptTokenCount)
		fmt.Printf("Response tokens: %d\n", result.UsageMetadata.CandidatesTokenCount)
		fmt.Printf("Total tokens: %d\n", result.UsageMetadata.TotalTokenCount)
	}

	s.ChannelMessageSend(i.ChannelID, result.Text())
	return nil
}
