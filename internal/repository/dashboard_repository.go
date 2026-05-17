package repository

import (
	"database/sql"
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

func (r *DashboardRepository) StatCards(orgID uuid.UUID) (models.DashboardStatCards, error) {
	var s models.DashboardStatCards
	today := time.Now().Format("2006-01-02")

	err := r.db.QueryRow(`
		SELECT
		    COUNT(*) FILTER (WHERE status = 'pending') AS new_bookings_this_month,
		    COUNT(*) FILTER (WHERE check_in::date = $1::date AND status IN ('confirmed','checked_in')) AS checkins_today,
		    COUNT(*) FILTER (WHERE check_out::date = $1::date AND status = 'checked_in') AS checkouts_today
		FROM bookings
		WHERE org_id = $2`,
		today, orgID,
	).Scan(&s.NewBookingsThisMonth, &s.CheckInsToday, &s.CheckOutsToday)
	return s, err
}

func (r *DashboardRepository) RoomSummary(orgID uuid.UUID) (models.DashboardRoomSummary, error) {
	var s models.DashboardRoomSummary
	err := r.db.QueryRow(`
		SELECT
		    COUNT(*) FILTER (WHERE EXISTS (
		        SELECT 1 FROM bookings b
		        WHERE b.room_id = r.id AND b.status = 'checked_in'
		    )) AS occupied,
		    COUNT(*) FILTER (WHERE EXISTS (
		        SELECT 1 FROM bookings b
		        WHERE b.room_id = r.id AND b.status = 'confirmed'
		    )) AS reserved,
		    COUNT(*) FILTER (WHERE r.is_available = TRUE AND NOT EXISTS (
		        SELECT 1 FROM bookings b
		        WHERE b.room_id = r.id AND b.status IN ('confirmed','checked_in')
		    )) AS available,
		    COUNT(*) FILTER (WHERE r.is_available = FALSE AND NOT EXISTS (
		        SELECT 1 FROM bookings b
		        WHERE b.room_id = r.id AND b.status IN ('confirmed','checked_in')
		    )) AS not_ready
		FROM rooms r
		WHERE r.org_id = $1`,
		orgID,
	).Scan(&s.Occupied, &s.Reserved, &s.Available, &s.NotReady)
	return s, err
}

func (r *DashboardRepository) RevenueByMonth(orgID uuid.UUID, months int) ([]models.DashboardRevenuePoint, error) {
	rows, err := r.db.Query(`
		SELECT
		    TO_CHAR(DATE_TRUNC('month', i.created_at), 'YYYY-MM') AS month,
		    COALESCE(SUM(i.total), 0) AS revenue
		FROM invoices i
		WHERE i.org_id = $1
		  AND i.status IN ('issued', 'paid')
		  AND i.created_at >= DATE_TRUNC('month', NOW()) - ($2 - 1) * INTERVAL '1 month'
		GROUP BY DATE_TRUNC('month', i.created_at)
		ORDER BY DATE_TRUNC('month', i.created_at) ASC`,
		orgID, months,
	)
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

func (r *DashboardRepository) ReservationsByDay(orgID uuid.UUID, days int) ([]models.DashboardReservationPoint, error) {
	rows, err := r.db.Query(`
		SELECT
		    TO_CHAR(day::date, 'YYYY-MM-DD') AS day,
		    COUNT(*) FILTER (WHERE b.status != 'cancelled') AS booked,
		    COUNT(*) FILTER (WHERE b.status = 'cancelled') AS cancelled
		FROM generate_series(
		    (CURRENT_DATE - ($1 - 1) * INTERVAL '1 day')::date,
		    CURRENT_DATE,
		    '1 day'::interval
		) AS day
		LEFT JOIN bookings b ON b.created_at::date = day::date AND b.org_id = $2
		GROUP BY day
		ORDER BY day ASC`,
		days, orgID,
	)
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

func (r *DashboardRepository) RecentBookings(orgID uuid.UUID, limit int) ([]models.DashboardRecentBooking, error) {
	rows, err := r.db.Query(`
		SELECT
		    b.id,
		    b.booking_number,
		    CASE b.client_type
		        WHEN 'individual' THEN ip.full_name
		        WHEN 'corporate'  THEN cp.company_name
		    END AS client_name,
		    r.name  AS room_name,
		    r.type  AS room_type,
		    TO_CHAR(b.check_in,  'Mon DD, YYYY') AS check_in,
		    TO_CHAR(b.check_out, 'Mon DD, YYYY') AS check_out,
		    b.status
		FROM bookings b
		JOIN rooms r                  ON r.id = b.room_id
		LEFT JOIN individual_profiles ip ON b.client_type = 'individual' AND ip.id = b.client_id
		LEFT JOIN corporate_profiles  cp ON b.client_type = 'corporate'  AND cp.id = b.client_id
		WHERE b.org_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2`,
		orgID, limit,
	)
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
