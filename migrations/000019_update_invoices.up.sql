-- Add cancelled to invoice_status enum
ALTER TYPE invoice_status ADD VALUE IF NOT EXISTS 'cancelled';

-- Add missing columns to invoices table
ALTER TABLE invoices
    ADD COLUMN invoice_number VARCHAR(50)    UNIQUE,
    ADD COLUMN tax_rate       NUMERIC(5, 2)  NOT NULL DEFAULT 0 CHECK (tax_rate >= 0),
    ADD COLUMN paid_date      TIMESTAMPTZ,
    ADD COLUMN notes          TEXT;

-- Backfill invoice_number for any existing rows
WITH numbered AS (
    SELECT id, ROW_NUMBER() OVER (ORDER BY created_at) AS rn, created_at
    FROM invoices
)
UPDATE invoices
SET invoice_number = 'INV-' || TO_CHAR(numbered.created_at, 'YYYY') || '-' || LPAD(CAST(numbered.rn AS TEXT), 4, '0')
FROM numbered
WHERE invoices.id = numbered.id;

-- Now enforce NOT NULL
ALTER TABLE invoices ALTER COLUMN invoice_number SET NOT NULL;

CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
