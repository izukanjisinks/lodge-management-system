CREATE TABLE IF NOT EXISTS leave_balances (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id      UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id    UUID NOT NULL REFERENCES leave_types(id) ON DELETE RESTRICT,
    year             INT NOT NULL,
    total_entitled   INT NOT NULL DEFAULT 0,
    used             INT NOT NULL DEFAULT 0,
    pending          INT NOT NULL DEFAULT 0,
    carried_forward  INT NOT NULL DEFAULT 0,
    adjustment       INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_leave_balance_employee_type_year UNIQUE (employee_id, leave_type_id, year)
);

CREATE INDEX idx_leave_balances_employee ON leave_balances(employee_id);
CREATE INDEX idx_leave_balances_year ON leave_balances(year);
