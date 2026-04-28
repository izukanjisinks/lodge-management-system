package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{db: database.DB}
}

// Create inserts a new order and all its items in a single transaction.
// Menu item prices are snapshotted at order time.
func (r *OrderRepository) Create(o *models.Order, items []models.PlaceOrderItemRequest, orgID uuid.UUID) (*models.Order, error) {
	o.ID = uuid.New()
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	o.OrgID = orgID

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.QueryRow(`
		INSERT INTO orders (id, org_id, booking_id, type, notes, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING order_number`,
		o.ID, orgID, o.BookingID, o.Type, o.Notes, now, now,
	).Scan(&o.OrderNumber)
	if err != nil {
		return nil, err
	}

	for _, req := range items {
		var price float64
		err = tx.QueryRow(`SELECT price FROM menu_items WHERE id=$1 AND org_id=$2`, req.MenuItemID, orgID).Scan(&price)
		if err != nil {
			return nil, fmt.Errorf("menu item %s not found", req.MenuItemID)
		}
		subtotal := price * float64(req.Quantity)
		_, err = tx.Exec(`
			INSERT INTO order_items (id, order_id, menu_item_id, quantity, unit_price, subtotal, notes, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			uuid.New(), o.ID, req.MenuItemID, req.Quantity, price, subtotal, req.Notes, now,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(o.ID, orgID)
}

func (r *OrderRepository) GetByID(id uuid.UUID, orgID uuid.UUID) (*models.Order, error) {
	var o models.Order
	var bookingID uuid.NullUUID
	var notes sql.NullString

	err := r.db.QueryRow(`
		SELECT id, org_id, booking_id, order_number, type, notes, created_at, updated_at
		FROM orders WHERE id=$1 AND org_id=$2`, id, orgID).
		Scan(&o.ID, &o.OrgID, &bookingID, &o.OrderNumber, &o.Type, &notes, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if bookingID.Valid {
		o.BookingID = &bookingID.UUID
	}
	if notes.Valid {
		o.Notes = notes.String
	}
	o.Items, err = r.fetchItems(id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) List(orgID uuid.UUID, orderType string, bookingID *uuid.UUID, page, pageSize int) ([]models.Order, int, error) {
	args := []interface{}{orgID}
	where := []string{"org_id = $1"}
	i := 2

	if orderType != "" {
		where = append(where, fmt.Sprintf("type = $%d", i))
		args = append(args, orderType)
		i++
	}
	if bookingID != nil {
		where = append(where, fmt.Sprintf("booking_id = $%d", i))
		args = append(args, *bookingID)
		i++
	}

	whereStr := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM orders WHERE %s`, whereStr), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, org_id, booking_id, order_number, type, notes, created_at, updated_at
		FROM orders WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var bid uuid.NullUUID
		var notes sql.NullString
		if err := rows.Scan(&o.ID, &o.OrgID, &bid, &o.OrderNumber, &o.Type, &notes, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if bid.Valid {
			o.BookingID = &bid.UUID
		}
		if notes.Valid {
			o.Notes = notes.String
		}
		orders = append(orders, o)
	}
	return orders, total, rows.Err()
}

// AddItems appends more items to an existing order, snapshotting prices at the time of addition.
func (r *OrderRepository) AddItems(orderID uuid.UUID, orgID uuid.UUID, items []models.PlaceOrderItemRequest) (*models.Order, error) {
	// Verify the order exists and belongs to this org
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1 AND org_id=$2)`, orderID, orgID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	for _, req := range items {
		var price float64
		err = tx.QueryRow(`SELECT price FROM menu_items WHERE id=$1 AND org_id=$2`, req.MenuItemID, orgID).Scan(&price)
		if err != nil {
			return nil, fmt.Errorf("menu item %s not found", req.MenuItemID)
		}
		subtotal := price * float64(req.Quantity)
		_, err = tx.Exec(`
			INSERT INTO order_items (id, order_id, menu_item_id, quantity, unit_price, subtotal, notes, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			uuid.New(), orderID, req.MenuItemID, req.Quantity, price, subtotal, req.Notes, now,
		)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(`UPDATE orders SET updated_at=$1 WHERE id=$2`, now, orderID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(orderID, orgID)
}

// GetItemsTotal returns the sum of all item subtotals for an order.
func (r *OrderRepository) GetItemsTotal(orderID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.QueryRow(`SELECT COALESCE(SUM(subtotal), 0) FROM order_items WHERE order_id=$1`, orderID).Scan(&total)
	return total, err
}

func (r *OrderRepository) fetchItems(orderID uuid.UUID) ([]models.OrderItem, error) {
	rows, err := r.db.Query(`
		SELECT oi.id, oi.order_id, oi.menu_item_id, mi.name, oi.quantity, oi.unit_price, oi.subtotal, oi.notes, oi.created_at
		FROM order_items oi
		JOIN menu_items mi ON mi.id = oi.menu_item_id
		WHERE oi.order_id=$1
		ORDER BY oi.created_at ASC`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var notes sql.NullString
		if err := rows.Scan(&item.ID, &item.OrderID, &item.MenuItemID, &item.ItemName, &item.Quantity, &item.UnitPrice, &item.Subtotal, &notes, &item.CreatedAt); err != nil {
			return nil, err
		}
		if notes.Valid {
			item.Notes = notes.String
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
