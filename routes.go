package main

import (
	"net/http"
	"time"

	"github.com/Proluxe/proluxe-common-api/api"
	middleware "github.com/Proluxe/proluxe-common-api/middleware"
	"github.com/Proluxe/proluxe-common-api/salesforce"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
	u "github.com/scottraio/go-utils"
)

func SetupRoutes(router *gin.Engine, app *util.App) {

	// Health Check
	router.GET("/", apiRoute(StatusOk, app))

	// Meta endpoint
	router.GET("/meta", apiRoute(api.GET_META, app))

	// Home endpoint
	router.GET("/home", apiRoute(api.GET_HOME, app))

	// Algolia Products (Basic Auth)
	basicUsername := u.GetDotEnvVariable("ALGOLIA_AUTH_USER")
	basicPassword := u.GetDotEnvVariable("ALGOLIA_AUTH_PASS")

	router.GET("/algolia/products",
		middleware.BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_PRODUCTS, app),
	)

	router.GET("/algolia/customers",
		middleware.BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_CUSTOMERS, app),
	)

	router.GET("/algolia/contacts",
		middleware.BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_CONTACTS, app),
	)

	router.GET("/algolia/parts",
		middleware.BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_PARTS, app),
	)

	router.GET("/algolia/vendors",
		middleware.BasicAuthMiddleware(basicUsername, basicPassword),
		apiRoute(api.GET_ALGOLIA_VENDORS, app),
	)

	// Messages
	router.GET("/messages/search/:email", apiRoute(api.GET_MESSAGES, app))
	router.GET("/messages/:email/:id", apiRoute(api.GET_MESSAGE_DETAILS, app))

	// Users
	router.POST("/users/pinned_links", apiRoute(api.POST_CREATE_PINNED_LINK, app))
	router.DELETE("/users/pinned_links/:id", apiRoute(api.DELETE_PINNED_LINK, app))
	router.GET("/users/pinned_links", apiRoute(api.GET_PINNED_LINKS, app))

	// Events
	router.POST("/events", apiRoute(api.CREATE_EVENT, app))
	router.POST("/events/:id", apiRoute(api.UPDATE_EVENT, app))
	router.GET("/events", apiRoute(api.GET_EVENTS, app))
	router.GET("/events/:id", apiRoute(api.GET_EVENT_DETAILS, app))
	router.DELETE("/events/:id", apiRoute(api.DELETE_EVENT, app))

	// Files
	router.GET("/files", apiRoute(api.GET_BUCKET_CONTENTS, app))
	router.POST("/files/make_public", apiRoute(api.POST_MAKE_PUBLIC, app))
	router.POST("/files/make_private", apiRoute(api.POST_MAKE_PRIVATE, app))
	router.POST("/files/upload", apiRoute(api.UPLOAD_FILE_TO_BUCKET, app))
	router.POST("/files/create_folder", apiRoute(api.CREATE_FOLDER_IN_BUCKET, app))
	router.GET("/files/download/*path", api.SERVE_FILE_FROM_BUCKET)
	router.DELETE("/files", apiRoute(api.DELETE_FILE_FROM_BUCKET, app))
	router.POST("/files/share", apiRoute(api.POST_SHARE_FILES, app))
	router.POST("/files/send", apiRoute(api.POST_SEND_FILES, app))

	// Comments
	router.GET("/comments/:recordID", apiRoute(api.GET_COMMENTS, app))
	router.POST("/comments/:recordID", apiRoute(api.POST_COMMENT, app))
	router.DELETE("/comments/:commentID", apiRoute(api.DELETE_COMMENT, app))

	// Issues
	router.GET("/issues", apiRoute(api.GET_ISSUES, app))
	router.GET("/issues/:id", apiRoute(api.GET_ISSUE, app))
	router.POST("/issues/:id", apiRoute(api.POST_UPDATE_ISSUE, app))
	router.POST("/issues", apiRoute(api.POST_CREATE_ISSUE, app))
	router.POST("/issues/:id/close", apiRoute(api.POST_CLOSE_ISSUE, app))

	// Users
	router.GET("/users/current", apiRoute(api.GET_DEFAULTS, app))
	router.POST("/users/:id/settings", apiRoute(api.POST_UPDATE_USER, app))

	router.POST("/users/:id/notes", apiRoute(api.POST_SAVE_USER_NOTES, app))
}

// Global funcs

func HandleRequest(handlerFunc func(*gin.Context, *util.App) interface{}, App *util.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		checkClient(App.SF, c) // Check client credentials or other preprocessing
		App.SF.QueryCount = 0  // Resetting the count for each request

		// Call the specific handler for this route
		result := handlerFunc(c, App)
		c.JSON(http.StatusOK, result)
	}
}

func apiRoute(handlerFunc func(*gin.Context, *util.App), App *util.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		checkClient(App.SF, c) // Check client credentials or other preprocessing
		handlerFunc(c, App)
	}
}

func clientErrorResponse(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "client not initialized"})
}

func checkClient(SF *salesforce.SF, c *gin.Context) {
	if SF.Client == nil {
		clientErrorResponse(c)
		return
	}
}

func parseDate(endDate string, c *gin.Context) {
	_, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		// If the date format is incorrect, return an error response
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Please use YYYY-MM-DD."})
		return
	}
}
