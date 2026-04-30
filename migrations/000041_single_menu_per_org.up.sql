-- Drop the old name+org uniqueness constraint and replace with one-menu-per-org
ALTER TABLE menus DROP CONSTRAINT IF EXISTS uq_menus_name_org;
ALTER TABLE menus DROP CONSTRAINT IF EXISTS menus_org_id_fkey;

-- Allow org_id to be NULL so the system default row has no org
ALTER TABLE menus ALTER COLUMN org_id DROP NOT NULL;

-- One menu per org (NULL org_id = system default, allowed once via partial index)
ALTER TABLE menus ADD CONSTRAINT uq_menus_org UNIQUE (org_id);

-- Re-add the foreign key as nullable
ALTER TABLE menus
    ADD CONSTRAINT menus_org_id_fkey
    FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- Seed the system default menu
INSERT INTO menus (name, description, is_active)
VALUES ('Default Menu', 'System default menu', TRUE);
