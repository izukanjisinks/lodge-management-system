ALTER TABLE organizations
    ADD COLUMN street_address TEXT,
    ADD COLUMN city           TEXT,
    ADD COLUMN country        TEXT;

UPDATE organizations SET street_address = address WHERE address IS NOT NULL;

ALTER TABLE organizations DROP COLUMN address;
