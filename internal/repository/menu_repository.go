package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type MenuRepository struct {
	db *sql.DB
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{db: database.DB}
}

// ── Menus ─────────────────────────────────────────────────────────────────────

// GetMenu returns the branch-scoped menu if branchID is set, otherwise the
// org-scoped menu, falling back to the system default (org_id IS NULL).
func (r *MenuRepository) GetMenu(orgID uuid.UUID, branchID *uuid.UUID) (*models.Menu, error) {
	var m models.Menu
	var oid uuid.NullUUID
	var bid uuid.NullUUID
	var description sql.NullString

	if branchID != nil {
		err := r.db.QueryRow(`
			SELECT id, org_id, branch_id, name, description, is_active, created_at, updated_at
			FROM menus
			WHERE org_id = $1 AND branch_id = $2
			LIMIT 1`, orgID, branchID).
			Scan(&m.ID, &oid, &bid, &m.Name, &description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == nil {
			if oid.Valid {
				m.OrgID = oid.UUID
			}
			if bid.Valid {
				m.BranchID = &bid.UUID
			}
			if description.Valid {
				m.Description = description.String
			}
			return &m, nil
		}
		// Fall through to org-level menu if no branch menu exists
	}

	err := r.db.QueryRow(`
		SELECT id, org_id, branch_id, name, description, is_active, created_at, updated_at
		FROM menus
		WHERE (org_id = $1 AND branch_id IS NULL) OR org_id IS NULL
		ORDER BY org_id NULLS LAST
		LIMIT 1`, orgID).
		Scan(&m.ID, &oid, &bid, &m.Name, &description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if oid.Valid {
		m.OrgID = oid.UUID
	}
	if bid.Valid {
		m.BranchID = &bid.UUID
	}
	if description.Valid {
		m.Description = description.String
	}
	return &m, nil
}

// UpsertMenu creates or updates the branch- or org-scoped menu row.
// The system default row (org_id IS NULL) is never touched by this method.
func (r *MenuRepository) UpsertMenu(orgID uuid.UUID, branchID *uuid.UUID, req *models.UpdateMenuRequest) (*models.Menu, error) {
	current, err := r.GetMenu(orgID, branchID)
	if err != nil {
		return nil, err
	}

	name := current.Name
	description := current.Description
	isActive := current.IsActive

	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	now := time.Now()
	var m models.Menu
	var oid, bid uuid.NullUUID
	var desc sql.NullString

	err = r.db.QueryRow(`
		INSERT INTO menus (org_id, branch_id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		ON CONFLICT (org_id, branch_id) DO UPDATE
		    SET name        = EXCLUDED.name,
		        description = EXCLUDED.description,
		        is_active   = EXCLUDED.is_active,
		        updated_at  = EXCLUDED.updated_at
		RETURNING id, org_id, branch_id, name, description, is_active, created_at, updated_at`,
		orgID, branchID, name, description, isActive, now,
	).Scan(&m.ID, &oid, &bid, &m.Name, &desc, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if oid.Valid {
		m.OrgID = oid.UUID
	}
	if bid.Valid {
		m.BranchID = &bid.UUID
	}
	if desc.Valid {
		m.Description = desc.String
	}
	return &m, nil
}

// ── Menu Items ────────────────────────────────────────────────────────────────

func (r *MenuRepository) CreateMenuItem(item *models.MenuItem, menuID uuid.UUID, orgID uuid.UUID, branchID *uuid.UUID) error {
	item.ID = uuid.New()
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	item.MenuID = menuID
	item.OrgID = orgID
	item.BranchID = branchID

	_, err := r.db.Exec(`
		INSERT INTO menu_items (id, menu_id, org_id, branch_id, name, description, category, image_url, price, is_available, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		item.ID, menuID, orgID, branchID, item.Name, item.Description, item.Category, item.ImageURL, item.Price, item.IsAvailable,
		item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *MenuRepository) GetMenuItemByID(id uuid.UUID, orgID uuid.UUID) (*models.MenuItem, error) {
	var item models.MenuItem
	var description, category, imageURL sql.NullString
	var branchID uuid.NullUUID
	err := r.db.QueryRow(`
		SELECT id, menu_id, org_id, branch_id, name, description, category, image_url, price, is_available, created_at, updated_at
		FROM menu_items WHERE id=$1 AND org_id=$2`, id, orgID).
		Scan(&item.ID, &item.MenuID, &item.OrgID, &branchID, &item.Name, &description, &category, &imageURL, &item.Price, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if branchID.Valid {
		item.BranchID = &branchID.UUID
	}
	if description.Valid {
		item.Description = description.String
	}
	if category.Valid {
		item.Category = category.String
	}
	if imageURL.Valid {
		item.ImageURL = &imageURL.String
	}
	return &item, nil
}

// ListAvailableMenuItems returns only available items — used by the guest endpoint.
func (r *MenuRepository) ListAvailableMenuItems(menuID uuid.UUID, category string, page, pageSize int) ([]models.MenuItem, int, error) {
	args := []interface{}{menuID}
	where := "menu_id=$1 AND is_available=TRUE"
	if category != "" {
		args = append(args, category)
		where += fmt.Sprintf(" AND category=$%d", len(args))
	}

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM menu_items WHERE %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, menu_id, org_id, branch_id, name, description, category, image_url, price, is_available, created_at, updated_at
		FROM menu_items WHERE %s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items, err := scanMenuItems(rows)
	return items, total, err
}

func (r *MenuRepository) ListMenuItems(menuID uuid.UUID, orgID uuid.UUID, category string, page, pageSize int) ([]models.MenuItem, int, error) {
	args := []interface{}{menuID, orgID}
	where := "menu_id=$1 AND org_id=$2"
	if category != "" {
		args = append(args, category)
		where += fmt.Sprintf(" AND category=$%d", len(args))
	}

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM menu_items WHERE %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, menu_id, org_id, branch_id, name, description, category, image_url, price, is_available, created_at, updated_at
		FROM menu_items WHERE %s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items, err := scanMenuItems(rows)
	return items, total, err
}

func (r *MenuRepository) UpdateMenuItem(id uuid.UUID, orgID uuid.UUID, req *models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	item, err := r.GetMenuItemByID(id, orgID)
	if err != nil {
		return nil, fmt.Errorf("menu item not found")
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.ImageURL != nil {
		item.ImageURL = req.ImageURL
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}
	item.UpdatedAt = time.Now()

	_, err = r.db.Exec(`
		UPDATE menu_items SET name=$1, description=$2, category=$3, image_url=$4, price=$5, is_available=$6, updated_at=$7
		WHERE id=$8 AND org_id=$9`,
		item.Name, item.Description, item.Category, item.ImageURL, item.Price, item.IsAvailable, item.UpdatedAt, id, orgID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetMenuItemByID(id, orgID)
}

func scanMenuItems(rows *sql.Rows) ([]models.MenuItem, error) {
	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var description, category, imageURL sql.NullString
		var branchID uuid.NullUUID
		if err := rows.Scan(&item.ID, &item.MenuID, &item.OrgID, &branchID, &item.Name, &description, &category, &imageURL, &item.Price, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		if branchID.Valid {
			item.BranchID = &branchID.UUID
		}
		if description.Valid {
			item.Description = description.String
		}
		if category.Valid {
			item.Category = category.String
		}
		if imageURL.Valid {
			item.ImageURL = &imageURL.String
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.MenuItem{}
	}
	return items, rows.Err()
}

func (r *MenuRepository) DeleteMenuItem(id uuid.UUID, orgID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM menu_items WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("menu item not found")
	}
	return nil
}

// GuestGetMenu returns the active menu for an org with only available items — for public display.
// Falls back to the system default menu if the org has no custom menu.
func (r *MenuRepository) GuestGetMenu(orgID uuid.UUID, branchID *uuid.UUID) (*models.Menu, error) {
	var m models.Menu
	var oid uuid.NullUUID
	var bid uuid.NullUUID
	var description sql.NullString

	if branchID != nil {
		err := r.db.QueryRow(`
			SELECT id, org_id, branch_id, name, description, is_active, created_at, updated_at
			FROM menus
			WHERE org_id = $1 AND branch_id = $2 AND is_active = TRUE
			LIMIT 1`, orgID, branchID).
			Scan(&m.ID, &oid, &bid, &m.Name, &description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == nil {
			if oid.Valid {
				m.OrgID = oid.UUID
			}
			if bid.Valid {
				m.BranchID = &bid.UUID
			}
			if description.Valid {
				m.Description = description.String
			}
			return &m, nil
		}
		// Fall through to org-level menu
	}

	err := r.db.QueryRow(`
		SELECT id, org_id, branch_id, name, description, is_active, created_at, updated_at
		FROM menus
		WHERE (org_id = $1 OR org_id IS NULL) AND branch_id IS NULL AND is_active = TRUE
		ORDER BY org_id NULLS LAST
		LIMIT 1`, orgID).
		Scan(&m.ID, &oid, &bid, &m.Name, &description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if oid.Valid {
		m.OrgID = oid.UUID
	}
	if bid.Valid {
		m.BranchID = &bid.UUID
	}
	if description.Valid {
		m.Description = description.String
	}
	return &m, nil
}
