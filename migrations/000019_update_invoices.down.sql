DROP INDEX IF EXISTS idx_invoices_invoice_number;

ALTER TABLE invoices
    DROP COLUMN IF EXISTS invoice_number,
    DROP COLUMN IF EXISTS tax_rate,
    DROP COLUMN IF EXISTS paid_date,
    DROP COLUMN IF EXISTS notes;
