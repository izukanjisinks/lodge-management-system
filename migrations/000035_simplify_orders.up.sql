-- Remove order status — orders are append-only records; invoice is updated immediately on placement.
ALTER TABLE orders DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS order_status;
