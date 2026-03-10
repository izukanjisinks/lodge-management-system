package jobs

import (
	"log"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

// MonthlyLeaveAccrualJob credits 2 days of Annual Leave to every active employee
// on the 1st of each month. It runs as a background goroutine.
//
// Accrual logic:
//   - Only employees with employment_status = "active" are credited.
//   - 2 days are added to the earned_leave_days column of the leave_balances row.
//   - If an employee has no balance record for the current year yet, one is
//     initialised first (with zero entitlement) before the credit is applied.
//   - The Annual Leave type is identified by its code "AL".
type MonthlyLeaveAccrualJob struct {
	empRepo        *repository.EmployeeRepository
	lbRepo         *repository.LeaveBalanceRepository
	ltRepo         *repository.LeaveTypeRepository
	daysPerAccrual int
}

func NewMonthlyLeaveAccrualJob(
	empRepo *repository.EmployeeRepository,
	lbRepo *repository.LeaveBalanceRepository,
	ltRepo *repository.LeaveTypeRepository,
) *MonthlyLeaveAccrualJob {
	return &MonthlyLeaveAccrualJob{
		empRepo:        empRepo,
		lbRepo:         lbRepo,
		ltRepo:         ltRepo,
		daysPerAccrual: 2,
	}
}

// Start launches the accrual loop as a background goroutine. Call once from
// main after the database connection is established.
func (j *MonthlyLeaveAccrualJob) Start() {
	go j.loop()
}

func (j *MonthlyLeaveAccrualJob) loop() {
	// Block until the 1st of the next month at midnight UTC. Yearly loop
	waitUntilNextRun()

	log.Println("[LeaveAccrual] Running monthly leave accrual")
	j.run()

	// After the first aligned run, tick every 30 days.
	// For exact cron-style scheduling consider robfig/cron — this approach
	// has zero external dependencies and stays close to the 1st of each month.
	ticker := time.NewTicker(30 * 24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		log.Println("[LeaveAccrual] Running monthly leave accrual")
		j.run()
	}
}

// run executes one accrual cycle.
func (j *MonthlyLeaveAccrualJob) run() {
	year := time.Now().Year()

	// 1. Resolve the Annual Leave type by its well-known code.
	lt, err := j.ltRepo.GetByCode("AL")
	if err != nil {
		log.Printf("[LeaveAccrual] ERROR: could not find leave type AL: %v", err)
		return
	}

	// 2. Fetch all active employees (page size 10 000; increase if needed).
	emps, _, err := j.empRepo.List(interfaces.EmployeeFilter{EmploymentStatus: "active"}, 1, 10000)
	if err != nil {
		log.Printf("[LeaveAccrual] ERROR: could not list employees: %v", err)
		return
	}
	if len(emps) == 0 {
		log.Println("[LeaveAccrual] No active employees found, skipping")
		return
	}

	credited := 0
	skipped := 0

	for _, emp := range emps {
		// 3. Ensure a balance row exists for this employee / leave type / year.
		lb, err := j.lbRepo.GetByEmployeeTypeYear(emp.ID, lt.ID, year)
		if err != nil {
			// Row is missing — insert a blank one so we can credit it.
			blank := &models.LeaveBalance{
				ID:          uuid.New(),
				EmployeeID:  emp.ID,
				LeaveTypeID: lt.ID,
				Year:        year,
			}
			if upsertErr := j.lbRepo.Upsert(blank); upsertErr != nil {
				log.Printf("[LeaveAccrual] WARN: could not initialise balance for %s: %v", emp.EmployeeNumber, upsertErr)
				skipped++
				continue
			}
			lb, err = j.lbRepo.GetByEmployeeTypeYear(emp.ID, lt.ID, year)
			if err != nil {
				log.Printf("[LeaveAccrual] WARN: re-fetch failed for %s: %v", emp.EmployeeNumber, err)
				skipped++
				continue
			}
		}

		// 4. Credit 2 days via the earned_leave_days column.
		if err := j.lbRepo.Adjust(lb.ID, j.daysPerAccrual); err != nil {
			log.Printf("[LeaveAccrual] WARN: failed to credit %s: %v", emp.EmployeeNumber, err)
			skipped++
			continue
		}
		credited++
	}

	log.Printf("[LeaveAccrual] Done — credited %d employees, skipped %d", credited, skipped)
}

// waitUntilNextRun blocks until midnight on the 1st of the next month (UTC).
func waitUntilNextRun() {
	now := time.Now().UTC()
	firstOfNext := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	wait := time.Until(firstOfNext)
	log.Printf("[LeaveAccrual] Next run scheduled in %s (on %s UTC)", wait.Round(time.Minute), firstOfNext.Format("2006-01-02"))
	time.Sleep(wait)
}
