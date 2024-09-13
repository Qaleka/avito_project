package repository

import (
	// "avito_project/internal/app/ds"
	"errors"
	// "strings"
	"fmt"

	"gorm.io/gorm"

	"avito_project/internal/app/ds"
	"avito_project/internal/app/schemes"
)

func (r *Repository) GetAllTenders(serviceType string, limit int, offset int) ([]ds.Tender, error) {
	var tenders []ds.Tender
	if serviceType != ds.COSTRUCTION && serviceType != ds.DELIVERY && serviceType != ds.MANUFACTURE && serviceType!= "" {
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

func (r *Repository) AddTender(tender *ds.Tender, request schemes.AddTenderRequest) error {
	var employee ds.Employee
	if err := r.db.Where("username = ?", request.CreatorUsername).First(&employee).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("пользователь не найден")
		}
		return err
	}

	var organization ds.Organization
	if err := r.db.Where("id = ?", request.OrganizationID).First(&organization).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("организация не найдена")
		}
		return err
	}

	var orgResp ds.OrganizationResponsible
	if err := r.db.Where("organization_id = ? AND user_id = ?", request.OrganizationID, employee.ID).First(&orgResp).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return fmt.Errorf("организация не совпадает")
		}
		return err
	}

	tender.Name = request.Name
	tender.Description = request.Description
	tender.ServiceType = request.ServiceType
	tender.OrganizationID = request.OrganizationID
	tender.CreatorID = employee.ID
	tender.Status = ds.CREATED
	tender.Version = 1

	if err := r.db.Create(&tender).Error; err != nil {
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
	query := r.db.Where("creator_id = ?",employee.ID).Order("name ASC")
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
	var tender ds.Tender

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
	if err := r.db.Where("id = ?", tender_id).First(&tender).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return "", fmt.Errorf("тендер не найден")
		}
		return "", err
	}

	// Проверка прав доступа: если пользователь указан, проверяем, совпадает ли его ID с creator_id тендера
	if username != "" && tender.CreatorID != employee.ID {
		return "", fmt.Errorf("нет права доступа")
	}

	// Получение статуса тендера
	err := r.db.Model(&ds.Tender{}).Select("status").Where("id = ?", tender_id).Row().Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *Repository) GetTenderById(tender_id string, user_id string) (*ds.Tender, error) {
	var tender ds.Tender

	if err := r.db.Where("id = ?", tender_id).First(&tender).Error; err != nil {
		return nil, fmt.Errorf("тендер с id '%s' не найден", tender_id)
	}

	var orgResp ds.OrganizationResponsible
	if err := r.db.Where("organization_id = ? AND user_id = ?", tender.OrganizationID, user_id).First(&orgResp).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("организация не совпадает")
		}
		return nil, err
	}

	return &tender, nil
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

func (r *Repository) GetTenderNewVersion(tenderId string, version int) (*ds.TenderVersion, error) {
	var tenderVersion ds.TenderVersion
	if err := r.db.Where("tender_id = ? AND version = ?", tenderId, version).First(&tenderVersion).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, fmt.Errorf("версия тендера не найдена")
		}
		return nil, err
	}
	return &tenderVersion, nil
}

func (r *Repository) SaveTenderVersion(tender *ds.Tender) error {
	tenderVersion := ds.TenderVersion{
		TenderId: tender.ID,
		Version:     tender.Version,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:		 tender.Status,
		OrganizationID: tender.OrganizationID,
		CreatorID: tender.CreatorID,
		CreatedAt:   tender.CreatedAt,
	}
	if err := r.db.Create(&tenderVersion).Error; err != nil {
		return err
	}
	err := r.db.Save(tenderVersion).Error
	if err != nil {
		return err
	}
	return nil
}