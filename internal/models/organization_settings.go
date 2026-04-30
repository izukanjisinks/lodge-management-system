package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationSettings struct {
	ID                   uuid.UUID  `json:"id"`
	OrgID                *uuid.UUID `json:"org_id,omitempty"`
	AutoCloseOrders      bool       `json:"auto_close_orders"`
	AutoExtendCheckout   bool       `json:"auto_extend_checkout"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type UpdateOrganizationSettingsRequest struct {
	AutoCloseOrders    *bool `json:"auto_close_orders"`
	AutoExtendCheckout *bool `json:"auto_extend_checkout"`
}
