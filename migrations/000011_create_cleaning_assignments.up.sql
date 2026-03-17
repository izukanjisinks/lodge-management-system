CREATE TYPE cleaning_status AS ENUM ('pending', 'in_progress', 'completed', 'skipped');

CREATE TABLE IF NOT EXISTS cleaning_assignments (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id        UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    assigned_to    UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    assigned_by    UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    scheduled_date DATE NOT NULL,
    status         cleaning_status NOT NULL DEFAULT 'pending',
    notes          TEXT,
    completed_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cleaning_assignments_room_id ON cleaning_assignments(room_id);
CREATE INDEX idx_cleaning_assignments_assigned_to ON cleaning_assignments(assigned_to);
CREATE INDEX idx_cleaning_assignments_scheduled_date ON cleaning_assignments(scheduled_date);
CREATE INDEX idx_cleaning_assignments_status ON cleaning_assignments(status);
