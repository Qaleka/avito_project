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
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат данных: " + err.Error(),
		})
		return
	}
	fmt.Println(request)
	var bid ds.Bid
	bid.CreatedAt = time.Now()
	if err := app.repo.AddBid(&bid, request); err != nil {
		switch err.Error() {
		case "пользователь не найден":
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": "Пользователь с ID '" + request.AuthorID + "' не найден",
			})
		case "тендер не найден":
			c.JSON(http.StatusNotFound, gin.H{
				"reason": "Тендер с ID '" + request.TenderId + "' не найден",
			})
		// case "организация не совпадает":
		// 	c.JSON(http.StatusForbidden, gin.H{
		// 		"reason": "Пользователь '" + request.CreatorUsername + "' не принадлежит к организации с ID '" + request.OrganizationID + "'",
		// 	})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"reason": err.Error(),
			})
		}
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

func (app *Application) GetBid(c *gin.Context) {
	var request schemes.GetUserBidsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат запроса или его параметры: " + err.Error(),
		})
		return
	}
	bids, err := app.repo.GetUserBids(request.Limit, request.Offset, request.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func (app *Application) GetTenderBids(c *gin.Context) {
	var request schemes.GetTenderBidsRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Неверный формат URI параметров: '%s'", err.Error()),
		})
		return
	}
	if err := c.ShouldBindQuery(&request.Query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат запроса или его параметры: " + err.Error(),
		})
		return
	}
	fmt.Println(request)
	bids, err := app.repo.GetTenderBids(request)
	if err != nil {
		if err.Error() == "пользователь не найден" {
			c.JSON(http.StatusForbidden, gin.H{
				"reason": fmt.Sprintf("Пользователь '%s' не найден", request.Query.Username),
			})
			return
		}
		if err.Error() == "тендер не найден" {
			c.JSON(http.StatusNotFound, gin.H{
				"reason": fmt.Sprintf("Тендер с id '%s' не найден", request.URI.TenderId),
			})
			return
		}
		if err.Error() == "нет права доступа" {
			c.JSON(http.StatusForbidden, gin.H{
				"reason": fmt.Sprintf("У пользователя '%s' нет права доступа", request.Query.Username),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, bids)
}

func (app *Application) GetBidStatus(c *gin.Context) {
	var request schemes.GetBidStatusRequest

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
	status, err := app.repo.GetBidStatus(request.BidId, request.Username)
	if err != nil {
		if err.Error() == "пользователь не найден" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": fmt.Sprintf("Пользователь с именем '%s' не найден", request.Username),
			})
			return
		} else if err.Error() == "нет права доступа" {
			c.JSON(http.StatusForbidden, gin.H{
				"reason": "Пользователь не имеет прав на доступ к предложению",
			})
			return
		} else if err.Error() == "тендер не найден" {
			c.JSON(http.StatusNotFound, gin.H{
				"reason": fmt.Sprintf("Предложение с id '%s' не найдено", request.BidId),
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

func (app *Application) ChangeBidStatus(c *gin.Context) {
	var request schemes.ChangeBidStatusRequest

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
	if request.Parameters.Status != ds.CREATED && request.Parameters.Status != ds.PUBLISHED && request.Parameters.Status != ds.CANCELED {
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

	bid, err := app.repo.GetBidById(request.URI.BidId, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if bid.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не имеет прав на редактирование тендера", request.Parameters.Username),
		})
		return
	}

	bid.Status = request.Parameters.Status
	if err := app.repo.SaveBid(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении тендера",
		})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func (app *Application) ChangeBid(c *gin.Context) {
	var request schemes.ChangeBidRequest

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
	
	bid, err := app.repo.GetBidById(request.URI.BidId, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if bid.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не имеет прав на редактирование предложения", *request.Query.Username),
		})
		return
	}

	if request.Parameters.Name != "" {
		bid.Name = request.Parameters.Name
	}
	if request.Parameters.Description != "" {
		bid.Description = request.Parameters.Description
	}

	if err := app.repo.SaveBid(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении тендера",
		})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func (app *Application) SubmitBid(c *gin.Context) {
	var request schemes.SubmitBidRequest

	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат URI параметров: " + err.Error(),
		})
		return
	}

	if err := c.ShouldBindQuery(&request.Query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат параметров тела запроса: " + err.Error(),
		})
		return
	}

	user, err := app.repo.GetUserByUsername(request.Query.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при проверке пользователя",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не найден", request.Query.Username),
		})
		return
	}

	bid, err := app.repo.SubmitBid(request.URI.BidId, user.ID)
	if err != nil {
		if err.Error() == "нет прав"  {
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": fmt.Sprintf("Пользователь '%s' не имеет прав на решение по предложению", request.Query.Username),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}
	if request.Query.Decision != "Approved" && request.Query.Decision != "Rejected" {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Некоректное решение: '%s'", request.Query.Decision),
		})
		return
	}
	tender := ds.Tender{ID:bid.TenderID}
	if request.Query.Decision == "Approved" {
		bid.Status = ds.PUBLISHED
		tender.Status = ds.CLOSED
	} else if request.Query.Decision == "Rejected" {
		bid.Status = ds.CANCELED
	}
	if err := app.repo.SaveBid(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении тендера",
		})
		return
	}

	if err := app.repo.SaveTender(&tender); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении тендера",
		})
		return
	}
	c.JSON(http.StatusOK, bid)
}