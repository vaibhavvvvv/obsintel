// token counting and cost accumulation
package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/vaibhavvvvv/obsintel/config"
)

type ModelPricing struct {
	InputPricePerToken  float64
	OutputPricePerToken float64
}

type CostTracker struct {
	pricing map[string]ModelPricing // keyed by model name
	mu      sync.RWMutex
}

type OpenRouterResponse struct {
	Data struct {
		Pricing struct {
			Input  string `json:"prompt"`
			Output string `json:"completion"`
		} `json:"pricing"`
	} `json:"data"`
}

func NewCostTracker() *CostTracker {
	return &CostTracker{pricing: make(map[string]ModelPricing)}
}

func (ct *CostTracker) FetchAndCachePricing(ctx context.Context) error {

	apiUrl := fmt.Sprintf("https://openrouter.ai/api/v1/model/%s/%s", "google", config.C.AI_MODEL)
	fmt.Println(apiUrl)

	response, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch pricing: %w", err)
	}
	defer response.Body.Close()

	var apiResponse OpenRouterResponse
	if err = json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return 	fmt.Errorf("JSON parsing error: %v\n", err)
	}

	inputPrice, err := strconv.ParseFloat(apiResponse.Data.Pricing.Input, 64)
	if err != nil {
		return fmt.Errorf("failed to parse input price float: %w", err)
	}
	outputPrice, err := strconv.ParseFloat(apiResponse.Data.Pricing.Output, 64)
	if err != nil {
		return fmt.Errorf("failed to parse output price float: %w", err)
	}
	// Perform a thread-safe write assignment using a Write Lock
	ct.mu.Lock()
	ct.pricing[config.C.AI_MODEL] = ModelPricing{
		InputPricePerToken:  inputPrice,
		OutputPricePerToken: outputPrice,
	}
	ct.mu.Unlock()

	return nil
}

func (ct *CostTracker) Calculate(model string, promptTokens, responseTokens int) float64 {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	modelPricing, exists := ct.pricing[model]
	if !exists {
		return 0.0
	}
	return (modelPricing.InputPricePerToken*float64(promptTokens) + modelPricing.OutputPricePerToken*float64(responseTokens))
}

// StartAutoRefresh boots a background routine that updates rates based on specifed time duration
func (ct *CostTracker) StartAutoRefresh(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return // Terminate the background routine if the app context closes
			case <-ticker.C:
				// Refresh the token cache
				if err := ct.FetchAndCachePricing(ctx); err != nil {
					log.Printf("Background cost refresh failed: %v", err)
				}
			}
		}
	}()
}
