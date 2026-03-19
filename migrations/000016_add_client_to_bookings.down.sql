DROP INDEX IF EXISTS idx_bookings_client_id;
DROP INDEX IF EXISTS idx_bookings_client_type;

ALTER TABLE bookings
    DROP COLUMN IF EXISTS client_id,
    DROP COLUMN IF EXISTS client_type;
