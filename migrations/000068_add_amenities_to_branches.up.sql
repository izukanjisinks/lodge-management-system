ALTER TABLE branches
    ADD COLUMN parking       BOOLEAN      NOT NULL DEFAULT FALSE,
    ADD COLUMN restaurant    BOOLEAN      NOT NULL DEFAULT FALSE,
    ADD COLUMN check_in_time VARCHAR(5)   NULL,
    ADD COLUMN check_out_time VARCHAR(5)  NULL;
