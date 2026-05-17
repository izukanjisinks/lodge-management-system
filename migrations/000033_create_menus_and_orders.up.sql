-- Phase 1: Introduce menu/order system.
-- meal_plans and booking_meal_plans are kept untouched in this migration
-- so existing data is safe. They are dropped in 000034.

-- ── Menus ────────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS menus (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT  uq_menus_name_org UNIQUE (org_id, name)
);

CREATE INDEX idx_menus_org_id   ON menus(org_id);
CREATE INDEX idx_menus_is_active ON menus(is_active);

-- ── Menu Items ────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS menu_items (
    id           UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_id      UUID           NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    org_id       UUID           NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name         VARCHAR(255)   NOT NULL,
    description  TEXT,
    price        NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    is_available BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    CONSTRAINT   uq_menu_items_name_menu UNIQUE (menu_id, name)
);

CREATE INDEX idx_menu_items_menu_id      ON menu_items(menu_id);
CREATE INDEX idx_menu_items_org_id       ON menu_items(org_id);
CREATE INDEX idx_menu_items_is_available ON menu_items(is_available);

-- ── Order number sequence (per org would be ideal but global is simpler) ──────

CREATE SEQUENCE IF NOT EXISTS order_number_seq START 1000;

-- ── Orders ───────────────────────────────────────────────────────────────────

CREATE TYPE order_type   AS ENUM ('in_house', 'walk_in');
CREATE TYPE order_status AS ENUM ('open', 'closed', 'voided');

CREATE TABLE IF NOT EXISTS orders (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    booking_id   UUID         REFERENCES bookings(id) ON DELETE SET NULL,  -- nullable for walk-ins
    order_number VARCHAR(50)  NOT NULL UNIQUE DEFAULT 'ORD-' || LPAD(CAST(nextval('order_number_seq') AS TEXT), 6, '0'),
    type         order_type   NOT NULL DEFAULT 'in_house',
    status       order_status NOT NULL DEFAULT 'open',
    notes        TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_org_id     ON orders(org_id);
CREATE INDEX idx_orders_booking_id ON orders(booking_id);
CREATE INDEX idx_orders_status     ON orders(status);
CREATE INDEX idx_orders_type       ON orders(type);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- ── Order Items ───────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS order_items (
    id             UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id       UUID           NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id   UUID           NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
    quantity       INT            NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price     NUMERIC(10, 2) NOT NULL CHECK (unit_price >= 0),  -- snapshotted at order time
    subtotal       NUMERIC(10, 2) NOT NULL CHECK (subtotal >= 0),
    notes          TEXT,                                               -- e.g. "no onions"
    created_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_items_order_id     ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item_id ON order_items(menu_item_id);

-- ── Link invoice line items back to their source order ────────────────────────

ALTER TABLE invoice_line_items
    ADD COLUMN IF NOT EXISTS order_id UUID REFERENCES orders(id) ON DELETE SET NULL;

CREATE INDEX idx_invoice_line_items_order_id ON invoice_line_items(order_id);
