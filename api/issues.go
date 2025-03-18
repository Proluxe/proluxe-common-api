package api

import (
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_ISSUES(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	client := App.SF.Client
	issues := model.FetchIssues(client, "Closed__c = FALSE ORDER BY CreatedDate DESC")

	c.JSON(http.StatusOK, issues)
}

func GET_ISSUE(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	client := App.SF.Client
	id := c.Param("id")
	issue := model.FetchIssues(client, "Id = '"+id+"'")[0]

	issue.AttachRelatedObjects(client)

	c.JSON(http.StatusOK, issue)
}

func POST_CREATE_ISSUE(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	client := App.SF.Client

	var issue model.Issue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if issue.Create(client) == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create issue"})
		return
	}

	issue.SendNewIssueNotification(client)

	c.JSON(http.StatusOK, issue)
}

func POST_UPDATE_ISSUE(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	client := App.SF.Client

	var issue model.Issue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	issue.Id = c.Param("id")

	if issue.Update(client) == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update issue"})
		return
	}

	c.JSON(http.StatusOK, issue)
}

func POST_CLOSE_ISSUE(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	client := App.SF.Client

	var issue model.Issue
	issue.Id = c.Param("id")

	if issue.Close(client) == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close issue"})
		return
	}

	issue.SendClosedIssueNotification(client)

	c.JSON(http.StatusOK, issue)
}
