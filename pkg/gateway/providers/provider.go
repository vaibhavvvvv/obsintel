package providers

import "context"

type Provider interface{
	Chat(ctx context.Context, message string) (ChatResult, error)
}