package models

// type DashboardStats struct {
// 	TotalEmployees   int `json:"total_employees"`
// 	EmployeesOnLeave int `json:"employees_on_leave"`
// }

type EmployeeDashboardStats struct {
	Employee EmployeDetails `json:"employee_details"`

	Holidays Holidays `json:"holidays_this_month"`

	LeaveDaysThisMonth int `json:"leave_days_this_month"` //leave days earned this month
	YearlyEntitlement  int `json:"yearly_entitlement"`    //total entitled for the year, excluding carried forward and earned leave days
	LeaveRequests      int `json:"leave_requests"`
}

type EmployeDetails struct {
	EmployeeName     string `json:"employee_name"`
	Address          string `json:"address"`
	Role             string `json:"role"`
	Position         string `json:"position"`
	EmploymentPeriod int    `json:"employment_period"` //hire date to current date in months
	Department       string `json:"department"`
	Supervisor       string `json:"supervisor"` //supervisor name or "None" if no supervisor comes from manager id from employee table
	PositionCode     string `json:"position_code"`
}

type Holidays struct {
	Total   int              `json:"total"`
	Details []HolidayDetails `json:"details"`
}

type HolidayDetails struct {
	Name string `json:"name"`
	Date string `json:"date"` // Format: "2026-01-01"
}
