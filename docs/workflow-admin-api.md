# Workflow Admin API Documentation

This document describes the administrative API endpoints for managing workflow templates, steps, and transitions.

## Access Control

All admin endpoints require authentication and the `hr_manager` role or higher.

**Base URL:** `/api/v1/admin/workflows`

---

## Workflow Template Management

### 1. Get All Workflows

**Endpoint:** `GET /api/v1/admin/workflows`

**Description:** Retrieves all active workflow templates in the system.

**Response:**
```json
{
  "workflows": [
    {
      "id": "uuid",
      "name": "Leave Request Approval",
      "description": "Standard leave approval workflow",
      "is_active": true,
      "created_by": "user-id",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

### 2. Get Workflow by ID

**Endpoint:** `GET /api/v1/admin/workflows/{id}`

**Description:** Retrieves a specific workflow template by ID.

**Response:**
```json
{
  "id": "uuid",
  "name": "Leave Request Approval",
  "description": "Standard leave approval workflow",
  "is_active": true,
  "created_by": "user-id",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

### 3. Create Workflow

**Endpoint:** `POST /api/v1/admin/workflows`

**Description:** Creates a new workflow template.

**Request Body:**
```json
{
  "name": "Equipment Request Approval",
  "description": "Workflow for equipment procurement requests"
}
```

**Response:**
```json
{
  "message": "Workflow created successfully",
  "workflow": {
    "id": "uuid",
    "name": "Equipment Request Approval",
    "description": "Workflow for equipment procurement requests",
    "is_active": true,
    "created_by": "user-id",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

---

## Workflow Steps Management

### 4. Get Workflow Steps

**Endpoint:** `GET /api/v1/admin/workflows/{id}/steps`

**Description:** Retrieves all steps for a specific workflow.

**Response:**
```json
{
  "steps": [
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "step_name": "Manager Approval",
      "step_order": 2,
      "initial": false,
      "final": false,
      "allowed_roles": ["manager"],
      "requires_all_approvers": false,
      "min_approvals": 1
    }
  ],
  "count": 3
}
```

---

### 5. Get Step by ID

**Endpoint:** `GET /api/v1/admin/workflows/steps/{step_id}`

**Description:** Retrieves a specific workflow step.

**Response:**
```json
{
  "id": "uuid",
  "workflow_id": "uuid",
  "step_name": "Manager Approval",
  "step_order": 2,
  "initial": false,
  "final": false,
  "allowed_roles": ["manager"],
  "requires_all_approvers": false,
  "min_approvals": 1
}
```

---

### 6. Create Workflow Step

**Endpoint:** `POST /api/v1/admin/workflows/steps`

**Description:** Creates a new step in a workflow.

**Request Body:**
```json
{
  "workflow_id": "uuid",
  "step_name": "CFO Approval",
  "step_order": 3,
  "initial": false,
  "final": false,
  "allowed_roles": ["cfo", "finance_manager"],
  "requires_all_approvers": false,
  "min_approvals": 1
}
```

**Response:**
```json
{
  "message": "Step creation endpoint - implement CreateStep in repository",
  "request": {
    "workflow_id": "uuid",
    "step_name": "CFO Approval",
    "step_order": 3,
    "initial": false,
    "final": false,
    "allowed_roles": ["cfo", "finance_manager"],
    "requires_all_approvers": false,
    "min_approvals": 1
  }
}
```

**Note:** The `CreateStep` method needs to be implemented in `WorkflowRepository`.

---

## Workflow Transitions Management

### 7. Get Workflow Transitions

**Endpoint:** `GET /api/v1/admin/workflows/{id}/transitions`

**Description:** Retrieves all transitions for a specific workflow.

**Response:**
```json
{
  "transitions": [
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "from_step_id": "uuid",
      "to_step_id": "uuid",
      "action_name": "approve",
      "condition_type": "none",
      "condition_value": null
    }
  ],
  "count": 5
}
```

---

### 8. Get Valid Transitions from Step

**Endpoint:** `GET /api/v1/admin/workflows/steps/{step_id}/transitions`

**Description:** Retrieves all valid transitions from a specific step.

**Response:**
```json
{
  "transitions": [
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "from_step_id": "uuid",
      "to_step_id": "uuid",
      "action_name": "approve",
      "condition_type": "none",
      "condition_value": null
    },
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "from_step_id": "uuid",
      "to_step_id": "uuid",
      "action_name": "reject",
      "condition_type": "none",
      "condition_value": null
    }
  ],
  "count": 2
}
```

---

### 9. Create Workflow Transition

**Endpoint:** `POST /api/v1/admin/workflows/transitions`

**Description:** Creates a new transition between workflow steps.

**Request Body:**
```json
{
  "workflow_id": "uuid",
  "from_step_id": "uuid",
  "to_step_id": "uuid",
  "action_name": "approve",
  "condition_type": "none",
  "condition_value": ""
}
```

**Condition Types:**
- `none` - No conditions
- `role` - Requires specific role
- `custom` - Custom condition logic

**Response:**
```json
{
  "message": "Transition creation endpoint - implement CreateTransition in repository",
  "request": {
    "workflow_id": "uuid",
    "from_step_id": "uuid",
    "to_step_id": "uuid",
    "action_name": "approve",
    "condition_type": "none",
    "condition_value": ""
  }
}
```

**Note:** The `CreateTransition` method needs to be implemented in `WorkflowRepository`.

---

## Workflow Structure Overview

### 10. Get Complete Workflow Structure

**Endpoint:** `GET /api/v1/admin/workflows/{id}/structure`

**Description:** Retrieves the complete workflow structure including template, steps, and transitions.

**Response:**
```json
{
  "workflow": {
    "id": "uuid",
    "name": "Leave Request Approval",
    "description": "Standard leave approval workflow",
    "is_active": true,
    "created_by": "user-id",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  },
  "steps": [
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "step_name": "Start",
      "step_order": 1,
      "initial": true,
      "final": false,
      "allowed_roles": [],
      "requires_all_approvers": false,
      "min_approvals": 0
    },
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "step_name": "Manager Approval",
      "step_order": 2,
      "initial": false,
      "final": false,
      "allowed_roles": ["manager"],
      "requires_all_approvers": false,
      "min_approvals": 1
    },
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "step_name": "Completed",
      "step_order": 4,
      "initial": false,
      "final": true,
      "allowed_roles": [],
      "requires_all_approvers": false,
      "min_approvals": 0
    }
  ],
  "transitions": [
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "from_step_id": "uuid",
      "to_step_id": "uuid",
      "action_name": "submit",
      "condition_type": "none",
      "condition_value": null
    },
    {
      "id": "uuid",
      "workflow_id": "uuid",
      "from_step_id": "uuid",
      "to_step_id": "uuid",
      "action_name": "approve",
      "condition_type": "none",
      "condition_value": null
    }
  ]
}
```

---

## Common Workflow Patterns

### Leave Approval Workflow

```
Start → Manager Approval → HR Review → Completed
                ↓              ↓
              Rejected      Rejected
```

### Equipment Request Workflow

```
Start → Manager Approval → Department Head → Finance Approval → Completed
                ↓                 ↓                  ↓
              Rejected         Rejected          Rejected
```

### Onboarding Workflow

```
Start → HR Setup → IT Setup → Manager Assignment → Completed
```

---

## Error Responses

All endpoints return standard error responses:

```json
{
  "error": "Error message description"
}
```

**HTTP Status Codes:**
- `400 Bad Request` - Invalid request body or parameters
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Workflow, step, or transition not found
- `500 Internal Server Error` - Server error

---

## Next Steps

To complete the workflow admin system:

1. **Implement Repository Methods:**
   - `WorkflowRepository.CreateStep()`
   - `WorkflowRepository.CreateTransition()`
   - `WorkflowRepository.UpdateWorkflow()`
   - `WorkflowRepository.DeleteWorkflow()` (soft delete)

2. **Add Update Endpoints:**
   - `PUT /api/v1/admin/workflows/{id}` - Update workflow
   - `PUT /api/v1/admin/workflows/steps/{step_id}` - Update step
   - `PUT /api/v1/admin/workflows/transitions/{transition_id}` - Update transition

3. **Add Delete Endpoints:**
   - `DELETE /api/v1/admin/workflows/{id}` - Deactivate workflow
   - `DELETE /api/v1/admin/workflows/steps/{step_id}` - Remove step
   - `DELETE /api/v1/admin/workflows/transitions/{transition_id}` - Remove transition

4. **Add Validation:**
   - Validate workflow structure (must have initial and final steps)
   - Validate transitions (no orphaned steps)
   - Validate step order uniqueness
   - Validate role names against system roles

5. **Testing:**
   - Test workflow creation flow
   - Test step management
   - Test transition creation
   - Test complete workflow structure retrieval