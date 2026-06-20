package repository

import (
	"database/sql"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type BookingEventRepository struct {
	db *sql.DB
}

func NewBookingEventRepository() *BookingEventRepository {
	return &BookingEventRepository{db: database.DB}
}

func (r *BookingEventRepository) CreateInTx(tx *sql.Tx, e *models.BookingEvent) error {
	e.ID = uuid.New()
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now

	// start_time/end_time are TIME columns; pass NULL when blank so Postgres doesn't
	// choke on an empty string.
	startTime := nullableTime(e.StartTime)
	endTime := nullableTime(e.EndTime)

	return tx.QueryRow(`
		INSERT INTO booking_events (
			id, booking_id, venue_id, event_type,
			start_date, end_date, start_time, end_time,
			pax_count, price, catering_required, notes,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id`,
		e.ID, e.BookingID, e.VenueID, e.EventType,
		e.StartDate, e.EndDate, startTime, endTime,
		e.PaxCount, e.Price, e.CateringRequired, e.Notes,
		e.CreatedAt, e.UpdatedAt,
	).Scan(&e.ID)
}

// nullableTime returns a sql.NullString for a "HH:MM" time, NULL when blank.
func nullableTime(t string) sql.NullString {
	if t == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: t, Valid: true}
}

func (r *BookingEventRepository) ListByBookingID(bookingID uuid.UUID) ([]models.BookingEvent, error) {
	rows, err := r.db.Query(`
		SELECT e.id, e.booking_id, e.venue_id, e.event_type,
		       e.start_date, e.end_date,
		       TO_CHAR(e.start_time, 'HH24:MI'), TO_CHAR(e.end_time, 'HH24:MI'),
		       e.pax_count, e.price, e.catering_required, COALESCE(e.notes, ''),
		       e.created_at, e.updated_at,
		       COALESCE(v.name, '') AS venue_name
		FROM booking_events e
		LEFT JOIN venues v ON v.id = e.venue_id
		WHERE e.booking_id = $1
		ORDER BY e.start_date ASC`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.BookingEvent
	for rows.Next() {
		var e models.BookingEvent
		var venueID uuid.NullUUID
		var startTime, endTime sql.NullString
		if err := rows.Scan(
			&e.ID, &e.BookingID, &venueID, &e.EventType,
			&e.StartDate, &e.EndDate, &startTime, &endTime,
			&e.PaxCount, &e.Price, &e.CateringRequired, &e.Notes,
			&e.CreatedAt, &e.UpdatedAt, &e.VenueName,
		); err != nil {
			return nil, err
		}
		if venueID.Valid {
			e.VenueID = &venueID.UUID
		}
		if startTime.Valid {
			e.StartTime = startTime.String
		}
		if endTime.Valid {
			e.EndTime = endTime.String
		}
		days := int(e.EndDate.Sub(e.StartDate).Hours()/24) + 1
		if days < 1 {
			days = 1
		}
		e.Days = days
		events = append(events, e)
	}
	return events, rows.Err()
}

// SumForBooking returns the total event hire charge for a booking (price × days, summed).
func (r *BookingEventRepository) SumForBooking(bookingID uuid.UUID) (float64, error) {
	row := r.db.QueryRow(`
		SELECT COALESCE(SUM(price * (GREATEST((end_date - start_date), 0) + 1)), 0)
		FROM booking_events
		WHERE booking_id = $1`, bookingID)
	var total float64
	err := row.Scan(&total)
	return total, err
}
