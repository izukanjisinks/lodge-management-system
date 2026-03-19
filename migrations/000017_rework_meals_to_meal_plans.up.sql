-- Drop old meals table (was: id, name, type meal_type, price, description, is_available)
-- and replace with meal_plans matching the frontend MealPlan model.
-- booking_meals references meals(id) via ON DELETE RESTRICT, so drop it first then recreate.

DROP TABLE IF EXISTS booking_meals;
DROP TABLE IF EXISTS meals;
DROP TYPE IF EXISTS meal_type;

CREATE TABLE IF NOT EXISTS meal_plans (
    id                       UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    name                     VARCHAR(255)   NOT NULL,
    price_per_person_per_night NUMERIC(10,2) NOT NULL CHECK (price_per_person_per_night >= 0),
    includes                 TEXT[]         NOT NULL DEFAULT '{}',
    description              TEXT,
    is_active                BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at               TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_meal_plans_is_active ON meal_plans(is_active);

-- Recreate join table referencing meal_plans
CREATE TABLE IF NOT EXISTS booking_meal_plans (
    booking_id   UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    meal_plan_id UUID NOT NULL REFERENCES meal_plans(id) ON DELETE RESTRICT,
    guests       INT  NOT NULL DEFAULT 1 CHECK (guests > 0),
    PRIMARY KEY (booking_id, meal_plan_id)
);

CREATE INDEX idx_booking_meal_plans_booking_id   ON booking_meal_plans(booking_id);
CREATE INDEX idx_booking_meal_plans_meal_plan_id ON booking_meal_plans(meal_plan_id);
