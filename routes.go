package main

import (
	"net/http"
	"time"

	"github.com/Proluxe/proluxe-common-api/salesforce"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

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
