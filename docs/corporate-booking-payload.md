# Corporate Booking — Complete API Payload Contract

## Design Rules
- Every top-level section is **always present**. When a service (accommodation, events, meals) is not enabled, it is sent as **`null`**, not omitted.
- Every field within an enabled section is **always present**. Empty strings become `null`; unset numbers become `null`.
- `attendants` is always an array of at least one record (the booker/rep).
- `sessions` (events and meals) are **flattened** — one object per session per calendar date.
- Company fields are **snapshots** — copied from the profile at booking time and editable by the rep before submission. The `selected_*_id` fields are reference pointers only; the snapshot fields are the source of truth for billing and reporting.
- Corporate accommodation does **not** include individual room assignments. The property team assigns specific rooms based on availability prior to the event. Only room type preference and count are submitted.

---

## Full Payload — All Services Enabled (Detailed Mode)

```json
{
  "org_id": "5",
  "branch_id": "12",
  "booking_type": "corporate",
  "source": "web",
  "currency": "ZMW",

  "participant_mode": "detailed",
  "participant_count": null,

  "company": {
    "selected_company_id": "88",
    "selected_branch_id": "14",
    "corporate_profile_id": "31",
    "name": "Acme Corporation Ltd",
    "tpin": "1003456789",
    "industry": "Mining",
    "email": "billing@acme.co.zm",
    "phone": "+260971000001",
    "city": "Lusaka",
    "street_address": "Plot 1234, Cairo Road",
    "branch_name": "Lusaka Head Office",
    "department_name": "Human Resources",
    "cost_center": "CC-HR-001",
    "gl_code": "GL-4500"
  },

  "booked_by": {
    "name": "Martin Sinkolongo",
    "email": "martin@acme.co.zm",
    "phone": "+260971234567",
    "job_title": "HR Manager"
  },

  "approver": {
    "name": "Sandra Mwale",
    "email": "sandra.mwale@acme.co.zm",
    "phone": "+260971000099",
    "title": "Finance Director"
  },

  "attendants": [
    {
      "full_name": "Martin Sinkolongo",
      "email": "martin@acme.co.zm",
      "phone": "+260971234567",
      "id_number": null,
      "dietary_notes": null,
      "company": "Acme Corporation Ltd",
      "is_lead_contact": true
    },
    {
      "full_name": "Jane Banda",
      "email": "jane.banda@acme.co.zm",
      "phone": null,
      "id_number": "234567/89/1",
      "dietary_notes": "Vegetarian, no nuts",
      "company": "Acme Corporation Ltd",
      "is_lead_contact": false
    },
    {
      "full_name": "David Phiri",
      "email": "d.phiri@partner.co.zm",
      "phone": "+260972000002",
      "id_number": null,
      "dietary_notes": null,
      "company": "Partner Consulting Ltd",
      "is_lead_contact": false
    }
  ],

  "accommodation": {
    "reason_for_booking": "Overnight stay for out-of-town delegates",
    "room_type": "double",
    "room_count": 5,
    "check_in": "2026-08-10",
    "check_out": "2026-08-12",
    "notes": "Prefer rooms on the same floor. Early check-in from 12:00 if possible."
  },

  "events": {
    "reason_for_booking": "Annual HR Strategy Workshop",
    "start_date": "2026-08-11",
    "end_date": "2026-08-12",
    "schedule_mode": "uniform",
    "sessions": [
      {
        "event_name": "Morning Plenary",
        "event_type": "workshop",
        "event_date": "2026-08-11",
        "start_time": "08:00",
        "end_time": "12:00",
        "expected_attendees": 40,
        "setup_type": "theatre",
        "venue_id": "9",
        "venue_name": "The Baobab Hall",
        "venue_capacity": 200,
        "pricing_basis": "half_day",
        "special_requirements": "Projector, PA system, company banner stand"
      },
      {
        "event_name": "Afternoon Breakout",
        "event_type": "workshop",
        "event_date": "2026-08-11",
        "start_time": "13:00",
        "end_time": "17:00",
        "expected_attendees": 20,
        "setup_type": "boardroom",
        "venue_id": "11",
        "venue_name": "Boardroom A",
        "venue_capacity": 30,
        "pricing_basis": "half_day",
        "special_requirements": null
      },
      {
        "event_name": "Morning Plenary",
        "event_type": "workshop",
        "event_date": "2026-08-12",
        "start_time": "08:00",
        "end_time": "12:00",
        "expected_attendees": 40,
        "setup_type": "theatre",
        "venue_id": "9",
        "venue_name": "The Baobab Hall",
        "venue_capacity": 200,
        "pricing_basis": "half_day",
        "special_requirements": "Projector, PA system, company banner stand"
      },
      {
        "event_name": "Afternoon Breakout",
        "event_type": "workshop",
        "event_date": "2026-08-12",
        "start_time": "13:00",
        "end_time": "17:00",
        "expected_attendees": 20,
        "setup_type": "boardroom",
        "venue_id": "11",
        "venue_name": "Boardroom A",
        "venue_capacity": 30,
        "pricing_basis": "half_day",
        "special_requirements": null
      }
    ]
  },

  "meals": {
    "reason_for_booking": null,
    "meal_mode": "event_linked",
    "schedule_mode": "uniform",
    "sessions": [
      {
        "session_name": null,
        "meal_date": "2026-08-11",
        "meal_period": "breakfast",
        "service_type": "buffet",
        "pax_count": 40,
        "linked_master_session_index": 0,
        "dietary_notes": "Two vegetarian guests, one gluten-free",
        "arrangements_notes": "Set up in the dining area from 07:00",
        "individual_orders": []
      },
      {
        "session_name": null,
        "meal_date": "2026-08-11",
        "meal_period": "lunch",
        "service_type": "buffet",
        "pax_count": 40,
        "linked_master_session_index": 0,
        "dietary_notes": null,
        "arrangements_notes": null,
        "individual_orders": []
      },
      {
        "session_name": "VIP Dinner",
        "meal_date": "2026-08-11",
        "meal_period": "dinner",
        "service_type": "individual_order",
        "pax_count": 3,
        "linked_master_session_index": null,
        "dietary_notes": null,
        "arrangements_notes": "Reserved table on the terrace",
        "individual_orders": [
          {
            "attendant_idx": 0,
            "menu_item_id": "menu_305",
            "quantity": 1,
            "notes": null
          },
          {
            "attendant_idx": 1,
            "menu_item_id": "menu_410",
            "quantity": 1,
            "notes": "No garlic"
          },
          {
            "attendant_idx": 2,
            "menu_item_id": "menu_305",
            "quantity": 1,
            "notes": null
          }
        ]
      },
      {
        "session_name": null,
        "meal_date": "2026-08-12",
        "meal_period": "breakfast",
        "service_type": "buffet",
        "pax_count": 40,
        "linked_master_session_index": 0,
        "dietary_notes": "Two vegetarian guests, one gluten-free",
        "arrangements_notes": "Set up in the dining area from 07:00",
        "individual_orders": []
      },
      {
        "session_name": null,
        "meal_date": "2026-08-12",
        "meal_period": "lunch",
        "service_type": "buffet",
        "pax_count": 40,
        "linked_master_session_index": 0,
        "dietary_notes": null,
        "arrangements_notes": null,
        "individual_orders": []
      }
    ]
  },

  "notes": "Please ensure all signage is removed after the event. Invoice to be addressed to the Finance Director."
}
```

---

## Minimal Payload — Headcount Mode, No Services

```json
{
  "org_id": "5",
  "branch_id": "12",
  "booking_type": "corporate",
  "source": "web",
  "currency": "ZMW",

  "participant_mode": "headcount",
  "participant_count": 25,

  "company": {
    "selected_company_id": "88",
    "selected_branch_id": null,
    "corporate_profile_id": null,
    "name": "Acme Corporation Ltd",
    "tpin": null,
    "industry": null,
    "email": null,
    "phone": null,
    "city": null,
    "street_address": null,
    "branch_name": null,
    "department_name": null,
    "cost_center": null,
    "gl_code": null
  },

  "booked_by": {
    "name": "Martin Sinkolongo",
    "email": "martin@acme.co.zm",
    "phone": null,
    "job_title": null
  },

  "approver": {
    "name": null,
    "email": null,
    "phone": null,
    "title": null
  },

  "attendants": [
    {
      "full_name": "Martin Sinkolongo",
      "email": "martin@acme.co.zm",
      "phone": null,
      "id_number": null,
      "dietary_notes": null,
      "company": null,
      "is_lead_contact": true
    }
  ],

  "accommodation": null,
  "events": null,
  "meals": null,

  "notes": null
}
```

---

## Field Reference

### Root Object

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `org_id` | `string` | ✅ | Lodge / property identifier |
| `branch_id` | `string \| null` | ✅ | Lodge branch identifier; null if lodge has no branches |
| `booking_type` | `"corporate"` | ✅ | Fixed value |
| `source` | `string` | ✅ | `"web"` \| `"mobile"` \| `"reception"` |
| `currency` | `string` | ✅ | ISO 4217 — e.g. `"ZMW"` |
| `participant_mode` | `string` | ✅ | `"headcount"` \| `"detailed"` |
| `participant_count` | `integer \| null` | ✅ | Total delegates incl. booker. Non-null **only** when `participant_mode = "headcount"`; null in detailed mode |
| `company` | `object` | ✅ | See below. Always present; inner fields may be null |
| `booked_by` | `object` | ✅ | The company rep submitting the booking. See below |
| `approver` | `object` | ✅ | The authorising approver. Always present; all fields null if not provided |
| `attendants` | `array` | ✅ | Min 1 record (the booker/rep). See below |
| `accommodation` | `object \| null` | ✅ | null when accommodation is not selected |
| `events` | `object \| null` | ✅ | null when events are not selected |
| `meals` | `object \| null` | ✅ | null when meals are not selected |
| `notes` | `string \| null` | ✅ | General booking notes |

---

### `company`

Snapshot fields are copied from the corporate profile at booking time. The rep may override them before submission. `selected_*_id` fields are reference pointers for CRM linkage only — the snapshot fields drive billing and reporting.

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `selected_company_id` | `string \| null` | ✅ | ID of the company record selected from the CRM |
| `selected_branch_id` | `string \| null` | ✅ | ID of the company branch selected; null if no branch chosen |
| `corporate_profile_id` | `string \| null` | ✅ | ID of the billing/department profile used; null if no profile chosen |
| `name` | `string \| null` | ✅ | Company trading name |
| `tpin` | `string \| null` | ✅ | Tax / registration number |
| `industry` | `string \| null` | ✅ | Industry sector |
| `email` | `string \| null` | ✅ | Company billing email |
| `phone` | `string \| null` | ✅ | Company / branch phone |
| `city` | `string \| null` | ✅ | Company city |
| `street_address` | `string \| null` | ✅ | Company street address |
| `branch_name` | `string \| null` | ✅ | Name of the company branch |
| `department_name` | `string \| null` | ✅ | Booking department |
| `cost_center` | `string \| null` | ✅ | Internal cost centre code |
| `gl_code` | `string \| null` | ✅ | General ledger code for invoice allocation |

---

### `booked_by`

The company representative who is submitting the booking. Auto-filled from the authenticated user account.

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `name` | `string` | ✅ | Full name — required |
| `email` | `string \| null` | ✅ | Work email |
| `phone` | `string \| null` | ✅ | |
| `job_title` | `string \| null` | ✅ | Rep's job title |

---

### `approver`

The person authorising this booking on behalf of the company. Always present as an object; all fields are null if the rep leaves the approver section blank.

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `name` | `string \| null` | ✅ | |
| `email` | `string \| null` | ✅ | |
| `phone` | `string \| null` | ✅ | |
| `title` | `string \| null` | ✅ | Job title / authority level |

---

### `attendants[]`

In **headcount mode**, the array contains exactly one record — the booker/rep — with `is_lead_contact: true` and all optional fields as `null`. The full delegate list is communicated via `participant_count`.

In **detailed mode**, the array contains one record per registered delegate, with the booker's record first.

The `company` field on each attendant supports multi-company events where delegates come from different organisations.

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `full_name` | `string` | ✅ | |
| `email` | `string \| null` | ✅ | |
| `phone` | `string \| null` | ✅ | |
| `id_number` | `string \| null` | ✅ | National ID / passport number |
| `dietary_notes` | `string \| null` | ✅ | |
| `company` | `string \| null` | ✅ | Attendant's own company name; useful for multi-company events |
| `is_lead_contact` | `boolean` | ✅ | Exactly one attendant carries `true` |

---

### `accommodation` (when not null)

Corporate accommodation does **not** include per-room or per-attendant assignments. The property team assigns specific rooms prior to the event based on availability. The booking specifies only the preferred room type and the number of rooms required.

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `reason_for_booking` | `string \| null` | ✅ | e.g. "Out-of-town delegates" |
| `room_type` | `string \| null` | ✅ | Preferred room type: `"single"` \| `"double"` \| `"twin"` \| `"suite"` \| `null` if no preference |
| `room_count` | `integer` | ✅ | Total number of rooms required. Min 1 |
| `check_in` | `string \| null` | ✅ | ISO 8601 date `YYYY-MM-DD` |
| `check_out` | `string \| null` | ✅ | ISO 8601 date `YYYY-MM-DD` |
| `notes` | `string \| null` | ✅ | Special requests — floor preferences, accessibility, early check-in, etc. |

---

### `events` (when not null)

Sessions are **flattened by date**. `schedule_mode = "uniform"` with 2 master sessions over 3 days produces 6 session objects (3 days × 2 sessions). Days excluded by the rep via per-day overrides are absent from the array entirely.

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `reason_for_booking` | `string \| null` | ✅ | Purpose of the event booking |
| `start_date` | `string \| null` | ✅ | ISO 8601 date — first event day; null if dates not set |
| `end_date` | `string \| null` | ✅ | ISO 8601 date — last event day; null if dates not set |
| `schedule_mode` | `string` | ✅ | `"uniform"` \| `"per_day"` |
| `sessions` | `array` | ✅ | Flattened session list; see below |

#### `events.sessions[]`

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `event_name` | `string \| null` | ✅ | Optional session label / name |
| `event_type` | `string` | ✅ | See **Event Types** enum |
| `event_date` | `string \| null` | ✅ | `YYYY-MM-DD`; null if no date range was configured |
| `start_time` | `string` | ✅ | `HH:MM` 24-hour |
| `end_time` | `string` | ✅ | `HH:MM` 24-hour |
| `expected_attendees` | `integer` | ✅ | Number of people for this session |
| `setup_type` | `string` | ✅ | See **Setup Types** enum |
| `venue_id` | `string \| null` | ✅ | null if no venue selected |
| `venue_name` | `string \| null` | ✅ | Snapshot of venue name at booking time |
| `venue_capacity` | `integer \| null` | ✅ | Snapshot of max capacity at booking time |
| `pricing_basis` | `string` | ✅ | See **Pricing Basis** enum |
| `special_requirements` | `string \| null` | ✅ | AV, signage, branding, staging requests |

---

### `meals` (when not null)

When `meal_mode = "event_linked"`, the meal date range is derived from `events.start_date` / `events.end_date`. When `meal_mode = "standalone"`, meals carry their own independent date range (below).

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `reason_for_booking` | `string \| null` | ✅ | |
| `meal_mode` | `string` | ✅ | `"event_linked"` \| `"standalone"` |
| `schedule_mode` | `string` | ✅ | `"uniform"` \| `"per_day"` |
| `start_date` | `string \| null` | ✅ | ISO 8601 date. Non-null **only** when `meal_mode = "standalone"` |
| `end_date` | `string \| null` | ✅ | ISO 8601 date. Non-null **only** when `meal_mode = "standalone"` |
| `sessions` | `array` | ✅ | Flattened; see below |

#### `meals.sessions[]`

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `session_name` | `string \| null` | ✅ | Optional label, e.g. "VIP Dinner" |
| `meal_date` | `string \| null` | ✅ | `YYYY-MM-DD`; null if no date range |
| `meal_period` | `string` | ✅ | See **Meal Periods** enum |
| `service_type` | `string` | ✅ | See **Service Types** enum |
| `pax_count` | `integer` | ✅ | Total covers for this meal |
| `linked_master_session_index` | `integer \| null` | ✅ | Index into the master sessions list this meal is tied to; null when `meal_mode = "standalone"` or meal is not linked to a specific session |
| `dietary_notes` | `string \| null` | ✅ | Group-level dietary requirements |
| `arrangements_notes` | `string \| null` | ✅ | Service logistics — timing, table layout, service notes |
| `individual_orders` | `array` | ✅ | Always `[]` when `service_type = "buffet"` or `participant_mode = "headcount"`. See below when applicable |

#### `meals.sessions[].individual_orders[]`

Populated only when `service_type = "individual_order"` or `"mixed"` **and** `participant_mode = "detailed"`. Otherwise always an empty array `[]`.

| Field | Type | Always Present | Notes |
|---|---|---|---|
| `attendant_idx` | `integer` | ✅ | Index into the root `attendants` array |
| `menu_item_id` | `string` | ✅ | |
| `quantity` | `integer` | ✅ | Min 1 |
| `notes` | `string \| null` | ✅ | Per-item preparation or allergy notes |

---

## Enum Reference

### `event_type`
| Value | Label |
|---|---|
| `conference` | Conference |
| `seminar` | Seminar |
| `workshop` | Workshop |
| `gala` | Gala / Dinner |
| `wedding` | Wedding |
| `training` | Training |

### `setup_type`
| Value | Label |
|---|---|
| `boardroom` | Boardroom |
| `theatre` | Theatre |
| `classroom` | Classroom |
| `u_shape` | U-Shape |
| `banquet` | Banquet |
| `cocktail` | Cocktail |

### `pricing_basis`
| Value | Label |
|---|---|
| `half_day` | Half Day |
| `full_day` | Full Day |
| `hourly` | Hourly |
| `flat_rate` | Flat Rate |

### `meal_period`
| Value | Label |
|---|---|
| `breakfast` | Breakfast |
| `lunch` | Lunch |
| `dinner` | Dinner |
| `tea_break` | Tea Break |
| `cocktail` | Cocktail |

### `service_type`
| Value | Label | `individual_orders` populated? |
|---|---|---|
| `buffet` | Buffet | No — always `[]` |
| `individual_order` | Individual Orders | Yes — when `participant_mode = "detailed"` |
| `mixed` | Mixed (Buffet + Exceptions) | Yes — exception orders only, when `participant_mode = "detailed"` |

### `room_type` (accommodation preference)
| Value | Label |
|---|---|
| `single` | Single |
| `double` | Double |
| `twin` | Twin |
| `suite` | Suite |
| `null` | No preference — property decides |

---

## Key Differences from Individual Booking

| Aspect | Individual | Corporate |
|---|---|---|
| `booking_type` | `"individual"` | `"corporate"` |
| Company section | ❌ Not present | ✅ Full `company` object |
| Approver section | ❌ Not present | ✅ `approver` object |
| Booker `job_title` | ❌ Not present | ✅ Present |
| Attendant `company` field | ❌ Not present | ✅ Present |
| Accommodation model | Per-room / per-attendant assignments | Room type + count only; property assigns rooms |
| Accommodation `reason_for_booking` | ❌ Not present | ✅ Present |
| Events `start_date` / `end_date` | Not at events level (derived from sessions) | ✅ Explicit on `events` object |
| Meals `start_date` / `end_date` (standalone) | Present | Present |
| API endpoint | `POST /guest/bookings/individual` | `POST /guest/bookings/corporate-event` |

---

## Fields Added to Contract vs. Current `submit()` Implementation

The following fields belong in the full contract but are **not yet sent** by the current frontend `submit()` function. These must be added before going live:

| Field path | Currently sent | Action required |
|---|---|---|
| `source` | ❌ | Add hardcoded `"web"` |
| `currency` | ❌ | Add hardcoded `"ZMW"` |
| `participant_mode` | ❌ | Add from store |
| `participant_count` | ✅ (headcount only) | Send `null` in detailed mode instead of omitting |
| `company` (nested object) | ❌ Sent flat at root | Nest all company/approver fields under `company` object |
| `approver` (nested object) | ❌ Sent flat at root as `approver_name`, etc. | Nest under `approver` object |
| `events.start_date` | ❌ | Add from `events.value.startDate` |
| `events.end_date` | ❌ | Add from `events.value.endDate` |
| `events.sessions[].venue_name` | ❌ | Stored as `s.venueName`; add to session payload |
| `events.sessions[].venue_capacity` | ❌ | Stored as `s.venueCapacity`; add to session payload |
| `meals.start_date` | ❌ | Add from `meals.value.startDate` (standalone mode) |
| `meals.end_date` | ❌ | Add from `meals.value.endDate` (standalone mode) |
| `meals.meal_mode` | ❌ | Add from `meals.value.mealMode` |
| `meals.sessions[].arrangements_notes` | ✅ Already sent | No change needed |
| `accommodation` | Omitted when disabled | Send `null` instead |
| `events` | Omitted when disabled | Send `null` instead |
| `meals` | Omitted when disabled | Send `null` instead |
