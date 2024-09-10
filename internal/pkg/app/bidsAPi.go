package app

import (
	"fmt"
	"net/http"

	"time"

	// "avito_project/internal/app/ds"
	// "avito_project/internal/app/schemes"

	"avito_project/internal/app/ds"
	"avito_project/internal/app/schemes"

	"github.com/gin-gonic/gin"
)

func (app *Application) GetAllBids (c *gin.Context) {

}

func (app *Application) AddBid(c *gin.Context) {
	var request schemes.AddBidRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(request)
	bid := ds.Bid(request.Bid)
	bid.CreatedAt = time.Now()
	bid.UpdatedAt = time.Now()
	if err := app.repo.AddBid(&bid); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if err := app.repo.SaveBid(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bid)
}