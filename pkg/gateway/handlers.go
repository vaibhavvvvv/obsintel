package gateway

import (
	"context"
	"fmt"
	"log"
	"strings"
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

	apiKey, _ := c.Get("api_key")

	cachedResponse, similarity, hit, err := g.cache.Get(g.ctx, apiKey.(string), req.Message)
	if err != nil {
		log.Printf("Cache lookup failed: %v", err)
	}
	if hit {
		c.JSON(200, gin.H{
			"response":   cachedResponse,
			"cached":     true,
			"similarity": similarity,
		})
		return
	}

	start := time.Now()
	response, err := g.provider.Chat(g.ctx, req.Message)
	latencyMs := int(time.Since(start).Milliseconds())

	// store in cache after successful LLM call
	if err == nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if cacheErr := g.cache.Set(cacheCtx, apiKey.(string), req.Message, response.Text, config.C.AI_MODEL); cacheErr != nil {
				log.Printf("Cache set failed: %v", cacheErr)
			}
		}()
	}

	// db logging goes on seperetely w/o making whole api wait.
	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if logErr := queries.LogRequest(logCtx, queries.RequestLog{
			APIKey:         apiKey.(string), // type asserting as string cuz c.Get("api_key") returns an any
			Model:          config.C.AI_MODEL,
			LatencyMs:      latencyMs,
			PromptTokens:   response.PromptTokens,
			ResponseTokens: response.ResponseTokens,
			Cost:           g.ct.Calculate(config.C.AI_MODEL, response.PromptTokens, response.ResponseTokens),
			Error: func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
			Success: err == nil,
		}); logErr != nil {
			log.Printf("Failed to log request: %v", logErr)
		}
	}()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"response":   response.Text,
		"cached":     false,
		"similarity": nil,
	})
}

func (g *Gateway) handleChatStream(c *gin.Context) {
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

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	apiKey, _ := c.Get("api_key")

	cachedResponse, similarity, hit, err := g.cache.Get(g.ctx, apiKey.(string), req.Message)
	if err != nil {
		log.Printf("Cache lookup failed: %v", err)
	}
	if hit {
		// send it as a single SSE event
		fmt.Fprintf(c.Writer, "data: %s\n\n", cachedResponse)
		c.Writer.Flush()
		c.SSEvent("done", gin.H{"status": "complete", "cached": true, "similarity": similarity})
		c.Writer.Flush()
		return
	}

	start := time.Now()
	response, err := g.provider.ChatStream(c.Request.Context(), req.Message, func(token string) {
		lines := strings.Split(token, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", line) //\n\n is mandatory SSE format — two newlines after each data line. Without it clients won't parse the events correctly.
		}
		c.Writer.Flush() //Without this, chunks may sit in a buffer and arrive in batches, defeating the purpose of streaming.it pushes data to client directly.
	})

	if err != nil {
		// send error through the SSE stream itself, not as a new HTTP response
		c.SSEvent("error", gin.H{"error": "stream failed"})
		c.Writer.Flush()
	} else {
		c.SSEvent("done", gin.H{"status": "complete", "cached": false, "similarity": nil})
		c.Writer.Flush()
	}

	latencyMs := int(time.Since(start).Milliseconds())

	// store in cache after successful LLM call
	if err == nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if cacheErr := g.cache.Set(cacheCtx, apiKey.(string), req.Message, response.Text, config.C.AI_MODEL); cacheErr != nil {
				log.Printf("Cache set failed: %v", cacheErr)
			}
		}()
	}
	// db logging goes on seperetely w/o making whole api wait.
	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if logErr := queries.LogRequest(logCtx, queries.RequestLog{
			APIKey:         apiKey.(string), // type asserting as string cuz c.Get("api_key") returns an any
			Model:          config.C.AI_MODEL,
			LatencyMs:      latencyMs,
			PromptTokens:   response.PromptTokens,
			ResponseTokens: response.ResponseTokens,
			Cost:           g.ct.Calculate(config.C.AI_MODEL, response.PromptTokens, response.ResponseTokens),
			Error: func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
			Success: err == nil,
		}); logErr != nil {
			log.Printf("Failed to log request: %v", logErr)
		}
	}()

}
