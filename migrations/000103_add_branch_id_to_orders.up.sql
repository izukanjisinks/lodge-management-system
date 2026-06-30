ALTER TABLE orders
    ADD COLUMN branch_id UUID REFERENCES branches(id) ON DELETE SET NULL;

CREATE INDEX idx_orders_branch_id ON orders(branch_id);
