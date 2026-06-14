-- booking_documents was dropped in 000082 along with the old bookings table.
-- Recreate it linked to the new bookings table instead of corporate_profiles.
CREATE TABLE booking_documents (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    org_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    urls       TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_booking_documents_booking_id ON booking_documents(booking_id);
CREATE INDEX idx_booking_documents_org_id     ON booking_documents(org_id);
