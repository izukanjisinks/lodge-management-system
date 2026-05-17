ALTER TABLE invoice_line_items DROP COLUMN IF EXISTS order_id;

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS menus;

DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS order_type;

DROP SEQUENCE IF EXISTS order_number_seq;
