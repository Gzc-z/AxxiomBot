// Package temputils is a temporarily
package temputils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"axxiom/src/interactions"

	"github.com/bwmarrin/discordgo"
)

var (
	dirResp string = "src/ui/"
	dirData string = "src/data/"
	file    string
	err     error
)

type modelInterface interface {
	interactions.GroupTags | interactions.Tag | interactions.InternalUniqueValue
}

func GetInputs[model modelInterface](data discordgo.ModalSubmitInteractionData) (*model, error) {
	keyInput := make(map[string]any)
	for _, row := range data.Components {
		for _, component := range row.(*discordgo.ActionsRow).Components {
			input, ok := component.(*discordgo.TextInput)
			if !ok {
				fmt.Println("WARN: textInput possible error")
				continue
			}
			// TODO: map inputValue -> decision structure
			if input.CustomID == "tagSelected" {
				value, _ := strconv.Atoi(input.Value)
				keyInput[input.CustomID] = value
				continue
			}
			keyInput[input.CustomID] = input.Value
		}
	}
	m, _ := json.Marshal(keyInput)

	var values model
	err := json.Unmarshal(m, &values)
	if err != nil {
		return nil, err
	}
	return &values, nil
}

func OpenGroupTag() []*interactions.GroupTags {
	file = dirData + "groupTags.json"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.Mkdir(dirData, 0o755)
		os.WriteFile(file, []byte("[]"), 0o644)
		fmt.Println("file dont exist: add '[]' at the end of json")
	}

	groupsContent, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	if len(groupsContent) == 0 {
		os.WriteFile(file, []byte("[]"), 0o644)
		fmt.Println("add '[]' at the end of json")
		groupsContent, _ = os.ReadFile(file) // bad error handler
	}
	var groupTags []*interactions.GroupTags

	json.Unmarshal(groupsContent, &groupTags)
	return groupTags
}
