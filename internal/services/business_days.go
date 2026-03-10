package services

import (
	"time"
)

// CountBusinessDays counts working days between start and end (inclusive),
// excluding weekends and the provided holiday dates (map of "YYYY-MM-DD" → true).
func CountBusinessDays(start, end time.Time, holidays map[string]bool) int {
	if end.Before(start) {
		return 0
	}
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		wd := d.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			continue
		}
		if holidays[d.Format("2006-01-02")] {
			continue
		}
		count++
	}
	return count
}

// IsWeekend returns true for Saturday and Sunday.
func IsWeekend(d time.Time) bool {
	wd := d.Weekday()
	return wd == time.Saturday || wd == time.Sunday
}

// ProrateEntitlement computes the prorated entitled days for an employee hired mid-year.
// It returns: floor(defaultDays * remainingMonths / 12)
func ProrateEntitlement(defaultDays int, hireDate time.Time, year int) int {
	yearEnd := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	if hireDate.After(yearEnd) {
		return 0
	}
	// Remaining months including the hire month
	remainingMonths := int(yearEnd.Month()) - int(hireDate.Month()) + 1
	if remainingMonths <= 0 {
		return 0
	}
	return defaultDays * remainingMonths / 12
}
