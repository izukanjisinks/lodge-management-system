DROP INDEX IF EXISTS idx_orders_scheduled_for;

ALTER TABLE orders
    DROP COLUMN IF EXISTS scheduled_for,
    DROP COLUMN IF EXISTS meal_period;
