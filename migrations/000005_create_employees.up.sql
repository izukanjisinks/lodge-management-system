CREATE TYPE employment_type AS ENUM ('full_time', 'part_time', 'contract', 'intern');
CREATE TYPE employment_status AS ENUM ('active', 'on_leave', 'suspended', 'terminated', 'resigned');

CREATE TABLE IF NOT EXISTS employees (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
    employee_number     VARCHAR(30) UNIQUE NOT NULL,
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,
    email               VARCHAR(255) NOT NULL,
    personal_email      VARCHAR(255) NOT NULL DEFAULT '',
    phone               VARCHAR(50) NOT NULL DEFAULT '',
    date_of_birth       DATE NULL,
    gender              VARCHAR(20) NOT NULL DEFAULT '',
    national_id         VARCHAR(100) NOT NULL DEFAULT '',
    marital_status      VARCHAR(30) NOT NULL DEFAULT '',
    address             TEXT NOT NULL DEFAULT '',
    city                VARCHAR(100) NOT NULL DEFAULT '',
    state               VARCHAR(100) NOT NULL DEFAULT '',
    country             VARCHAR(100) NOT NULL DEFAULT '',
    department_id       UUID NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    position_id         UUID NOT NULL REFERENCES positions(id) ON DELETE RESTRICT,
    manager_id          UUID NULL REFERENCES employees(id) ON DELETE SET NULL,
    hire_date           DATE NOT NULL,
    probation_end_date  DATE NULL,
    employment_type     employment_type NOT NULL DEFAULT 'full_time',
    employment_status   employment_status NOT NULL DEFAULT 'active',
    termination_date    DATE NULL,
    termination_reason  TEXT NOT NULL DEFAULT '',
    profile_photo_url   TEXT NOT NULL DEFAULT '',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ NULL
);

-- Unique email per active employee (excluding soft-deleted)
CREATE UNIQUE INDEX idx_employees_email_active ON employees(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_employees_number ON employees(employee_number);
CREATE INDEX idx_employees_department ON employees(department_id);
CREATE INDEX idx_employees_position ON employees(position_id);
CREATE INDEX idx_employees_manager ON employees(manager_id);
CREATE INDEX idx_employees_status ON employees(employment_status);
CREATE INDEX idx_employees_deleted ON employees(deleted_at);

-- Add FK from departments.manager_id to employees
ALTER TABLE departments
    ADD CONSTRAINT fk_departments_manager
    FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL;
