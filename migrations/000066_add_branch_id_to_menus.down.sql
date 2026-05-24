ALTER TABLE menus DROP CONSTRAINT IF EXISTS uq_menus_org_branch;
ALTER TABLE menus ADD CONSTRAINT uq_menus_name_org UNIQUE (org_id, name);

ALTER TABLE menu_items DROP COLUMN IF EXISTS branch_id;
ALTER TABLE menus DROP COLUMN IF EXISTS branch_id;
