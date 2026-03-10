package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveRequestStatus string

const (
	LeaveStatusPending   LeaveRequestStatus = "pending"
	LeaveStatusApproved  LeaveRequestStatus = "approved"
	LeaveStatusRejected  LeaveRequestStatus = "rejected"
	LeaveStatusCancelled LeaveRequestStatus = "cancelled"
)

type LeaveRequest struct {
	ID            uuid.UUID          `json:"id"`
	EmployeeID    uuid.UUID          `json:"employee_id"`
	LeaveTypeID   uuid.UUID          `json:"leave_type_id"`
	StartDate     time.Time          `json:"start_date"`
	EndDate       time.Time          `json:"end_date"`
	TotalDays     int                `json:"total_days"`
	Reason        string             `json:"reason"`
	Status        LeaveRequestStatus `json:"status"`
	ReviewedBy    *uuid.UUID         `json:"reviewed_by,omitempty"`
	ReviewedAt    *time.Time         `json:"reviewed_at,omitempty"`
	ReviewComment string             `json:"review_comment"`
	AttachmentURL string             `json:"attachment_url"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`

	// Relations (populated on demand)
	LeaveType *LeaveType `json:"leave_type,omitempty"`
	Employee  *Employee  `json:"employee,omitempty"`
}
