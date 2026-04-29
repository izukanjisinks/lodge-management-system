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

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{db: database.DB}
}

func (r *BookingRepository) Create(b *models.Booking, orgID uuid.UUID) error {
	b.ID = uuid.New()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

	_, err := r.db.Exec(`
		INSERT INTO bookings
		    (id, user_id, room_id, client_id, client_type, check_in, check_out, guests, status, special_requests, org_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		b.ID, b.UserID, b.RoomID, b.ClientID, b.ClientType,
		b.CheckIn, b.CheckOut, b.Guests, b.Status, b.SpecialRequests,
		orgID, b.CreatedAt, b.UpdatedAt,
	)
	return err
}

// GetByIDUnscoped fetches a booking by ID with no org filter — use only in guest/review flows
// where the caller has already verified ownership via client_id.
func (r *BookingRepository) GetByIDUnscoped(id uuid.UUID) (*models.Booking, error) {
	row := r.db.QueryRow(`
		SELECT b.id, b.booking_number, b.user_id, b.room_id, r.name AS room_name,
		       b.client_id, b.client_type,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       b.check_in, b.check_out, b.guests,
		       GREATEST(b.check_out - b.check_in, 1) AS nights,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS room_cost,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS total_amount,
		       b.status, b.overstayed, b.special_requests,
		       b.created_at, b.updated_at
		FROM bookings b
		JOIN rooms                    r   ON r.id = b.room_id
		LEFT JOIN individual_profiles ip  ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp  ON b.client_type = 'corporate'  AND cp.id = b.client_id
		WHERE b.id = $1`, id)
	return scanBooking(row)
}

func (r *BookingRepository) GetByID(id uuid.UUID, orgID uuid.UUID) (*models.Booking, error) {
	row := r.db.QueryRow(`
		SELECT b.id, b.booking_number, b.user_id, b.room_id, r.name AS room_name,
		       b.client_id, b.client_type,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       b.check_in, b.check_out, b.guests,
		       GREATEST(b.check_out - b.check_in, 1) AS nights,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS room_cost,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS total_amount,
		       b.status, b.overstayed, b.special_requests,
		       b.created_at, b.updated_at
		FROM bookings b
		JOIN rooms                    r   ON r.id = b.room_id
		LEFT JOIN individual_profiles ip  ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp  ON b.client_type = 'corporate'  AND cp.id = b.client_id
		WHERE b.id = $1 AND b.org_id = $2`, id, orgID)
	return scanBooking(row)
}

func (r *BookingRepository) List(orgID uuid.UUID, status, clientType string, clientID *uuid.UUID, page, pageSize int) ([]models.Booking, int, error) {
	args := []interface{}{}
	where := []string{}
	i := 1

	if orgID != uuid.Nil {
		where = append(where, fmt.Sprintf("b.org_id = $%d", i))
		args = append(args, orgID)
		i++
	}
	if status != "" {
		where = append(where, fmt.Sprintf("b.status = $%d", i))
		args = append(args, status)
		i++
	}
	if clientType != "" {
		where = append(where, fmt.Sprintf("b.client_type = $%d", i))
		args = append(args, clientType)
		i++
	}
	if clientID != nil {
		where = append(where, fmt.Sprintf("b.client_id = $%d", i))
		args = append(args, *clientID)
		i++
	}

	whereStr := "b.id IS NOT NULL"
	if len(where) > 0 {
		whereStr = strings.Join(where, " AND ")
	}

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM bookings b WHERE %s`, whereStr), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT b.id, b.booking_number, b.user_id, b.room_id, r.name AS room_name,
		       b.client_id, b.client_type,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       b.check_in, b.check_out, b.guests,
		       GREATEST(b.check_out - b.check_in, 1) AS nights,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS room_cost,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS total_amount,
		       b.status, b.overstayed, b.special_requests,
		       b.created_at, b.updated_at
		FROM bookings b
		JOIN rooms                    r   ON r.id = b.room_id
		LEFT JOIN individual_profiles ip  ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp  ON b.client_type = 'corporate'  AND cp.id = b.client_id
		WHERE %s
		ORDER BY b.created_at DESC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		b, err := scanBooking(rows)
		if err != nil {
			return nil, 0, err
		}
		bookings = append(bookings, *b)
	}
	return bookings, total, rows.Err()
}

func (r *BookingRepository) Update(b *models.Booking, orgID uuid.UUID) error {
	b.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE bookings
		SET check_in=$1, check_out=$2, guests=$3, special_requests=$4, updated_at=$5
		WHERE id=$6 AND org_id=$7`,
		b.CheckIn, b.CheckOut, b.Guests, b.SpecialRequests, b.UpdatedAt, b.ID, orgID,
	)
	return err
}

// UpdateStatusTx updates the booking status and the room's availability atomically.
// confirmed  → room becomes unavailable
// checked_out / cancelled → room becomes available again
func (r *BookingRepository) UpdateStatusTx(id uuid.UUID, orgID uuid.UUID, newStatus string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Fetch current booking to get room_id, scoped to org
	var roomID uuid.UUID
	var currentStatus string
	err = tx.QueryRow(`SELECT room_id, status FROM bookings WHERE id=$1 AND org_id=$2`, id, orgID).Scan(&roomID, &currentStatus)
	if err != nil {
		return fmt.Errorf("booking not found")
	}

	// Update booking status
	_, err = tx.Exec(`UPDATE bookings SET status=$1, updated_at=$2 WHERE id=$3 AND org_id=$4`,
		newStatus, time.Now(), id, orgID)
	if err != nil {
		return err
	}

	// Sync room availability
	switch newStatus {
	case models.BookingStatusConfirmed:
		_, err = tx.Exec(`UPDATE rooms SET is_available=FALSE, updated_at=$1 WHERE id=$2`, time.Now(), roomID)
	case models.BookingStatusCheckedOut, models.BookingStatusCancelled:
		_, err = tx.Exec(`UPDATE rooms SET is_available=TRUE, updated_at=$1 WHERE id=$2`, time.Now(), roomID)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *BookingRepository) Delete(id uuid.UUID, orgID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM bookings WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("booking not found")
	}
	return nil
}

// IsRoomAvailable returns true if the room has no active bookings (pending/confirmed/checked_in)
// overlapping [checkIn, checkOut), optionally excluding a booking by ID (for updates).
func (r *BookingRepository) IsRoomAvailable(roomID uuid.UUID, checkIn, checkOut time.Time, excludeID *uuid.UUID) (bool, error) {
	args := []interface{}{roomID, checkOut, checkIn}
	excludeClause := ""
	if excludeID != nil {
		args = append(args, *excludeID)
		excludeClause = fmt.Sprintf(" AND id != $%d", len(args))
	}

	var count int
	err := r.db.QueryRow(fmt.Sprintf(`
		SELECT COUNT(*) FROM bookings
		WHERE room_id = $1
		  AND status IN ('pending', 'confirmed', 'checked_in')
		  AND check_in  < $2
		  AND check_out > $3%s`, excludeClause), args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// HasActiveBookingForClient returns true if the client already has a pending/confirmed/checked_in
// booking on the same room overlapping [checkIn, checkOut).
func (r *BookingRepository) HasActiveBookingForClient(clientID, roomID uuid.UUID, checkIn, checkOut time.Time) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM bookings
		WHERE client_id = $1
		  AND room_id   = $2
		  AND status IN ('pending', 'confirmed', 'checked_in')
		  AND check_in  < $3
		  AND check_out > $4`,
		clientID, roomID, checkOut, checkIn).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// OverdueBookingRef is a lightweight projection used by the nightly overdue-checkout job.
type OverdueBookingRef struct {
	ID    uuid.UUID
	OrgID uuid.UUID
}

// MarkOverstayed sets overstayed=TRUE on a booking. Called only by the nightly job.
func (r *BookingRepository) MarkOverstayed(id uuid.UUID, orgID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE bookings SET overstayed=TRUE, updated_at=$1 WHERE id=$2 AND org_id=$3`,
		time.Now(), id, orgID,
	)
	return err
}

// ClearOverstayed sets overstayed=FALSE on a booking. Called when staff manually resolves the flag.
func (r *BookingRepository) ClearOverstayed(id uuid.UUID, orgID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE bookings SET overstayed=FALSE, updated_at=$1 WHERE id=$2 AND org_id=$3`,
		time.Now(), id, orgID,
	)
	return err
}

// FindOverdueCheckouts returns all checked_in bookings whose check_out date is before today.
func (r *BookingRepository) FindOverdueCheckouts() ([]OverdueBookingRef, error) {
	rows, err := r.db.Query(`
		SELECT id, org_id FROM bookings
		WHERE status = 'checked_in'
		  AND check_out < CURRENT_DATE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []OverdueBookingRef
	for rows.Next() {
		var ref OverdueBookingRef
		if err := rows.Scan(&ref.ID, &ref.OrgID); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, rows.Err()
}

// ExtendCheckout updates a booking's check_out to newDate.
// Used by the nightly job to roll forward overdue guests to today.
func (r *BookingRepository) ExtendCheckout(id uuid.UUID, orgID uuid.UUID, newDate time.Time) error {
	_, err := r.db.Exec(`
		UPDATE bookings SET check_out=$1, updated_at=$2 WHERE id=$3 AND org_id=$4`,
		newDate, time.Now(), id, orgID,
	)
	return err
}

type bookingScanner interface {
	Scan(dest ...interface{}) error
}

func scanBooking(row bookingScanner) (*models.Booking, error) {
	var b models.Booking
	var roomName, clientName, specialRequests sql.NullString
	err := row.Scan(
		&b.ID, &b.BookingNumber, &b.UserID, &b.RoomID, &roomName,
		&b.ClientID, &b.ClientType, &clientName,
		&b.CheckIn, &b.CheckOut, &b.Guests,
		&b.Nights, &b.RoomCost, &b.TotalAmount,
		&b.Status, &b.Overstayed, &specialRequests,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if roomName.Valid {
		b.RoomName = roomName.String
	}
	if clientName.Valid {
		b.ClientName = clientName.String
	}
	if specialRequests.Valid {
		b.SpecialRequests = specialRequests.String
	}
	return &b, nil
}
