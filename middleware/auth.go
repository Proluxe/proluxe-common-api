// Package middleware provides authentication middleware for Gin.
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	u "github.com/scottraio/go-utils"
)

// AuthMiddleware validates requests using either a JWT token (preferred) or the legacy EXTERNAL_API_KEY (deprecated).
// Order of precedence:
//  1. JWT token in URL (?token=...)
//  2. JWT token in Authorization header (Bearer ...)
//  3. Legacy EXTERNAL_API_KEY in api_key query param or X-External-Api-Key header (deprecated, for backwards compatibility)
//
// If a valid JWT is found, user claims are set in context as "User".
// If neither is present or valid, the request is rejected with 401 Unauthorized.
//
// SECURITY NOTE: Passing JWTs in URLs is discouraged for new clients due to logging risks. Prefer Authorization header.

// Middleware to validate JWT except on the root path
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" ||
			c.Request.URL.Path == "/" ||
			c.Request.URL.Path == "/notify/build_order_changes" ||
			c.Request.URL.Path == "/mrp/timeline/prototype" ||
			c.Request.URL.Path == "/webhook/leads" ||
			c.Request.URL.Path == "/algolia/products" ||
			c.Request.URL.Path == "/algolia/contacts" ||
			c.Request.URL.Path == "/algolia/parts" ||
			c.Request.URL.Path == "/algolia/vendors" ||
			c.Request.URL.Path == "/algolia/customers" {
			c.Next()
			return
		}

		apiKey := c.Query("api_key")
		if apiKey == u.GetDotEnvVariable("EXTERNAL_API_KEY") {
			c.Next()
			return
		}

		// Authenticate with EXTERNAL_API_KEY from headers.
		externalApiKeyHeader := c.GetHeader("X-External-Api-Key")
		if externalApiKeyHeader != "" && externalApiKeyHeader == u.GetDotEnvVariable("EXTERNAL_API_KEY") {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			fmt.Println("Authorization token not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		fmt.Printf("Raw JWT token string: %s\n", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			fmt.Printf("JWT parse error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		if !token.Valid {
			fmt.Println("JWT token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("JWT claims are not of type MapClaims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims type"})
			return
		}
		fmt.Printf("User: %v\n", claims)
		c.Set("User", claims)

		c.Next()
	}
}
