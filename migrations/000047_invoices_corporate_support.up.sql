ALTER TABLE invoices
    ALTER COLUMN booking_id DROP NOT NULL,
    ADD COLUMN corporate_client_id UUID REFERENCES corporate_profiles(id) ON DELETE SET NULL;

CREATE INDEX idx_invoices_corporate_client_id ON invoices(corporate_client_id) WHERE corporate_client_id IS NOT NULL;
