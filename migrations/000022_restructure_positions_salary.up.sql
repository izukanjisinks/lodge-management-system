-- Remove min_salary and max_salary, add base_salary and calculated allowance/tax fields
ALTER TABLE positions
    DROP COLUMN min_salary,
    DROP COLUMN max_salary,
    ADD COLUMN base_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN housing_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN transport_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN medical_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    ADD COLUMN income_tax NUMERIC(15,2) NOT NULL DEFAULT 0;
