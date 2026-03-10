# Workflow System - Quick Start Guide

## 🚀 TL;DR

The workflow system provides a flexible approval process framework. Instead of hardcoding approval logic, define workflows once and reuse them.

## 📦 What You Need

```go
import (
    "hr-system/internal/models"
    "hr-system/internal/repository"
    "hr-system/internal/services"
)

// Initialize repositories
workflowRepo := repository.NewWorkflowRepository()
instanceRepo := repository.NewWorkflowInstanceRepository()
taskRepo := repository.NewAssignedTaskRepository()
historyRepo := repository.NewWorkflowHistoryRepository()
userRepo := repository.NewUserRepository()

// Create service
workflowService := services.NewWorkflowService(
    workflowRepo,
    instanceRepo,
    taskRepo,
    historyRepo,
    userRepo,
)
```

## 🎬 Common Use Cases

### 1. Start a Leave Request Workflow

```go
instance, err := workflowService.InitiateWorkflow(
    "Leave Request Approval",  // workflow name
    models.TaskDetails{
        TaskID:          leaveRequest.ID,
        TaskType:        "leave_request",
        TaskDescription: fmt.Sprintf("%s: %d days", leaveType, days),
        SenderDetails: models.SenderDetails{
            SenderID:   employee.ID,
            SenderName: employee.FirstName + " " + employee.LastName,
            Position:   employee.Position.Name,
            Department: employee.Department.Name,
        },
    },
    employee.UserID,  // who initiated
    "medium",         // priority
    &dueDate,         // optional deadline
)
```

### 2. HR Manager Approves

```go
err := workflowService.ProcessAction(
    instance.ID,
    "approve",
    hrManagerID,
    "Leave balance verified. Approved for manager review.",
)
```

### 3. Department Head Rejects

```go
err := workflowService.ProcessAction(
    instance.ID,
    "reject",
    deptHeadID,
    "Insufficient coverage during requested period. Please reschedule.",
)
```

### 4. Get My Pending Tasks

```go
// All pending tasks
tasks, err := workflowService.GetMyTasks(userID, "pending")

// All tasks (any status)
allTasks, err := workflowService.GetMyTasks(userID, "")
```

### 5. View Workflow History

```go
history, err := workflowService.GetInstanceHistory(instanceID)

for _, entry := range history {
    fmt.Printf("%s: %s -> %s by %s\n",
        entry.Timestamp,
        entry.ActionTaken,
        entry.ToStepID,
        entry.PerformedByName,
    )
}
```

## 🔧 Database Setup

```bash
# Run migrations
psql -U postgres -d hr_system -f migrations/009_create_workflow_tables.sql
psql -U postgres -d hr_system -f migrations/010_seed_leave_approval_workflow.sql
```

## 📊 Query Examples

```sql
-- Get all my pending tasks
SELECT * FROM assigned_tasks
WHERE assigned_to = 'user-id' AND status = 'pending';

-- Get workflow progress
SELECT i.*, w.name as workflow_name, s.step_name as current_step
FROM workflow_instances i
JOIN workflows w ON i.workflow_id = w.id
JOIN workflow_steps s ON i.current_step_id = s.id
WHERE i.id = 'instance-id';

-- Get full audit trail
SELECT * FROM workflow_history
WHERE instance_id = 'instance-id'
ORDER BY timestamp;
```

## 🎯 Leave Request Integration Example

```go
func (s *LeaveRequestService) SubmitLeaveRequest(req *CreateLeaveRequestInput, employeeID string) error {
    // 1. Create leave request
    leaveRequest := &models.LeaveRequest{
        EmployeeID:  employeeID,
        LeaveTypeID: req.LeaveTypeID,
        StartDate:   req.StartDate,
        EndDate:     req.EndDate,
        Reason:      req.Reason,
        Status:      "pending",
    }

    if err := s.leaveRequestRepo.Create(leaveRequest); err != nil {
        return err
    }

    // 2. Start workflow
    employee, _ := s.employeeRepo.GetByID(employeeID)
    leaveType, _ := s.leaveTypeRepo.GetByID(req.LeaveTypeID)

    instance, err := s.workflowService.InitiateWorkflow(
        "Leave Request Approval",
        models.TaskDetails{
            TaskID:   leaveRequest.ID,
            TaskType: "leave_request",
            TaskDescription: fmt.Sprintf(
                "%s: %s to %s (%d days)",
                leaveType.Name,
                req.StartDate.Format("2006-01-02"),
                req.EndDate.Format("2006-01-02"),
                calculateDays(req.StartDate, req.EndDate),
            ),
            SenderDetails: models.SenderDetails{
                SenderID:   employee.ID,
                SenderName: employee.FirstName + " " + employee.LastName,
                Position:   employee.Position.Name,
                Department: employee.Department.Name,
            },
        },
        employee.UserID,
        "medium",
        nil,
    )

    if err != nil {
        return err
    }

    // 3. Link workflow instance to leave request
    leaveRequest.WorkflowInstanceID = &instance.ID
    return s.leaveRequestRepo.Update(leaveRequest)
}
```

## 🛡️ Permission Checking

The service automatically checks permissions:

```go
// This will fail if user doesn't have the right role
err := workflowService.ProcessAction(instanceID, "approve", employeeID, "...")
// Error: "user employee@example.com does not have permission to perform this action"
```

**Step Configuration:**
```sql
-- Only HR Managers and Super Admins can approve at HR Review step
SELECT step_name, allowed_roles FROM workflow_steps WHERE step_name = 'HR Review';

-- Result:
-- HR Review | ["HR_MANAGER", "SUPER_ADMIN"]
```

## 📱 Dashboard Integration

```go
// Count pending tasks for badge
count, err := taskRepo.CountPendingByAssignee(userID)

// Get tasks with instance details
tasks, _ := taskRepo.GetPendingByAssignee(userID)

for _, task := range tasks {
    instance, _ := instanceRepo.GetByID(task.InstanceID)

    fmt.Printf("Task: %s\n", task.StepName)
    fmt.Printf("From: %s\n", instance.TaskDetails.SenderDetails.SenderName)
    fmt.Printf("Description: %s\n", instance.TaskDetails.TaskDescription)
    fmt.Printf("Priority: %s\n", instance.Priority)
}
```

## 🔄 Workflow States

**Instance Status:**
- `pending` - Created but not started (rarely used)
- `in_progress` - Currently being processed
- `completed` - Successfully finished
- `rejected` - Denied (final)
- `cancelled` - Cancelled by initiator or admin

**Task Status:**
- `pending` - Waiting for assignee action
- `in_progress` - Assignee is working on it
- `completed` - Assignee finished
- `skipped` - Skipped (e.g., parallel approval not needed)

## ⚠️ Important Notes

1. **Always check errors** - Workflow operations can fail
2. **Transaction safety** - Consider using DB transactions for critical operations
3. **Idempotency** - Don't process same action twice
4. **Audit trail** - History is automatically maintained
5. **Performance** - Use pagination for task lists in production

## 🐛 Common Issues

**Issue:** "no valid transitions found"
```go
// Solution: Check workflow has transitions defined
transitions, _ := workflowRepo.GetValidTransitions(currentStepID)
```

**Issue:** "user does not have permission"
```go
// Solution: Verify user's role is in step's allowed_roles
step, _ := workflowRepo.GetStepByID(currentStepID)
fmt.Println(step.AllowedRoles) // Check if user's role is here
```

**Issue:** "action 'xyz' is not valid from current step"
```go
// Solution: Check available actions
transitions, _ := workflowRepo.GetValidTransitions(currentStepID)
for _, tr := range transitions {
    fmt.Println("Available action:", tr.ActionName)
}
```

## 📚 Further Reading

- [Workflow System Documentation](./workflow-system.md) - Detailed architecture
- [Implementation Summary](./workflow-implementation-summary.md) - What's built
- [Models](../internal/models/workflow.go) - Data structures
- [Service](../internal/services/workflow_service.go) - Business logic

## 🎓 Pro Tips

1. **Workflow names are unique** - Use descriptive names
2. **Task metadata** - Use `TaskDetails.Metadata` for custom data
3. **History comments** - Provide meaningful comments for audit
4. **Due dates** - Set reasonable deadlines for SLA tracking
5. **Priority levels** - Use consistently across your app

---

**Need Help?** Check the detailed documentation or review the seeded "Leave Request Approval" workflow as a reference implementation.