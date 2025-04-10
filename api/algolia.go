package api

import (
	model "github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

func GET_ALGOLIA_PRODUCTS(c *gin.Context, App *util.App) {
	products := model.FetchAlgoliaProducts(App.SF.Client)

	c.JSON(200, products)
}
