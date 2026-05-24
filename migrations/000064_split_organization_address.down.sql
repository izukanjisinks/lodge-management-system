ALTER TABLE organizations
    ADD COLUMN address TEXT;

UPDATE organizations SET address = street_address WHERE street_address IS NOT NULL;

ALTER TABLE organizations
    DROP COLUMN IF EXISTS street_address,
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS country;
