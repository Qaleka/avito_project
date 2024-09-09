package app

import (
	// "fmt"
	"fmt"
	"net/http"
	// "time"

	
	"avito_project/internal/app/ds"
	"avito_project/internal/app/schemes"

	"github.com/gin-gonic/gin"
)

func (app *Application) GetAllTenders (c *gin.Context) {
	var request schemes.GetAllTendersRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}
	fmt.Println(request)
	tenders, err := app.repo.GetAllTenders(request.ServiceType, request.Limit, request.Offset)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response := schemes.TenderOutput{Tenders:tenders}
	c.JSON(http.StatusOK, response)
}

func (app *Application) AddTender (c *gin.Context) {
	var request schemes.AddTenderRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		return
	}
	fmt.Println(request)
	tender := ds.Tender(request.Tender)
	
	if err := app.repo.AddTender(&tender); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := app.repo.SaveTender(&tender); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, request)
}