CREATE TABLE individual_booking_requests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    web_user_id UUID REFERENCES web_users(id) ON DELETE SET NULL,
    booker_name  VARCHAR(255) NOT NULL,
    booker_email VARCHAR(255),
    booker_phone VARCHAR(50),
    booking_type VARCHAR(20) NOT NULL DEFAULT 'room' CHECK (booking_type IN ('room')),
    status       VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    notes        TEXT,
    documents    TEXT[] NOT NULL DEFAULT '{}',
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_individual_booking_requests_org_id      ON individual_booking_requests(org_id);
CREATE INDEX idx_individual_booking_requests_web_user_id ON individual_booking_requests(web_user_id);
CREATE INDEX idx_individual_booking_requests_status      ON individual_booking_requests(status);
CREATE INDEX idx_individual_booking_requests_booking_type ON individual_booking_requests(booking_type);
