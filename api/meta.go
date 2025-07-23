package api

import (
	"net/http"

	"github.com/Proluxe/proluxe-common-api/model"
	"github.com/Proluxe/proluxe-common-api/util"
	"github.com/gin-gonic/gin"
)

// GET_META handles the /meta endpoint.
// It returns all metadata needed by the UI as a JSON payload.
func GET_META(c *gin.Context, App *util.App) {
	client := App.SF.Client

	// Fetch all meta data
	terms := model.FetchTerms(client)
	carriers := model.FetchCarriers(client)
	fobCodes := model.FetchFOBCodes(client, "Id != ''")
	shipMethods := model.FetchShipMethods(client)
	reps := model.FetchReps(client)

	freightTerms := model.FetchFreightTerms(client, "Id != ''")

	meta := model.Meta{
		Terms:        terms,
		Carriers:     carriers,
		FOBCodes:     fobCodes,
		ShipMethods:  shipMethods,
		FreightTerms: freightTerms,
		Reps:         reps,
	}

	c.JSON(http.StatusOK, meta)
}
