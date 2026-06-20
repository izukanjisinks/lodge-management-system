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

func (r *BookingRepository) Begin() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *BookingRepository) Create(tx *sql.Tx, b *models.Booking) error {
	b.ID = uuid.New()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

	return tx.QueryRow(`
		INSERT INTO bookings (
			id, org_id, branch_id,
			booking_type, booker_type,
			booker_name, booker_email, booker_phone,
			web_user_id, cor_profile_id, company_id, request_id, venue_id,
			total_amount, status, special_requests, overstayed,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19
		) RETURNING booking_number`,
		b.ID, b.OrgID, b.BranchID,
		b.BookingType, b.BookerType,
		b.BookerName, b.BookerEmail, b.BookerPhone,
		b.WebUserID, b.CorProfileID, b.CompanyID, b.RequestID, b.VenueID,
		b.TotalAmount, b.Status, b.SpecialRequests, b.Overstayed,
		now, now,
	).Scan(&b.BookingNumber)
}

func (r *BookingRepository) GetByID(id, orgID uuid.UUID) (*models.Booking, error) {
	row := r.db.QueryRow(`
		SELECT b.id, b.booking_number, b.org_id, b.branch_id,
		       b.booking_type, b.booker_type,
		       b.booker_name, b.booker_email, b.booker_phone,
		       b.web_user_id, b.cor_profile_id, b.company_id, b.request_id, b.venue_id,
		       COALESCE(cd.company_name, '')      AS company_name,
		       COALESCE(cp.first_name || ' ' || cp.last_name, '') AS profile_name,
		       COALESCE(v.name, '')               AS venue_name,
		       b.total_amount, b.status, b.special_requests, b.overstayed,
		       b.created_at, b.updated_at
		FROM bookings b
		LEFT JOIN cor_company_details cd ON cd.id = b.company_id
		LEFT JOIN cor_profiles        cp ON cp.id = b.cor_profile_id
		LEFT JOIN venues              v  ON v.id  = b.venue_id
		WHERE b.id = $1 AND b.org_id = $2`, id, orgID)
	return scanBooking(row)
}

func (r *BookingRepository) List(orgID uuid.UUID, bookerType, bookingType, status string, page, pageSize int) ([]models.Booking, int, error) {
	args := []interface{}{orgID}
	where := []string{"b.org_id = $1"}
	i := 2

	if bookerType != "" {
		where = append(where, fmt.Sprintf("b.booker_type = $%d", i))
		args = append(args, bookerType)
		i++
	}
	if bookingType != "" {
		where = append(where, fmt.Sprintf("b.booking_type = $%d", i))
		args = append(args, bookingType)
		i++
	}
	if status != "" {
		where = append(where, fmt.Sprintf("b.status = $%d", i))
		args = append(args, status)
		i++
	}

	whereStr := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM bookings b WHERE %s`, whereStr), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT b.id, b.booking_number, b.org_id, b.branch_id,
		       b.booking_type, b.booker_type,
		       b.booker_name, b.booker_email, b.booker_phone,
		       b.web_user_id, b.cor_profile_id, b.company_id, b.request_id, b.venue_id,
		       COALESCE(cd.company_name, '')      AS company_name,
		       COALESCE(cp.first_name || ' ' || cp.last_name, '') AS profile_name,
		       COALESCE(v.name, '')               AS venue_name,
		       b.total_amount, b.status, b.special_requests, b.overstayed,
		       b.created_at, b.updated_at,
		       asg.id, asg.room_id, asg.check_in, asg.check_out, asg.status, COALESCE(r.name, '')
		FROM bookings b
		LEFT JOIN cor_company_details cd ON cd.id = b.company_id
		LEFT JOIN cor_profiles        cp ON cp.id = b.cor_profile_id
		LEFT JOIN venues              v  ON v.id  = b.venue_id
		LEFT JOIN LATERAL (
		    SELECT bra.id, bra.room_id, bra.check_in, bra.check_out, bra.status
		    FROM booking_room_assignments bra
		    WHERE bra.booking_id = b.id
		    ORDER BY bra.check_in ASC LIMIT 1
		) asg ON TRUE
		LEFT JOIN rooms r ON r.id = asg.room_id
		WHERE %s
		ORDER BY b.created_at DESC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		var branchID, webUserID, corProfileID, companyID, requestID, venueID uuid.NullUUID
		var bookerEmail, bookerPhone, specialRequests sql.NullString
		var companyName, profileName, venueName sql.NullString
		// lead assignment columns (nullable — individual only)
		var asgID uuid.NullUUID
		var asgRoomID uuid.NullUUID
		var asgCheckIn, asgCheckOut sql.NullTime
		var asgStatus, asgRoomName sql.NullString

		if err := rows.Scan(
			&b.ID, &b.BookingNumber, &b.OrgID, &branchID,
			&b.BookingType, &b.BookerType,
			&b.BookerName, &bookerEmail, &bookerPhone,
			&webUserID, &corProfileID, &companyID, &requestID, &venueID,
			&companyName, &profileName, &venueName,
			&b.TotalAmount, &b.Status, &specialRequests, &b.Overstayed,
			&b.CreatedAt, &b.UpdatedAt,
			&asgID, &asgRoomID, &asgCheckIn, &asgCheckOut, &asgStatus, &asgRoomName,
		); err != nil {
			return nil, 0, err
		}
		if branchID.Valid { b.BranchID = &branchID.UUID }
		if webUserID.Valid { b.WebUserID = &webUserID.UUID }
		if corProfileID.Valid { b.CorProfileID = &corProfileID.UUID }
		if companyID.Valid { b.CompanyID = &companyID.UUID }
		if requestID.Valid { b.RequestID = &requestID.UUID }
		if venueID.Valid { b.VenueID = &venueID.UUID }
		if bookerEmail.Valid { b.BookerEmail = bookerEmail.String }
		if bookerPhone.Valid { b.BookerPhone = bookerPhone.String }
		if specialRequests.Valid { b.SpecialRequests = specialRequests.String }
		if companyName.Valid { b.CompanyName = companyName.String }
		if profileName.Valid { b.ProfileName = profileName.String }
		if venueName.Valid { b.VenueName = venueName.String }

		if asgID.Valid {
			nights := 0
			if asgCheckIn.Valid && asgCheckOut.Valid {
				nights = int(asgCheckOut.Time.Sub(asgCheckIn.Time).Hours() / 24)
			}
			a := models.BookingRoomAssignment{
				ID:       asgID.UUID,
				BookingID: b.ID,
				Status:   asgStatus.String,
				RoomName: asgRoomName.String,
				Nights:   nights,
			}
			if asgRoomID.Valid { a.RoomID = asgRoomID.UUID }
			if asgCheckIn.Valid { a.CheckIn = asgCheckIn.Time }
			if asgCheckOut.Valid { a.CheckOut = asgCheckOut.Time }
			b.Assignments = []models.BookingRoomAssignment{a}
		}

		bookings = append(bookings, b)
	}
	return bookings, total, rows.Err()
}

func (r *BookingRepository) GetByIDUnscoped(id uuid.UUID) (*models.Booking, error) {
	row := r.db.QueryRow(`
		SELECT b.id, b.booking_number, b.org_id, b.branch_id,
		       b.booking_type, b.booker_type,
		       b.booker_name, b.booker_email, b.booker_phone,
		       b.web_user_id, b.cor_profile_id, b.company_id, b.request_id, b.venue_id,
		       COALESCE(cd.company_name, '')                   AS company_name,
		       COALESCE(cp.first_name || ' ' || cp.last_name, '') AS profile_name,
		       COALESCE(v.name, '')                            AS venue_name,
		       b.total_amount, b.status, b.special_requests, b.overstayed,
		       b.created_at, b.updated_at
		FROM bookings b
		LEFT JOIN cor_company_details cd ON cd.id = b.company_id
		LEFT JOIN cor_profiles        cp ON cp.id = b.cor_profile_id
		LEFT JOIN venues              v  ON v.id  = b.venue_id
		WHERE b.id = $1`, id)
	return scanBooking(row)
}

func (r *BookingRepository) ListByWebUserID(webUserID uuid.UUID, page, pageSize int) ([]models.Booking, int, error) {
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM bookings WHERE web_user_id = $1`, webUserID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT b.id, b.booking_number, b.org_id, b.branch_id,
		       b.booking_type, b.booker_type,
		       b.booker_name, b.booker_email, b.booker_phone,
		       b.web_user_id, b.cor_profile_id, b.company_id, b.request_id, b.venue_id,
		       COALESCE(cd.company_name, '')                   AS company_name,
		       COALESCE(cp.first_name || ' ' || cp.last_name, '') AS profile_name,
		       COALESCE(v.name, '')                            AS venue_name,
		       b.total_amount, b.status, b.special_requests, b.overstayed,
		       b.created_at, b.updated_at
		FROM bookings b
		LEFT JOIN cor_company_details cd ON cd.id = b.company_id
		LEFT JOIN cor_profiles        cp ON cp.id = b.cor_profile_id
		LEFT JOIN venues              v  ON v.id  = b.venue_id
		WHERE b.web_user_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3`, webUserID, pageSize, (page-1)*pageSize)
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

func (r *BookingRepository) UpdateStatus(tx *sql.Tx, id, orgID uuid.UUID, status string) error {
	_, err := tx.Exec(`
		UPDATE bookings SET status=$1, updated_at=$2 WHERE id=$3 AND org_id=$4`,
		status, time.Now(), id, orgID)
	return err
}

func (r *BookingRepository) UpdateTotalAmount(tx *sql.Tx, id, orgID uuid.UUID, amount float64) error {
	_, err := tx.Exec(`
		UPDATE bookings SET total_amount=$1, updated_at=$2 WHERE id=$3 AND org_id=$4`,
		amount, time.Now(), id, orgID)
	return err
}

// UpdateVenueAndTotal pins the chosen venue and hire charge on an event booking.
// The booking row is inserted before the venue is resolved, so this writes both
// back in one update during materialisation.
func (r *BookingRepository) UpdateVenueAndTotal(tx *sql.Tx, id, orgID, venueID uuid.UUID, amount float64) error {
	_, err := tx.Exec(`
		UPDATE bookings SET venue_id=$1, total_amount=$2, updated_at=$3 WHERE id=$4 AND org_id=$5`,
		venueID, amount, time.Now(), id, orgID)
	return err
}

func (r *BookingRepository) ExtendCheckout(id, orgID uuid.UUID, newDate time.Time) error {
	_, err := r.db.Exec(`
		UPDATE bookings SET updated_at=$1 WHERE id=$2 AND org_id=$3`,
		time.Now(), id, orgID)
	return err
}

func (r *BookingRepository) MarkOverstayed(id, orgID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE bookings SET overstayed=TRUE, updated_at=$1 WHERE id=$2 AND org_id=$3`,
		time.Now(), id, orgID)
	return err
}

// OverdueBookingRef carries details the nightly job needs to flag overstayed bookings.
type OverdueBookingRef struct {
	ID            uuid.UUID
	OrgID         uuid.UUID
	BookingNumber string
	BookerName    string
	OriginalCheckOut time.Time
}

func (r *BookingRepository) FindOverdueCheckouts(orgIDs []uuid.UUID) ([]OverdueBookingRef, error) {
	rows, err := r.db.Query(`
		SELECT b.id, b.org_id, b.booking_number, b.booker_name,
		       MAX(a.check_out) AS latest_checkout
		FROM bookings b
		JOIN booking_room_assignments a ON a.booking_id = b.id
		WHERE b.status = 'checked_in'
		  AND a.check_out < CURRENT_DATE
		  AND b.org_id = ANY($1)
		GROUP BY b.id, b.org_id, b.booking_number, b.booker_name`, orgIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []OverdueBookingRef
	for rows.Next() {
		var ref OverdueBookingRef
		if err := rows.Scan(&ref.ID, &ref.OrgID, &ref.BookingNumber, &ref.BookerName, &ref.OriginalCheckOut); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, rows.Err()
}

type bookingScanner interface {
	Scan(dest ...interface{}) error
}

func scanBooking(row bookingScanner) (*models.Booking, error) {
	var b models.Booking
	var branchID, webUserID, corProfileID, companyID, requestID, venueID uuid.NullUUID
	var bookerEmail, bookerPhone, specialRequests sql.NullString
	var companyName, profileName, venueName sql.NullString

	err := row.Scan(
		&b.ID, &b.BookingNumber, &b.OrgID, &branchID,
		&b.BookingType, &b.BookerType,
		&b.BookerName, &bookerEmail, &bookerPhone,
		&webUserID, &corProfileID, &companyID, &requestID, &venueID,
		&companyName, &profileName, &venueName,
		&b.TotalAmount, &b.Status, &specialRequests, &b.Overstayed,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if branchID.Valid {
		b.BranchID = &branchID.UUID
	}
	if webUserID.Valid {
		b.WebUserID = &webUserID.UUID
	}
	if corProfileID.Valid {
		b.CorProfileID = &corProfileID.UUID
	}
	if companyID.Valid {
		b.CompanyID = &companyID.UUID
	}
	if requestID.Valid {
		b.RequestID = &requestID.UUID
	}
	if venueID.Valid {
		b.VenueID = &venueID.UUID
	}
	if bookerEmail.Valid {
		b.BookerEmail = bookerEmail.String
	}
	if bookerPhone.Valid {
		b.BookerPhone = bookerPhone.String
	}
	if specialRequests.Valid {
		b.SpecialRequests = specialRequests.String
	}
	if companyName.Valid {
		b.CompanyName = companyName.String
	}
	if profileName.Valid {
		b.ProfileName = profileName.String
	}
	if venueName.Valid {
		b.VenueName = venueName.String
	}
	return &b, nil
}
