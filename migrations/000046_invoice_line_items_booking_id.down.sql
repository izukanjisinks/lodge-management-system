DROP INDEX IF EXISTS idx_invoice_line_items_booking_id;

ALTER TABLE invoice_line_items DROP COLUMN IF EXISTS booking_id;
