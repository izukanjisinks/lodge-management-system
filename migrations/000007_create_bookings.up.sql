CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'checked_in', 'checked_out', 'cancelled');

CREATE TABLE IF NOT EXISTS bookings (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    room_id          UUID NOT NULL REFERENCES rooms(id) ON DELETE RESTRICT,
    check_in         DATE NOT NULL,
    check_out        DATE NOT NULL,
    guests           INT NOT NULL CHECK (guests > 0),
    status           booking_status NOT NULL DEFAULT 'pending',
    special_requests TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_booking_dates CHECK (check_out > check_in)
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_room_id ON bookings(room_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_check_in ON bookings(check_in);
CREATE INDEX idx_bookings_check_out ON bookings(check_out);
