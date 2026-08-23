package ui

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func randomColor() int {
	return rand.Intn(0xffffff)
}

func MembersResponse(v *discordgo.Member) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{
			{
				Color:       randomColor(),
				Type:        "rich",
				Title:       "User",
				Description: v.User.Username,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "",
				},
				Fields: []*discordgo.MessageEmbedField{
					// {
					// 	Name:   "ID",
					// 	Value:  v.User.ID,
					// 	Inline: false,
					// },
				},
				Image: &discordgo.MessageEmbedImage{
					URL: v.AvatarURL(""),
				},
			},
		},
	}
}

func UserResponse(v *discordgo.User) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{
			{
				Color: randomColor(),
				Type:  "rich",
				Title: "Me",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Nome global",
						Value:  v.GlobalName,
						Inline: false,
					},
				},
				Image: &discordgo.MessageEmbedImage{
					URL: v.AvatarURL(""),
				},
			},
		},
	}
}

func CommandsResponse(v *discordgo.ApplicationCommand) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: fmt.Sprintf("Nome: %s\nDescrição: %s\n%v",
			v.Name,
			v.Description,
			v.Version,
		),
	}
}
