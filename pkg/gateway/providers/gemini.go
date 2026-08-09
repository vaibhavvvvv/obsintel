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

type ChatResult struct {
	Text           string
	PromptTokens   int
	ResponseTokens int
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

func (g *GeminiProvider) Chat(ctx context.Context, message string) (ChatResult, error) {
    var chatresponse ChatResult
	result, err := g.chat.SendMessage(ctx, genai.Part{Text: message})
	if err != nil {
		return chatresponse, err
	}

    chatresponse.PromptTokens = int(result.UsageMetadata.PromptTokenCount)
    chatresponse.ResponseTokens = int(result.UsageMetadata.CandidatesTokenCount)
    chatresponse.Text = result.Text()
	return chatresponse, nil
}