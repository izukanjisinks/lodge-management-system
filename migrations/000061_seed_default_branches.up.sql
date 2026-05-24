-- Create one default branch per existing organization.
-- branch_code is derived from the org name: uppercase, spaces→underscores, max 20 chars.
INSERT INTO branches (id, org_id, name, branch_code, address, location, created_at, updated_at)
SELECT
    gen_random_uuid(),
    o.id,
    o.name || ' (Main Branch)',
    UPPER(SUBSTRING(REGEXP_REPLACE(o.name, '\s+', '_', 'g'), 1, 20)),
    o.address,
    NULL,
    NOW(),
    NOW()
FROM organizations o
WHERE NOT EXISTS (
    SELECT 1 FROM branches b WHERE b.org_id = o.id
);

-- Assign existing rooms with no branch to their org's default branch.
UPDATE rooms r
SET branch_id = b.id
FROM branches b
WHERE b.org_id = r.org_id
  AND r.branch_id IS NULL;

-- Assign existing bookings with no branch to their org's default branch.
UPDATE bookings bk
SET branch_id = b.id
FROM branches b
WHERE b.org_id = bk.org_id
  AND bk.branch_id IS NULL;

-- Assign existing users with no branch to their org's default branch.
UPDATE users u
SET branch_id = b.id
FROM branches b
WHERE b.org_id = u.org_id
  AND u.branch_id IS NULL;

