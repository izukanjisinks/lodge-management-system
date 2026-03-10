package repository

import (
	"database/sql"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type AttendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository() *AttendanceRepository {
	return &AttendanceRepository{db: database.DB}
}

func (r *AttendanceRepository) Create(a *models.Attendance) error {
	a.ID = uuid.New()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO attendance
		(id, employee_id, date, clock_in, clock_out, total_hours, status, overtime_hours, notes, source, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		a.ID, a.EmployeeID, a.Date, a.ClockIn, a.ClockOut, a.TotalHours, a.Status,
		a.OvertimeHours, a.Notes, a.Source, a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *AttendanceRepository) GetByEmployeeAndDate(employeeID uuid.UUID, date time.Time) (*models.Attendance, error) {
	return r.scanOne(r.db.QueryRow(`
		SELECT id, employee_id, date, clock_in, clock_out, total_hours, status, overtime_hours, notes, source, created_at, updated_at
		FROM attendance WHERE employee_id=$1 AND date=$2`,
		employeeID, date.Format("2006-01-02")))
}

func (r *AttendanceRepository) GetByID(id uuid.UUID) (*models.Attendance, error) {
	return r.scanOne(r.db.QueryRow(`
		SELECT id, employee_id, date, clock_in, clock_out, total_hours, status, overtime_hours, notes, source, created_at, updated_at
		FROM attendance WHERE id=$1`, id))
}

func (r *AttendanceRepository) ListByEmployee(employeeID uuid.UUID, from, to time.Time) ([]models.Attendance, error) {
	rows, err := r.db.Query(`
		SELECT id, employee_id, date, clock_in, clock_out, total_hours, status, overtime_hours, notes, source, created_at, updated_at
		FROM attendance WHERE employee_id=$1 AND date>=$2 AND date<=$3 ORDER BY date`,
		employeeID, from.Format("2006-01-02"), to.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *AttendanceRepository) ListByDepartmentAndDate(departmentID uuid.UUID, date time.Time) ([]models.Attendance, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.employee_id, a.date, a.clock_in, a.clock_out, a.total_hours, a.status,
		       a.overtime_hours, a.notes, a.source, a.created_at, a.updated_at
		FROM attendance a
		JOIN employees e ON a.employee_id=e.id
		WHERE e.department_id=$1 AND a.date=$2 AND e.deleted_at IS NULL
		ORDER BY e.last_name`,
		departmentID, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanRows(rows)
}

func (r *AttendanceRepository) Update(a *models.Attendance) error {
	a.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE attendance SET clock_in=$1, clock_out=$2, total_hours=$3, status=$4,
		overtime_hours=$5, notes=$6, source=$7, updated_at=$8 WHERE id=$9`,
		a.ClockIn, a.ClockOut, a.TotalHours, a.Status,
		a.OvertimeHours, a.Notes, a.Source, a.UpdatedAt, a.ID,
	)
	return err
}

func (r *AttendanceRepository) SetClockIn(employeeID uuid.UUID, date time.Time, clockIn time.Time) (*models.Attendance, error) {
	existing, err := r.GetByEmployeeAndDate(employeeID, date)
	if err == sql.ErrNoRows {
		a := &models.Attendance{
			EmployeeID: employeeID,
			Date:       date,
			ClockIn:    &clockIn,
			Status:     models.AttendanceStatusPresent,
			Source:     models.AttendanceSourceSystem,
		}
		return a, r.Create(a)
	}
	if err != nil {
		return nil, err
	}
	existing.ClockIn = &clockIn
	existing.Status = models.AttendanceStatusPresent
	return existing, r.Update(existing)
}

func (r *AttendanceRepository) SetClockOut(employeeID uuid.UUID, date time.Time, clockOut time.Time) (*models.Attendance, error) {
	existing, err := r.GetByEmployeeAndDate(employeeID, date)
	if err != nil {
		return nil, err
	}
	existing.ClockOut = &clockOut
	// Compute hours
	if existing.ClockIn != nil {
		hours := clockOut.Sub(*existing.ClockIn).Hours()
		existing.TotalHours = hours
		const stdHours = 8.0
		if hours > stdHours {
			existing.OvertimeHours = hours - stdHours
		}
	}
	return existing, r.Update(existing)
}

func (r *AttendanceRepository) GetMonthlySummary(employeeID uuid.UUID, month, year int) (map[string]interface{}, error) {
	row := r.db.QueryRow(`
		SELECT
			COUNT(*) FILTER (WHERE status='present') AS days_present,
			COUNT(*) FILTER (WHERE status='absent') AS days_absent,
			COUNT(*) FILTER (WHERE status='on_leave') AS days_on_leave,
			COALESCE(SUM(total_hours),0) AS total_hours,
			COALESCE(SUM(overtime_hours),0) AS overtime_hours
		FROM attendance
		WHERE employee_id=$1
		AND EXTRACT(MONTH FROM date)=$2
		AND EXTRACT(YEAR FROM date)=$3`,
		employeeID, month, year)

	var present, absent, onLeave int
	var totalHours, overtimeHours float64
	err := row.Scan(&present, &absent, &onLeave, &totalHours, &overtimeHours)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"employee_id":    employeeID,
		"month":          month,
		"year":           year,
		"days_present":   present,
		"days_absent":    absent,
		"days_on_leave":  onLeave,
		"total_hours":    totalHours,
		"overtime_hours": overtimeHours,
	}, nil
}

func (r *AttendanceRepository) scanRows(rows *sql.Rows) ([]models.Attendance, error) {
	var out []models.Attendance
	for rows.Next() {
		a, err := r.scanOne(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return out, rows.Err()
}

func (r *AttendanceRepository) scanOne(row rowScanner) (*models.Attendance, error) {
	var a models.Attendance
	var clockIn, clockOut sql.NullTime
	err := row.Scan(
		&a.ID, &a.EmployeeID, &a.Date, &clockIn, &clockOut, &a.TotalHours, &a.Status,
		&a.OvertimeHours, &a.Notes, &a.Source, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if clockIn.Valid {
		a.ClockIn = &clockIn.Time
	}
	if clockOut.Valid {
		a.ClockOut = &clockOut.Time
	}
	return &a, nil
}
