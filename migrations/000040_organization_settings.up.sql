CREATE TABLE IF NOT EXISTS organization_settings (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id                UUID        UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    auto_close_orders     BOOLEAN     NOT NULL DEFAULT TRUE,
    auto_extend_checkout  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed the system default row (org_id = NULL means it applies to all orgs unless overridden)
INSERT INTO organization_settings (auto_close_orders, auto_extend_checkout)
VALUES (TRUE, TRUE);
