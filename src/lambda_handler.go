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
	fmt.Println("Processing request data for request: ", request.Body)

	var discordInteraction discord.Interaction
	err := json.Unmarshal([]byte(request.Body), &discordInteraction)
	if err != nil {
		return handleError("Encountered error unmarshalling Interaction: " + err.Error()), nil
	}

	// Verify the request is indeed from discord'
	if !discord.VerifyInteraction(request.Body, request.Headers, discordApplicationPublicKey) {
		return handleUnauthorized("invalid request signature"), nil
	}

	// See details of handling discord interactions: https://discord.com/developers/docs/interactions/slash-commands#receiving-an-interaction
	var response events.APIGatewayProxyResponse
	switch discordInteraction.Type {
	case discord.Ping:
		response = handleResponse(discord.InteractionResponse{
			Type: discord.Pong,
		})
	case discord.ApplicationCommand:
		// TODO: Just echoing back for now
		response = handleResponse(discord.InteractionResponse{
			Type: discord.ChannelMessageWithSource,
			Data: discord.InteractionApplicationCommandCallbackData{
				TTS:     false,
				Content: "Congrats on sending a message!",
			},
		})
	default:
		response = handleError(fmt.Sprintf("Unknown Interaction Type '%v'", discordInteraction.Type))
	}

	fmt.Println("Responding with: ", response)
	return response, nil
}

func handleResponse(body discord.InteractionResponse) (events.APIGatewayProxyResponse) {
	response, err := json.Marshal(body)
	if err != nil {
		return handleError("Encountered error marshalling InteractionResponse: " + err.Error())
	}
	fmt.Println("success response: ", response)
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}
}

func handleUnauthorized(errMessage string) (events.APIGatewayProxyResponse) {
	return events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 401}
}

func handleError(errMessage string) (events.APIGatewayProxyResponse) {
	// FIXME: Find a way to elegantly handle different typed errors as defined in src/error/*
	return events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 500}
}

func main() {
	lambda.Start(handleRequest)
}
