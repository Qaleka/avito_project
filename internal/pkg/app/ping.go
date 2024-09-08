package app

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (app *Application) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}