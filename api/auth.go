package api

import (
	"net/http"

	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_VERIFY_EMAIL(c *gin.Context, App *util.App) {
	client := App.SF.Client
	email := c.Param("email")

	resp := model.VerifyEmail(client, email)

	c.JSON(http.StatusOK, resp)
}
