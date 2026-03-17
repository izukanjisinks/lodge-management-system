# Lodge Management System — Backend Implementation Phases

> **Context:** This document tracks the rework of the HR System Go backend into a full Lodge Management backend. Phase 1 is complete. Phases 2–10 are the roadmap.

---

## Stack

| Layer | Technology |
|---|---|
| Language | Go |
| HTTP | `net/http` (stdlib) |
| Database | PostgreSQL |
| Migrations | `golang-migrate` |
| Auth | JWT (HS256) via `golang-jwt/jwt/v4` |
| Password hashing | bcrypt via `golang.org/x/crypto` |
| Architecture | Routes → Handlers → Services → Repositories → Models |

---

## What Was Kept from the HR System Foundation

| Component | Kept |
|---|---|
| JWT auth middleware | ✅ |
| RBAC middleware (`RequireRole`, `RequireAnyRole`) | ✅ |
| Password policy engine (complexity, history, lockout, expiry) | ✅ |
| Workflow engine (generic — steps, transitions, instances, tasks, history) | ✅ |
| User repository & service | ✅ |
| Email service (SMTP) | ✅ |
| Pagination utilities | ✅ |
| Config, database, Docker, Makefile | ✅ |

---

## Roles

| Role | Description |
|---|---|
| `admin` | Full system access |
| `manager` | Approves bookings, views reports, manages rooms |
| `receptionist` | Handles bookings, clients, invoices |
| `cleaner` | Views assigned rooms and cleaning schedule |

---

## Phase Status

| Phase | Feature | Status |
|---|---|---|
| 1 | Strip & Restructure | ✅ Done |
| 2 | Rooms Module | 🔲 Pending |
| 3 | Client Profiles | 🔲 Pending |
| 4 | Bookings Module | 🔲 Pending |
| 5 | Meals Module | 🔲 Pending |
| 6 | Invoices Module | 🔲 Pending |
| 7 | Cleaning Assignments | 🔲 Pending |
| 8 | Dashboard & KPIs | 🔲 Pending |
| 9 | Reporting | 🔲 Pending |
| 10 | Booking Approval Workflow | 🔲 Pending |

---

## Phase 1 — Strip & Restructure ✅

**Goal:** Remove all HR-specific code, rename the module, wire up lodge roles and new migrations.

### Done
- Deleted all HR handlers, services, repos, models, routes, interfaces, background jobs
- Renamed Go module `hr-system` → `lodge-system`
- Replaced 4 HR roles with lodge roles: `admin`, `manager`, `receptionist`, `cleaner`
- Wrote 14 fresh migrations (001–014):
  - 001 users, 002 roles, 003 rooms, 004 individual_profiles, 005 corporate_profiles
  - 006 meals, 007 bookings, 008 booking_meals, 009 invoices, 010 invoice_line_items
  - 011 cleaning_assignments, 012 workflow_tables, 013 password_policies, 014 seed_dev_users
- Updated `.env` → `DB_NAME=lodge-management-system`
- Kept: auth, user management, password policy, workflow engine, email service
- Created Postman collection covering all Phase 1 endpoints

### Active Endpoints After Phase 1

| Method | Path | Roles |
|---|---|---|
| GET | `/health` | public |
| POST | `/api/v1/auth/login` | public |
| GET | `/api/v1/auth/me` | all |
| POST | `/api/v1/auth/logout` | all |
| POST | `/api/v1/auth/register` | admin |
| POST | `/api/v1/auth/change-password` | all |
| GET | `/api/v1/auth/generate-password` | all |
| GET | `/api/v1/profile` | all |
| GET | `/api/v1/admin/users` | admin |
| GET | `/api/v1/admin/users/{id}` | admin |
| POST | `/api/v1/admin/users/{id}/role` | admin |
| POST | `/api/v1/admin/users/{id}/lock` | admin |
| POST | `/api/v1/admin/users/{id}/unlock` | admin |
| POST | `/api/v1/admin/users/{id}/reset-password` | admin |
| DELETE | `/api/v1/admin/users/{id}` | admin |
| GET | `/api/v1/password-policy` | admin |
| PUT | `/api/v1/password-policy` | admin |
| GET | `/api/v1/admin/workflow-types` | admin |
| POST | `/api/v1/admin/workflows` | admin |
| GET | `/api/v1/admin/workflows` | admin |
| GET | `/api/v1/admin/workflows/{id}` | admin |
| GET | `/api/v1/admin/workflows/{id}/structure` | admin |
| GET | `/api/v1/admin/workflows/{id}/steps` | admin |
| GET | `/api/v1/admin/workflows/{id}/transitions` | admin |
| PUT | `/api/v1/admin/workflows/{id}` | admin |
| DELETE | `/api/v1/admin/workflows/{id}/deactivate` | admin |
| DELETE | `/api/v1/admin/workflows/{id}` | admin |
| POST | `/api/v1/admin/workflow-steps` | admin |
| GET | `/api/v1/admin/workflow-steps/{step_id}` | admin |
| PUT | `/api/v1/admin/workflow-steps/{step_id}` | admin |
| DELETE | `/api/v1/admin/workflow-steps/{step_id}` | admin |
| GET | `/api/v1/admin/workflow-steps/{step_id}/transitions` | admin |
| POST | `/api/v1/admin/workflow-transitions` | admin |
| PUT | `/api/v1/admin/workflow-transitions/{transition_id}` | admin |
| DELETE | `/api/v1/admin/workflow-transitions/{transition_id}` | admin |
| GET | `/api/v1/workflow/my-tasks` | admin, manager, receptionist |
| GET | `/api/v1/workflow/my-tasks/pending` | admin, manager, receptionist |
| GET | `/api/v1/workflow/tasks/{id}` | admin, manager, receptionist |
| POST | `/api/v1/workflow/instances` | admin, manager |
| GET | `/api/v1/workflow/task/{task_id}/instance` | admin, manager, receptionist |
| POST | `/api/v1/workflow/instances/{id}/action` | admin, manager |
| GET | `/api/v1/workflow/instances/{id}/history` | admin, manager, receptionist |

---

## Phase 2 — Rooms Module 🔲

**Goal:** Full room inventory management. Admins and managers can create, edit, delete rooms. All staff can view. Room availability filter supports the booking wizard.

### Migrations needed
None — `000003_create_rooms` already exists.

### Files to create

```
internal/models/room.go
internal/repository/room_repository.go
internal/services/room_service.go
internal/handlers/room_handler.go
internal/routes/room_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| GET | `/api/v1/rooms` | all | List all rooms (paginated, filterable by type/availability) |
| GET | `/api/v1/rooms/available` | all | List available rooms for a date range (`?check_in=&check_out=`) |
| GET | `/api/v1/rooms/{id}` | all | Get single room |
| POST | `/api/v1/rooms` | admin, manager | Create room |
| PUT | `/api/v1/rooms/{id}` | admin, manager | Update room |
| PATCH | `/api/v1/rooms/{id}/availability` | admin, manager | Toggle availability |
| DELETE | `/api/v1/rooms/{id}` | admin, manager | Delete room |

### Room model

```go
type Room struct {
    ID            uuid.UUID `json:"id"`
    Name          string    `json:"name"`
    Type          string    `json:"type"`      // single|double|suite|cabin|conference
    Capacity      int       `json:"capacity"`
    PricePerNight float64   `json:"price_per_night"`
    Amenities     []string  `json:"amenities"`
    IsAvailable   bool      `json:"is_available"`
    Description   string    `json:"description,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

### Availability query logic
A room is considered available for `[check_in, check_out)` if no confirmed or checked_in booking overlaps that window:

```sql
SELECT r.* FROM rooms r
WHERE r.is_available = TRUE
  AND r.id NOT IN (
    SELECT room_id FROM bookings
    WHERE status IN ('confirmed', 'checked_in')
      AND check_in  < $check_out
      AND check_out > $check_in
  )
```

---

## Phase 3 — Client Profiles 🔲

**Goal:** Clients (individuals and corporates) can register and manage their profiles. Staff can view and create client profiles on behalf of walk-ins.

### Migrations needed
None — `000004_create_individual_profiles` and `000005_create_corporate_profiles` already exist.

### Files to create

```
internal/models/client.go
internal/repository/client_repository.go
internal/services/client_service.go
internal/handlers/client_handler.go
internal/routes/client_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| POST | `/api/v1/auth/register` | public (or admin) | Register with individual or corporate profile in one transaction |
| GET | `/api/v1/clients/me` | authenticated user (client) | Get own profile |
| PUT | `/api/v1/clients/me` | authenticated user (client) | Update own profile |
| GET | `/api/v1/admin/clients` | admin, manager, receptionist | List all clients (individual + corporate) |
| GET | `/api/v1/admin/clients/individual` | admin, manager, receptionist | List individual clients |
| GET | `/api/v1/admin/clients/corporate` | admin, manager, receptionist | List corporate clients |
| GET | `/api/v1/admin/clients/{id}` | admin, manager, receptionist | Get client profile |
| POST | `/api/v1/admin/clients/individual` | admin, manager, receptionist | Create individual profile |
| POST | `/api/v1/admin/clients/corporate` | admin, manager, receptionist | Create corporate profile |
| PUT | `/api/v1/admin/clients/{id}` | admin, manager, receptionist | Update client profile |
| DELETE | `/api/v1/admin/clients/{id}` | admin | Delete client |

### Notes
- Registration creates a `users` row + matching `individual_profiles` or `corporate_profiles` row in a single transaction.
- `auth/register` will need to be extended to support a `client_type` field and profile payload.
- Role assigned on register: derive a `client` role or leave as no-role (TBD in Phase 3).

---

## Phase 4 — Bookings Module 🔲

**Goal:** Clients can create bookings. Staff can manage status transitions (confirm, check-in, check-out, cancel).

**Depends on:** Phase 2 (rooms), Phase 3 (client profiles)

### Migrations needed
None — `000007_create_bookings` and `000008_create_booking_meals` already exist.

### Files to create

```
internal/models/booking.go
internal/repository/booking_repository.go
internal/services/booking_service.go
internal/handlers/booking_handler.go
internal/routes/booking_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| POST | `/api/v1/bookings` | authenticated | Create booking |
| GET | `/api/v1/bookings/my` | authenticated | List own bookings |
| GET | `/api/v1/bookings/{id}` | authenticated | Get booking detail |
| PATCH | `/api/v1/bookings/{id}/cancel` | authenticated (own) | Cancel own booking |
| GET | `/api/v1/admin/bookings` | admin, manager, receptionist | List all bookings (filterable) |
| PATCH | `/api/v1/admin/bookings/{id}/confirm` | admin, manager | Confirm booking |
| PATCH | `/api/v1/admin/bookings/{id}/cancel` | admin, manager, receptionist | Cancel booking |
| PATCH | `/api/v1/admin/bookings/{id}/check-in` | admin, manager, receptionist | Mark checked in |
| PATCH | `/api/v1/admin/bookings/{id}/check-out` | admin, manager, receptionist | Mark checked out |

### Status transitions

```
pending → confirmed  (admin/manager confirm)
pending → cancelled  (client or staff cancel)
confirmed → checked_in   (receptionist/manager check-in)
confirmed → cancelled    (staff cancel)
checked_in → checked_out (receptionist/manager check-out)
```

### Booking service responsibilities
- Validate room availability for requested dates before creating
- Calculate `room_cost = nights × price_per_night`
- Store selected meal IDs in `booking_meals` join table
- On `confirm`: auto-trigger invoice creation (Phase 6)
- On `check-out`: auto-trigger cleaning assignment creation (Phase 7)

---

## Phase 5 — Meals Module 🔲

**Goal:** Admin/manager manage the meal catalogue. All authenticated users can view available meals (for booking wizard).

**Depends on:** Phase 2 (can run in parallel)

### Migrations needed
None — `000006_create_meals` already exists.

### Files to create

```
internal/models/meal.go
internal/repository/meal_repository.go
internal/services/meal_service.go
internal/handlers/meal_handler.go
internal/routes/meal_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| GET | `/api/v1/meals` | all | List available meals (filterable by type) |
| GET | `/api/v1/meals/{id}` | all | Get single meal |
| POST | `/api/v1/admin/meals` | admin, manager | Create meal |
| PUT | `/api/v1/admin/meals/{id}` | admin, manager | Update meal |
| PATCH | `/api/v1/admin/meals/{id}/availability` | admin, manager | Toggle availability |
| DELETE | `/api/v1/admin/meals/{id}` | admin, manager | Delete meal |

---

## Phase 6 — Invoices Module 🔲

**Goal:** Invoices are auto-generated when a booking is confirmed. Clients can view and download their invoices. Staff can manage invoice status.

**Depends on:** Phase 4 (bookings)

### Migrations needed
None — `000009_create_invoices` and `000010_create_invoice_line_items` already exist.

### Files to create

```
internal/models/invoice.go
internal/repository/invoice_repository.go
internal/services/invoice_service.go
internal/handlers/invoice_handler.go
internal/routes/invoice_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| GET | `/api/v1/invoices/my` | authenticated | List own invoices |
| GET | `/api/v1/invoices/{id}` | authenticated (own) | Get invoice detail with line items |
| GET | `/api/v1/bookings/{id}/invoice` | authenticated (own) | Get invoice for a specific booking |
| GET | `/api/v1/admin/invoices` | admin, manager, receptionist | List all invoices (filterable by status) |
| PATCH | `/api/v1/admin/invoices/{id}/mark-paid` | admin, manager | Mark invoice as paid |
| PATCH | `/api/v1/admin/invoices/{id}/issue` | admin, manager, receptionist | Issue a draft invoice |
| DELETE | `/api/v1/admin/invoices/{id}` | admin | Delete invoice |

### Invoice auto-generation
Triggered inside `booking_service.ConfirmBooking()`:

```
Line items:
  - Room: "{room_name} × {nights} nights"  → room_cost
  - Meal: "{meal_name}" (one line per meal) → meal.price × nights (or per stay TBD)
Subtotal = sum of line items
Tax      = subtotal × tax_rate (configurable, default 0)
Total    = subtotal + tax
Status   = 'issued'
Due date = issued_at + 14 days
```

---

## Phase 7 — Cleaning Assignments 🔲

**Goal:** After check-out, a cleaning assignment is auto-created for the room. Cleaners see their daily schedule. Managers can assign and track.

**Depends on:** Phase 2 (rooms), Phase 4 (bookings — check-out trigger)

### Migrations needed
None — `000011_create_cleaning_assignments` already exists.

### Files to create

```
internal/models/cleaning.go
internal/repository/cleaning_repository.go
internal/services/cleaning_service.go
internal/handlers/cleaning_handler.go
internal/routes/cleaning_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| GET | `/api/v1/cleaning/my` | cleaner | Get own cleaning assignments for today |
| PATCH | `/api/v1/cleaning/{id}/complete` | cleaner | Mark assignment as completed |
| PATCH | `/api/v1/cleaning/{id}/start` | cleaner | Mark assignment as in_progress |
| GET | `/api/v1/admin/cleaning` | admin, manager | List all assignments (filterable by date/status/room) |
| POST | `/api/v1/admin/cleaning` | admin, manager | Create manual cleaning assignment |
| PUT | `/api/v1/admin/cleaning/{id}` | admin, manager | Update assignment (reassign, reschedule) |
| DELETE | `/api/v1/admin/cleaning/{id}` | admin, manager | Delete assignment |

### Auto-assignment on check-out
Triggered inside `booking_service.CheckOut()`:

```go
cleaningService.CreateAssignment(CleaningAssignment{
    RoomID:        booking.RoomID,
    AssignedBy:    staffUserID,
    ScheduledDate: today,
    // AssignedTo: load-balance across available cleaners
})
```

---

## Phase 8 — Dashboard & KPIs 🔲

**Goal:** Role-aware dashboard endpoint. Admin/manager see operational KPIs. Clients see booking summary.

**Depends on:** Phases 4, 5, 6

### Files to create

```
internal/models/dashboard.go
internal/repository/dashboard_repository.go
internal/services/dashboard_service.go
internal/handlers/dashboard_handler.go
internal/routes/dashboard_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| GET | `/api/v1/admin/dashboard` | admin, manager | Staff KPI dashboard |
| GET | `/api/v1/dashboard` | authenticated | Client booking summary |

### Admin dashboard payload

```json
{
  "today": {
    "check_ins":  4,
    "check_outs": 2,
    "new_bookings": 1
  },
  "rooms": {
    "total": 20,
    "occupied": 12,
    "available": 7,
    "not_ready": 1
  },
  "revenue": {
    "this_month": 45000.00,
    "last_month": 38000.00
  },
  "bookings": {
    "pending_confirmation": 3,
    "active": 12
  }
}
```

### Client dashboard payload

```json
{
  "active_bookings": [...],
  "upcoming_check_in": { "booking_id": "...", "days_until": 3 },
  "outstanding_invoices": 1,
  "total_stays": 5
}
```

---

## Phase 9 — Reporting 🔲

**Goal:** Admin and manager can query aggregated data for occupancy, revenue, and meal demand.

**Depends on:** Phases 4, 5, 6

### Files to create

```
internal/handlers/report_handler.go
internal/services/report_service.go
internal/repository/report_repository.go
internal/routes/report_routes.go
```

### Endpoints

| Method | Path | Roles | Description |
|---|---|---|---|
| GET | `/api/v1/admin/reports/occupancy` | admin, manager | Room occupancy by date range (`?from=&to=`) |
| GET | `/api/v1/admin/reports/revenue` | admin, manager | Revenue breakdown by date range |
| GET | `/api/v1/admin/reports/meals` | admin, manager | Meal demand — most ordered meals |
| GET | `/api/v1/admin/reports/bookings` | admin, manager | Bookings by status over time |

---

## Phase 10 — Booking Approval Workflow 🔲

**Goal:** Wire the existing generic workflow engine to the bookings domain. Certain bookings (e.g. corporate, long-stay, high-value) require manager approval before confirmation.

**Depends on:** Phases 4, workflow engine (Phase 1)

### No new files needed — uses existing workflow engine

### What needs to be added
1. **Seed a booking approval workflow** — new migration `000015_seed_booking_approval_workflow.up.sql`:
   - Step 1 (initial): "Submitted" — created by system
   - Step 2: "Pending Manager Review" — `allowed_roles: [manager, admin]`
   - Step 3 (final): "Approved"
   - Transitions: submit → step 2, approve → step 3, reject → rejected

2. **Trigger workflow from booking service** — in `booking_service.CreateBooking()`:
   ```go
   if shouldRequireApproval(booking) {
       workflowService.InitiateWorkflow("booking_approval", taskDetails, userID, "medium", nil)
   }
   ```

3. **Update booking confirm endpoint** to also accept workflow-driven approval (i.e. workflow completion triggers booking confirmation).

### Approval trigger criteria (configurable)
- Booking is for a corporate client
- Total value > threshold (e.g. > $5000)
- Stay duration > 7 nights

---

## Parallel Phase Map

```
Phase 1 ──────────────────────────────── ✅ Done
           │
     ┌─────┴──────┐
     ▼            ▼
  Phase 2       Phase 5
  (Rooms)       (Meals)
     │            │
     └─────┬──────┘
           ▼
        Phase 3
     (Client Profiles)
           │
           ▼
        Phase 4
        (Bookings) ──────────────────────────────── Phase 10
           │                                         (Workflow)
     ┌─────┴──────────────┐
     ▼                    ▼
  Phase 6             Phase 7
  (Invoices)          (Cleaning)
     │
     └─────┬──────┐
           ▼      ▼
        Phase 8  Phase 9
       (Dashboard) (Reports)
```
