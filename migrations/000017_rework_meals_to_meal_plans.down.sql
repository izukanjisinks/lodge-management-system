DROP TABLE IF EXISTS booking_meal_plans;
DROP TABLE IF EXISTS meal_plans;

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

CREATE TABLE IF NOT EXISTS booking_meals (
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    meal_id    UUID NOT NULL REFERENCES meals(id) ON DELETE RESTRICT,
    PRIMARY KEY (booking_id, meal_id)
);
