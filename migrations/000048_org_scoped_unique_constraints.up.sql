-- Resolve duplicate names/emails across orgs before adding scoped unique constraints.
-- Each duplicate gets a numeric suffix (e.g. "Room 101 (2)") to avoid conflicts.

-- ── rooms.name ────────────────────────────────────────────────────────────────
WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY org_id, name ORDER BY created_at ASC) AS rn
    FROM rooms
)
UPDATE rooms r
SET name = r.name || ' (' || d.rn::text || ')'
FROM duplicates d
WHERE r.id = d.id AND d.rn > 1;

ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_name_key;
ALTER TABLE rooms ADD CONSTRAINT uq_rooms_name_org UNIQUE (org_id, name);

-- ── roles.name ────────────────────────────────────────────────────────────────
WITH duplicates AS (
    SELECT role_id,
           ROW_NUMBER() OVER (PARTITION BY org_id, name ORDER BY created_at ASC) AS rn
    FROM roles
)
UPDATE roles r
SET name = r.name || ' (' || d.rn::text || ')'
FROM duplicates d
WHERE r.role_id = d.role_id AND d.rn > 1;

ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_name_key;
ALTER TABLE roles ADD CONSTRAINT uq_roles_name_org UNIQUE (org_id, name);

-- ── individual_profiles.email ─────────────────────────────────────────────────
WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY org_id, email ORDER BY created_at ASC) AS rn
    FROM individual_profiles
    WHERE email IS NOT NULL AND email != ''
)
UPDATE individual_profiles ip
SET email = ip.email || '.dup' || d.rn::text
FROM duplicates d
WHERE ip.id = d.id AND d.rn > 1;

ALTER TABLE individual_profiles DROP CONSTRAINT IF EXISTS individual_profiles_email_key;
ALTER TABLE individual_profiles ADD CONSTRAINT uq_individual_profiles_email_org UNIQUE (org_id, email);

-- ── corporate_profiles.email ──────────────────────────────────────────────────
WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY org_id, email ORDER BY created_at ASC) AS rn
    FROM corporate_profiles
    WHERE email IS NOT NULL AND email != ''
)
UPDATE corporate_profiles cp
SET email = cp.email || '.dup' || d.rn::text
FROM duplicates d
WHERE cp.id = d.id AND d.rn > 1;

ALTER TABLE corporate_profiles DROP CONSTRAINT IF EXISTS corporate_profiles_email_key;
ALTER TABLE corporate_profiles ADD CONSTRAINT uq_corporate_profiles_email_org UNIQUE (org_id, email);

-- ── bookings.booking_number ───────────────────────────────────────────────────
WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY org_id, booking_number ORDER BY created_at ASC) AS rn
    FROM bookings
)
UPDATE bookings b
SET booking_number = b.booking_number || '-' || d.rn::text
FROM duplicates d
WHERE b.id = d.id AND d.rn > 1;

ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_booking_number_key;
ALTER TABLE bookings ADD CONSTRAINT uq_bookings_booking_number_org UNIQUE (org_id, booking_number);

-- ── orders.order_number ───────────────────────────────────────────────────────
WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY org_id, order_number ORDER BY created_at ASC) AS rn
    FROM orders
)
UPDATE orders o
SET order_number = o.order_number || '-' || d.rn::text
FROM duplicates d
WHERE o.id = d.id AND d.rn > 1;

ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_order_number_key;
ALTER TABLE orders ADD CONSTRAINT uq_orders_order_number_org UNIQUE (org_id, order_number);
