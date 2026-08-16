package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vaibhavvvvv/obsintel/config"
)

type SemanticCache struct {
	pool       *pgxpool.Pool
	threshhold float64 // similarity threshhold, 0.95
}

type OllamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func New(pool *pgxpool.Pool, threshold float64) *SemanticCache {
	return &SemanticCache{
		pool:       pool,
		threshhold: threshold,
	}
}

// Get returns cached response if similar query exists, empty string if miss
func (sc *SemanticCache) Get(ctx context.Context, apiKey, query string) (string, float64, bool, error) {

	embedding, err := sc.embed(ctx, query)
	if err != nil {
		return "", 0.0, false, fmt.Errorf("failed to generate embedding: %w", err)
	}

	var cachedResponse string
	var similarity float64
	var rowId string

	err = sc.pool.QueryRow(ctx, `SELECT response_text, 1 - (query_embedding <=> $1::float4[]::vector) as similarity, id
	FROM semantic_cache
	WHERE api_key = $2
	AND 1 - (query_embedding <=> $1::float4[]::vector) > $3
	ORDER BY query_embedding <=> $1::float4[]::vector
	LIMIT 1;`, embedding, apiKey, sc.threshhold).Scan(&cachedResponse, &similarity, &rowId)
	fmt.Println(query, ": ", similarity)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", 0.0, false, nil // cache miss
	}
	if err != nil {
		return "", 0.0, false, fmt.Errorf("failed to query database: %w", err)
	}

	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err = sc.pool.Exec(logCtx,
			"UPDATE semantic_cache SET hit_count = hit_count +1, last_accessed_at = NOW() where id=$1;",
			rowId); err != nil {
			log.Printf("Failed to log request: %v", err)
		}
	}()

	return cachedResponse, similarity, true, nil

}

func (sc *SemanticCache) Set(ctx context.Context, apiKey, query, response, model string) error {

	query_embedding, err := sc.embed(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to generate embedding for storage: %w", err)
	}

	_, err = sc.pool.Exec(ctx,
		"INSERT INTO semantic_cache (api_key, query_text, query_embedding, response_text, model, hit_count) VALUES($1,$2,$3::float4[]::vector,$4,$5,$6)",
		apiKey, query, query_embedding, response, model, 0)
	if err != nil {
		return fmt.Errorf("Failed to insert Cache Entry: %w", err)
	}
	return nil
}

// embed calls Ollama API to generate embedding for text
func (sc *SemanticCache) embed(ctx context.Context, text string) ([]float32, error) {
	reqBody := OllamaEmbedRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	ollamaEmbeddingUrl := config.C.OLLAMA_URL + "/api/embeddings"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaEmbeddingUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()

	var res ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if len(res.Embedding) == 0 {
		return nil, errors.New("received empty embedding vector from Ollama")
	}

	return res.Embedding, nil

}
