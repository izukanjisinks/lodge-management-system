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

type PayslipRepository struct {
	db *sql.DB
}

func NewPayslipRepository() *PayslipRepository {
	return &PayslipRepository{db: database.DB}
}

func (r *PayslipRepository) Create(p *models.Payslip) error {
	p.ID = uuid.New()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO payslips (id, employee_id, month, year, base_salary, housing_allowance, transport_allowance, medical_allowance, gross_salary, income_tax, leave_days, net_salary, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		p.ID, p.EmployeeID, p.Month, p.Year, p.BaseSalary, p.HousingAllowance,
		p.TransportAllowance, p.MedicalAllowance, p.GrossSalary, p.IncomeTax,
		p.LeaveDays, p.NetSalary, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PayslipRepository) GetByID(id uuid.UUID) (*models.Payslip, error) {
	row := r.db.QueryRow(`
		SELECT p.id, p.employee_id, p.month, p.year, p.base_salary, p.housing_allowance,
		       p.transport_allowance, p.medical_allowance, p.gross_salary, p.income_tax,
		       p.leave_days, p.net_salary, p.created_at, p.updated_at,
		       CONCAT(e.first_name, ' ', e.last_name) AS employee_name,
		       COALESCE(pos.title, '') AS position_name
		FROM payslips p
		JOIN employees e ON p.employee_id = e.id
		LEFT JOIN positions pos ON e.position_id = pos.id
		WHERE p.id=$1`, id)
	return r.scanOne(row)
}

func (r *PayslipRepository) GetByEmployeeAndPeriod(employeeID uuid.UUID, month, year int) (*models.Payslip, error) {
	row := r.db.QueryRow(`
		SELECT p.id, p.employee_id, p.month, p.year, p.base_salary, p.housing_allowance,
		       p.transport_allowance, p.medical_allowance, p.gross_salary, p.income_tax,
		       p.leave_days, p.net_salary, p.created_at, p.updated_at,
		       CONCAT(e.first_name, ' ', e.last_name) AS employee_name,
		       COALESCE(pos.title, '') AS position_name
		FROM payslips p
		JOIN employees e ON p.employee_id = e.id
		LEFT JOIN positions pos ON e.position_id = pos.id
		WHERE p.employee_id=$1 AND p.month=$2 AND p.year=$3`, employeeID, month, year)
	return r.scanOne(row)
}

func (r *PayslipRepository) List(employeeID *uuid.UUID, month *int, year *int, page, pageSize int) ([]models.Payslip, int, error) {
	args := []interface{}{}
	where := []string{}
	i := 1

	if employeeID != nil {
		where = append(where, fmt.Sprintf("p.employee_id=$%d", i))
		args = append(args, *employeeID)
		i++
	}
	if month != nil {
		where = append(where, fmt.Sprintf("p.month=$%d", i))
		args = append(args, *month)
		i++
	}
	if year != nil {
		where = append(where, fmt.Sprintf("p.year=$%d", i))
		args = append(args, *year)
		i++
	}

	whereStr := "1=1"
	if len(where) > 0 {
		whereStr = strings.Join(where, " AND ")
	}

	var total int
	err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM payslips p WHERE %s`, whereStr), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT p.id, p.employee_id, p.month, p.year, p.base_salary, p.housing_allowance,
		       p.transport_allowance, p.medical_allowance, p.gross_salary, p.income_tax,
		       p.leave_days, p.net_salary, p.created_at, p.updated_at,
		       CONCAT(e.first_name, ' ', e.last_name) AS employee_name,
		       COALESCE(pos.title, '') AS position_name
		FROM payslips p
		JOIN employees e ON p.employee_id = e.id
		LEFT JOIN positions pos ON e.position_id = pos.id
		WHERE %s
		ORDER BY p.year DESC, p.month DESC, e.last_name
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payslips []models.Payslip
	for rows.Next() {
		p, err := r.scanRow(rows)
		if err != nil {
			return nil, 0, err
		}
		payslips = append(payslips, *p)
	}
	return payslips, total, rows.Err()
}

func (r *PayslipRepository) Delete(id uuid.UUID) error {
	result, err := r.db.Exec(`DELETE FROM payslips WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("payslip not found")
	}
	return nil
}

func (r *PayslipRepository) scanOne(row *sql.Row) (*models.Payslip, error) {
	return r.scanRow(row)
}

func (r *PayslipRepository) scanRow(row rowScanner) (*models.Payslip, error) {
	var p models.Payslip
	err := row.Scan(
		&p.ID, &p.EmployeeID, &p.Month, &p.Year, &p.BaseSalary, &p.HousingAllowance,
		&p.TransportAllowance, &p.MedicalAllowance, &p.GrossSalary, &p.IncomeTax,
		&p.LeaveDays, &p.NetSalary, &p.CreatedAt, &p.UpdatedAt,
		&p.EmployeeName, &p.PositionName,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
