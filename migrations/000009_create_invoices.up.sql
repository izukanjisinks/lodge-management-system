CREATE TYPE invoice_status AS ENUM ('draft', 'issued', 'paid', 'overdue');

CREATE TABLE IF NOT EXISTS invoices (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE RESTRICT,
    subtotal   NUMERIC(10, 2) NOT NULL CHECK (subtotal >= 0),
    tax        NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (tax >= 0),
    total      NUMERIC(10, 2) NOT NULL CHECK (total >= 0),
    status     invoice_status NOT NULL DEFAULT 'draft',
    issued_at  TIMESTAMPTZ,
    due_date   TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoices_booking_id ON invoices(booking_id);
CREATE INDEX idx_invoices_status ON invoices(status);
