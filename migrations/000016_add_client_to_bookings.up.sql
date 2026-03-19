ALTER TABLE bookings
    ADD COLUMN client_id   UUID        NOT NULL DEFAULT gen_random_uuid(),
    ADD COLUMN client_type VARCHAR(20) NOT NULL DEFAULT 'individual'
        CHECK (client_type IN ('individual', 'corporate'));

-- Remove the defaults after adding (they were only needed to satisfy NOT NULL on existing rows)
ALTER TABLE bookings
    ALTER COLUMN client_id   DROP DEFAULT,
    ALTER COLUMN client_type DROP DEFAULT;

CREATE INDEX idx_bookings_client_id   ON bookings(client_id);
CREATE INDEX idx_bookings_client_type ON bookings(client_type);
