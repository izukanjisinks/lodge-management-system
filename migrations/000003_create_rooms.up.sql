CREATE TYPE room_type AS ENUM ('single', 'double', 'suite', 'cabin', 'conference');

CREATE TABLE IF NOT EXISTS rooms (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(100) UNIQUE NOT NULL,
    type           room_type NOT NULL,
    capacity       INT NOT NULL CHECK (capacity > 0),
    price_per_night NUMERIC(10, 2) NOT NULL CHECK (price_per_night >= 0),
    amenities      TEXT[] NOT NULL DEFAULT '{}',
    is_available   BOOLEAN NOT NULL DEFAULT TRUE,
    description    TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rooms_type ON rooms(type);
CREATE INDEX idx_rooms_is_available ON rooms(is_available);
