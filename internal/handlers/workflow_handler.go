package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/services"
)

type WorkflowHandler struct {
	service *services.WorkflowService
}

func NewWorkflowHandler(service *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		service: service,
	}
}

// GetMyTasks retrieves all tasks assigned to the authenticated user
func (h *WorkflowHandler) GetMyTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get status filter from query params (optional)
	status := r.URL.Query().Get("status") // "pending", "completed", etc.

	tasks, err := h.service.GetMyTasks(userID.String(), status)
	if err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetMyPendingTasks retrieves only pending tasks for the authenticated user
func (h *WorkflowHandler) GetMyPendingTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.service.GetMyTasks(userID.String(), "pending")
	if err != nil {
		http.Error(w, "Failed to retrieve pending tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// ProcessActionRequest represents the request body for processing an action
type ProcessActionRequest struct {
	Action   string `json:"action"`   // "approve", "reject", "submit"
	Comments string `json:"comments"` // Optional comments
}

// ProcessAction handles workflow actions (approve, reject, etc.)
func (h *WorkflowHandler) ProcessAction(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get instance ID from path
	instanceID := r.PathValue("id")
	if instanceID == "" {
		http.Error(w, "Instance ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req ProcessActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate action
	if req.Action == "" {
		http.Error(w, "Action is required", http.StatusBadRequest)
		return
	}

	// Process the action
	err := h.service.ProcessAction(instanceID, req.Action, userID.String(), req.Comments)
	if err != nil {
		// Check for specific errors
		if err.Error() == "workflow instance is already closed" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if err.Error() == "user does not have permission to perform this action" ||
		   err.Error()[:4] == "user" && err.Error()[len(err.Error())-31:] == "does not have permission to perform this action" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Action processed successfully",
		"action":  req.Action,
	})
}

// GetInstanceHistory retrieves the complete history of a workflow instance
func (h *WorkflowHandler) GetInstanceHistory(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")
	if instanceID == "" {
		http.Error(w, "Instance ID is required", http.StatusBadRequest)
		return
	}

	history, err := h.service.GetInstanceHistory(instanceID)
	if err != nil {
		http.Error(w, "Failed to retrieve history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"instance_id": instanceID,
		"history":     history,
		"count":       len(history),
	})
}

// InitiateWorkflowRequest represents the request to start a new workflow
type InitiateWorkflowRequest struct {
	WorkflowType models.WorkflowType `json:"workflow_type"` // e.g., "LEAVE_REQUEST", "EMPLOYEE_ONBOARDING"
	TaskDetails  models.TaskDetails  `json:"task_details"`
	Priority     string              `json:"priority"`  // "low", "medium", "high", "urgent"
	DueDate      *time.Time          `json:"due_date"`  // Optional
}

// InitiateWorkflow starts a new workflow instance (used internally by other services)
func (h *WorkflowHandler) InitiateWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req InitiateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WorkflowType == "" {
		http.Error(w, "Workflow type is required", http.StatusBadRequest)
		return
	}
	if req.TaskDetails.TaskID == "" || req.TaskDetails.TaskType == "" {
		http.Error(w, "Task details (task_id and task_type) are required", http.StatusBadRequest)
		return
	}

	// Default priority
	if req.Priority == "" {
		req.Priority = "medium"
	}

	// Initiate workflow
	instance, err := h.service.InitiateWorkflow(
		req.WorkflowType,
		req.TaskDetails,
		userID.String(),
		req.Priority,
		req.DueDate,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Workflow initiated successfully",
		"instance": instance,
	})
}

// GetInstanceByTaskID retrieves a workflow instance by the task ID
func (h *WorkflowHandler) GetInstanceByTaskID(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	instance, err := h.service.GetInstanceByTaskID(taskID)
	if err != nil {
		http.Error(w, "Workflow instance not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instance)
}

// GetTaskDetails retrieves detailed information about a specific task
func (h *WorkflowHandler) GetTaskDetails(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all tasks for the user to find this one
	tasks, err := h.service.GetMyTasks(userID.String(), "")
	if err != nil {
		http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		return
	}

	// Find the specific task
	var foundTask *models.AssignedTask
	for _, task := range tasks {
		if task.ID == taskID {
			foundTask = &task
			break
		}
	}

	if foundTask == nil {
		http.Error(w, "Task not found or not assigned to you", http.StatusNotFound)
		return
	}

	// Get the workflow instance for more context
	instance, err := h.service.GetInstanceByTaskID(foundTask.InstanceID)
	if err != nil {
		// Task found but instance not found - just return task
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(foundTask)
		return
	}

	// Get history for full context
	history, _ := h.service.GetInstanceHistory(instance.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task":     foundTask,
		"instance": instance,
		"history":  history,
	})
}