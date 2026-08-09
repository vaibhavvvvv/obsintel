// main gateway struct and interface

package gateway

import (
	"context"
	"log"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vaibhavvvvv/obsintel/config"
	"github.com/vaibhavvvvv/obsintel/pkg/gateway/middleware"
	"github.com/vaibhavvvvv/obsintel/pkg/gateway/providers"
)

type Gateway struct {
	router   *gin.Engine
	provider providers.Provider // interface, not concrete type -> Chat(ctx,message)
	ctx      context.Context
}

func New(ctx context.Context) *Gateway {
	provider, err := providers.NewGeminiProvider(ctx, config.C.GEMINI_API_KEY)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	g := &Gateway{
		router:   gin.Default(),
		provider: provider,
		ctx:      ctx,
	}

	authApiKeys := middleware.NewAuthMiddleware(strings.Split(config.C.VALID_API_KEYS, ","))
	limiter := middleware.NewRateLimiter(3,time.Minute)

	g.router.Use(authApiKeys.Middleware())	
	g.router.Use(limiter.Middleware())
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