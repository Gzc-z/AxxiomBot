package interactions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	dirResp string = "bot/interactions/responses/"
	dirData string = "bot/interactions/data/"
	file    string
)

type iResponseDataStruct struct {
	Content    string                    `json:"content,omitempty"`
	Embeds     []*discordgo.MessageEmbed `json:"embeds"`
	Components []discordgo.ActionsRow    `json:"components"`
	CustomID   string                    `json:"custom_id"`
	Title      string                    `json:"title"`
}

type tag struct {
	User        discordgo.User
	Description string
	Value       int
}

type groupTags struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []*tag `json:"tags,omitempty"`
}

func PtsCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content, err := os.ReadFile(dirResp + "ptsResponse.json")
	if err != nil {
		fmt.Println(err)
	}

	var data iResponseDataStruct
	json.Unmarshal(content, &data)

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
	file = dirResp + "groupTagsResponse.json"
	content, err := os.ReadFile(file)
	if err != nil {
		log.Println(err)
	}
	var dataStruct iResponseDataStruct
	json.Unmarshal(content, &dataStruct)

	var components []discordgo.MessageComponent
	for _, comp := range dataStruct.Components {
		components = append(components, comp)
	}

	switch data.CustomID {
	case "newGroupTag":
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID:   dataStruct.CustomID,
				Components: components,
				Content:    dataStruct.Content,
				Title:      dataStruct.Title,
			},
		})
	case "delGroupTag":
		fmt.Println(data.CustomID)
	}
	if err != nil {
		log.Println(err)
	}
}

// TODO: make this func polymorphic
func attrInputValues(values *[]groupTags, chanInput chan discordgo.TextInput) *[]groupTags {
	input := <-chanInput
	keyInput := map[string]groupTags{
		"Name":        {Name: input.Value},
		"Description": {Description: input.Value},
	}
	*values = append(*values, keyInput[input.CustomID])
	return values
}

func SubmitNewGrouptag(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	file = dirData + "groupTags.json"
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	values := make([]groupTags, 0, 2)
	json.Unmarshal(content, &values)
	ch := make(chan discordgo.TextInput, 2)
	for _, row := range data.Components {
		for _, component := range row.(*discordgo.ActionsRow).Components {
			if input, ok := component.(discordgo.TextInput); ok {
				ch <- input
			}
		}
		attrInputValues(&values, ch)
	}
	// for _, i := range values {
	// 	fmt.Println(i.Name)
	// }
	newContext, _ := json.MarshalIndent(values, "", "	")
	fmt.Println(string(newContext))

	// os.WriteFile(file, newContext, 0o644)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "criado com sucesso",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func AxiomTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "ainda em desenvolvimento!!",
		},
	})
}
