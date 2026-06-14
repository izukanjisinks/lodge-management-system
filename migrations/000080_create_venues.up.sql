CREATE TABLE venues (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    branch_id    UUID REFERENCES branches(id) ON DELETE SET NULL,
    name         VARCHAR(255) NOT NULL,
    venue_type   VARCHAR(50) NOT NULL CHECK (venue_type IN ('conference_hall', 'event_space', 'boardroom', 'outdoor', 'dining')),
    capacity     INT NOT NULL CHECK (capacity > 0),
    area_sqm     NUMERIC(8, 2),
    floor        VARCHAR(50),
    base_rate    NUMERIC(10, 2) NOT NULL DEFAULT 0,
    rate_type    VARCHAR(10) NOT NULL DEFAULT 'daily' CHECK (rate_type IN ('hourly', 'daily')),
    amenities    TEXT[] NOT NULL DEFAULT '{}',
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    notes        TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_venues_name_org UNIQUE (org_id, name)
);

CREATE INDEX idx_venues_org_id      ON venues(org_id);
CREATE INDEX idx_venues_branch_id   ON venues(branch_id);
CREATE INDEX idx_venues_venue_type  ON venues(venue_type);
CREATE INDEX idx_venues_is_available ON venues(is_available);
