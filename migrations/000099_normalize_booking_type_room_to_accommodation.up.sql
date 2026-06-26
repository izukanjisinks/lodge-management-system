-- Normalize booking_type 'room' -> 'accommodation' to match the canonical
-- model values (accommodation, meals, conference, event). The existing CHECK
-- constraints still reference 'room', so they must be dropped before the data
-- update and recreated with the new allowed set afterwards.

-- individual_booking_requests
ALTER TABLE individual_booking_requests DROP CONSTRAINT IF EXISTS individual_booking_requests_booking_type_check;
UPDATE individual_booking_requests SET booking_type = 'accommodation' WHERE booking_type = 'room';
ALTER TABLE individual_booking_requests ADD CONSTRAINT individual_booking_requests_booking_type_check
    CHECK (booking_type IN ('accommodation', 'meals', 'conference', 'event'));

-- bookings
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_booking_type_check;
UPDATE bookings SET booking_type = 'accommodation' WHERE booking_type = 'room';
ALTER TABLE bookings ADD CONSTRAINT bookings_booking_type_check
    CHECK (booking_type IN ('accommodation', 'meals', 'conference', 'event'));
