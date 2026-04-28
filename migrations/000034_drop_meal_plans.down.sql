-- Restore meal_plans and booking_meal_plans if rolling back
CREATE TABLE IF NOT EXISTS meal_plans (
    id                         UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id                     UUID           REFERENCES organizations(id),
    name                       VARCHAR(255)   NOT NULL,
    price_per_person_per_night NUMERIC(10, 2) NOT NULL CHECK (price_per_person_per_night >= 0),
    includes                   TEXT[]         NOT NULL DEFAULT '{}',
    description                TEXT,
    is_active                  BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at                 TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS booking_meal_plans (
    booking_id   UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    meal_plan_id UUID NOT NULL REFERENCES meal_plans(id) ON DELETE RESTRICT,
    guests       INT  NOT NULL DEFAULT 1 CHECK (guests > 0),
    PRIMARY KEY (booking_id, meal_plan_id)
);
