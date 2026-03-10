package services

import (
	"fmt"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"
	"hr-system/pkg/utils"

	"github.com/google/uuid"
)

type DashboardService struct {
	empRepo        *repository.EmployeeRepository
	posRepo        *repository.PositionRepository
	deptRepo       *repository.DepartmentRepository
	lbRepo         *repository.LeaveBalanceRepository
	lrRepo         *repository.LeaveRequestRepository
	adminDashRepo  *repository.AdminDashboardRepository
}

func NewDashboardService(
	empRepo *repository.EmployeeRepository,
	posRepo *repository.PositionRepository,
	deptRepo *repository.DepartmentRepository,
	lbRepo *repository.LeaveBalanceRepository,
	lrRepo *repository.LeaveRequestRepository,
	adminDashRepo *repository.AdminDashboardRepository,
) *DashboardService {
	return &DashboardService{
		empRepo:       empRepo,
		posRepo:       posRepo,
		deptRepo:      deptRepo,
		lbRepo:        lbRepo,
		lrRepo:        lrRepo,
		adminDashRepo: adminDashRepo,
	}
}

// GetEmployeeDashboard builds the EmployeeDashboardStats for the employee
// linked to the given userID.
func (s *DashboardService) GetEmployeeDashboard(userID uuid.UUID) (*models.EmployeeDashboardStats, error) {
	// 1. Find the employee record linked to this user.
	// emps, _, err := s.empRepo.List(interfaces.EmployeeFilter{}, 1, 10000)
	// if err != nil {
	// 	return nil, err
	// }
	emp, err := s.empRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// for i := range emps {
	// 	if emps[i].UserID != nil && *emps[i].UserID == userID {
	// 		emp = &emps[i]
	// 		break
	// 	}
	// }
	// if emp == nil {
	// 	return nil, nil // caller handles the not-found case
	// }

	now := time.Now()
	year := now.Year()

	// 2. Resolve position.
	posTitle := ""
	posCode := ""
	pos, err := s.posRepo.GetByID(emp.PositionID)
	if err == nil {
		posTitle = pos.Title
		posCode = pos.Code
	}

	// 3. Resolve department.
	deptName := ""
	dept, err := s.deptRepo.GetByID(emp.DepartmentID)
	if err == nil {
		deptName = dept.Name
	}

	// 4. Resolve supervisor name.
	supervisorName := "None"
	if emp.ManagerID != nil {
		mgr, err := s.empRepo.GetByID(*emp.ManagerID)
		if err == nil {
			supervisorName = mgr.FirstName + " " + mgr.LastName
		}
	}

	// 5. Employment period in months (hire date → today).
	employmentMonths := monthsBetween(emp.HireDate, now)

	// 6. Holidays this month (Zambian holidays only).
	zambianHolidays := utils.GetHolidaysForMonth(now.Month())
	thisMonthHolidays := make([]models.HolidayDetails, 0, len(zambianHolidays))
	for _, zh := range zambianHolidays {
		thisMonthHolidays = append(thisMonthHolidays, models.HolidayDetails{
			Name: zh.Name,
			Date: fmt.Sprintf("%d-%02d-%02d", now.Year(), zh.Month, zh.Day),
		})
	}

	// 7. Leave days earned this month = earned_leave_days for the AL balance
	//    divided by the number of months elapsed so far this year.
	leaveDaysThisMonth := 0
	alBalance, err := s.lbRepo.GetByEmployeeAndYear(emp.ID, year)
	yearlyEntitlement := 0
	if err == nil {
		for _, lb := range alBalance {
			if lb.LeaveType != nil && lb.LeaveType.Code == "AL" {
				// Yearly entitlement = base days from leave type (e.g., 24 for AL)
				yearlyEntitlement = lb.TotalEntitled
				// earned_leave_days accumulates +2 per month; divide by months elapsed
				elapsed := int(now.Month())
				if elapsed > 0 {
					leaveDaysThisMonth = lb.EarnedLeaveDays / elapsed
				}
				break
			}
		}
	}

	// 8. Total leave requests for this employee (all statuses).
	_, leaveRequestCount, err := s.lrRepo.List(
		interfaces.LeaveRequestFilter{EmployeeID: &emp.ID}, 1, 1,
	)
	if err != nil {
		leaveRequestCount = 0
	}

	return &models.EmployeeDashboardStats{
		Employee: models.EmployeDetails{
			EmployeeName:     emp.FirstName + " " + emp.LastName,
			Address:          emp.Address,
			Role:             string(emp.EmploymentType),
			Position:         posTitle,
			EmploymentPeriod: employmentMonths,
			Department:       deptName,
			Supervisor:       supervisorName,
			PositionCode:     posCode,
		},
		Holidays: models.Holidays{
			Total:   len(thisMonthHolidays),
			Details: thisMonthHolidays,
		},
		LeaveDaysThisMonth: leaveDaysThisMonth,
		YearlyEntitlement:  yearlyEntitlement,
		LeaveRequests:      leaveRequestCount,
	}, nil
}

// GetAdminDashboard returns aggregate stats for the admin/HR dashboard.
func (s *DashboardService) GetAdminDashboard(from, to *time.Time) (*models.AdminDashboard, error) {
	return s.adminDashRepo.GetStats(from, to)
}

// monthsBetween returns the number of whole months between from and to.
func monthsBetween(from, to time.Time) int {
	months := (to.Year()-from.Year())*12 + int(to.Month()) - int(from.Month())
	if months < 0 {
		return 0
	}
	return months
}
