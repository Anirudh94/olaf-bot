package image

import (
	"log"
	"strings"

	"github.com/Anirudh94/olaf-bot/util/discord"
	"github.com/Anirudh94/olaf-bot/util/unsplash"
	"github.com/bwmarrin/discordgo"
)

const (
	// CommandToken for image handler
	CommandToken  string = "image"
	keywordOption string = "keywords"
)

// HandleCommand for image handler
func HandleCommand(data *discord.ApplicationCommandInteractionData, unsplashClientID string) *discord.InteractionApplicationCommandCallbackData {
	var keywords []string
	for _, option := range data.Options {
		if option.Name == keywordOption {
			keywords = strings.Split(option.Value, " ")
		}
	}
	message := "Here's a thing"
	photoURL, err := unsplash.GetRandomPhoto(unsplashClientID, keywords, "")
	if err != nil {
		log.Println("Got error getting photo from unsplash: ", err)
		message = "Better luck next time..."
		photoURL = unsplash.ImageNotFoundURL
	}
	return &discord.InteractionApplicationCommandCallbackData{
		Content: message,
		Embeds: []discordgo.MessageEmbed{
			{
				Image: &discordgo.MessageEmbedImage{
					URL: photoURL,
				},
			},
		},
	}
}
