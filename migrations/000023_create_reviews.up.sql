CREATE TABLE IF NOT EXISTS reviews (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id  UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    guest_id    UUID NOT NULL REFERENCES individual_profiles(id) ON DELETE CASCADE,
    facilities  NUMERIC(2,1) NOT NULL CHECK (facilities  BETWEEN 0 AND 5),
    cleanliness NUMERIC(2,1) NOT NULL CHECK (cleanliness BETWEEN 0 AND 5),
    services    NUMERIC(2,1) NOT NULL CHECK (services    BETWEEN 0 AND 5),
    comfort     NUMERIC(2,1) NOT NULL CHECK (comfort     BETWEEN 0 AND 5),
    location    NUMERIC(2,1) NOT NULL CHECK (location    BETWEEN 0 AND 5),
    comment     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT one_review_per_booking UNIQUE (booking_id)
);

CREATE INDEX idx_reviews_guest_id ON reviews(guest_id);
