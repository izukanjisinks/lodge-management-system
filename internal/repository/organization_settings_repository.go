package repository

import (
	"database/sql"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type OrganizationSettingsRepository struct {
	db *sql.DB
}

func NewOrganizationSettingsRepository() *OrganizationSettingsRepository {
	return &OrganizationSettingsRepository{db: database.DB}
}

// GetForOrg returns the settings for a specific org, falling back to the system default
// (org_id IS NULL) if no org-specific row exists yet.
func (r *OrganizationSettingsRepository) GetForOrg(orgID uuid.UUID) (*models.OrganizationSettings, error) {
	var s models.OrganizationSettings
	var oid uuid.NullUUID

	err := r.db.QueryRow(`
		SELECT id, org_id, auto_close_orders, auto_extend_checkout, created_at, updated_at
		FROM organization_settings
		WHERE org_id = $1 OR org_id IS NULL
		ORDER BY org_id NULLS LAST
		LIMIT 1`, orgID,
	).Scan(&s.ID, &oid, &s.AutoCloseOrders, &s.AutoExtendCheckout, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if oid.Valid {
		s.OrgID = &oid.UUID
	}
	return &s, nil
}

// Upsert creates or updates the org-scoped settings row.
// The system default row (org_id IS NULL) is never touched by this method.
func (r *OrganizationSettingsRepository) Upsert(orgID uuid.UUID, req *models.UpdateOrganizationSettingsRequest) (*models.OrganizationSettings, error) {
	now := time.Now()

	// Resolve current effective values to use as defaults on first insert
	current, err := r.GetForOrg(orgID)
	if err != nil {
		return nil, err
	}

	autoCloseOrders := current.AutoCloseOrders
	autoExtendCheckout := current.AutoExtendCheckout

	if req.AutoCloseOrders != nil {
		autoCloseOrders = *req.AutoCloseOrders
	}
	if req.AutoExtendCheckout != nil {
		autoExtendCheckout = *req.AutoExtendCheckout
	}

	var s models.OrganizationSettings
	var oid uuid.NullUUID

	err = r.db.QueryRow(`
		INSERT INTO organization_settings (org_id, auto_close_orders, auto_extend_checkout, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		ON CONFLICT (org_id) DO UPDATE
		    SET auto_close_orders    = EXCLUDED.auto_close_orders,
		        auto_extend_checkout = EXCLUDED.auto_extend_checkout,
		        updated_at           = EXCLUDED.updated_at
		RETURNING id, org_id, auto_close_orders, auto_extend_checkout, created_at, updated_at`,
		orgID, autoCloseOrders, autoExtendCheckout, now,
	).Scan(&s.ID, &oid, &s.AutoCloseOrders, &s.AutoExtendCheckout, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if oid.Valid {
		s.OrgID = &oid.UUID
	}
	return &s, nil
}

// ListEnabledForJob returns the org IDs where a given job flag is TRUE.
// Falls back to the system default for orgs that have no custom row.
// jobColumn must be a trusted internal constant — never user input.
func (r *OrganizationSettingsRepository) ListEnabledOrgsForJob(jobColumn string) ([]uuid.UUID, error) {
	rows, err := r.db.Query(`
		SELECT o.id
		FROM organizations o
		WHERE o.is_active = TRUE
		  AND COALESCE(
		      (SELECT s.`+jobColumn+` FROM organization_settings s WHERE s.org_id = o.id),
		      (SELECT s.`+jobColumn+` FROM organization_settings s WHERE s.org_id IS NULL)
		  ) = TRUE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
