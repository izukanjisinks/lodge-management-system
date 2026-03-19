CREATE TABLE IF NOT EXISTS corporate_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name        VARCHAR(255) NOT NULL,
    contact_person      VARCHAR(255) NOT NULL,
    email               VARCHAR(255) UNIQUE NOT NULL,
    phone               VARCHAR(50) NOT NULL,
    company_reg_number  VARCHAR(100) NOT NULL,
    industry            VARCHAR(100),
    status              VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_corporate_profiles_email  ON corporate_profiles(email);
CREATE INDEX idx_corporate_profiles_status ON corporate_profiles(status);
