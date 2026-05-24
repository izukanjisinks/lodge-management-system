ALTER TABLE branches
    DROP COLUMN IF EXISTS address,
    ADD COLUMN street_address TEXT,
    ADD COLUMN city           TEXT,
    ADD COLUMN country        TEXT;
