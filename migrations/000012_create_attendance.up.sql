CREATE TYPE attendance_status AS ENUM ('present', 'absent', 'half_day', 'on_leave', 'holiday', 'weekend');
CREATE TYPE attendance_source AS ENUM ('manual', 'system', 'biometric');

CREATE TABLE IF NOT EXISTS attendance (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id    UUID NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
    date           DATE NOT NULL,
    clock_in       TIMESTAMPTZ NULL,
    clock_out      TIMESTAMPTZ NULL,
    total_hours    NUMERIC(5,2) NOT NULL DEFAULT 0,
    status         attendance_status NOT NULL DEFAULT 'absent',
    overtime_hours NUMERIC(5,2) NOT NULL DEFAULT 0,
    notes          TEXT NOT NULL DEFAULT '',
    source         attendance_source NOT NULL DEFAULT 'system',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_attendance_employee_date UNIQUE (employee_id, date)
);

CREATE INDEX idx_attendance_employee ON attendance(employee_id);
CREATE INDEX idx_attendance_date ON attendance(date);
