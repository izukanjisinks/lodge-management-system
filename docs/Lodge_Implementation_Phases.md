# Lodge Management System — Phased Implementation Plan

> **Context:** This document maps the existing HR System Frontend codebase to the Lodge Administration target described in `Lodge_Administration_VueJS_Guide.md`. The goal is to **reuse and rework** — not rewrite from scratch. Every shadcn component, the Pinia setup, the API client, the auth store, router guards, composables, and layout system all carry over unchanged or with minimal edits.

---

## What Stays, What Goes, What Changes

| Layer | Keep As-Is | Rework | Delete |
|---|---|---|---|
| `src/components/ui/**` | All 40+ shadcn components | — | — |
| `src/lib/utils.ts` | `cn()` utility | — | — |
| `src/composables/*` | All three composables | — | — |
| `src/layouts/AuthenticatedLayout.vue` | Sidebar wiring, global dialogs | — | — |
| `src/services/api/client.ts` | Bearer token, 403 handler, base URL | — | — |
| `src/stores/auth.ts` | Token storage key rename | Role enum, role labels | — |
| `src/types/auth.ts` | `LoginCredentials`, `LoginResponse`, `ApiError` | `AuthUser`, `UserRole` | — |
| `src/views/auth/LoginView.vue` | Form logic | Branding/logo | — |
| `src/components/layout/AppSidebar.vue` | Sidebar shell | All nav items | — |
| `src/router/index.ts` | Guard logic | All routes | — |
| `src/styles/global.css` | Tailwind + CSS vars | — | — |
| All HR-specific views | — | — | Replace with lodge views |
| All HR-specific stores | — | — | Replace with lodge stores |
| All HR-specific types | — | — | Replace with lodge types |
| All HR-specific API services | — | — | Replace with lodge API services |

---

## Phase 1 — Foundation & Auth Cleanup

**Goal:** Remove all HR-specific code, update identity layer (roles, branding, token key), and confirm login still works end-to-end against the existing backend auth routes.

### 1.1 Update Auth Types

File: `src/types/auth.ts`

```typescript
export type UserRole = 'admin' | 'client_individual' | 'client_corporate'

export interface AuthUser {
  id: number
  email: string
  role: UserRole
  full_name: string
  created_at: string
  change_password: boolean
  is_active: boolean
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: AuthUser
}

export interface ApiError {
  error: { code: string; message: string; details?: unknown }
}
```

### 1.2 Update Auth Store

File: `src/stores/auth.ts`

- Change localStorage key from `hr_token` to `lodge_token`
- Update `roleLabel` computed to return `'Administrator' | 'Individual Client' | 'Corporate Client'`
- Keep all actions: `login`, `fetchCurrentUser`, `logout`, `setToken`, `clearAuth`

### 1.3 Update Login View Branding

File: `src/views/auth/LoginView.vue`

- Replace HR logo with lodge logo (or placeholder)
- Update page title / app name string to "Lodge Management"
- Keep all form logic, error handling, and `ChangePasswordDialog` integration

### 1.4 Clear All HR Routes

File: `src/router/index.ts`

- Keep `/login` route (guestOnly)
- Keep the `AuthenticatedLayout` wrapper
- Keep the navigation guard logic (`requiresAuth`, `guestOnly`, token check, `fetchCurrentUser`)
- Remove all HR-specific route entries — replace temporarily with just `/dashboard`
- Routes will be added progressively in later phases

### 1.5 Clear AppSidebar Navigation

File: `src/components/layout/AppSidebar.vue`

- Keep the sidebar shell, user menu (profile, logout, change password)
- Remove all HR nav groups
- Add placeholder nav groups for lodge (will be fleshed out per phase)

### 1.6 Delete HR-Specific Files

Remove the following directories and their contents (they will be replaced in later phases):

```
src/views/employees/
src/views/departments/
src/views/positions/
src/views/leave/
src/views/attendance/
src/views/holidays/
src/views/payroll/
src/views/recruitment/
src/views/performance/
src/views/workflows/
src/views/approvals/
src/views/system/
src/components/approvals/
src/components/departments/
src/components/employees/
src/components/holidays/
src/components/leave/
src/components/payroll/
src/components/performance/
src/components/positions/
src/components/recruitment/
src/components/workflow/
src/components/users/
src/components/dashboard/
src/stores/workflow.ts
src/services/api/dashboard.ts
src/services/api/department.ts
src/services/api/employee.ts
src/services/api/leave.ts
src/services/api/password.ts
src/services/api/payroll.ts
src/services/api/position.ts
src/services/api/role.ts
src/services/api/user.ts
src/services/api/workflow.ts
src/types/dashboard.ts
src/types/employee.ts
src/types/leave.ts
src/types/password.ts
src/types/payroll.ts
src/types/role.ts
src/types/user.ts
src/types/workflow.ts
src/lib/salary.ts
```

### 1.7 Create Stub Dashboard

File: `src/views/dashboard/DashboardView.vue`

Temporary placeholder using shadcn `Card` and `UnderDevelopment` component. This will be replaced in Phase 5.

### Phase 1 Deliverable

The app boots, login works against the existing backend, authenticated users land on a stub dashboard, and no broken imports remain.

---

## Phase 2 — Rooms Module (Admin)

**Goal:** Admins can view, create, edit, and toggle availability of rooms.

### 2.1 Types

File: `src/types/room.ts`

```typescript
export type RoomType = 'single' | 'double' | 'suite' | 'conference'

export interface Room {
  id: number
  name: string
  type: RoomType
  capacity: number
  price_per_night: number
  amenities: string[]
  is_available: boolean
  images: string[]
  created_at: string
}

export interface RoomPayload {
  name: string
  type: RoomType
  capacity: number
  price_per_night: number
  amenities: string[]
  is_available: boolean
}
```

### 2.2 API Service

File: `src/services/api/room.ts`

```typescript
import { client } from './client'
import type { Room, RoomPayload } from '@/types/room'

export const roomApi = {
  list: () => client.get<Room[]>('/rooms'),
  get: (id: number) => client.get<Room>(`/rooms/${id}`),
  create: (payload: RoomPayload) => client.post<Room>('/rooms', payload),
  update: (id: number, payload: Partial<RoomPayload>) => client.put<Room>(`/rooms/${id}`, payload),
  toggleAvailability: (id: number, available: boolean) =>
    client.patch<Room>(`/rooms/${id}`, { is_available: available }),
  delete: (id: number) => client.delete(`/rooms/${id}`),
}
```

### 2.3 Pinia Store

File: `src/stores/rooms.ts`

State: `rooms[]`, `loading`, `error`
Actions: `fetchRooms`, `createRoom`, `updateRoom`, `toggleAvailability`, `deleteRoom`

### 2.4 Components

- `src/components/rooms/RoomDialog.vue` — Create/edit form using shadcn `Dialog`, `Input`, `Select`, `Switch`, `Checkbox`
- `src/components/rooms/RoomCard.vue` — Card display of a room (image, name, type badge, price, availability switch)

### 2.5 Views

- `src/views/rooms/RoomsView.vue` — TanStack table with columns: Name, Type, Capacity, Price, Available (Switch), Actions (Edit / Delete)

### 2.6 Routes

```typescript
{ path: '/rooms', component: RoomsView, meta: { requiresAuth: true } }
```

### 2.7 Sidebar

Add "Rooms" nav item under "Management" group (admin only).

### Phase 2 Deliverable

Admin can fully manage the room inventory.

---

## Phase 3 — Client Registration & Profile

**Goal:** Clients can register (individual or corporate), log in, and manage their profile.

### 3.1 Types

File: `src/types/client.ts`

```typescript
export interface IndividualProfile {
  id: number
  user_id: number
  full_name: string
  phone: string
  id_number: string
}

export interface CorporateProfile {
  id: number
  user_id: number
  company_name: string
  tax_id: string
  contact_person: string
  phone: string
  cost_centre?: string
}

export type ClientProfile = IndividualProfile | CorporateProfile
```

### 3.2 Registration View

File: `src/views/auth/RegisterView.vue`

- Step 1: Account type selection (Individual / Corporate) using shadcn `Card` toggle group
- Step 2: Account credentials (email, password)
- Step 3: Profile details (conditional fields based on account type) using shadcn `Stepper` component
- Uses `useResultDialog` composable on success/error

### 3.3 Profile View

File: `src/views/profile/ProfileView.vue`

- Rework existing `ProfileView.vue`
- Show editable profile fields depending on role
- Corporate clients see company billing fields

### 3.4 API Service

File: `src/services/api/client.ts`

```typescript
export const clientApi = {
  register: (payload) => client.post('/auth/register', payload),
  getProfile: () => client.get('/clients/me'),
  updateProfile: (payload) => client.put('/clients/me', payload),
}
```

### 3.5 Routes

```typescript
{ path: '/register', component: RegisterView, meta: { guestOnly: true } }
{ path: '/profile', component: ProfileView, meta: { requiresAuth: true } }
```

### Phase 3 Deliverable

New users can self-register, and authenticated clients can view and update their profile.

---

## Phase 4 — Booking Flow (Client)

**Goal:** Clients can search available rooms and complete a multi-step booking.

### 4.1 Types

File: `src/types/booking.ts`

```typescript
export type BookingStatus = 'pending' | 'confirmed' | 'cancelled' | 'completed'

export interface Booking {
  id: number
  user_id: number
  room_id: number
  room: Room
  check_in: string      // ISO date
  check_out: string     // ISO date
  guests: number
  status: BookingStatus
  meal_ids: number[]
  meals: Meal[]
  special_requests?: string
  invoice_id?: number
  created_at: string
}

export interface BookingPayload {
  room_id: number
  check_in: string
  check_out: string
  guests: number
  meal_ids: number[]
  special_requests?: string
}
```

### 4.2 Meal Types

File: `src/types/meal.ts`

```typescript
export type MealType = 'breakfast' | 'lunch' | 'dinner'

export interface Meal {
  id: number
  name: string
  type: MealType
  price: number
  description: string
  is_available: boolean
}
```

### 4.3 Booking Wizard Component

File: `src/components/bookings/BookingWizard.vue`

A 4-step flow using the existing shadcn `Stepper` component:

| Step | Content | Components Used |
|---|---|---|
| 1. Dates & Guests | DateRangePicker (RangeCalendar), guest count Input | `RangeCalendar`, `Input`, `Label` |
| 2. Room Selection | Filterable grid of `RoomCard` components | `Card`, `Badge`, `Button`, `Select` |
| 3. Meal Add-ons | Checkbox list of available meals by type | `Checkbox`, `Badge`, `Separator` |
| 4. Confirm & Book | Summary of selections, special requests Textarea, submit | `Card`, `Textarea`, `Button` |

### 4.4 API Services

File: `src/services/api/booking.ts`

```typescript
export const bookingApi = {
  checkAvailability: (checkIn: string, checkOut: string) =>
    client.get<Room[]>('/rooms/available', { check_in: checkIn, check_out: checkOut }),
  create: (payload: BookingPayload) => client.post<Booking>('/bookings', payload),
  list: () => client.get<Booking[]>('/bookings/my'),
  get: (id: number) => client.get<Booking>(`/bookings/${id}`),
  cancel: (id: number) => client.patch(`/bookings/${id}/cancel`),
}

export const mealApi = {
  list: () => client.get<Meal[]>('/meals'),
}
```

### 4.5 Pinia Store

File: `src/stores/bookings.ts`

State: `bookings[]`, `currentBooking`, `availableRooms[]`, `meals[]`, `loading`, `error`
Actions: `fetchMyBookings`, `checkAvailability`, `createBooking`, `cancelBooking`, `fetchMeals`

### 4.6 Views

- `src/views/bookings/NewBookingView.vue` — Hosts `BookingWizard.vue`
- `src/views/bookings/MyBookingsView.vue` — TanStack table of client's bookings with status badges and cancel action

### 4.7 Routes

```typescript
{ path: '/book', component: NewBookingView, meta: { requiresAuth: true } }
{ path: '/my-bookings', component: MyBookingsView, meta: { requiresAuth: true } }
{ path: '/my-bookings/:id', component: BookingDetailView, meta: { requiresAuth: true } }
```

### 4.8 Sidebar

Add "Book a Room" and "My Bookings" nav items to the client navigation group (hidden for admins).

### Phase 4 Deliverable

Clients can search for rooms by date, select a room, add meals, and complete a booking.

---

## Phase 5 — Admin Booking Management

**Goal:** Admins can view all bookings, confirm/cancel them, and manage meals.

### 5.1 Components

- `src/components/bookings/BookingActionDialog.vue` — Confirm or cancel a booking with optional notes (reuse pattern from `TaskActionDialog.vue`)
- `src/components/bookings/BookingDetailSheet.vue` — Slide-in `Sheet` showing full booking details

### 5.2 Admin Booking View

File: `src/views/admin/BookingsView.vue`

TanStack table columns: Booking ID, Client, Room, Check-in, Check-out, Guests, Status (Badge), Actions (Confirm / Cancel / View)

Filters: Status dropdown, date range, search by client name

### 5.3 Meal Management

File: `src/views/admin/MealsView.vue`

- Table of meals with type badges and availability switches
- `MealDialog.vue` for create/edit (name, type, price, description, availability)

API: `mealApi.list`, `mealApi.create`, `mealApi.update`, `mealApi.toggleAvailability`, `mealApi.delete`

### 5.4 Bookings API Extension (Admin)

Add to `src/services/api/booking.ts`:

```typescript
admin: {
  list: (params?) => client.get<Booking[]>('/admin/bookings', params),
  confirm: (id: number) => client.patch(`/admin/bookings/${id}/confirm`),
  cancel: (id: number, reason?: string) => client.patch(`/admin/bookings/${id}/cancel`, { reason }),
}
```

### 5.5 Routes

```typescript
{ path: '/admin/bookings', component: AdminBookingsView, meta: { requiresAuth: true, role: 'admin' } }
{ path: '/admin/meals', component: MealsView, meta: { requiresAuth: true, role: 'admin' } }
```

### 5.6 Sidebar Update

Add admin-only group: "Bookings", "Meals"

### Phase 5 Deliverable

Admins have full booking oversight and meal catalogue management.

---

## Phase 6 — Invoice Module

**Goal:** Invoices are auto-generated on booking confirmation and downloadable as PDF.

### 6.1 Types

File: `src/types/invoice.ts`

```typescript
export type InvoiceStatus = 'draft' | 'issued' | 'paid' | 'overdue'

export interface InvoiceLineItem {
  description: string
  quantity: number
  unit_price: number
  total: number
}

export interface Invoice {
  id: number
  booking_id: number
  booking: Booking
  client: AuthUser
  line_items: InvoiceLineItem[]
  subtotal: number
  tax: number
  total: number
  status: InvoiceStatus
  issued_at: string
  due_date: string
}
```

### 6.2 API Service

File: `src/services/api/invoice.ts`

```typescript
export const invoiceApi = {
  getByBooking: (bookingId: number) => client.get<Invoice>(`/bookings/${bookingId}/invoice`),
  myInvoices: () => client.get<Invoice[]>('/invoices/my'),
  adminList: (params?) => client.get<Invoice[]>('/admin/invoices', params),
  markPaid: (id: number) => client.patch(`/admin/invoices/${id}/mark-paid`),
}
```

### 6.3 Invoice Document Component

File: `src/components/invoices/InvoiceDocument.vue`

Reuse the pattern from `PayslipDocument.vue` (print/PDF layout):

- Lodge header (name, address, logo)
- Client details
- Booking summary (room, dates, nights)
- Line items table (room cost, meal costs, extras)
- Subtotal, tax, total
- Payment status badge
- Download button using `@ceereals/vue-pdf` (already installed as dev dep)

### 6.4 Views

- `src/views/invoices/MyInvoicesView.vue` — Client's invoice list with download buttons
- `src/views/admin/InvoicesView.vue` — Admin invoice list with status filters and "Mark Paid" action

### 6.5 Pinia Store

File: `src/stores/invoices.ts`

State: `invoices[]`, `currentInvoice`, `loading`
Actions: `fetchMyInvoices`, `fetchAdminInvoices`, `getByBooking`, `markPaid`

### 6.6 Routes

```typescript
{ path: '/my-invoices', component: MyInvoicesView, meta: { requiresAuth: true } }
{ path: '/admin/invoices', component: AdminInvoicesView, meta: { requiresAuth: true, role: 'admin' } }
```

### Phase 6 Deliverable

Clients can view and download their invoices. Admins can manage invoice status.

---

## Phase 7 — Dashboard (Role-Aware)

**Goal:** Replace stub dashboard with role-specific summary views.

### 7.1 Admin Dashboard

File: `src/components/dashboard/AdminDashboard.vue`

KPI Cards (shadcn `Card`):
- Today's check-ins / check-outs
- Rooms occupied vs available
- Revenue this month
- Pending booking confirmations

Charts (using existing `@unovis` setup):
- Occupancy heatmap (room × day grid)
- Revenue bar chart (monthly)
- Meal demand (most ordered meals)

### 7.2 Client Dashboard

File: `src/components/dashboard/ClientDashboard.vue`

- Active bookings summary cards
- Quick action: "Book a Room" button
- Upcoming check-in countdown
- Outstanding invoices alert

### 7.3 Dashboard View

File: `src/views/dashboard/DashboardView.vue`

```vue
<template>
  <AdminDashboard v-if="authStore.userRole === 'admin'" />
  <ClientDashboard v-else />
</template>
```

### 7.4 API Service

File: `src/services/api/dashboard.ts`

```typescript
export const dashboardApi = {
  adminSummary: () => client.get('/admin/dashboard'),
  clientSummary: () => client.get('/dashboard'),
}
```

### Phase 7 Deliverable

All users land on a meaningful dashboard after login.

---

## Phase 8 — Admin System Management

**Goal:** Admins can manage system users and their roles.

### 8.1 System Users View

File: `src/views/admin/SystemUsersView.vue`

- Table of all registered users with role badges
- `ChangeRoleDialog.vue` — reuse and rework existing `ChangeRoleDialog` from HR system, updating role options to lodge roles
- Activate / deactivate toggle

### 8.2 Types

File: `src/types/user.ts`

```typescript
export interface SystemUser {
  id: number
  email: string
  role: UserRole
  full_name: string
  is_active: boolean
  created_at: string
}
```

### 8.3 API Service

File: `src/services/api/systemUser.ts`

```typescript
export const systemUserApi = {
  list: () => client.get<SystemUser[]>('/admin/users'),
  updateRole: (id: number, role: UserRole) => client.patch(`/admin/users/${id}/role`, { role }),
  toggleActive: (id: number, active: boolean) =>
    client.patch(`/admin/users/${id}/active`, { is_active: active }),
}
```

### 8.4 Routes & Sidebar

```typescript
{ path: '/admin/users', component: SystemUsersView, meta: { requiresAuth: true, role: 'admin' } }
```

Sidebar: Add "Users" under admin group.

### Phase 8 Deliverable

Admins can manage user accounts and roles.

---

## Phase 9 — Reporting

**Goal:** Admins see visual reports on occupancy, revenue, and meal demand.

### 9.1 Report Views

File: `src/views/admin/ReportsView.vue`

Tabs (shadcn `Tabs`):
- **Occupancy** — Room × date heatmap using `@unovis` `HeatmapChart`
- **Revenue** — Monthly revenue bar chart with date range filter
- **Meal Demand** — Top meals ordered (bar or pie chart)
- **Bookings** — Bookings by status over time (line chart)

### 9.2 API Service

File: `src/services/api/reports.ts`

```typescript
export const reportsApi = {
  occupancy: (from: string, to: string) => client.get('/admin/reports/occupancy', { from, to }),
  revenue: (from: string, to: string) => client.get('/admin/reports/revenue', { from, to }),
  mealDemand: (from: string, to: string) => client.get('/admin/reports/meals', { from, to }),
}
```

### 9.3 Routes & Sidebar

```typescript
{ path: '/admin/reports', component: ReportsView, meta: { requiresAuth: true, role: 'admin' } }
```

Sidebar: Add "Reports" under admin group.

### Phase 9 Deliverable

Admins have a fully functional analytics dashboard.

---

## Final Router Structure

After all phases, `src/router/index.ts` will contain:

```typescript
// Public
{ path: '/login', guestOnly: true }
{ path: '/register', guestOnly: true }

// Authenticated (all roles)
{ path: '/dashboard', requiresAuth: true }
{ path: '/profile', requiresAuth: true }

// Client only
{ path: '/book', requiresAuth: true }
{ path: '/my-bookings', requiresAuth: true }
{ path: '/my-bookings/:id', requiresAuth: true }
{ path: '/my-invoices', requiresAuth: true }

// Admin only
{ path: '/rooms', requiresAuth: true, role: 'admin' }
{ path: '/admin/bookings', requiresAuth: true, role: 'admin' }
{ path: '/admin/meals', requiresAuth: true, role: 'admin' }
{ path: '/admin/invoices', requiresAuth: true, role: 'admin' }
{ path: '/admin/users', requiresAuth: true, role: 'admin' }
{ path: '/admin/reports', requiresAuth: true, role: 'admin' }
```

---

## Final Sidebar Structure

```
Lodge Admin (admin role)
├── Overview
│   └── Dashboard
├── Management
│   ├── Rooms
│   ├── Bookings
│   ├── Meals
│   └── Invoices
├── Reports
│   └── Reports
└── System
    └── Users

Client (individual / corporate)
├── Overview
│   └── Dashboard
├── Bookings
│   ├── Book a Room
│   └── My Bookings
└── Billing
    └── My Invoices
```

---

## Reuse Map — Existing → Lodge Equivalent

| Existing HR File | Lodge Equivalent | Action |
|---|---|---|
| `TaskActionDialog.vue` | `BookingActionDialog.vue` | Copy & adapt |
| `ChangeRoleDialog.vue` | `ChangeRoleDialog.vue` | Update role options |
| `EmployeeDialog.vue` | `RoomDialog.vue` | Rewrite with room fields |
| `PayslipDocument.vue` | `InvoiceDocument.vue` | Rewrite with invoice fields |
| `LeaveCalendar.vue` | `BookingWizard` step 1 | RangeCalendar already installed |
| `AdminDashboard.vue` | `AdminDashboard.vue` | Rewrite KPIs + charts |
| `EmployeeDashboard.vue` | `ClientDashboard.vue` | Rewrite |
| `DepartmentsView.vue` (table pattern) | `RoomsView.vue` | Copy table shell, replace columns |
| `EmployeesView.vue` (table pattern) | `AdminBookingsView.vue` | Copy table shell, replace columns |
| `stores/workflow.ts` (store pattern) | `stores/rooms.ts`, `stores/bookings.ts` | Rewrite following same pattern |

---

## Phase Sequence Summary

| Phase | Feature | Effort |
|---|---|---|
| 1 | Foundation & Auth Cleanup | Low |
| 2 | Rooms Module (Admin) | Medium |
| 3 | Client Registration & Profile | Medium |
| 4 | Booking Flow (Client) | High |
| 5 | Admin Booking Management | Medium |
| 6 | Invoice Module | Medium |
| 7 | Dashboard (Role-Aware) | Medium |
| 8 | Admin System Management | Low |
| 9 | Reporting | Medium |

Start with Phase 1 before touching any feature phase. Phases 2–3 can proceed in parallel once Phase 1 is stable. Phase 4 depends on Phase 2 (rooms) and Phase 3 (auth roles). Phases 5–6 depend on Phase 4. Phases 7–9 depend on all prior phases.
