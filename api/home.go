package api

import (
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_HOME(c *gin.Context, App *util.App) {
	client := App.SF.Client

	_, email, _ := GetCurrentUser(c)

	// Fetch leads
	leads := model.FetchHomeLeads(client)

	// Fetch opportunities
	opportunities := model.FetchHomeOpportunities(client)

	// Fetch issues
	issues := model.FetchIssues(client, "Closed__c = FALSE ORDER BY CreatedDate DESC")

	// Fetch events
	events := model.FetchEvents(client, "Name != null AND End_Date_Time__c > TODAY ORDER BY End_Date_Time__c ASC")

	// Tasks
	tasks := model.FetchHomeTasks(client, email)

	// Fetch current user
	users := model.FetchUsers(client, "Email__c = '"+email+"' LIMIT 1")

	var currentUser interface{}
	if len(users) > 0 {
		currentUser = users[0]
	} else {
		currentUser = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"CurrentUser":   currentUser,
		"Leads":         leads,
		"Opportunities": opportunities,
		"Events":        events,
		"Issues":        issues,
		"Tasks":         tasks,
	})
}
