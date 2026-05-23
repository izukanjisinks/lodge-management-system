package models

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID         uuid.UUID `json:"id"`
	OrgID      uuid.UUID `json:"org_id"`
	Name       string    `json:"name"`
	BranchCode string    `json:"branch_code"`
	Address    string    `json:"address,omitempty"`
	Location   string    `json:"location,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateBranchRequest struct {
	Name       string `json:"name"`
	BranchCode string `json:"branch_code"`
	Address    string `json:"address,omitempty"`
	Location   string `json:"location,omitempty"`
}

type UpdateBranchRequest struct {
	Name     *string `json:"name,omitempty"`
	Address  *string `json:"address,omitempty"`
	Location *string `json:"location,omitempty"`
}
