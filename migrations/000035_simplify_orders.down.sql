CREATE TYPE order_status AS ENUM ('open', 'closed', 'voided');

ALTER TABLE orders ADD COLUMN IF NOT EXISTS status order_status NOT NULL DEFAULT 'open';
