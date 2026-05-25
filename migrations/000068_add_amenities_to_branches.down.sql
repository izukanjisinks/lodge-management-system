ALTER TABLE branches
    DROP COLUMN IF EXISTS parking,
    DROP COLUMN IF EXISTS restaurant,
    DROP COLUMN IF EXISTS check_in_time,
    DROP COLUMN IF EXISTS check_out_time;
