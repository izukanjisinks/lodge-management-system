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

func (r *BookingRepository) Create(b *models.Booking) error {
	b.ID = uuid.New()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

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
		INSERT INTO bookings
		    (id, user_id, room_id, client_id, client_type, check_in, check_out, guests, status, special_requests, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		b.ID, b.UserID, b.RoomID, b.ClientID, b.ClientType,
		b.CheckIn, b.CheckOut, b.Guests, b.Status, b.SpecialRequests,
		b.CreatedAt, b.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if b.MealPlanID != nil {
		_, err = tx.Exec(`
			INSERT INTO booking_meal_plans (booking_id, meal_plan_id, guests)
			VALUES ($1, $2, $3)`,
			b.ID, *b.MealPlanID, b.Guests,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *BookingRepository) GetByID(id uuid.UUID) (*models.Booking, error) {
	row := r.db.QueryRow(`
		SELECT b.id, b.user_id, b.room_id, r.name AS room_name,
		       b.client_id, b.client_type,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       bmp.meal_plan_id, mp.name AS meal_plan_name,
		       b.check_in, b.check_out, b.guests,
		       GREATEST(b.check_out - b.check_in, 1) AS nights,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS room_cost,
		       COALESCE(GREATEST(b.check_out - b.check_in, 1) * bmp.guests * mp.price_per_person_per_night, 0) AS meal_cost,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night +
		           COALESCE(GREATEST(b.check_out - b.check_in, 1) * bmp.guests * mp.price_per_person_per_night, 0) AS total_amount,
		       b.status, b.special_requests,
		       b.created_at, b.updated_at
		FROM bookings b
		JOIN rooms                    r   ON r.id = b.room_id
		LEFT JOIN individual_profiles ip  ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp  ON b.client_type = 'corporate'  AND cp.id = b.client_id
		LEFT JOIN booking_meal_plans  bmp ON bmp.booking_id = b.id
		LEFT JOIN meal_plans          mp  ON mp.id = bmp.meal_plan_id
		WHERE b.id = $1`, id)
	return scanBooking(row)
}

func (r *BookingRepository) List(status, clientType string, clientID *uuid.UUID, page, pageSize int) ([]models.Booking, int, error) {
	args := []interface{}{}
	where := []string{}
	i := 1

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
		SELECT b.id, b.user_id, b.room_id, r.name AS room_name,
		       b.client_id, b.client_type,
		       CASE b.client_type
		           WHEN 'individual' THEN ip.full_name
		           WHEN 'corporate'  THEN cp.company_name
		       END AS client_name,
		       bmp.meal_plan_id, mp.name AS meal_plan_name,
		       b.check_in, b.check_out, b.guests,
		       GREATEST(b.check_out - b.check_in, 1) AS nights,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night AS room_cost,
		       COALESCE(GREATEST(b.check_out - b.check_in, 1) * bmp.guests * mp.price_per_person_per_night, 0) AS meal_cost,
		       GREATEST(b.check_out - b.check_in, 1) * r.price_per_night +
		           COALESCE(GREATEST(b.check_out - b.check_in, 1) * bmp.guests * mp.price_per_person_per_night, 0) AS total_amount,
		       b.status, b.special_requests,
		       b.created_at, b.updated_at
		FROM bookings b
		JOIN rooms                    r   ON r.id = b.room_id
		LEFT JOIN individual_profiles ip  ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp  ON b.client_type = 'corporate'  AND cp.id = b.client_id
		LEFT JOIN booking_meal_plans  bmp ON bmp.booking_id = b.id
		LEFT JOIN meal_plans          mp  ON mp.id = bmp.meal_plan_id
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

func (r *BookingRepository) Update(b *models.Booking) error {
	b.UpdatedAt = time.Now()

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
		UPDATE bookings
		SET check_in=$1, check_out=$2, guests=$3, special_requests=$4, updated_at=$5
		WHERE id=$6`,
		b.CheckIn, b.CheckOut, b.Guests, b.SpecialRequests, b.UpdatedAt, b.ID,
	)
	if err != nil {
		return err
	}

	// Replace meal plan — delete existing then insert new if provided
	_, err = tx.Exec(`DELETE FROM booking_meal_plans WHERE booking_id=$1`, b.ID)
	if err != nil {
		return err
	}
	if b.MealPlanID != nil {
		_, err = tx.Exec(`
			INSERT INTO booking_meal_plans (booking_id, meal_plan_id, guests)
			VALUES ($1, $2, $3)`,
			b.ID, *b.MealPlanID, b.Guests,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateStatusTx updates the booking status and the room's availability atomically.
// confirmed  → room becomes unavailable
// checked_out / cancelled → room becomes available again
func (r *BookingRepository) UpdateStatusTx(id uuid.UUID, newStatus string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Fetch current booking to get room_id
	var roomID uuid.UUID
	var currentStatus string
	err = tx.QueryRow(`SELECT room_id, status FROM bookings WHERE id=$1`, id).Scan(&roomID, &currentStatus)
	if err != nil {
		return fmt.Errorf("booking not found")
	}

	// Update booking status
	_, err = tx.Exec(`UPDATE bookings SET status=$1, updated_at=$2 WHERE id=$3`,
		newStatus, time.Now(), id)
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

func (r *BookingRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM bookings WHERE id=$1`, id)
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

type bookingScanner interface {
	Scan(dest ...interface{}) error
}

func scanBooking(row bookingScanner) (*models.Booking, error) {
	var b models.Booking
	var roomName, clientName, mealPlanName, specialRequests sql.NullString
	var mealPlanID uuid.NullUUID
	err := row.Scan(
		&b.ID, &b.UserID, &b.RoomID, &roomName,
		&b.ClientID, &b.ClientType, &clientName,
		&mealPlanID, &mealPlanName,
		&b.CheckIn, &b.CheckOut, &b.Guests,
		&b.Nights, &b.RoomCost, &b.MealCost, &b.TotalAmount,
		&b.Status, &specialRequests,
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
	if mealPlanID.Valid {
		b.MealPlanID = &mealPlanID.UUID
	}
	if mealPlanName.Valid {
		b.MealPlanName = mealPlanName.String
	}
	if specialRequests.Valid {
		b.SpecialRequests = specialRequests.String
	}
	return &b, nil
}
