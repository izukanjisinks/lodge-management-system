DROP INDEX IF EXISTS idx_positions_role;

ALTER TABLE positions
DROP COLUMN role_id;
