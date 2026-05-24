UPDATE invoices SET branch_id = NULL
WHERE branch_id IN (SELECT id FROM branches WHERE name LIKE '% (Main Branch)');

ALTER TABLE invoices DROP COLUMN IF EXISTS branch_id;
