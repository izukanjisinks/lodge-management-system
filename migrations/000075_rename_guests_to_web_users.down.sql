-- Restore guests table and individual_profiles.
-- Note: individual_profiles data cannot be recovered from this rollback alone.

ALTER TABLE web_users RENAME TO guests;

ALTER INDEX idx_web_users_email RENAME TO idx_guests_email;

CREATE TABLE IF NOT EXISTS individual_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    guest_id            UUID REFERENCES guests(id) ON DELETE SET NULL,
    org_id              UUID REFERENCES organizations(id),
    full_name           VARCHAR(255) NOT NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    phone               VARCHAR(50) NOT NULL,
    id_passport_number  VARCHAR(100),
    nationality         VARCHAR(100),
    status              VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
