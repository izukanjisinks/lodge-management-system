CREATE TABLE corporate_booking_requests (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id               UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    branch_id            UUID REFERENCES branches(id) ON DELETE SET NULL,
    cor_profile_id       UUID REFERENCES cor_profiles(id) ON DELETE SET NULL,
    company_id           UUID REFERENCES cor_company_details(id) ON DELETE SET NULL,
    booking_type         VARCHAR(20) NOT NULL CHECK (booking_type IN ('accommodation', 'meals', 'conference', 'event')),
    status               VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    reason_for_booking   TEXT,
    notes                TEXT,
    authoriser_name      VARCHAR(255),
    authoriser_email     VARCHAR(255),
    authoriser_phone     VARCHAR(50),
    authoriser_title     VARCHAR(100),
    authoriser_department VARCHAR(100),
    authoriser_gl_code   VARCHAR(50),
    documents            TEXT[] NOT NULL DEFAULT '{}',
    payload              JSONB,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_corporate_booking_requests_org_id       ON corporate_booking_requests(org_id);
CREATE INDEX idx_corporate_booking_requests_branch_id    ON corporate_booking_requests(branch_id);
CREATE INDEX idx_corporate_booking_requests_cor_profile  ON corporate_booking_requests(cor_profile_id);
CREATE INDEX idx_corporate_booking_requests_company_id   ON corporate_booking_requests(company_id);
CREATE INDEX idx_corporate_booking_requests_status       ON corporate_booking_requests(status);
CREATE INDEX idx_corporate_booking_requests_booking_type ON corporate_booking_requests(booking_type);




