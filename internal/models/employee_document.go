package models

import (
	"time"

	"github.com/google/uuid"
)

type DocumentType string

const (
	DocumentTypeContract    DocumentType = "contract"
	DocumentTypeIDDocument  DocumentType = "id_document"
	DocumentTypeCertification DocumentType = "certification"
	DocumentTypeOfferLetter DocumentType = "offer_letter"
	DocumentTypeWarningLetter DocumentType = "warning_letter"
	DocumentTypeOther       DocumentType = "other"
)

type EmployeeDocument struct {
	ID           uuid.UUID    `json:"id"`
	EmployeeID   uuid.UUID    `json:"employee_id"`
	DocumentType DocumentType `json:"document_type"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	FileURL      string       `json:"file_url"`
	FileName     string       `json:"file_name"`
	FileSize     int64        `json:"file_size"`
	MimeType     string       `json:"mime_type"`
	UploadedBy   uuid.UUID    `json:"uploaded_by"`
	ExpiryDate   *time.Time   `json:"expiry_date,omitempty"`
	IsVerified   bool         `json:"is_verified"`
	VerifiedBy   *uuid.UUID   `json:"verified_by,omitempty"`
	VerifiedAt   *time.Time   `json:"verified_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty"`
}
