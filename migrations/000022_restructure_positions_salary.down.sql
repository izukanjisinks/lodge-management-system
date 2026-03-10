-- Revert: restore min_salary and max_salary, remove new salary fields
ALTER TABLE positions
    DROP COLUMN base_salary,
    DROP COLUMN housing_allowance,
    DROP COLUMN transport_allowance,
    DROP COLUMN medical_allowance,
    DROP COLUMN income_tax,
    ADD COLUMN min_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN max_salary NUMERIC(15,2) NOT NULL DEFAULT 0;
