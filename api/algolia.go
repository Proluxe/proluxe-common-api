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

func GET_ALGOLIA_CUSTOMERS(c *gin.Context, App *util.App) {
	products := model.FetchAlgoliaCustomers(App.SF.Client)

	c.JSON(200, products)
}

func GET_ALGOLIA_CONTACTS(c *gin.Context, App *util.App) {
	products := model.FetchAlgoliaContacts(App.SF.Client)

	c.JSON(200, products)
}

func GET_ALGOLIA_PARTS(c *gin.Context, App *util.App) {
	parts := model.FetchAlgoliaParts(App.SF.Client)

	c.JSON(200, parts)
}

func GET_ALGOLIA_VENDORS(c *gin.Context, App *util.App) {
	vendors := model.FetchAlgoliaVendors(App.SF.Client)

	c.JSON(200, vendors)
}

func GET_ALGOLIA_CATALOG(c *gin.Context, App *util.App) {
	catalog := model.FetchAlgoliaCatalog(App.SF.Client)

	c.JSON(200, catalog)
}
