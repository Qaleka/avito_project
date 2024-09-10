package app

import (
	// "fmt"
	// "fmt"
	"fmt"
	"net/http"
	"time"

	// "time"

	"avito_project/internal/app/ds"
	"avito_project/internal/app/schemes"

	"github.com/gin-gonic/gin"
)

	func (app *Application) GetAllTenders(c *gin.Context) {
		var request schemes.GetAllTendersRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"reason": "Неверный формат запроса или его параметры: " + err.Error(),
			})
			return
		}
		tenders, err := app.repo.GetAllTenders(request.ServiceType, request.Limit, request.Offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"reason": "Неверный формат запроса или его параметры: " + err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, tenders)
	}

func (app *Application) AddTender(c *gin.Context) {
	var request schemes.AddTenderRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(request)
	tender := ds.Tender(request.Tender)
	tender.CreatedAt = time.Now()
	tender.UpdatedAt = time.Now()
	if err := app.repo.AddTender(&tender); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if err := app.repo.SaveTender(&tender); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tender)
}

func (app *Application) GetTender(c *gin.Context) {
	var request schemes.GetUserTendersRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат запроса или его параметры: " + err.Error(),
		})
	}
	tenders, err := app.repo.GetUserTenders(request.Limit, request.Offset, request.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

func (app *Application) GetTenderStatus(c *gin.Context) {
	var request schemes.GetTenderStatusRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(request)
	status, err := app.repo.GetTenderStatus(request.TenderId, request.Username)
	if status == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (app *Application) ChangeTenderStatus(c *gin.Context) {
	var request schemes.ChangeTenderStatusRequest
	if err := c.ShouldBind(&request.Parameters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
	}
	if *request.Parameters.Status != ds.CREATED && *request.Parameters.Status != ds.PUBLISHED && *request.Parameters.Status != ds.CLOSED {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Некоректный статус '%s'",*request.Parameters.Status),
		})
		return
	}
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(request)
	user, err := app.repo.GetUserByUsername(*request.Parameters.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при проверке пользователя",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не найден", *request.Parameters.Username),
		})
		return
	}
	tender, err := app.repo.GetTenderById(request.TenderId, *request.Parameters.Username)
	if tender == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if request.Parameters.Status != nil {
		tender.Status = *request.Parameters.Status
	}
	if err := app.repo.SaveTender(tender); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, tender)
}

func (app *Application) ChangeTender(c *gin.Context) {
	var request schemes.ChangeTenderRequest
	if err := c.ShouldBind(&request.Parameters); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBindQuery(&request.Query); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(request)
	if *request.Parameters.ServiceType != "Construction" && *request.Parameters.ServiceType != "Delivery" && *request.Parameters.ServiceType != "Manufacture" {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Некоректный тип сервиса '%s'",*request.Parameters.ServiceType),
		})
		return
	}
	user, err := app.repo.GetUserByUsername(*request.Query.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при проверке пользователя",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не найден", *request.Query.Username),
		})
		return
	}
	tender, err := app.repo.GetTenderById(request.TenderId, *request.Query.Username)
	if tender == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if request.Parameters.Name != nil {
		tender.Name = *request.Parameters.Name
	}

	if request.Parameters.Description != nil {
		tender.Description = *request.Parameters.Description
	}

	if request.Parameters.ServiceType != nil {
		tender.ServiceType = *request.Parameters.ServiceType
	}
	if err := app.repo.SaveTender(tender); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, tender)
}