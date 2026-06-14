package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository() *DashboardRepository {
	return &DashboardRepository{db: database.DB}
}

func (r *DashboardRepository) StatCards(orgID uuid.UUID, branchID *uuid.UUID) (models.DashboardStatCards, error) {
	var s models.DashboardStatCards
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT
		    (SELECT COUNT(*) FROM bookings b
		        WHERE b.org_id = $2 AND b.status = 'pending'%[1]s) AS new_bookings_this_month,
		    (SELECT COUNT(*) FROM booking_room_assignments bra
		        JOIN bookings b ON b.id = bra.booking_id
		        WHERE b.org_id = $2 AND bra.check_in::date = $1::date
		          AND bra.status IN ('confirmed','checked_in')%[1]s) AS checkins_today,
		    (SELECT COUNT(*) FROM booking_room_assignments bra
		        JOIN bookings b ON b.id = bra.booking_id
		        WHERE b.org_id = $2 AND bra.check_out::date = $1::date
		          AND bra.status = 'checked_in'%[1]s) AS checkouts_today`
	args := []interface{}{today, orgID}
	branchFilter := ""
	if branchID != nil {
		args = append(args, *branchID)
		branchFilter = fmt.Sprintf(" AND b.branch_id = $%d", len(args))
	}
	query = fmt.Sprintf(query, branchFilter)
	err := r.db.QueryRow(query, args...).Scan(&s.NewBookingsThisMonth, &s.CheckInsToday, &s.CheckOutsToday)
	return s, err
}

func (r *DashboardRepository) RoomSummary(orgID uuid.UUID, branchID *uuid.UUID) (models.DashboardRoomSummary, error) {
	var s models.DashboardRoomSummary
	query := `
		SELECT
		    COUNT(*) FILTER (WHERE EXISTS (
		        SELECT 1 FROM booking_room_assignments bra
		        WHERE bra.room_id = r.id AND bra.status = 'checked_in'
		    )) AS occupied,
		    COUNT(*) FILTER (WHERE EXISTS (
		        SELECT 1 FROM booking_room_assignments bra
		        WHERE bra.room_id = r.id AND bra.status = 'confirmed'
		    )) AS reserved,
		    COUNT(*) FILTER (WHERE r.is_available = TRUE AND NOT EXISTS (
		        SELECT 1 FROM booking_room_assignments bra
		        WHERE bra.room_id = r.id AND bra.status IN ('confirmed','checked_in')
		    )) AS available,
		    COUNT(*) FILTER (WHERE r.is_available = FALSE AND NOT EXISTS (
		        SELECT 1 FROM booking_room_assignments bra
		        WHERE bra.room_id = r.id AND bra.status IN ('confirmed','checked_in')
		    )) AS not_ready
		FROM rooms r
		WHERE r.org_id = $1`
	args := []interface{}{orgID}
	if branchID != nil {
		args = append(args, *branchID)
		query += fmt.Sprintf(" AND r.branch_id = $%d", len(args))
	}
	err := r.db.QueryRow(query, args...).Scan(&s.Occupied, &s.Reserved, &s.Available, &s.NotReady)
	return s, err
}

func (r *DashboardRepository) RevenueByMonth(orgID uuid.UUID, branchID *uuid.UUID, months int) ([]models.DashboardRevenuePoint, error) {
	args := []interface{}{orgID, months}
	query := `
		SELECT
		    TO_CHAR(DATE_TRUNC('month', i.created_at), 'YYYY-MM') AS month,
		    COALESCE(SUM(i.total), 0) AS revenue
		FROM invoices i
		WHERE i.org_id = $1
		  AND i.status IN ('issued', 'paid')
		  AND i.created_at >= DATE_TRUNC('month', NOW()) - ($2 - 1) * INTERVAL '1 month'`
	if branchID != nil {
		args = append(args, *branchID)
		query += fmt.Sprintf(" AND i.branch_id = $%d", len(args))
	}
	query += ` GROUP BY DATE_TRUNC('month', i.created_at) ORDER BY DATE_TRUNC('month', i.created_at) ASC`
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []models.DashboardRevenuePoint
	for rows.Next() {
		var p models.DashboardRevenuePoint
		if err := rows.Scan(&p.Month, &p.Revenue); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	if points == nil {
		points = []models.DashboardRevenuePoint{}
	}
	return points, rows.Err()
}

func (r *DashboardRepository) ReservationsByDay(orgID uuid.UUID, branchID *uuid.UUID, days int) ([]models.DashboardReservationPoint, error) {
	args := []interface{}{days, orgID}
	branchFilter := ""
	if branchID != nil {
		args = append(args, *branchID)
		branchFilter = fmt.Sprintf(" AND b.branch_id = $%d", len(args))
	}
	query := fmt.Sprintf(`
		SELECT
		    TO_CHAR(day::date, 'YYYY-MM-DD') AS day,
		    COUNT(*) FILTER (WHERE b.status != 'cancelled') AS booked,
		    COUNT(*) FILTER (WHERE b.status = 'cancelled') AS cancelled
		FROM generate_series(
		    (CURRENT_DATE - ($1 - 1) * INTERVAL '1 day')::date,
		    CURRENT_DATE,
		    '1 day'::interval
		) AS day
		LEFT JOIN bookings b ON b.created_at::date = day::date AND b.org_id = $2%s
		GROUP BY day
		ORDER BY day ASC`, branchFilter)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []models.DashboardReservationPoint
	for rows.Next() {
		var p models.DashboardReservationPoint
		if err := rows.Scan(&p.Day, &p.Booked, &p.Cancelled); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	if points == nil {
		points = []models.DashboardReservationPoint{}
	}
	return points, rows.Err()
}

func (r *DashboardRepository) RecentBookings(orgID uuid.UUID, branchID *uuid.UUID, limit int) ([]models.DashboardRecentBooking, error) {
	args := []interface{}{orgID, limit}
	branchFilter := ""
	if branchID != nil {
		args = append(args, *branchID)
		branchFilter = fmt.Sprintf(" AND b.branch_id = $%d", len(args))
	}
	query := fmt.Sprintf(`
		SELECT
		    b.id,
		    b.booking_number,
		    b.booker_name AS client_name,
		    COALESCE(r.name, '') AS room_name,
		    COALESCE(r.type::text, '')  AS room_type,
		    COALESCE(TO_CHAR(asg.check_in,  'Mon DD, YYYY'), '') AS check_in,
		    COALESCE(TO_CHAR(asg.check_out, 'Mon DD, YYYY'), '') AS check_out,
		    b.status
		FROM bookings b
		LEFT JOIN LATERAL (
		    SELECT bra.room_id, bra.check_in, bra.check_out
		    FROM booking_room_assignments bra
		    WHERE bra.booking_id = b.id
		    ORDER BY bra.check_in ASC
		    LIMIT 1
		) asg ON TRUE
		LEFT JOIN rooms r ON r.id = asg.room_id
		WHERE b.org_id = $1%s
		ORDER BY b.created_at DESC
		LIMIT $2`, branchFilter)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.DashboardRecentBooking
	for rows.Next() {
		var b models.DashboardRecentBooking
		var clientName sql.NullString
		if err := rows.Scan(&b.ID, &b.BookingNumber, &clientName, &b.RoomName, &b.RoomType, &b.CheckIn, &b.CheckOut, &b.Status); err != nil {
			return nil, err
		}
		if clientName.Valid {
			b.ClientName = clientName.String
		}
		bookings = append(bookings, b)
	}
	if bookings == nil {
		bookings = []models.DashboardRecentBooking{}
	}
	return bookings, rows.Err()
}
