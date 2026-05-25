ALTER TABLE organizations
    DROP COLUMN IF EXISTS check_in_time,
    DROP COLUMN IF EXISTS check_out_time;
