package repository

import (
	"database/sql"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type BookingDocumentRepository struct {
	db *sql.DB
}

func NewBookingDocumentRepository() *BookingDocumentRepository {
	return &BookingDocumentRepository{db: database.DB}
}

// Upsert inserts or replaces the document URLs for a corporate client.
func (r *BookingDocumentRepository) Upsert(corporateClientID, orgID uuid.UUID, urls []string) (*models.BookingDocument, error) {
	if urls == nil {
		urls = []string{}
	}
	now := time.Now()
	var doc models.BookingDocument
	err := r.db.QueryRow(`
		INSERT INTO booking_documents (corporate_client_id, org_id, urls, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		ON CONFLICT (corporate_client_id, org_id)
		DO UPDATE SET urls = EXCLUDED.urls, updated_at = EXCLUDED.updated_at
		RETURNING id, corporate_client_id, org_id, urls, created_at, updated_at`,
		corporateClientID, orgID, pq.Array(urls), now,
	).Scan(&doc.ID, &doc.CorporateClientID, &doc.OrgID, pq.Array(&doc.URLs), &doc.CreatedAt, &doc.UpdatedAt)
	return &doc, err
}

// GetByCorporateClientID returns the document record for a corporate client, or nil if none.
func (r *BookingDocumentRepository) GetByCorporateClientID(corporateClientID, orgID uuid.UUID) (*models.BookingDocument, error) {
	var doc models.BookingDocument
	err := r.db.QueryRow(`
		SELECT id, corporate_client_id, org_id, urls, created_at, updated_at
		FROM booking_documents
		WHERE corporate_client_id = $1 AND org_id = $2`,
		corporateClientID, orgID,
	).Scan(&doc.ID, &doc.CorporateClientID, &doc.OrgID, pq.Array(&doc.URLs), &doc.CreatedAt, &doc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &doc, err
}
