-- =============================================================================
-- Demo seed data for Lodge Management System
-- Covers: rooms, individual clients, corporate clients, meal plans,
--         bookings (all statuses), booking_meal_plans, invoices, invoice_line_items
--
-- Depends on: migration 000014 (dev users already seeded)
-- Idempotent: all inserts use ON CONFLICT DO NOTHING
-- =============================================================================

DO $$
DECLARE
    -- Users
    v_admin_id        UUID;
    v_manager_id      UUID;
    v_receptionist_id UUID;

    -- Rooms
    v_room_101 UUID;
    v_room_102 UUID;
    v_room_201 UUID;
    v_room_202 UUID;
    v_suite_301 UUID;
    v_cabin_01  UUID;
    v_conf_a    UUID;

    -- Individual clients
    v_client_john   UUID;
    v_client_sarah  UUID;
    v_client_peter  UUID;

    -- Corporate clients
    v_corp_acme   UUID;
    v_corp_zesco  UUID;

    -- Meal plans
    v_meal_bb     UUID;
    v_meal_hb     UUID;
    v_meal_fb     UUID;

    -- Bookings
    v_booking_1 UUID;
    v_booking_2 UUID;
    v_booking_3 UUID;
    v_booking_4 UUID;
    v_booking_5 UUID;

    -- Invoices
    v_invoice_1 UUID;
    v_invoice_2 UUID;

BEGIN

-- ---------------------------------------------------------------------------
-- Resolve user IDs
-- ---------------------------------------------------------------------------
SELECT user_id INTO v_admin_id        FROM users WHERE email = 'admin@lodge.dev'        LIMIT 1;
SELECT user_id INTO v_manager_id      FROM users WHERE email = 'manager@lodge.dev'      LIMIT 1;
SELECT user_id INTO v_receptionist_id FROM users WHERE email = 'receptionist@lodge.dev' LIMIT 1;

IF v_admin_id IS NULL THEN
    RAISE EXCEPTION 'Admin user not found — ensure migration 000014 has run first';
END IF;

-- ---------------------------------------------------------------------------
-- Rooms
-- ---------------------------------------------------------------------------
INSERT INTO rooms (id, name, type, capacity, price_per_night, amenities, is_available, description) VALUES
    (gen_random_uuid(), 'Room 101', 'single', 1, 150.00, ARRAY['WiFi','Air Conditioning','En-suite Bathroom'],                               TRUE,  'Cosy single room on the ground floor with garden view'),
    (gen_random_uuid(), 'Room 102', 'single', 1, 150.00, ARRAY['WiFi','Air Conditioning','En-suite Bathroom'],                               TRUE,  'Bright single room overlooking the courtyard'),
    (gen_random_uuid(), 'Room 201', 'double', 2, 250.00, ARRAY['WiFi','Air Conditioning','En-suite Bathroom','TV'],                          TRUE,  'Spacious double room with king-size bed'),
    (gen_random_uuid(), 'Room 202', 'double', 2, 260.00, ARRAY['WiFi','Air Conditioning','En-suite Bathroom','TV','Mini Bar'],               TRUE,  'Double room with mini bar and pool view'),
    (gen_random_uuid(), 'Suite 301', 'suite', 4, 520.00, ARRAY['WiFi','Air Conditioning','En-suite Bathroom','TV','Mini Bar','Jacuzzi','Lounge Area'], TRUE, 'Luxury suite with separate lounge and jacuzzi'),
    (gen_random_uuid(), 'Cabin 01',  'cabin', 3, 380.00, ARRAY['WiFi','Fireplace','En-suite Bathroom','Kitchenette','Outdoor Deck'],         TRUE,  'Rustic cabin set in the lodge grounds with private deck'),
    (gen_random_uuid(), 'Conference A', 'conference', 20, 800.00, ARRAY['WiFi','Projector','Whiteboard','Air Conditioning','PA System'],     TRUE,  'Main conference room seating up to 20 delegates')
ON CONFLICT (name) DO NOTHING;

SELECT id INTO v_room_101  FROM rooms WHERE name = 'Room 101'      LIMIT 1;
SELECT id INTO v_room_102  FROM rooms WHERE name = 'Room 102'      LIMIT 1;
SELECT id INTO v_room_201  FROM rooms WHERE name = 'Room 201'      LIMIT 1;
SELECT id INTO v_room_202  FROM rooms WHERE name = 'Room 202'      LIMIT 1;
SELECT id INTO v_suite_301 FROM rooms WHERE name = 'Suite 301'     LIMIT 1;
SELECT id INTO v_cabin_01  FROM rooms WHERE name = 'Cabin 01'      LIMIT 1;
SELECT id INTO v_conf_a    FROM rooms WHERE name = 'Conference A'  LIMIT 1;

-- ---------------------------------------------------------------------------
-- Individual clients
-- ---------------------------------------------------------------------------
INSERT INTO individual_profiles (id, full_name, email, phone, id_passport_number, nationality, status, notes) VALUES
    (gen_random_uuid(), 'John Mwewa',    'john.mwewa@gmail.com',    '+260971234001', 'NRC100100/10/1', 'Zambian',   'active', 'Frequent guest, prefers ground floor'),
    (gen_random_uuid(), 'Sarah Banda',   'sarah.banda@gmail.com',   '+260971234002', 'NRC100100/10/2', 'Zambian',   'active', NULL),
    (gen_random_uuid(), 'Peter Okafor',  'peter.okafor@gmail.com',  '+260971234003', 'A12345678',      'Nigerian',  'active', 'Travelling on business')
ON CONFLICT (email) DO NOTHING;

SELECT id INTO v_client_john  FROM individual_profiles WHERE email = 'john.mwewa@gmail.com'  LIMIT 1;
SELECT id INTO v_client_sarah FROM individual_profiles WHERE email = 'sarah.banda@gmail.com'  LIMIT 1;
SELECT id INTO v_client_peter FROM individual_profiles WHERE email = 'peter.okafor@gmail.com' LIMIT 1;

-- ---------------------------------------------------------------------------
-- Corporate clients
-- ---------------------------------------------------------------------------
INSERT INTO corporate_profiles (id, company_name, contact_person, email, phone, company_reg_number, industry, status, notes) VALUES
    (gen_random_uuid(), 'Acme Zambia Ltd',  'Grace Tembo',  'bookings@acme.co.zm',  '+260211100001', 'ZM-2020-1001', 'Mining',   'active', 'Monthly corporate rate negotiated'),
    (gen_random_uuid(), 'ZESCO Limited',    'David Phiri',  'events@zesco.co.zm',   '+260211100002', 'ZM-2018-0042', 'Energy',   'active', 'Conference room bookings only')
ON CONFLICT (email) DO NOTHING;

SELECT id INTO v_corp_acme  FROM corporate_profiles WHERE email = 'bookings@acme.co.zm' LIMIT 1;
SELECT id INTO v_corp_zesco FROM corporate_profiles WHERE email = 'events@zesco.co.zm'  LIMIT 1;

-- ---------------------------------------------------------------------------
-- Meal plans
-- ---------------------------------------------------------------------------
INSERT INTO meal_plans (id, name, price_per_person_per_night, includes, description, is_active) VALUES
    (gen_random_uuid(), 'Bed & Breakfast', 25.00, ARRAY['Breakfast'],                    'Continental breakfast included',              TRUE),
    (gen_random_uuid(), 'Half Board',      55.00, ARRAY['Breakfast','Dinner'],            'Breakfast and dinner included',               TRUE),
    (gen_random_uuid(), 'Full Board',      85.00, ARRAY['Breakfast','Lunch','Dinner'],    'All three meals included throughout your stay', TRUE)
ON CONFLICT DO NOTHING;

SELECT id INTO v_meal_bb FROM meal_plans WHERE name = 'Bed & Breakfast' LIMIT 1;
SELECT id INTO v_meal_hb FROM meal_plans WHERE name = 'Half Board'      LIMIT 1;
SELECT id INTO v_meal_fb FROM meal_plans WHERE name = 'Full Board'      LIMIT 1;

-- ---------------------------------------------------------------------------
-- Bookings
-- Booking 1: John in Room 101 — checked_out (past stay)
-- Booking 2: Sarah in Room 201 — checked_in (current stay)
-- Booking 3: Peter in Suite 301 — confirmed (upcoming)
-- Booking 4: Acme Corp in Conference A — pending
-- Booking 5: ZESCO in Conference A — cancelled
-- ---------------------------------------------------------------------------
v_booking_1 := gen_random_uuid();
v_booking_2 := gen_random_uuid();
v_booking_3 := gen_random_uuid();
v_booking_4 := gen_random_uuid();
v_booking_5 := gen_random_uuid();

INSERT INTO bookings (id, user_id, room_id, client_id, client_type, check_in, check_out, guests, status, special_requests) VALUES
    (v_booking_1, v_receptionist_id, v_room_101,  v_client_john,  'individual', CURRENT_DATE - 7,  CURRENT_DATE - 3,  1, 'checked_out', 'Non-smoking room'),
    (v_booking_2, v_receptionist_id, v_room_201,  v_client_sarah, 'individual', CURRENT_DATE - 1,  CURRENT_DATE + 3,  2, 'checked_in',  'Extra pillows requested'),
    (v_booking_3, v_manager_id,      v_suite_301, v_client_peter, 'individual', CURRENT_DATE + 5,  CURRENT_DATE + 10, 2, 'confirmed',   'Late check-in, arriving after 10 PM'),
    (v_booking_4, v_receptionist_id, v_conf_a,    v_corp_acme,    'corporate',  CURRENT_DATE + 2,  CURRENT_DATE + 3,  15,'pending',     'AV setup required, catering for 15'),
    (v_booking_5, v_manager_id,      v_conf_a,    v_corp_zesco,   'corporate',  CURRENT_DATE - 14, CURRENT_DATE - 12, 10,'cancelled',   NULL)
ON CONFLICT (id) DO NOTHING;

-- Room availability: rooms with confirmed/checked_in bookings are unavailable
UPDATE rooms SET is_available = FALSE WHERE id IN (v_room_201, v_suite_301);

-- ---------------------------------------------------------------------------
-- Booking meal plans
-- Booking 1: John had Full Board
-- Booking 2: Sarah has Half Board
-- Booking 3: Peter has Bed & Breakfast
-- ---------------------------------------------------------------------------
INSERT INTO booking_meal_plans (booking_id, meal_plan_id, guests) VALUES
    (v_booking_1, v_meal_fb, 1),
    (v_booking_2, v_meal_hb, 2),
    (v_booking_3, v_meal_bb, 2)
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------------
-- Invoices — only for checked_out and checked_in bookings (confirmed + past)
-- Invoice 1: Booking 1 (John, checked_out) — paid
-- Invoice 2: Booking 2 (Sarah, checked_in) — issued
-- Invoice 3: Booking 3 (Peter, confirmed)  — draft
-- ---------------------------------------------------------------------------

-- Invoice 1: John — Room 101 (4 nights × 150) + Full Board (4 nights × 1 guest × 85)
INSERT INTO invoices (id, invoice_number, booking_id, subtotal, tax_rate, tax, total, status, issued_at, due_date, paid_date)
VALUES (
    gen_random_uuid(),
    'INV-' || TO_CHAR(NOW(), 'YYYY') || '-0001',
    v_booking_1,
    940.00, 16.00, 150.40, 1090.40,
    'paid', CURRENT_DATE - 3, CURRENT_DATE + 27, CURRENT_DATE - 3
)
ON CONFLICT (booking_id) DO NOTHING;

SELECT id INTO v_invoice_1 FROM invoices WHERE booking_id = v_booking_1 LIMIT 1;

INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, total)
SELECT gen_random_uuid(), v_invoice_1, 'Room 101 — 4 night(s) @ 150.00/night',                      4, 150.00, 600.00
WHERE NOT EXISTS (SELECT 1 FROM invoice_line_items WHERE invoice_id = v_invoice_1);

INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, total)
SELECT gen_random_uuid(), v_invoice_1, 'Full Board — 4 night(s) × 1 guest(s) @ 85.00/person/night', 4,  85.00, 340.00
WHERE NOT EXISTS (SELECT 1 FROM invoice_line_items WHERE invoice_id = v_invoice_1 AND description LIKE 'Full Board%');

-- Invoice 2: Sarah — Room 201 (4 nights × 250) + Half Board (4 nights × 2 guests × 55)
INSERT INTO invoices (id, invoice_number, booking_id, subtotal, tax_rate, tax, total, status, issued_at, due_date)
VALUES (
    gen_random_uuid(),
    'INV-' || TO_CHAR(NOW(), 'YYYY') || '-0002',
    v_booking_2,
    1440.00, 16.00, 230.40, 1670.40,
    'issued', CURRENT_DATE, CURRENT_DATE + 30
)
ON CONFLICT (booking_id) DO NOTHING;

SELECT id INTO v_invoice_2 FROM invoices WHERE booking_id = v_booking_2 LIMIT 1;

INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, total)
SELECT gen_random_uuid(), v_invoice_2, 'Room 201 — 4 night(s) @ 250.00/night',                       4, 250.00, 1000.00
WHERE NOT EXISTS (SELECT 1 FROM invoice_line_items WHERE invoice_id = v_invoice_2);

INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, total)
SELECT gen_random_uuid(), v_invoice_2, 'Half Board — 4 night(s) × 2 guest(s) @ 55.00/person/night',  8,  55.00,  440.00
WHERE NOT EXISTS (SELECT 1 FROM invoice_line_items WHERE invoice_id = v_invoice_2 AND description LIKE 'Half Board%');

-- Invoice 3: Peter — Suite 301 (5 nights × 520) + Bed & Breakfast (5 nights × 2 guests × 25)
INSERT INTO invoices (id, invoice_number, booking_id, subtotal, tax_rate, tax, total, status, issued_at, due_date)
VALUES (
    gen_random_uuid(),
    'INV-' || TO_CHAR(NOW(), 'YYYY') || '-0003',
    v_booking_3,
    2850.00, 16.00, 456.00, 3306.00,
    'draft', CURRENT_DATE, CURRENT_DATE + 30
)
ON CONFLICT (booking_id) DO NOTHING;

INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, total)
SELECT gen_random_uuid(), i.id, 'Suite 301 — 5 night(s) @ 520.00/night', 5, 520.00, 2600.00
FROM invoices i WHERE i.booking_id = v_booking_3
AND NOT EXISTS (SELECT 1 FROM invoice_line_items WHERE invoice_id = i.id);

INSERT INTO invoice_line_items (id, invoice_id, description, quantity, unit_price, total)
SELECT gen_random_uuid(), i.id, 'Bed & Breakfast — 5 night(s) × 2 guest(s) @ 25.00/person/night', 10, 25.00, 250.00
FROM invoices i WHERE i.booking_id = v_booking_3
AND NOT EXISTS (SELECT 1 FROM invoice_line_items WHERE invoice_id = i.id AND description LIKE 'Bed & Breakfast%');

RAISE NOTICE 'Demo seed data inserted successfully';

END $$;
