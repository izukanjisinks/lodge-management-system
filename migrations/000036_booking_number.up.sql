CREATE SEQUENCE IF NOT EXISTS booking_number_seq START 1000;

ALTER TABLE bookings
    ADD COLUMN IF NOT EXISTS booking_number VARCHAR(50)
        NOT NULL UNIQUE DEFAULT 'BK-' || LPAD(CAST(nextval('booking_number_seq') AS TEXT), 6, '0');
