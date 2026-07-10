// Package temputils is a temporarily
package temputils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"axiom/bot/interactions/models"

	"github.com/bwmarrin/discordgo"
)

type modelInterface interface {
	models.GroupTags | models.Tag
}

var (
	dirResp string = "bot/interactions/ui/"
	dirData string = "bot/interactions/data/"
	file    string
	err     error
)

func GetInputs[model modelInterface](data discordgo.ModalSubmitInteractionData) (*model, error) {
	keyInput := make(map[string]any)
	var in float32
	for _, row := range data.Components {
		for _, component := range row.(*discordgo.ActionsRow).Components {
			input, ok := component.(*discordgo.TextInput)
			if !ok {
				continue
			}
			// TODO: map inputValue -> decision structure
			if input.CustomID == "TagValue" {
				input.Value = strings.Replace(input.Value, ",", ".", 1)
				input.Value = strings.TrimSpace(input.Value)
				inputParsed, err := strconv.ParseFloat(input.Value, 32)
				if err != nil {
					return nil, err
				}
				in = float32(inputParsed)
				keyInput[input.CustomID] = in
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

func OpenGroupTag() []byte {
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
	return groupsContent
}
