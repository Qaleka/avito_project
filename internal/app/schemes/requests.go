package schemes

// import "avito_project/internal/app/ds"

// "avito_project/internal/app/ds"

// "mime/multipart"
// "time"

type GetAllTendersRequest struct {
	ServiceType string `form:"service_type"`
	Limit int `form:"limit"`
	Offset int `form:"offset"`
}

type AddTenderRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	ServiceType string `form:"servicetype" json:"servicetype" binding:"required"`
	OrganizationID string `form:"organizationId" json:"organizationId" binding:"required"`
	CreatorUsername string `form:"creatorUsername" json:"creatorUsername" binding:"required"`
}


type GetUserTendersRequest struct {
	Username string `form:"username" binding:"required"`
	Limit int `form:"limit"`
	Offset int `form:"offset"`
}

type GetTenderStatusRequest struct {
	TenderId string `uri:"tenderId" binding:"required"`
	Username string `form:"username"`
}

type ChangeTenderStatusRequest struct {
	URI struct {
		TenderId string `uri:"tenderId" binding:"required"`
	}
	Parameters struct {
		Status string `form:"status" binding:"required"`
		Username string `form:"username" binding:"required"`
	}
}

type ChangeTenderRequest struct {
	Parameters struct {
		Name string `form:"name" json:"name" binding:"omitempty,max=75"`
		Description string `form:"description" json:"description" binding:"omitempty,max=75"`
		ServiceType string `form:"serviceType" json:"serviceType" binding:"omitempty,max=75"`

	}
	Query struct {
		Username *string `form:"username" json:"username" binding:"required"`
	}
	URI struct {
		TenderId string `uri:"tenderId" binding:"required"`
	}
}

type AddBidRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	TenderId string `form:"tender_id" json:"tenderId" binding:"required"`
	AuthorType string `form:"authorType" json:"authorType" binding:"required"`
	AuthorID string `form:"authorId" json:"authorId" binding:"required"`
}

type GetUserBidsRequest struct {
	Username string `form:"username" binding:"required"`
	Limit int `form:"limit"`
	Offset int `form:"offset"`
}

type GetTenderBidsRequest struct {
	URI struct {
		TenderId string `uri:"bidId" binding:"required"`
	}
	Query struct {
		Username string `form:"username" binding:"required"`
		Limit int `form:"limit"`
		Offset int `form:"offset"`
	}
}

type GetBidStatusRequest struct {
	BidId string `uri:"bidId" binding:"required"`
	Username string `form:"username"`
}

type ChangeBidStatusRequest struct {
	URI struct {
		BidId string `uri:"bidId" binding:"required"`
	}
	Parameters struct {
		Status string `form:"status" binding:"required"`
		Username string `form:"username" binding:"required"`
	}
}

type ChangeBidRequest struct {
	Parameters struct {
		Name string `form:"name" json:"name" binding:"omitempty,max=75"`
		Description string `form:"description" json:"description" binding:"omitempty,max=75"`
	}
	Query struct {
		Username *string `form:"username" json:"username" binding:"required"`
	}
	URI struct {
		BidId string `uri:"bidId" binding:"required"`
	}
}

type SubmitBidRequest struct {
	URI struct {
		BidId string `uri:"bidId" binding:"required"`
	}
	Query struct {
		Decision string `form:"decision" binding:"required"`
		Username string `form:"username" binding:"required"`
	}
}