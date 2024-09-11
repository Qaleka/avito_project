package ds

import (
	"time"
)

type Employee struct {
	ID        string      `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	FirstName string    `gorm:"size:50" json:"first_name"`
	LastName  string    `gorm:"size:50" json:"last_name"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
}

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

type Organization struct {
	ID          string            `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"size:100;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Type        OrganizationType `gorm:"type:string" json:"type"`
	CreatedAt   time.Time       `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"type:timestamp" json:"updated_at"`
}

type OrganizationResponsible struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	OrganizationID string   `gorm:"not null" json:"organization_id"`
	UserID        string   `gorm:"not null" json:"user_id"`
	
	// Связи
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	User         Employee     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
}

const CREATED string = "Created"
const PUBLISHED string = "Published"
const CLOSED string = "Closed"

type Tender struct {
	ID             string         `gorm:"primaryKey;default:gen_random_uuid()" json:"id" binding:"-"`
	Name           string       `gorm:"size:100;not null" form:"name" json:"name" binding:"-"`
	Description    string       `gorm:"type:text" form:"description" json:"description" binding:"-"`
	ServiceType    string  `gorm:"type:string;not null" form:"service_type" json:"service_type" binding:"-"`
	Status         string `gorm:"type:string;not null" form:"status" json:"status" binding:"-"`
	OrganizationID string         `gorm:"not null" form:"organization_id" json:"organization_id" binding:"-"`
	CreatorUsername string      `gorm:"size:50;not null" form:"creator_username" json:"creator_username" binding:"-"`
	CreatedAt   time.Time       `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"type:timestamp" json:"updated_at"`
	// Связь с организацией
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"-" binding:"-"`
}

type ProposalStatus string

const (
	Submitted ProposalStatus = "Submitted"
	Accepted  ProposalStatus = "Accepted"
	Rejected  ProposalStatus = "Rejected"
)

type Bid struct {
	ID             string            `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string          `gorm:"size:100;not null" json:"name"`
	Description    string          `gorm:"type:text" json:"description"`
	Status         ProposalStatus  `gorm:"type:string;not null" json:"status"`
	TenderID       string            `gorm:"not null" json:"tender_id"`
	OrganizationID string            `gorm:"not null" json:"organization_id"`
	CreatorUsername string         `gorm:"size:50;not null" json:"creator_username"`
	CreatedAt   time.Time       `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"type:timestamp" json:"updated_at"`

	// Связь с тендером
	Tender        Tender       `gorm:"foreignKey:TenderID;constraint:OnDelete:CASCADE" json:"-"`
	Organization  Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"-"`
}