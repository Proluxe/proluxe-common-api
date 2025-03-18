package api

import (
	"fmt"
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_EVENTS(c *gin.Context, App *util.App) {
	// Fetch the Events
	client := App.SF.Client
	events := model.FetchEvents(client, "Name != null ORDER BY Start_Date_Time__c ASC")

	sortedEvents := model.SortedEvents(events)

	c.JSON(http.StatusOK, sortedEvents)
}

func GET_EVENT_DETAILS(c *gin.Context, App *util.App) {
	// Fetch event details by ID
	id := c.Param("id")
	client := App.SF.Client

	where := fmt.Sprintf("Id = '%s'", id)
	events := model.FetchEvents(client, where)
	event := events[0]

	c.JSON(http.StatusOK, event)
}

func CREATE_EVENT(c *gin.Context, App *util.App) {
	// Parse the event details from the request body
	var newEvent model.Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the new event in Salesforce
	client := App.SF.Client
	sobj := model.CreateEvent(client, newEvent)

	c.JSON(http.StatusCreated, sobj)
}

func UPDATE_EVENT(c *gin.Context, App *util.App) {
	eventId := c.Param("id")
	// Parse the event details from the request body
	var updatedEvent model.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = eventId

	// Update the event in Salesforce
	client := App.SF.Client
	sobj := model.UpdateEvent(client, updatedEvent)

	c.JSON(http.StatusOK, sobj)
}

func DELETE_EVENT(c *gin.Context, App *util.App) {
	// Fetch event details by ID
	id := c.Param("id")
	client := App.SF.Client
	// Delete the event
	client.SObject("Events__c").Set("Id", id).Delete()

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
