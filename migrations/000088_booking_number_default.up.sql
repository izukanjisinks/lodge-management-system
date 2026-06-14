ALTER TABLE bookings
    ALTER COLUMN booking_number SET DEFAULT 'BK-' || LPAD(CAST(nextval('booking_number_seq') AS TEXT), 6, '0');
