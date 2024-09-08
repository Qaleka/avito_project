package repository

import (
	// "avito_project/internal/app/ds"
	// "errors"
	"strings"

	// "gorm.io/gorm"

	"avito_project/internal/app/ds"
)

func (r *Repository) GetAllTenders(serviceType string) ([]ds.Tender, error) {
	var tenders []ds.Tender
	query := r.db.Where("service_type LIKE ?", "%" + strings.ToLower(serviceType) + "%")
	if err := query.Find(&tenders).Error; err != nil {
		return nil, err
	}
	return tenders, nil
}

