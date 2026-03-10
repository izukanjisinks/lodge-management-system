CREATE TABLE IF NOT EXISTS payslips (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id         UUID NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
    month               INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
    year                INTEGER NOT NULL,
    base_salary         NUMERIC(15,2) NOT NULL DEFAULT 0,
    housing_allowance   NUMERIC(15,2) NOT NULL DEFAULT 0,
    transport_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    medical_allowance   NUMERIC(15,2) NOT NULL DEFAULT 0,
    gross_salary        NUMERIC(15,2) NOT NULL DEFAULT 0,
    income_tax          NUMERIC(15,2) NOT NULL DEFAULT 0,
    leave_days          NUMERIC(15,2) NOT NULL DEFAULT 0,
    net_salary          NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(employee_id, month, year)
);

CREATE INDEX idx_payslips_employee ON payslips(employee_id);
CREATE INDEX idx_payslips_period ON payslips(year, month);
