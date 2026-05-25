ALTER TABLE organizations
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS parking,
    DROP COLUMN IF EXISTS restaurant;
