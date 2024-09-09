package schemes

import "avito_project/internal/app/ds"

// "avito_project/internal/app/ds"

// "mime/multipart"
// "time"

type GetAllTendersRequest struct {
	ServiceType string `form:"service_type"`
	Limit int `form:"limit"`
	Offset int `form:"offset"`
}

type AddTenderRequest struct {
	ds.Tender
}