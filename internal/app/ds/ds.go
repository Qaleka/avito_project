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
}

const IE string = "IE"
const LLC string = "LLC"
const JSC string = "JSC"

type Organization struct {
	ID          string            `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"size:100;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Type        string `gorm:"type:string" json:"type"`
	CreatedAt   time.Time       `gorm:"type:timestamp" json:"created_at"`
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
const CANCELED string = "Canceled"

const COSTRUCTION string = "Construction"
const DELIVERY string = "Delivery"
const MANUFACTURE string = "Manufacture"

type Tender struct {
	ID             string         `gorm:"primaryKey;default:gen_random_uuid()" json:"id" binding:"-"`
	Name           string       `gorm:"size:100;not null" form:"name" json:"name" binding:"-"`
	Description    string       `gorm:"type:text" form:"description" json:"description" binding:"-"`
	ServiceType    string  `gorm:"type:string;not null" form:"service_type" json:"serviceType" binding:"-"`
	Status         string `gorm:"type:string;not null" form:"status" json:"status" binding:"-"`
	OrganizationID string         `gorm:"not null" form:"organization_id" json:"organizationId" binding:"-"`
	Version 		int `gorm:"not null" form:"version" json:"version" binding:"-"`
	CreatorID string      `gorm:"size:50;not null" form:"creator_id" json:"-" binding:"-"`
	CreatedAt   time.Time       `gorm:"type:timestamp" json:"created_at"`

	// Связь с организацией
	Employee Employee `gorm:"foreignKey:CreatorID" json:"-" binding:"-"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-" binding:"-"`
}


type Bid struct {
	ID             string            `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string          `gorm:"size:100;not null" json:"name"`
	Description    string          `gorm:"type:text" json:"description"`
	Status         string  `gorm:"type:string;not null" json:"status"`
	TenderID       string            `gorm:"not null" json:"-"`
	AuthorType     string            `gorm:"not null" json:"authorType"`
	AuthorID       string            `gorm:"not null" json:"authorId"`
	Version 		int `gorm:"not null" form:"version" json:"version" binding:"-"`
	CreatedAt   time.Time       `gorm:"type:timestamp" json:"created_at"`

	// Связь с тендером
	Employee        Employee       `gorm:"foreignKey:AuthorID" json:"-"`
	Tender  Tender `gorm:"foreignKey:TenderID" json:"-"`
}