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
		case "предложение не найдено":
			c.JSON(http.StatusNotFound, gin.H{
				"reason": "Предложение с ID '" + request.TenderId + "' не найден",
			})
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

	if err := app.repo.SaveBidVersion(&bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении версии предложения: " + err.Error(),
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
		if err.Error() == "предложение не найдено" {
			c.JSON(http.StatusNotFound, gin.H{
				"reason": fmt.Sprintf("Предложение с id '%s' не найдено", request.URI.TenderId),
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
		} else if err.Error() == "предложение не найдено" {
			c.JSON(http.StatusNotFound, gin.H{
				"reason": fmt.Sprintf("Предложение с id '%s' не найдено", request.BidId),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (app *Application) ChangeBidStatus(c *gin.Context) {
	var request schemes.ChangeBidStatusRequest

	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Неверный формат URI параметров: '%s'", err.Error()),
		})
		return
	}

	if err := c.ShouldBindQuery(&request.Parameters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Неверный формат параметров запроса: '%s'", err.Error()),
		})
		return
	}

	if request.Parameters.Status != ds.CREATED && request.Parameters.Status != ds.PUBLISHED && request.Parameters.Status != ds.CANCELED {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Некорректный статус '%s'", request.Parameters.Status),
		})
		return
	}

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
			"reason": fmt.Sprintf("Пользователь '%s' не имеет прав на редактирование предложения", request.Parameters.Username),
		})
		return
	}

	bid.Status = request.Parameters.Status
	if err := app.repo.SaveBid(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении предложения",
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
	bid.Version++
	if err := app.repo.SaveBid(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении предложения",
		})
		return
	}

	if err := app.repo.SaveBidVersion(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении версии предложения: " + err.Error(),
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
	tender, err := app.repo.GetTenderById(bid.TenderID, bid.AuthorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if tender.Status == ds.CLOSED {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf("Тендер уже закрыт"),
		})
		return
	}

	if bid.Status == ds.CANCELED || bid.Status == ds.PUBLISHED {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Предложение уже рассмотрено",
		})
	}

	if request.Query.Decision == "Approved" {
		bid.Status = ds.PUBLISHED
		tender.Status = ds.CLOSED
	} else if request.Query.Decision == "Rejected" {
		bid.Status = ds.CANCELED
	}
	if err := app.repo.SaveBid(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении предложения",
		})
		return
	}

	if err := app.repo.SaveBidVersion(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении версии предложения: " + err.Error(),
		})
		return
	}

	if err := app.repo.SaveTender(tender); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении тендера",
		})
		return
	}
	c.JSON(http.StatusOK, bid)
}

func (app *Application) ChangeBidVersion(c *gin.Context) {
	var request schemes.ChangeBidVersionRequest

	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат URI параметров: " + err.Error(),
		})
		return
	}

	if err := c.ShouldBindQuery(&request.Query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат Query параметров: " + err.Error(),
		})
		return
	}


	user, err := app.repo.GetUserByUsername(request.Query.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
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


	bidVersion, err := app.repo.GetBidNewVersion(request.URI.BidId, request.URI.Version)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": fmt.Sprintf("Предложение с id '%s' версии '%d' не найдено", request.URI.BidId, request.URI.Version),
		})
		return
	}

	bid, err := app.repo.GetBidById(request.URI.BidId, user.ID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"reason": err,
		})
		return
	}

	if bid.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"reason": "Пользователь не имеет доступа к этому предложению",
		})
		return
	}

	bid.Name = bidVersion.Name
	bid.Description = bidVersion.Description
	bid.AuthorID = bidVersion.AuthorID
	bid.AuthorType = bidVersion.AuthorType
	bid.Status = bidVersion.Status
	bid.Version++

	if err := app.repo.SaveBid(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Не удалось сохранить изменения предложения",
		})
		return
	}

	if err := app.repo.SaveBidVersion(bid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при сохранении версии предложения: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func (app *Application) AddBidFeedback(c *gin.Context) {
	var request schemes.AddBidFeedback

	
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат URI параметров: " + err.Error(),
		})
		return
	}

	if err := c.ShouldBindQuery(&request.Query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат Query параметров: " + err.Error(),
		})
		return
	}

	user, err := app.repo.GetUserByUsername(request.Query.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
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

	bid, err := app.repo.GetBidById(request.URI.BidId, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err,
		})
		return
	}

	feedback,err := app.repo.AddBidFeedback(bid, request.Query.BidFeedback, user.ID)
	if err != nil {
		if err.Error() == "не создатель тендера" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": fmt.Sprintf("Пользователь '%s' не является автором тендера, на который пишется отвзы", user.Username),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err,
		})
		return
	}
	err = app.repo.SaveFeedback(feedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": err,
		})
		return
	}
	c.JSON(http.StatusOK, bid)
}

func (app *Application) GetReviews(c *gin.Context) {
	var request schemes.GetReviewsRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат URI параметров: " + err.Error(),
		})
		return
	}
	fmt.Println(request)
	if err := c.ShouldBindQuery(&request.Query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Неверный формат Query параметров: " + err.Error(),
		})
		return
	}
	fmt.Println(request)
	requesterUser, err := app.repo.GetUserByUsername(request.Query.RequesterUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при проверке пользователя: " + err.Error(),
		})
		return
	}
	if requesterUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не найден", request.Query.RequesterUsername),
		})
		return
	}


	tender, err := app.repo.GetTenderById(request.URI.TenderId, requesterUser.ID)
	if err != nil {
		if err.Error() == "организация не совпадает" {
			c.JSON(http.StatusForbidden, gin.H{
				"reason": fmt.Sprintf("Пользователь '%s' не принадлежит к организации, связанной с тендером", request.Query.RequesterUsername),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if tender.CreatorID != requesterUser.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"reason": fmt.Sprintf("Пользователь '%s' не имеет прав на просмотр отзывов для этого тендера", request.Query.RequesterUsername),
		})
		return
	}

	authorUser, err := app.repo.GetUserByUsername(request.Query.AuthorUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при проверке автора отзывов: " + err.Error(),
		})
		return
	}

	if authorUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": fmt.Sprintf("Автор '%s' не найден", request.Query.AuthorUsername),
		})
		return
	}

	bids, err := app.repo.GetBidsByAuthorAndTender(tender.ID, authorUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"reason": "Ошибка при получении предложений: " + err.Error(),
		})
		return
	}

	var reviews []ds.Feedback
	for _, bid := range bids {
		bidReviews, err := app.repo.GetReviewsByBid(bid.ID, request.Query.Limit, request.Query.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"reason": "Ошибка при получении отзывов: " + err.Error(),
			})
			return
		}
		reviews = append(reviews, bidReviews...)
	}

	c.JSON(http.StatusOK, reviews)
}