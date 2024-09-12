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
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат данных: " + err.Error(),
		})
		return
	}
	fmt.Println(request)
	var tender ds.Tender
	tender.CreatedAt = time.Now()
	if err := app.repo.AddTender(&tender, request); err != nil {
		switch err.Error() {
		case "пользователь не найден":
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": "Пользователь с именем '" + request.CreatorUsername + "' не найден",
			})
		case "организация не найдена":
			c.JSON(http.StatusNotFound, gin.H{
				"reason": "Организация с ID '" + request.OrganizationID + "' не найдена",
			})
		case "организация не совпадает":
			c.JSON(http.StatusForbidden, gin.H{
				"reason": "Пользователь '" + request.CreatorUsername + "' не принадлежит к организации с ID '" + request.OrganizationID + "'",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"reason": err.Error(),
			})
		}
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
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат запроса или его параметры: " + err.Error(),
		})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат URI параметров: " + err.Error(),
		})
		return
	}

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат query параметров: " + err.Error(),
		})
		return
	}

	// Получение статуса тендера
	status, err := app.repo.GetTenderStatus(request.TenderId, request.Username)
	if err != nil {
		if err.Error() == "пользователь не найден" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": fmt.Sprintf("Пользователь с именем '%s' не найден", request.Username),
			})
			return
		} else if err.Error() == "нет права доступа" {
			c.JSON(http.StatusForbidden, gin.H{
				"reason": "Пользователь не имеет прав на доступ к тендеру",
			})
			return
		} else if err.Error() == "тендер не найден" {
			c.JSON(http.StatusNotFound, gin.H{
				"reason": fmt.Sprintf("Тендер с id '%s' не найден", request.TenderId),
			})
			return
		}
		// Общая ошибка
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": err.Error(),
		})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, status)
}

func (app *Application) ChangeTenderStatus(c *gin.Context) {
	var request schemes.ChangeTenderStatusRequest

	// Обработка URI параметров
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Неверный формат URI параметров: '%s'", err.Error()),
		})
		return
	}

	// Обработка query параметров
	if err := c.ShouldBindQuery(&request.Parameters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Неверный формат параметров запроса: '%s'", err.Error()),
		})
		return
	}

	// Проверка статуса
	if request.Parameters.Status != ds.CREATED && request.Parameters.Status != ds.PUBLISHED && request.Parameters.Status != ds.CLOSED {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Некорректный статус '%s'", request.Parameters.Status),
		})
		return
	}

	// Проверка пользователя
	user, err := app.repo.GetUserByUsername(request.Parameters.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при проверке пользователя",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не найден", request.Parameters.Username),
		})
		return
	}

	tender, err := app.repo.GetTenderById(request.URI.TenderId, user.ID)
	if err != nil {
		if err.Error() == "организация не совпадает" {
			c.JSON(http.StatusForbidden, gin.H{
				"reason": fmt.Sprintf("Пользователь '%s' не принадлежит к организации, связанной с тендером", request.Parameters.Username),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if tender.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не имеет прав на редактирование тендера", request.Parameters.Username),
		})
		return
	}

	tender.Status = request.Parameters.Status
	if err := app.repo.SaveTender(tender); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении тендера",
		})
		return
	}

	c.JSON(http.StatusOK, tender)
}


func (app *Application) ChangeTender(c *gin.Context) {
	var request schemes.ChangeTenderRequest

	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат URI параметров: " + err.Error(),
		})
		return
	}

	if err := c.ShouldBindQuery(&request.Query); err != nil || request.Query.Username == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат параметров запроса или отсутствует username",
		})
		return
	}

	if err := c.ShouldBind(&request.Parameters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат параметров тела запроса: " + err.Error(),
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

	tender, err := app.repo.GetTenderById(request.URI.TenderId, user.ID)
	if err != nil {
		if err.Error() == "организация не совпадает" {
			c.JSON(http.StatusForbidden, gin.H{
				"reason": fmt.Sprintf("Пользователь '%s' не принадлежит к организации, связанной с тендером", *request.Query.Username),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if tender.CreatorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не имеет прав на редактирование тендера", *request.Query.Username),
		})
		return
	}

	if request.Parameters.ServiceType != "" {
		if request.Parameters.ServiceType != ds.COSTRUCTION && request.Parameters.ServiceType != ds.DELIVERY && request.Parameters.ServiceType != ds.MANUFACTURE {
			c.JSON(http.StatusBadRequest, gin.H{
				"reason": fmt.Sprintf("Некорректный тип сервиса '%s'", request.Parameters.ServiceType),
			})
			return
		}
		tender.ServiceType = request.Parameters.ServiceType
	}

	if request.Parameters.Name != "" {
		tender.Name = request.Parameters.Name
	}
	if request.Parameters.Description != "" {
		tender.Description = request.Parameters.Description
	}
	tender.Version++
	if err := app.repo.SaveTender(tender); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении тендера",
		})
		return
	}

	c.JSON(http.StatusOK, tender)
}
