package models

type AdminDashboard struct {
	TotalEmployees     int                   `json:"total_employees"`
	TotalDepartments   int                   `json:"total_departments"`
	ActivePayrolls     int                   `json:"active_payrolls"`
	RecentHires        []RecentHire          `json:"recent_hires"`
	LeaveRequests      LeaveRequestsOverview `json:"leave_requests"`
	MonthlyPayrollCost []MonthlyPayrollCost  `json:"monthly_payroll_cost"`
	HiringTrend        []HiringTrend         `json:"hiring_trend"`
}

type MonthlyPayrollCost struct {
	Month          string  `json:"month"`
	Year           int     `json:"year"`
	TotalNetSalary float64 `json:"total_net_salary"`
}

type HiringTrend struct {
	Month    string `json:"month"`
	Year     int    `json:"year"`
	NewHires int    `json:"new_hires"`
}

type RecentHire struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Position  string `json:"position"`
	HireDate  string `json:"hire_date"`
}

type LeaveRequestsOverview struct {
	PendingRequests  int `json:"pending_requests"`
	ApprovedRequests int `json:"approved_requests"`
	RejectedRequests int `json:"rejected_requests"`
}
