# Accommodation Booking — Individual Mode API Payload

**Endpoint:** `POST /api/v1/guest/bookings/accommodation`  
**Auth:** Bearer token required  
**Content-Type:** `application/json`

---

## Headcount Mode

Use when the booker does not register individual guest details — only the total party size and the room(s) they need.

```json
{
  "org_id": "eb3b35d7-e5b4-4ffd-85d4-82ed7bad37be",
  "branch_id": "20168397-c85e-437c-a1be-b3266c8e5dd1",
  "booking_type": "accommodation",
  "source": "web",
  "currency": "ZMW",
  "booking_context": "individual",

  "participant_mode": "headcount",
  "participant_count": 3,

  "booked_by": {
    "name": "Martin Sinkolongo",
    "email": "martin@example.com",
    "phone": "+260971234567",
    "job_title": null
  },

  "attendants": [
    {
      "full_name": "Martin Sinkolongo",
      "email": "martin@example.com",
      "phone": "+260971234567",
      "id_number": null,
      "dietary_notes": null,
      "company": null,
      "is_lead_contact": true
    }
  ],

  "company": null,
  "approver": null,

  "accommodation": {
    "check_in": "2026-07-15",
    "check_out": "2026-07-18",
    "notes": "Early check-in requested if possible",
    "room_count": null,
    "room_type_preference": null,
    "rooms": [
      {
        "slot_index": 0,
        "attendant_idx": 0,
        "room_id": "6e7728d8-ab3d-467b-a275-a84d469fddab",
        "room_name": "Safari Suite",
        "room_type": "suite",
        "rate_per_night": 1250.00
      },
      {
        "slot_index": 1,
        "attendant_idx": 1,
        "room_id": "4530e720-500a-4b30-9427-cfc206fc9f6d",
        "room_name": "Savannah Twin",
        "room_type": "twin",
        "rate_per_night": 850.00
      }
    ]
  }
}
```

**Notes — Headcount Mode:**
- `participant_count` is the total party size including the booker.
- `attendants` always contains exactly **one record** — the booker — with all optional fields as `null`.
- `rooms` contains one entry per room slot the booker created. Multiple rooms may be added independently of guest count.
- `slot_index` and `attendant_idx` both reflect the room's position in the slot list (0-based). When there is only one room, both are `0`.

---

## Detailed (Individual Records) Mode

Use when the booker registers each guest individually with their name and contact details. Each guest can be assigned their own room.

```json
{
  "org_id": "eb3b35d7-e5b4-4ffd-85d4-82ed7bad37be",
  "branch_id": "20168397-c85e-437c-a1be-b3266c8e5dd1",
  "booking_type": "accommodation",
  "source": "web",
  "currency": "ZMW",
  "booking_context": "individual",

  "participant_mode": "detailed",
  "participant_count": null,

  "booked_by": {
    "name": "Martin Sinkolongo",
    "email": "martin@example.com",
    "phone": "+260971234567",
    "job_title": null
  },

  "attendants": [
    {
      "full_name": "Martin Sinkolongo",
      "email": "martin@example.com",
      "phone": "+260971234567",
      "id_number": "123456/78/1",
      "dietary_notes": null,
      "company": null,
      "is_lead_contact": true
    },
    {
      "full_name": "Jane Doe",
      "email": "jane@example.com",
      "phone": "+260976543210",
      "id_number": "654321/10/2",
      "dietary_notes": "Vegetarian, no nuts",
      "company": null,
      "is_lead_contact": false
    },
    {
      "full_name": "John Banda",
      "email": null,
      "phone": null,
      "id_number": null,
      "dietary_notes": null,
      "company": "Acme Ltd",
      "is_lead_contact": false
    }
  ],

  "company": null,
  "approver": null,

  "accommodation": {
    "check_in": "2026-07-15",
    "check_out": "2026-07-18",
    "notes": "Ground floor rooms preferred",
    "room_count": null,
    "room_type_preference": null,
    "rooms": [
      {
        "slot_index": 0,
        "attendant_idx": 0,
        "room_id": "6e7728d8-ab3d-467b-a275-a84d469fddab",
        "room_name": "Safari Suite",
        "room_type": "suite",
        "rate_per_night": 1250.00
      },
      {
        "slot_index": 1,
        "attendant_idx": 1,
        "room_id": "4530e720-500a-4b30-9427-cfc206fc9f6d",
        "room_name": "Savannah Twin",
        "room_type": "twin",
        "rate_per_night": 850.00
      },
      {
        "slot_index": 2,
        "attendant_idx": 2,
        "room_id": "816af3dc-88d7-4f96-a18e-87a7185d8819",
        "room_name": "Riverside Double",
        "room_type": "double",
        "rate_per_night": 920.00
      }
    ]
  }
}
```

**Notes — Detailed Mode:**
- `participant_count` is always `null` in detailed mode; the guest count is inferred from `attendants.length`.
- The `attendants` array contains one record per guest. The booker appears first with `is_lead_contact: true`.
- `rooms[].slot_index` and `rooms[].attendant_idx` both equal the attendant's position in the `attendants` array (0-based).
- Not every attendant is required to have a room. If an attendant has no room selected, no entry appears in `rooms` for their index.

---

## Minimum Viable Payload

The smallest valid payload — headcount mode, single room, no notes:

```json
{
  "org_id": "eb3b35d7-e5b4-4ffd-85d4-82ed7bad37be",
  "branch_id": null,
  "booking_type": "accommodation",
  "source": "web",
  "currency": "ZMW",
  "booking_context": "individual",

  "participant_mode": "headcount",
  "participant_count": 1,

  "booked_by": {
    "name": "Martin Sinkolongo",
    "email": "martin@example.com",
    "phone": null,
    "job_title": null
  },

  "attendants": [
    {
      "full_name": "Martin Sinkolongo",
      "email": "martin@example.com",
      "phone": null,
      "id_number": null,
      "dietary_notes": null,
      "company": null,
      "is_lead_contact": true
    }
  ],

  "company": null,
  "approver": null,

  "accommodation": {
    "check_in": "2026-07-15",
    "check_out": "2026-07-16",
    "notes": null,
    "room_count": null,
    "room_type_preference": null,
    "rooms": [
      {
        "slot_index": 0,
        "attendant_idx": 0,
        "room_id": "6e7728d8-ab3d-467b-a275-a84d469fddab",
        "room_name": "Safari Suite",
        "room_type": "suite",
        "rate_per_night": 1250.00
      }
    ]
  }
}
```

---

## Field Reference

### Root

| Field | Type | Notes |
|---|---|---|
| `org_id` | `string` | Lodge / property UUID |
| `branch_id` | `string \| null` | Branch UUID; `null` if not applicable |
| `booking_type` | `"accommodation"` | Fixed value |
| `source` | `"web"` | Fixed value from this client |
| `currency` | `"ZMW"` | Fixed value |
| `booking_context` | `"individual"` | Fixed for this flow; corporate has its own flow |
| `participant_mode` | `"headcount" \| "detailed"` | Determines attendants and room indexing |
| `participant_count` | `integer \| null` | Set when `participant_mode = "headcount"`; `null` in detailed mode |
| `booked_by` | `object` | Person submitting the booking |
| `attendants` | `array` | Min 1 element. In headcount: booker only. In detailed: all registered guests |
| `company` | `null` | Always `null` for individual bookings |
| `approver` | `null` | Always `null` for individual bookings |
| `accommodation` | `object` | Accommodation details — always present in this flow |

---

### `booked_by`

| Field | Type | Notes |
|---|---|---|
| `name` | `string` | Required — full name |
| `email` | `string \| null` | Required |
| `phone` | `string \| null` | Optional |
| `job_title` | `null` | Always `null` for individual bookings |

---

### `attendants[]`

| Field | Type | Notes |
|---|---|---|
| `full_name` | `string \| null` | |
| `email` | `string \| null` | |
| `phone` | `string \| null` | |
| `id_number` | `string \| null` | National ID or passport number |
| `dietary_notes` | `string \| null` | Free text |
| `company` | `string \| null` | Guest's employer — used for delegation lists |
| `is_lead_contact` | `boolean` | Exactly one attendant per booking carries `true` |

---

### `accommodation`

| Field | Type | Notes |
|---|---|---|
| `check_in` | `string` | `YYYY-MM-DD` |
| `check_out` | `string` | `YYYY-MM-DD` |
| `notes` | `string \| null` | Special requests — accessibility, floor, quiet room, etc. |
| `room_count` | `null` | Always `null` for individual bookings (used by corporate flow only) |
| `room_type_preference` | `null` | Always `null` for individual bookings (used by corporate flow only) |
| `rooms` | `array` | One entry per room selected. See below |

---

### `accommodation.rooms[]`

| Field | Type | Notes |
|---|---|---|
| `slot_index` | `integer` | 0-based position of this room in the selection list |
| `attendant_idx` | `integer` | Index of the guest in `attendants` this room is assigned to. In headcount mode mirrors `slot_index`. |
| `room_id` | `string` | UUID from the availability API |
| `room_name` | `string \| null` | Snapshot of room name at booking time |
| `room_type` | `string \| null` | e.g. `"suite"`, `"double"`, `"twin"`, `"single"` |
| `rate_per_night` | `number \| null` | ZMW rate at time of booking |
