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
	var notes, bookingNumber, roomName, clientName sql.NullString

	err := r.db.QueryRow(`
		SELECT o.id, o.org_id, o.booking_id, o.order_number, o.type, o.notes,
		       COALESCE((SELECT SUM(subtotal) FROM order_items WHERE order_id = o.id), 0) AS total,
		       b.booking_number,
		       r.name AS room_name,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       o.created_at, o.updated_at
		FROM orders o
		LEFT JOIN bookings            b  ON b.id = o.booking_id
		LEFT JOIN rooms               r  ON r.id = b.room_id
		LEFT JOIN individual_profiles ip ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp ON b.client_type = 'corporate'  AND cp.id = b.client_id
		WHERE o.id=$1 AND o.org_id=$2`, id, orgID).
		Scan(&o.ID, &o.OrgID, &bookingID, &o.OrderNumber, &o.Type, &notes, &o.Total,
			&bookingNumber, &roomName, &clientName, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if bookingID.Valid {
		o.BookingID = &bookingID.UUID
	}
	if notes.Valid {
		o.Notes = notes.String
	}
	if bookingNumber.Valid {
		o.BookingNumber = bookingNumber.String
	}
	if roomName.Valid {
		o.RoomName = roomName.String
	}
	if clientName.Valid {
		o.ClientName = clientName.String
	}
	o.Items, err = r.fetchItems(id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) List(orgID uuid.UUID, orderType string, bookingID *uuid.UUID, page, pageSize int) ([]models.Order, int, error) {
	args := []interface{}{orgID}
	extraFilters := []string{}
	i := 2

	if orderType != "" {
		extraFilters = append(extraFilters, fmt.Sprintf("o.type = $%d", i))
		args = append(args, orderType)
		i++
	}
	if bookingID != nil {
		extraFilters = append(extraFilters, fmt.Sprintf("o.booking_id = $%d", i))
		args = append(args, *bookingID)
		i++
	}

	extraWhere := ""
	if len(extraFilters) > 0 {
		extraWhere = " AND " + strings.Join(extraFilters, " AND ")
	}

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM orders o WHERE o.org_id = $1%s`, extraWhere), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT o.id, o.org_id, o.booking_id, o.order_number, o.type, o.notes,
		       COALESCE((SELECT SUM(subtotal) FROM order_items WHERE order_id = o.id), 0) AS total,
		       b.booking_number,
		       r.name AS room_name,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       o.created_at, o.updated_at
		FROM orders o
		LEFT JOIN bookings            b  ON b.id = o.booking_id
		LEFT JOIN rooms               r  ON r.id = b.room_id
		LEFT JOIN individual_profiles ip ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp ON b.client_type = 'corporate'  AND cp.id = b.client_id
		WHERE o.org_id = $1%s
		ORDER BY o.created_at DESC
		LIMIT $%d OFFSET $%d`, extraWhere, i, i+1), append(args, pageSize, (page-1)*pageSize)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var bid uuid.NullUUID
		var notes, bookingNumber, roomName, clientName sql.NullString
		if err := rows.Scan(&o.ID, &o.OrgID, &bid, &o.OrderNumber, &o.Type, &notes, &o.Total,
			&bookingNumber, &roomName, &clientName, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if bid.Valid {
			o.BookingID = &bid.UUID
		}
		if notes.Valid {
			o.Notes = notes.String
		}
		if bookingNumber.Valid {
			o.BookingNumber = bookingNumber.String
		}
		if roomName.Valid {
			o.RoomName = roomName.String
		}
		if clientName.Valid {
			o.ClientName = clientName.String
		}
		orders = append(orders, o)
	}
	return orders, total, rows.Err()
}

// AddItems appends more items to an existing order, snapshotting prices at the time of addition.
// Returns the updated order and the newly inserted items so the caller can append them to the invoice.
func (r *OrderRepository) AddItems(orderID uuid.UUID, orgID uuid.UUID, items []models.PlaceOrderItemRequest) (*models.Order, []models.OrderItem, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1 AND org_id=$2)`, orderID, orgID).Scan(&exists)
	if err != nil {
		return nil, nil, err
	}
	if !exists {
		return nil, nil, fmt.Errorf("order not found")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	var newItems []models.OrderItem
	for _, req := range items {
		var price float64
		var itemName string
		err = tx.QueryRow(`SELECT price, name FROM menu_items WHERE id=$1 AND org_id=$2`, req.MenuItemID, orgID).Scan(&price, &itemName)
		if err != nil {
			return nil, nil, fmt.Errorf("menu item %s not found", req.MenuItemID)
		}
		subtotal := price * float64(req.Quantity)
		itemID := uuid.New()
		_, err = tx.Exec(`
			INSERT INTO order_items (id, order_id, menu_item_id, quantity, unit_price, subtotal, notes, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			itemID, orderID, req.MenuItemID, req.Quantity, price, subtotal, req.Notes, now,
		)
		if err != nil {
			return nil, nil, err
		}
		newItems = append(newItems, models.OrderItem{
			ID:         itemID,
			OrderID:    orderID,
			MenuItemID: req.MenuItemID,
			ItemName:   itemName,
			Quantity:   req.Quantity,
			UnitPrice:  price,
			Subtotal:   subtotal,
			Notes:      req.Notes,
			CreatedAt:  now,
		})
	}

	_, err = tx.Exec(`UPDATE orders SET updated_at=$1 WHERE id=$2`, now, orderID)
	if err != nil {
		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}
	order, err := r.GetByID(orderID, orgID)
	return order, newItems, err
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
