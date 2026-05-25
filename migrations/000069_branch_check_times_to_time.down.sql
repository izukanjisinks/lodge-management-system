ALTER TABLE branches
    ALTER COLUMN check_in_time  TYPE VARCHAR(5) USING check_in_time::VARCHAR,
    ALTER COLUMN check_out_time TYPE VARCHAR(5) USING check_out_time::VARCHAR;
