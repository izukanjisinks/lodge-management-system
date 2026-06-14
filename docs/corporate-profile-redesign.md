# Corporate Profile Redesign

## Overview

The original `corporate_profiles` table conflated company-level data, branch-level data, and individual contact data into a single row, causing duplication when multiple individuals from the same company (or different branches/departments) were added.

This redesign splits the concerns into four tables:

- `cor_company_details` — the company record
- `cor_branch_details` — offices/branches of a company
- `cor_profiles` — the individual representative from a department
- `corporate_guests` — employees/guests the representative is booking on behalf of

---

## Table Definitions

### `cor_company_details`
Holds company-level information. One row per company per organization.

| Column | Type | Notes |
|---|---|---|
| id | UUID PK | |
| org_id | UUID FK → organizations | |
| company_name | VARCHAR(255) NOT NULL | |
| tpin | VARCHAR(100) | Tax Payer Identification Number |
| reg_number | VARCHAR(100) | Company registration number |
| industry | VARCHAR(100) | |
| country | VARCHAR(100) | |
| status | VARCHAR(20) | `active` \| `inactive`, default `active` |
| meta_data | JSONB | Free-form extra fields |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |
| | UNIQUE(org_id, reg_number) | |

---

### `cor_branch_details`
An office or branch of a company. One company can have many branches.

| Column | Type | Notes |
|---|---|---|
| id | UUID PK | |
| company_id | UUID FK → cor_company_details | CASCADE DELETE |
| name | VARCHAR(255) NOT NULL | e.g. "Lusaka Office", "Head Office" |
| address | TEXT | |
| phone | VARCHAR(50) | |
| meta_data | JSONB | Free-form extra fields |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

---

### `cor_profiles`
The individual representative from a company department. Replaces `corporate_profiles`.
One representative per department is enforced at the application level.

| Column | Type | Notes |
|---|---|---|
| id | UUID PK | |
| org_id | UUID FK → organizations | |
| company_id | UUID FK → cor_company_details | |
| branch_id | UUID FK → cor_branch_details | Nullable — not all reps are branch-specific |
| first_name | VARCHAR(100) NOT NULL | |
| last_name | VARCHAR(100) NOT NULL | |
| email | VARCHAR(255) | |
| phone | VARCHAR(50) | |
| job_title | VARCHAR(100) | |
| department | VARCHAR(100) | One rep per department (app-level rule) |
| status | VARCHAR(20) | `active` \| `inactive`, default `active` |
| meta_data | JSONB | Free-form extra fields |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |
| | UNIQUE(org_id, email) | |

---

### `corporate_guests`
Employees or guests that a representative is making bookings for.
Guests belong to the same department as their representative.

| Column | Type | Notes |
|---|---|---|
| id | UUID PK | |
| corporate_profile_id | UUID FK → cor_profiles | CASCADE DELETE |
| first_name | VARCHAR(100) NOT NULL | |
| last_name | VARCHAR(100) NOT NULL | |
| phone | VARCHAR(50) | |
| email | VARCHAR(255) | Optional |
| identification_card | VARCHAR(100) NOT NULL | NRC or Passport number |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

---

## Relationship Map

```
cor_company_details (1)
        │
        ├──── (many) cor_branch_details
        │
        └──── (many) cor_profiles
                          │
                          │   [one rep per department enforced at app level]
                          │
                          └──── (many) corporate_guests
```

---

## Example

```
cor_company_details: { company_name: "Zambeef", reg_number: "ZM-001", tpin: "1234567", industry: "Agriculture" }
        │
        ├── cor_branch_details: { name: "Lusaka Office",  address: "Cairo Rd", phone: "+260 211 000" }
        ├── cor_branch_details: { name: "Ndola Office",   address: "Broadway",  phone: "+260 212 000" }
        │
        ├── cor_profiles: { first_name: "John",  last_name: "Banda", department: "Finance",  branch_id: Lusaka  }
        │       └── corporate_guests: { first_name: "Alice", identification_card: "123456/78/1" }
        │       └── corporate_guests: { first_name: "Bob",   identification_card: "234567/89/2" }
        │
        └── cor_profiles: { first_name: "Mary",  last_name: "Phiri", department: "HR",       branch_id: Ndola   }
                └── corporate_guests: { first_name: "Carol", identification_card: "345678/90/3" }
```

---

## Website Users & Booking Flow

### `website_users` *(renamed from `individual_profiles`)*
The primary table for anyone interacting with the website — individual bookers and corporate representatives alike. Not scoped to any lodge.

| Column | Type | Notes |
|---|---|---|
| id | UUID PK | |
| first_name | VARCHAR(100) NOT NULL | |
| last_name | VARCHAR(100) NOT NULL | |
| email | VARCHAR(255) | UNIQUE |
| phone | VARCHAR(50) | |
| identification_card | VARCHAR(100) | NRC or Passport |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

---

### `corporate_booking_requests`
Staging table for corporate bookings submitted via the website. A request sits here while workflow is in progress and staff are assigning rooms. One row per corporate booking submission.

| Column | Type | Notes |
|---|---|---|
| id | UUID PK | |
| org_id | UUID FK → organizations | Which lodge this is for |
| branch_id | UUID FK → branches | Nullable |
| website_user_id | UUID FK → website_users | The rep who submitted |
| cor_profile_id | UUID FK → cor_profiles | Nullable — set if rep already exists |
| company_id | UUID FK → cor_company_details | Nullable — set if company already exists |
| booking_type | VARCHAR(20) | `accommodation` \| `meals` \| `conference` \| `event` |
| status | VARCHAR(20) | `pending` \| `approved` \| `rejected` \| `cancelled` |
| reason_for_booking | TEXT | |
| cost_center | VARCHAR(100) | |
| notes | TEXT | |
| authoriser_name | VARCHAR(255) | |
| authoriser_email | VARCHAR(255) | |
| authoriser_phone | VARCHAR(50) | |
| authoriser_title | VARCHAR(100) | |
| authoriser_department | VARCHAR(100) | |
| authoriser_gl_code | VARCHAR(50) | |
| documents | TEXT[] | Uploaded supporting documents |
| payload | JSONB | Full request snapshot — guest list, dates, preferences, etc. |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

**Example `payload` for an accommodation request:**
```json
{
  "guests": [
    { "first_name": "Alice", "last_name": "Mwale", "identification_card": "123456/78/1", "check_in": "2026-07-01", "check_out": "2026-07-05", "room_type": "double" },
    { "first_name": "Bob",   "last_name": "Phiri", "identification_card": "234567/89/2", "check_in": "2026-07-01", "check_out": "2026-07-03", "room_type": "single" }
  ]
}
```

---

### `venues`
Bookable spaces that are not rooms — conference halls, event spaces, boardrooms, outdoor areas, and dining halls. Amenities follow the same `TEXT[]` pattern as rooms.

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

**Example venue response:**
```json
{
  "id": "uuid",
  "name": "The Riverside Restaurant",
  "venue_type": "dining",
  "capacity": 80,
  "area_sqm": 120.0,
  "floor": "Ground Floor",
  "base_rate": 0.00,
  "rate_type": "daily",
  "is_available": true,
  "amenities": ["Air Conditioning", "WiFi", "Outdoor Seating", "Bar"]
}
```

**Venue types and their booking use:**

| venue_type | booking_type | Notes |
|---|---|---|
| `conference_hall` | `conference` | Formal meetings, large groups |
| `boardroom` | `conference` | Smaller meeting rooms |
| `event_space` | `event` | Weddings, galas, functions |
| `outdoor` | `event` | Garden, poolside, open-air |
| `dining` | `meals` | Restaurant/dining hall for meal bookings |

**Note:** The `conference` room_type on the existing `rooms` enum will be removed and those records migrated to `venues` during implementation.

---

### `bookings` *(reworked)*
Every confirmed booking — individual or corporate — ends up here. One row per assignment.

- Accommodation bookings → `room_id` set, `venue_id` null
- Conference/event bookings → `venue_id` set, `room_id` null
- Meals bookings → `venue_id` points to a `dining` venue, `room_id` null

| Column | Type | Notes |
|---|---|---|
| id | UUID PK | |
| booking_number | VARCHAR | Auto-generated reference |
| org_id | UUID FK → organizations | |
| branch_id | UUID FK → branches | Nullable |
| booking_type | VARCHAR(20) | `room` \| `meals` \| `conference` \| `event` |
| **Booker info (denormalised)** | | |
| booker_type | VARCHAR(20) | `individual` \| `corporate` |
| booker_name | VARCHAR(255) | Full name of whoever made the booking |
| booker_email | VARCHAR(255) | |
| booker_phone | VARCHAR(50) | |
| website_user_id | UUID FK → website_users | Nullable — set for individual bookers |
| cor_profile_id | UUID FK → cor_profiles | Nullable — set for corporate reps |
| corporate_request_id | UUID FK → corporate_booking_requests | Nullable — the originating request |
| **Guest info** | | |
| guest_name | VARCHAR(255) | For corporate: the employee being accommodated |
| guest_identification | VARCHAR(100) | NRC/Passport of the actual guest |
| corporate_guest_id | UUID FK → corporate_guests | Nullable — links to the guest record |
| **Space assignment** | | |
| room_id | UUID FK → rooms | Set for `room` booking type only |
| venue_id | UUID FK → venues | Set for `meals`, `conference`, `event` types |
| check_in | DATE | |
| check_out | DATE | |
| guests | INT | Number of people in this booking slot |
| **Financials** | | |
| room_cost | DECIMAL | |
| total_amount | DECIMAL | |
| status | VARCHAR(20) | `pending` \| `confirmed` \| `checked_in` \| `checked_out` \| `cancelled` |
| special_requests | TEXT | |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

**Why booker info is denormalised:** Individual bookers (`website_users`) are public-facing and not lodge records. Storing their contact info directly on the booking ensures invoices and records remain correct even if the user later changes their details or deletes their account.

---

## Meals Booking — Orders Integration

Corporate meals bookings reuse the existing `orders` + `order_items` system rather than introducing new tables. The kitchen sees the same order view regardless of whether the order came from an in-house guest or a corporate booking request.

**`orders.type` gains a third value:** `corporate` — so kitchen/restaurant staff can filter by origin.

### Flow
```
corporate_booking_requests (payload has guest list + meal_items per guest)
        │
        └──► workflow → staff approves
                        │
                        └──► bookings (1 row per guest, booking_type = 'meals', venue_id → dining)
                                        │
                                        └──► orders (1 order per guest, booking_id = booking.id, type = 'corporate')
                                                        │
                                                        └──► order_items (1 row per menu item)
                                                                  menu_item_id, quantity, unit_price (snapshotted)
```

### What each layer holds

| Layer | Table | Key fields |
|---|---|---|
| Request snapshot | `corporate_booking_requests.payload` | Full guest + meal item list as submitted |
| Per-guest booking | `bookings` | `booking_type = 'meals'`, `guest_name`, `venue_id` |
| Kitchen order | `orders` | `booking_id`, `type = 'corporate'`, `status` |
| Line items | `order_items` | `menu_item_id`, `quantity`, `unit_price` (snapshotted) |

The kitchen filters `orders` by `org_id + status = 'open'` — same as today. The `type = 'corporate'` flag lets them optionally filter to see only corporate meal orders.

---

## Full Relationship Map

```
website_users
    │
    ├──► bookings (individual room booking)
    │         booking_type = 'room', room_id set, booker_type = 'individual'
    │         └──► orders (in_house) → order_items
    │
    ├──► bookings (individual conference/event)
    │         booking_type = 'conference'|'event', venue_id set
    │
    └──► corporate_booking_requests (staging for all corporate types)
                    │
                    └──► workflow → staff processes
                                        │
                                        ├── accommodation → bookings (1 per guest, room_id set)
                                        │
                                        ├── meals → bookings (1 per guest, venue_id → dining)
                                        │               └──► orders (type='corporate') → order_items
                                        │
                                        └── conference/event → bookings (1 for group, venue_id set)

cor_company_details (1)
        │
        ├──── (many) cor_branch_details
        │
        └──── (many) cor_profiles
                          │
                          ├──── (many) corporate_guests
                          └──► corporate_booking_requests (cor_profile_id)

rooms   ──► bookings.room_id    (accommodation)
venues  ──► bookings.venue_id   (meals → dining, conference, event)
```

---

## Booking Flow Summary

### Individual Room Booking
```
website_users submits
        └──► bookings (1 row, booking_type = 'room', room_id set, booker_type = 'individual')
                        └──► workflow → staff confirms → status: confirmed
```

### Corporate Accommodation Booking
```
website_users (rep) submits
        └──► corporate_booking_requests (booking_type = 'accommodation', guests in payload)
                        └──► workflow task created
                                        └──► staff assigns rooms per guest
                                                        └──► bookings (1 row per guest)
                                                                  booking_type = 'room'
                                                                  room_id set per assignment
                                                                  corporate_guest_id = guest
```

### Corporate Meals Booking
```
website_users (rep) submits
        └──► corporate_booking_requests (booking_type = 'meals', guests + meal_items in payload)
                        └──► workflow task created
                                        └──► staff approves
                                                        └──► bookings (1 row per guest)
                                                                  booking_type = 'meals'
                                                                  venue_id → dining venue
                                                                        └──► orders (type = 'corporate')
                                                                                └──► order_items (snapshotted)
```

### Corporate Conference / Event Booking
```
website_users (rep) submits
        └──► corporate_booking_requests (booking_type = 'conference'|'event', details in payload)
                        └──► workflow task created
                                        └──► staff assigns venue
                                                        └──► bookings (1 row for the group)
                                                                  booking_type = 'conference'|'event'
                                                                  venue_id → conference_hall/event_space
```

---

## Migration Notes

- `corporate_profiles` is replaced by `cor_profiles`
- `individual_profiles` is replaced by `website_users`
- Existing references to `corporate_profiles.id` on `bookings` and `booking_documents` will need updating to `cor_profiles.id`
- `cor_company_details`, `cor_branch_details`, and `corporate_booking_requests` are new tables
- Data migration: existing `corporate_profiles` rows should be split — company fields into `cor_company_details`, individual fields into `cor_profiles`
- Existing `bookings` rows will need `booker_type`, `booker_name`, `booker_email`, `booker_phone`, `booking_type` backfilled
