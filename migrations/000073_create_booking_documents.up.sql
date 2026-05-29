CREATE TABLE booking_documents (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    corporate_client_id UUID NOT NULL REFERENCES corporate_profiles(id) ON DELETE CASCADE,
    org_id              UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    urls                TEXT[] NOT NULL DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (corporate_client_id, org_id)
);
