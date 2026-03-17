# Lodge Administration Web Application — Vue.js Development Guide

## Table of Contents

1. [System Analysis](#1-system-analysis)
2. [Process Flows](#2-process-flows)
3. [Application Architecture](#3-application-architecture)
4. [Project Setup](#4-project-setup)
5. [Data Models](#5-data-models)
6. [Routing & Navigation](#6-routing--navigation)
7. [State Management (Pinia)](#7-state-management-pinia)
8. [Page-by-Page Implementation](#8-page-by-page-implementation)
9. [Component Library](#9-component-library)
10. [Authentication & Authorization](#10-authentication--authorization)
11. [Invoice Generation](#11-invoice-generation)
12. [Reporting Dashboard](#12-reporting-dashboard)
13. [API Contract (Expected Backend)](#13-api-contract-expected-backend)
14. [Deployment Checklist](#14-deployment-checklist)

---

## 1. System Analysis

### 1.1 User Roles

The application serves two distinct user groups, each with separate interfaces and permissions.

**Lodge Administrator** — a single or small team of staff members who manage the day-to-day lodge operations. They have full access to rooms, bookings, meals, invoices, and reporting.

**Clients** — external users who book accommodation and meals. They come in two flavours:

- **Individual clients** who book for themselves.
- **Corporate clients** who book on behalf of employees. Corporate accounts carry additional billing fields (company name, tax ID, cost centre, etc.).

### 1.2 Core Entities

| Entity | Description |
|--------|-------------|
| **User / Account** | Authentication identity. Linked to either an individual profile or a corporate profile. |
| **Room** | A physical room with a type, capacity, price, and availability flag. |
| **Booking** | A reservation tying a client (or company + guest list) to one or more rooms for a date range. |
| **Meal** | A meal request tied to a booking, specifying type (breakfast/lunch/dinner), date, quantity, and dietary notes. |
| **Invoice** | An auto-generated financial document summarising room charges, meal charges, and totals. |
| **Report** | Admin-only aggregated views: occupancy, revenue, meal demand. |

### 1.3 Key Business Rules

1. **No double-booking** — a room cannot be assigned to two overlapping date ranges.
2. **Meals are optional** — guests may add meals during or after booking.
3. **Invoices are generated automatically** once a booking is confirmed but may be regenerated if the booking changes.
4. **Corporate bookings** may list multiple guest names against a single company account.
5. **Only admins** can confirm/modify bookings and view reports.

---

## 2. Process Flows

### 2.1 Client Registration Flow

```
┌─────────────┐
│  Home Page   │
└──────┬──────┘
       │ clicks "Register"
       ▼
┌──────────────────┐
│ Registration Page │
│                  │
│ Step 1: Choose   │
│ account type     │
│ (Individual /    │
│  Corporate)      │
└──────┬───────────┘
       │
       ▼
┌──────────────────────────────────────────────┐
│ Step 2: Fill details                         │
│                                              │
│ Individual:                                  │
│   - Full name, email, phone, ID/passport,    │
│     address, emergency contact               │
│                                              │
│ Corporate:                                   │
│   - Company name, reg number, address,       │
│     phone, email                             │
│   - Primary contact (name, title, phone,     │
│     email)                                   │
│   - Billing info (contact, email, address,   │
│     tax ID)                                  │
└──────┬───────────────────────────────────────┘
       │
       ▼
┌──────────────────┐
│ Step 3: Create   │
│ login credentials│
│ (email + password│
│  or username)    │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Account created  │
│ → Redirect to    │
│   Client         │
│   Dashboard      │
└──────────────────┘
```

### 2.2 Client Booking Flow (Core Workflow)

This is the primary revenue-generating flow and the heart of the system.

```
┌────────────────────┐
│ Client logs in     │
└──────┬─────────────┘
       ▼
┌────────────────────┐
│ Client Dashboard   │
│ (overview of       │
│  active bookings,  │
│  upcoming stays,   │
│  quick actions)    │
└──────┬─────────────┘
       │ clicks "New Booking"
       ▼
┌────────────────────────────────────────────────────┐
│ STEP 1 — Select Dates                             │
│                                                    │
│  • Check-in date picker                            │
│  • Check-out date picker                           │
│  • Number of guests                                │
│  • (Corporate only) Select/add guest names         │
└──────┬─────────────────────────────────────────────┘
       │ dates validated (check-out > check-in)
       ▼
┌────────────────────────────────────────────────────┐
│ STEP 2 — Browse Available Rooms                    │
│                                                    │
│  System queries rooms NOT booked for those dates.  │
│  Display:                                          │
│  • Room number, type, capacity, price/night        │
│  • Photo / description                             │
│  • "Select" button per room                        │
│                                                    │
│  Filters: room type, max price, capacity           │
└──────┬─────────────────────────────────────────────┘
       │ room(s) selected
       ▼
┌────────────────────────────────────────────────────┐
│ STEP 3 — Meal Selection (Optional)                 │
│                                                    │
│  For each day of the stay, the guest can add:      │
│  • Breakfast  (quantity, dietary notes)             │
│  • Lunch      (quantity, dietary notes)             │
│  • Dinner     (quantity, dietary notes)             │
│                                                    │
│  Menu items shown from admin's daily menu.         │
│  Guest can skip this step entirely.                │
└──────┬─────────────────────────────────────────────┘
       │
       ▼
┌────────────────────────────────────────────────────┐
│ STEP 4 — Booking Summary & Confirmation            │
│                                                    │
│  • Room(s), dates, nightly rate                    │
│  • Meals selected with prices                      │
│  • Total estimated cost                            │
│  • Special requests text field                     │
│  • "Confirm Booking" button                        │
└──────┬─────────────────────────────────────────────┘
       │ confirmed
       ▼
┌────────────────────────────────────────────────────┐
│ SYSTEM ACTIONS (behind the scenes)                 │
│                                                    │
│  1. Create Booking record (status: PENDING)        │
│  2. Create Meal records linked to booking          │
│  3. Generate Invoice (status: UNPAID)              │
│  4. Notify admin of new booking                    │
└──────┬─────────────────────────────────────────────┘
       │
       ▼
┌────────────────────┐
│ Booking confirmed  │
│ screen with        │
│ booking reference, │
│ invoice link       │
└────────────────────┘
```

### 2.3 Admin Booking Management Flow

```
┌────────────────────────┐
│ Admin logs in          │
└──────┬─────────────────┘
       ▼
┌────────────────────────┐
│ Admin Dashboard        │
│ • Today's check-ins    │
│ • Today's check-outs   │
│ • Pending bookings     │
│ • Occupancy summary    │
│ • Meal summary         │
└──────┬─────────────────┘
       │ clicks "Bookings"
       ▼
┌────────────────────────────────────────────────────┐
│ Booking Management Page                            │
│                                                    │
│ Table/list of all bookings with filters:           │
│ • Status (Pending / Confirmed / Checked-in /       │
│          Checked-out / Cancelled)                  │
│ • Date range                                       │
│ • Client name / company                            │
│                                                    │
│ Actions per booking:                               │
│ • View details                                     │
│ • Confirm booking (Pending → Confirmed)            │
│ • Modify dates / room                              │
│ • Cancel booking                                   │
│ • Regenerate invoice                               │
└────────────────────────────────────────────────────┘
```

### 2.4 Room Management Flow (Admin)

```
Admin Dashboard → Room Management Page

┌───────────────────────────────────────────┐
│ Room Management                           │
│                                           │
│ ┌─────────┐  ┌─────────────────────────┐  │
│ │ Add Room │  │ Search / Filter rooms   │  │
│ └─────────┘  └─────────────────────────┘  │
│                                           │
│ Room list / grid:                         │
│ ┌────────────────────────────────────┐    │
│ │ Room 101 │ Single │ K50/night     │    │
│ │ Status: Available │ Capacity: 1   │    │
│ │ [Edit] [Toggle Availability]      │    │
│ └────────────────────────────────────┘    │
│                                           │
│ Add/Edit Room form:                       │
│ • Room number                             │
│ • Room type (Single/Double/Suite/etc.)    │
│ • Capacity                                │
│ • Price per night                         │
│ • Description                             │
│ • Available (yes/no toggle)               │
└───────────────────────────────────────────┘
```

### 2.5 Meal Management Flow (Admin)

```
Admin Dashboard → Meal Management Page

┌───────────────────────────────────────────────────┐
│ Meal Management                                   │
│                                                   │
│ Tab 1: Daily Menu                                 │
│ ┌───────────────────────────────────────────────┐ │
│ │ Select Date: [date picker]                    │ │
│ │                                               │ │
│ │ Breakfast items:  [Add / Edit / Remove]       │ │
│ │ Lunch items:      [Add / Edit / Remove]       │ │
│ │ Dinner items:     [Add / Edit / Remove]       │ │
│ │                                               │ │
│ │ Each item: name, description, price,          │ │
│ │            available (toggle)                  │ │
│ └───────────────────────────────────────────────┘ │
│                                                   │
│ Tab 2: Meal Requests                              │
│ ┌───────────────────────────────────────────────┐ │
│ │ Filter by date / meal type                    │ │
│ │                                               │ │
│ │ Table:                                        │ │
│ │ Guest Name | Room | Meal Type | Qty | Diet    │ │
│ │ ─────────────────────────────────────────     │ │
│ │ J. Smith   | 101  | Breakfast | 1   | None    │ │
│ │ A. Banda   | 205  | Dinner    | 2   | Vegan   │ │
│ │                                               │ │
│ │ Summary: Breakfast: 12, Lunch: 8, Dinner: 15  │ │
│ └───────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────┘
```

### 2.6 Invoice Lifecycle

```
Booking Confirmed
       │
       ▼
┌──────────────────────────────────┐
│ Invoice auto-generated           │
│ Status: UNPAID                   │
│                                  │
│ Contains:                        │
│ • Invoice number (auto)          │
│ • Client / company name          │
│ • Dept / cost centre (corporate) │
│ • Room charges (rate × nights)   │
│ • Meal charges (itemised)        │
│ • Total amount due               │
│ • Date                           │
└──────────┬───────────────────────┘
           │
     ┌─────┴─────┐
     ▼           ▼
  Client       Admin
  views /      views /
  downloads    manages
  invoice      invoice
               │
               ▼
        ┌──────────────┐
        │ Mark as PAID │
        │ or generate  │
        │ statement    │
        └──────────────┘
```

### 2.7 Reporting Flow (Admin)

```
Admin Dashboard → Reports Page

Available reports:
┌─────────────────────────────────────────┐
│ 1. Daily Bookings                       │
│    → Select date → table of bookings    │
│                                         │
│ 2. Monthly Revenue                      │
│    → Select month → bar/line chart      │
│    → Breakdown: rooms vs meals          │
│                                         │
│ 3. Occupancy Rate                       │
│    → Select date range → percentage     │
│    → Visual calendar heatmap            │
│                                         │
│ 4. Meal Demand                          │
│    → Select date range → chart showing  │
│      breakfast / lunch / dinner counts  │
└─────────────────────────────────────────┘
```

---

## 3. Application Architecture

### 3.1 Technology Stack

| Layer | Technology |
|-------|-----------|
| Framework | Vue 3 (Composition API + `<script setup>`) |
| Build tool | Vite |
| Routing | Vue Router 4 |
| State management | Pinia |
| UI components | PrimeVue or Vuetify 3 (or custom with Tailwind CSS) |
| HTTP client | Axios |
| Form validation | VeeValidate + Yup |
| Charts | Chart.js via vue-chartjs |
| PDF generation | jsPDF (client-side invoice download) |
| Date handling | date-fns |
| Authentication | JWT tokens stored in httpOnly cookies (preferred) or localStorage |

### 3.2 Folder Structure

```
src/
├── assets/                  # Static images, global CSS
├── components/
│   ├── common/              # Shared UI: AppHeader, AppSidebar, DataTable, Modal, etc.
│   ├── booking/             # BookingWizard, RoomCard, DateRangePicker
│   ├── meals/               # MealSelector, MenuEditor, MealRequestTable
│   ├── invoices/            # InvoicePreview, InvoiceTable
│   └── reports/             # RevenueChart, OccupancyHeatmap, MealDemandChart
├── composables/             # Reusable logic (useAuth, useBooking, useRooms, etc.)
├── layouts/
│   ├── PublicLayout.vue     # Header + footer, no sidebar
│   ├── ClientLayout.vue     # Header + client sidebar
│   └── AdminLayout.vue      # Header + admin sidebar
├── pages/                   # One file per route (or nested folders)
│   ├── public/
│   │   ├── HomePage.vue
│   │   ├── LoginPage.vue
│   │   └── RegisterPage.vue
│   ├── client/
│   │   ├── DashboardPage.vue
│   │   ├── BookingPage.vue
│   │   ├── MealSelectionPage.vue
│   │   ├── BookingHistoryPage.vue
│   │   ├── InvoicePage.vue
│   │   └── ProfilePage.vue
│   └── admin/
│       ├── DashboardPage.vue
│       ├── RoomManagementPage.vue
│       ├── BookingManagementPage.vue
│       ├── MealManagementPage.vue
│       ├── InvoiceManagementPage.vue
│       └── ReportsPage.vue
├── router/
│   └── index.js             # Route definitions + guards
├── stores/
│   ├── auth.js
│   ├── rooms.js
│   ├── bookings.js
│   ├── meals.js
│   ├── invoices.js
│   └── reports.js
├── services/
│   ├── api.js               # Axios instance with interceptors
│   ├── authService.js
│   ├── roomService.js
│   ├── bookingService.js
│   ├── mealService.js
│   ├── invoiceService.js
│   └── reportService.js
├── utils/
│   ├── constants.js          # Enums, status labels, config
│   ├── validators.js         # Yup schemas
│   └── formatters.js         # Currency, date formatting helpers
├── App.vue
└── main.js
```

---

## 4. Project Setup

### 4.1 Scaffold the project

```bash
npm create vite@latest lodge-admin -- --template vue
cd lodge-admin
npm install
```

### 4.2 Install dependencies

```bash
# Core
npm install vue-router@4 pinia axios

# UI (choose one)
npm install primevue primeicons
# OR
npm install vuetify

# Forms
npm install vee-validate yup

# Utilities
npm install date-fns chart.js vue-chartjs jspdf
```

### 4.3 Configure main.js

```javascript
// src/main.js
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// If using PrimeVue:
import PrimeVue from 'primevue/config'
import 'primevue/resources/themes/lara-light-blue/theme.css'
import 'primeicons/primeicons.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(PrimeVue)

app.mount('#app')
```

---

## 5. Data Models

These TypeScript-style interfaces describe the shape of data flowing through the frontend. Even if you write plain JavaScript, use these as a reference.

### 5.1 User & Account

```typescript
interface User {
  id: string
  email: string
  role: 'admin' | 'individual' | 'corporate'
  profile: IndividualProfile | CorporateProfile
  createdAt: string
}

interface IndividualProfile {
  fullName: string
  nationalId?: string
  phone: string
  address: string
  nationality?: string
  emergencyContactName?: string
  emergencyContactPhone?: string
}

interface CorporateProfile {
  companyName: string
  registrationNumber?: string
  companyAddress: string
  companyPhone: string
  companyEmail: string
  primaryContact: {
    fullName: string
    jobTitle: string
    phone: string
    email: string
  }
  billing: {
    contactPerson: string
    email: string
    address: string
    taxId?: string
  }
  authorizedEmployees: Employee[]
}

interface Employee {
  id: string
  fullName: string
  email?: string
  phone?: string
}
```

### 5.2 Room

```typescript
interface Room {
  id: string
  roomNumber: string
  roomType: 'single' | 'double' | 'twin' | 'suite' | 'dormitory'
  capacity: number
  pricePerNight: number
  description: string
  isAvailable: boolean
  imageUrl?: string
}
```

### 5.3 Booking

```typescript
interface Booking {
  id: string
  bookingReference: string
  clientId: string
  clientType: 'individual' | 'corporate'
  companyName?: string
  department?: string
  costCentre?: string
  guests: GuestInfo[]
  roomId: string
  roomNumber: string
  checkIn: string        // ISO date
  checkOut: string       // ISO date
  numberOfGuests: number
  specialRequests?: string
  status: 'pending' | 'confirmed' | 'checked_in' | 'checked_out' | 'cancelled'
  meals: MealRequest[]
  invoiceId?: string
  createdAt: string
}

interface GuestInfo {
  fullName: string
  email?: string
  phone?: string
}
```

### 5.4 Meal

```typescript
interface MenuItem {
  id: string
  name: string
  description: string
  price: number
  mealType: 'breakfast' | 'lunch' | 'dinner'
  isAvailable: boolean
  date: string           // The date this menu item is offered
}

interface MealRequest {
  id: string
  bookingId: string
  guestName: string
  mealType: 'breakfast' | 'lunch' | 'dinner'
  date: string
  quantity: number
  dietaryRequirements?: string
  menuItemId?: string
  price: number
}
```

### 5.5 Invoice

```typescript
interface Invoice {
  id: string
  invoiceNumber: string
  bookingReference: string
  clientName: string
  companyName?: string
  department?: string
  costCentre?: string
  items: InvoiceLineItem[]
  totalAmount: number
  status: 'unpaid' | 'paid'
  invoiceDate: string
  dueDate?: string
}

interface InvoiceLineItem {
  description: string    // e.g. "Room 101 — Single (3 nights)"
  unitPrice: number
  quantity: number
  total: number
}
```

---

## 6. Routing & Navigation

### 6.1 Route Definitions

```javascript
// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Lazy-loaded pages
const HomePage         = () => import('@/pages/public/HomePage.vue')
const LoginPage        = () => import('@/pages/public/LoginPage.vue')
const RegisterPage     = () => import('@/pages/public/RegisterPage.vue')

const ClientDashboard  = () => import('@/pages/client/DashboardPage.vue')
const BookingPage      = () => import('@/pages/client/BookingPage.vue')
const MealSelection    = () => import('@/pages/client/MealSelectionPage.vue')
const BookingHistory   = () => import('@/pages/client/BookingHistoryPage.vue')
const InvoicePage      = () => import('@/pages/client/InvoicePage.vue')
const ProfilePage      = () => import('@/pages/client/ProfilePage.vue')

const AdminDashboard   = () => import('@/pages/admin/DashboardPage.vue')
const RoomManagement   = () => import('@/pages/admin/RoomManagementPage.vue')
const BookingMgmt      = () => import('@/pages/admin/BookingManagementPage.vue')
const MealManagement   = () => import('@/pages/admin/MealManagementPage.vue')
const InvoiceMgmt      = () => import('@/pages/admin/InvoiceManagementPage.vue')
const ReportsPage      = () => import('@/pages/admin/ReportsPage.vue')

const routes = [
  // ── Public ──
  { path: '/',          name: 'home',     component: HomePage,     meta: { layout: 'public' } },
  { path: '/login',     name: 'login',    component: LoginPage,    meta: { layout: 'public' } },
  { path: '/register',  name: 'register', component: RegisterPage, meta: { layout: 'public' } },

  // ── Client ──
  {
    path: '/client',
    meta: { layout: 'client', requiresAuth: true, role: 'client' },
    children: [
      { path: '',               name: 'client-dashboard', component: ClientDashboard },
      { path: 'book',           name: 'client-book',      component: BookingPage },
      { path: 'meals/:bookingId', name: 'client-meals',   component: MealSelection, props: true },
      { path: 'history',        name: 'client-history',   component: BookingHistory },
      { path: 'invoices',       name: 'client-invoices',  component: InvoicePage },
      { path: 'profile',        name: 'client-profile',   component: ProfilePage },
    ]
  },

  // ── Admin ──
  {
    path: '/admin',
    meta: { layout: 'admin', requiresAuth: true, role: 'admin' },
    children: [
      { path: '',           name: 'admin-dashboard', component: AdminDashboard },
      { path: 'rooms',      name: 'admin-rooms',     component: RoomManagement },
      { path: 'bookings',   name: 'admin-bookings',  component: BookingMgmt },
      { path: 'meals',      name: 'admin-meals',     component: MealManagement },
      { path: 'invoices',   name: 'admin-invoices',  component: InvoiceMgmt },
      { path: 'reports',    name: 'admin-reports',    component: ReportsPage },
    ]
  },

  // ── Catch-all ──
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// ── Navigation Guard ──
router.beforeEach((to, from, next) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next({ name: 'login', query: { redirect: to.fullPath } })
  }

  if (to.meta.role === 'admin' && auth.user?.role !== 'admin') {
    return next({ name: 'home' })
  }

  if (to.meta.role === 'client' && auth.user?.role === 'admin') {
    return next({ name: 'admin-dashboard' })
  }

  next()
})

export default router
```

### 6.2 Layout Switching in App.vue

```vue
<!-- src/App.vue -->
<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import PublicLayout from '@/layouts/PublicLayout.vue'
import ClientLayout from '@/layouts/ClientLayout.vue'
import AdminLayout from '@/layouts/AdminLayout.vue'

const route = useRoute()

const layoutComponent = computed(() => {
  const layout = route.meta.layout
  if (layout === 'admin') return AdminLayout
  if (layout === 'client') return ClientLayout
  return PublicLayout
})
</script>

<template>
  <component :is="layoutComponent">
    <RouterView />
  </component>
</template>
```

---

## 7. State Management (Pinia)

### 7.1 Auth Store

```javascript
// src/stores/auth.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import authService from '@/services/authService'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isCorpClient = computed(() => user.value?.role === 'corporate')

  async function login(credentials) {
    const response = await authService.login(credentials)
    token.value = response.token
    user.value = response.user
    localStorage.setItem('token', response.token)
  }

  async function register(data) {
    const response = await authService.register(data)
    token.value = response.token
    user.value = response.user
    localStorage.setItem('token', response.token)
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  async function fetchProfile() {
    if (!token.value) return
    user.value = await authService.getProfile()
  }

  return { user, token, isAuthenticated, isAdmin, isCorpClient, login, register, logout, fetchProfile }
})
```

### 7.2 Rooms Store

```javascript
// src/stores/rooms.js
import { defineStore } from 'pinia'
import { ref } from 'vue'
import roomService from '@/services/roomService'

export const useRoomStore = defineStore('rooms', () => {
  const rooms = ref([])
  const loading = ref(false)

  async function fetchRooms() {
    loading.value = true
    rooms.value = await roomService.getAll()
    loading.value = false
  }

  async function fetchAvailable(checkIn, checkOut) {
    loading.value = true
    rooms.value = await roomService.getAvailable(checkIn, checkOut)
    loading.value = false
  }

  async function createRoom(room) {
    const newRoom = await roomService.create(room)
    rooms.value.push(newRoom)
  }

  async function updateRoom(id, data) {
    const updated = await roomService.update(id, data)
    const idx = rooms.value.findIndex(r => r.id === id)
    if (idx !== -1) rooms.value[idx] = updated
  }

  async function toggleAvailability(id) {
    const room = rooms.value.find(r => r.id === id)
    if (room) await updateRoom(id, { isAvailable: !room.isAvailable })
  }

  return { rooms, loading, fetchRooms, fetchAvailable, createRoom, updateRoom, toggleAvailability }
})
```

### 7.3 Bookings Store

```javascript
// src/stores/bookings.js
import { defineStore } from 'pinia'
import { ref } from 'vue'
import bookingService from '@/services/bookingService'

export const useBookingStore = defineStore('bookings', () => {
  const bookings = ref([])
  const currentBooking = ref(null)
  const loading = ref(false)

  // Wizard state — holds in-progress booking data across steps
  const wizard = ref({
    step: 1,
    checkIn: null,
    checkOut: null,
    numberOfGuests: 1,
    guests: [],
    selectedRoom: null,
    meals: [],
    specialRequests: ''
  })

  function resetWizard() {
    wizard.value = {
      step: 1, checkIn: null, checkOut: null,
      numberOfGuests: 1, guests: [], selectedRoom: null,
      meals: [], specialRequests: ''
    }
  }

  async function fetchBookings(filters = {}) {
    loading.value = true
    bookings.value = await bookingService.getAll(filters)
    loading.value = false
  }

  async function fetchMyBookings() {
    loading.value = true
    bookings.value = await bookingService.getMine()
    loading.value = false
  }

  async function createBooking() {
    const payload = {
      checkIn: wizard.value.checkIn,
      checkOut: wizard.value.checkOut,
      numberOfGuests: wizard.value.numberOfGuests,
      guests: wizard.value.guests,
      roomId: wizard.value.selectedRoom.id,
      meals: wizard.value.meals,
      specialRequests: wizard.value.specialRequests
    }
    const booking = await bookingService.create(payload)
    bookings.value.unshift(booking)
    currentBooking.value = booking
    resetWizard()
    return booking
  }

  async function updateStatus(id, status) {
    const updated = await bookingService.updateStatus(id, status)
    const idx = bookings.value.findIndex(b => b.id === id)
    if (idx !== -1) bookings.value[idx] = updated
  }

  return {
    bookings, currentBooking, loading, wizard,
    resetWizard, fetchBookings, fetchMyBookings, createBooking, updateStatus
  }
})
```

---

## 8. Page-by-Page Implementation

### 8.1 Public — Home Page

A landing page introducing the lodge with a call-to-action to book or register.

```vue
<!-- src/pages/public/HomePage.vue -->
<script setup>
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
  <section class="hero">
    <h1>Welcome to [Lodge Name]</h1>
    <p>Comfortable accommodation and dining for individuals and corporate groups.</p>
    <div class="hero-actions">
      <button @click="router.push({ name: 'register' })">Create Account</button>
      <button @click="router.push({ name: 'login' })" class="secondary">Sign In</button>
    </div>
  </section>

  <section class="features">
    <!-- Room highlights, meal info, corporate booking info -->
  </section>
</template>
```

### 8.2 Client — Booking Page (Multi-Step Wizard)

This is the most complex client page. Use a stepper component to guide the user through the four steps outlined in the booking flow.

```vue
<!-- src/pages/client/BookingPage.vue -->
<script setup>
import { computed } from 'vue'
import { useBookingStore } from '@/stores/bookings'
import { useRoomStore } from '@/stores/rooms'
import { useRouter } from 'vue-router'

import StepDates from '@/components/booking/StepDates.vue'
import StepRooms from '@/components/booking/StepRooms.vue'
import StepMeals from '@/components/booking/StepMeals.vue'
import StepConfirm from '@/components/booking/StepConfirm.vue'

const bookingStore = useBookingStore()
const roomStore = useRoomStore()
const router = useRouter()

const step = computed(() => bookingStore.wizard.step)

function nextStep() {
  bookingStore.wizard.step++
}

function prevStep() {
  if (bookingStore.wizard.step > 1) bookingStore.wizard.step--
}

async function onDatesSelected() {
  // Fetch available rooms for chosen dates before advancing
  await roomStore.fetchAvailable(
    bookingStore.wizard.checkIn,
    bookingStore.wizard.checkOut
  )
  nextStep()
}

async function confirmBooking() {
  const booking = await bookingStore.createBooking()
  router.push({ name: 'client-history' })
}
</script>

<template>
  <div class="booking-wizard">
    <!-- Stepper indicator -->
    <div class="stepper">
      <span :class="{ active: step >= 1 }">1. Dates</span>
      <span :class="{ active: step >= 2 }">2. Room</span>
      <span :class="{ active: step >= 3 }">3. Meals</span>
      <span :class="{ active: step >= 4 }">4. Confirm</span>
    </div>

    <StepDates    v-if="step === 1" @next="onDatesSelected" />
    <StepRooms    v-if="step === 2" @next="nextStep" @back="prevStep" />
    <StepMeals    v-if="step === 3" @next="nextStep" @back="prevStep" />
    <StepConfirm  v-if="step === 4" @confirm="confirmBooking" @back="prevStep" />
  </div>
</template>
```

### 8.3 Admin — Dashboard Page

```vue
<!-- src/pages/admin/DashboardPage.vue -->
<script setup>
import { onMounted, ref } from 'vue'
import bookingService from '@/services/bookingService'
import roomService from '@/services/roomService'

const stats = ref({
  todayCheckIns: 0,
  todayCheckOuts: 0,
  pendingBookings: 0,
  occupancyRate: 0,
  totalMealsToday: 0
})

onMounted(async () => {
  const [bookingStats, roomStats] = await Promise.all([
    bookingService.getTodayStats(),
    roomService.getOccupancyStats()
  ])

  stats.value = {
    todayCheckIns: bookingStats.checkIns,
    todayCheckOuts: bookingStats.checkOuts,
    pendingBookings: bookingStats.pending,
    occupancyRate: roomStats.occupancyRate,
    totalMealsToday: bookingStats.totalMeals
  }
})
</script>

<template>
  <h1>Admin Dashboard</h1>

  <div class="stats-grid">
    <div class="stat-card">
      <span class="stat-value">{{ stats.todayCheckIns }}</span>
      <span class="stat-label">Check-ins Today</span>
    </div>
    <div class="stat-card">
      <span class="stat-value">{{ stats.todayCheckOuts }}</span>
      <span class="stat-label">Check-outs Today</span>
    </div>
    <div class="stat-card">
      <span class="stat-value">{{ stats.pendingBookings }}</span>
      <span class="stat-label">Pending Bookings</span>
    </div>
    <div class="stat-card">
      <span class="stat-value">{{ stats.occupancyRate }}%</span>
      <span class="stat-label">Occupancy Rate</span>
    </div>
    <div class="stat-card">
      <span class="stat-value">{{ stats.totalMealsToday }}</span>
      <span class="stat-label">Meals Today</span>
    </div>
  </div>

  <!-- Quick-access tables: recent bookings, upcoming check-ins -->
</template>
```

### 8.4 Admin — Room Management Page

```vue
<!-- src/pages/admin/RoomManagementPage.vue -->
<script setup>
import { onMounted, ref } from 'vue'
import { useRoomStore } from '@/stores/rooms'
import RoomFormModal from '@/components/common/RoomFormModal.vue'

const roomStore = useRoomStore()
const showModal = ref(false)
const editingRoom = ref(null)

onMounted(() => roomStore.fetchRooms())

function openAdd() {
  editingRoom.value = null
  showModal.value = true
}

function openEdit(room) {
  editingRoom.value = { ...room }
  showModal.value = true
}

async function saveRoom(formData) {
  if (editingRoom.value) {
    await roomStore.updateRoom(editingRoom.value.id, formData)
  } else {
    await roomStore.createRoom(formData)
  }
  showModal.value = false
}
</script>

<template>
  <div class="page-header">
    <h1>Room Management</h1>
    <button @click="openAdd">+ Add Room</button>
  </div>

  <table class="data-table">
    <thead>
      <tr>
        <th>Room #</th><th>Type</th><th>Capacity</th>
        <th>Price/Night</th><th>Status</th><th>Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="room in roomStore.rooms" :key="room.id">
        <td>{{ room.roomNumber }}</td>
        <td>{{ room.roomType }}</td>
        <td>{{ room.capacity }}</td>
        <td>{{ room.pricePerNight }}</td>
        <td>
          <span :class="room.isAvailable ? 'badge-green' : 'badge-red'">
            {{ room.isAvailable ? 'Available' : 'Unavailable' }}
          </span>
        </td>
        <td>
          <button @click="openEdit(room)">Edit</button>
          <button @click="roomStore.toggleAvailability(room.id)">
            {{ room.isAvailable ? 'Disable' : 'Enable' }}
          </button>
        </td>
      </tr>
    </tbody>
  </table>

  <RoomFormModal
    v-if="showModal"
    :room="editingRoom"
    @save="saveRoom"
    @close="showModal = false"
  />
</template>
```

---

## 9. Component Library

Key reusable components to build:

### 9.1 DateRangePicker

A wrapper around two date inputs enforcing check-out > check-in.

```vue
<!-- src/components/booking/DateRangePicker.vue -->
<script setup>
import { computed } from 'vue'

const props = defineProps({
  checkIn: String,
  checkOut: String
})

const emit = defineEmits(['update:checkIn', 'update:checkOut'])

const minCheckOut = computed(() => {
  if (!props.checkIn) return undefined
  const d = new Date(props.checkIn)
  d.setDate(d.getDate() + 1)
  return d.toISOString().split('T')[0]
})
</script>

<template>
  <div class="date-range">
    <label>
      Check-in
      <input
        type="date"
        :value="checkIn"
        :min="new Date().toISOString().split('T')[0]"
        @input="emit('update:checkIn', $event.target.value)"
      />
    </label>
    <label>
      Check-out
      <input
        type="date"
        :value="checkOut"
        :min="minCheckOut"
        @input="emit('update:checkOut', $event.target.value)"
      />
    </label>
  </div>
</template>
```

### 9.2 RoomCard

Displays a single room in the available-rooms grid.

```vue
<!-- src/components/booking/RoomCard.vue -->
<script setup>
defineProps({
  room: Object,
  selected: Boolean
})

defineEmits(['select'])
</script>

<template>
  <div class="room-card" :class="{ 'room-card--selected': selected }" @click="$emit('select', room)">
    <img v-if="room.imageUrl" :src="room.imageUrl" :alt="room.roomNumber" />
    <div class="room-card__info">
      <h3>Room {{ room.roomNumber }}</h3>
      <p class="room-type">{{ room.roomType }} — up to {{ room.capacity }} guest(s)</p>
      <p class="room-price">K{{ room.pricePerNight }} / night</p>
      <p class="room-desc">{{ room.description }}</p>
    </div>
  </div>
</template>
```

### 9.3 MealSelector

Allows a guest to pick meals for each day of their stay.

```vue
<!-- src/components/meals/MealSelector.vue -->
<script setup>
import { ref, computed } from 'vue'
import { eachDayOfInterval, format } from 'date-fns'

const props = defineProps({
  checkIn: String,
  checkOut: String,
  menuItems: Array   // All available MenuItem objects for the date range
})

const emit = defineEmits(['update:meals'])

const stayDays = computed(() =>
  eachDayOfInterval({
    start: new Date(props.checkIn),
    end: new Date(props.checkOut)
  }).slice(0, -1) // exclude checkout day
)

const selections = ref({})  // { 'YYYY-MM-DD': { breakfast: { qty, diet }, lunch: {...}, dinner: {...} } }

function toggleMeal(date, type) {
  const key = format(date, 'yyyy-MM-dd')
  if (!selections.value[key]) selections.value[key] = {}
  if (selections.value[key][type]) {
    delete selections.value[key][type]
  } else {
    selections.value[key][type] = { quantity: 1, dietaryRequirements: '' }
  }
  emitAll()
}

function emitAll() {
  const meals = []
  for (const [date, types] of Object.entries(selections.value)) {
    for (const [type, data] of Object.entries(types)) {
      meals.push({ date, mealType: type, ...data })
    }
  }
  emit('update:meals', meals)
}
</script>

<template>
  <div class="meal-selector">
    <div v-for="day in stayDays" :key="day" class="meal-day">
      <h4>{{ format(day, 'EEE, dd MMM yyyy') }}</h4>
      <div class="meal-options">
        <label v-for="type in ['breakfast', 'lunch', 'dinner']" :key="type">
          <input
            type="checkbox"
            :checked="selections[format(day, 'yyyy-MM-dd')]?.[type]"
            @change="toggleMeal(day, type)"
          />
          {{ type }}
        </label>
      </div>
    </div>
  </div>
</template>
```

---

## 10. Authentication & Authorization

### 10.1 Axios Interceptor

```javascript
// src/services/api.js
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' }
})

// Attach token to every request
api.interceptors.request.use(config => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

// Handle 401 globally
api.interceptors.response.use(
  res => res,
  error => {
    if (error.response?.status === 401) {
      const auth = useAuthStore()
      auth.logout()
      router.push({ name: 'login' })
    }
    return Promise.reject(error)
  }
)

export default api
```

### 10.2 Auth Service

```javascript
// src/services/authService.js
import api from './api'

export default {
  login: (credentials) => api.post('/auth/login', credentials).then(r => r.data),
  register: (data)     => api.post('/auth/register', data).then(r => r.data),
  getProfile: ()       => api.get('/auth/me').then(r => r.data),
  resetPassword: (email) => api.post('/auth/reset-password', { email }).then(r => r.data),
  updateProfile: (data) => api.put('/auth/profile', data).then(r => r.data),
}
```

### 10.3 Route Guard Logic Summary

| Route group | Requires auth? | Role check |
|-------------|---------------|------------|
| `/`, `/login`, `/register` | No | None |
| `/client/*` | Yes | Must NOT be admin |
| `/admin/*` | Yes | Must be admin |

---

## 11. Invoice Generation

### 11.1 Client-Side PDF Download

Use jsPDF to generate a downloadable invoice from the invoice data object.

```javascript
// src/utils/invoicePdf.js
import jsPDF from 'jspdf'

export function generateInvoicePdf(invoice) {
  const doc = new jsPDF()
  const margin = 20
  let y = margin

  // Header
  doc.setFontSize(20)
  doc.text('INVOICE', margin, y)
  y += 10

  doc.setFontSize(10)
  doc.text(`Invoice #: ${invoice.invoiceNumber}`, margin, y)
  doc.text(`Date: ${invoice.invoiceDate}`, 140, y)
  y += 7
  doc.text(`Booking Ref: ${invoice.bookingReference}`, margin, y)
  y += 7
  doc.text(`Client: ${invoice.clientName}`, margin, y)

  if (invoice.companyName) {
    y += 7
    doc.text(`Company: ${invoice.companyName}`, margin, y)
  }
  if (invoice.department) {
    y += 7
    doc.text(`Department / Cost Centre: ${invoice.department} / ${invoice.costCentre || '—'}`, margin, y)
  }

  y += 14

  // Line items table header
  doc.setFontSize(10)
  doc.setFont(undefined, 'bold')
  doc.text('Description', margin, y)
  doc.text('Qty', 110, y)
  doc.text('Unit Price', 130, y)
  doc.text('Total', 165, y)
  y += 2
  doc.line(margin, y, 190, y)
  y += 6
  doc.setFont(undefined, 'normal')

  // Line items
  for (const item of invoice.items) {
    doc.text(item.description, margin, y)
    doc.text(String(item.quantity), 110, y)
    doc.text(`K${item.unitPrice.toFixed(2)}`, 130, y)
    doc.text(`K${item.total.toFixed(2)}`, 165, y)
    y += 7
  }

  // Total
  y += 4
  doc.line(margin, y, 190, y)
  y += 8
  doc.setFont(undefined, 'bold')
  doc.setFontSize(12)
  doc.text(`Total Due: K${invoice.totalAmount.toFixed(2)}`, 130, y)

  // Status
  y += 10
  doc.setFontSize(10)
  doc.text(`Payment Status: ${invoice.status.toUpperCase()}`, margin, y)

  doc.save(`Invoice_${invoice.invoiceNumber}.pdf`)
}
```

---

## 12. Reporting Dashboard

### 12.1 Revenue Chart Component

```vue
<!-- src/components/reports/RevenueChart.vue -->
<script setup>
import { ref, onMounted } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS, CategoryScale, LinearScale,
  BarElement, Title, Tooltip, Legend
} from 'chart.js'
import reportService from '@/services/reportService'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

const chartData = ref(null)

onMounted(async () => {
  const data = await reportService.getMonthlyRevenue()

  chartData.value = {
    labels: data.map(d => d.month),
    datasets: [
      {
        label: 'Room Revenue',
        data: data.map(d => d.roomRevenue),
        backgroundColor: '#4F46E5'
      },
      {
        label: 'Meal Revenue',
        data: data.map(d => d.mealRevenue),
        backgroundColor: '#10B981'
      }
    ]
  }
})

const options = {
  responsive: true,
  plugins: { title: { display: true, text: 'Monthly Revenue Breakdown' } },
  scales: { y: { beginAtZero: true } }
}
</script>

<template>
  <Bar v-if="chartData" :data="chartData" :options="options" />
  <p v-else>Loading revenue data…</p>
</template>
```

### 12.2 Occupancy Heatmap

Build a calendar-style grid where each cell is a day, coloured by occupancy percentage (green → low, red → high). Use a simple CSS grid with computed background colours.

### 12.3 Meal Demand Chart

A stacked bar chart (breakfast/lunch/dinner) per day for a selected date range, using the same Chart.js approach as the revenue chart.

---

## 13. API Contract (Expected Backend)

The frontend expects a RESTful JSON API. Below is the minimum set of endpoints.

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Create account (individual or corporate) |
| POST | `/api/auth/login` | Returns JWT + user object |
| GET  | `/api/auth/me` | Current user profile |
| PUT  | `/api/auth/profile` | Update profile |
| POST | `/api/auth/reset-password` | Send reset link |

### Rooms

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET  | `/api/rooms` | All rooms (admin) |
| GET  | `/api/rooms/available?checkIn=&checkOut=` | Available rooms for date range |
| POST | `/api/rooms` | Create room (admin) |
| PUT  | `/api/rooms/:id` | Update room (admin) |
| DELETE | `/api/rooms/:id` | Delete room (admin) |

### Bookings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET  | `/api/bookings` | All bookings with filters (admin) |
| GET  | `/api/bookings/mine` | Current client's bookings |
| GET  | `/api/bookings/:id` | Single booking detail |
| POST | `/api/bookings` | Create booking |
| PATCH | `/api/bookings/:id/status` | Update status (admin) |
| PUT  | `/api/bookings/:id` | Modify booking (admin) |

### Meals

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET  | `/api/meals/menu?date=` | Menu items for a date |
| POST | `/api/meals/menu` | Add/update menu item (admin) |
| DELETE | `/api/meals/menu/:id` | Remove menu item (admin) |
| GET  | `/api/meals/requests?date=` | All meal requests for a date (admin) |
| POST | `/api/meals/requests` | Submit meal requests (part of booking) |

### Invoices

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET  | `/api/invoices` | All invoices (admin) |
| GET  | `/api/invoices/mine` | Client's invoices |
| GET  | `/api/invoices/:id` | Single invoice |
| PATCH | `/api/invoices/:id/status` | Mark paid/unpaid (admin) |
| POST | `/api/invoices/:bookingId/regenerate` | Regenerate invoice (admin) |

### Reports (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET  | `/api/reports/daily-bookings?date=` | Bookings for a day |
| GET  | `/api/reports/monthly-revenue?month=` | Revenue breakdown |
| GET  | `/api/reports/occupancy?from=&to=` | Occupancy rates |
| GET  | `/api/reports/meal-demand?from=&to=` | Meal counts |

---

## 14. Deployment Checklist

Before going to production, make sure the following items are addressed:

1. **Environment variables** — set `VITE_API_BASE_URL` for the production API.
2. **Build** — run `npm run build` and serve the `dist/` folder from a static host or CDN.
3. **CORS** — ensure the backend allows requests from the frontend's domain.
4. **HTTPS** — all traffic must be encrypted; JWT tokens travel in headers.
5. **Error boundaries** — add a global Vue error handler and friendly error pages (404, 500).
6. **Loading states** — every async operation should show a spinner or skeleton screen.
7. **Form validation** — all user inputs validated on the client (VeeValidate + Yup) AND on the server.
8. **Accessibility** — semantic HTML, ARIA labels, keyboard navigation on modals and forms.
9. **Responsive design** — test on mobile, tablet, and desktop breakpoints.
10. **Testing** — unit tests for stores and composables (Vitest), component tests (Vue Test Utils), and end-to-end tests (Cypress or Playwright).

---

## Summary of Page ↔ Process Mapping

| Page | Primary Process | Key Store(s) |
|------|----------------|-------------|
| Home | Marketing / entry point | — |
| Login | Authentication | auth |
| Register | Account creation (individual / corporate) | auth |
| Client Dashboard | Overview of active bookings | bookings |
| Booking Page | 4-step booking wizard | bookings, rooms, meals |
| Meal Selection | Add/edit meals for existing booking | meals |
| Booking History | View past and upcoming bookings | bookings |
| Invoice Page | View / download invoices | invoices |
| Profile Page | Update personal or company details | auth |
| Admin Dashboard | Operational summary | bookings, rooms |
| Room Management | CRUD rooms, toggle availability | rooms |
| Booking Management | Confirm, modify, cancel bookings | bookings |
| Meal Management | Update daily menu, view requests | meals |
| Invoice Management | View all invoices, mark paid | invoices |
| Reports | Revenue, occupancy, meal demand charts | reports |

This guide provides the full blueprint for building the Lodge Administration frontend in Vue.js. Each section can be tackled independently by a developer or team, and the process flows ensure that every user story from the original requirements document is accounted for.
