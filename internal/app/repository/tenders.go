package repository

import (
	// "avito_project/internal/app/ds"
	"errors"
	// "strings"
	"fmt"
	"gorm.io/gorm"

	"avito_project/internal/app/ds"
)

func (r *Repository) GetAllTenders(serviceType string, limit int, offset int) ([]ds.Tender, error) {
	var tenders []ds.Tender
	if serviceType != "Construction" && serviceType != "Delivery" && serviceType != "Manufacture" {
		return nil, fmt.Errorf("введен некоректный сервис: %s", serviceType)
	}
	query := r.db.Where("service_type LIKE ?", "%" + serviceType + "%").Order("name ASC")
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

func (r *Repository) AddTender(tender *ds.Tender) error {
	var employee ds.Employee
	err := r.db.Where("username = ?", tender.CreatorUsername).First(&employee).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("пользователь с именем '%s' не найден", tender.CreatorUsername)
		}
		return err
	}

	err = r.db.Create(&tender).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveTender(tender *ds.Tender) error {
	err := r.db.Save(tender).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetUserTenders(limit int, offset int, username string) ([]ds.Tender, error) {
	var tenders []ds.Tender
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

func (r *Repository) GetTenderStatus(tender_id string, username string) (string, error) {
	var status string
	var employee ds.Employee
	err := r.db.Where("username = ?", username).First(&employee).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return "", fmt.Errorf("пользователь с именем '%s' не найден", username)
		}
		return "", err
	}
	query := r.db.Model(&ds.Tender{}).Select("status").Where("id = ?", tender_id)
	if username != "" {
		query = query.Where("creator_username = ?", username)
	}
	if err := query.Row().Scan(&status); err != nil {
		return "", fmt.Errorf("тендер с id '%s' не найден", tender_id)
	}

	return status, nil
}

func (r *Repository) GetTenderById(tender_id string, username string) (*ds.Tender, error) {
	tender := &ds.Tender{ID: tender_id}
	err := r.db.First(tender, "creator_username = ?", username).Error
	if err != nil {
		return nil, fmt.Errorf("тендер с id '%s' не найден", tender_id)
	}
	return tender, nil
}

func (r *Repository) GetUserByUsername(username string) (*ds.Employee, error) {
	var user ds.Employee

	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}