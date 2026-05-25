ALTER TABLE branches
    ALTER COLUMN check_in_time  TYPE TIME USING check_in_time::TIME,
    ALTER COLUMN check_out_time TYPE TIME USING check_out_time::TIME;
