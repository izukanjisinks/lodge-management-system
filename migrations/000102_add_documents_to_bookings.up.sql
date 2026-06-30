-- Single-phase booking flow: customer submissions now land directly in bookings as
-- pending rows. The old *_booking_requests tables carried a documents TEXT[] column
-- (POPs, authorisation letters, etc.); bring that onto bookings so the documents
-- submitted with a booking are preserved without a separate table.
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS documents TEXT[] NOT NULL DEFAULT '{}';
