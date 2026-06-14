-- Drop old bookings table and all its dependents cleanly
DROP TABLE IF EXISTS booking_meals CASCADE;
DROP TABLE IF EXISTS booking_documents CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;

CREATE TABLE bookings (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_number   VARCHAR(50) NOT NULL UNIQUE,
    org_id           UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    branch_id        UUID REFERENCES branches(id) ON DELETE SET NULL,

    booking_type     VARCHAR(20) NOT NULL DEFAULT 'room'
                       CHECK (booking_type IN ('room', 'meals', 'conference', 'event')),
    booker_type      VARCHAR(20) NOT NULL DEFAULT 'individual'
                       CHECK (booker_type IN ('individual', 'corporate')),

    -- Denormalised booker identity (survives account deletion)
    booker_name      VARCHAR(255),
    booker_email     VARCHAR(255),
    booker_phone     VARCHAR(50),

    -- Soft links (all nullable, SET NULL on delete)
    web_user_id      UUID REFERENCES web_users(id) ON DELETE SET NULL,
    cor_profile_id   UUID REFERENCES cor_profiles(id) ON DELETE SET NULL,
    company_id       UUID REFERENCES cor_company_details(id) ON DELETE SET NULL,
    request_id       UUID REFERENCES corporate_booking_requests(id) ON DELETE SET NULL,
    venue_id         UUID REFERENCES venues(id) ON DELETE SET NULL,

    total_amount     NUMERIC(10, 2) NOT NULL DEFAULT 0,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending'
                       CHECK (status IN ('pending', 'confirmed', 'checked_in', 'checked_out', 'cancelled')),
    special_requests TEXT,
    overstayed       BOOLEAN NOT NULL DEFAULT FALSE,

    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bookings_org_id        ON bookings(org_id);
CREATE INDEX idx_bookings_branch_id     ON bookings(branch_id);
CREATE INDEX idx_bookings_status        ON bookings(status);
CREATE INDEX idx_bookings_booker_type   ON bookings(booker_type);
CREATE INDEX idx_bookings_web_user_id   ON bookings(web_user_id)    WHERE web_user_id IS NOT NULL;
CREATE INDEX idx_bookings_cor_profile   ON bookings(cor_profile_id) WHERE cor_profile_id IS NOT NULL;
CREATE INDEX idx_bookings_company_id    ON bookings(company_id)     WHERE company_id IS NOT NULL;
CREATE INDEX idx_bookings_request_id    ON bookings(request_id)     WHERE request_id IS NOT NULL;
