-- Add back organization_id column
ALTER TABLE password_policies ADD COLUMN organization_id UUID;

-- Recreate the organization-specific unique index
CREATE UNIQUE INDEX idx_password_policies_organization
    ON password_policies(organization_id)
    WHERE organization_id IS NOT NULL;
