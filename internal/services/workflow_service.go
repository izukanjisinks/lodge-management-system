package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"hr-system/internal/models"
	"hr-system/internal/repository"
	"hr-system/internal/utils/email"
)

type WorkflowService struct {
	workflowRepo        *repository.WorkflowRepository
	instanceRepo        *repository.WorkflowInstanceRepository
	taskRepo            *repository.AssignedTaskRepository
	historyRepo         *repository.WorkflowHistoryRepository
	userRepo            *repository.UserRepository
	employeeRepo        *repository.EmployeeRepository
	leaveRequestRepo    *repository.LeaveRequestRepository
	leaveBalanceService *LeaveBalanceService
	emailService        *email.EmailService
}

func NewWorkflowService(
	workflowRepo *repository.WorkflowRepository,
	instanceRepo *repository.WorkflowInstanceRepository,
	taskRepo *repository.AssignedTaskRepository,
	historyRepo *repository.WorkflowHistoryRepository,
	userRepo *repository.UserRepository,
	employeeRepo *repository.EmployeeRepository,
	leaveRequestRepo *repository.LeaveRequestRepository,
	leaveBalanceService *LeaveBalanceService,
	emailService *email.EmailService,
) *WorkflowService {
	return &WorkflowService{
		workflowRepo:        workflowRepo,
		instanceRepo:        instanceRepo,
		taskRepo:            taskRepo,
		historyRepo:         historyRepo,
		userRepo:            userRepo,
		employeeRepo:        employeeRepo,
		leaveRequestRepo:    leaveRequestRepo,
		leaveBalanceService: leaveBalanceService,
		emailService:        emailService,
	}
}

// InitiateWorkflow starts a new workflow instance using workflow type
func (s *WorkflowService) InitiateWorkflow(
	workflowType models.WorkflowType,
	taskDetails models.TaskDetails,
	initiatorID string,
	priority string,
	dueDate *time.Time,
) (*models.WorkflowInstance, error) {
	// Get workflow template by type
	workflow, err := s.workflowRepo.GetByType(workflowType)
	if err != nil {
		return nil, fmt.Errorf("workflow not found for type %s: %w", workflowType, err)
	}

	// Get the first action step (the step after submission, e.g., "Pending Review")
	// This also returns the initial step ID for history tracking
	// This is optimized to do in one query instead of: get initial -> get transitions -> find submit -> get next step
	firstActionStep, initialStepID, err := s.workflowRepo.GetFirstActionStep(workflow.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get first action step: %w", err)
	}

	// Create workflow instance at the first action step (not the initial step)
	// The instance starts at "Pending Review" or whatever the first action step is
	instance := &models.WorkflowInstance{
		WorkflowID:    workflow.ID,
		CurrentStepID: firstActionStep.ID,
		Status:        "in_progress",
		TaskDetails:   taskDetails,
		CreatedBy:     initiatorID,
		Priority:      priority,
		DueDate:       dueDate,
	}

	if err := s.instanceRepo.Create(instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	// Determine who to assign the task to (for the first action step)
	assigneeID, err := s.determineAssignee(firstActionStep, taskDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to determine assignee: %w", err)
	}

	// Create assigned task for the first action step
	task := &models.AssignedTask{
		InstanceID: instance.ID,
		StepID:     firstActionStep.ID,
		StepName:   firstActionStep.StepName,
		AssignedTo: assigneeID,
		AssignedBy: initiatorID,
		Status:     "pending",
		DueDate:    dueDate,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Send email notification to assignee
	assigneeUUID, _ := uuid.Parse(assigneeID)
	if assignee, err := s.userRepo.GetUserByID(assigneeUUID); err == nil {
		task.TaskDetails = &taskDetails // Add task details for email
		go s.notifyTaskAssignment(task, assignee, instance)
	}

	// Get initiator name for history
	initiatorUUID, err := uuid.Parse(initiatorID)
	if err != nil {
		return nil, fmt.Errorf("invalid initiator ID: %w", err)
	}

	initiator, err := s.employeeRepo.GetByUserID(initiatorUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get initiator: %w", err)
	}

	// Create history entry (from initial step -> first action step via "submit")
	history := &models.WorkflowHistory{
		InstanceID:      instance.ID,
		FromStepID:      &initialStepID,
		ToStepID:        firstActionStep.ID,
		ActionTaken:     "submit",
		PerformedBy:     initiatorID,
		PerformedByName: initiator.FirstName + " " + initiator.LastName,
		Comments:        fmt.Sprintf("Initiated %s workflow", workflow.Name),
	}

	if err := s.historyRepo.Create(history); err != nil {
		return nil, fmt.Errorf("failed to create history: %w", err)
	}

	return instance, nil
}

// ProcessAction processes an action on a workflow instance
func (s *WorkflowService) ProcessAction(
	instanceID string,
	action string,
	performedByID string,
	comments string,
) error {
	// Get the workflow instance
	instance, err := s.instanceRepo.GetByID(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	// Check if instance is still in progress
	if instance.Status == "completed" || instance.Status == "cancelled" {
		return errors.New("workflow instance is already closed")
	}

	// Get current step
	currentStep, err := s.workflowRepo.GetStepByID(instance.CurrentStepID)
	if err != nil {
		return fmt.Errorf("current step not found: %w", err)
	}

	// Check if user has permission to perform this action
	if err := s.checkPermission(performedByID, currentStep); err != nil {
		return err
	}

	// Check if this is a rejection action
	isRejection := action == "reject" || action == "rejected" || action == "deny"

	// Check if we're performing an action FROM the final step (completing the workflow)
	// Final steps don't need transitions - any action from a final step completes the workflow
	isCompletingFinalStep := currentStep.Final && !isRejection

	var nextStep *models.WorkflowStep

	if isCompletingFinalStep {
		// For final steps, we don't need a transition - the next step is the same step
		// The workflow will be marked as completed below
		nextStep = currentStep
	} else {
		// Get valid transition for non-final steps
		transition, err := s.workflowRepo.GetTransitionByAction(currentStep.ID, action)
		if err == sql.ErrNoRows {
			return fmt.Errorf("action '%s' is not valid from current step '%s'", action, currentStep.StepName)
		}
		if err != nil {
			return fmt.Errorf("failed to get transition: %w", err)
		}

		// Get next step
		nextStep, err = s.workflowRepo.GetStepByID(transition.ToStepID)
		if err != nil {
			return fmt.Errorf("next step not found: %w", err)
		}
	}

	// Determine the new status
	// - If rejected: mark as "rejected"
	// - If completing from final step: mark as "completed"
	// - Otherwise: keep as "in_progress"
	newStatus := "in_progress"
	if isRejection {
		newStatus = "rejected"
	} else if isCompletingFinalStep {
		newStatus = "completed"
	}

	if err := s.instanceRepo.UpdateStep(instanceID, nextStep.ID, newStatus); err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	// Complete current task
	activeTask, err := s.taskRepo.GetActiveTaskForInstance(instanceID)
	if err == nil && activeTask != nil {
		if err := s.taskRepo.Complete(activeTask.ID); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}
	}

	// If rejected or completing final step, mark instance as completed and stop workflow
	if isRejection || isCompletingFinalStep {
		if err := s.instanceRepo.Complete(instanceID); err != nil {
			return fmt.Errorf("failed to complete instance: %w", err)
		}
	} else {
		// Create new task for the next step (including when moving TO final steps)
		// The final step still needs someone assigned to perform the final action
		assigneeID, err := s.determineAssignee(nextStep, instance.TaskDetails)
		if err != nil {
			return fmt.Errorf("failed to determine assignee: %w", err)
		}

		newTask := &models.AssignedTask{
			InstanceID: instanceID,
			StepID:     nextStep.ID,
			StepName:   nextStep.StepName,
			AssignedTo: assigneeID,
			AssignedBy: performedByID,
			Status:     "pending",
			DueDate:    instance.DueDate,
		}

		if err := s.taskRepo.Create(newTask); err != nil {
			return fmt.Errorf("failed to create new task: %w", err)
		}

		// Send email notification to new assignee
		assigneeUUID, _ := uuid.Parse(assigneeID)
		if assignee, err := s.userRepo.GetUserByID(assigneeUUID); err == nil {
			newTask.TaskDetails = &instance.TaskDetails // Add task details for email
			go s.notifyTaskAssignment(newTask, assignee, instance)
		}
	}

	// Get performer name for history
	performerUUID, err := uuid.Parse(performedByID)
	if err != nil {
		return fmt.Errorf("invalid performer ID: %w", err)
	}

	performer, err := s.employeeRepo.GetByUserID(performerUUID)
	if err != nil {
		return fmt.Errorf("failed to get performer: %w", err)
	}

	// Create history entry
	history := &models.WorkflowHistory{
		InstanceID:      instanceID,
		FromStepID:      &currentStep.ID,
		ToStepID:        nextStep.ID,
		ActionTaken:     action,
		PerformedBy:     performedByID,
		PerformedByName: performer.FirstName + " " + performer.LastName,
		Comments:        comments,
	}

	if err := s.historyRepo.Create(history); err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}

	// Update the underlying task (e.g., leave request) status if workflow is completed
	if isRejection || isCompletingFinalStep {
		if err := s.updateUnderlyingTaskStatus(instance, isRejection, performerUUID); err != nil {
			// Log error but don't fail the workflow - the workflow has already been completed
			fmt.Printf("Warning: Failed to update underlying task status: %v\n", err)
		}

		// Send email notification for workflow outcome
		if instance.TaskDetails.TaskType == "leave_request" {
			leaveRequestID, _ := uuid.Parse(instance.TaskDetails.TaskID)
			reviewerName := performer.FirstName + " " + performer.LastName
			go s.notifyLeaveRequestOutcome(leaveRequestID, !isRejection, reviewerName)
		}
	}

	return nil
}

// GetMyTasks retrieves all tasks assigned to a user
func (s *WorkflowService) GetMyTasks(userID string, statusFilter string) ([]models.AssignedTask, error) {
	if statusFilter != "" {
		return s.taskRepo.GetByAssignee(userID, statusFilter)
	}
	return s.taskRepo.GetByAssignee(userID)
}

// GetInstanceHistory retrieves the complete history of a workflow instance
func (s *WorkflowService) GetInstanceHistory(instanceID string) ([]models.WorkflowHistory, error) {
	return s.historyRepo.GetByInstanceID(instanceID)
}

// GetInstanceByTaskID retrieves a workflow instance by the associated task ID
func (s *WorkflowService) GetInstanceByTaskID(taskID string) (*models.WorkflowInstance, error) {
	return s.instanceRepo.GetByTaskID(taskID)
}

// Helper: Determine who to assign the task to based on step configuration
// Uses intelligent load balancing - assigns to the user with the required role who has the fewest pending tasks
func (s *WorkflowService) determineAssignee(step *models.WorkflowStep, taskDetails models.TaskDetails) (string, error) {
	if len(step.AllowedRoles) == 0 {
		return "", errors.New("no allowed roles defined for step")
	}

	// Try each allowed role and find the user with the fewest pending tasks
	for _, roleName := range step.AllowedRoles {
		user, err := s.userRepo.GetUserWithFewestTasksByRole(roleName)
		if err == nil && user != nil {
			// Found a user with this role - return their user_id as string
			return user.UserID.String(), nil
		}
	}

	return "", fmt.Errorf("no active user found with any of the allowed roles: %v", step.AllowedRoles)
}

// Helper: Check if user has permission to perform action on a step
func (s *WorkflowService) checkPermission(userID string, step *models.WorkflowStep) error {
	// Get user
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.GetUserByID(userUUID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if user has a role
	if user.Role == nil {
		return fmt.Errorf("user %s has no role assigned", user.Email)
	}

	// Check if user's role is in allowed roles
	for _, allowedRole := range step.AllowedRoles {
		if user.Role.Name == allowedRole {
			return nil
		}
	}

	return fmt.Errorf("user %s does not have permission to perform this action", user.Email)
}

// Helper: Update the underlying task (e.g., leave request) status when workflow completes
func (s *WorkflowService) updateUnderlyingTaskStatus(instance *models.WorkflowInstance, isRejection bool, reviewerID uuid.UUID) error {
	// Get employee from reviewer's user_id
	reviewer, err := s.employeeRepo.GetByUserID(reviewerID)
	if err != nil {
		return fmt.Errorf("failed to get reviewer employee: %w", err)
	}

	// Handle based on task type
	switch instance.TaskDetails.TaskType {
	case "leave_request":
		// Parse the leave request ID from task details
		leaveRequestID, err := uuid.Parse(instance.TaskDetails.TaskID)
		if err != nil {
			return fmt.Errorf("invalid leave request ID: %w", err)
		}

		// Update leave request status based on workflow outcome
		if isRejection {
			// Reject the leave request
			// Note: We need a reference to LeaveRequestService
			// For now, we'll update the database directly through a repository
			// TODO: Consider refactoring to avoid circular dependencies
			return s.updateLeaveRequestStatus(leaveRequestID, reviewer.ID, "rejected")
		} else {
			// Approve the leave request
			return s.updateLeaveRequestStatus(leaveRequestID, reviewer.ID, "approved")
		}

	default:
		// For other task types, do nothing for now
		return nil
	}
}

// Helper: Update leave request status directly
func (s *WorkflowService) updateLeaveRequestStatus(leaveRequestID, reviewerEmployeeID uuid.UUID, status string) error {
	// Get the leave request to access employee and leave type info
	req, err := s.leaveRequestRepo.GetByID(leaveRequestID)
	if err != nil {
		return fmt.Errorf("failed to get leave request: %w", err)
	}

	year := req.StartDate.Year()

	// Convert string status to LeaveRequestStatus type and update balances
	var leaveStatus models.LeaveRequestStatus
	switch status {
	case "approved":
		leaveStatus = models.LeaveStatusApproved
		// Update the leave request status
		if err := s.leaveRequestRepo.UpdateStatus(leaveRequestID, leaveStatus, &reviewerEmployeeID, ""); err != nil {
			return err
		}
		// Update balance: move from pending to used
		return s.leaveBalanceService.ApproveLeave(req.EmployeeID, req.LeaveTypeID, year, req.TotalDays)

	case "rejected":
		leaveStatus = models.LeaveStatusRejected
		// Update the leave request status
		if err := s.leaveRequestRepo.UpdateStatus(leaveRequestID, leaveStatus, &reviewerEmployeeID, ""); err != nil {
			return err
		}
		// Update balance: decrement pending (return the days)
		return s.leaveBalanceService.DecrementPending(req.EmployeeID, req.LeaveTypeID, year, req.TotalDays)

	default:
		return fmt.Errorf("invalid leave request status: %s", status)
	}
}

// Helper: Send email notification when task is assigned
func (s *WorkflowService) notifyTaskAssignment(task *models.AssignedTask, assignee *models.User, instance *models.WorkflowInstance) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in notifyTaskAssignment: %v\n", r)
		}
	}()

	if s.emailService == nil {
		return // Email service not configured
	}

	// Get assignee employee details
	assigneeEmployee, err := s.employeeRepo.GetByUserID(assignee.UserID)
	if err != nil {
		fmt.Printf("Warning: Failed to get assignee employee for email notification: %v\n", err)
		return
	}

	// Build email based on task type
	var subject string
	var htmlBody string

	if task.TaskDetails != nil && task.TaskDetails.TaskType == "leave_request" {
		// For leave request workflow
		subject = "New Leave Request Assigned for Review"

		// Get leave request details
		leaveReq, err := s.leaveRequestRepo.GetByID(uuid.MustParse(task.TaskDetails.TaskID))
		if err != nil {
			fmt.Printf("Warning: Failed to get leave request for email: %v\n", err)
			return
		}

		// Get leave type name
		var leaveTypeName string
		if leaveReq.LeaveType != nil {
			leaveTypeName = leaveReq.LeaveType.Name
		}

		htmlBody = email.LeaveRequestAssignedTemplate(
			task.TaskDetails.SenderDetails.SenderName,
			leaveTypeName,
			leaveReq.TotalDays,
			leaveReq.StartDate.Format("2006-01-02"),
			leaveReq.EndDate.Format("2006-01-02"),
		)
	} else {
		// Generic task assignment
		subject = fmt.Sprintf("New Task Assigned: %s", task.StepName)
		htmlBody = email.GenericTaskAssignedTemplate(
			assigneeEmployee.FirstName,
			task.StepName,
			task.TaskDetails.TaskDescription,
		)
	}

	// Send email
	if err := s.emailService.SendEmail([]string{assignee.Email}, subject, htmlBody); err != nil {
		fmt.Printf("Warning: Failed to send task assignment email: %v\n", err)
	}
}

// Helper: Send email notification when leave request is approved/rejected
func (s *WorkflowService) notifyLeaveRequestOutcome(leaveRequestID uuid.UUID, isApproved bool, reviewerName string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in notifyLeaveRequestOutcome: %v\n", r)
		}
	}()

	if s.emailService == nil {
		return // Email service not configured
	}

	// Get leave request
	leaveReq, err := s.leaveRequestRepo.GetByID(leaveRequestID)
	if err != nil {
		fmt.Printf("Warning: Failed to get leave request for outcome email: %v\n", err)
		return
	}

	// Get employee
	employee, err := s.employeeRepo.GetByID(leaveReq.EmployeeID)
	if err != nil {
		fmt.Printf("Warning: Failed to get employee for outcome email: %v\n", err)
		return
	}

	// Get user email
	user, err := s.userRepo.GetUserByID(*employee.UserID)
	if err != nil {
		fmt.Printf("Warning: Failed to get user for outcome email: %v\n", err)
		return
	}

	// Get leave type name
	var leaveTypeName string
	if leaveReq.LeaveType != nil {
		leaveTypeName = leaveReq.LeaveType.Name
	}

	var subject string
	var htmlBody string

	if isApproved {
		subject = "Leave Request Approved"
		htmlBody = email.LeaveRequestApprovedTemplate(
			employee.FirstName,
			leaveTypeName,
			leaveReq.TotalDays,
			leaveReq.StartDate.Format("2006-01-02"),
			leaveReq.EndDate.Format("2006-01-02"),
			reviewerName,
		)
	} else {
		subject = "Leave Request Not Approved"
		htmlBody = email.LeaveRequestRejectedTemplate(
			employee.FirstName,
			leaveTypeName,
			leaveReq.TotalDays,
			leaveReq.StartDate.Format("2006-01-02"),
			leaveReq.EndDate.Format("2006-01-02"),
			reviewerName,
			"", // No reason for now, can be added later
		)
	}

	// Send email
	if err := s.emailService.SendEmail([]string{user.Email}, subject, htmlBody); err != nil {
		fmt.Printf("Warning: Failed to send leave request outcome email: %v\n", err)
	}
}
