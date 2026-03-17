CREATE TABLE IF NOT EXISTS corporate_profiles (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    company_name   VARCHAR(255) NOT NULL,
    tax_id         VARCHAR(100),
    contact_person VARCHAR(255) NOT NULL,
    phone          VARCHAR(50),
    cost_centre    VARCHAR(100),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_corporate_profiles_user_id ON corporate_profiles(user_id);
