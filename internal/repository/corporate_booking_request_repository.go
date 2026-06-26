package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CorporateBookingRequestRepository struct {
	db *sql.DB
}

func NewCorporateBookingRequestRepository() *CorporateBookingRequestRepository {
	return &CorporateBookingRequestRepository{db: database.DB}
}

func (r *CorporateBookingRequestRepository) Create(req *models.CorporateBookingRequest) error {
	var payloadBytes []byte
	if req.Payload != nil {
		payloadBytes = []byte(req.Payload)
	}
	// documents is NOT NULL in the DB; a nil slice would marshal to SQL NULL.
	documents := req.Documents
	if documents == nil {
		documents = []string{}
	}
	return r.db.QueryRow(`
		INSERT INTO corporate_booking_requests
			(org_id, branch_id, cor_profile_id, company_id, web_user_id, booking_type, status,
			 reason_for_booking, notes,
			 authoriser_name, authoriser_email, authoriser_phone,
			 authoriser_title, authoriser_department, authoriser_gl_code,
			 documents, payload)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		RETURNING id, created_at, updated_at`,
		req.OrgID, req.BranchID, req.CorProfileID, req.CompanyID, req.WebUserID,
		req.BookingType, req.Status,
		req.ReasonForBooking, req.Notes,
		req.AuthoriserName, req.AuthoriserEmail, req.AuthoriserPhone,
		req.AuthoriserTitle, req.AuthoriserDepartment, req.AuthoriserGLCode,
		pq.Array(documents), payloadBytes,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *CorporateBookingRequestRepository) GetByID(id, orgID uuid.UUID) (*models.CorporateBookingRequest, error) {
	req := &models.CorporateBookingRequest{}
	var payloadBytes []byte
	var branchName, profileName, companyName sql.NullString

	err := r.db.QueryRow(`
		SELECT
			cbr.id, cbr.org_id, cbr.branch_id, cbr.cor_profile_id, cbr.company_id, cbr.web_user_id,
			cbr.booking_type, cbr.status,
			cbr.reason_for_booking, cbr.notes,
			cbr.authoriser_name, cbr.authoriser_email, cbr.authoriser_phone,
			cbr.authoriser_title, cbr.authoriser_department, cbr.authoriser_gl_code,
			cbr.documents, cbr.payload, cbr.created_at, cbr.updated_at,
			c.company_name, b.name, CONCAT(p.first_name, ' ', p.last_name)
		FROM corporate_booking_requests cbr
		LEFT JOIN cor_company_details c  ON c.id = cbr.company_id
		LEFT JOIN cor_branch_details b   ON b.id = cbr.branch_id
		LEFT JOIN cor_profiles p         ON p.id = cbr.cor_profile_id
		WHERE cbr.id = $1 AND cbr.org_id = $2`, id, orgID,
	).Scan(
		&req.ID, &req.OrgID, &req.BranchID, &req.CorProfileID, &req.CompanyID, &req.WebUserID,
		&req.BookingType, &req.Status,
		&req.ReasonForBooking, &req.Notes,
		&req.AuthoriserName, &req.AuthoriserEmail, &req.AuthoriserPhone,
		&req.AuthoriserTitle, &req.AuthoriserDepartment, &req.AuthoriserGLCode,
		pq.Array(&req.Documents), &payloadBytes, &req.CreatedAt, &req.UpdatedAt,
		&companyName, &branchName, &profileName,
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
	req.CompanyName = companyName.String
	req.BranchName = branchName.String
	req.ProfileName = profileName.String
	return req, nil
}

func (r *CorporateBookingRequestRepository) List(orgID uuid.UUID, bookingType, status string, page, pageSize int) ([]models.CorporateBookingRequest, int, error) {
	offset := (page - 1) * pageSize

	args := []interface{}{orgID}
	where := "cbr.org_id = $1"
	i := 2
	if bookingType != "" {
		where += " AND cbr.booking_type = $" + itoa(i)
		args = append(args, bookingType)
		i++
	}
	if status != "" {
		where += " AND cbr.status = $" + itoa(i)
		args = append(args, status)
		i++
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM corporate_booking_requests cbr WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, offset)
	rows, err := r.db.Query(`
		SELECT
			cbr.id, cbr.org_id, cbr.branch_id, cbr.cor_profile_id, cbr.company_id, cbr.web_user_id,
			cbr.booking_type, cbr.status,
			cbr.reason_for_booking, cbr.notes,
			cbr.authoriser_name, cbr.authoriser_email, cbr.authoriser_phone,
			cbr.authoriser_title, cbr.authoriser_department, cbr.authoriser_gl_code,
			cbr.documents, cbr.created_at, cbr.updated_at,
			c.company_name, b.name, CONCAT(p.first_name, ' ', p.last_name)
		FROM corporate_booking_requests cbr
		LEFT JOIN cor_company_details c ON c.id = cbr.company_id
		LEFT JOIN cor_branch_details b  ON b.id = cbr.branch_id
		LEFT JOIN cor_profiles p        ON p.id = cbr.cor_profile_id
		WHERE `+where+`
		ORDER BY cbr.created_at DESC
		LIMIT $`+itoa(i)+` OFFSET $`+itoa(i+1), args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reqs []models.CorporateBookingRequest
	for rows.Next() {
		var req models.CorporateBookingRequest
		var companyName, branchName, profileName sql.NullString
		if err := rows.Scan(
			&req.ID, &req.OrgID, &req.BranchID, &req.CorProfileID, &req.CompanyID, &req.WebUserID,
			&req.BookingType, &req.Status,
			&req.ReasonForBooking, &req.Notes,
			&req.AuthoriserName, &req.AuthoriserEmail, &req.AuthoriserPhone,
			&req.AuthoriserTitle, &req.AuthoriserDepartment, &req.AuthoriserGLCode,
			pq.Array(&req.Documents), &req.CreatedAt, &req.UpdatedAt,
			&companyName, &branchName, &profileName,
		); err != nil {
			return nil, 0, err
		}
		req.CompanyName = companyName.String
		req.BranchName = branchName.String
		req.ProfileName = profileName.String
		reqs = append(reqs, req)
	}
	return reqs, total, nil
}

func (r *CorporateBookingRequestRepository) UpdateStatus(id, orgID uuid.UUID, status string) error {
	res, err := r.db.Exec(`
		UPDATE corporate_booking_requests SET status = $1, updated_at = NOW()
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

func itoa(i int) string {
	return strconv.Itoa(i)
}
