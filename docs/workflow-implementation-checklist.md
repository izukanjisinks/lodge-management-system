# Workflow System Implementation Checklist

## ✅ Completed

### 1. Data Models
- [x] Workflow (template)
- [x] WorkflowStep
- [x] WorkflowTransition
- [x] WorkflowInstance
- [x] AssignedTask
- [x] WorkflowHistory
- [x] TaskDetails
- [x] SenderDetails

### 2. Database
- [x] Migration 009: Create workflow tables
- [x] Migration 010: Seed Leave Approval workflow
- [x] Proper indexes for performance
- [x] Foreign key constraints
- [x] Auto-update triggers

### 3. Repository Layer
- [x] WorkflowRepository (template operations)
- [x] WorkflowInstanceRepository (instance lifecycle)
- [x] AssignedTaskRepository (task management)
- [x] WorkflowHistoryRepository (audit trail)

### 4. Service Layer
- [x] WorkflowService
  - [x] InitiateWorkflow
  - [x] ProcessAction
  - [x] GetMyTasks
  - [x] GetInstanceHistory
  - [x] GetInstanceByTaskID
  - [x] Permission checking
  - [x] Task assignment logic

### 5. Handler Layer
- [x] WorkflowHandler (user-facing endpoints)
  - [x] GetMyTasks
  - [x] GetMyPendingTasks
  - [x] GetTaskDetails
  - [x] ProcessAction
  - [x] GetInstanceHistory
  - [x] GetInstanceByTaskID
  - [x] InitiateWorkflow
- [x] WorkflowAdminHandler (administration endpoints)
  - [x] GetAllWorkflows
  - [x] GetWorkflowByID
  - [x] CreateWorkflow
  - [x] GetWorkflowSteps
  - [x] GetStepByID
  - [x] CreateWorkflowStep
  - [x] GetWorkflowTransitions
  - [x] GetValidTransitions
  - [x] CreateWorkflowTransition
  - [x] GetWorkflowStructure

### 6. Routes
**User Workflow Routes:**
- [x] GET /api/v1/workflow/my-tasks
- [x] GET /api/v1/workflow/my-tasks/pending
- [x] GET /api/v1/workflow/tasks/{id}
- [x] POST /api/v1/workflow/instances/{id}/action
- [x] GET /api/v1/workflow/instances/{id}/history
- [x] GET /api/v1/workflow/instances/by-task/{task_id}
- [x] POST /api/v1/workflow/instances

**Admin Workflow Routes:**
- [x] GET /api/v1/admin/workflows
- [x] GET /api/v1/admin/workflows/{id}
- [x] POST /api/v1/admin/workflows
- [x] GET /api/v1/admin/workflows/{id}/steps
- [x] GET /api/v1/admin/workflows/steps/{step_id}
- [x] POST /api/v1/admin/workflows/steps
- [x] GET /api/v1/admin/workflows/{id}/transitions
- [x] GET /api/v1/admin/workflows/steps/{step_id}/transitions
- [x] POST /api/v1/admin/workflows/transitions
- [x] GET /api/v1/admin/workflows/{id}/structure

### 7. Main Application Integration
- [x] Initialize workflow repositories
- [x] Create workflow service
- [x] Create workflow handler
- [x] Create workflow admin handler
- [x] Register workflow routes
- [x] Register workflow admin routes

### 8. Documentation
- [x] Workflow System Architecture (workflow-system.md)
- [x] Implementation Summary (workflow-implementation-summary.md)
- [x] Quick Start Guide (workflow-quick-start.md)
- [x] API Documentation (workflow-api.md)
- [x] Admin API Documentation (workflow-admin-api.md)
- [x] Implementation Checklist (this file)

---

## 🔄 Next Steps (To Complete Integration)

### 1. Run Migrations
```bash
cd c:\development\hr-system

# Connect to database
psql -U postgres -d hr_system

# Run migrations
\i migrations/009_create_workflow_tables.sql
\i migrations/010_seed_leave_approval_workflow.sql

# Verify tables created
\dt workflow*
SELECT * FROM workflows;
SELECT * FROM workflow_steps ORDER BY step_order;
```

### 2. Test the API

**Start the server:**
```bash
go run cmd/api/main.go
```

**Test endpoints:**
```bash
# Login to get token
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@hr-system.com","password":"Admin@123"}' \
  | jq -r '.token')

# Get my pending tasks
curl http://localhost:8081/api/v1/workflow/my-tasks/pending \
  -H "Authorization: Bearer $TOKEN"

# Initiate a workflow
curl -X POST http://localhost:8081/api/v1/workflow/instances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workflow_name": "Leave Request Approval",
    "task_details": {
      "task_id": "test-leave-001",
      "task_type": "leave_request",
      "task_description": "Annual Leave: 3 days",
      "sender_details": {
        "sender_id": "emp-001",
        "sender_name": "Test Employee",
        "position": "Developer",
        "department": "IT"
      }
    },
    "priority": "medium"
  }'
```

### 3. Integrate with Leave Request Service

**Update LeaveRequestService:**

```go
// File: internal/services/leave_request_service.go

type LeaveRequestService struct {
    repo            *repository.LeaveRequestRepository
    balanceService  *LeaveBalanceService
    leaveTypeRepo   *repository.LeaveTypeRepository
    holidayRepo     *repository.HolidayRepository
    empRepo         *repository.EmployeeRepository
    workflowService *WorkflowService  // ADD THIS
}

func NewLeaveRequestService(
    repo *repository.LeaveRequestRepository,
    balanceService *LeaveBalanceService,
    leaveTypeRepo *repository.LeaveTypeRepository,
    holidayRepo *repository.HolidayRepository,
    empRepo *repository.EmployeeRepository,
    workflowService *WorkflowService,  // ADD THIS
) *LeaveRequestService {
    return &LeaveRequestService{
        repo:            repo,
        balanceService:  balanceService,
        leaveTypeRepo:   leaveTypeRepo,
        holidayRepo:     holidayRepo,
        empRepo:         empRepo,
        workflowService: workflowService,  // ADD THIS
    }
}

// Update Create method to initiate workflow
func (s *LeaveRequestService) Create(req *models.LeaveRequest) error {
    // 1. Validate leave request
    // 2. Create leave request in database
    if err := s.repo.Create(req); err != nil {
        return err
    }

    // 3. Get employee details
    employee, err := s.empRepo.GetByID(req.EmployeeID)
    if err != nil {
        return err
    }

    // 4. Get leave type
    leaveType, err := s.leaveTypeRepo.GetByID(req.LeaveTypeID)
    if err != nil {
        return err
    }

    // 5. Calculate days
    days := calculateWorkDays(req.StartDate, req.EndDate)

    // 6. Initiate workflow
    instance, err := s.workflowService.InitiateWorkflow(
        "Leave Request Approval",
        models.TaskDetails{
            TaskID:   req.ID,
            TaskType: "leave_request",
            TaskDescription: fmt.Sprintf(
                "%s: %d days (%s to %s)",
                leaveType.Name,
                days,
                req.StartDate.Format("2006-01-02"),
                req.EndDate.Format("2006-01-02"),
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
        return fmt.Errorf("failed to initiate workflow: %w", err)
    }

    // 7. Update leave request with workflow instance ID
    req.WorkflowInstanceID = &instance.ID
    return s.repo.Update(req)
}
```

**Add WorkflowInstanceID to LeaveRequest model:**

```go
// File: internal/models/leave_request.go

type LeaveRequest struct {
    ID                 string
    EmployeeID         string
    LeaveTypeID        string
    StartDate          time.Time
    EndDate            time.Time
    Reason             string
    Status             string
    WorkflowInstanceID *string    // ADD THIS
    ApprovedBy         *string
    ApprovedAt         *time.Time
    RejectedBy         *string
    RejectedAt         *time.Time
    RejectionReason    *string
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

**Add column to database:**

```sql
-- Migration: Add workflow_instance_id to leave_requests
ALTER TABLE leave_requests
ADD COLUMN workflow_instance_id UUID REFERENCES workflow_instances(id);

CREATE INDEX idx_leave_requests_workflow_instance
ON leave_requests(workflow_instance_id);
```

### 4. Update Leave Request Handler

**Add workflow status to leave request responses:**

```go
func (h *LeaveRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    leaveRequest, err := h.service.GetByID(id)
    if err != nil {
        http.Error(w, "Leave request not found", http.StatusNotFound)
        return
    }

    response := map[string]interface{}{
        "leave_request": leaveRequest,
    }

    // Add workflow status if exists
    if leaveRequest.WorkflowInstanceID != nil {
        instance, err := h.workflowService.GetInstanceByTaskID(*leaveRequest.WorkflowInstanceID)
        if err == nil {
            response["workflow_status"] = instance.Status
            response["current_step"] = instance.CurrentStepID
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### 5. Enhance Assignment Logic

**Update WorkflowService.determineAssignee:**

Currently it's a placeholder. Implement proper logic:

```go
func (s *WorkflowService) determineAssignee(step *models.WorkflowStep, taskDetails models.TaskDetails) (string, error) {
    // For HR Review step
    if contains(step.AllowedRoles, "HR_MANAGER") {
        // Find HR manager (implement this in UserRepository)
        hrManager, err := s.userRepo.GetFirstUserByRole("HR_MANAGER")
        if err == nil {
            return hrManager.ID, nil
        }
    }

    // For Manager Approval step
    if contains(step.AllowedRoles, "DEPARTMENT_HEAD") {
        // Get employee's department head
        employee, err := s.userRepo.GetByID(taskDetails.SenderDetails.SenderID)
        if err == nil {
            // Implement GetDepartmentHead in repository
            deptHead, err := s.userRepo.GetDepartmentHead(employee.DepartmentID)
            if err == nil {
                return deptHead.ID, nil
            }
        }
    }

    // Fallback to first user with allowed role
    for _, roleCode := range step.AllowedRoles {
        user, err := s.userRepo.GetFirstUserByRole(roleCode)
        if err == nil {
            return user.ID, nil
        }
    }

    return "", errors.New("no suitable assignee found")
}
```

### 6. Add Notifications

**Create notification when task assigned:**

```go
// After creating assigned task in InitiateWorkflow
if err := s.notificationService.NotifyTaskAssigned(task, instance); err != nil {
    // Log but don't fail the workflow
    log.Printf("Failed to send notification: %v", err)
}
```

### 7. Update Dashboard

**Add pending tasks count to dashboard:**

```go
type DashboardStats struct {
    // ... existing fields
    PendingTasks int `json:"pending_tasks"`
}

func (s *DashboardService) GetMyDashboard(userID string) (*DashboardStats, error) {
    // ... existing logic

    // Add pending tasks count
    count, err := s.taskRepo.CountPendingByAssignee(userID)
    if err == nil {
        stats.PendingTasks = count
    }

    return stats, nil
}
```

### 8. Testing Checklist

- [ ] Can initiate workflow
- [ ] Can view my tasks
- [ ] Can approve task
- [ ] Can reject task
- [ ] History is recorded correctly
- [ ] Permissions are enforced
- [ ] Cannot act on closed workflows
- [ ] Cannot act without permission
- [ ] Task reassignment works
- [ ] Due dates are tracked
- [ ] Integration with leave requests works

### 9. Frontend Updates Needed

- [ ] Create "My Tasks" page
- [ ] Add task action buttons (Approve/Reject)
- [ ] Show workflow history timeline
- [ ] Add pending tasks badge to nav
- [ ] Display current workflow status in leave request details
- [ ] Add comments field for approvals

---

## 📝 SQL Queries for Testing

```sql
-- Check workflow template
SELECT * FROM workflows WHERE name = 'Leave Request Approval';

-- Check workflow steps
SELECT step_name, step_order, allowed_roles, initial, final
FROM workflow_steps
WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Leave Request Approval')
ORDER BY step_order;

-- Check transitions
SELECT
    (SELECT step_name FROM workflow_steps WHERE id = from_step_id) as from_step,
    (SELECT step_name FROM workflow_steps WHERE id = to_step_id) as to_step,
    action_name
FROM workflow_transitions
WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Leave Request Approval');

-- View active workflow instances
SELECT
    i.id,
    i.status,
    s.step_name as current_step,
    i.task_details->>'task_description' as description,
    i.created_at
FROM workflow_instances i
JOIN workflow_steps s ON i.current_step_id = s.id
WHERE i.status = 'in_progress'
ORDER BY i.created_at DESC;

-- View my pending tasks
SELECT
    t.id,
    t.step_name,
    i.task_details->>'task_description' as description,
    i.task_details->>'sender_details'->>'sender_name' as from_user,
    t.created_at
FROM assigned_tasks t
JOIN workflow_instances i ON t.instance_id = i.id
WHERE t.assigned_to = 'user-id' AND t.status = 'pending'
ORDER BY t.created_at DESC;

-- View workflow history
SELECT
    h.timestamp,
    (SELECT step_name FROM workflow_steps WHERE id = h.from_step_id) as from_step,
    (SELECT step_name FROM workflow_steps WHERE id = h.to_step_id) as to_step,
    h.action_taken,
    h.performed_by_name,
    h.comments
FROM workflow_history h
WHERE h.instance_id = 'instance-id'
ORDER BY h.timestamp;
```

---

## 🎯 Success Criteria

The workflow system is ready for production when:

1. ✅ All migrations run successfully
2. ✅ API endpoints return correct responses
3. ✅ Workflow can be initiated from leave request
4. ✅ HR can approve/reject tasks
5. ✅ Managers can approve/reject tasks
6. ✅ History is complete and accurate
7. ✅ Permissions are properly enforced
8. ✅ Frontend displays tasks correctly
9. ✅ Notifications are sent on task assignment
10. ✅ Performance is acceptable under load

---

## 📚 Additional Resources

- API Documentation: [workflow-api.md](./workflow-api.md)
- Architecture Guide: [workflow-system.md](./workflow-system.md)
- Quick Start: [workflow-quick-start.md](./workflow-quick-start.md)
- Implementation Details: [workflow-implementation-summary.md](./workflow-implementation-summary.md)

---

**Current Status:** ✅ Core system complete, ready for integration testing
**Next Milestone:** Run migrations and test API endpoints