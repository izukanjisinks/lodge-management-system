ALTER TABLE booking_events
    ADD COLUMN IF NOT EXISTS start_time        TIME,
    ADD COLUMN IF NOT EXISTS end_time          TIME,
    ADD COLUMN IF NOT EXISTS catering_required BOOLEAN NOT NULL DEFAULT FALSE;
