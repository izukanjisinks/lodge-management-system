CREATE TABLE IF NOT EXISTS positions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title         VARCHAR(255) NOT NULL,
    code          VARCHAR(50) UNIQUE NOT NULL,
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    grade_level   VARCHAR(50) NOT NULL DEFAULT '',
    min_salary    NUMERIC(15,2) NOT NULL DEFAULT 0,
    max_salary    NUMERIC(15,2) NOT NULL DEFAULT 0,
    description   TEXT NOT NULL DEFAULT '',
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ NULL
);

CREATE INDEX idx_positions_code ON positions(code);
CREATE INDEX idx_positions_department ON positions(department_id);
CREATE INDEX idx_positions_deleted ON positions(deleted_at);
