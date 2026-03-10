package services

import (
	"errors"
	"fmt"
	"time"

	"hr-system/internal/interfaces"
	"hr-system/internal/models"
	"hr-system/internal/repository"

	"github.com/google/uuid"
)

type LeaveRequestService struct {
	repo            *repository.LeaveRequestRepository
	balanceService  *LeaveBalanceService
	leaveTypeRepo   *repository.LeaveTypeRepository
	holidayRepo     *repository.HolidayRepository
	empRepo         *repository.EmployeeRepository
	workflowService *WorkflowService
}

func NewLeaveRequestService(
	repo *repository.LeaveRequestRepository,
	balanceSvc *LeaveBalanceService,
	ltRepo *repository.LeaveTypeRepository,
	holidayRepo *repository.HolidayRepository,
	empRepo *repository.EmployeeRepository,
	workflowSvc *WorkflowService,
) *LeaveRequestService {
	return &LeaveRequestService{
		repo:            repo,
		balanceService:  balanceSvc,
		leaveTypeRepo:   ltRepo,
		holidayRepo:     holidayRepo,
		empRepo:         empRepo,
		workflowService: workflowSvc,
	}
}

func (s *LeaveRequestService) Create(req *models.LeaveRequest) error {
	// Validate dates
	today := time.Now().Truncate(24 * time.Hour)
	if req.StartDate.Before(today) {
		return errors.New("start_date must be today or in the future")
	}
	if req.EndDate.Before(req.StartDate) {
		return errors.New("end_date must be >= start_date")
	}

	// Verify employee exists
	if _, err := s.empRepo.GetByID(req.EmployeeID); err != nil {
		return errors.New("employee not found")
	}

	// Verify leave type
	lt, err := s.leaveTypeRepo.GetByID(req.LeaveTypeID)
	if err != nil {
		return errors.New("leave type not found")
	}
	if !lt.IsActive {
		return errors.New("leave type is inactive")
	}

	// Count business days
	holidays, err := s.holidayRepo.GetHolidaysInRange(req.StartDate, req.EndDate, "")
	if err != nil {
		return err
	}
	req.TotalDays = CountBusinessDays(req.StartDate, req.EndDate, holidays)
	if req.TotalDays == 0 {
		return errors.New("leave period contains no working days")
	}

	// Check balance
	year := req.StartDate.Year()
	ok, err := s.balanceService.HasSufficientBalance(req.EmployeeID, req.LeaveTypeID, year, req.TotalDays)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("insufficient leave balance")
	}

	// Check overlap
	overlap, err := s.repo.HasOverlap(req.EmployeeID, req.StartDate, req.EndDate, nil)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("overlapping leave request already exists")
	}

	// Warn if document required (don't block)
	if lt.RequiresDocument && req.AttachmentURL == "" {
		// just continue — caller may log/warn
	}

	// Create and increment pending balance
	if err := s.repo.Create(req); err != nil {
		return err
	}
	if err := s.balanceService.IncrementPending(req.EmployeeID, req.LeaveTypeID, year, req.TotalDays); err != nil {
		return err
	}

	// Trigger workflow for leave request approval
	if s.workflowService != nil {
		if err := s.initiateLeaveRequestWorkflow(req); err != nil {
			// Log the error but don't fail the leave request creation
			// The workflow can be manually triggered later if needed
			fmt.Printf("Warning: Failed to initiate workflow for leave request %s: %v\n", req.ID, err)
		}
	}

	return nil
}

// initiateLeaveRequestWorkflow creates a workflow instance for the leave request
func (s *LeaveRequestService) initiateLeaveRequestWorkflow(req *models.LeaveRequest) error {
	// Get employee details for the task
	employee, err := s.empRepo.GetByID(req.EmployeeID)
	if err != nil {
		return fmt.Errorf("failed to get employee details: %w", err)
	}

	// Extract position and department names
	positionName := ""
	if employee.Position != nil {
		positionName = employee.Position.Title
	}

	departmentName := ""
	if employee.Department != nil {
		departmentName = employee.Department.Name
	}

	// Create task details
	taskDetails := models.TaskDetails{
		TaskID:          req.ID.String(),
		TaskType:        "leave_request",
		TaskDescription: fmt.Sprintf("%s %s has requested %d days of leave from %s to %s",
			employee.FirstName, employee.LastName,
			req.TotalDays,
			req.StartDate.Format("2006-01-02"),
			req.EndDate.Format("2006-01-02")),
		SenderDetails: models.SenderDetails{
			SenderID:   employee.UserID.String(),
			SenderName: fmt.Sprintf("%s %s", employee.FirstName, employee.LastName),
			Position:   positionName,
			Department: departmentName,
		},
	}

	// Calculate due date (e.g., 3 days from now for approval)
	dueDate := time.Now().Add(72 * time.Hour)

	// Initiate the workflow using the LEAVE_REQUEST type constant
	_, err = s.workflowService.InitiateWorkflow(
		models.WorkflowTypeLeaveRequest,
		taskDetails,
		employee.UserID.String(),
		"medium",
		&dueDate,
	)

	return err
}

func (s *LeaveRequestService) GetByID(id uuid.UUID) (*models.LeaveRequest, error) {
	return s.repo.GetByID(id)
}

func (s *LeaveRequestService) List(filter interfaces.LeaveRequestFilter, page, pageSize int) ([]models.LeaveRequest, int, error) {
	return s.repo.List(filter, page, pageSize)
}

func (s *LeaveRequestService) Cancel(id, requestorID uuid.UUID) error {
	req, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("leave request not found")
	}
	if req.EmployeeID != requestorID {
		return errors.New("you can only cancel your own leave requests")
	}
	switch req.Status {
	case models.LeaveStatusPending:
		// always allowed
	case models.LeaveStatusApproved:
		if !req.StartDate.After(time.Now()) {
			return errors.New("cannot cancel an approved request that has already started")
		}
	default:
		return errors.New("only pending or approved requests can be cancelled")
	}

	year := req.StartDate.Year()
	if err := s.repo.UpdateStatus(id, models.LeaveStatusCancelled, nil, ""); err != nil {
		return err
	}
	return s.balanceService.DecrementPending(req.EmployeeID, req.LeaveTypeID, year, req.TotalDays)
}

func (s *LeaveRequestService) Approve(id, reviewerEmployeeID uuid.UUID, comment string) error {
	req, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("leave request not found")
	}
	if req.Status != models.LeaveStatusPending {
		return errors.New("only pending requests can be approved")
	}

	year := req.StartDate.Year()
	if err := s.repo.UpdateStatus(id, models.LeaveStatusApproved, &reviewerEmployeeID, comment); err != nil {
		return err
	}
	return s.balanceService.ApproveLeave(req.EmployeeID, req.LeaveTypeID, year, req.TotalDays)
}

func (s *LeaveRequestService) Reject(id, reviewerEmployeeID uuid.UUID, comment string) error {
	req, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("leave request not found")
	}
	if req.Status != models.LeaveStatusPending {
		return errors.New("only pending requests can be rejected")
	}

	year := req.StartDate.Year()
	if err := s.repo.UpdateStatus(id, models.LeaveStatusRejected, &reviewerEmployeeID, comment); err != nil {
		return err
	}
	return s.balanceService.DecrementPending(req.EmployeeID, req.LeaveTypeID, year, req.TotalDays)
}
