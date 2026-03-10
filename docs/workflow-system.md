# Workflow System Documentation

## Overview

The workflow system provides a flexible, reusable framework for managing multi-step approval processes in the HR system. It separates workflow templates (blueprints) from workflow instances (actual executions).

## Architecture

### Core Concepts

1. **Workflow Template** - The blueprint defining the approval process
2. **Workflow Instance** - A specific execution (e.g., one leave request going through approval)
3. **Assigned Task** - An action item for a specific user
4. **Workflow History** - Complete audit trail of all actions

## Data Models

### 1. Workflow (Template)
The reusable blueprint for a process.

```go
type Workflow struct {
    ID          string    // Unique identifier
    Name        string    // e.g., "Leave Request Approval"
    Description string
    IsActive    bool
    CreatedBy   string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 2. WorkflowStep
Individual stages in the workflow template.

```go
type WorkflowStep struct {
    ID                   string
    WorkflowID           string
    StepName             string    // e.g., "HR Review"
    StepOrder            int       // Position in sequence
    Initial              bool      // First step?
    Final                bool      // Last step?
    AllowedRoles         []string  // Roles that can act
    RequiresAllApprovers bool      // All must approve?
    MinApprovals         int       // Minimum approvals needed
    CreatedAt            time.Time
}
```

### 3. WorkflowTransition
Valid movements between steps.

```go
type WorkflowTransition struct {
    ID             string
    WorkflowID     string
    FromStepID     string
    ToStepID       string
    ActionName     string    // "approve", "reject", "submit"
    ConditionType  string    // "user_role", "assigned_user_only"
    ConditionValue string    // JSON for complex conditions
    CreatedAt      time.Time
}
```

### 4. WorkflowInstance
A single execution of a workflow.

```go
type WorkflowInstance struct {
    ID            string
    WorkflowID    string      // References template
    CurrentStepID string      // Where is it now?
    Status        string      // "pending", "in_progress", "completed", "rejected", "cancelled"
    TaskDetails   TaskDetails // The actual data
    CreatedBy     string
    CreatedAt     time.Time
    UpdatedAt     time.Time
    CompletedAt   *time.Time  // Null until finished
    DueDate       *time.Time
    Priority      string      // "low", "medium", "high", "urgent"
}
```

### 5. AssignedTask
Action item for a specific user.

```go
type AssignedTask struct {
    ID          string
    InstanceID  string     // References WorkflowInstance
    StepID      string
    StepName    string     // Denormalized
    AssignedTo  string     // User who needs to act
    AssignedBy  string     // Who assigned it
    Status      string     // "pending", "in_progress", "completed", "skipped"
    DueDate     *time.Time
    CompletedAt *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 6. WorkflowHistory
Complete audit trail.

```go
type WorkflowHistory struct {
    ID              string
    InstanceID      string
    FromStepID      *string   // Null for creation
    ToStepID        string
    ActionTaken     string    // "submit", "approve", "reject"
    PerformedBy     string
    PerformedByName string    // Denormalized
    Comments        string
    Metadata        string    // JSON
    Timestamp       time.Time
}
```

## Example: Leave Request Flow

### Setup (One-time)

1. **Create Workflow Template**: "Leave Request Approval"
2. **Define Steps**:
   - Submit (initial)
   - HR Review
   - Manager Approval
   - Completed (final)
3. **Define Transitions**:
   - Submit → HR Review (action: submit)
   - HR Review → Manager Approval (action: approve)
   - HR Review → Submit (action: reject)
   - Manager Approval → Completed (action: approve)
   - Manager Approval → HR Review (action: reject)

### Execution (Per Request)

**Step 1: Employee Submits Leave Request**
```json
{
  "workflow_instance": {
    "workflow_id": "leave-approval-workflow",
    "current_step_id": "hr-review-step",
    "status": "in_progress",
    "task_details": {
      "task_id": "leave-request-123",
      "task_type": "leave_request",
      "task_description": "Annual Leave: 5 days (2026-03-01 to 2026-03-05)",
      "sender_details": {
        "sender_id": "emp-456",
        "sender_name": "John Doe",
        "position": "Software Engineer",
        "department": "Engineering"
      }
    },
    "created_by": "emp-456",
    "priority": "medium"
  },
  "assigned_task": {
    "instance_id": "instance-001",
    "step_id": "hr-review-step",
    "step_name": "HR Review",
    "assigned_to": "hr-manager-789",
    "assigned_by": "system",
    "status": "pending"
  },
  "history": {
    "instance_id": "instance-001",
    "from_step_id": null,
    "to_step_id": "hr-review-step",
    "action_taken": "submit",
    "performed_by": "emp-456",
    "performed_by_name": "John Doe",
    "comments": "Submitted leave request for annual leave"
  }
}
```

**Step 2: HR Manager Approves**
```json
{
  "workflow_instance": {
    "current_step_id": "manager-approval-step",
    "status": "in_progress"
  },
  "previous_task": {
    "id": "task-001",
    "status": "completed",
    "completed_at": "2026-02-20T10:30:00Z"
  },
  "new_task": {
    "instance_id": "instance-001",
    "step_id": "manager-approval-step",
    "step_name": "Manager Approval",
    "assigned_to": "dept-head-999",
    "assigned_by": "hr-manager-789",
    "status": "pending"
  },
  "history": {
    "instance_id": "instance-001",
    "from_step_id": "hr-review-step",
    "to_step_id": "manager-approval-step",
    "action_taken": "approve",
    "performed_by": "hr-manager-789",
    "performed_by_name": "Jane Smith (HR Manager)",
    "comments": "Leave balance verified, approved for manager review"
  }
}
```

**Step 3: Department Head Approves**
```json
{
  "workflow_instance": {
    "current_step_id": "completed-step",
    "status": "completed",
    "completed_at": "2026-02-20T14:15:00Z"
  },
  "task": {
    "id": "task-002",
    "status": "completed",
    "completed_at": "2026-02-20T14:15:00Z"
  },
  "history": {
    "instance_id": "instance-001",
    "from_step_id": "manager-approval-step",
    "to_step_id": "completed-step",
    "action_taken": "approve",
    "performed_by": "dept-head-999",
    "performed_by_name": "Bob Johnson (Dept Head)",
    "comments": "Team coverage confirmed, approved"
  }
}
```

## Key Features

### 1. Separation of Template and Instance
- **One template** → **Many instances**
- Easy to modify workflow without affecting active requests
- Can create different workflows for different leave types

### 2. Flexible Routing
- Role-based: "Only HR Managers can approve this step"
- Conditional: "If amount > $1000, requires CFO approval"
- Parallel approvals: "Requires 2 out of 3 managers to approve"

### 3. Complete Audit Trail
- Every action is logged
- Who did what, when, and why
- Comments preserved for compliance

### 4. Task Assignment
- Clear action items for each user
- Separate "my tasks" from "workflow state"
- Support for reassignment

### 5. Status Tracking
- Instance-level: Overall workflow progress
- Task-level: Individual action items
- Easy to query: "Show me all my pending tasks"

## Database Schema

See `migrations/009_create_workflow_tables.sql` for complete schema.

### Key Indexes
- `assigned_tasks.assigned_to` - Fast "my tasks" queries
- `workflow_instances.status` - Filter by status
- `workflow_history.instance_id` - Quick audit trail retrieval

## Use Cases

1. **Leave Requests** - Submit → HR → Manager → Approve
2. **Employee Onboarding** - Request → HR Setup → IT Setup → Manager Assignment
3. **Equipment Requests** - Submit → Manager → IT Approval → Procurement
4. **Performance Reviews** - Self-Assessment → Manager Review → HR Review → Finalize
5. **Expense Reimbursements** - Submit → Manager → Finance → Payment

## Next Steps

1. ✅ Models defined
2. ✅ Database migrations created
3. ✅ Seed data for Leave Approval workflow
4. 🔄 Create repository layer
5. 🔄 Create service layer
6. 🔄 Create handlers and routes
7. 🔄 Integrate with existing leave request system

## Migration Commands

```bash
# Run migrations
psql -U your_user -d hr_system -f migrations/009_create_workflow_tables.sql
psql -U your_user -d hr_system -f migrations/010_seed_leave_approval_workflow.sql

# Verify
psql -U your_user -d hr_system -c "SELECT * FROM workflows;"
psql -U your_user -d hr_system -c "SELECT * FROM workflow_steps ORDER BY step_order;"
```
