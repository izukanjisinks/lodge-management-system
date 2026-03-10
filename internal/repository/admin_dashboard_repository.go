package repository

import (
	"database/sql"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"
)

type AdminDashboardRepository struct {
	db *sql.DB
}

func NewAdminDashboardRepository() *AdminDashboardRepository {
	return &AdminDashboardRepository{db: database.DB}
}

func (r *AdminDashboardRepository) GetStats(from, to *time.Time) (*models.AdminDashboard, error) {
	dash := &models.AdminDashboard{}

	// Default date range: last 12 months
	now := time.Now()
	defaultFrom := now.AddDate(-1, 0, 0)
	defaultTo := now

	if from == nil {
		from = &defaultFrom
	}
	if to == nil {
		to = &defaultTo
	}

	fromMonth := int(from.Month())
	fromYear := from.Year()
	toMonth := int(to.Month())
	toYear := to.Year()

	// Total active employees
	r.db.QueryRow(`SELECT COUNT(*) FROM employees WHERE employment_status = 'active'`).Scan(&dash.TotalEmployees)

	// Total active departments
	r.db.QueryRow(`SELECT COUNT(*) FROM departments WHERE is_active = true`).Scan(&dash.TotalDepartments)

	// Active payrolls (OPEN or PROCESSING)
	r.db.QueryRow(`SELECT COUNT(*) FROM payrolls WHERE status IN ('OPEN', 'PROCESSING')`).Scan(&dash.ActivePayrolls)

	// Leave requests overview (within date range)
	r.db.QueryRow(`SELECT COUNT(*) FROM leave_requests WHERE status = 'pending' AND created_at >= $1 AND created_at <= $2`, from, to).Scan(&dash.LeaveRequests.PendingRequests)
	r.db.QueryRow(`SELECT COUNT(*) FROM leave_requests WHERE status = 'approved' AND created_at >= $1 AND created_at <= $2`, from, to).Scan(&dash.LeaveRequests.ApprovedRequests)
	r.db.QueryRow(`SELECT COUNT(*) FROM leave_requests WHERE status = 'rejected' AND created_at >= $1 AND created_at <= $2`, from, to).Scan(&dash.LeaveRequests.RejectedRequests)

	// Recent hires (within date range, last 5)
	rows, err := r.db.Query(`
		SELECT e.first_name, e.last_name, COALESCE(p.title, '') AS position, e.hire_date::text
		FROM employees e
		LEFT JOIN positions p ON e.position_id = p.id
		WHERE e.employment_status = 'active' AND e.hire_date >= $1 AND e.hire_date <= $2
		ORDER BY e.hire_date DESC
		LIMIT 5`, from, to)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var h models.RecentHire
			if err := rows.Scan(&h.FirstName, &h.LastName, &h.Position, &h.HireDate); err != nil {
				continue
			}
			dash.RecentHires = append(dash.RecentHires, h)
		}
	}
	if dash.RecentHires == nil {
		dash.RecentHires = []models.RecentHire{}
	}

	// Monthly payroll cost (within date range)
	payrollRows, err := r.db.Query(`
		SELECT TO_CHAR(TO_DATE(month::text || '-' || year::text, 'MM-YYYY'), 'Month') AS month_name,
		       year, COALESCE(SUM(net_salary), 0) AS total_net_salary
		FROM payslips
		WHERE (year > $1 OR (year = $1 AND month >= $2))
		  AND (year < $3 OR (year = $3 AND month <= $4))
		GROUP BY year, month, month_name
		ORDER BY year, month`,
		fromYear, fromMonth, toYear, toMonth)
	if err == nil {
		defer payrollRows.Close()
		for payrollRows.Next() {
			var m models.MonthlyPayrollCost
			if err := payrollRows.Scan(&m.Month, &m.Year, &m.TotalNetSalary); err != nil {
				continue
			}
			dash.MonthlyPayrollCost = append(dash.MonthlyPayrollCost, m)
		}
	}
	if dash.MonthlyPayrollCost == nil {
		dash.MonthlyPayrollCost = []models.MonthlyPayrollCost{}
	}

	// Hiring trend (within date range)
	hiringRows, err := r.db.Query(`
		SELECT TO_CHAR(hire_date, 'Month') AS month_name,
		       EXTRACT(YEAR FROM hire_date)::int AS year,
		       COUNT(*) AS new_hires
		FROM employees
		WHERE hire_date >= $1 AND hire_date <= $2
		GROUP BY EXTRACT(YEAR FROM hire_date), EXTRACT(MONTH FROM hire_date), TO_CHAR(hire_date, 'Month')
		ORDER BY EXTRACT(YEAR FROM hire_date), EXTRACT(MONTH FROM hire_date)`,
		from, to)
	if err == nil {
		defer hiringRows.Close()
		for hiringRows.Next() {
			var h models.HiringTrend
			if err := hiringRows.Scan(&h.Month, &h.Year, &h.NewHires); err != nil {
				continue
			}
			dash.HiringTrend = append(dash.HiringTrend, h)
		}
	}
	if dash.HiringTrend == nil {
		dash.HiringTrend = []models.HiringTrend{}
	}

	return dash, nil
}
