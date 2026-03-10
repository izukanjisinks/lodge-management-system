package models

import "time"

// WorkflowType represents the type/purpose of a workflow
type WorkflowType string

// Workflow type constants
const (
	WorkflowTypeLeaveRequest WorkflowType = "LEAVE_REQUEST"
)

// GetAllWorkflowTypes returns all available workflow types
func GetAllWorkflowTypes() []WorkflowType {
	return []WorkflowType{
		WorkflowTypeLeaveRequest,
	}
}

// WorkflowTypeInfo provides metadata about a workflow type
type WorkflowTypeInfo struct {
	Type        WorkflowType `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
}

// GetWorkflowTypeInfo returns metadata for all workflow types
func GetWorkflowTypeInfo() []WorkflowTypeInfo {
	return []WorkflowTypeInfo{
		{
			Type:        WorkflowTypeLeaveRequest,
			Name:        "Leave Request",
			Description: "Workflow for managing employee leave requests and approvals",
		},
	}
}

// Workflow represents a workflow template (the blueprint)
type Workflow struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	WorkflowType *WorkflowType `json:"workflow_type,omitempty"` // Unique identifier for workflow purpose
	IsActive     bool          `json:"is_active"`
	CreatedBy    string        `json:"created_by"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// WorkflowWithCounts represents a workflow template with step and transition counts
type WorkflowWithCounts struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	WorkflowType    *WorkflowType `json:"workflow_type,omitempty"`
	IsActive        bool          `json:"is_active"`
	CreatedBy       string        `json:"created_by"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	StepCount       int           `json:"step_count"`
	TransitionCount int           `json:"transition_count"`
}

// WorkflowStep represents a step/stage in a workflow template
type WorkflowStep struct {
	ID                   string    `json:"id"`
	WorkflowID           string    `json:"workflow_id"`
	StepName             string    `json:"step_name"`
	StepOrder            int       `json:"step_order"`
	Initial              bool      `json:"initial"`              // First step in workflow
	Final                bool      `json:"final"`                // Last step in workflow
	AllowedRoles         []string  `json:"allowed_roles"`        // Roles that can act on this step
	RequiresAllApprovers bool      `json:"requires_all_approvers"` // true = all must approve, false = any one
	MinApprovals         int       `json:"min_approvals"`        // Minimum approvals needed (0 = not used)
	CreatedAt            time.Time `json:"created_at"`
}

// WorkflowTransition represents a transition between steps
type WorkflowTransition struct {
	ID             string    `json:"id"`
	WorkflowID     string    `json:"workflow_id"`
	FromStepID     string    `json:"from_step_id"`
	ToStepID       string    `json:"to_step_id"`
	ActionName     string    `json:"action_name"`      // "submit", "approve", "reject", "reassign"
	ConditionType  *string   `json:"condition_type"`   // e.g., "user_role", "assigned_user_only"
	ConditionValue *string   `json:"condition_value"`  // JSON for complex conditions
	CreatedAt      time.Time `json:"created_at"`
}

// WorkflowInstance represents a single execution of a workflow template
// This tracks the overall progress of one specific case (e.g., one leave request)
type WorkflowInstance struct {
	ID            string      `json:"id"`
	WorkflowID    string      `json:"workflow_id"`    // References the Workflow template
	CurrentStepID string      `json:"current_step_id"` // Where is this instance currently?
	Status        string      `json:"status"`          // "pending", "in_progress", "completed", "rejected", "cancelled"
	TaskDetails   TaskDetails `json:"task_details"`    // The actual data for this instance
	CreatedBy     string      `json:"created_by"`      // Who initiated this workflow
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	CompletedAt   *time.Time  `json:"completed_at,omitempty"` // When workflow finished (null if ongoing)
	DueDate       *time.Time  `json:"due_date,omitempty"`     // Optional deadline
	Priority      string      `json:"priority"`               // "low", "medium", "high", "urgent"
}

// AssignedTask represents an action item assigned to a specific user for a workflow instance
// Multiple tasks can exist for the same instance (e.g., parallel approvals)
type AssignedTask struct {
	ID          string       `json:"id"`
	InstanceID  string       `json:"instance_id"` // References WorkflowInstance
	StepID      string       `json:"step_id"`     // Which step is this task for?
	StepName    string       `json:"step_name"`   // Denormalized for easy display
	AssignedTo  string       `json:"assigned_to"` // User ID who needs to act
	AssignedBy  string       `json:"assigned_by"` // User ID who assigned this task
	Status      string       `json:"status"`      // "pending", "in_progress", "completed", "skipped"
	DueDate     *time.Time   `json:"due_date,omitempty"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	TaskDetails *TaskDetails `json:"task_details,omitempty"` // Details from the workflow instance
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// TaskDetails contains the context and data for a workflow instance
type TaskDetails struct {
	TaskID          string        `json:"task_id"`   // e.g., leave_request_id, employee_onboarding_id
	TaskType        string        `json:"task_type"` // e.g., "leave_request", "employee_onboarding"
	TaskDescription string        `json:"task_description"`
	SenderDetails   SenderDetails `json:"sender_details"`
	Metadata        string        `json:"metadata,omitempty"` // JSON for additional flexible data
}

// SenderDetails contains information about who initiated the workflow
type SenderDetails struct {
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Position   string `json:"position"`
	Department string `json:"department"`
}

// WorkflowHistory represents the audit trail of a workflow instance
type WorkflowHistory struct {
	ID              string    `json:"id"`
	InstanceID      string    `json:"instance_id"`
	FromStepID      *string   `json:"from_step_id"`      // Nullable for initial creation
	ToStepID        string    `json:"to_step_id"`
	ActionTaken     string    `json:"action_taken"`      // "submit", "approve", "reject", "reassign"
	PerformedBy     string    `json:"performed_by"`      // User ID
	PerformedByName string    `json:"performed_by_name"` // Denormalized for display
	Comments        string    `json:"comments"`
	Metadata        string    `json:"metadata,omitempty"` // JSON for additional context
	Timestamp       time.Time `json:"timestamp"`
}
