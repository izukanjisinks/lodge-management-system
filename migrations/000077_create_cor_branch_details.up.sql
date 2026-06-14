CREATE TABLE cor_branch_details (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES cor_company_details(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    address    TEXT,
    phone      VARCHAR(50),
    meta_data  JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cor_branch_details_company_id ON cor_branch_details(company_id);

ALTER TABLE cor_branch_details ADD CONSTRAINT uq_cor_branch_details_company_name UNIQUE (company_id, name);
