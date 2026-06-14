CREATE TABLE booking_attendees (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id          UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    corporate_guest_id  UUID REFERENCES corporate_guests(id) ON DELETE SET NULL,

    full_name           VARCHAR(255) NOT NULL,
    email               VARCHAR(255),
    phone               VARCHAR(50),
    identification_card VARCHAR(100),
    dietary_notes       TEXT,
    special_needs       TEXT,
    is_lead_contact     BOOLEAN NOT NULL DEFAULT FALSE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_booking_attendees_booking_id ON booking_attendees(booking_id);
CREATE INDEX idx_booking_attendees_corporate_guest ON booking_attendees(corporate_guest_id)
    WHERE corporate_guest_id IS NOT NULL;
