package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"lodge-system/internal/models"
	"lodge-system/internal/repository"
	"lodge-system/internal/utils/email"
)

type WorkflowService struct {
	workflowRepo *repository.WorkflowRepository
	instanceRepo *repository.WorkflowInstanceRepository
	taskRepo     *repository.AssignedTaskRepository
	historyRepo  *repository.WorkflowHistoryRepository
	userRepo     *repository.UserRepository
	clientRepo   *repository.ClientRepository
	emailService *email.EmailService
}

func NewWorkflowService(
	workflowRepo *repository.WorkflowRepository,
	instanceRepo *repository.WorkflowInstanceRepository,
	taskRepo *repository.AssignedTaskRepository,
	historyRepo *repository.WorkflowHistoryRepository,
	userRepo *repository.UserRepository,
	clientRepo *repository.ClientRepository,
	emailService *email.EmailService,
) *WorkflowService {
	return &WorkflowService{
		workflowRepo: workflowRepo,
		instanceRepo: instanceRepo,
		taskRepo:     taskRepo,
		historyRepo:  historyRepo,
		userRepo:     userRepo,
		clientRepo:   clientRepo,
		emailService: emailService,
	}
}

// InitiateWorkflow starts a new workflow instance using workflow type.
// orgID scopes the workflow template lookup and stamps the instance.
func (s *WorkflowService) InitiateWorkflow(
	workflowType models.WorkflowType,
	taskDetails models.TaskDetails,
	initiatorID string,
	priority string,
	dueDate interface{},
	orgID string,
) (*models.WorkflowInstance, error) {
	workflow, err := s.workflowRepo.GetByType(workflowType, orgID)
	if err != nil {
		return nil, fmt.Errorf("workflow not found for type %s: %w", workflowType, err)
	}

	firstActionStep, initialStepID, err := s.workflowRepo.GetFirstActionStep(workflow.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get first action step: %w", err)
	}

	instance := &models.WorkflowInstance{
		OrgID:         orgID,
		WorkflowID:    workflow.ID,
		CurrentStepID: firstActionStep.ID,
		Status:        "in_progress",
		TaskDetails:   taskDetails,
		CreatedBy:     initiatorID,
		Priority:      priority,
	}

	if err := s.instanceRepo.Create(instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	assigneeID, err := s.determineAssignee(orgID, firstActionStep, taskDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to determine assignee: %w", err)
	}

	task := &models.AssignedTask{
		OrgID:      orgID,
		InstanceID: instance.ID,
		StepID:     firstActionStep.ID,
		StepName:   firstActionStep.StepName,
		AssignedTo: assigneeID,
		AssignedBy: initiatorID,
		Status:     "pending",
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	assigneeUUID, _ := uuid.Parse(assigneeID)
	if assignee, err := s.userRepo.GetUserByID(assigneeUUID); err == nil {
		task.TaskDetails = &taskDetails
		go s.notifyTaskAssignment(task, assignee, instance)
	}

	initiatorUUID, err := uuid.Parse(initiatorID)
	if err != nil {
		return nil, fmt.Errorf("invalid initiator ID: %w", err)
	}

	initiatorUser, err := s.userRepo.GetUserByID(initiatorUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get initiator: %w", err)
	}

	history := &models.WorkflowHistory{
		InstanceID:      instance.ID,
		FromStepID:      &initialStepID,
		ToStepID:        firstActionStep.ID,
		ActionTaken:     "submit",
		PerformedBy:     initiatorID,
		PerformedByName: initiatorUser.Email,
		Comments:        fmt.Sprintf("Initiated %s workflow", workflow.Name),
	}

	if err := s.historyRepo.Create(history); err != nil {
		return nil, fmt.Errorf("failed to create history: %w", err)
	}

	return instance, nil
}

// ProcessAction processes an action on a workflow instance, scoped to org.
func (s *WorkflowService) ProcessAction(
	instanceID string,
	action string,
	performedByID string,
	comments string,
	orgID string,
) error {
	instance, err := s.instanceRepo.GetByID(instanceID, orgID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	if instance.Status == "completed" || instance.Status == "cancelled" {
		return errors.New("workflow instance is already closed")
	}

	currentStep, err := s.workflowRepo.GetStepByID(instance.CurrentStepID)
	if err != nil {
		return fmt.Errorf("current step not found: %w", err)
	}

	if err := s.checkPermission(performedByID, currentStep); err != nil {
		return err
	}

	isRejection := action == "reject" || action == "rejected" || action == "deny"
	isCompletingFinalStep := currentStep.Final && !isRejection

	var nextStep *models.WorkflowStep

	if isCompletingFinalStep {
		nextStep = currentStep
	} else {
		transition, err := s.workflowRepo.GetTransitionByAction(currentStep.ID, action)
		if err == sql.ErrNoRows {
			return fmt.Errorf("action '%s' is not valid from current step '%s'", action, currentStep.StepName)
		}
		if err != nil {
			return fmt.Errorf("failed to get transition: %w", err)
		}

		nextStep, err = s.workflowRepo.GetStepByID(transition.ToStepID)
		if err != nil {
			return fmt.Errorf("next step not found: %w", err)
		}
	}

	newStatus := "in_progress"
	if isRejection {
		newStatus = "rejected"
	} else if isCompletingFinalStep {
		newStatus = "completed"
	}

	if err := s.instanceRepo.UpdateStep(instanceID, nextStep.ID, newStatus, instance.OrgID); err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	activeTask, err := s.taskRepo.GetActiveTaskForInstance(instanceID, instance.OrgID)
	if err == nil && activeTask != nil {
		if err := s.taskRepo.Complete(activeTask.ID, instance.OrgID); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}
	}

	if isRejection || isCompletingFinalStep {
		if err := s.instanceRepo.Complete(instanceID, instance.OrgID); err != nil {
			return fmt.Errorf("failed to complete instance: %w", err)
		}

		// Update the real entity status based on workflow outcome
		if instance.TaskDetails.TaskType != "" && instance.TaskDetails.TaskID != "" {
			var entityStatus string
			switch instance.TaskDetails.TaskType {
			case "booking":
				if isCompletingFinalStep {
					entityStatus = models.BookingStatusConfirmed
				} else {
					entityStatus = models.BookingStatusCancelled
				}
			}
			if entityStatus != "" {
				if err := s.workflowRepo.UpdateEntityStatus(instance.TaskDetails.TaskType, instance.TaskDetails.TaskID, entityStatus, instance.OrgID); err != nil {
					fmt.Printf("warning: failed to update %s status after workflow outcome: %v\n", instance.TaskDetails.TaskType, err)
				} else {
					go s.notifyGuestBookingOutcome(instance.OrgID, instance.TaskDetails, entityStatus)
				}
			}
		}
	} else {
		assigneeID, err := s.determineAssignee(instance.OrgID, nextStep, instance.TaskDetails)
		if err != nil {
			return fmt.Errorf("failed to determine assignee: %w", err)
		}

		newTask := &models.AssignedTask{
			OrgID:      instance.OrgID,
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

		assigneeUUID, _ := uuid.Parse(assigneeID)
		if assignee, err := s.userRepo.GetUserByID(assigneeUUID); err == nil {
			newTask.TaskDetails = &instance.TaskDetails
			go s.notifyTaskAssignment(newTask, assignee, instance)
		}
	}

	performerUUID, err := uuid.Parse(performedByID)
	if err != nil {
		return fmt.Errorf("invalid performer ID: %w", err)
	}

	performer, err := s.userRepo.GetUserByID(performerUUID)
	if err != nil {
		return fmt.Errorf("failed to get performer: %w", err)
	}

	history := &models.WorkflowHistory{
		InstanceID:      instanceID,
		FromStepID:      &currentStep.ID,
		ToStepID:        nextStep.ID,
		ActionTaken:     action,
		PerformedBy:     performedByID,
		PerformedByName: performer.Email,
		Comments:        comments,
	}

	if err := s.historyRepo.Create(history); err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}

	return nil
}

// GetMyTasks retrieves all tasks assigned to a user, scoped to org.
func (s *WorkflowService) GetMyTasks(orgID, userID string, statusFilter string) ([]models.AssignedTask, error) {
	if statusFilter != "" {
		return s.taskRepo.GetByAssignee(orgID, userID, statusFilter)
	}
	return s.taskRepo.GetByAssignee(orgID, userID)
}

// GetInstanceHistory retrieves the complete history of a workflow instance, scoped to org.
func (s *WorkflowService) GetInstanceHistory(instanceID, orgID string) ([]models.WorkflowHistory, error) {
	if _, err := s.instanceRepo.GetByID(instanceID, orgID); err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}
	return s.historyRepo.GetByInstanceID(instanceID)
}

// GetInstanceByTaskID retrieves a workflow instance by the associated task ID, scoped to org.
func (s *WorkflowService) GetInstanceByTaskID(taskID, orgID string) (*models.WorkflowInstance, error) {
	return s.instanceRepo.GetByTaskID(taskID, orgID)
}

// determineAssignee finds the user with the required role who has the fewest pending tasks,
// scoped to the given org so tasks are never assigned cross-org.
func (s *WorkflowService) determineAssignee(orgID string, step *models.WorkflowStep, taskDetails models.TaskDetails) (string, error) {
	if len(step.AllowedRoles) == 0 {
		return "", errors.New("no allowed roles defined for step")
	}

	orgUUID, _ := uuid.Parse(orgID)

	for _, roleName := range step.AllowedRoles {
		user, err := s.userRepo.GetUserWithFewestTasksByRole(orgUUID, roleName)
		if err == nil && user != nil {
			return user.UserID.String(), nil
		}
	}

	return "", fmt.Errorf("no active user found with any of the allowed roles: %v", step.AllowedRoles)
}

// checkPermission verifies the user's role is allowed to act on the current step
func (s *WorkflowService) checkPermission(userID string, step *models.WorkflowStep) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.GetUserByID(userUUID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.Role == nil {
		return fmt.Errorf("user %s has no role assigned", user.Email)
	}

	for _, allowedRole := range step.AllowedRoles {
		if user.Role.Name == allowedRole {
			return nil
		}
	}

	return fmt.Errorf("user %s does not have permission to perform this action", user.Email)
}

// notifyTaskAssignment sends an email when a task is assigned
// notifyGuestBookingOutcome emails the guest when their booking is approved or rejected.
func (s *WorkflowService) notifyGuestBookingOutcome(orgID string, details models.TaskDetails, entityStatus string) {
	if s.emailService == nil || s.clientRepo == nil {
		return
	}

	clientID, err := uuid.Parse(details.SenderDetails.SenderID)
	if err != nil {
		fmt.Printf("warning: invalid client ID in task details: %v\n", err)
		return
	}

	orgUUID, _ := uuid.Parse(orgID)
	profile, err := s.clientRepo.GetIndividualByID(clientID, orgUUID)
	if err != nil {
		fmt.Printf("warning: could not find guest profile for booking outcome email: %v\n", err)
		return
	}

	var subject, htmlBody string
	switch entityStatus {
	case models.BookingStatusConfirmed:
		subject = "Your Booking is Confirmed — The Sanctuary"
		htmlBody = email.BookingApprovedTemplate(profile.FullName, details.TaskID, details.TaskDescription)
	case models.BookingStatusCancelled:
		subject = "Booking Update — The Sanctuary"
		htmlBody = email.BookingRejectedTemplate(profile.FullName, details.TaskID, details.TaskDescription)
	default:
		return
	}

	if err := s.emailService.SendEmail([]string{profile.Email}, subject, htmlBody); err != nil {
		fmt.Printf("warning: failed to send booking outcome email to %s: %v\n", profile.Email, err)
	}
}

func (s *WorkflowService) notifyTaskAssignment(task *models.AssignedTask, assignee *models.User, instance *models.WorkflowInstance) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in notifyTaskAssignment: %v\n", r)
		}
	}()

	if s.emailService == nil {
		return
	}

	var subject, htmlBody string

	if task.TaskDetails != nil && task.TaskDetails.TaskType == "booking" {
		subject = fmt.Sprintf("Booking Approval Required — %s", task.TaskDetails.SenderDetails.SenderName)
		htmlBody = email.BookingTaskAssignedTemplate(
			assignee.FullName,
			task.TaskDetails.TaskID,
			task.TaskDetails.TaskDescription,
			task.TaskDetails.SenderDetails.SenderName,
			task.TaskDetails.SenderDetails.Position,
		)
	} else {
		subject = fmt.Sprintf("New Task Assigned: %s", task.StepName)
		description := ""
		if task.TaskDetails != nil {
			description = task.TaskDetails.TaskDescription
		}
		htmlBody = email.GenericTaskAssignedTemplate(assignee.FullName, task.StepName, description)
	}

	if err := s.emailService.SendEmail([]string{assignee.Email}, subject, htmlBody); err != nil {
		fmt.Printf("Warning: Failed to send task assignment email: %v\n", err)
	}
}
