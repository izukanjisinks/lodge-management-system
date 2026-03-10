# HR System Architecture

## Overview

This is a REST API built in Go using the standard `net/http` package (Go 1.22+). It follows a layered architecture:

```
HTTP Request → Routes → Middleware → Handlers → Services → Repositories → PostgreSQL
```

- **Routes** — register URL patterns and wrap handlers with auth/role middleware
- **Handlers** — parse requests, call services, write responses
- **Services** — business logic and validation
- **Repositories** — SQL queries and database interaction
- **Models** — shared data structures

---

## Users vs Employees

These are two distinct concepts in the system. Understanding the difference is important.

### Users (`users` table)

A **User** is a system account used for authentication. It answers the question: *"who can log in?"*

| Field      | Type      | Description                          |
|------------|-----------|--------------------------------------|
| `user_id`  | UUID      | Primary key                          |
| `email`    | string    | Login email (unique)                 |
| `password` | string    | bcrypt-hashed password (never returned in API responses) |
| `role_id`  | UUID (FK) | References the `roles` table         |
| `is_active`| bool      | Whether the account is enabled       |
| `created_at` / `updated_at` | timestamp | Audit fields |

A user exists purely for authentication and authorization. It has no knowledge of HR data — no name, no department, no hire date.

### Employees (`employees` table)

An **Employee** is an HR record representing a person who works (or worked) at the company. It answers the question: *"who is employed here and what do we know about them?"*

**Personal information:**
- `first_name`, `last_name`, `email`, `personal_email`, `phone`
- `date_of_birth`, `gender`, `national_id`, `marital_status`
- `address`, `city`, `state`, `country`
- `profile_photo_url`

**Employment information:**
- `employee_number` — unique identifier (e.g. `EMP-001`)
- `department_id` — which department they belong to
- `position_id` — their job position/title
- `manager_id` — their direct manager (self-referencing FK to `employees`)
- `hire_date`, `probation_end_date`
- `employment_type` — `full_time`, `part_time`, `contract`, or `intern`
- `employment_status` — `active`, `on_leave`, `suspended`, `terminated`, or `resigned`
- `termination_date`, `termination_reason`

**Audit:**
- `created_at`, `updated_at`, `deleted_at` (soft delete)

### The Link Between Them

An employee record has an **optional** `user_id` foreign key:

```
employees.user_id → users.user_id (nullable)
```

This means:

- An employee **can exist without a user account** — for example, a new hire whose system account hasn't been created yet, or a terminated employee whose account was deactivated.
- A user account **can exist without an employee record** — for example, a system administrator who manages the platform but is not an actual employee in the HR sense.
- When both exist and are linked, the user can log in and perform self-service actions (clock in/out, apply for leave, view their own balances) tied to their employee record.

**In practice, most active employees will have a linked user account.**

---

## Roles and Access Control

There are four predefined roles:

| Role          | Description                                               |
|---------------|-----------------------------------------------------------|
| `super_admin` | Full system access — can do everything                    |
| `hr_manager`  | Manage employees, leave, payroll, recruitment             |
| `manager`     | View team, approve leave, read attendance and performance |
| `employee`    | Self-service only — own profile, own leave, clock in/out  |

Roles are stored in a `roles` table and referenced by `users.role_id`. Enforcement happens in route middleware (`withAuthAndRole`) and in the `User.HasPermission()` method.

---

## Related HR Records

These tables extend the employee's record:

| Table               | Purpose                                                  |
|---------------------|----------------------------------------------------------|
| `emergency_contacts`| Emergency contacts for an employee                       |
| `employee_documents`| Uploaded HR documents (contracts, IDs, certificates)     |
| `leave_types`       | Types of leave available (annual, sick, maternity, etc.) |
| `leave_balances`    | How many days of each leave type an employee has per year|
| `leave_requests`    | Leave applications (pending → approved/rejected/cancelled)|
| `attendance`        | Clock-in/clock-out records per employee per day          |
| `holidays`          | Company-wide or location-specific public holidays        |
| `departments`       | Organisational departments (supports parent/child tree)  |
| `positions`         | Job titles/positions within the company                  |

---

## Authentication Flow

1. Client sends `POST /api/v1/auth/login` with `email` + `password`
2. Server looks up the user by email, verifies the bcrypt password hash
3. On success, returns a signed JWT (HS256, 24-hour expiry) containing `userId` and `email`
4. All subsequent requests must include `Authorization: Bearer <token>`
5. The `AuthMiddleware` validates the token and attaches the `User` object to the request context
6. Handlers retrieve the user/employee from context to perform actions

---

## API Structure

All endpoints are prefixed with `/api/v1`.

| Prefix                        | Area                  | Auth required |
|-------------------------------|-----------------------|---------------|
| `/api/v1/auth/...`            | Login, register, me   | Partial       |
| `/api/v1/hr/departments/...`  | Departments           | Yes           |
| `/api/v1/hr/positions/...`    | Positions             | Yes           |
| `/api/v1/hr/employees/...`    | Employees + documents + emergency contacts | Yes |
| `/api/v1/hr/leave/types/...`  | Leave types           | Yes           |
| `/api/v1/hr/leave/balances/...` | Leave balances      | Yes           |
| `/api/v1/hr/leave/requests/...` | Leave requests      | Yes           |
| `/api/v1/hr/attendance/...`   | Attendance / clock-in | Yes           |
| `/api/v1/hr/holidays/...`     | Public holidays       | Yes           |
