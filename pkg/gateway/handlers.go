package gateway

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vaibhavvvvv/obsintel/config"
	"github.com/vaibhavvvvv/obsintel/store/queries"
)

func (g *Gateway) handleChat(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	if req.Message == "" {
		c.JSON(400, gin.H{"error": "message is required"})
		return
	}

	start := time.Now()
	response, err := g.provider.Chat(g.ctx, req.Message)
	latencyMs := int(time.Since(start).Milliseconds())
	apiKey, _ := c.Get("api_key")

	// db logging goes on seperetely w/o making whole api wait.
	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if logErr := queries.LogRequest(logCtx, queries.RequestLog{
			APIKey:    apiKey.(string), // type asserting as string cuz c.Get("api_key") returns an any
			Model:     config.C.AI_MODEL,
			LatencyMs: latencyMs,
			PromptTokens: response.PromptTokens,
			ResponseTokens: response.ResponseTokens,
			Success:   err == nil,
		}); logErr != nil {
			log.Printf("Failed to log request: %v", logErr)
		}
	}()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"response": response.Text})
}
