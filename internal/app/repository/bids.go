package repository

import (
	// "avito_project/internal/app/ds"
	// "errors"
	// "strings"
	"fmt"
	"gorm.io/gorm"

	"avito_project/internal/app/ds"
)

// func (r *Repository) GetAllBids(serviceType string) ([]ds.Bid, error) {
// 	var tenders []ds.Tender
// 	query := r.db.Where("service_type LIKE ?", "%" + strings.ToLower(serviceType) + "%")
// 	if err := query.Find(&tenders).Error; err != nil {
// 		return nil, err
// 	}
// 	return tenders, nil
// }

func (r *Repository) AddBid(bid *ds.Bid) error {
	var employee ds.Employee
	var tender ds.Tender
	err := r.db.Where("username = ?", bid.CreatorUsername).First(&employee).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("пользователь с именем '%s' не найден", bid.CreatorUsername)
		}
		return err
	}

	err = r.db.Where("organization_id = ?", bid.OrganizationID).First(&tender).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("тендер с id '%s' не найден", bid.OrganizationID)
		}
		return err
	}

	err = r.db.Create(&bid).Error
	if err != nil {
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
	var tenders []ds.Bid
	var employee ds.Employee
	err := r.db.Where("username = ?", username).First(&employee).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("пользователь с именем '%s' не найден", username)
		}
		return nil, err
	}
	query := r.db.Where("creator_username = ?",username).Order("name ASC")
	if limit > 0 {
		query = r.db.Limit(limit)
	}
	if offset > 0 {
		query = r.db.Offset(offset)
	}
	if err := query.Find(&tenders).Error; err != nil {
		return nil, err
	}
	return tenders, nil
}