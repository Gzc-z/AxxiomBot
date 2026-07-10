package ui

import (
	"fmt"

	"axiom/bot/interactions/models"

	"github.com/bwmarrin/discordgo"
)

type ptr interface {
	int | bool | discordgo.SeparatorSpacingSize
}

func pointer[pointer ptr](i pointer) *pointer {
	return &i
}

var (
	largeSpace = discordgo.SeparatorSpacingSizeLarge
	smallSpace = discordgo.SeparatorSpacingSizeSmall
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

func GroupsTagResponse() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Content:  "",
			CustomID: "submitNewGroupTag",
			Title:    "Criando novo grupo",
			Components: []discordgo.MessageComponent{
				newActionRow(
					discordgo.TextInput{
						CustomID:    "name",
						Label:       "Nome",
						Placeholder: "Name",
						Style:       1,
						Required:    true,
						MinLength:   0,
						MaxLength:   45,
					},
				),
				newActionRow(
					discordgo.TextInput{
						CustomID:    "description",
						Label:       "Descrição",
						Placeholder: "Description",
						Style:       2,
						Required:    false,
						MinLength:   0,
						MaxLength:   145,
					},
				),
			},
		},
	}
}

func TagSelectMenuResponse(group models.GroupTags) *discordgo.InteractionResponse {
	comp := []discordgo.MessageComponent{
		newContainer(
			pointer(0x4c4bff), // AccentColor
			textDisplay(fmt.Sprintf("**Editando:** `#%s`", group.Name)),
			separator(true, smallSpace),

			textDisplay("**Descrição**"),
			textDisplay("> "+group.Description),

			separator(true, largeSpace),
			newActionRow(
				discordgo.Button{
					Label:    "Criar Tag",
					Style:    discordgo.SuccessButton,
					CustomID: "createTag",
					Emoji:    &discordgo.ComponentEmoji{Name: "➕"},
				},
				// discordgo.Button{
				// 	Label:    "editar",
				// 	Style:    discordgo.PrimaryButton,
				// 	CustomID: "manageTags",
				// },
				// discordgo.Button{
				// 	Label:    "Grupo",
				// 	Style:    discordgo.SecondaryButton,
				// 	CustomID: "manageGroup",
				// },
			),
		),
	}
	if len(group.Tags) > 0 {
		comp = append(comp,
			newContainer(
				pointer(0x4c4bff), // AccentColor
				displayListModel(group.Tags)...,
			),
		)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: comp,
		},
	}
}

func ModalTag() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Content:  "",
			CustomID: "submitNewTag",
			Title:    "Criando nova Tag",
			Components: []discordgo.MessageComponent{
				newActionRow(
					discordgo.TextInput{
						CustomID:    "tagName",
						Label:       "Nome",
						Placeholder: "Tag",
						Style:       1,
						Required:    true,
						MinLength:   0,
						MaxLength:   45,
					},
				),
				newActionRow(
					discordgo.TextInput{
						CustomID:    "tagDescription",
						Label:       "Descrição",
						Placeholder: "Description",
						Style:       2,
						Required:    false,
						MinLength:   0,
						MaxLength:   145,
					},
				),
				newActionRow(
					discordgo.TextInput{
						CustomID:    "tagValue",
						Label:       "valor",
						Placeholder: "0.8   .3   10   20",
						Style:       1,
						Required:    false,
						MinLength:   0,
						MaxLength:   8, // float32 max int
					},
				),
			},
		},
	}
}

//	func ResponseTag() *discordgo.InteractionResponse {
//		return &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Components: []discordgo.MessageComponent{
//					newActionRow(
//						discordgo.Button{},
//					),
//				},
//			},
//		}
//	}
//

// TODO: do polymorphism
// 'see-more Tags' whether lenght about above a set number
func displayListModel(tags []*models.Tag) []discordgo.MessageComponent {
	var model []discordgo.MessageComponent
	for _, tag := range tags {
		// model = []discordgo.MessageComponent{
		// 	discordgo.Button{},
		// 	textDisplay(tag.TagName),
		// }
		model = append(model,
			textDisplay("**Nome** "+tag.TagName),
			textDisplay("**Descrição** "+tag.TagDescription),
			textDisplay(fmt.Sprint("**Valor** ", tag.TagValue)),
			separator(true, smallSpace),
		)
	}
	// row := newSection(model...)
	return model
}

func separator(b bool, space discordgo.SeparatorSpacingSize) discordgo.Separator {
	return discordgo.Separator{
		Divider: pointer(b),
		Spacing: pointer(space),
	}
}

// TODO: verify and split into pieces
// search for " " after or before
func textDisplay[T any](txt T) discordgo.TextDisplay {
	return discordgo.TextDisplay{
		Content: fmt.Sprint(txt),
	}
}

func newSection(ac discordgo.Button, txt ...discordgo.TextDisplay) discordgo.Section {
	if len(txt) > 3 {
		fmt.Println("WARN: muitos TextDisplay")
	}
	var comp []discordgo.MessageComponent
	for _, c := range txt {
		comp = append(comp, c)
	}
	return discordgo.Section{
		Components: comp,
		Accessory:  ac,
	}
}

func newActionRow(comp ...discordgo.MessageComponent) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: comp,
	}
}

func newContainer(AC *int, children ...discordgo.MessageComponent) discordgo.Container {
	return discordgo.Container{
		AccentColor: AC,
		Components:  children,
	}
}
