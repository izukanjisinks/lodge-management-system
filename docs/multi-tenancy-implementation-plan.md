# Multi-Tenancy Implementation Plan

## Overview

Each lodge/property is an **organization**. All staff data (users, rooms, bookings, clients,
invoices, meal plans, workflows) is scoped to an organization via `org_id`. Guests are
**standalone** — they live in their own table, are not org-scoped, and can book across any
organization.

The system has two distinct user classes above org staff:

- **Backoffice users** — platform-level super admins. No org, no role. Manage tenant
  organizations and other backoffice users through a separate backoffice portal.
- **Org admins** — the first admin user created per organization. Scoped to their org.
  Created automatically when a backoffice user provisions a new organization.

---

## Phase 1 — Database: Organizations, Backoffice Users & Guest Table

### 1.1 Create `organizations` table (migration 000026)

```sql
CREATE TABLE organizations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    logo_url    TEXT,
    address     TEXT,
    phone       VARCHAR(50),
    email       VARCHAR(255),
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 1.2 Create `backoffice_users` table (migration 000026, same file)

Backoffice users are not org-scoped and have no role. They are the platform operators.

```sql
CREATE TABLE backoffice_users (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name             VARCHAR(255) NOT NULL DEFAULT '',
    email                 VARCHAR(255) UNIQUE NOT NULL,
    password              VARCHAR(255) NOT NULL,
    is_active             BOOLEAN NOT NULL DEFAULT TRUE,
    change_password       BOOLEAN NOT NULL DEFAULT FALSE,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    is_locked             BOOLEAN NOT NULL DEFAULT FALSE,
    locked_until          TIMESTAMPTZ,
    last_login_at         TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_backoffice_users_email ON backoffice_users(email);
```

### 1.3 Add `org_id` to all staff tables (migration 000027)

Add `org_id UUID NOT NULL REFERENCES organizations(id)` to:
- `users`
- `roles`
- `rooms`
- `bookings`
- `individual_profiles`
- `corporate_profiles`
- `meal_plans`
- `invoices`
- `workflows`
- `workflow_instances`
- `assigned_tasks`

Add indexes: `CREATE INDEX idx_<table>_org_id ON <table>(org_id)`.

### 1.4 Create `guests` table (migration 000028)

Guests are not org-scoped. They self-register on the public website and can book at any lodge.

```sql
CREATE TABLE guests (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                 VARCHAR(255) UNIQUE NOT NULL,
    password              VARCHAR(255) NOT NULL,
    full_name             VARCHAR(255) NOT NULL DEFAULT '',
    phone                 VARCHAR(50),
    is_active             BOOLEAN NOT NULL DEFAULT TRUE,
    change_password       BOOLEAN NOT NULL DEFAULT FALSE,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    is_locked             BOOLEAN NOT NULL DEFAULT FALSE,
    locked_until          TIMESTAMPTZ,
    last_login_at         TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_guests_email ON guests(email);
```

`individual_profiles` gets a `guest_id UUID NULL REFERENCES guests(id) ON DELETE SET NULL`
column (replaces the current `user_id` link).

### 1.5 Seed default organization & backoffice user (migration 000029)

```sql
-- Default org for existing dev data
INSERT INTO organizations (name, email)
VALUES ('The Sanctuary Lodge', 'admin@lodge.dev')
ON CONFLICT DO NOTHING;

-- Default backoffice user
INSERT INTO backoffice_users (full_name, email, password)
VALUES ('Platform Admin', 'backoffice@lodge.dev', '<bcrypt of generated password>')
ON CONFLICT DO NOTHING;
```

Backfill `org_id` on all existing seeded rows to point to the default org.

---

## Phase 2 — Auth: Staff Login, Backoffice Login & Guest Login

Three completely separate login endpoints for three completely separate user classes.

### 2.1 Staff login (org-scoped)

```
POST /api/v1/auth/login  { email, password, org_id? }
```

**Step 1** — look up all users with that email across all organizations.

- **0 matches** → `401 Invalid credentials`
- **1 match** → verify password → return JWT with `org_id` embedded
- **2+ matches** → return `300 Multiple Organizations` (no password check yet):

```json
{
  "requires_org_selection": true,
  "organizations": [
    { "org_id": "...", "name": "The Sanctuary Lodge" },
    { "org_id": "...", "name": "Mountain View Lodge" }
  ]
}
```

**Step 2 (when `org_id` provided)** — find user by `(email, org_id)` → verify password → return JWT.

Staff JWT claims: `user_id`, `org_id`, `role`.

### 2.2 Backoffice login

```
POST /api/v1/backoffice/auth/login  { email, password }
```

Hits `backoffice_users` table only. No org resolution needed.

Backoffice JWT claims: `backoffice_user_id`, `role: "backoffice"` (no `org_id`).

### 2.3 Guest login

```
POST /api/v1/guest/auth/login  { email, password }
```

Hits `guests` table only. No org resolution needed.

Guest JWT claims: `guest_id`, `role: "guest"` (no `org_id`).

---

## Phase 3 — Go Models & Repository Layer

### 3.1 New models

- `internal/models/organization.go` — `Organization` struct
- `internal/models/backoffice_user.go` — `BackofficeUser` struct (no org, no role)
- `internal/models/guest.go` — `Guest` struct (no org, no role)

### 3.2 Updated `User` model

Add `OrgID uuid.UUID` field.

### 3.3 New repositories

- `OrganizationRepository` — CRUD for organizations
- `BackofficeUserRepository` — CRUD for backoffice_users, password reset
- `GuestRepository` — CRUD for guests (moves inline SQL out of `guest_auth_service.go`)

### 3.4 Updated `UserRepository`

- `GetByEmail(email)` → returns `[]User` (slice) so auth service can detect multi-org
- `GetByEmailAndOrg(email, orgID)` → single user, used for final login step
- All other methods accept `orgID uuid.UUID` and apply `WHERE org_id = $n`

### 3.5 All other repositories

Add `orgID uuid.UUID` to every query method touching org-scoped tables.

---

## Phase 4 — Middleware: Org & Backoffice Context

### 4.1 Staff middleware

Extract `org_id` from JWT, inject into context alongside `user_id` and `role`.

### 4.2 Backoffice middleware

Separate middleware that validates the backoffice JWT and injects `backoffice_user_id`.
Only applied to `/api/v1/backoffice/...` routes.

### 4.3 Helpers

```go
func GetOrgIDFromContext(ctx context.Context) (uuid.UUID, bool)
func GetBackofficeUserIDFromContext(ctx context.Context) (uuid.UUID, bool)
```

---

## Phase 5 — Services

- All org-scoped service methods accept `orgID uuid.UUID` and pass it to the repository.
- `AuthService.Login` implements the multi-org flow.
- New `BackofficeAuthService` — login, password reset for backoffice users.
- New `BackofficeOrganizationService` — provisions orgs + their first admin user in a
  single transaction (see Phase 6 for the payload shape).
- `GuestAuthService` refactored to use `GuestRepository` instead of inline SQL.
- Workflow service scopes instances and tasks to `org_id`.

---

## Phase 6 — Routes & Handlers

### Staff routes (unchanged URLs, now org-scoped internally)

All existing `/api/v1/...` routes remain. Middleware injects `org_id`; handlers pass it
to services. No URL changes visible to the frontend.

### Guest routes (standalone)

```
POST /api/v1/guest/auth/register
POST /api/v1/guest/auth/login
GET  /api/v1/guest/me
PUT  /api/v1/guest/me
POST /api/v1/guest/bookings
GET  /api/v1/guest/bookings
GET  /api/v1/guest/bookings/{id}
PATCH /api/v1/guest/bookings/{id}/cancel
```

### Backoffice routes (new, backoffice JWT required)

#### Backoffice auth
```
POST /api/v1/backoffice/auth/login
POST /api/v1/backoffice/auth/change-password
```

#### Backoffice user management
```
GET    /api/v1/backoffice/users
POST   /api/v1/backoffice/users
GET    /api/v1/backoffice/users/{id}
PUT    /api/v1/backoffice/users/{id}
DELETE /api/v1/backoffice/users/{id}
POST   /api/v1/backoffice/users/{id}/reset-password
POST   /api/v1/backoffice/users/{id}/lock
POST   /api/v1/backoffice/users/{id}/unlock
```

#### Organization management
```
GET    /api/v1/backoffice/organizations
GET    /api/v1/backoffice/organizations/{id}
PUT    /api/v1/backoffice/organizations/{id}
DELETE /api/v1/backoffice/organizations/{id}
POST   /api/v1/backoffice/organizations/provision   ← creates org + first admin in one call
```

### Organization provisioning payload

The frontend uses a stepper (org details → admin details → confirm). The final submit
sends a single payload:

```json
{
  "organization": {
    "name": "Mountain View Lodge",
    "email": "info@mountainview.com",
    "phone": "+260 97 000 0000",
    "address": "123 Mountain Rd, Lusaka"
  },
  "admin": {
    "full_name": "Jane Banda",
    "email": "jane@mountainview.com"
  }
}
```

The backend:
1. Creates the organization in a transaction
2. Creates the admin user with a randomly generated password
3. Assigns the `admin` role scoped to that org
4. Sends a welcome email to the admin with their generated password and a prompt to change it
5. Returns the created org + admin (password excluded)

---

## Phase 7 — Seed Data Updates

- Migration 000014 (`seed_dev_users`) — backfill `org_id` from default org
- Migration 000020 (`seed_demo_data`) — backfill `org_id` from default org
- Migration 000021 (`seed_workflow_data`) — backfill `org_id` from default org

---

## Key Design Decisions

| Decision | Rationale |
|---|---|
| Backoffice users in their own table | No org, no role — structurally different from staff users. Clean separation. |
| Guests in their own table | Not staff, not backoffice. No org scoping. Simpler auth flow. |
| `org_id` in JWT, not in URL | Clean URLs. Org context is part of identity, not routing. |
| Multi-org login returns org list before password check | UX: user picks org first, then one password check. |
| Org + admin created in one transaction | Atomic — no org without an admin, no orphaned admin. |
| Admin password generated on backend | Security — never transmitted in the request, only sent via email. |
| Three separate login endpoints | Each user class hits its own table. No shared auth path, no role confusion. |

---

## Implementation Order

```
Phase 1 → Phase 2 → Phase 3 → Phase 4 → Phase 5 → Phase 6 → Phase 7
```

Phases 3–5 can be partially parallelized (models + repos first, then services, then handlers).
Each phase should be a separate PR to keep diffs reviewable.
