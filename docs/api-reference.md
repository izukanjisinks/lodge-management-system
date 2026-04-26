# API Reference — Lodge Management System

## Multi-Tenancy Overview

The backend has shifted to a full multi-tenant architecture with three distinct user classes. Every piece of data is scoped to an organisation; no data leaks across org boundaries.

### User Classes

| Class | Token | Scope | Base path |
|---|---|---|---|
| **Staff** | `Bearer <token>` | Org-scoped — JWT contains `org_id`, every query filters by it | `/api/v1/` |
| **Backoffice** | `Bearer <backoffice_token>` | Platform-level — no org, manages all orgs | `/api/v1/backoffice/` |
| **Guest** | `Bearer <guest_token>` | Standalone — no org; bookings derive org from room | `/api/v1/guest/` |

### Key Concepts for the Frontend

- **Org scoping is invisible to staff.** A logged-in staff user only ever sees their own org's data. No `org_id` needs to be sent in request bodies — it is read from the JWT on every request.
- **Multi-org staff login is a two-step flow.** If a staff email exists in more than one org, step 1 returns an `orgs` array instead of a token. The frontend must present an org picker and re-send with `org_id` to complete login.
- **Backoffice has no org.** It is a platform-admin view that can see and manage all organisations. It provisions new orgs (which auto-creates an admin user and emails them credentials).
- **`change_password: true`** is returned in the login response for any user whose password must be changed on first login. The frontend should redirect to a change-password screen in this case.

---

## Staff API — `/api/v1/`

All staff endpoints require `Authorization: Bearer <token>` unless marked **Public**.

Roles: `admin`, `manager`, `receptionist`

### Authentication

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/v1/auth/login` | Public | Step 1: email + password. Returns token or org list for multi-org users. Send `org_id` in body for step 2. |
| `GET` | `/api/v1/auth/me` | Staff | Returns the authenticated user's profile |
| `POST` | `/api/v1/auth/logout` | Staff | Stateless logout (client discards token) |
| `GET` | `/api/v1/profile` | Staff | Alias for `/auth/me` — returns profile of authenticated user |

#### Login Request
```json
{ "email": "string", "password": "string", "org_id": "uuid (optional, step 2 only)" }
```

#### Login Response — single org
```json
{
  "token": "string",
  "user": {
    "user_id": "uuid",
    "full_name": "string",
    "email": "string",
    "role": "admin | manager | receptionist",
    "org_id": "uuid",
    "change_password": false
  }
}
```

#### Login Response — multi-org (step 1)
```json
{
  "requires_org_selection": true,
  "orgs": [
    { "org_id": "uuid", "org_name": "string" }
  ]
}
```

---

### Password

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `POST` | `/api/v1/auth/change-password` | Staff | All | Change own password. Body: `{ old_password, new_password }` |
| `GET` | `/api/v1/password-policy` | Staff | `admin` | Get current password policy |
| `PUT` | `/api/v1/password-policy` | Staff | `admin` | Update password policy |
| `GET` | `/api/v1/auth/generate-password` | Staff | All | Generate a password that satisfies the current policy |
| `POST` | `/api/v1/admin/users/{id}/reset-password` | Staff | `admin` | Reset another user's password (sends email with temp password) |

---

### Staff User Management

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `POST` | `/api/v1/users` | Staff | `admin` | Create a new staff user in the org |
| `GET` | `/api/v1/users` | Staff | `admin` | List staff users. Query: `search`, `role_id`, `is_active`, `page`, `page_size` |
| `GET` | `/api/v1/users/{id}` | Staff | `admin` | Get a staff user by ID |
| `PUT` | `/api/v1/users/{id}` | Staff | `admin` | Update a staff user (full_name, email, password, role, status) |
| `DELETE` | `/api/v1/users/{id}` | Staff | `admin` | Delete a staff user |
| `POST` | `/api/v1/users/{id}/lock` | Staff | `admin` | Lock a user account |
| `POST` | `/api/v1/users/{id}/unlock` | Staff | `admin` | Unlock a user account |
| `POST` | `/api/v1/admin/users/{id}/role` | Staff | `admin` | Change a user's role. Body: `{ role_id }` |

---

### Rooms

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/rooms` | Staff | All | List org rooms. Query: `type`, `is_available`, `page`, `page_size` |
| `GET` | `/api/v1/rooms/available` | Staff | All | List available rooms for a date range. Query: `check_in`, `check_out` (YYYY-MM-DD), `type` |
| `GET` | `/api/v1/rooms/{id}` | Staff | All | Get room by ID |
| `POST` | `/api/v1/rooms` | Staff | `admin`, `manager` | Create a room |
| `PUT` | `/api/v1/rooms/{id}` | Staff | `admin`, `manager` | Update room details |
| `PUT` | `/api/v1/rooms/{id}/images` | Staff | `admin`, `manager` | Replace room image list. Body: `{ images: string[] }` |
| `PATCH` | `/api/v1/rooms/{id}/availability` | Staff | `admin`, `manager` | Toggle availability. Body: `{ is_available: bool }` |
| `DELETE` | `/api/v1/rooms/{id}` | Staff | `admin`, `manager` | Delete a room |

---

### Meal Plans

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/meal-plans` | Staff | All | List org meal plans. Query: `is_active`, `page`, `page_size` |
| `GET` | `/api/v1/meal-plans/{id}` | Staff | All | Get meal plan by ID |
| `POST` | `/api/v1/meal-plans` | Staff | `admin`, `manager` | Create a meal plan |
| `PUT` | `/api/v1/meal-plans/{id}` | Staff | `admin`, `manager` | Update a meal plan |
| `DELETE` | `/api/v1/meal-plans/{id}` | Staff | `admin`, `manager` | Delete a meal plan |

---

### Clients

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/clients/individual` | Staff | All | List individual clients. Query: `search`, `page`, `page_size` |
| `GET` | `/api/v1/clients/individual/{id}` | Staff | All | Get individual client by ID |
| `POST` | `/api/v1/clients/individual` | Staff | `admin`, `manager`, `receptionist` | Create individual client |
| `PUT` | `/api/v1/clients/individual/{id}` | Staff | `admin`, `manager`, `receptionist` | Update individual client |
| `DELETE` | `/api/v1/clients/individual/{id}` | Staff | `admin`, `manager` | Delete individual client |
| `GET` | `/api/v1/clients/corporate` | Staff | All | List corporate clients. Query: `search`, `page`, `page_size` |
| `GET` | `/api/v1/clients/corporate/{id}` | Staff | All | Get corporate client by ID |
| `POST` | `/api/v1/clients/corporate` | Staff | `admin`, `manager`, `receptionist` | Create corporate client |
| `PUT` | `/api/v1/clients/corporate/{id}` | Staff | `admin`, `manager`, `receptionist` | Update corporate client |
| `DELETE` | `/api/v1/clients/corporate/{id}` | Staff | `admin`, `manager` | Delete corporate client |

---

### Bookings

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/bookings` | Staff | All | List org bookings. Query: `status`, `client_id`, `room_id`, `page`, `page_size` |
| `GET` | `/api/v1/bookings/{id}` | Staff | All | Get booking by ID |
| `POST` | `/api/v1/bookings` | Staff | `admin`, `manager`, `receptionist` | Create a booking (staff-initiated, auto-confirms without workflow) |
| `PUT` | `/api/v1/bookings/{id}` | Staff | `admin`, `manager`, `receptionist` | Update booking details |
| `PATCH` | `/api/v1/bookings/{id}/status` | Staff | `admin`, `manager`, `receptionist` | Update booking status. Body: `{ status }` |
| `DELETE` | `/api/v1/bookings/{id}` | Staff | `admin`, `manager`, `receptionist` | Delete a booking |

#### Booking Status Lifecycle
`pending → confirmed → checked_in → checked_out` or `→ cancelled`

> Guest-initiated bookings go through the approval workflow before reaching `confirmed`. Staff-initiated bookings skip the workflow.

---

### Invoices

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/invoices` | Staff | `admin`, `manager` | List org invoices |
| `GET` | `/api/v1/invoices/{id}` | Staff | `admin`, `manager`, `receptionist` | Get invoice by ID |
| `GET` | `/api/v1/invoices/booking/{booking_id}` | Staff | All | Get invoice for a booking |
| `PATCH` | `/api/v1/invoices/{id}/status` | Staff | `admin`, `manager` | Update invoice status. Body: `{ status }` |

#### Invoice Status Lifecycle
`draft → issued → paid` or `→ cancelled`

---

### Dashboard

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/dashboard/stats` | Staff | `admin`, `manager`, `receptionist` | Org dashboard stats (occupancy, revenue, bookings summary) |

---

### Workflow — Task Actions

Used by staff to view and process their assigned approval tasks.

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/workflow/my-tasks` | Staff | All | List tasks assigned to me. Query: `status` (pending/completed) |
| `GET` | `/api/v1/workflow/my-tasks/pending` | Staff | All | Shortcut — pending tasks only |
| `GET` | `/api/v1/workflow/tasks/{id}` | Staff | All | Get task details with full workflow context |
| `GET` | `/api/v1/workflow/task/{task_id}/instance` | Staff | All | Get workflow instance for a task |
| `GET` | `/api/v1/workflow/instances/{id}/history` | Staff | All | Full audit history of a workflow instance |
| `POST` | `/api/v1/workflow/instances/{id}/action` | Staff | `admin`, `manager` | Process an action on a workflow instance. Body: `{ action, comments }` |
| `POST` | `/api/v1/workflow/instances` | Staff | `admin`, `manager` | Manually initiate a workflow (rarely needed; guest bookings do this automatically) |

#### Actions
- `approve` — approve the current step and advance the workflow
- `reject` — reject and close the workflow (booking → cancelled)

---

### Workflow — Template Administration

Used by org admins to configure the booking approval workflow structure.

| Method | Path | Auth | Roles | Description |
|---|---|---|---|---|
| `GET` | `/api/v1/admin/workflow-types` | Staff | `admin` | List available workflow type constants |
| `GET` | `/api/v1/admin/workflows` | Staff | `admin` | List all org workflow templates |
| `POST` | `/api/v1/admin/workflows` | Staff | `admin` | Create a workflow template |
| `GET` | `/api/v1/admin/workflows/{id}` | Staff | `admin` | Get a workflow template |
| `PUT` | `/api/v1/admin/workflows/{id}` | Staff | `admin` | Update a workflow template |
| `DELETE` | `/api/v1/admin/workflows/{id}/deactivate` | Staff | `admin` | Soft-deactivate a workflow |
| `DELETE` | `/api/v1/admin/workflows/{id}` | Staff | `admin` | Permanently delete a workflow |
| `GET` | `/api/v1/admin/workflows/{id}/structure` | Staff | `admin` | Full structure: steps + transitions |
| `GET` | `/api/v1/admin/workflows/{id}/steps` | Staff | `admin` | List steps for a workflow |
| `GET` | `/api/v1/admin/workflows/{id}/transitions` | Staff | `admin` | List transitions for a workflow |
| `POST` | `/api/v1/admin/workflow-steps` | Staff | `admin` | Create a workflow step |
| `GET` | `/api/v1/admin/workflow-steps/{step_id}` | Staff | `admin` | Get a step |
| `PUT` | `/api/v1/admin/workflow-steps/{step_id}` | Staff | `admin` | Update a step |
| `DELETE` | `/api/v1/admin/workflow-steps/{step_id}` | Staff | `admin` | Delete a step |
| `GET` | `/api/v1/admin/workflow-steps/{step_id}/transitions` | Staff | `admin` | List valid transitions from a step |
| `POST` | `/api/v1/admin/workflow-transitions` | Staff | `admin` | Create a transition between steps |
| `PUT` | `/api/v1/admin/workflow-transitions/{transition_id}` | Staff | `admin` | Update a transition |
| `DELETE` | `/api/v1/admin/workflow-transitions/{transition_id}` | Staff | `admin` | Delete a transition |

---

## Backoffice API — `/api/v1/backoffice/`

All backoffice endpoints require `Authorization: Bearer <backoffice_token>` unless marked **Public**. Backoffice users have no org — they see all organisations.

### Authentication

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/v1/backoffice/auth/login` | Public | Login with email + password. Returns `backoffice_token`. |
| `GET` | `/api/v1/backoffice/auth/me` | Backoffice | Get own profile |
| `POST` | `/api/v1/backoffice/auth/change-password` | Backoffice | Change own password. Body: `{ current_password, new_password }` |

---

### Backoffice User Management

Manage who has backoffice (platform admin) access.

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/v1/backoffice/users` | Backoffice | List all backoffice users |
| `POST` | `/api/v1/backoffice/users` | Backoffice | Create a backoffice user |
| `GET` | `/api/v1/backoffice/users/{id}` | Backoffice | Get a backoffice user |
| `PUT` | `/api/v1/backoffice/users/{id}` | Backoffice | Update a backoffice user |
| `DELETE` | `/api/v1/backoffice/users/{id}` | Backoffice | Delete a backoffice user |
| `POST` | `/api/v1/backoffice/users/{id}/reset-password` | Backoffice | Reset a backoffice user's password |
| `POST` | `/api/v1/backoffice/users/{id}/lock` | Backoffice | Lock a backoffice user account |
| `POST` | `/api/v1/backoffice/users/{id}/unlock` | Backoffice | Unlock a backoffice user account |

---

### Organisation Management

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/v1/backoffice/organizations` | Backoffice | List all organisations |
| `POST` | `/api/v1/backoffice/organizations/provision` | Backoffice | Create an org + its first admin user in one step. Emails the admin their credentials. |
| `GET` | `/api/v1/backoffice/organizations/{id}` | Backoffice | Get an organisation |
| `PUT` | `/api/v1/backoffice/organizations/{id}` | Backoffice | Update org details (name, email, phone, address, logo_url) |
| `DELETE` | `/api/v1/backoffice/organizations/{id}` | Backoffice | Delete an organisation |

#### Provision Request Body
```json
{
  "organization": {
    "name": "string",
    "email": "string",
    "phone": "string",
    "address": "string",
    "logo_url": "string"
  },
  "admin": {
    "full_name": "string",
    "email": "string"
  }
}
```

The admin's password is auto-generated and emailed to them. They are flagged `change_password: true` and must change it on first login.

---

## Public Endpoints

No authentication required.

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/api/v1/reviews/summary` | Public lodge rating summary (no org scope — aggregated) |
