CREATE TABLE IF NOT EXISTS leave_types (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                    VARCHAR(100) NOT NULL,
    code                    VARCHAR(10) UNIQUE NOT NULL,
    description             TEXT NOT NULL DEFAULT '',
    default_days_per_year   INT NOT NULL DEFAULT 0,
    is_paid                 BOOLEAN NOT NULL DEFAULT TRUE,
    is_carry_forward_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    max_carry_forward_days  INT NOT NULL DEFAULT 0,
    requires_approval       BOOLEAN NOT NULL DEFAULT TRUE,
    requires_document       BOOLEAN NOT NULL DEFAULT FALSE,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed default leave types
INSERT INTO leave_types (code, name, default_days_per_year, is_paid, is_carry_forward_allowed, max_carry_forward_days, requires_approval, requires_document) VALUES
    ('AL', 'Annual Leave',       21, TRUE,  TRUE,  5, TRUE, FALSE),
    ('SL', 'Sick Leave',         15, TRUE,  FALSE, 0, TRUE, TRUE),
    ('PL', 'Parental Leave',     90, TRUE,  FALSE, 0, TRUE, FALSE),
    ('UL', 'Unpaid Leave',        0, FALSE, FALSE, 0, TRUE, FALSE),
    ('CL', 'Compassionate Leave', 5, TRUE,  FALSE, 0, TRUE, FALSE),
    ('ML', 'Marriage Leave',      3, TRUE,  FALSE, 0, TRUE, FALSE)
ON CONFLICT (code) DO NOTHING;

CREATE INDEX idx_leave_types_code ON leave_types(code);
