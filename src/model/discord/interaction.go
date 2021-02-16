// Package discord models the types defined here: https://discord.com/developers/docs/interactions/slash-commands#interaction
package discord

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"

	"github.com/aws/aws-lambda-go/events"
)

type InteractionType int64

const (
	Ping               InteractionType = 1
	ApplicationCommand InteractionType = 2
)

type Interaction struct {
	Type      InteractionType                   `json:"type"`
	Token     string                            `json:"token"`
	Member    GuildMember                       `json:"member"`
	ID        string                            `json:"id"`
	GuildID   string                            `json:"guild_id"`
	Data      ApplicationCommandInteractionData `json:"data"`
	ChannelID string                            `json:"channel_id"`
	Version   int64                             `json:"version"`
}

type ApplicationCommandInteractionData struct {
	Options []ApplicationCommandInteractionDataOption `json:"options"`
	Name    string                                    `json:"name"`
	ID      string                                    `json:"id"`
}

type ApplicationCommandInteractionDataOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GuildMember struct {
	User         User        `json:"user"`
	Roles        []string    `json:"roles"`
	PremiumSince interface{} `json:"premium_since"`
	Permissions  string      `json:"permissions"`
	Pending      bool        `json:"pending"`
	Nick         interface{} `json:"nick"`
	Mute         bool        `json:"mute"`
	JoinedAt     string      `json:"joined_at"`
	IsPending    bool        `json:"is_pending"`
	Deaf         bool        `json:"deaf"`
}

type User struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	PublicFlags   int64  `json:"public_flags"`
}

type InteractionResponseType int64

const (
	// Pong ACK a Ping
	Pong InteractionResponseType = 1
	// Acknowledge ACK a command without sending a message, eating the user's input
	Acknowledge InteractionResponseType = 2
	// ChannelMessage respond with a message, eating the user's input
	ChannelMessage InteractionResponseType = 3
	// ChannelMessageWithSource respond with a message, showing the user's input
	ChannelMessageWithSource InteractionResponseType = 4
	// AcknowledgeWithSource ACK a command without sending a message, showing the user's input
	AcknowledgeWithSource InteractionResponseType = 5
)

type InteractionResponse struct {
	Type InteractionResponseType                   `json:"type"`
	Data InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}

type InteractionApplicationCommandCallbackData struct {
	TTS             bool          `json:"tts,omitempty"`
	Content         string        `json:"content"`
	Embeds          []interface{} `json:"embeds,omitempty"`
	AllowedMentions []interface{} `json:"allowed_mentions,omitempty"`
}

// VerifyInteraction does AuthN/Z on the request: https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
// This is mostly copy paste from discordgo: https://github.com/bwmarrin/discordgo/blob/ad76e324502b76c7507178ed07b242841c0724a4/interactions.go
func VerifyInteraction(r *events.APIGatewayProxyRequest, hexPublicKey string) bool {
	var msg bytes.Buffer

	signature := r.Headers["X-Signature-Ed25519"]
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	timestamp := r.Headers["X-Signature-Timestamp"]
	if timestamp == "" {
		return false
	}

	key, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return false
	}

	msg.WriteString(timestamp + r.Body)

	return ed25519.Verify(key, msg.Bytes(), sig)
}
