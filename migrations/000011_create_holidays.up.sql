CREATE TABLE IF NOT EXISTS holidays (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    date         DATE NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    location     VARCHAR(255) NOT NULL DEFAULT '',
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_holidays_date ON holidays(date);
CREATE INDEX idx_holidays_location ON holidays(location);
