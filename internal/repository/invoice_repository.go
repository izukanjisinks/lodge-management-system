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

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{db: database.DB}
}

// Create inserts the invoice and all its line items in a single transaction.
func (r *InvoiceRepository) Create(inv *models.Invoice) error {
	inv.ID = uuid.New()
	now := time.Now()
	inv.CreatedAt = now
	inv.UpdatedAt = now

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO invoices
		    (id, invoice_number, booking_id, subtotal, tax_rate, tax, total, status, issued_at, due_date, notes, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		inv.ID, inv.InvoiceNumber, inv.BookingID,
		inv.Subtotal, inv.TaxRate, inv.TaxAmount, inv.Total,
		inv.Status, inv.IssuedDate, inv.DueDate, inv.Notes,
		inv.CreatedAt, inv.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for i := range inv.LineItems {
		inv.LineItems[i].ID = uuid.New()
		inv.LineItems[i].InvoiceID = inv.ID
		inv.LineItems[i].CreatedAt = now
		_, err = tx.Exec(`
			INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, total, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			inv.LineItems[i].ID, inv.ID,
			inv.LineItems[i].Description, inv.LineItems[i].Quantity,
			inv.LineItems[i].UnitPrice, inv.LineItems[i].Total,
			inv.LineItems[i].CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *InvoiceRepository) GetByID(id uuid.UUID) (*models.Invoice, error) {
	return r.fetchOne(`WHERE i.id = $1`, id)
}

func (r *InvoiceRepository) GetByBookingID(bookingID uuid.UUID) (*models.Invoice, error) {
	return r.fetchOne(`WHERE i.booking_id = $1`, bookingID)
}

func (r *InvoiceRepository) List(status string, page, pageSize int) ([]models.Invoice, int, error) {
	where := []string{}
	args := []interface{}{}
	i := 1

	if status != "" {
		where = append(where, fmt.Sprintf("i.status = $%d", i))
		args = append(args, status)
		i++
	}

	whereStr := "1=1"
	if len(where) > 0 {
		whereStr = strings.Join(where, " AND ")
	}

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM invoices i WHERE %s`, whereStr), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.booking_id,
		       b.client_id, b.client_type,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.email
		           WHEN 'corporate'  THEN cp.email
		       END AS client_email,
		       i.subtotal, i.tax_rate, i.tax, i.total,
		       i.status, i.issued_at, i.due_date, i.paid_date, i.notes,
		       i.created_at, i.updated_at
		FROM invoices i
		JOIN bookings b             ON b.id = i.booking_id
		LEFT JOIN individual_profiles ip ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp ON b.client_type = 'corporate'  AND cp.id = b.client_id
		WHERE %s
		ORDER BY i.created_at DESC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		lineItems, err := r.fetchLineItems(inv.ID)
		if err != nil {
			return nil, 0, err
		}
		inv.LineItems = lineItems
		invoices = append(invoices, *inv)
	}
	return invoices, total, rows.Err()
}

func (r *InvoiceRepository) UpdateStatus(id uuid.UUID, status string, paidDate *time.Time, notes *string) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE invoices
		SET status=$1, paid_date=$2, notes=COALESCE($3, notes), updated_at=$4
		WHERE id=$5`,
		status, paidDate, notes, now, id,
	)
	return err
}

// fetchOne is a shared helper used by GetByID and GetByBookingID.
func (r *InvoiceRepository) fetchOne(whereClause string, arg interface{}) (*models.Invoice, error) {
	row := r.db.QueryRow(fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.booking_id,
		       b.client_id, b.client_type,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.email
		           WHEN 'corporate'  THEN cp.email
		       END AS client_email,
		       i.subtotal, i.tax_rate, i.tax, i.total,
		       i.status, i.issued_at, i.due_date, i.paid_date, i.notes,
		       i.created_at, i.updated_at
		FROM invoices i
		JOIN bookings b             ON b.id = i.booking_id
		LEFT JOIN individual_profiles ip ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp ON b.client_type = 'corporate'  AND cp.id = b.client_id
		%s`, whereClause), arg)

	inv, err := scanInvoice(row)
	if err != nil {
		return nil, err
	}

	inv.LineItems, err = r.fetchLineItems(inv.ID)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *InvoiceRepository) fetchLineItems(invoiceID uuid.UUID) ([]models.InvoiceLineItem, error) {
	rows, err := r.db.Query(`
		SELECT id, invoice_id, description, quantity, unit_price, total, created_at
		FROM invoice_line_items WHERE invoice_id = $1
		ORDER BY created_at ASC`, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InvoiceLineItem
	for rows.Next() {
		var item models.InvoiceLineItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.Quantity, &item.UnitPrice, &item.Total, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []models.InvoiceLineItem{}
	}
	return items, rows.Err()
}

// GenerateInvoiceNumber produces a sequential human-readable number like INV-2026-0001.
func (r *InvoiceRepository) GenerateInvoiceNumber() (string, error) {
	var count int
	year := time.Now().Year()
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM invoices
		WHERE EXTRACT(YEAR FROM created_at) = $1`, year).Scan(&count)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("INV-%d-%04d", year, count+1), nil
}

type invoiceScanner interface {
	Scan(dest ...interface{}) error
}

func scanInvoice(row invoiceScanner) (*models.Invoice, error) {
	var inv models.Invoice
	var clientEmail, notes sql.NullString
	var issuedAt, dueDate, paidDate sql.NullTime
	err := row.Scan(
		&inv.ID, &inv.InvoiceNumber, &inv.BookingID,
		&inv.ClientID, &inv.ClientType, &inv.ClientName, &clientEmail,
		&inv.Subtotal, &inv.TaxRate, &inv.TaxAmount, &inv.Total,
		&inv.Status, &issuedAt, &dueDate, &paidDate, &notes,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if clientEmail.Valid {
		inv.ClientEmail = clientEmail.String
	}
	if notes.Valid {
		inv.Notes = notes.String
	}
	if issuedAt.Valid {
		inv.IssuedDate = &issuedAt.Time
	}
	if dueDate.Valid {
		inv.DueDate = &dueDate.Time
	}
	if paidDate.Valid {
		inv.PaidDate = &paidDate.Time
	}
	return &inv, nil
}
