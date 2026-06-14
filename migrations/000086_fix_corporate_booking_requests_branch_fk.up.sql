-- branch_id was referencing branches(id) (lodge branches) but should reference
-- cor_branch_details(id) (corporate client branches from ResolveChain).
ALTER TABLE corporate_booking_requests
    DROP CONSTRAINT IF EXISTS corporate_booking_requests_branch_id_fkey;

ALTER TABLE corporate_booking_requests
    ADD CONSTRAINT corporate_booking_requests_branch_id_fkey
    FOREIGN KEY (branch_id) REFERENCES cor_branch_details(id) ON DELETE SET NULL;
