-- Revert to room-only. Note: this will fail if any event/meals rows exist.
ALTER TABLE individual_booking_requests
    DROP CONSTRAINT IF EXISTS individual_booking_requests_booking_type_check;

ALTER TABLE individual_booking_requests
    ADD CONSTRAINT individual_booking_requests_booking_type_check
    CHECK (booking_type IN ('room'));
