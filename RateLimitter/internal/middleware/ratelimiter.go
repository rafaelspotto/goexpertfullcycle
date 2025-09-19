package middleware

import (
	"net"
	"strings"

	"rate-limiter/internal/limiter"

	"github.com/gin-gonic/gin"
)

// RateLimiterMiddleware creates a rate limiter middleware
func RateLimiterMiddleware(rateLimiter *limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract IP address
		ip := getClientIP(c)

		// Extract token from API_KEY header
		token := c.GetHeader("API_KEY")

		// Check rate limit
		result, err := rateLimiter.CheckRequest(c.Request.Context(), ip, token)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		// If request is not allowed, return 429
		if !result.Allowed {
			c.JSON(429, gin.H{
				"error":  "you have reached the maximum number of requests or actions allowed within a certain time frame",
				"reason": result.Reason,
			})
			c.Abort()
			return
		}

		// Add rate limit info to headers
		remaining, err := rateLimiter.GetRemainingRequests(c.Request.Context(), ip, token)
		if err == nil {
			c.Header("X-RateLimit-Remaining", string(rune(remaining)))
		}

		// Continue to next handler
		c.Next()
	}
}

// getClientIP extracts the real client IP from the request
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		if net.ParseIP(xri) != nil {
			return xri
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}
