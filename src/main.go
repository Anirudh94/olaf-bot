package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Anirudh94/olaf-bot/handlers/image"
	"github.com/Anirudh94/olaf-bot/util/discord"
)

// Environment variables
var (
	discordApplicationPublicKey string = os.Getenv("DISCORD_APPLICATION_PUBLIC_KEY")
	unsplashClientID            string = os.Getenv("UNSPLASH_CLIENT_ID")
)

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Processing request data for request: ", request.Body)

	var interaction discord.Interaction
	err := json.Unmarshal([]byte(request.Body), &interaction)
	if err != nil {
		return *errorResp("Encountered error unmarshalling Interaction: " + err.Error()), nil
	}

	// Verify the request is indeed from discord
	if !discord.VerifyInteraction(request.Body, request.Headers, discordApplicationPublicKey) {
		return *unauthorizedResp("invalid request signature"), nil
	}

	var response *events.APIGatewayProxyResponse

	// See details of handling discord interactions: https://discord.com/developers/docs/interactions/slash-commands#receiving-an-interaction
	switch interaction.Type {
	case discord.Ping:
		response = successResp(&discord.InteractionResponse{
			Type: discord.Pong,
		})
	case discord.ApplicationCommand:
		callbackData := routeToCorrectHandler(&interaction.Data)
		response = successResp(&discord.InteractionResponse{
			Type: discord.ChannelMessageWithSource,
			Data: *callbackData,
		})
	default:
		response = errorResp(fmt.Sprintf("Unknown Interaction Type '%v'", interaction.Type))
	}

	fmt.Println("Responding with: ", response)
	return *response, nil
}

func routeToCorrectHandler(interactionData *discord.ApplicationCommandInteractionData) *discord.InteractionApplicationCommandCallbackData {
	var callbackData discord.InteractionApplicationCommandCallbackData
	switch interactionData.Name {
	case image.CommandToken:
		callbackData = *image.HandleCommand(interactionData, unsplashClientID)
	default:
		callbackData = discord.InteractionApplicationCommandCallbackData{
			TTS:     false,
			Content: "Congrats on sending a message!",
		}
	}
	return &callbackData
}

func successResp(body *discord.InteractionResponse) *events.APIGatewayProxyResponse {
	response, err := json.Marshal(*body)
	if err != nil {
		return errorResp("Encountered error marshalling InteractionResponse: " + err.Error())
	}
	return &events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}
}

func unauthorizedResp(errMessage string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 401}
}

func errorResp(errMessage string) *events.APIGatewayProxyResponse {
	// FIXME: Find a way to elegantly handle different typed errors as defined in src/error/*
	return &events.APIGatewayProxyResponse{Body: errMessage, StatusCode: 500}
}
