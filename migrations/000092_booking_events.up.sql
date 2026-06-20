-- A booking_event is the venue reservation underpinning a conference/event booking.
-- One row per event booking (single venue, single date range — per current scope).
-- The booking's own status field holds the check-in/out state; events carry pricing.
CREATE TABLE booking_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id      UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    venue_id        UUID REFERENCES venues(id) ON DELETE SET NULL,

    event_type      VARCHAR(50) NOT NULL DEFAULT 'conference',
                       -- conference | seminar | workshop | gala | wedding | training | event

    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,

    pax_count       INT NOT NULL DEFAULT 0,

    -- Hire price locked at materialise time (copied from venue.base_rate). Per-day;
    -- the invoice multiplies by the number of days in [start_date, end_date].
    price           NUMERIC(10, 2) NOT NULL DEFAULT 0,

    notes           TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_booking_events_booking_id ON booking_events(booking_id);
CREATE INDEX idx_booking_events_venue_id   ON booking_events(venue_id) WHERE venue_id IS NOT NULL;
