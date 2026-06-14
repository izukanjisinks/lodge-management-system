CREATE TABLE booking_room_assignments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id   UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    room_id      UUID NOT NULL REFERENCES rooms(id) ON DELETE RESTRICT,
    attendee_id  UUID REFERENCES booking_attendees(id) ON DELETE SET NULL,

    check_in     DATE NOT NULL,
    check_out    DATE NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending'
                   CHECK (status IN ('pending', 'confirmed', 'checked_in', 'checked_out', 'cancelled')),

    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_room_assignment_dates CHECK (check_out > check_in)
);

CREATE INDEX idx_room_assignments_booking_id  ON booking_room_assignments(booking_id);
CREATE INDEX idx_room_assignments_room_id     ON booking_room_assignments(room_id);
CREATE INDEX idx_room_assignments_attendee_id ON booking_room_assignments(attendee_id)
    WHERE attendee_id IS NOT NULL;
CREATE INDEX idx_room_assignments_status      ON booking_room_assignments(status);
CREATE INDEX idx_room_assignments_check_in    ON booking_room_assignments(check_in);
CREATE INDEX idx_room_assignments_check_out   ON booking_room_assignments(check_out);
