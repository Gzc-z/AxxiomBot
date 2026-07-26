package grouptags

import (
	"encoding/json"
	"fmt"
	"os"

	"axiom/src/interactions"
	temputils "axiom/src/tempUtils"
	"axiom/src/ui"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
)

var (
	dirResp  string = "src/ui/"
	dirData  string = "src/data/"
	file     string
	response *discordgo.InteractionResponse

	groupID snowflake.ID
	page    uint8
)

// precisa de permanência temporária
func PtsCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ptsResponse := ui.PtsResponse()
	response := groupTagsUpdate(*ptsResponse)

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return err
	}
	return nil
}

func PtsGroupTagResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var err error
	var customID string

	data := i.MessageComponentData()
	if data.CustomID == "groupOptions" {
		customID = data.Values[0]
	} else {
		customID = data.CustomID
	}

	// TODO: control flow with map
	switch customID {
	case "newGroupTag":
		response = ui.GroupsTagResponse()
	case "delGroupTag":
		response = ui.DelModalGroup()
	case "selectGroupTag":
		response = selectGroupTag(data)
	case "createTag":
		response = ui.ModalTag()
	case "delTag":
		response = ui.DelModalTag()
	}
	if response == nil {
		return nil
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return err
	}
	return nil
}

func SubmitNewTag(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	gtCAP := 40
	file = dirData + "groupTags.json"

	gtValues := make([]*interactions.GroupTags, 0, gtCAP)
	err := json.Unmarshal(temputils.OpenGroupTag(), &gtValues)
	if err != nil || len(gtValues) >= gtCAP {
		return err
	}
	data := i.ModalSubmitData()

	inputs, err := temputils.GetInputs[interactions.Tag](data)
	if err != nil {
		return err
	}
	group := groupPagination(gtValues, 0)
	inputs.ID = group.Tags[len(group.Tags)-1].ID + 1
	group.Tags = append(group.Tags, inputs)

	response := ui.TagSelectMenuResponse(*group, page)
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return err
	}
	newContext, _ := json.MarshalIndent(gtValues, "", "	")
	err = os.WriteFile(file, newContext, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func SubmitDelTag(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var gtValues []*interactions.GroupTags
	data := i.ModalSubmitData()
	file = dirData + "groupTags.json"

	err := json.Unmarshal(temputils.OpenGroupTag(), &gtValues)
	if err != nil {
		return err
	}
	input, err := temputils.GetInputs[interactions.InternalUniqueValue](data)
	if err != nil {
		return err
	}
	// FIXIT: should select group tag before this
	for _, group := range gtValues {
		if group.ID == groupID {
			if len(group.Tags) > 128 { // set up std range
				fmt.Println("out of range: muitas tags")
				return nil
			}
			for i, tag := range group.Tags {
				if tag.ID == input.AsValue() {
					group.Tags = append(group.Tags[:i], group.Tags[i+1:]...)
					newContext, _ := json.MarshalIndent(gtValues, "", "	")
					os.WriteFile(file, newContext, 0o644)
				}
			}
			err = PtsCommandResponse(s, i)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func SubmitNewGrouptag(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	file = dirData + "groupTags.json"
	// groupTag defn and limit
	gtCAP := 40
	groupTags := make([]interactions.GroupTags, 0, gtCAP)

	err := json.Unmarshal(temputils.OpenGroupTag(), &groupTags)
	if err != nil || len(groupTags) >= gtCAP {
		return err
	}

	data := i.ModalSubmitData()
	inputs, err := temputils.GetInputs[interactions.GroupTags](data)
	if err != nil {
		return err
	}
	node, _ := snowflake.NewNode(1)
	inputs.ID = snowflake.ID(node.Generate().Int64())

	groupTags = append(groupTags, *inputs)
	newContext, _ := json.MarshalIndent(groupTags, "", "	")
	os.WriteFile(file, newContext, 0o644)

	// TODO: delete previous interaction
	err = PtsCommandResponse(s, i)
	if err != nil {
		return err
	}
	return nil
}

func DelGroupTag(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var gtValues []*interactions.GroupTags
	file = dirData + "groupTags.json"

	err := json.Unmarshal(temputils.OpenGroupTag(), &gtValues)
	if err != nil {
		return err
	}
	for idx, group := range gtValues {
		if group.ID == groupID {
			gtValues = append(gtValues[:idx], gtValues[idx+1:]...)
			newContext, _ := json.MarshalIndent(gtValues, "", "	")
			os.WriteFile(file, newContext, 0o644)

			ptsResponse := ui.PtsResponse()
			response = groupTagsUpdate(*ptsResponse)
		}
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		return err
	}
	return nil
}

func IncrementPage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if page > 10 {
		return nil
	}
	page += 1
	data := i.MessageComponentData()
	selectGroupTag(data)
	return nil
}

func DecrementPage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if page < 1 {
		return nil
	}
	page -= 1
	data := i.MessageComponentData()
	selectGroupTag(data)
	return nil
}

// TODO: pagination 2 bellow
func groupPagination(gtValues []*interactions.GroupTags, page int) *interactions.GroupTags {
	for _, group := range gtValues {
		if group.ID == groupID {
			if len(group.Tags) > 128 {
				fmt.Println("out of range: muitas tags")
				return nil
			}
			for i := 0; i < len(group.Tags); i += 5 {
				end := i + 5
				if end > len(group.Tags) {
					end = len(group.Tags)
				}
				group.Tags = group.Tags[i:end]
			}
			return group
		}
	}
	return nil
}

func selectGroupTag(data discordgo.MessageComponentInteractionData) *discordgo.InteractionResponse {
	if len(data.Values) == 0 {
		fmt.Println("no data found; returning nil")
		return nil
	}
	var groupTags []*interactions.GroupTags

	json.Unmarshal(temputils.OpenGroupTag(), &groupTags)

	for _, group := range groupTags {
		if group.ID.String() == data.Values[0] { // selectMenu Options
			// group := groupPagination(groupTags, 0)
			response = ui.TagSelectMenuResponse(*group, page)
			groupID = group.ID
		}
	}
	return response
}

// this one append groupTags from the groupTags.json file to ptsResponse.Components index 0
func groupTagsUpdate(slice discordgo.InteractionResponse) *discordgo.InteractionResponse {
	var groupTags []interactions.GroupTags

	json.Unmarshal(temputils.OpenGroupTag(), &groupTags)
	if len(groupTags) != 0 {
		gtSelectMenu := []discordgo.MessageComponent{groupTagSelectMenu(groupTags)}
		slice.Data.Components = append(gtSelectMenu, slice.Data.Components...)
	}
	return &slice
}

func groupTagSelectMenu(groups []interactions.GroupTags) discordgo.ActionsRow {
	var opts []discordgo.SelectMenuOption
	for _, opt := range groups {
		opts = append([]discordgo.SelectMenuOption{
			{
				Label:       opt.Name,
				Description: opt.Description,
				Value:       opt.ID.String(),
			},
		}, opts...)
	}
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    "selectGroupTag",
				Placeholder: "Selecionar grupo",
				Options:     opts,
				Disabled:    false,
			},
		},
	}
}
