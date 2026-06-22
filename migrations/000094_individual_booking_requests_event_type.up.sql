-- Allow standalone event (and forward-looking meals) requests in the individual
-- booking pipeline. Originally the CHECK only permitted 'room'.
ALTER TABLE individual_booking_requests
    DROP CONSTRAINT IF EXISTS individual_booking_requests_booking_type_check;

ALTER TABLE individual_booking_requests
    ADD CONSTRAINT individual_booking_requests_booking_type_check
    CHECK (booking_type IN ('room', 'event', 'meals'));
