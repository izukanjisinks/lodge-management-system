DROP INDEX IF EXISTS idx_orders_attendee_id;
ALTER TABLE orders DROP COLUMN IF EXISTS attendee_id;
