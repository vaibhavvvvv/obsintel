package gateway

import "github.com/gin-gonic/gin"

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
	response, err := g.provider.Chat(g.ctx, req.Message)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"response": response})
}
