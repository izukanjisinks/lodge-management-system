package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type PayrollRepository struct {
	db *sql.DB
}

func NewPayrollRepository() *PayrollRepository {
	return &PayrollRepository{db: database.DB}
}

func (r *PayrollRepository) Create(p *models.Payroll) error {
	p.ID = uuid.New()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO payrolls (id, start_date, end_date, status, processed_by, processed_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		p.ID, p.StartDate, p.EndDate, p.Status, p.ProcessedBy, p.ProcessedAt, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PayrollRepository) GetByID(id uuid.UUID) (*models.Payroll, error) {
	var p models.Payroll
	err := r.db.QueryRow(`
		SELECT p.id, p.start_date, p.end_date, p.status, p.processed_by, p.processed_at, p.created_at, p.updated_at,
		       COALESCE(u.email, '') AS processed_by_name
		FROM payrolls p
		LEFT JOIN users u ON p.processed_by = u.user_id
		WHERE p.id=$1`, id,
	).Scan(&p.ID, &p.StartDate, &p.EndDate, &p.Status, &p.ProcessedBy, &p.ProcessedAt,
		&p.CreatedAt, &p.UpdatedAt, &p.ProcessedByName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PayrollRepository) List(status string, page, pageSize int) ([]models.Payroll, int, error) {
	args := []interface{}{}
	where := []string{}
	i := 1

	if status != "" {
		where = append(where, fmt.Sprintf("p.status=$%d", i))
		args = append(args, status)
		i++
	}

	whereStr := "1=1"
	if len(where) > 0 {
		whereStr = strings.Join(where, " AND ")
	}

	var total int
	err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM payrolls p WHERE %s`, whereStr), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT p.id, p.start_date, p.end_date, p.status, p.processed_by, p.processed_at, p.created_at, p.updated_at,
		       COALESCE(u.email, '') AS processed_by_name
		FROM payrolls p
		LEFT JOIN users u ON p.processed_by = u.user_id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payrolls []models.Payroll
	for rows.Next() {
		var p models.Payroll
		if err := rows.Scan(&p.ID, &p.StartDate, &p.EndDate, &p.Status, &p.ProcessedBy, &p.ProcessedAt,
			&p.CreatedAt, &p.UpdatedAt, &p.ProcessedByName); err != nil {
			return nil, 0, err
		}
		payrolls = append(payrolls, p)
	}
	return payrolls, total, rows.Err()
}

func (r *PayrollRepository) Update(p *models.Payroll) error {
	p.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE payrolls SET status=$1, processed_by=$2, processed_at=$3, updated_at=$4
		WHERE id=$5`,
		p.Status, p.ProcessedBy, p.ProcessedAt, p.UpdatedAt, p.ID,
	)
	return err
}

func (r *PayrollRepository) Delete(id uuid.UUID) error {
	result, err := r.db.Exec(`DELETE FROM payrolls WHERE id=$1 AND status='OPEN'`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("payroll not found or not in OPEN status")
	}
	return nil
}
