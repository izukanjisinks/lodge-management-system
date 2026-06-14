CREATE TABLE cor_profiles (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES cor_company_details(id) ON DELETE CASCADE,
    branch_id  UUID REFERENCES cor_branch_details(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100) NOT NULL,
    email      VARCHAR(255),
    phone      VARCHAR(50),
    job_title  VARCHAR(100),
    department VARCHAR(100),
    status     VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    meta_data  JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_cor_profiles_email_org UNIQUE (org_id, email)
);

CREATE INDEX idx_cor_profiles_org_id     ON cor_profiles(org_id);
CREATE INDEX idx_cor_profiles_company_id ON cor_profiles(company_id);
CREATE INDEX idx_cor_profiles_branch_id  ON cor_profiles(branch_id);
CREATE INDEX idx_cor_profiles_status     ON cor_profiles(status);
