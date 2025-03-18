package api

import (
	"fmt"
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_DEFAULTS(c *gin.Context, App *util.App) {
	_, email, _ := GetCurrentUser(c)
	client := App.SF.Client

	links, err := model.FetchPinnedLinks(client, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pinned links", "details": err.Error()})
		return
	}

	issues := model.FetchIssues(client, "Closed__c = FALSE ORDER BY CreatedDate DESC")

	events := model.FetchEvents(client, "Name != null AND End_Date_Time__c > TODAY ORDER BY End_Date_Time__c ASC")

	users := model.FetchUsers(client, fmt.Sprintf("rstk__syusr_empl_email__c = '%s' LIMIT 1", email))

	c.JSON(http.StatusOK, gin.H{
		"PinnedLinks": links,
		"Issues":      issues,
		"Events":      events,
		"User":        users[0],
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
