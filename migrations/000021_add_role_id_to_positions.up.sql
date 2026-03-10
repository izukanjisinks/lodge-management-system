ALTER TABLE positions
ADD COLUMN role_id UUID NULL REFERENCES roles(role_id) ON DELETE SET NULL;

CREATE INDEX idx_positions_role ON positions(role_id);
