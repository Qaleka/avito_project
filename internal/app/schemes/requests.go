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


type GetUserTendersRequest struct {
	Username string `form:"username"`
	Limit int `form:"limit"`
	Offset int `form:"offset"`
}

type GetTenderStatusRequest struct {
	TenderId string `uri:"tender_id" binding:"required"`
	Username string `form:"username"`
}

type ChangeTenderStatusRequest struct {
	Parameters struct {
		Status *string `form:"status" json:"status" binding:"required"`
		Username *string `form:"username" json:"username" binding:"required"`
	}
	TenderId string `uri:"tender_id" binding:"required,uuid"`
}

type ChangeTenderRequest struct {
	Parameters struct {
		Name *string `form:"name" json:"name" binding:"omitempty,max=75"`
		Description *string `form:"description" json:"description" binding:"omitempty,max=75"`
		ServiceType *string `form:"service_type" json:"service_type" binding:"omitempty,max=75"`

	}
	Query struct {
		Username *string `form:"username" json:"username" binding:"required"`
	}
	TenderId string `uri:"tender_id" binding:"required"`
}

type AddBidRequest struct {
	ds.Bid
}