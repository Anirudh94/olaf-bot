package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Anirudh94/olaf-bot/model/discord"
)

const discordApplicationPublicKey = "8cab6f3901ac6dddf7c7d5db2e9292998e784915106f7e96c70ccc061cbd3e3d"

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.Body)

	var discordInteraction discord.Interaction
	err := json.Unmarshal([]byte(request.Body), &discordInteraction)
	if err != nil {
		return handleError("Encountered error unmarshalling Interaction: " + err.Error())
	}

	// Verify the request is indeed from discord'
	if !discord.VerifyInteraction(&request, discordApplicationPublicKey) {
		return handleUnauthorized("Invalid request signature")
	}

	// See details of handling discord interactions: https://discord.com/developers/docs/interactions/slash-commands#receiving-an-interaction
	switch discordInteraction.Type {
	case discord.Ping:
		// ACK a PING message
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("{ type: %v }", discord.Pong), StatusCode: 200}, nil
	case discord.ApplicationCommand:
		// TODO: Just echoing back for now
		response, err := json.Marshal(discord.InteractionResponse{
			Type: discord.ChannelMessageWithSource,
			Data: discord.InteractionApplicationCommandCallbackData{
				TTS:     false,
				Content: "Congrats on sending a message!",
			},
		})
		if err != nil {
			return handleError("Encountered error marshalling InteractionResponse: " + err.Error())
		}
		return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
	default:
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Unknown Interaction Type '%v'", discordInteraction.Type), StatusCode: 200}, nil
	}
}

func handleUnauthorized(errMessage string) (events.APIGatewayProxyResponse, error) {
	fmt.Println(errMessage)
	return events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 401}, nil
}

func handleError(errMessage string) (events.APIGatewayProxyResponse, error) {
	// FIXME: Find a way to elegantly handle different typed errors such as ServiceFailureError
	fmt.Println(errMessage)
	return events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 500}, nil
}

func main() {
	lambda.Start(handleRequest)
}
