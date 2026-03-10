package repository

import (
	"database/sql"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type LeaveBalanceRepository struct {
	db *sql.DB
}

func NewLeaveBalanceRepository() *LeaveBalanceRepository {
	return &LeaveBalanceRepository{db: database.DB}
}

func (r *LeaveBalanceRepository) GetByEmployeeAndYear(employeeID uuid.UUID, year int) ([]models.LeaveBalance, error) {
	rows, err := r.db.Query(`
		SELECT lb.id, lb.employee_id, lb.leave_type_id, lb.year, lb.total_entitled, lb.used, lb.pending,
		       lb.carried_forward, lb.earned_leave_days, lb.created_at, lb.updated_at,
		       lt.id, lt.name, lt.code
		FROM leave_balances lb
		JOIN leave_types lt ON lb.leave_type_id = lt.id
		WHERE lb.employee_id=$1 AND lb.year=$2
		ORDER BY lt.name`, employeeID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *LeaveBalanceRepository) GetByEmployeeTypeYear(employeeID, leaveTypeID uuid.UUID, year int) (*models.LeaveBalance, error) {
	row := r.db.QueryRow(`
		SELECT lb.id, lb.employee_id, lb.leave_type_id, lb.year, lb.total_entitled, lb.used, lb.pending,
		       lb.carried_forward, lb.earned_leave_days, lb.created_at, lb.updated_at,
		       lt.id, lt.name, lt.code
		FROM leave_balances lb
		JOIN leave_types lt ON lb.leave_type_id = lt.id
		WHERE lb.employee_id=$1 AND lb.leave_type_id=$2 AND lb.year=$3`,
		employeeID, leaveTypeID, year)
	return r.scanOne(row)
}

// GetAllByYear returns every leave_balance row for the given year across all employees,
// joining leave_types so carry-forward rules are available.
func (r *LeaveBalanceRepository) GetAllByYear(year int) ([]models.LeaveBalance, error) {
	rows, err := r.db.Query(`
		SELECT lb.id, lb.employee_id, lb.leave_type_id, lb.year, lb.total_entitled, lb.used, lb.pending,
		       lb.carried_forward, lb.earned_leave_days, lb.created_at, lb.updated_at,
		       lt.id, lt.name, lt.code
		FROM leave_balances lb
		JOIN leave_types lt ON lb.leave_type_id = lt.id
		WHERE lb.year=$1
		ORDER BY lb.employee_id, lt.name`, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *LeaveBalanceRepository) SetCarriedForward(id uuid.UUID, days int) error {
	_, err := r.db.Exec(`
		UPDATE leave_balances SET carried_forward=$1, updated_at=NOW() WHERE id=$2`,
		days, id)
	return err
}

func (r *LeaveBalanceRepository) Upsert(lb *models.LeaveBalance) error {
	if lb.ID == uuid.Nil {
		lb.ID = uuid.New()
	}
	now := time.Now()
	_, err := r.db.Exec(`
		INSERT INTO leave_balances
		(id, employee_id, leave_type_id, year, total_entitled, used, pending, carried_forward, earned_leave_days, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (employee_id, leave_type_id, year) DO UPDATE SET
		  total_entitled=EXCLUDED.total_entitled,
		  carried_forward=EXCLUDED.carried_forward,
		  earned_leave_days=EXCLUDED.earned_leave_days,
		  updated_at=EXCLUDED.updated_at`,
		lb.ID, lb.EmployeeID, lb.LeaveTypeID, lb.Year, lb.TotalEntitled,
		lb.Used, lb.Pending, lb.CarriedForward, lb.EarnedLeaveDays, now, now,
	)
	return err
}

func (r *LeaveBalanceRepository) IncrementPending(employeeID, leaveTypeID uuid.UUID, year, days int) error {
	_, err := r.db.Exec(`
		UPDATE leave_balances SET pending=pending+$1, updated_at=NOW()
		WHERE employee_id=$2 AND leave_type_id=$3 AND year=$4`,
		days, employeeID, leaveTypeID, year)
	return err
}

func (r *LeaveBalanceRepository) DecrementPending(employeeID, leaveTypeID uuid.UUID, year, days int) error {
	_, err := r.db.Exec(`
		UPDATE leave_balances SET pending=GREATEST(0, pending-$1), updated_at=NOW()
		WHERE employee_id=$2 AND leave_type_id=$3 AND year=$4`,
		days, employeeID, leaveTypeID, year)
	return err
}

// ApproveLeave moves days from pending → used
func (r *LeaveBalanceRepository) ApproveLeave(employeeID, leaveTypeID uuid.UUID, year, days int) error {
	_, err := r.db.Exec(`
		UPDATE leave_balances
		SET pending=GREATEST(0, pending-$1), used=used+$1, updated_at=NOW()
		WHERE employee_id=$2 AND leave_type_id=$3 AND year=$4`,
		days, employeeID, leaveTypeID, year)
	return err
}

func (r *LeaveBalanceRepository) Adjust(id uuid.UUID, delta int) error {
	_, err := r.db.Exec(`
		UPDATE leave_balances SET earned_leave_days=earned_leave_days+$1, updated_at=NOW() WHERE id=$2`,
		delta, id)
	return err
}

func (r *LeaveBalanceRepository) scanRows(rows *sql.Rows) ([]models.LeaveBalance, error) {
	var out []models.LeaveBalance
	for rows.Next() {
		lb, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *lb)
	}
	return out, rows.Err()
}

func (r *LeaveBalanceRepository) scanOne(row *sql.Row) (*models.LeaveBalance, error) {
	return r.scanRow(row)
}

func (r *LeaveBalanceRepository) scanRow(row rowScanner) (*models.LeaveBalance, error) {
	var lb models.LeaveBalance
	var ltID uuid.UUID
	var ltName, ltCode string
	err := row.Scan(
		&lb.ID, &lb.EmployeeID, &lb.LeaveTypeID, &lb.Year, &lb.TotalEntitled, &lb.Used, &lb.Pending,
		&lb.CarriedForward, &lb.EarnedLeaveDays, &lb.CreatedAt, &lb.UpdatedAt,
		&ltID, &ltName, &ltCode,
	)
	if err != nil {
		return nil, err
	}
	lb.Balance = lb.TotalEntitled + lb.CarriedForward + lb.EarnedLeaveDays - lb.Used - lb.Pending
	lb.LeaveType = &models.LeaveType{ID: ltID, Name: ltName, Code: ltCode}
	return &lb, nil
}
