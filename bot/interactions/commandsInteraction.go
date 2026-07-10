package interactions

import (
	"encoding/json"
	"fmt"
	"os"

	"axiom/bot/interactions/models"
	"axiom/bot/interactions/ui"
	temputils "axiom/bot/tempUtils"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
)

var (
	dirResp string = "bot/interactions/ui/"
	dirData string = "bot/interactions/data/"
	file    string

	// Data = make(map[UserKey]CommandsValue)
	// shr  = func(comp *[]discordgo.MessageComponent, i *discordgo.InteractionCreate) {
	// 	key := UserKey{
	// 		User:    i.User,
	// 		GuidID:  config.GetGuildID(),
	// 		Channel: i.ChannelID,
	// 	}
	// 	value := CommandsValue{
	// 		Interaction: i.Interaction,
	// 		// Components:  comp,
	// 	}
	// 	Data = map[UserKey]CommandsValue{
	// 		key: value,
	// 	}
	// }
	// TEMP
	groupID snowflake.ID
)

type UserKey struct {
	User    *discordgo.User
	GuidID  string
	Channel string
}
type CommandsValue struct {
	Interaction *discordgo.Interaction
	// Components  *[]discordgo.MessageComponent
}

// INFO: other handlers needs to call ApplicationsCommand /pts to read Interaction
// infos precisa de permanência temporária
func PtsCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	response := ui.PtsResponse()
	ptsResponse := groupTagsUpdate(*response)

	// go shr(&ptsResponse.Components, i)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    ptsResponse.Content,
			Embeds:     ptsResponse.Embeds,
			Components: ptsResponse.Components,
			Title:      ptsResponse.Title,
			// Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func PtsGroupTagResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var err error
	var response *discordgo.InteractionResponse

	data := i.MessageComponentData()

	// TODO: control flow with map
	switch data.CustomID {
	case "newGroupTag":
		response = ui.GroupsTagResponse()
	case "selectGroupTag":
		if len(data.Values) == 0 {
			fmt.Println("no data found; returning nil")
			return nil
		}
		var groupTags []models.GroupTags

		json.Unmarshal(temputils.OpenGroupTag(), &groupTags)

		for i, r := range groupTags {
			if r.ID.String() == data.Values[0] { // selectMenu Options
				response = ui.TagSelectMenuResponse(groupTags[i])
				groupID = r.ID
			}
		}
		if response == nil {
			fmt.Println("no groupTags.ID found; returning nil")
			return nil
		}
	case "createTag":
		response = ui.ModalTag()
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

	gtValues := make([]*models.GroupTags, 0, gtCAP)
	err := json.Unmarshal(temputils.OpenGroupTag(), &gtValues)
	if err != nil || len(gtValues) >= gtCAP {
		return err
	}
	data := i.ModalSubmitData()

	inputs, err := temputils.GetInputs[models.Tag](data)
	if err != nil {
		return err
	}
	for _, group := range gtValues {
		if group.ID == groupID {
			if len(group.Tags) > 128 {
				fmt.Println("out of range: muitas tags")
				return nil
			}
			inputs.ID = len(group.Tags) + 1
			group.Tags = append(group.Tags, inputs)
		}
	}
	newContext, _ := json.MarshalIndent(gtValues, "", "	")
	err = os.WriteFile(file, newContext, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func SubmitNewGrouptag(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	file = dirData + "groupTags.json"
	// groupTag defn and limit
	gtCAP := 40
	groupTags := make([]models.GroupTags, 0, gtCAP)

	err := json.Unmarshal(temputils.OpenGroupTag(), &groupTags)
	if err != nil || len(groupTags) >= gtCAP {
		return err
	}

	data := i.ModalSubmitData()
	inputs, err := temputils.GetInputs[models.GroupTags](data)
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

// this block append groupTags from the groupTags.json file to ptsResponse.Components index 0
func groupTagsUpdate(slice discordgo.InteractionResponseData) discordgo.InteractionResponseData {
	var groupTags []models.GroupTags

	json.Unmarshal(temputils.OpenGroupTag(), &groupTags)
	if len(groupTags) != 0 {
		gtSelectMenu := []discordgo.MessageComponent{groupTagSelectMenu(groupTags)}
		slice.Components = append(gtSelectMenu, slice.Components...)
	}
	return slice
}

func groupTagSelectMenu(groups []models.GroupTags) discordgo.ActionsRow {
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
