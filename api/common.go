package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

/*
GetCurrentUser extracts the user's email from the X-User-Email header.

Returns the email as a string, or an empty string if not present.

Usage:

	email := GetCurrentUser(c)
*/
func GetCurrentUser(c *gin.Context) string {
	email := c.GetHeader("X-User-Email")
	if email == "" {
		log.Println("X-User-Email header not found")
	}
	return email
}

func parseFloat(value interface{}) float64 {
	if value == nil {
		return 0
	}
	if floatValue, ok := value.(float64); ok {
		return floatValue
	}
	return 0
}
