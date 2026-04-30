package jobs

import (
	"encoding/json"
	"log"
	"time"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"
)

// OverdueCheckoutJob runs nightly and extends the check_out date (and invoice due_date)
// for any checked_in booking whose check_out has already passed.
type OverdueCheckoutJob struct {
	bookingRepo  *repository.BookingRepository
	invoiceRepo  *repository.InvoiceRepository
	auditLogRepo *repository.AuditLogRepository
	settingsRepo *repository.OrganizationSettingsRepository
}

func NewOverdueCheckoutJob(
	bookingRepo *repository.BookingRepository,
	invoiceRepo *repository.InvoiceRepository,
	auditLogRepo *repository.AuditLogRepository,
	settingsRepo *repository.OrganizationSettingsRepository,
) *OverdueCheckoutJob {
	return &OverdueCheckoutJob{
		bookingRepo:  bookingRepo,
		invoiceRepo:  invoiceRepo,
		auditLogRepo: auditLogRepo,
		settingsRepo: settingsRepo,
	}
}

// Start launches the job in a background goroutine, firing once at the next midnight UTC
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

	enabledOrgIDs, err := j.settingsRepo.ListEnabledOrgsForJob("auto_extend_checkout")
	if err != nil {
		log.Printf("[overdue-checkout] failed to fetch enabled orgs: %v", err)
		return
	}
	if len(enabledOrgIDs) == 0 {
		log.Println("[overdue-checkout] no orgs have auto_extend_checkout enabled")
		return
	}

	refs, err := j.bookingRepo.FindOverdueCheckouts(enabledOrgIDs)
	if err != nil {
		log.Printf("[overdue-checkout] failed to query overdue bookings: %v", err)
		return
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	extended := 0

	for _, b := range refs {
		if err := j.bookingRepo.ExtendCheckout(b.ID, b.OrgID, today); err != nil {
			log.Printf("[overdue-checkout] failed to extend booking %s: %v", b.ID, err)
			continue
		}

		if err := j.bookingRepo.MarkOverstayed(b.ID, b.OrgID); err != nil {
			log.Printf("[overdue-checkout] failed to mark overstayed for booking %s: %v", b.ID, err)
		}

		if err := j.invoiceRepo.UpdateDueDate(b.ID, b.OrgID, today); err != nil {
			log.Printf("[overdue-checkout] failed to update invoice due_date for booking %s: %v", b.ID, err)
		}

		j.writeAuditLog(b, today)
		extended++
	}

	log.Printf("[overdue-checkout] extended %d overdue booking(s)", extended)
}

func (j *OverdueCheckoutJob) writeAuditLog(b repository.OverdueBookingRef, extendedTo time.Time) {
	payload := models.OverstayedPayload{
		BookingNumber:    b.BookingNumber,
		RoomName:         b.RoomName,
		ClientName:       b.ClientName,
		OriginalCheckOut: b.OriginalCheckOut.Format("2006-01-02"),
		ExtendedTo:       extendedTo.Format("2006-01-02"),
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[overdue-checkout] failed to marshal audit payload for booking %s: %v", b.ID, err)
		return
	}

	entry := &models.AuditLog{
		OrgID:      b.OrgID,
		ActorType:  models.AuditActorSystem,
		ActorName:  "overdue-checkout-job",
		Action:     models.AuditActionBookingOverstayed,
		EntityType: models.AuditEntityBooking,
		EntityID:   b.ID,
		Payload:    raw,
	}
	if err := j.auditLogRepo.Insert(entry); err != nil {
		log.Printf("[overdue-checkout] failed to write audit log for booking %s: %v", b.ID, err)
	}
}

// durationUntilMidnight returns the duration from now until the next midnight UTC.
func durationUntilMidnight() time.Duration {
	now := time.Now().UTC()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return time.Until(next)
}
