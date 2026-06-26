package ui

import (
	componentsv2 "github.com/OGCraft-Eu/discordgo-componentsv2"
	"github.com/bwmarrin/discordgo"
)

func PtsResponse() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		CustomID: "pts_command",
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:        "rich",
				Title:       "pts command",
				Description: "## grupo de informações",
				Color:       16777215,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Criar grupo",
						Value:  "Cria um novo grupo :)",
						Inline: false,
					},
					{
						Name:   "Selecionar grupo",
						Value:  "seleção e configuração do grupo\naparece somente se tiver um grupo",
						Inline: false,
					},
					{
						Name:   "Deletar grupo",
						Value:  "deleta obviamente o grupo",
						Inline: false,
					},
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "newGroupTag",
						Style:    3,
						Label:    "Criar Grupo",
					},
				},
			},
		},
	}
}

func GroupsTagResponse() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Content:  "",
		CustomID: "submitNewGroupTag",
		Title:    "Criando novo grupo",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "Name",
						Label:       "Nome",
						Placeholder: "Name",
						Style:       1,
						Required:    true,
						MinLength:   0,
						MaxLength:   45,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "Description",
						Label:       "Descrição",
						Placeholder: "Description",
						Style:       2,
						Required:    false,
						MinLength:   0,
						MaxLength:   145,
					},
				},
			},
		},
	}
}

func TagSelectMenuResponse() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		componentsv2.NewContainerBuilder().
			SetAccentColor(0x00ff00).
			AddComponent(
				componentsv2.NewTextDisplayBuilder().
					SetContent("**Lorem ipsum sit dolor**").
					Build(),
			).
			AddComponent(
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "adicionar",
							Style:    discordgo.SecondaryButton,
							CustomID: "add",
							Emoji:    &discordgo.ComponentEmoji{Name: "🎵 "},
						},
					},
				},
			).Build(),
	}
}
