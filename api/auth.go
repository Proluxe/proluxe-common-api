package api

import (
	"net/http"

	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_VERIFY_EMAIL(c *gin.Context, App *util.App) {
	// client := App.SF.Client
	// email := c.Param("email")

	// resp := model.VerifyEmail(client, email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verification is not implemented yet.", // "data":    resp,
	})
}
