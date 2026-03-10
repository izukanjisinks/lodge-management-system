-- Remove organization_id from password_policies table
ALTER TABLE password_policies DROP COLUMN IF EXISTS organization_id;

-- Drop the organization-specific unique index
DROP INDEX IF EXISTS idx_password_policies_organization;

-- The global unique index ensures only one policy exists
-- No need to modify it since it already ensures single policy
