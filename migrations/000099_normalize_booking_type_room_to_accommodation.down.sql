-- Revert booking_type 'accommodation' -> 'room' and restore the original
-- CHECK constraints that referenced 'room'.

-- individual_booking_requests
ALTER TABLE individual_booking_requests DROP CONSTRAINT IF EXISTS individual_booking_requests_booking_type_check;
UPDATE individual_booking_requests SET booking_type = 'room' WHERE booking_type = 'accommodation';
ALTER TABLE individual_booking_requests ADD CONSTRAINT individual_booking_requests_booking_type_check
    CHECK (booking_type IN ('room', 'event', 'meals'));

-- bookings
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_booking_type_check;
UPDATE bookings SET booking_type = 'room' WHERE booking_type = 'accommodation';
ALTER TABLE bookings ADD CONSTRAINT bookings_booking_type_check
    CHECK (booking_type IN ('room', 'meals', 'conference', 'event'));
