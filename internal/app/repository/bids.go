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
	err := r.db.Where("username = ?", bid.CreatorUsername).First(&employee).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("пользователь с именем '%s' не найден", bid.CreatorUsername)
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