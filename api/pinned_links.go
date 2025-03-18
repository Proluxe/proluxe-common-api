package api

import (
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

// POST_CREATE_PINNED_LINK handles saving a pinned link
func POST_CREATE_PINNED_LINK(c *gin.Context, app *util.App) {
	var pinnedLink model.PinnedLink

	if err := c.ShouldBindJSON(&pinnedLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if pinnedLink.Email == "" || pinnedLink.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: Email or Path"})
		return
	}

	client := app.SF.Client
	link, err := model.CreatePinnedLink(client, pinnedLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pinned link", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, link)
}

// GET_PINNED_LINKS retrieves all pinned links for a user
func GET_PINNED_LINKS(c *gin.Context, app *util.App) {
	_, email, _ := GetCurrentUser(c)
	client := app.SF.Client

	links, err := model.FetchPinnedLinks(client, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pinned links", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": links})
}

func DELETE_PINNED_LINK(c *gin.Context, app *util.App) {
	id := c.Param("id")
	client := app.SF.Client

	err := model.DeletePinnedLink(client, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pinned link", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pinned link deleted successfully"})
}
