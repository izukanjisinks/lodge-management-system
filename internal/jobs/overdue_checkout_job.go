package jobs

import (
	"log"
	"time"

	"lodge-system/internal/repository"
)

// OverdueCheckoutJob runs nightly and extends the check_out date (and invoice due_date)
// for any checked_in booking whose check_out has already passed.
// This covers guests who stayed beyond their original departure date.
type OverdueCheckoutJob struct {
	bookingRepo *repository.BookingRepository
	invoiceRepo *repository.InvoiceRepository
}

func NewOverdueCheckoutJob(
	bookingRepo *repository.BookingRepository,
	invoiceRepo *repository.InvoiceRepository,
) *OverdueCheckoutJob {
	return &OverdueCheckoutJob{bookingRepo: bookingRepo, invoiceRepo: invoiceRepo}
}

// Start launches the job in a background goroutine, firing once at the next midnight
// and then every 24 hours thereafter.
func (j *OverdueCheckoutJob) Start() {
	go func() {
		time.Sleep(durationUntilMidnight())
		j.run()

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			j.run()
		}
	}()
}

func (j *OverdueCheckoutJob) run() {
	log.Println("[overdue-checkout] running overdue checkout scan")

	rows, err := j.bookingRepo.FindOverdueCheckouts()
	if err != nil {
		log.Printf("[overdue-checkout] failed to query overdue bookings: %v", err)
		return
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	extended := 0

	for _, b := range rows {
		if err := j.bookingRepo.ExtendCheckout(b.ID, b.OrgID, today); err != nil {
			log.Printf("[overdue-checkout] failed to extend booking %s: %v", b.ID, err)
			continue
		}
		if err := j.bookingRepo.MarkOverstayed(b.ID, b.OrgID); err != nil {
			log.Printf("[overdue-checkout] failed to mark overstayed for booking %s: %v", b.ID, err)
		}
		// Keep invoice due date in sync
		if err := j.invoiceRepo.UpdateDueDate(b.ID, b.OrgID, today); err != nil {
			log.Printf("[overdue-checkout] failed to update invoice due_date for booking %s: %v", b.ID, err)
		}
		extended++
	}

	log.Printf("[overdue-checkout] extended %d overdue booking(s)", extended)
}

// durationUntilMidnight returns the duration from now until the next midnight UTC.
func durationUntilMidnight() time.Duration {
	now := time.Now().UTC()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return time.Until(next)
}
