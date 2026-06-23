ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS scheduled_for DATE,
    ADD COLUMN IF NOT EXISTS meal_period   VARCHAR(20);

CREATE INDEX idx_orders_scheduled_for ON orders(scheduled_for) WHERE scheduled_for IS NOT NULL;
