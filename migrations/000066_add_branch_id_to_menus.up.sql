ALTER TABLE menus
    ADD COLUMN branch_id UUID REFERENCES branches(id) ON DELETE SET NULL;

ALTER TABLE menu_items
    ADD COLUMN branch_id UUID REFERENCES branches(id) ON DELETE SET NULL;

-- Drop the org-level unique constraint on menus so each branch can have its own menu.
ALTER TABLE menus DROP CONSTRAINT IF EXISTS uq_menus_name_org;
ALTER TABLE menus ADD CONSTRAINT uq_menus_org_branch UNIQUE (org_id, branch_id);
