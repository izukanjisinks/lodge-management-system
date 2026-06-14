ALTER TABLE booking_documents
    DROP CONSTRAINT IF EXISTS booking_documents_booking_id_fkey,
    DROP COLUMN IF EXISTS booking_id,
    ADD COLUMN corporate_client_id UUID NOT NULL REFERENCES corporate_profiles(id) ON DELETE CASCADE,
    ADD CONSTRAINT booking_documents_corporate_client_id_org_id_key UNIQUE (corporate_client_id, org_id);
