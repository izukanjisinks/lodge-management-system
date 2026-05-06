DROP INDEX IF EXISTS idx_bookings_corporate_client_id;

ALTER TABLE bookings DROP COLUMN IF EXISTS corporate_client_id;
