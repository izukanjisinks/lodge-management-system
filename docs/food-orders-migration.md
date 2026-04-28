# Food & Orders System — Migration Notes

## What Changed

The old system attached a **meal plan** to a booking at creation time. A fixed nightly rate per person was billed when the invoice was generated at checkout.

The new system replaces this with a **menu/orders model**: staff place individual food orders during a guest's stay and the invoice updates in real time.

---

## Invoicing Behaviour

### Old
- Booking created with optional `meal_plan_id`
- On confirmation, invoice generated with two line items: room cost + meal plan cost
- Meal cost was fixed: `nights × guests × price_per_person_per_night`

### New
- Booking created with no meal plan — room only
- On confirmation, invoice generated with **one line item**: room cost
- As food orders are placed during the stay, **each order appends its own line item** to the invoice immediately
- Invoice totals (subtotal, tax, total) are recalculated atomically on every append
- Walk-in orders (no booking) do **not** touch any invoice

**There is no "close order" step.** The invoice is always up to date.

---

## Database Changes (Migrations to Run)

| Migration | What it does |
|-----------|-------------|
| `000033_create_menus_and_orders` | Creates `menus`, `menu_items`, `orders`, `order_items`; adds `order_id` column to `invoice_line_items` |
| `000034_drop_meal_plans` | Drops `booking_meal_plans` and `meal_plans` tables |
| `000035_simplify_orders` | Drops `status` column and `order_status` enum from `orders` |

Run all three in sequence.

---

## API Changes

### Removed endpoints
```
GET  /api/v1/meal-plans
GET  /api/v1/meal-plans/{id}
POST /api/v1/meal-plans
PUT  /api/v1/meal-plans/{id}
DELETE /api/v1/meal-plans/{id}

GET  /api/v1/guest/meal-plans
GET  /api/v1/guest/meal-plans/{id}
```

### New staff endpoints
```
# Menus (admin/manager write, all staff read)
GET    /api/v1/menus
GET    /api/v1/menus/{id}
POST   /api/v1/menus
PUT    /api/v1/menus/{id}
DELETE /api/v1/menus/{id}

# Menu items (nested under menu)
POST   /api/v1/menus/{id}/items
PUT    /api/v1/menus/{id}/items/{item_id}
DELETE /api/v1/menus/{id}/items/{item_id}

# Orders (admin/manager/receptionist)
GET    /api/v1/orders                    # list, filter by ?type=in_house|walk_in&booking_id=
GET    /api/v1/orders/{id}
POST   /api/v1/orders                    # in-house order (requires booking_id)
POST   /api/v1/orders/walk-in            # walk-in order (no booking)
POST   /api/v1/orders/{id}/items         # add more items to existing order
```

### New public endpoint
```
GET /api/v1/guest/menus?org_id=<uuid>    # optional org filter, returns active menus + available items
```

---

## Request Body Changes

### Create Booking — `meal_plan_id` removed
```json
// Before
{
  "room_id": "...",
  "client_id": "...",
  "client_type": "individual",
  "check_in": "2026-05-01T00:00:00Z",
  "check_out": "2026-05-05T00:00:00Z",
  "guests": 2,
  "meal_plan_id": "..."       // ← removed
}

// Now
{
  "room_id": "...",
  "client_id": "...",
  "client_type": "individual",
  "check_in": "2026-05-01T00:00:00Z",
  "check_out": "2026-05-05T00:00:00Z",
  "guests": 2
}
```

### Place In-House Order — new
```json
POST /api/v1/orders
{
  "booking_id": "<uuid>",
  "notes": "Table 4",
  "items": [
    { "menu_item_id": "<uuid>", "quantity": 2, "notes": "no onions" },
    { "menu_item_id": "<uuid>", "quantity": 1 }
  ]
}
```
Booking must be `confirmed` or `checked_in`. Invoice updated immediately.

### Place Walk-In Order — new
```json
POST /api/v1/orders/walk-in
{
  "notes": "Counter seat 2",
  "items": [
    { "menu_item_id": "<uuid>", "quantity": 1 }
  ]
}
```

### Add Items to Existing Order — new
```json
POST /api/v1/orders/{id}/items
{
  "items": [
    { "menu_item_id": "<uuid>", "quantity": 1, "notes": "extra sauce" }
  ]
}
```

---

## Invoice Line Item Shape

Each line item now includes an optional `order_id` field:

```json
{
  "id": "...",
  "invoice_id": "...",
  "order_id": "<uuid or null>",
  "description": "Food & Beverage — Order ORD-001000 (3 item(s))",
  "quantity": 1,
  "unit_price": 135.00,
  "total": 135.00,
  "created_at": "..."
}
```

The room line item will always have `order_id: null`. Food line items will always have an `order_id`.

---

## Front-End Checklist

- [ ] Remove all meal plan selection UI from the booking creation flow
- [ ] Remove `meal_plan_id` from any booking request payloads
- [ ] Remove meal plan display from booking detail views
- [ ] Replace meal plan list/detail pages with menu list page (`GET /api/v1/guest/menus?org_id=`)
- [ ] Add staff order placement UI (select booking → pick items from menu → POST /orders)
- [ ] Add "add more items" UI on order detail (POST /orders/{id}/items)
- [ ] Invoice detail view: render `order_id` line items as food charges (already works — just new line items appearing)
- [ ] Invoice totals now update live as orders are placed — no need to refresh on checkout
