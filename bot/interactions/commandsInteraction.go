package interactions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"axiom/bot/config"
	"axiom/bot/interactions/ui"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
)

var (
	dirResp string = "bot/interactions/ui/"
	dirData string = "bot/interactions/data/"
	file    string

	Data = make(map[UserKey]CommandsValueCtx)
	shr  = func(comp *[]discordgo.MessageComponent, i *discordgo.InteractionCreate) {
		key := UserKey{
			User:    i.User,
			GuidID:  config.GetGuildID(),
			Channel: i.ChannelID,
		}
		value := CommandsValueCtx{
			Interaction: i.Interaction,
			// Components:  comp,
		}
		Data = map[UserKey]CommandsValueCtx{
			key: value,
		}
	}
)

type UserKey struct {
	User    *discordgo.User
	GuidID  string
	Channel string
}
type CommandsValueCtx struct {
	Interaction *discordgo.Interaction
	// Components  *[]discordgo.MessageComponent
}

type tag struct {
	User        discordgo.User
	Description string
	Value       int
}
type groupTags struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Tags        []*tag       `json:"tags,omitempty"`
	ID          snowflake.ID `json:"ID"`
}

func PtsCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	response := ui.PtsResponse()
	ptsResponse := groupTagsUpdate(*response)

	go shr(&ptsResponse.Components, i)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    ptsResponse.Content,
			Embeds:     ptsResponse.Embeds,
			Components: ptsResponse.Components,
			Title:      ptsResponse.Title,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func PtsGroupTagResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var err error

	// INFO: has to call ApplicationsCommand /pts to read infos
	// FIXIT:infos precisa de permanência temporária
	if len(Data) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "use o comando /pts primeiro\nFale pro Gzc arrumar sa porra aqui",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		time.Sleep(time.Second * 6)
		s.InteractionResponseDelete(i.Interaction)
		return nil //
	}

	data := i.MessageComponentData()
	switch data.CustomID {
	case "newGroupTag":
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: ui.GroupsTagResponse(),
		})
	case "selectGroupTag":
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: ui.TagSelectMenuResponse(),
				Flags:      discordgo.MessageFlagsIsComponentsV2,
			},
		})
	}
	if err != nil {
		log.Println(err)
	}
	return nil
}

func SubmitNewGrouptag(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	gtCAP := 40
	file = dirData + "groupTags.json"
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// groupTag defn and limit
	gtValues := make([]groupTags, 0, gtCAP)
	err = json.Unmarshal(content, &gtValues)
	if err != nil || len(gtValues) >= gtCAP {
		return err
	}

	// TODO: send inputs then wait for response to write in file
	inputs := getInputValues(*i)
	gtValues = append(gtValues, inputs)
	newContext, _ := json.MarshalIndent(gtValues, "", "	")
	os.WriteFile(file, newContext, 0o644)

	key := UserKey{
		User:    i.User,
		GuidID:  config.GetGuildID(),
		Channel: i.ChannelID,
	}
	data := Data[key]

	response := ui.PtsResponse()
	ptsResponse := groupTagsUpdate(*response)

	_, err = s.InteractionResponseEdit(data.Interaction, &discordgo.WebhookEdit{
		Components: &ptsResponse.Components,
	})
	if err != nil {
		return err
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	s.InteractionResponseDelete(i.Interaction)
	return nil
}

// this block append groupTags from the groupTags.json file to ptsResponse.Components index 0
func groupTagsUpdate(slice discordgo.InteractionResponseData) discordgo.InteractionResponseData {
	var groupTags []groupTags

	file = dirData + "groupTags.json"
	groupsContent, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(groupsContent, &groupTags)

	if len(groupTags) != 0 {
		gtSelectMenu := []discordgo.MessageComponent{groupTagSelectMenu(groupTags)}
		slice.Components = append(gtSelectMenu, slice.Components...)
	}
	return slice
}

// TODO: make this func polymorphic
// WARN:: nedted code here
func getInputValues(i discordgo.InteractionCreate) groupTags {
	data := i.ModalSubmitData()
	node, _ := snowflake.NewNode(1)

	keyInput := make(map[string]any)
	for _, row := range data.Components {
		for _, component := range row.(*discordgo.ActionsRow).Components {
			input, ok := component.(*discordgo.TextInput)
			if !ok {
				continue
			}
			keyInput[input.CustomID] = input.Value
		}
	}
	m, _ := json.Marshal(keyInput)

	var values groupTags
	json.Unmarshal(m, &values)
	values.ID = snowflake.ID(node.Generate().Int64())

	return values
}

func groupTagSelectMenu(groups []groupTags) discordgo.ActionsRow {
	var opts []discordgo.SelectMenuOption
	for _, opt := range groups {
		opts = append(opts, discordgo.SelectMenuOption{
			Label:       opt.Name,
			Description: opt.Description,
			Value:       opt.ID.String(),
		})
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
