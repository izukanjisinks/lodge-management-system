# Workflow API - Postman Quick Reference

This guide provides quick examples for testing the workflow system using the updated Postman collection.

## Prerequisites

1. **Import Collection**: Import `postman/hr-system.postman_collection.json` into Postman
2. **Login**: Run the "Login" request first - it automatically sets `{{token}}` variable
3. **Get Workflow ID**: Run "Get All Workflows" to get the Leave Request Approval workflow ID

## Testing Workflow Flow

### Step 1: Get the Seeded Workflow

**Request:** `GET /api/v1/admin/workflows`

**Expected Response:**
```json
{
  "count": 1,
  "workflows": [
    {
      "id": "132ddcc9-7385-4c37-9ad1-cba4a939c7b5",
      "name": "Leave Request Approval",
      "description": "Standard workflow for employee leave request approvals...",
      "is_active": true
    }
  ]
}
```

Save the workflow `id` for later use.

---

### Step 2: View Workflow Structure

**Request:** `GET /api/v1/admin/workflows/{workflow_id}/structure`

**What to look for:**
- 4 steps: Submit → HR Review → Manager Approval → Completed
- 5 transitions defining allowed actions
- Allowed roles for each step

---

### Step 3: Initiate a Workflow Instance

**Request:** `POST /api/v1/workflow/instances`

**Body:**
```json
{
  "workflow_name": "Leave Request Approval",
  "task_details": {
    "task_id": "leave-req-001",
    "task_type": "leave_request",
    "task_description": "Annual Leave: 5 days (2026-03-10 to 2026-03-14)",
    "sender_details": {
      "sender_id": "b94804ea-b107-47ca-823f-235ec6b6f1d1",
      "sender_name": "Admin User",
      "position": "System Administrator",
      "department": "IT"
    }
  },
  "priority": "medium"
}
```

**Expected Response:**
```json
{
  "message": "Workflow initiated successfully",
  "instance": {
    "id": "instance-uuid-here",
    "workflow_id": "workflow-uuid",
    "current_step_id": "hr-review-step-id",
    "status": "in_progress",
    "created_at": "2026-02-22T..."
  }
}
```

Save the `instance.id` for the next steps.

---

### Step 4: Check Pending Tasks

**Request:** `GET /api/v1/workflow/my-tasks/pending`

**Expected Response:**
```json
{
  "tasks": [
    {
      "id": "task-uuid",
      "instance_id": "instance-uuid",
      "step_name": "HR Review",
      "status": "pending",
      "assigned_to": "current-user-uuid"
    }
  ],
  "count": 1
}
```

---

### Step 5: Process an Action (Approve)

**Request:** `POST /api/v1/workflow/instances/{instance_id}/action`

**Body:**
```json
{
  "action": "approve",
  "comments": "Leave request looks good. Approved for HR step."
}
```

**Expected Response:**
```json
{
  "message": "Action processed successfully",
  "next_step": "Manager Approval"
}
```

The workflow will now move to the Manager Approval step.

---

### Step 6: View Workflow History

**Request:** `GET /api/v1/workflow/instances/{instance_id}/history`

**Expected Response:**
```json
{
  "history": [
    {
      "id": "history-uuid-1",
      "from_step_id": "submit-step-id",
      "to_step_id": "hr-review-step-id",
      "action_taken": "submit",
      "performed_by_name": "Admin User",
      "comments": "Initiated Leave Request Approval workflow",
      "timestamp": "2026-02-22T..."
    },
    {
      "id": "history-uuid-2",
      "from_step_id": "hr-review-step-id",
      "to_step_id": "manager-approval-step-id",
      "action_taken": "approve",
      "performed_by_name": "Admin User",
      "comments": "Leave request looks good. Approved for HR step.",
      "timestamp": "2026-02-22T..."
    }
  ],
  "count": 2
}
```

---

## Common Actions

### Approve a Task
```json
{
  "action": "approve",
  "comments": "Approved. Everything looks good."
}
```

### Reject a Task
```json
{
  "action": "reject",
  "comments": "Insufficient leave balance. Please revise."
}
```

---

## Workflow States

| Status | Description |
|--------|-------------|
| `pending` | Workflow created but not yet started |
| `in_progress` | Workflow is actively moving through steps |
| `completed` | Workflow reached the final step |
| `rejected` | Workflow was rejected |
| `cancelled` | Workflow was cancelled |

---

## Task States

| Status | Description |
|--------|-------------|
| `pending` | Task assigned but not yet started |
| `in_progress` | Task is being worked on |
| `completed` | Task has been completed |
| `skipped` | Task was skipped |

---

## Available Actions

Depends on the workflow configuration, but typically:

- `submit` - Initial submission of workflow
- `approve` - Approve current step
- `reject` - Reject and send back
- `reassign` - Reassign to another user

The valid actions for each step are defined in the workflow transitions.

---

## Admin Workflow Management

### Create a Custom Workflow

1. **Create Workflow Template:**
   ```
   POST /api/v1/admin/workflows
   {
     "name": "Equipment Request",
     "description": "Workflow for equipment procurement"
   }
   ```

2. **Create Steps:**
   ```
   POST /api/v1/admin/workflow-steps
   {
     "workflow_id": "workflow-uuid",
     "step_name": "Submit",
     "step_order": 1,
     "initial": true,
     "allowed_roles": ["employee"]
   }
   ```

3. **Create Transitions:**
   ```
   POST /api/v1/admin/workflow-transitions
   {
     "workflow_id": "workflow-uuid",
     "from_step_id": "submit-step-id",
     "to_step_id": "manager-review-step-id",
     "action_name": "submit"
   }
   ```

---

## Troubleshooting

### "Forbidden: role required"
- Make sure you're logged in as `super_admin` or `hr_manager` for admin endpoints
- Check that the `{{token}}` variable is set (run Login first)

### "Workflow not found"
- Verify migrations were run: `make migrate-up`
- Check workflow exists: `GET /api/v1/admin/workflows`

### "No suitable assignee found"
- The workflow step configuration requires specific roles
- Current logic assigns to users with matching roles
- For testing, super_admin can perform all actions

### "Action not valid from current step"
- Check valid transitions: `GET /api/v1/admin/workflows/{id}/structure`
- The action must be defined as a transition from the current step

---

## Next Steps

1. **Integrate with Leave Requests**: Update LeaveRequestService to initiate workflows
2. **Frontend Integration**: Build UI for pending tasks and workflow actions
3. **Notifications**: Add email/push notifications when tasks are assigned
4. **Advanced Assignment**: Implement department-based and manager-based task assignment
5. **Workflow Builder**: Create UI for designing custom workflows

---

## Related Documentation

- [Workflow System Architecture](./workflow-system.md)
- [Workflow API Documentation](./workflow-api.md)
- [Admin API Documentation](./workflow-admin-api.md)
- [Implementation Checklist](./workflow-implementation-checklist.md)
