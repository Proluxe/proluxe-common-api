package api

import (
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_COMMENTS(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	recordID := c.Param("recordID")
	whereCondition := "Record_ID__c = '" + recordID + "' ORDER BY CreatedDate ASC"

	comments := model.FetchComments(App.SF.Client, whereCondition)
	commentMentions := model.FetchMentions(App.SF.Client, whereCondition)

	response := gin.H{
		"Comments":        comments,
		"CommentMentions": commentMentions,
	}

	c.JSON(http.StatusOK, response)
}

func POST_COMMENT(c *gin.Context, App *util.App) {
	client := App.SF.Client
	recordID := c.Param("recordID")

	_, email, _ := GetCurrentUser(c)
	comment := model.Comment{
		RecordID:  recordID,
		CreatedBy: email,
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	comment.Create(client)
	comment.SendNotificationEmail(client)

	c.JSON(http.StatusCreated, gin.H{})
}

func DELETE_COMMENT(c *gin.Context, App *util.App) {
	commentID := c.Param("commentID")

	comment := model.Comment{
		Id: commentID,
	}

	if err := comment.Delete(App.SF.Client); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
