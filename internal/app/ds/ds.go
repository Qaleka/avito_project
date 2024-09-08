package ds

import (
	"time"
)

type Employee struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	FirstName string    `gorm:"size:50" json:"first_name"`
	LastName  string    `gorm:"size:50" json:"last_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

type Organization struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"size:100;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Type        OrganizationType `gorm:"type:organization_type" json:"type"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type OrganizationResponsible struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	OrganizationID uint   `gorm:"not null" json:"organization_id"`
	UserID        uint   `gorm:"not null" json:"user_id"`
	
	// Связи
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	User         Employee     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

type TenderStatus string

const (
	Open       TenderStatus = "Open"
	Closed     TenderStatus = "Closed"
	InProgress TenderStatus = "InProgress"
)

type Tender struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	Name           string       `gorm:"size:100;not null" json:"name"`
	Description    string       `gorm:"type:text" json:"description"`
	ServiceType    string  `gorm:"type:string;not null" json:"service_type"`
	Status         TenderStatus `gorm:"type:string;not null" json:"status"`
	OrganizationID uint         `gorm:"not null" json:"organization_id"`
	CreatorUsername string      `gorm:"size:50;not null" json:"creator_username"`

	// Связь с организацией
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
}

type ProposalStatus string

const (
	Submitted ProposalStatus = "Submitted"
	Accepted  ProposalStatus = "Accepted"
	Rejected  ProposalStatus = "Rejected"
)

type Proposal struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	Name           string          `gorm:"size:100;not null" json:"name"`
	Description    string          `gorm:"type:text" json:"description"`
	Status         ProposalStatus  `gorm:"type:string;not null" json:"status"`
	TenderID       uint            `gorm:"not null" json:"tender_id"`
	OrganizationID uint            `gorm:"not null" json:"organization_id"`
	CreatorUsername string         `gorm:"size:50;not null" json:"creator_username"`

	// Связь с тендером
	Tender        Tender       `gorm:"foreignKey:TenderID;constraint:OnDelete:CASCADE" json:"tender"`
	Organization  Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
}