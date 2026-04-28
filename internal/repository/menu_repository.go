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

func (r *MenuRepository) CreateMenu(m *models.Menu, orgID uuid.UUID) error {
	m.ID = uuid.New()
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.OrgID = orgID

	_, err := r.db.Exec(`
		INSERT INTO menus (id, org_id, name, description, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		m.ID, orgID, m.Name, m.Description, m.IsActive, m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *MenuRepository) GetMenuByID(id uuid.UUID, orgID uuid.UUID) (*models.Menu, error) {
	var m models.Menu
	var description sql.NullString
	err := r.db.QueryRow(`
		SELECT id, org_id, name, description, is_active, created_at, updated_at
		FROM menus WHERE id=$1 AND org_id=$2`, id, orgID).
		Scan(&m.ID, &m.OrgID, &m.Name, &description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		m.Description = description.String
	}
	items, err := r.ListMenuItems(id, orgID)
	if err != nil {
		return nil, err
	}
	m.Items = items
	return &m, nil
}

func (r *MenuRepository) ListMenus(orgID uuid.UUID, page, pageSize int) ([]models.Menu, int, error) {
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM menus WHERE org_id=$1`, orgID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, org_id, name, description, is_active, created_at, updated_at
		FROM menus WHERE org_id=$1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3`, orgID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var m models.Menu
		var description sql.NullString
		if err := rows.Scan(&m.ID, &m.OrgID, &m.Name, &description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if description.Valid {
			m.Description = description.String
		}
		menus = append(menus, m)
	}
	return menus, total, rows.Err()
}

func (r *MenuRepository) UpdateMenu(id uuid.UUID, orgID uuid.UUID, req *models.UpdateMenuRequest) (*models.Menu, error) {
	m, err := r.GetMenuByID(id, orgID)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = *req.Description
	}
	if req.IsActive != nil {
		m.IsActive = *req.IsActive
	}
	m.UpdatedAt = time.Now()

	_, err = r.db.Exec(`
		UPDATE menus SET name=$1, description=$2, is_active=$3, updated_at=$4
		WHERE id=$5 AND org_id=$6`,
		m.Name, m.Description, m.IsActive, m.UpdatedAt, id, orgID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetMenuByID(id, orgID)
}

func (r *MenuRepository) DeleteMenu(id uuid.UUID, orgID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM menus WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("menu not found")
	}
	return nil
}

// ── Menu Items ────────────────────────────────────────────────────────────────

func (r *MenuRepository) CreateMenuItem(item *models.MenuItem, menuID uuid.UUID, orgID uuid.UUID) error {
	item.ID = uuid.New()
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	item.MenuID = menuID
	item.OrgID = orgID

	_, err := r.db.Exec(`
		INSERT INTO menu_items (id, menu_id, org_id, name, description, price, is_available, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		item.ID, menuID, orgID, item.Name, item.Description, item.Price, item.IsAvailable,
		item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *MenuRepository) GetMenuItemByID(id uuid.UUID, orgID uuid.UUID) (*models.MenuItem, error) {
	var item models.MenuItem
	var description sql.NullString
	err := r.db.QueryRow(`
		SELECT id, menu_id, org_id, name, description, price, is_available, created_at, updated_at
		FROM menu_items WHERE id=$1 AND org_id=$2`, id, orgID).
		Scan(&item.ID, &item.MenuID, &item.OrgID, &item.Name, &description, &item.Price, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		item.Description = description.String
	}
	return &item, nil
}

func (r *MenuRepository) ListMenuItems(menuID uuid.UUID, orgID uuid.UUID) ([]models.MenuItem, error) {
	rows, err := r.db.Query(`
		SELECT id, menu_id, org_id, name, description, price, is_available, created_at, updated_at
		FROM menu_items WHERE menu_id=$1 AND org_id=$2
		ORDER BY name ASC`, menuID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var description sql.NullString
		if err := rows.Scan(&item.ID, &item.MenuID, &item.OrgID, &item.Name, &description, &item.Price, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		if description.Valid {
			item.Description = description.String
		}
		items = append(items, item)
	}
	return items, rows.Err()
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
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}
	item.UpdatedAt = time.Now()

	_, err = r.db.Exec(`
		UPDATE menu_items SET name=$1, description=$2, price=$3, is_available=$4, updated_at=$5
		WHERE id=$6 AND org_id=$7`,
		item.Name, item.Description, item.Price, item.IsAvailable, item.UpdatedAt, id, orgID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetMenuItemByID(id, orgID)
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

// GuestListMenus returns active menus with their available items for public display.
// orgID nil means all orgs.
func (r *MenuRepository) GuestListMenus(orgID *uuid.UUID, page, pageSize int) ([]models.Menu, int, error) {
	var (
		args  []interface{}
		where = "m.is_active = TRUE"
	)
	if orgID != nil {
		where += " AND m.org_id = $1"
		args = append(args, *orgID)
	}

	countArgs := append([]interface{}{}, args...)
	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM menus m WHERE %s`, where), countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitIdx := len(args) + 1
	offsetIdx := limitIdx + 1
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, org_id, name, description, is_active, created_at, updated_at
		FROM menus m WHERE %s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, where, limitIdx, offsetIdx), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var m models.Menu
		var description sql.NullString
		if err := rows.Scan(&m.ID, &m.OrgID, &m.Name, &description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if description.Valid {
			m.Description = description.String
		}
		// Load only available items for guest view
		itemRows, err := r.db.Query(`
			SELECT id, menu_id, org_id, name, description, price, is_available, created_at, updated_at
			FROM menu_items WHERE menu_id=$1 AND is_available=TRUE
			ORDER BY name ASC`, m.ID)
		if err != nil {
			return nil, 0, err
		}
		for itemRows.Next() {
			var item models.MenuItem
			var desc sql.NullString
			if err := itemRows.Scan(&item.ID, &item.MenuID, &item.OrgID, &item.Name, &desc, &item.Price, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt); err != nil {
				itemRows.Close()
				return nil, 0, err
			}
			if desc.Valid {
				item.Description = desc.String
			}
			m.Items = append(m.Items, item)
		}
		itemRows.Close()
		menus = append(menus, m)
	}
	return menus, total, rows.Err()
}
