package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type IndividualBookingRequestRepository struct {
	db *sql.DB
}

func NewIndividualBookingRequestRepository() *IndividualBookingRequestRepository {
	return &IndividualBookingRequestRepository{db: database.DB}
}

func (r *IndividualBookingRequestRepository) Create(req *models.IndividualBookingRequest) error {
	return r.db.QueryRow(`
		INSERT INTO individual_booking_requests
			(org_id, web_user_id, booker_name, booker_email, booker_phone,
			 booking_type, status, notes, documents, payload)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, created_at, updated_at`,
		req.OrgID, req.WebUserID, req.BookerName, req.BookerEmail, req.BookerPhone,
		req.BookingType, req.Status, req.Notes,
		pq.Array(req.Documents), []byte(req.Payload),
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *IndividualBookingRequestRepository) GetByID(id, orgID uuid.UUID) (*models.IndividualBookingRequest, error) {
	req := &models.IndividualBookingRequest{}
	var payloadBytes []byte
	var roomName sql.NullString

	err := r.db.QueryRow(`
		SELECT
			ibr.id, ibr.org_id, ibr.web_user_id,
			ibr.booker_name, ibr.booker_email, ibr.booker_phone,
			ibr.booking_type, ibr.status, ibr.notes,
			ibr.documents, ibr.payload, ibr.created_at, ibr.updated_at,
			r.name
		FROM individual_booking_requests ibr
		LEFT JOIN rooms r ON r.id = (ibr.payload->>'room_id')::uuid
		WHERE ibr.id = $1 AND ibr.org_id = $2`, id, orgID,
	).Scan(
		&req.ID, &req.OrgID, &req.WebUserID,
		&req.BookerName, &req.BookerEmail, &req.BookerPhone,
		&req.BookingType, &req.Status, &req.Notes,
		pq.Array(&req.Documents), &payloadBytes, &req.CreatedAt, &req.UpdatedAt,
		&roomName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("booking request not found")
	}
	if err != nil {
		return nil, err
	}

	if len(payloadBytes) > 0 {
		req.Payload = json.RawMessage(payloadBytes)
	}
	req.RoomName = roomName.String
	return req, nil
}

func (r *IndividualBookingRequestRepository) GetByIDForWebUser(id, webUserID uuid.UUID) (*models.IndividualBookingRequest, error) {
	req := &models.IndividualBookingRequest{}
	var payloadBytes []byte
	var roomName sql.NullString

	err := r.db.QueryRow(`
		SELECT
			ibr.id, ibr.org_id, ibr.web_user_id,
			ibr.booker_name, ibr.booker_email, ibr.booker_phone,
			ibr.booking_type, ibr.status, ibr.notes,
			ibr.documents, ibr.payload, ibr.created_at, ibr.updated_at,
			r.name
		FROM individual_booking_requests ibr
		LEFT JOIN rooms r ON r.id = (ibr.payload->>'room_id')::uuid
		WHERE ibr.id = $1 AND ibr.web_user_id = $2`, id, webUserID,
	).Scan(
		&req.ID, &req.OrgID, &req.WebUserID,
		&req.BookerName, &req.BookerEmail, &req.BookerPhone,
		&req.BookingType, &req.Status, &req.Notes,
		pq.Array(&req.Documents), &payloadBytes, &req.CreatedAt, &req.UpdatedAt,
		&roomName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("booking request not found")
	}
	if err != nil {
		return nil, err
	}

	if len(payloadBytes) > 0 {
		req.Payload = json.RawMessage(payloadBytes)
	}
	req.RoomName = roomName.String
	return req, nil
}

func (r *IndividualBookingRequestRepository) ListByWebUser(webUserID uuid.UUID, page, pageSize int) ([]models.IndividualBookingRequest, int, error) {
	offset := (page - 1) * pageSize

	var total int
	if err := r.db.QueryRow(
		`SELECT COUNT(*) FROM individual_booking_requests WHERE web_user_id = $1`, webUserID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			ibr.id, ibr.org_id, ibr.web_user_id,
			ibr.booker_name, ibr.booker_email, ibr.booker_phone,
			ibr.booking_type, ibr.status, ibr.notes,
			ibr.documents, ibr.created_at, ibr.updated_at,
			r.name
		FROM individual_booking_requests ibr
		LEFT JOIN rooms r ON r.id = (ibr.payload->>'room_id')::uuid
		WHERE ibr.web_user_id = $1
		ORDER BY ibr.created_at DESC
		LIMIT $2 OFFSET $3`, webUserID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reqs []models.IndividualBookingRequest
	for rows.Next() {
		var req models.IndividualBookingRequest
		var roomName sql.NullString
		if err := rows.Scan(
			&req.ID, &req.OrgID, &req.WebUserID,
			&req.BookerName, &req.BookerEmail, &req.BookerPhone,
			&req.BookingType, &req.Status, &req.Notes,
			pq.Array(&req.Documents), &req.CreatedAt, &req.UpdatedAt,
			&roomName,
		); err != nil {
			return nil, 0, err
		}
		req.RoomName = roomName.String
		reqs = append(reqs, req)
	}
	return reqs, total, nil
}

func (r *IndividualBookingRequestRepository) List(orgID uuid.UUID, status string, page, pageSize int) ([]models.IndividualBookingRequest, int, error) {
	offset := (page - 1) * pageSize

	args := []interface{}{orgID}
	where := "ibr.org_id = $1"
	i := 2
	if status != "" {
		where += " AND ibr.status = $" + itoa(i)
		args = append(args, status)
		i++
	}

	var total int
	if err := r.db.QueryRow(
		`SELECT COUNT(*) FROM individual_booking_requests ibr WHERE `+where, args...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, offset)
	rows, err := r.db.Query(`
		SELECT
			ibr.id, ibr.org_id, ibr.web_user_id,
			ibr.booker_name, ibr.booker_email, ibr.booker_phone,
			ibr.booking_type, ibr.status, ibr.notes,
			ibr.documents, ibr.created_at, ibr.updated_at,
			r.name
		FROM individual_booking_requests ibr
		LEFT JOIN rooms r ON r.id = (ibr.payload->>'room_id')::uuid
		WHERE `+where+`
		ORDER BY ibr.created_at DESC
		LIMIT $`+itoa(i)+` OFFSET $`+itoa(i+1), args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reqs []models.IndividualBookingRequest
	for rows.Next() {
		var req models.IndividualBookingRequest
		var roomName sql.NullString
		if err := rows.Scan(
			&req.ID, &req.OrgID, &req.WebUserID,
			&req.BookerName, &req.BookerEmail, &req.BookerPhone,
			&req.BookingType, &req.Status, &req.Notes,
			pq.Array(&req.Documents), &req.CreatedAt, &req.UpdatedAt,
			&roomName,
		); err != nil {
			return nil, 0, err
		}
		req.RoomName = roomName.String
		reqs = append(reqs, req)
	}
	return reqs, total, nil
}

func (r *IndividualBookingRequestRepository) UpdateStatus(id, orgID uuid.UUID, status string) error {
	res, err := r.db.Exec(`
		UPDATE individual_booking_requests SET status = $1, updated_at = NOW()
		WHERE id = $2 AND org_id = $3`, status, id, orgID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("booking request not found")
	}
	return nil
}
