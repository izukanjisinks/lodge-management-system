#!/usr/bin/env bash
#
# End-to-end test for the corporate MEALS booking flow.
#
#   submit (public) -> view enriched task (staff) -> approve (staff, auto-creates
#   booking + orders + draft invoice) -> verify booking / orders / invoice.
#
# Prerequisites:
#   - API running with the meals code (restart after the latest changes).
#   - A staff JWT (admin / branch_admin / manager) for the menu + approve calls.
#   - jq installed (https://jqlang.github.io/jq/).
#
# Fill in the three vars below, then: bash scripts/test_meals_flow.sh
# Optionally pass MODE=buffet to test the buffet path instead of itemised:
#   MODE=buffet bash scripts/test_meals_flow.sh
#
set -euo pipefail

# ─── Fill these in ────────────────────────────────────────────────────────────
API="http://localhost:8080"                 # API base URL
TOKEN="PASTE_STAFF_JWT_HERE"                 # back-office JWT (admin/branch_admin/manager)
ORG_ID="PASTE_ORG_ID_HERE"                   # organization UUID
# ──────────────────────────────────────────────────────────────────────────────

MODE="${MODE:-itemised}"                     # itemised | buffet
AUTH=(-H "Authorization: Bearer ${TOKEN}")
JSON=(-H "Content-Type: application/json")

say()  { printf '\n\033[1;36m== %s\033[0m\n' "$*"; }
fail() { printf '\033[1;31m!! %s\033[0m\n' "$*" >&2; exit 1; }

[[ "$TOKEN" == PASTE_* ]] && fail "Set TOKEN at the top of the script."
[[ "$ORG_ID" == PASTE_* ]] && fail "Set ORG_ID at the top of the script."
command -v jq >/dev/null || fail "jq is required."

# ─── Step 1: ensure a menu exists ─────────────────────────────────────────────
say "Step 1: upsert menu"
curl -fsS -X PUT "${API}/api/v1/menu" "${AUTH[@]}" "${JSON[@]}" \
  -d '{"name":"Main Menu","is_active":true}' | jq -c '{id, name}'

# ─── Step 2: create menu items, capture their ids ─────────────────────────────
create_item() {  # name price -> prints new item id
  curl -fsS -X POST "${API}/api/v1/menu/items" "${AUTH[@]}" "${JSON[@]}" \
    -d "{\"name\":\"$1\",\"category\":\"$3\",\"price\":$2}" | jq -r '.id'
}
say "Step 2: create menu items"
CHICKEN_ID="$(create_item "Chicken & Rice" 120 "Main")"
VEG_ID="$(create_item "Vegetarian Plate" 100 "Main")"
BUFFET_ID="$(create_item "Lunch Buffet" 250 "Buffet")"
echo "  chicken = ${CHICKEN_ID}"
echo "  veg     = ${VEG_ID}"
echo "  buffet  = ${BUFFET_ID}"

# ─── Step 3: submit the meals request (PUBLIC — no token) ──────────────────────
say "Step 3: submit meals request (${MODE})"
if [[ "$MODE" == "buffet" ]]; then
  BODY=$(cat <<JSON
{
  "company": { "company_name": "ABC Corporation Ltd", "tpin": "1002345678" },
  "booked_by": { "first_name": "Martin", "last_name": "Phiri", "email": "martin@abccorp.com" },
  "from": "2026-07-10", "to": "2026-07-10",
  "headcount": 120,
  "items": [ { "menu_item_id": "${BUFFET_ID}", "quantity": 120 } ]
}
JSON
)
else
  BODY=$(cat <<JSON
{
  "company": { "company_name": "ABC Corporation Ltd", "tpin": "1002345678" },
  "booked_by": { "first_name": "Martin", "last_name": "Phiri", "email": "martin@abccorp.com" },
  "from": "2026-07-10", "to": "2026-07-10",
  "dietary_notes": "One vegetarian",
  "guests": [
    { "first_name": "John",  "last_name": "Banda", "identification_card": "ZM-9981",
      "items": [ { "menu_item_id": "${CHICKEN_ID}", "quantity": 1 } ] },
    { "first_name": "Grace", "last_name": "Phiri", "identification_card": "ZM-9982",
      "items": [ { "menu_item_id": "${VEG_ID}", "quantity": 1 } ] }
  ]
}
JSON
)
fi

REQUEST_ID="$(curl -fsS -X POST \
  "${API}/api/v1/web/bookings/corporate?type=meals&org_id=${ORG_ID}" \
  "${JSON[@]}" -d "${BODY}" | jq -r '.id')"
[[ -n "$REQUEST_ID" && "$REQUEST_ID" != "null" ]] || fail "submit did not return a request id"
echo "  request_id = ${REQUEST_ID}"

# ─── Step 4: view the enriched task (meals_summary with names + prices) ────────
say "Step 4: GET request — meals_summary (the back-office task view)"
curl -fsS "${API}/api/v1/bookings/requests/${REQUEST_ID}" "${AUTH[@]}" \
  | jq '.meals_summary'

# ─── Step 5: approve -> auto-creates booking + orders + draft invoice ──────────
say "Step 5: approve request"
curl -fsS -X PUT "${API}/api/v1/booking-requests/${REQUEST_ID}/approve" "${AUTH[@]}" \
  | jq -c '.'

# ─── Step 6: verify booking / orders / invoice ────────────────────────────────
say "Step 6a: find the materialised booking"
BOOKING_ID="$(curl -fsS "${API}/api/v1/bookings?booking_type=meals&page=1&page_size=20" "${AUTH[@]}" \
  | jq -r --arg rid "$REQUEST_ID" '.data[] | select(.request_id == $rid) | .id' | head -n1)"
[[ -n "$BOOKING_ID" && "$BOOKING_ID" != "null" ]] || fail "no meals booking linked to request ${REQUEST_ID}"
echo "  booking_id = ${BOOKING_ID}"

say "Step 6b: orders for the booking"
curl -fsS "${API}/api/v1/orders?booking_id=${BOOKING_ID}&status=open" "${AUTH[@]}" \
  | jq '.data[]? // .[]? | {order_number, attendee_name, total, items: [.items[]? | {item_name, quantity, unit_price, subtotal}]}'

say "Step 6c: draft invoice for the booking"
curl -fsS "${API}/api/v1/invoices/booking/${BOOKING_ID}" "${AUTH[@]}" \
  | jq '{invoice_number, status, subtotal, tax_amount, total, line_items: [.line_items[]? | {description, quantity, unit_price, total}]}'

say "DONE — meals flow verified end to end."
