package api

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) (string, string, string) {
	user, exists := c.Get("User")
	if !exists {
		log.Println("User not found in context")
		return "", "", ""
	}

	// Type assertion to convert user to map
	userMap, ok := user.(jwt.MapClaims)
	if !ok {
		log.Println("User type assertion failed")
		return "", "", ""
	}

	// Extract user details
	name, nameOk := userMap["name"].(string)
	email, emailOk := userMap["id"].(string)
	avatar, avatarOk := userMap["avatar"].(string)

	if !nameOk {
		log.Println("Failed to extract user name")
	}

	if !emailOk {
		log.Println("Failed to extract user email")
	}

	if !avatarOk {
		log.Println("Failed to extract user avatar")
	}

	return name, email, avatar
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
