package schemes

import (
	"avito_project/internal/app/ds"
)

type TenderOutput struct {
	Tenders []ds.Tender `json:"tenders"`
}