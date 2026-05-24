ALTER TABLE branches
    DROP COLUMN IF EXISTS street_address,
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS country,
    ADD COLUMN address TEXT;
