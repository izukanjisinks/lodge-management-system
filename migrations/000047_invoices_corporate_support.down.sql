DROP INDEX IF EXISTS idx_invoices_corporate_client_id;

ALTER TABLE invoices
    DROP COLUMN IF EXISTS corporate_client_id,
    ALTER COLUMN booking_id SET NOT NULL;
