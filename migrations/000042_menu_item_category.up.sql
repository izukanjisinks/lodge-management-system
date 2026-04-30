ALTER TABLE menu_items
    ADD COLUMN IF NOT EXISTS category VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_menu_items_category ON menu_items(category);
