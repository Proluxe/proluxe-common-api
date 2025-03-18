package api

import (
	"net/http"

	"github.com/Proluxe/proluxe-common-api/google"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_MESSAGES(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	email := c.Param("email")

	messages, _ := google.GetEmails(email)

	c.JSON(http.StatusOK, messages)
}

func GET_MESSAGE_DETAILS(c *gin.Context, App *util.App) {
	// Fetch the PO Lines
	email := c.Param("email")
	id := c.Param("id")

	message, _ := google.GetEmailByID(email, id)

	c.JSON(http.StatusOK, message)
}
