CREATE TABLE IF NOT EXISTS individual_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name           VARCHAR(255) NOT NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    phone               VARCHAR(50) NOT NULL,
    id_passport_number  VARCHAR(100) NOT NULL,
    nationality         VARCHAR(100),
    status              VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_individual_profiles_email  ON individual_profiles(email);
CREATE INDEX idx_individual_profiles_status ON individual_profiles(status);
