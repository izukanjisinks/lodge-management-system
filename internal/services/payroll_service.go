package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"
	"hr-system/internal/utils/email"

	"github.com/google/uuid"
)

type PayrollService struct {
	repo           *repository.PayrollRepository
	payslipService *PayslipService
	empRepo        *repository.EmployeeRepository
	emailService   *email.EmailService
}

func NewPayrollService(
	repo *repository.PayrollRepository,
	payslipService *PayslipService,
	empRepo *repository.EmployeeRepository,
	emailService *email.EmailService,
) *PayrollService {
	return &PayrollService{
		repo:           repo,
		payslipService: payslipService,
		empRepo:        empRepo,
		emailService:   emailService,
	}
}

// Create opens a new payroll period
func (s *PayrollService) Create(startDate, endDate time.Time) (*models.Payroll, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	payroll := &models.Payroll{
		StartDate: startDate,
		EndDate:   endDate,
		Status:    models.PayrollStatusOpen,
	}

	if err := s.repo.Create(payroll); err != nil {
		return nil, fmt.Errorf("failed to create payroll: %w", err)
	}

	return s.repo.GetByID(payroll.ID)
}

// Process validates the payroll, marks it as PROCESSING, and kicks off payslip
// generation for all active employees in the background. The caller gets an
// immediate response — poll GET /payrolls/{id} to check completion.
func (s *PayrollService) Process(payrollID uuid.UUID, processedBy uuid.UUID) (*models.Payroll, error) {
	payroll, err := s.repo.GetByID(payrollID)
	if err != nil {
		return nil, errors.New("payroll not found")
	}

	if payroll.Status != models.PayrollStatusOpen {
		return nil, fmt.Errorf("payroll is already %s", payroll.Status)
	}

	// Mark as processing
	payroll.Status = models.PayrollStatusProcessing
	if err := s.repo.Update(payroll); err != nil {
		return nil, err
	}

	// Run payslip generation in the background
	go s.processPayslips(payrollID, processedBy)

	// Return the payroll immediately with PROCESSING status
	return s.repo.GetByID(payrollID)
}

// processPayslips generates payslips for all active employees and marks the payroll as completed.
func (s *PayrollService) processPayslips(payrollID uuid.UUID, processedBy uuid.UUID) {
	payroll, err := s.repo.GetByID(payrollID)
	if err != nil {
		log.Printf("payroll processing error: failed to fetch payroll %s: %v", payrollID, err)
		return
	}

	month := int(payroll.EndDate.Month())
	year := payroll.EndDate.Year()

	filter := interfaces.EmployeeFilter{EmploymentStatus: "active"}
	employees, _, err := s.empRepo.List(filter, 1, 10000)
	if err != nil {
		log.Printf("payroll processing error: failed to fetch employees: %v", err)
		s.markFailed(payroll)
		return
	}

	period := fmt.Sprintf("%s %d", time.Month(month), year)

	generated := 0
	for _, emp := range employees {
		_, err := s.payslipService.Generate(emp.ID, month, year)
		if err != nil {
			log.Printf("payroll %s: skipped employee %s (%s): %s", payrollID, emp.ID, emp.FullName(), err.Error())
			continue
		}
		generated++

		// Send payslip notification email
		if s.emailService != nil {
			empEmail := emp.Email
			empName := emp.FirstName
			go func() {
				htmlBody := email.PayslipReadyTemplate(empName, period)
				subject := fmt.Sprintf("Your Payslip for %s is Ready", period)
				if err := s.emailService.SendEmail([]string{empEmail}, subject, htmlBody); err != nil {
					log.Printf("payroll %s: failed to send payslip email to %s: %v", payrollID, empEmail, err)
				}
			}()
		}
	}

	now := time.Now()
	payroll.Status = models.PayrollStatusCompleted
	payroll.ProcessedBy = &processedBy
	payroll.ProcessedAt = &now
	if err := s.repo.Update(payroll); err != nil {
		log.Printf("payroll processing error: failed to mark payroll %s as completed: %v", payrollID, err)
		return
	}

	log.Printf("payroll %s processed: %d payslips generated", payrollID, generated)
}

// markFailed reverts a payroll back to OPEN if background processing fails before generating any payslips.
func (s *PayrollService) markFailed(payroll *models.Payroll) {
	payroll.Status = models.PayrollStatusOpen
	if err := s.repo.Update(payroll); err != nil {
		log.Printf("payroll processing error: failed to revert payroll %s to OPEN: %v", payroll.ID, err)
	}
}

func (s *PayrollService) GetByID(id uuid.UUID) (*models.Payroll, error) {
	payroll, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Derive month/year and load payslip summary
	month := int(payroll.EndDate.Month())
	year := payroll.EndDate.Year()
	m := month
	y := year
	payslips, total, _ := s.payslipService.List(nil, &m, &y, 1, 10000)
	payroll.Payslips = payslips
	payroll.EmployeeCount = total
	for _, p := range payslips {
		payroll.TotalNetSalary += p.NetSalary
	}

	return payroll, nil
}

func (s *PayrollService) List(status string, page, pageSize int) ([]models.Payroll, int, error) {
	return s.repo.List(status, page, pageSize)
}

func (s *PayrollService) Cancel(id uuid.UUID) error {
	payroll, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("payroll not found")
	}
	if payroll.Status != models.PayrollStatusOpen {
		return fmt.Errorf("can only cancel payrolls in OPEN status, current: %s", payroll.Status)
	}
	payroll.Status = models.PayrollStatusCancelled
	return s.repo.Update(payroll)
}

func (s *PayrollService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

