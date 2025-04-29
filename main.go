package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Proluxe/proluxe-common-api/api"
	"github.com/Proluxe/proluxe-common-api/salesforce"
	"github.com/Proluxe/proluxe-common-api/services"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	u "github.com/scottraio/go-utils"
)

// JWT middleware
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" ||
			c.Request.URL.Path == "/" ||
			c.Request.URL.Path == "/notify/build_order_changes" ||
			c.Request.URL.Path == "/mrp/timeline/prototype" ||
			c.Request.URL.Path == "/webhook/leads" ||
			c.Request.URL.Path == "/algolia/products" ||
			c.Request.URL.Path == "/algolia/contacts" ||
			c.Request.URL.Path == "/algolia/parts" ||
			c.Request.URL.Path == "/algolia/customers" {
			c.Next()
			return
		}

		apiKey := c.Query("api_key")
		if apiKey == u.GetDotEnvVariable("EXTERNAL_API_KEY") {
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		fmt.Printf("User: %v\n", claims)
		c.Set("User", claims)

		c.Next()
	}
}

// Basic Auth middleware for specific routes
func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, hasAuth := c.Request.BasicAuth()
		if !hasAuth || user != username || pass != password {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()

	// Setup JWT middleware
	secret := u.GetDotEnvVariable("JWT_SECRET")
	router.Use(JWTAuthMiddleware(secret))

	// Setup error handling middleware
	router.Use(services.ErrorHandling())

	// Setup CORS
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	SF := salesforce.NewSF()
	app := util.App{SF: SF}

	// Health Check
	router.GET("/", apiRoute(StatusOk, &app))

	// Algolia Products (Basic Auth)
	basicUsername := u.GetDotEnvVariable("ALGOLIA_AUTH_USER")
	basicPassword := u.GetDotEnvVariable("ALGOLIA_AUTH_PASS")

	router.GET("/algolia/products",
		BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_PRODUCTS, &app),
	)

	router.GET("/algolia/customers",
		BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_CUSTOMERS, &app),
	)

	router.GET("/algolia/contacts",
		BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_CONTACTS, &app),
	)

	router.GET("/algolia/parts",
		BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_PARTS, &app),
	)

	// Messages
	router.GET("/messages/search/:email", apiRoute(api.GET_MESSAGES, &app))
	router.GET("/messages/:email/:id", apiRoute(api.GET_MESSAGE_DETAILS, &app))

	// Users
	router.POST("/users/pinned_links", apiRoute(api.POST_CREATE_PINNED_LINK, &app))
	router.DELETE("/users/pinned_links/:id", apiRoute(api.DELETE_PINNED_LINK, &app))
	router.GET("/users/pinned_links", apiRoute(api.GET_PINNED_LINKS, &app))

	// Events
	router.POST("/events", apiRoute(api.CREATE_EVENT, &app))
	router.POST("/events/:id", apiRoute(api.UPDATE_EVENT, &app))
	router.GET("/events", apiRoute(api.GET_EVENTS, &app))
	router.GET("/events/:id", apiRoute(api.GET_EVENT_DETAILS, &app))
	router.DELETE("/events/:id", apiRoute(api.DELETE_EVENT, &app))

	// Files
	router.GET("/files", apiRoute(api.GET_BUCKET_CONTENTS, &app))
	router.POST("/files/make_public", apiRoute(api.POST_MAKE_PUBLIC, &app))
	router.POST("/files/make_private", apiRoute(api.POST_MAKE_PRIVATE, &app))
	router.POST("/files/upload", apiRoute(api.UPLOAD_FILE_TO_BUCKET, &app))
	router.POST("/files/create_folder", apiRoute(api.CREATE_FOLDER_IN_BUCKET, &app))
	router.GET("/files/download/*path", api.SERVE_FILE_FROM_BUCKET)
	router.DELETE("/files", apiRoute(api.DELETE_FILE_FROM_BUCKET, &app))
	router.POST("/files/share", apiRoute(api.POST_SHARE_FILES, &app))
	router.POST("/files/send", apiRoute(api.POST_SEND_FILES, &app))

	// Comments
	router.GET("/comments/:recordID", apiRoute(api.GET_COMMENTS, &app))
	router.POST("/comments/:recordID", apiRoute(api.POST_COMMENT, &app))
	router.DELETE("/comments/:commentID", apiRoute(api.DELETE_COMMENT, &app))

	// Issues
	router.GET("/issues", apiRoute(api.GET_ISSUES, &app))
	router.GET("/issues/:id", apiRoute(api.GET_ISSUE, &app))
	router.POST("/issues/:id", apiRoute(api.POST_UPDATE_ISSUE, &app))
	router.POST("/issues", apiRoute(api.POST_CREATE_ISSUE, &app))
	router.POST("/issues/:id/close", apiRoute(api.POST_CLOSE_ISSUE, &app))

	// Users
	router.GET("/users/current", apiRoute(api.GET_DEFAULTS, &app))
	router.POST("/users/:id/settings", apiRoute(api.POST_UPDATE_USER, &app))

	// Start server
	router.Run(":" + u.GetDotEnvVariable("PORT"))
}

func StatusOk(c *gin.Context, App *util.App) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
