package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/interfaces"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type LeaveRequestRepository struct {
	db *sql.DB
}

func NewLeaveRequestRepository() *LeaveRequestRepository {
	return &LeaveRequestRepository{db: database.DB}
}

func (r *LeaveRequestRepository) Create(req *models.LeaveRequest) error {
	req.ID = uuid.New()
	now := time.Now()
	req.CreatedAt = now
	req.UpdatedAt = now
	req.Status = models.LeaveStatusPending
	_, err := r.db.Exec(`
		INSERT INTO leave_requests
		(id, employee_id, leave_type_id, start_date, end_date, total_days, reason, status,
		 attachment_url, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		req.ID, req.EmployeeID, req.LeaveTypeID, req.StartDate, req.EndDate, req.TotalDays,
		req.Reason, req.Status, req.AttachmentURL, req.CreatedAt, req.UpdatedAt,
	)
	return err
}

func (r *LeaveRequestRepository) GetByID(id uuid.UUID) (*models.LeaveRequest, error) {
	row := r.db.QueryRow(`
		SELECT lr.id, lr.employee_id, lr.leave_type_id, lr.start_date, lr.end_date, lr.total_days,
		       lr.reason, lr.status, lr.reviewed_by, lr.reviewed_at, lr.review_comment,
		       lr.attachment_url, lr.created_at, lr.updated_at,
		       lt.id, lt.name, lt.code
		FROM leave_requests lr
		JOIN leave_types lt ON lr.leave_type_id=lt.id
		WHERE lr.id=$1`, id)
	return r.scanOne(row)
}

func (r *LeaveRequestRepository) List(filter interfaces.LeaveRequestFilter, page, pageSize int) ([]models.LeaveRequest, int, error) {
	args := []interface{}{}
	where := []string{}
	i := 1

	if filter.EmployeeID != nil {
		where = append(where, fmt.Sprintf("lr.employee_id=$%d", i))
		args = append(args, *filter.EmployeeID)
		i++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("lr.status=$%d", i))
		args = append(args, filter.Status)
		i++
	}
	if filter.StartDateGTE != nil {
		where = append(where, fmt.Sprintf("lr.start_date>=$%d", i))
		args = append(args, *filter.StartDateGTE)
		i++
	}
	if filter.EndDateLTE != nil {
		where = append(where, fmt.Sprintf("lr.end_date<=$%d", i))
		args = append(args, *filter.EndDateLTE)
		i++
	}

	whereStr := "1=1"
	if len(where) > 0 {
		whereStr = strings.Join(where, " AND ")
	}

	var total int
	err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM leave_requests lr WHERE %s`, whereStr), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT lr.id, lr.employee_id, lr.leave_type_id, lr.start_date, lr.end_date, lr.total_days,
		       lr.reason, lr.status, lr.reviewed_by, lr.reviewed_at, lr.review_comment,
		       lr.attachment_url, lr.created_at, lr.updated_at,
		       lt.id, lt.name, lt.code
		FROM leave_requests lr
		JOIN leave_types lt ON lr.leave_type_id=lt.id
		WHERE %s ORDER BY lr.created_at DESC LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []models.LeaveRequest
	for rows.Next() {
		req, err := r.scanOne(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *req)
	}
	return out, total, rows.Err()
}

func (r *LeaveRequestRepository) UpdateStatus(id uuid.UUID, status models.LeaveRequestStatus, reviewedBy *uuid.UUID, comment string) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE leave_requests SET status=$1, reviewed_by=$2, reviewed_at=$3, review_comment=$4, updated_at=$3
		WHERE id=$5`,
		status, reviewedBy, now, comment, id)
	return err
}

func (r *LeaveRequestRepository) HasOverlap(employeeID uuid.UUID, start, end time.Time, excludeID *uuid.UUID) (bool, error) {
	q := `SELECT COUNT(1) FROM leave_requests
		  WHERE employee_id=$1 AND status NOT IN ('rejected','cancelled')
		  AND start_date <= $2 AND end_date >= $3`
	args := []interface{}{employeeID, end, start}
	if excludeID != nil {
		q += " AND id!=$4"
		args = append(args, *excludeID)
	}
	var count int
	err := r.db.QueryRow(q, args...).Scan(&count)
	return count > 0, err
}

func (r *LeaveRequestRepository) scanOne(row rowScanner) (*models.LeaveRequest, error) {
	var req models.LeaveRequest
	var reviewedBy sql.NullString
	var reviewedAt sql.NullTime
	var ltID uuid.UUID
	var ltName, ltCode string

	err := row.Scan(
		&req.ID, &req.EmployeeID, &req.LeaveTypeID, &req.StartDate, &req.EndDate, &req.TotalDays,
		&req.Reason, &req.Status, &reviewedBy, &reviewedAt, &req.ReviewComment,
		&req.AttachmentURL, &req.CreatedAt, &req.UpdatedAt,
		&ltID, &ltName, &ltCode,
	)
	if err != nil {
		return nil, err
	}
	if reviewedBy.Valid {
		id, _ := uuid.Parse(reviewedBy.String)
		req.ReviewedBy = &id
	}
	if reviewedAt.Valid {
		req.ReviewedAt = &reviewedAt.Time
	}
	req.LeaveType = &models.LeaveType{ID: ltID, Name: ltName, Code: ltCode}
	return &req, nil
}
