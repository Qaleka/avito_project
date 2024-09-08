package app

import (
	// "fmt"
	"net/http"
	// "time"

	// "avito_project/internal/app/ds"
	"avito_project/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

func (app *Application) GetAllTenders (c *gin.Context) {
	var request schemes.GetAllTendersRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}

	tenders, err := app.repo.GetAllTenders(request.ServiceType)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response := schemes.TenderOutput{Tenders:tenders}
	c.JSON(http.StatusOK, response)
}