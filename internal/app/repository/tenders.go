package repository

import (
	// "avito_project/internal/app/ds"
	// "errors"
	// "strings"

	// "gorm.io/gorm"

	"avito_project/internal/app/ds"
)

func (r *Repository) GetAllTenders(serviceType string, limit int, offset int) ([]ds.Tender, error) {
	var tenders []ds.Tender
	query := r.db.Where("service_type LIKE ?", "%" + serviceType + "%")
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
	err := r.db.Create(&tender).Error
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