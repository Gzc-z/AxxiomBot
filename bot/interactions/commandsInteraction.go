package interactions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	dirResp string = "bot/interactions/responses/"
	dirData string = "bot/interactions/data/"
)

type iResponseDataStruct struct {
	Content    string                    `json:"content,omitempty"`
	Embeds     []*discordgo.MessageEmbed `json:"embeds"`
	Components []discordgo.ActionsRow    `json:"components"`
	CustomID   string                    `json:"customID"`
	Title      string                    `json:"title"`
}

type tag struct {
	User        discordgo.User
	Description string
	Value       int
}

type groupTags struct {
	Name        string
	Description string
	Tags        []*tag
}

func PtsCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content, err := os.ReadFile(dirResp + "pts.json")
	if err != nil {
		fmt.Println(err)
	}

	var data iResponseDataStruct
	err = json.Unmarshal(content, &data)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(groupTagcfg())
	var components []discordgo.MessageComponent
	for _, comp := range data.Components {
		components = append(components, comp)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    data.Content,
			Embeds:     data.Embeds,
			Components: components,
		},
	})
}

func PtsNewGroupTag(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	switch data.CustomID {
	case "newGroupTag":
	case "delGroupTag":
		fmt.Println(data.CustomID)
	default:
	}
}

func AxiomTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "ainda em desenvolvimento!!",
		},
	})
}
