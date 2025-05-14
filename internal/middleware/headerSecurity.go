package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware sets various security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Prevent Clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// 2. Prevent MIME-type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// 3. Enable XSS Protection in older browsers
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// 4. Strict Transport Security (HSTS)
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 5. Content Security Policy (CSP)
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self';")

		// 6. Referrer Policy
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")

		// 7. Permissions Policy
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")

		// 8. Cross-Origin Resource Sharing (CORS) Headers (optional)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Proceed to the next middleware or handler
		c.Next()
	}
}
