CREATE TABLE IF NOT EXISTS departments (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                 VARCHAR(255) NOT NULL,
    code                 VARCHAR(50) UNIQUE NOT NULL,
    description          TEXT NOT NULL DEFAULT '',
    parent_department_id UUID NULL REFERENCES departments(id) ON DELETE SET NULL,
    manager_id           UUID NULL,
    is_active            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ NULL
);

CREATE INDEX idx_departments_code ON departments(code);
CREATE INDEX idx_departments_parent ON departments(parent_department_id);
CREATE INDEX idx_departments_deleted ON departments(deleted_at);
