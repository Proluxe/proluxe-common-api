package api

import (
	"fmt"
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_USERS(c *gin.Context, App *util.App) {
	client := App.SF.Client

	// Fetch all users with a specific condition
	whereCondition := "Id != NULL"
	users := model.FetchAppUsers(client, whereCondition)

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active users found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func POST_SAVE_USER_NOTES(c *gin.Context, App *util.App) {
	client := App.SF.Client

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Id = c.Param("id")

	if user.Notes == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notes cannot be empty"})
		return
	}

	user.SaveNotes(client)

	c.JSON(http.StatusOK, gin.H{"message": "User notes saved successfully"})
}

func GET_DEFAULTS(c *gin.Context, App *util.App) {
	email := GetCurrentUser(c)
	client := App.SF.Client

	origin := c.Query("app")

	links, err := model.FetchPinnedLinks(client, email, origin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pinned links", "details": err.Error()})
		return
	}

	issues := model.FetchIssues(client, "Closed__c = FALSE ORDER BY CreatedDate DESC")

	events := model.FetchEvents(client, "Name != null AND End_Date_Time__c > TODAY ORDER BY End_Date_Time__c ASC")

	users := model.FetchUsers(client, fmt.Sprintf("Email__c = '%s' LIMIT 1", email))

	mentionableUsers := model.MentionableUsers(client)

	var userResp interface{}
	if len(users) > 0 {
		userResp = users[0]
	} else {
		userResp = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"PinnedLinks":      links,
		"Issues":           issues,
		"Events":           events,
		"User":             userResp,
		"MentionableUsers": mentionableUsers,
	})
}

func POST_UPDATE_USER(c *gin.Context, App *util.App) {
	client := App.SF.Client

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user.Id = c.Param("id")

	if user.Update(client) == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
