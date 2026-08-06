// main gateway struct and interface

package gateway

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vaibhavvvvv/obsintel/config"
	"github.com/vaibhavvvvv/obsintel/pkg/gateway/providers"
)

type Gateway struct {
	router   *gin.Engine
	provider providers.Provider // interface, not concrete type -> Chat(ctx,message)
	ctx      context.Context
}

func New(ctx context.Context) *Gateway {
	provider, err := providers.NewGeminiProvider(ctx, config.C.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	g := &Gateway{
		router:   gin.Default(),
		provider: provider,
		ctx:      ctx,
	}

	g.registerRoutes()
	return g

}

func (g *Gateway) Run(port string) {
	log.Printf("Server starting on port %s", port)
	g.router.Run(":" + port)
}

func (g *Gateway) registerRoutes() {
	g.router.POST("/chat", g.handleChat)
}

// func Init(){
// 	ctx := context.Background()

// 	// Create a new Chat.
// 	chat, err := client.Chats.Create(ctx, "gemini-3.5-flash-lite", nil, nil)
// 	if err != nil {
// 		fmt.Printf("Failed to create chat: %v", err)
// 	}

// 	router.POST("/chat", func(c *gin.Context) {
// 		var req struct {
// 			Message string `json:"message"`
// 		}
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(400, gin.H{"error": "invalid request"})
// 			return
// 		}
// 		if req.Message == "" {
// 			c.JSON(400, gin.H{"error": "message is required"})
// 			return
// 		}
// 		response, err := providers.ChatWithGemini(ctx, chat, req.Message)
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"response": response,
// 		})
// 	})

// }
