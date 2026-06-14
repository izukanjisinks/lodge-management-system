CREATE TABLE cor_company_details (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    tpin         VARCHAR(100),
    reg_number   VARCHAR(100),
    industry     VARCHAR(100),
    country      VARCHAR(100),
    status       VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    meta_data    JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_cor_company_details_reg_number_org UNIQUE (org_id, reg_number, tpin)
);

CREATE INDEX idx_cor_company_details_org_id     ON cor_company_details(org_id);
CREATE INDEX idx_cor_company_details_status     ON cor_company_details(status);
CREATE INDEX idx_cor_company_details_reg_number ON cor_company_details(reg_number);
