CREATE TABLE IF NOT EXISTS emergency_contacts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id  UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    relationship VARCHAR(100) NOT NULL,
    phone        VARCHAR(50) NOT NULL,
    email        VARCHAR(255) NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_emergency_contacts_employee ON emergency_contacts(employee_id);
