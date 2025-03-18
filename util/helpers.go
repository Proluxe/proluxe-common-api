package util

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/Proluxe/proluxe-common-api/salesforce"
	"github.com/gin-gonic/gin"
)

type App struct {
	SF *salesforce.SF
}

func HandleRequest(handlerFunc func(*gin.Context, *App) interface{}, App *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		checkClient(App.SF, c) // Check client credentials or other preprocessing
		App.SF.QueryCount = 0  // Resetting the count for each request

		// Call the specific handler for this route
		result := handlerFunc(c, App)

		log.Println("Query count: ", App.SF.QueryCount)
		c.JSON(http.StatusOK, result)
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

func ParseDate(endDate string, c *gin.Context) {
	_, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		// If the date format is incorrect, return an error response
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Please use YYYY-MM-DD."})
		return
	}
}

// RemoveDuplicateIDs filters out duplicate item IDs from the list.
func RemoveDuplicateIDs(itemIDs []string) []string {
	seen := make(map[string]bool)
	var uniqueItemIDs []string
	for _, id := range itemIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = true
			uniqueItemIDs = append(uniqueItemIDs, id)
		}
	}
	return uniqueItemIDs
}

func PSTDateToCivil(now time.Time) civil.Date {
	// Define the PST location, UTC-8
	pstLocation, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return civil.Date{}
	}

	// Convert current time to PST
	nowPST := now.In(pstLocation) // T
	return civil.DateOf(nowPST)
}

func PSTDateToTime(now time.Time) time.Time {
	// Define the PST location, UTC-8
	pstLocation, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Time{}
	}

	// Convert current time to PST
	nowPST := now.In(pstLocation) // T
	return nowPST
}
