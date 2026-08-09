package providers

import (
	"context"

	"github.com/vaibhavvvvv/obsintel/config"
	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
	chat   *genai.Chat
}

func NewGeminiProvider(ctx context.Context, apiKey string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}
    chat, err := client.Chats.Create(ctx, config.C.AI_MODEL, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeminiProvider{client: client, chat: chat}, nil
}

func (g *GeminiProvider) Chat(ctx context.Context, message string) (string, error) {
	result, err := g.chat.SendMessage(ctx, genai.Part{Text: message})
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}