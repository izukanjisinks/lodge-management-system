CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'dinner');

CREATE TABLE IF NOT EXISTS meals (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    type         meal_type NOT NULL,
    price        NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    description  TEXT,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_meals_type ON meals(type);
CREATE INDEX idx_meals_is_available ON meals(is_available);
