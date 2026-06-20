package interactions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	dirResp string = "bot/interactions/responses/"
	dirData string = "bot/interactions/data/"
	file    string
)

type responseDataStruct struct {
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
	Description string `json:"description,omitempty"`
	Tags        []*tag `json:"tags,omitempty"`
}

func PtsCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	groupTagsFile := dirData + "groupTags.json"
	groupContent, err := os.Open(groupTagsFile)
	defer groupContent.Close()
	if err != nil {
		fmt.Println(err)
	}
	file = dirResp + "ptsResponse.json"
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	var data responseDataStruct
	var components []discordgo.MessageComponent
	var groupTags []groupTags

	json.Unmarshal(content, &data)
	json.NewDecoder(groupContent).Decode(groupTags)
	// PERF: perf issue flag:  nested code and bad structure
	for _, row := range data.Components {
		for _, comp := range row.Components {
			if len(groupTags) == 0 {
				if comp.Type() == 0x3 { // discordgo textinput id
					continue
				}
			}
			components = append(components, row)
			break
		}
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    data.Content,
			Embeds:     data.Embeds,
			Components: components,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func PtsGroupTagResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	file = dirResp + "groupTagsResponse.json"
	content, err := os.ReadFile(file)
	if err != nil {
		log.Println(err)
	}

	var dataStruct responseDataStruct
	json.Unmarshal(content, &dataStruct)

	var components []discordgo.MessageComponent
	for _, row := range dataStruct.Components {
		components = append(components, row)
	}

	data := i.MessageComponentData()
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
	return nil
}

// TODO: make this func polymorphic
func getInputValues(ch chan<- *discordgo.TextInput, wg *sync.WaitGroup, i discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	for _, row := range data.Components {
		for _, component := range row.(*discordgo.ActionsRow).Components {
			input, ok := component.(*discordgo.TextInput)
			if !ok {
				continue
			}
			defer wg.Done()
			ch <- input
		}
	}
}

func SubmitNewGrouptag(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var wg sync.WaitGroup
	var input discordgo.TextInput
	var values *[]groupTags
	// data := i.ModalSubmitData()

	file = dirData + "groupTags.json"
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	// groupTag defn and limit
	gtValues := make([]groupTags, 40)
	err = json.Unmarshal(content, &gtValues)
	if err != nil {
		fmt.Println(err)
		return err
	}

	wg.Add(3)
	ch := make(chan *discordgo.TextInput, 3)

	go func() {
		wg.Wait()
		close(ch)
	}()

	go func() {
		getInputValues(ch, &wg, *i)
	}()

	// NOTE: n precisa disso.. e sim de uma struct
	keyInput := map[string]groupTags{
		"Name":        {Name: input.Value},
		"Description": {Description: input.Value},
	}
	for input := range ch {
		defer wg.Done()

		*values = append(*values, keyInput[input.CustomID])
		// fmt.Println(input.CustomID)
		// fmt.Println(keyInput[input.CustomID])
	}
	newContext, _ := json.MarshalIndent(gtValues, "", "	")
	fmt.Println(string(newContext))

	// os.WriteFile(file, newContext, 0o644)
	return nil
}
