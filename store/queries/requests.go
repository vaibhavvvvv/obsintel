package queries

import (
	"context"
	"fmt"
	"strings"

	"github.com/vaibhavvvvv/obsintel/store"
)

type RequestLog struct {
    APIKey         string
    Model          string
    PromptTokens   int
    ResponseTokens int
    LatencyMs      int
    Success        bool
}

func LogRequest(ctx context.Context, r RequestLog) error {
	if strings.TrimSpace(r.APIKey) == "" || strings.TrimSpace(r.Model) == "" {
		return fmt.Errorf("api_key and model are required.")
	}
	_, err := store.Pool.Exec(ctx, 
	"INSERT INTO requests (api_key, model, prompt_tokens, response_tokens, latency_ms, success) VALUES($1,$2,$3,$4,$5,$6)",
	r.APIKey,r.Model,r.PromptTokens,r.ResponseTokens,r.LatencyMs,r.Success)
	return err
}