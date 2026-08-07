// token bucket implementation for rate limiting

package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	visitors sync.Map
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {

	return &RateLimiter{
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()

	value, exists := rl.visitors.Load(ip)
	if !exists {
		rl.visitors.Store(ip, &bucket{
			count:       1,
			windowStart: now,
		})
	}

	b, ok  := value.(*bucket)  //type assertion telling Go that this any value is *bucket type
	if !ok || b == nil {
		// value was nil or wrong type — treat as new IP
		rl.visitors.Store(ip, &bucket{
			count:       1,
			windowStart: now,
		})
		return true
	}
	

	if now.Sub(b.windowStart) >= rl.window {
		// window expired — reset the bucket
		b.count = 1
		b.windowStart = now
		rl.visitors.Store(ip, b)
		return true
	}

	//if window not expired then checking the count
	if b.count >= rl.limit {
		return false
	}

	//if within the window and rate limit then alow the req and increase the count
	b.count++
	rl.visitors.Store(ip, b)
	return true
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": "too many requests, please wait before trying again.",
			})
			c.Abort() //stops request here, else goes to next handler even if now allowed
			return
		}

		c.Next() //allowing request to continue to next handler
	}
}
