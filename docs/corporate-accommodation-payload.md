# Corporate Accommodation Booking — API Payload Contract

**Endpoint:** `POST /api/v1/guest/bookings/accommodation`

## Design Rules

- `booking_context` is always `"corporate"` for these payloads.
- `company` and `approver` objects are always present and fully populated (never `null`) for corporate bookings.
- `accommodation.rooms` is always `null` for corporate — the property assigns specific rooms. Only `room_count` and `room_type_preference` are used.
- In **headcount mode**, `attendants` contains exactly one record — the person submitting the booking — and `participant_count` holds the total delegate count.
- In **detailed mode**, `attendants` contains one record per registered delegate; `participant_count` is `null`.
- The `company` field in each attendant record holds the delegate's **job title** (field name reuse; not an organisation name).
- Supporting documents (approval letters, LPOs, etc.) are uploaded separately as `multipart/form-data` — see the **Document Upload** section at the bottom of this file.

---

## Scenario 1 — Headcount Only

The company provides a total delegate count. Individual delegate details are not recorded.

```json
{
  "org_id": "5",
  "branch_id": "12",
  "booking_type": "accommodation",
  "source": "web",
  "currency": "ZMW",
  "booking_context": "corporate",

  "participant_mode": "headcount",
  "participant_count": 24,

  "booked_by": {
    "name": "Grace Mwansa",
    "email": "grace.mwansa@acmecorp.zm",
    "phone": "+260971234567",
    "job_title": "Procurement Officer"
  },

  "attendants": [
    {
      "full_name": "Grace Mwansa",
      "email": "grace.mwansa@acmecorp.zm",
      "phone": "+260971234567",
      "id_number": null,
      "dietary_notes": null,
      "company": null,
      "is_lead_contact": true
    }
  ],

  "company": {
    "name": "Acme Corporation Zambia Ltd",
    "tpin": "1003456789",
    "industry": "Mining & Extractives",
    "email": "accounts@acmecorp.zm",
    "phone": "+260211345678",
    "city": "Lusaka",
    "street_address": "Plot 1234 Cairo Road, Lusaka",
    "branch_name": "Head Office",
    "department_name": "Human Resources",
    "cost_center": "CC-HR-001",
    "gl_code": "GL-7700"
  },

  "approver": {
    "name": "Daniel Phiri",
    "email": "d.phiri@acmecorp.zm",
    "phone": "+260977654321",
    "title": "Finance Manager"
  },

  "accommodation": {
    "check_in": "2026-08-10",
    "check_out": "2026-08-14",
    "notes": "All delegates require non-smoking rooms. Please arrange early check-in for the group.",
    "room_count": 12,
    "room_type_preference": "twin",
    "rooms": null
  }
}
```

---

## Scenario 2 — Individual Records (Detailed Mode)

Each delegate is registered with their personal details. `participant_count` is `null`.

```json
{
  "org_id": "5",
  "branch_id": "12",
  "booking_type": "accommodation",
  "source": "web",
  "currency": "ZMW",
  "booking_context": "corporate",

  "participant_mode": "detailed",
  "participant_count": null,

  "booked_by": {
    "name": "Grace Mwansa",
    "email": "grace.mwansa@acmecorp.zm",
    "phone": "+260971234567",
    "job_title": "Procurement Officer"
  },

  "attendants": [
    {
      "full_name": "Grace Mwansa",
      "email": "grace.mwansa@acmecorp.zm",
      "phone": "+260971234567",
      "id_number": "123456/78/1",
      "dietary_notes": null,
      "company": "Procurement Officer",
      "is_lead_contact": true
    },
    {
      "full_name": "Kenneth Zulu",
      "email": "k.zulu@acmecorp.zm",
      "phone": "+260966112233",
      "id_number": "234567/89/1",
      "dietary_notes": "Vegetarian",
      "company": "Senior Engineer",
      "is_lead_contact": false
    },
    {
      "full_name": "Mutale Banda",
      "email": "m.banda@acmecorp.zm",
      "phone": null,
      "id_number": null,
      "dietary_notes": "Halal",
      "company": "Project Coordinator",
      "is_lead_contact": false
    },
    {
      "full_name": "Chileshe Nkowane",
      "email": null,
      "phone": "+260955887766",
      "id_number": "345678/90/1",
      "dietary_notes": null,
      "company": "Finance Analyst",
      "is_lead_contact": false
    }
  ],

  "company": {
    "name": "Acme Corporation Zambia Ltd",
    "tpin": "1003456789",
    "industry": "Mining & Extractives",
    "email": "accounts@acmecorp.zm",
    "phone": "+260211345678",
    "city": "Lusaka",
    "street_address": "Plot 1234 Cairo Road, Lusaka",
    "branch_name": "Head Office",
    "department_name": "Human Resources",
    "cost_center": "CC-HR-001",
    "gl_code": "GL-7700"
  },

  "approver": {
    "name": "Daniel Phiri",
    "email": "d.phiri@acmecorp.zm",
    "phone": "+260977654321",
    "title": "Finance Manager"
  },

  "accommodation": {
    "check_in": "2026-08-10",
    "check_out": "2026-08-14",
    "notes": "All delegates require non-smoking rooms. Please arrange early check-in for the group.",
    "room_count": 4,
    "room_type_preference": "double",
    "rooms": null
  }
}
```

---

## Field Reference

### Root Object

| Field | Type | Notes |
|---|---|---|
| `org_id` | `string` | Lodge / property identifier |
| `branch_id` | `string \| null` | Property branch; `null` if not applicable |
| `booking_type` | `"accommodation"` | Fixed value |
| `source` | `string` | `"web"` \| `"mobile"` \| `"reception"` |
| `currency` | `string` | ISO 4217 — e.g. `"ZMW"` |
| `booking_context` | `"corporate"` | Fixed value for corporate bookings |
| `participant_mode` | `string` | `"headcount"` \| `"detailed"` |
| `participant_count` | `integer \| null` | Total delegate count. Non-null only when `participant_mode = "headcount"` |
| `booked_by` | `object` | Person submitting the booking — see below |
| `attendants` | `array` | Delegate records — see below |
| `company` | `object` | Company snapshot — see below. Always present for corporate |
| `approver` | `object` | Authorising person — see below. Always present for corporate |
| `accommodation` | `object` | Always present for accommodation bookings — see below |

---

### `booked_by`

| Field | Type | Notes |
|---|---|---|
| `name` | `string` | Required |
| `email` | `string \| null` | |
| `phone` | `string \| null` | |
| `job_title` | `string \| null` | Role of the person submitting the booking |

---

### `attendants[]`

In **headcount mode**, the array contains exactly one record — the booker — with `company`, `id_number`, and `dietary_notes` all `null`.

In **detailed mode**, the array contains one object per registered delegate. The submitter's record appears first with `is_lead_contact: true`.

> **Note:** The `company` field on each attendant holds the delegate's **job title**, not a company name. This is a field reuse inherited from the shared attendant schema.

| Field | Type | Notes |
|---|---|---|
| `full_name` | `string \| null` | Required in detailed mode |
| `email` | `string \| null` | |
| `phone` | `string \| null` | |
| `id_number` | `string \| null` | National ID or passport number |
| `dietary_notes` | `string \| null` | e.g. `"Vegetarian"`, `"Halal"`, `"Nut allergy"` |
| `company` | `string \| null` | Delegate's job title. `null` in headcount mode |
| `is_lead_contact` | `boolean` | Exactly one delegate carries `true` (the booker) |

---

### `company`

Always fully present for corporate bookings. Represents a snapshot of the company details at the time of booking.

| Field | Type | Notes |
|---|---|---|
| `name` | `string` | Company / organisation name |
| `tpin` | `string \| null` | Zambia Revenue Authority Tax Payer Identification Number |
| `industry` | `string \| null` | One of the industry enum values — see below |
| `email` | `string \| null` | Billing / accounts email |
| `phone` | `string \| null` | Company switchboard number |
| `city` | `string \| null` | |
| `street_address` | `string \| null` | |
| `branch_name` | `string \| null` | Internal company branch or division |
| `department_name` | `string \| null` | Requesting department |
| `cost_center` | `string \| null` | Internal cost centre code |
| `gl_code` | `string \| null` | General ledger code for accounting allocation |

---

### `approver`

The person within the company who has authorised this booking. Always fully present for corporate bookings.

| Field | Type | Notes |
|---|---|---|
| `name` | `string` | Required |
| `email` | `string \| null` | |
| `phone` | `string \| null` | |
| `title` | `string \| null` | Job title / role of the approver |

---

### `accommodation`

| Field | Type | Notes |
|---|---|---|
| `check_in` | `string` | ISO 8601 date — `YYYY-MM-DD` |
| `check_out` | `string` | ISO 8601 date — `YYYY-MM-DD` |
| `notes` | `string \| null` | Special requests or group requirements |
| `room_count` | `integer` | Total rooms required. Always present for corporate |
| `room_type_preference` | `string \| null` | e.g. `"twin"`, `"double"`, `"suite"`. Null if no preference |
| `rooms` | `null` | Always `null` for corporate — property assigns specific rooms |

---

## Enum Reference

### `industry`

| Value |
|---|
| `Agriculture & Agribusiness` |
| `Banking & Finance` |
| `Construction & Infrastructure` |
| `Education & Training` |
| `Energy & Utilities` |
| `Government & Public Sector` |
| `Healthcare & Medical` |
| `Hospitality & Tourism` |
| `Information Technology` |
| `Legal & Professional Services` |
| `Manufacturing` |
| `Media & Communications` |
| `Mining & Extractives` |
| `NGO & Non-profit` |
| `Real Estate` |
| `Retail & Trade` |
| `Telecommunications` |
| `Transportation & Logistics` |
| `Other` |

---

## Document Upload

Supporting documents (LPOs, authorisation letters, budget approvals, travel requests) attached by the user are **not** included in the JSON payload above. They must be submitted as a separate `multipart/form-data` request or appended to the booking creation request as form fields alongside the JSON body, depending on the backend implementation.

**Accepted MIME types:** `application/pdf`, `application/msword`, `application/vnd.openxmlformats-officedocument.wordprocessingml.document`, `image/jpeg`, `image/png`

**Maximum size per file:** 10 MB

Suggested form field name: `approval_documents[]` (array of files)

> Coordinate the exact upload strategy (separate endpoint vs. multipart booking creation) with the backend team before implementation.

---

## Key Behavioural Differences vs. Individual Booking

| Behaviour | Individual | Corporate |
|---|---|---|
| `booking_context` | `"individual"` | `"corporate"` |
| `company` object | `null` | Always present |
| `approver` object | `null` | Always present |
| `booked_by.job_title` | Not sent | Always present |
| Room assignment | Per attendant (`rooms[]`) | Total count only (`room_count`) |
| `accommodation.rooms` | Array of room assignments | Always `null` |
| `room_count` | Always `null` | Integer ≥ 1 |
| `room_type_preference` | Always `null` | String or `null` |
