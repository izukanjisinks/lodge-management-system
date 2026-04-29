ALTER TABLE bookings DROP COLUMN IF EXISTS booking_number;

DROP SEQUENCE IF EXISTS booking_number_seq;
