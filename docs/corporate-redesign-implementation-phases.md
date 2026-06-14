# Corporate Redesign — Implementation Phases

Reference design: `corporate-profile-redesign.md`

---

## Phase 1 — Database Foundation
*New tables and schema changes. No application code yet.*

### New migrations (in order)

1. **Create `cor_company_details`**
   - id, org_id, company_name, tpin, reg_number, industry, country, status, meta_data, created_at, updated_at
   - UNIQUE(org_id, reg_number)

2. **Create `cor_branch_details`**
   - id, company_id → cor_company_details (CASCADE), name, address, phone, meta_data, created_at, updated_at

3. **Create `cor_profiles`** *(replaces `corporate_profiles`)*
   - id, org_id, company_id → cor_company_details, branch_id → cor_branch_details (nullable)
   - first_name, last_name, email, phone, job_title, department, status, meta_data
   - UNIQUE(org_id, email)

4. **Create `corporate_guests`**
   - id, corporate_profile_id → cor_profiles (CASCADE), first_name, last_name, phone, email, identification_card
   - created_at, updated_at

5. **Rename `individual_profiles` → `website_users`**
   - Add identification_card column if not present

6. **Create `venues`**
   ```sql
   CREATE TABLE venues (
       id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       org_id       UUID NOT NULL REFERENCES organizations(id),
       branch_id    UUID REFERENCES branches(id),
       name         VARCHAR(255) NOT NULL,
       venue_type   VARCHAR(50) NOT NULL CHECK (venue_type IN ('conference_hall', 'event_space', 'boardroom', 'outdoor', 'dining')),
       capacity     INT NOT NULL CHECK (capacity > 0),
       area_sqm     NUMERIC(8, 2),
       floor        VARCHAR(50),
       base_rate    NUMERIC(10, 2) NOT NULL DEFAULT 0,
       rate_type    VARCHAR(10) NOT NULL DEFAULT 'daily' CHECK (rate_type IN ('hourly', 'daily')),
       amenities    TEXT[] NOT NULL DEFAULT '{}',
       is_available BOOLEAN NOT NULL DEFAULT TRUE,
       notes        TEXT,
       created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
   );
   ```

7. **Create `corporate_booking_requests`**
   - id, org_id, branch_id, website_user_id, cor_profile_id, company_id, booking_type, status
   - reason_for_booking, cost_center, notes
   - authoriser_name, authoriser_email, authoriser_phone, authoriser_title, authoriser_department, authoriser_gl_code
   - documents TEXT[], payload JSONB
   - created_at, updated_at

8. **Alter `bookings`** — add new columns
   - booking_type VARCHAR(20) DEFAULT 'room'
   - booker_type VARCHAR(20)
   - booker_name VARCHAR(255)
   - booker_email VARCHAR(255)
   - booker_phone VARCHAR(50)
   - website_user_id UUID FK → website_users (nullable)
   - cor_profile_id UUID FK → cor_profiles (nullable)
   - corporate_request_id UUID FK → corporate_booking_requests (nullable)
   - guest_name VARCHAR(255)
   - guest_identification VARCHAR(100)
   - corporate_guest_id UUID FK → corporate_guests (nullable)
   - venue_id UUID FK → venues (nullable)

9. **Alter `orders`** — add `corporate` to order_type enum

10. **Data migration** — backfill existing bookings
    - Set booking_type = 'room' for all existing rows
    - Set booker_type from existing client_type column
    - Copy client_name → booker_name, etc.

### Deliverable
All tables exist. Existing data intact. No breaking changes to running application yet.

---

## Phase 2 — Corporate Profile Layer
*Repository, service, and handler for the new corporate profile tables.*

### Models
- `CorCompanyDetails`, `CorBranchDetails`, `CorProfile`, `CorporateGuest`
- Request/response structs for create and update

### Repository
- `CorCompanyRepository` — CRUD + list by org_id
- `CorBranchRepository` — CRUD + list by company_id
- `CorProfileRepository` — CRUD + list by company_id, upsert by email
- `CorporateGuestRepository` — CRUD + list by corporate_profile_id

### Service
- `CorProfileService` — wraps all four repos, enforces one-rep-per-department rule

### Handlers + Routes
```
POST   /api/v1/clients/companies                          → create company
GET    /api/v1/clients/companies                          → list companies
GET    /api/v1/clients/companies/{id}                     → get company + branches + profiles
PUT    /api/v1/clients/companies/{id}                     → update company

POST   /api/v1/clients/companies/{id}/branches            → create branch
GET    /api/v1/clients/companies/{id}/branches            → list branches
PUT    /api/v1/clients/companies/{id}/branches/{branch_id}→ update branch

POST   /api/v1/clients/companies/{id}/profiles            → create profile (rep)
GET    /api/v1/clients/companies/{id}/profiles            → list profiles
GET    /api/v1/clients/profiles/{id}                      → get profile + guests
PUT    /api/v1/clients/profiles/{id}                      → update profile

POST   /api/v1/clients/profiles/{id}/guests               → add guest
GET    /api/v1/clients/profiles/{id}/guests               → list guests
PUT    /api/v1/clients/profiles/{id}/guests/{guest_id}    → update guest
DELETE /api/v1/clients/profiles/{id}/guests/{guest_id}    → remove guest
```

### Deliverable
Backoffice staff can fully manage corporate companies, branches, profiles, and guests. Old `corporate_profiles` endpoints still live in parallel — not removed yet.

---

## Phase 3 — Website Users
*Rename and extend individual_profiles into website_users.*

### Changes
- Rename `IndividualProfile` model → `WebsiteUser`
- Update all repo, service, handler references
- Add `identification_card` field to model and scanner
- Update guest auth flow to use `website_users` table name

### Routes (unchanged, just updated internals)
```
POST  /api/v1/guest/auth/register
POST  /api/v1/guest/auth/login
GET   /api/v1/guest/profile
PUT   /api/v1/guest/profile
```

### Deliverable
`individual_profiles` fully replaced by `website_users` in code. Guest auth works as before.

---

## Phase 4 — Venues
*Full CRUD for venues. No booking integration yet.*

### Models
- `Venue` struct with all fields including `amenities []string`
- `CreateVenueRequest`, `UpdateVenueRequest`

### Repository
- `VenueRepository` — CRUD + list by org_id + filter by venue_type + availability check

### Service
- `VenueService` — wraps repo, checks capacity conflicts on date ranges

### Handlers + Routes
```
POST  /api/v1/venues                     → create venue (admin/branch_admin)
GET   /api/v1/venues                     → list venues (all staff)
GET   /api/v1/venues/{id}               → get venue
PUT   /api/v1/venues/{id}               → update venue
DELETE /api/v1/venues/{id}              → delete venue

GET   /api/v1/guest/venues              → public list (filter by org_id, venue_type)
GET   /api/v1/guest/venues/{id}         → public get
```

### Deliverable
Venues are manageable from backoffice and browsable by guests.

---

## Phase 5 — Corporate Booking Requests
*The staging layer for all corporate booking types submitted via the website.*

### Models
- `CorporateBookingRequest` struct
- Per-type payload structs: `AccommodationPayload`, `MealsPayload`, `ConferencePayload`, `EventPayload`

### Repository
- `CorporateBookingRequestRepository` — create, get by ID, list by org_id + status, update status

### Service
- `CorporateBookingRequestService`
  - `Submit(websiteUserID, req)` — validates, stores request, fires workflow task
  - `Approve(requestID, staffUserID)` — transitions status, triggers space assignment
  - `Reject(requestID, reason)` — transitions status
  - `List(orgID, filters)` — for backoffice view

### Guest Routes
```
POST  /api/v1/guest/bookings/corporate?type=accommodation
POST  /api/v1/guest/bookings/corporate?type=meals
POST  /api/v1/guest/bookings/corporate?type=conference
POST  /api/v1/guest/bookings/corporate?type=event
GET   /api/v1/guest/bookings/corporate          → list own requests
GET   /api/v1/guest/bookings/corporate/{id}     → get own request
```

### Backoffice Routes
```
GET   /api/v1/bookings/requests                 → list all requests (filterable by type/status)
GET   /api/v1/bookings/requests/{id}            → get request + payload
PUT   /api/v1/bookings/requests/{id}/approve    → approve
PUT   /api/v1/bookings/requests/{id}/reject     → reject
```

### Deliverable
Corporate reps can submit all four booking types. Backoffice staff can see and act on requests through workflow.

---

## Phase 6 — Room & Venue Assignment (Post-Approval)
*Staff assigns spaces to guests after approving a corporate request. This is where bookings rows are created.*

### Accommodation assignment
- Staff selects a guest from `corporate_booking_requests.payload`
- Staff picks a room and dates
- System creates a `bookings` row: `booking_type = 'room'`, `room_id` set, `corporate_guest_id` set, `corporate_request_id` set
- Repeat per guest until all guests have rooms
- When all guests assigned → `corporate_booking_requests.status = 'approved'`

### Meals assignment
- Staff approves the meals request
- System creates one `bookings` row per guest: `booking_type = 'meals'`, `venue_id` → dining
- System creates one `orders` row per booking: `type = 'corporate'`
- System creates `order_items` from payload `meal_items` per guest (price snapshotted)
- Kitchen sees orders immediately

### Conference / Event assignment
- Staff picks a venue and confirms dates
- System creates one `bookings` row for the group: `booking_type = 'conference'|'event'`, `venue_id` set

### Routes
```
POST  /api/v1/bookings/requests/{id}/assign-room      → assign room to one guest
POST  /api/v1/bookings/requests/{id}/assign-venue     → assign venue (conference/event/meals)
GET   /api/v1/bookings/requests/{id}/assignments      → list current assignments
```

### Deliverable
Full end-to-end corporate booking flow working. Bookings table is the operational record for all confirmed assignments.

---

## Phase 7 — Cleanup & Deprecation
*Remove old tables and code once new system is stable.*

### Remove
- `corporate_profiles` table (after data migrated to `cor_profiles`)
- Old `corporate_profiles` endpoints and handlers
- `client_type` / `corporate_client_id` columns from `bookings` (replaced by new columns)
- `CorporateBookingResponse`, `CreateCorporateBookingRequest` structs from `booking.go`
- Old `GuestBookingService.CreateCorporate`, `CreateCorporateMeals`, `CreateCorporateConference` methods
- `booking_documents` table (documents now on `corporate_booking_requests.documents`)
- `conference` value from `room_type` enum — migrate any existing conference rooms to `venues`

### Deliverable
Codebase clean. No legacy corporate booking code remaining.

---

## Summary Table

| Phase | What | Touches DB | Touches API | Breaking |
|---|---|---|---|---|
| 1 | Schema — new tables + alter bookings | Yes | No | No |
| 2 | Corporate profile layer (companies, branches, reps, guests) | No | Yes | No |
| 3 | Website users rename | No | No | No |
| 4 | Venues CRUD | No | Yes | No |
| 5 | Corporate booking requests | No | Yes | No |
| 6 | Room/venue assignment post-approval | No | Yes | No |
| 7 | Remove legacy code | Yes | Yes | Yes — do last |
