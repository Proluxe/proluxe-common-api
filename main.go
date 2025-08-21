package main

import (
	"net/http"
	"strconv"

	middleware "github.com/Proluxe/proluxe-common-api/middleware"
	"github.com/Proluxe/proluxe-common-api/salesforce"
	"github.com/Proluxe/proluxe-common-api/services"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
	u "github.com/scottraio/go-utils"
)

func main() {
	router := gin.Default()
	env := u.GetDotEnvVariable("ENV")

	// Setup error handling middleware
	router.Use(services.ErrorHandling())

	SF := salesforce.NewSF()
	app := util.App{SF: SF}
	SetupRoutes(router, &app)

	port := util.Port()

	if env == "development" {
		middleware.StartDevMode(port, "COMMON")
	}

	router.Run(":" + strconv.Itoa(port))
}

func StatusOk(c *gin.Context, App *util.App) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
