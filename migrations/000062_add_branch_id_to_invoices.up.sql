ALTER TABLE invoices
    ADD COLUMN branch_id UUID REFERENCES branches(id) ON DELETE SET NULL;

-- Assign existing invoices with no branch to their org's default branch.
UPDATE invoices i
SET branch_id = b.id
FROM branches b
WHERE b.org_id = i.org_id
  AND i.branch_id IS NULL;
