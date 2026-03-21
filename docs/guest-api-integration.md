# Guest API Integration Guide — The Sanctuary Lodge

This document covers every backend endpoint the website frontend needs, from browsing rooms to managing reservations. All endpoints are prefixed with `/api/v1`.

Base URL (development): `http://localhost:8081/api/v1`

---

## Authentication

### How it works

- Login returns a JWT token. Include it on every authenticated request:
  ```
  Authorization: Bearer <token>
  ```
- Guests have role `guest`. The token payload includes the role so the frontend can distinguish guest users from staff.
- Token expiry is configurable. Check `expires_at` in the login response.

---

## 1. Register

Creates a new guest account and individual profile atomically. Returns a JWT — the guest is logged in immediately after registration.

**`POST /api/v1/guest/register`** — Public

### Request
```json
{
  "full_name": "John Mwewa",
  "email": "john@example.com",
  "password": "SecurePass@123",
  "phone": "+260971234567",
  "id_passport_number": "NRC123456/78/1",
  "nationality": "Zambian"
}
```

| Field | Required | Notes |
|---|---|---|
| `full_name` | Yes | |
| `email` | Yes | Must be unique |
| `password` | Yes | Subject to password policy |
| `phone` | Yes | |
| `id_passport_number` | No | NRC or passport |
| `nationality` | No | |

### Response `201`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2026-04-21T10:00:00Z",
  "user": {
    "id": "uuid",
    "full_name": "John Mwewa",
    "email": "john@example.com",
    "role": "guest",
    "status": "active",
    "created_at": "2026-03-21T10:00:00Z"
  }
}
```

---

## 2. Login

**`POST /api/v1/auth/login`** — Public

### Request
```json
{
  "email": "john@example.com",
  "password": "SecurePass@123"
}
```

### Response `200`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2026-04-21T10:00:00Z",
  "user": {
    "id": "uuid",
    "full_name": "John Mwewa",
    "email": "john@example.com",
    "role": "guest",
    "status": "active",
    "change_password": false
  }
}
```

---

## 3. Guest Profile

### Get my profile
**`GET /api/v1/guest/me`** — Requires auth (guest)

### Response `200`
```json
{
  "user": {
    "id": "uuid",
    "full_name": "John Mwewa",
    "email": "john@example.com",
    "role": "guest",
    "status": "active",
    "created_at": "2026-03-21T10:00:00Z"
  },
  "profile": {
    "id": "uuid",
    "full_name": "John Mwewa",
    "email": "john@example.com",
    "phone": "+260971234567",
    "id_passport_number": "NRC123456/78/1",
    "nationality": "Zambian",
    "status": "active",
    "created_at": "2026-03-21T10:00:00Z",
    "updated_at": "2026-03-21T10:00:00Z"
  }
}
```

### Update my profile
**`PUT /api/v1/guest/me`** — Requires auth (guest)

Only the fields below can be updated. Email cannot be changed (it is the login identity).

### Request
```json
{
  "full_name": "John K. Mwewa",
  "phone": "+260971234999",
  "id_passport_number": "NRC123456/78/1",
  "nationality": "Zambian"
}
```

All fields are optional — only non-empty values are applied.

### Response `200`
Returns the updated `IndividualClient` profile object (same shape as `profile` in the Get response above).

---

## 4. Rooms

All room endpoints are **public** — no authentication required.

### List all rooms
**`GET /api/v1/rooms`**

#### Query params
| Param | Type | Notes |
|---|---|---|
| `type` | string | Filter by room type: `single`, `double`, `suite`, `cabin`, `conference` |
| `is_available` | boolean | `true` or `false` |
| `page` | int | Default: 1 |
| `page_size` | int | Default: 20, max: 100 |

### Response `200`
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Suite 301",
      "type": "suite",
      "capacity": 4,
      "price_per_night": 520.00,
      "amenities": ["WiFi", "Air Conditioning", "TV", "Mini Bar", "Jacuzzi", "Lounge Area"],
      "is_available": true,
      "description": "Luxury suite with separate lounge and jacuzzi",
      "created_at": "2026-03-21T10:00:00Z",
      "updated_at": "2026-03-21T10:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 7
}
```

### Get a single room
**`GET /api/v1/rooms/{id}`**

### Response `200`
Returns a single room object (same shape as items in the list above).

### List available rooms for dates
**`GET /api/v1/rooms/available`**

#### Query params
| Param | Type | Required | Notes |
|---|---|---|---|
| `check_in` | string | Yes | Format: `YYYY-MM-DD` |
| `check_out` | string | Yes | Format: `YYYY-MM-DD` |
| `type` | string | No | Room type filter |

### Response `200`
Returns a plain array (not paginated) of available room objects.

---

## 5. Meal Plans

All meal plan endpoints are **public** — no authentication required.

### List all meal plans
**`GET /api/v1/meal-plans`**

#### Query params
| Param | Type | Notes |
|---|---|---|
| `is_active` | boolean | Filter active plans only (recommended: `true`) |
| `page` | int | Default: 1 |
| `page_size` | int | Default: 20 |

### Response `200`
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Half Board",
      "price_per_person_per_night": 55.00,
      "includes": ["Breakfast", "Dinner"],
      "description": "Breakfast and dinner included",
      "is_active": true,
      "created_at": "2026-03-21T10:00:00Z",
      "updated_at": "2026-03-21T10:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 3
}
```

### Get a single meal plan
**`GET /api/v1/meal-plans/{id}`**

### Response `200`
Returns a single meal plan object (same shape as items in the list above).

---

## 6. Bookings (Guest)

All booking endpoints require authentication (`guest` role).

### Create a booking
**`POST /api/v1/guest/bookings`** — Requires auth (guest)

`client_id` and `client_type` are resolved automatically from the JWT — the guest does not send them.

### Request
```json
{
  "room_id": "uuid",
  "check_in": "2026-04-01T00:00:00Z",
  "check_out": "2026-04-05T00:00:00Z",
  "guests": 2,
  "meal_plan_id": "uuid",
  "special_requests": "Late check-in, arriving after 10 PM"
}
```

| Field | Required | Notes |
|---|---|---|
| `room_id` | Yes | Must be an available room |
| `check_in` | Yes | ISO 8601 datetime |
| `check_out` | Yes | Must be after check_in |
| `guests` | Yes | Must not exceed room capacity |
| `meal_plan_id` | No | UUID of chosen meal plan |
| `special_requests` | No | Free text |

### Response `201`
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "room_id": "uuid",
  "room_name": "Suite 301",
  "client_id": "uuid",
  "client_type": "individual",
  "client_name": "John Mwewa",
  "meal_plan_id": "uuid",
  "meal_plan_name": "Half Board",
  "check_in": "2026-04-01T00:00:00Z",
  "check_out": "2026-04-05T00:00:00Z",
  "guests": 2,
  "nights": 4,
  "room_cost": 2080.00,
  "meal_cost": 440.00,
  "total_amount": 2520.00,
  "status": "pending",
  "special_requests": "Late check-in, arriving after 10 PM",
  "created_at": "2026-03-21T10:00:00Z",
  "updated_at": "2026-03-21T10:00:00Z"
}
```

### List my bookings
**`GET /api/v1/guest/bookings`** — Requires auth (guest)

Returns only bookings belonging to the logged-in guest.

#### Query params
| Param | Type | Notes |
|---|---|---|
| `page` | int | Default: 1 |
| `page_size` | int | Default: 20 |

### Response `200`
```json
{
  "data": [ /* array of booking objects */ ],
  "page": 1,
  "page_size": 20,
  "total": 3
}
```

### Get a single booking
**`GET /api/v1/guest/bookings/{id}`** — Requires auth (guest)

Returns `403` if the booking does not belong to the logged-in guest.

### Response `200`
Returns a single booking object (same shape as the create response).

### Cancel a booking
**`PATCH /api/v1/guest/bookings/{id}/cancel`** — Requires auth (guest)

Only bookings with status `pending` or `confirmed` can be cancelled.

### Response `200`
```json
{
  "message": "Booking cancelled successfully"
}
```

### Errors
| Status | Reason |
|---|---|
| `400` | Booking is in a status that cannot be cancelled (e.g. `checked_in`) |
| `403` | Booking belongs to a different guest |
| `404` | Booking not found |

---

## 7. Booking Status Lifecycle

```
pending → confirmed → checked_in → checked_out
pending → cancelled
confirmed → cancelled
```

Status transitions after `confirmed` are managed by lodge staff only. Guests can only cancel from `pending` or `confirmed`.

---

## 8. Error Responses

All errors follow this shape:

```json
{
  "error": {
    "message": "room is not available for the selected dates"
  }
}
```

---

## 9. Quick Reference

| Method | Endpoint | Auth | Purpose |
|---|---|---|---|
| `POST` | `/api/v1/guest/register` | Public | Register + auto-login |
| `POST` | `/api/v1/auth/login` | Public | Login |
| `GET` | `/api/v1/guest/me` | Guest | Get profile |
| `PUT` | `/api/v1/guest/me` | Guest | Update profile |
| `GET` | `/api/v1/rooms` | Public | Browse rooms |
| `GET` | `/api/v1/rooms/{id}` | Public | Room detail |
| `GET` | `/api/v1/rooms/available` | Public | Available rooms for dates |
| `GET` | `/api/v1/meal-plans` | Public | Browse meal plans |
| `GET` | `/api/v1/meal-plans/{id}` | Public | Meal plan detail |
| `POST` | `/api/v1/guest/bookings` | Guest | Create booking |
| `GET` | `/api/v1/guest/bookings` | Guest | My bookings |
| `GET` | `/api/v1/guest/bookings/{id}` | Guest | Single booking |
| `PATCH` | `/api/v1/guest/bookings/{id}/cancel` | Guest | Cancel booking |
