package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CorCompanyDetails struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	CompanyName string     `json:"company_name"`
	TPIN        string     `json:"tpin,omitempty"`
	RegNumber   string     `json:"reg_number,omitempty"`
	Industry    string     `json:"industry,omitempty"`
	Country     string     `json:"country,omitempty"`
	Status      string     `json:"status"`
	MetaData    json.RawMessage   `json:"meta_data,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CorBranchDetails struct {
	ID        uuid.UUID  `json:"id"`
	CompanyID uuid.UUID  `json:"company_id"`
	Name      string     `json:"name"`
	Address   string     `json:"address,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	MetaData  json.RawMessage   `json:"meta_data,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateCorCompanyRequest struct {
	CompanyName string   `json:"company_name"`
	TPIN        string   `json:"tpin,omitempty"`
	RegNumber   string   `json:"reg_number,omitempty"`
	Industry    string   `json:"industry,omitempty"`
	Country     string   `json:"country,omitempty"`
}

type UpdateCorCompanyRequest struct {
	CompanyName *string `json:"company_name,omitempty"`
	TPIN        *string `json:"tpin,omitempty"`
	Industry    *string `json:"industry,omitempty"`
	Country     *string `json:"country,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type CreateCorBranchRequest struct {
	Name    string `json:"name"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

type UpdateCorBranchRequest struct {
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
	Phone   *string `json:"phone,omitempty"`
}
