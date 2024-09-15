package repository

import (
	// "avito_project/internal/app/ds"
	// "errors"
	// "strings"
	"fmt"
	"time"
	"gorm.io/gorm"

	"avito_project/internal/app/ds"
	"avito_project/internal/app/schemes"
)

// func (r *Repository) GetAllBids(serviceType string) ([]ds.Bid, error) {
// 	var tenders []ds.Tender
// 	query := r.db.Where("service_type LIKE ?", "%" + strings.ToLower(serviceType) + "%")
// 	if err := query.Find(&tenders).Error; err != nil {
// 		return nil, err
// 	}
// 	return tenders, nil
// }

func (r *Repository) AddBid(bid *ds.Bid, request schemes.AddBidRequest) error {
	var employee ds.Employee
	if err := r.db.Where("id = ?", request.AuthorID).First(&employee).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("пользователь не найден")
		}
		return err
	}

	var tender ds.Tender
	if err := r.db.Where("id = ?", request.TenderId).First(&tender).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("тендер не найден")
		}
		return err
	}

	// var orgResp ds.OrganizationResponsible
	// if err := r.db.Where("organization_id = ? AND user_id = ?", request.OrganizationID, employee.ID).First(&orgResp).Error; err != nil {
	// 	if gorm.ErrRecordNotFound == err {
	// 		return fmt.Errorf("организация не совпадает")
	// 	}
	// 	return err
	// }

	bid.Name = request.Name
	bid.Description = request.Description
	bid.TenderID = request.TenderId
	bid.AuthorType = request.AuthorType
	bid.AuthorID = request.AuthorID
	bid.Status = ds.CREATED
	bid.Version = 1

	if err := r.db.Create(&bid).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveBid(bid *ds.Bid) error {
	err := r.db.Save(bid).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserBids(limit int, offset int, username string) ([]ds.Bid, error) {
	var bids []ds.Bid
	var employee ds.Employee
	err := r.db.Where("username = ?", username).First(&employee).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("пользователь с именем '%s' не найден", username)
		}
		return nil, err
	}
	query := r.db.Where("author_id = ?",employee.ID).Order("name ASC")
	if limit > 0 {
		query = r.db.Limit(limit)
	}
	if offset > 0 {
		query = r.db.Offset(offset)
	}
	if err := query.Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (r *Repository) GetTenderBids(request schemes.GetTenderBidsRequest) ([]ds.Bid, error) {
	var bids []ds.Bid
	var employee ds.Employee
	var tender ds.Tender
	if err := r.db.Where("username = ?", request.Query.Username).First(&employee).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, err
	}

	err := r.db.Where("id = ?", request.URI.TenderId).First(&tender).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("тендер не найден")
		}
		return nil, err
	}
	if employee.ID != tender.CreatorID {
		return nil, fmt.Errorf("нет права доступа")
	}
	query := r.db.Where("tender_id = ?",tender.ID).Order("name ASC")
	if request.Query.Limit > 0 {
		query = r.db.Limit(request.Query.Limit)
	}
	if request.Query.Offset > 0 {
		query = r.db.Offset(request.Query.Offset)
	}
	if err := query.Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (r *Repository) GetBidStatus(bid_id string, username string) (string, error) {
	var status string
	var employee ds.Employee
	var bid ds.Bid

	if username != "" {
		err := r.db.Where("username = ?", username).First(&employee).Error
		if err != nil {
			if gorm.ErrRecordNotFound == err {
				return "", fmt.Errorf("пользователь не найден")
			}
			return "", err
		}
	}

	// Проверка существования тендера по tender_id
	if err := r.db.Where("id = ?", bid_id).First(&bid).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return "", fmt.Errorf("тендер не найден")
		}
		return "", err
	}

	// Проверка прав доступа: если пользователь указан, проверяем, совпадает ли его ID с creator_id тендера
	if username != "" && bid.AuthorID != employee.ID {
		return "", fmt.Errorf("нет права доступа")
	}

	// Получение статуса тендера
	err := r.db.Model(&ds.Bid{}).Select("status").Where("id = ?", bid_id).Row().Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *Repository) GetBidById(bid_id string, user_id string) (*ds.Bid, error) {
	var bid ds.Bid

	if err := r.db.Where("id = ?", bid_id).First(&bid).Error; err != nil {
		return nil, fmt.Errorf("предложение с id '%s' не найден", bid_id)
	}
	
	return &bid, nil
}

func (r *Repository) SubmitBid(bid_id string, user_id string) (*ds.Bid, error) {
	var bid ds.Bid
	var tender ds.Tender

	if err := r.db.Where("id = ?", bid_id).First(&bid).Error; err != nil {
		return nil, fmt.Errorf("предложение с id '%s' не найден", bid_id)
	}

	if err := r.db.Where("id = ?", bid.TenderID).First(&tender).Error; err != nil {
		return nil, fmt.Errorf("тендер с id '%s' не найден", bid_id)
	}

	if tender.CreatorID != user_id {
		return nil, fmt.Errorf("нет прав")
	}
	return &bid, nil
}

func (r *Repository) GetBidNewVersion(bidId string, version int) (*ds.BidVersion, error) {
	var bidVersion ds.BidVersion
	if err := r.db.Where("bid_id = ? AND version = ?", bidId, version).First(&bidVersion).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("версия предложения не найдена")
		}
		return nil, err
	}
	return &bidVersion, nil
}

func (r *Repository) SaveBidVersion(bid *ds.Bid) error {
	BidVersion := ds.BidVersion{
		BidId: bid.ID,
		Version:     bid.Version,
		Name:        bid.Name,
		Description: bid.Description,
		AuthorID: bid.AuthorID,
		Status:		 bid.Status,
		AuthorType: bid.AuthorType,
		TenderId: bid.TenderID,
		CreatedAt:   bid.CreatedAt,
	}
	if err := r.db.Create(&BidVersion).Error; err != nil {
		return err
	}
	err := r.db.Save(BidVersion).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddBidFeedback(bid *ds.Bid, feedBack string, userId string) (*ds.Feedback, error) {
	var bidFeedback ds.Feedback
	var tender ds.Tender
	bidFeedback.BidId = bid.ID
	bidFeedback.Description = feedBack
	if err := r.db.Where("id = ?", bid.TenderID).First(&tender).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("тендер не найден")
		}
		return nil, err
	}
	if tender.CreatorID != userId {
		return nil, fmt.Errorf("не создатель тендера")
	}
	bidFeedback.AuthorID = userId
	bidFeedback.CreatedAt = time.Now()
	return &bidFeedback, nil
}

func (r *Repository) SaveFeedback(bidFeedback *ds.Feedback) error {
	if err := r.db.Create(&bidFeedback).Error; err != nil {
		return err
	}
	err := r.db.Save(bidFeedback).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetBidsByAuthorAndTender(tenderID string, authorId string) ([]ds.Bid, error) {
	var bids []ds.Bid

	query := r.db.Where("tender_id = ? AND author_id = ?", tenderID, authorId)

	if err := query.Find(&bids).Error; err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *Repository) GetReviewsByBid(bidID string, limit int, offset int) ([]ds.Feedback, error) {
	var feedbacks []ds.Feedback

	// Строим базовый запрос для поиска отзывов по BidId
	query := r.db.Where("bid_id = ?", bidID)

	// Применяем лимит, если он указан
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Применяем смещение, если оно указано
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Выполняем запрос и сохраняем результат в срез feedbacks
	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}