ALTER TABLE invoice_line_items
    ADD COLUMN order_item_id UUID REFERENCES order_items(id) ON DELETE SET NULL;
