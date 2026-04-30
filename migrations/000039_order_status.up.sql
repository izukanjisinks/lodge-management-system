ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS status VARCHAR(10) NOT NULL DEFAULT 'open'
        CHECK (status IN ('open', 'closed'));
