package jobs

import (
	"log"
	"time"

	"hr-system/internal/models"
	"hr-system/internal/repository"
)

// YearEndCarryForwardJob runs once a year on 31 December at 23:00 UTC.
// For every employee balance it:
//  1. Calculates the remaining days (balance) for the ending year.
//  2. If the leave type allows carry-forward, caps the remainder at
//     max_carry_forward_days and writes that value as carried_forward on
//     the new year's balance row (creating it if it doesn't exist yet).
//  3. Leave types with is_carry_forward_allowed=false are skipped (0 carries over).
type YearEndCarryForwardJob struct {
	lbRepo *repository.LeaveBalanceRepository
	ltRepo *repository.LeaveTypeRepository
}

func NewYearEndCarryForwardJob(
	lbRepo *repository.LeaveBalanceRepository,
	ltRepo *repository.LeaveTypeRepository,
) *YearEndCarryForwardJob {
	return &YearEndCarryForwardJob{lbRepo: lbRepo, ltRepo: ltRepo}
}

// Start launches the job as a background goroutine.
func (j *YearEndCarryForwardJob) Start() {
	go j.loop()
}

func (j *YearEndCarryForwardJob) loop() {
	waitUntilYearEnd()
	log.Println("[CarryForward] Running year-end carry-forward")
	j.run()

	// After the first run tick every 365 days to stay aligned annually.
	ticker := time.NewTicker(365 * 24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		log.Println("[CarryForward] Running year-end carry-forward")
		j.run()
	}
}

func (j *YearEndCarryForwardJob) run() {
	endingYear := time.Now().Year()
	newYear := endingYear + 1

	// Fetch all leave type carry-forward rules keyed by leave type ID.
	leaveTypes, err := j.ltRepo.List(false)
	if err != nil {
		log.Printf("[CarryForward] ERROR: could not list leave types: %v", err)
		return
	}
	ltRules := make(map[string]models.LeaveType, len(leaveTypes))
	for _, lt := range leaveTypes {
		ltRules[lt.ID.String()] = lt
	}

	// Fetch every balance row for the ending year.
	balances, err := j.lbRepo.GetAllByYear(endingYear)
	if err != nil {
		log.Printf("[CarryForward] ERROR: could not fetch balances for %d: %v", endingYear, err)
		return
	}
	if len(balances) == 0 {
		log.Printf("[CarryForward] No balances found for %d, skipping", endingYear)
		return
	}

	carried := 0
	skipped := 0

	for _, lb := range balances {
		lt, ok := ltRules[lb.LeaveTypeID.String()]
		if !ok || !lt.IsCarryForwardAllowed {
			skipped++
			continue
		}

		// Remaining days = balance (already computed by the repository scan).
		remaining := lb.Balance
		if remaining <= 0 {
			skipped++
			continue
		}

		// Cap at max_carry_forward_days.
		carryDays := remaining
		if lt.MaxCarryForwardDays > 0 && carryDays > lt.MaxCarryForwardDays {
			carryDays = lt.MaxCarryForwardDays
		}

		// Ensure a balance row exists for the new year.
		newLB, err := j.lbRepo.GetByEmployeeTypeYear(lb.EmployeeID, lb.LeaveTypeID, newYear)
		if err != nil {
			// Row doesn't exist yet — create a blank one.
			blank := &models.LeaveBalance{
				EmployeeID:  lb.EmployeeID,
				LeaveTypeID: lb.LeaveTypeID,
				Year:        newYear,
			}
			if upsertErr := j.lbRepo.Upsert(blank); upsertErr != nil {
				log.Printf("[CarryForward] WARN: could not create new-year row for employee %s / leave type %s: %v",
					lb.EmployeeID, lb.LeaveTypeID, upsertErr)
				skipped++
				continue
			}
			newLB, err = j.lbRepo.GetByEmployeeTypeYear(lb.EmployeeID, lb.LeaveTypeID, newYear)
			if err != nil {
				log.Printf("[CarryForward] WARN: re-fetch failed for employee %s: %v", lb.EmployeeID, err)
				skipped++
				continue
			}
		}

		// Write the carry-forward days onto the new year's row.
		if err := j.lbRepo.SetCarriedForward(newLB.ID, carryDays); err != nil {
			log.Printf("[CarryForward] WARN: could not set carried_forward for employee %s: %v", lb.EmployeeID, err)
			skipped++
			continue
		}
		carried++
	}

	log.Printf("[CarryForward] Done — carried forward %d balances, skipped %d", carried, skipped)
}

// waitUntilYearEnd blocks until 31 December at 23:00 UTC.
func waitUntilYearEnd() {
	now := time.Now().UTC()
	yearEnd := time.Date(now.Year(), 12, 31, 23, 0, 0, 0, time.UTC)
	// If we're already past this year's run time, target next year's.
	if now.After(yearEnd) {
		yearEnd = time.Date(now.Year()+1, 12, 31, 23, 0, 0, 0, time.UTC)
	}
	wait := time.Until(yearEnd)
	log.Printf("[CarryForward] Next run scheduled in %s (on %s UTC)", wait.Round(time.Minute), yearEnd.Format("2006-01-02 15:04"))
	time.Sleep(wait)
}
