ALTER TABLE corporate_booking_requests
    DROP CONSTRAINT IF EXISTS corporate_booking_requests_branch_id_fkey;

ALTER TABLE corporate_booking_requests
    ADD CONSTRAINT corporate_booking_requests_branch_id_fkey
    FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE SET NULL;
