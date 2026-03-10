# Workflow API Documentation

## Overview

The Workflow API provides endpoints for managing workflow instances, processing approvals, and tracking task assignments.

**Base URL:** `/api/v1/workflow`

**Authentication:** All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

---

## Endpoints

### 1. Get My Tasks

Retrieve all tasks assigned to the authenticated user.

**Endpoint:** `GET /api/v1/workflow/my-tasks`

**Query Parameters:**
- `status` (optional): Filter by status (`pending`, `completed`, `in_progress`, `skipped`)

**Example Request:**
```bash
# Get all tasks
curl -X GET http://localhost:8081/api/v1/workflow/my-tasks \
  -H "Authorization: Bearer <token>"

# Get only pending tasks
curl -X GET http://localhost:8081/api/v1/workflow/my-tasks?status=pending \
  -H "Authorization: Bearer <token>"
```

**Response:** `200 OK`
```json
{
  "tasks": [
    {
      "id": "task-uuid-1",
      "instance_id": "instance-uuid-1",
      "step_id": "step-uuid-1",
      "step_name": "HR Review",
      "assigned_to": "user-uuid-1",
      "assigned_by": "user-uuid-2",
      "status": "pending",
      "due_date": "2026-02-28T23:59:59Z",
      "completed_at": null,
      "created_at": "2026-02-20T10:00:00Z",
      "updated_at": "2026-02-20T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

### 2. Get My Pending Tasks

Retrieve only pending tasks for the authenticated user.

**Endpoint:** `GET /api/v1/workflow/my-tasks/pending`

**Example Request:**
```bash
curl -X GET http://localhost:8081/api/v1/workflow/my-tasks/pending \
  -H "Authorization: Bearer <token>"
```

**Response:** `200 OK` (same format as Get My Tasks)

---

### 3. Get Task Details

Get detailed information about a specific task including workflow instance and history.

**Endpoint:** `GET /api/v1/workflow/tasks/{id}`

**Path Parameters:**
- `id`: Task UUID

**Example Request:**
```bash
curl -X GET http://localhost:8081/api/v1/workflow/tasks/task-uuid-1 \
  -H "Authorization: Bearer <token>"
```

**Response:** `200 OK`
```json
{
  "task": {
    "id": "task-uuid-1",
    "instance_id": "instance-uuid-1",
    "step_id": "step-uuid-1",
    "step_name": "HR Review",
    "assigned_to": "user-uuid-1",
    "assigned_by": "user-uuid-2",
    "status": "pending",
    "due_date": "2026-02-28T23:59:59Z",
    "created_at": "2026-02-20T10:00:00Z"
  },
  "instance": {
    "id": "instance-uuid-1",
    "workflow_id": "workflow-uuid-1",
    "current_step_id": "step-uuid-1",
    "status": "in_progress",
    "task_details": {
      "task_id": "leave-request-123",
      "task_type": "leave_request",
      "task_description": "Annual Leave: 5 days (2026-03-01 to 2026-03-05)",
      "sender_details": {
        "sender_id": "emp-uuid-1",
        "sender_name": "John Doe",
        "position": "Software Engineer",
        "department": "Engineering"
      }
    },
    "created_by": "user-uuid-2",
    "priority": "medium",
    "created_at": "2026-02-20T10:00:00Z"
  },
  "history": [
    {
      "id": "history-uuid-1",
      "instance_id": "instance-uuid-1",
      "from_step_id": "step-uuid-0",
      "to_step_id": "step-uuid-1",
      "action_taken": "submit",
      "performed_by": "user-uuid-2",
      "performed_by_name": "John Doe",
      "comments": "Initiated Leave Request Approval workflow",
      "timestamp": "2026-02-20T10:00:00Z"
    }
  ]
}
```

---

### 4. Process Action

Process an action on a workflow instance (approve, reject, etc.).

**Endpoint:** `POST /api/v1/workflow/instances/{id}/action`

**Path Parameters:**
- `id`: Workflow Instance UUID

**Request Body:**
```json
{
  "action": "approve",
  "comments": "Leave balance verified. Approved for manager review."
}
```

**Fields:**
- `action` (required): Action to perform (`approve`, `reject`, `submit`)
- `comments` (optional): Additional comments for the action

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/workflow/instances/instance-uuid-1/action \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "comments": "Leave balance verified. Approved for manager review."
  }'
```

**Response:** `200 OK`
```json
{
  "message": "Action processed successfully",
  "action": "approve"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid action or missing required fields
```json
{
  "error": "action 'invalid' is not valid from current step"
}
```

- `403 Forbidden`: User doesn't have permission
```json
{
  "error": "Permission denied"
}
```

- `409 Conflict`: Workflow already closed
```json
{
  "error": "workflow instance is already closed"
}
```

---

### 5. Get Instance History

Retrieve the complete audit trail for a workflow instance.

**Endpoint:** `GET /api/v1/workflow/instances/{id}/history`

**Path Parameters:**
- `id`: Workflow Instance UUID

**Example Request:**
```bash
curl -X GET http://localhost:8081/api/v1/workflow/instances/instance-uuid-1/history \
  -H "Authorization: Bearer <token>"
```

**Response:** `200 OK`
```json
{
  "instance_id": "instance-uuid-1",
  "history": [
    {
      "id": "history-uuid-1",
      "instance_id": "instance-uuid-1",
      "from_step_id": null,
      "to_step_id": "step-uuid-hr",
      "action_taken": "submit",
      "performed_by": "user-uuid-employee",
      "performed_by_name": "John Doe",
      "comments": "Initiated Leave Request Approval workflow",
      "timestamp": "2026-02-20T10:00:00Z"
    },
    {
      "id": "history-uuid-2",
      "instance_id": "instance-uuid-1",
      "from_step_id": "step-uuid-hr",
      "to_step_id": "step-uuid-manager",
      "action_taken": "approve",
      "performed_by": "user-uuid-hr",
      "performed_by_name": "Jane Smith (HR Manager)",
      "comments": "Leave balance verified. Approved for manager review.",
      "timestamp": "2026-02-20T14:30:00Z"
    },
    {
      "id": "history-uuid-3",
      "instance_id": "instance-uuid-1",
      "from_step_id": "step-uuid-manager",
      "to_step_id": "step-uuid-completed",
      "action_taken": "approve",
      "performed_by": "user-uuid-manager",
      "performed_by_name": "Bob Johnson (Dept Head)",
      "comments": "Team coverage confirmed. Approved.",
      "timestamp": "2026-02-21T09:15:00Z"
    }
  ],
  "count": 3
}
```

---

### 6. Get Instance by Task ID

Retrieve a workflow instance using the associated task ID (e.g., leave request ID).

**Endpoint:** `GET /api/v1/workflow/instances/by-task/{task_id}`

**Path Parameters:**
- `task_id`: The ID from `task_details.task_id` (e.g., leave request ID)

**Example Request:**
```bash
curl -X GET http://localhost:8081/api/v1/workflow/instances/by-task/leave-request-123 \
  -H "Authorization: Bearer <token>"
```

**Response:** `200 OK`
```json
{
  "id": "instance-uuid-1",
  "workflow_id": "workflow-uuid-1",
  "current_step_id": "step-uuid-manager",
  "status": "in_progress",
  "task_details": {
    "task_id": "leave-request-123",
    "task_type": "leave_request",
    "task_description": "Annual Leave: 5 days",
    "sender_details": {
      "sender_id": "emp-uuid-1",
      "sender_name": "John Doe",
      "position": "Software Engineer",
      "department": "Engineering"
    }
  },
  "created_by": "user-uuid-2",
  "created_at": "2026-02-20T10:00:00Z",
  "updated_at": "2026-02-20T14:30:00Z",
  "completed_at": null,
  "due_date": "2026-02-28T23:59:59Z",
  "priority": "medium"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "Workflow instance not found"
}
```

---

### 7. Initiate Workflow

Start a new workflow instance. This is typically called by other services (e.g., when creating a leave request), not directly by end users.

**Endpoint:** `POST /api/v1/workflow/instances`

**Request Body:**
```json
{
  "workflow_name": "Leave Request Approval",
  "task_details": {
    "task_id": "leave-request-123",
    "task_type": "leave_request",
    "task_description": "Annual Leave: 5 days (2026-03-01 to 2026-03-05)",
    "sender_details": {
      "sender_id": "emp-uuid-1",
      "sender_name": "John Doe",
      "position": "Software Engineer",
      "department": "Engineering"
    },
    "metadata": "{\"leave_type\":\"annual\",\"days\":5}"
  },
  "priority": "medium",
  "due_date": "2026-02-28T23:59:59Z"
}
```

**Fields:**
- `workflow_name` (required): Name of the workflow template
- `task_details` (required): Details about the task
  - `task_id` (required): ID of the related entity
  - `task_type` (required): Type of task (e.g., "leave_request")
  - `task_description` (required): Human-readable description
  - `sender_details` (required): Information about who initiated
  - `metadata` (optional): Additional JSON data
- `priority` (optional): Priority level (default: "medium")
- `due_date` (optional): Deadline for completion

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/workflow/instances \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "workflow_name": "Leave Request Approval",
    "task_details": {
      "task_id": "leave-request-123",
      "task_type": "leave_request",
      "task_description": "Annual Leave: 5 days",
      "sender_details": {
        "sender_id": "emp-uuid-1",
        "sender_name": "John Doe",
        "position": "Software Engineer",
        "department": "Engineering"
      }
    },
    "priority": "medium"
  }'
```

**Response:** `201 Created`
```json
{
  "message": "Workflow initiated successfully",
  "instance": {
    "id": "instance-uuid-1",
    "workflow_id": "workflow-uuid-1",
    "current_step_id": "step-uuid-hr",
    "status": "in_progress",
    "task_details": { /* ... */ },
    "created_by": "user-uuid-2",
    "created_at": "2026-02-20T10:00:00Z",
    "priority": "medium"
  }
}
```

---

## Status Values

### Instance Status
- `pending`: Created but not started
- `in_progress`: Currently being processed
- `completed`: Successfully finished
- `rejected`: Denied
- `cancelled`: Cancelled

### Task Status
- `pending`: Waiting for assignee
- `in_progress`: Being worked on
- `completed`: Finished
- `skipped`: Not needed

### Priority Levels
- `low`: Low priority
- `medium`: Normal priority (default)
- `high`: High priority
- `urgent`: Urgent priority

---

## Common Actions

The available actions depend on the workflow configuration. Common actions include:

- `submit`: Submit for review
- `approve`: Approve and move to next step
- `reject`: Reject and send back
- `reassign`: Assign to different user
- `cancel`: Cancel the workflow

---

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input or action |
| 403 | Forbidden - Permission denied |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Workflow already closed |
| 500 | Internal Server Error |

---

## Integration Example

### Leave Request Submission

When an employee submits a leave request:

```javascript
// 1. Create leave request in your system
const leaveRequest = await createLeaveRequest({
  employee_id: employeeId,
  leave_type_id: leaveTypeId,
  start_date: startDate,
  end_date: endDate,
  reason: reason
});

// 2. Initiate workflow
const workflow = await fetch('/api/v1/workflow/instances', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    workflow_name: 'Leave Request Approval',
    task_details: {
      task_id: leaveRequest.id,
      task_type: 'leave_request',
      task_description: `${leaveType}: ${days} days`,
      sender_details: {
        sender_id: employee.id,
        sender_name: `${employee.first_name} ${employee.last_name}`,
        position: employee.position,
        department: employee.department
      }
    },
    priority: 'medium'
  })
});
```

### HR Manager Approval

```javascript
await fetch(`/api/v1/workflow/instances/${instanceId}/action`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    action: 'approve',
    comments: 'Leave balance verified. Approved for manager review.'
  })
});
```

---

## Best Practices

1. **Always check task ownership** - Users should only act on tasks assigned to them
2. **Provide meaningful comments** - Help maintain clear audit trail
3. **Handle errors gracefully** - Check for permission and validation errors
4. **Poll for updates** - Check task status periodically or use webhooks
5. **Display history** - Show users the full approval chain

---

**For more details, see:**
- [Workflow System Documentation](./workflow-system.md)
- [Quick Start Guide](./workflow-quick-start.md)
- [Implementation Summary](./workflow-implementation-summary.md)