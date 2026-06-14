ALTER TABLE orders
    ADD COLUMN attendee_id UUID REFERENCES booking_attendees(id) ON DELETE SET NULL;

CREATE INDEX idx_orders_attendee_id ON orders(attendee_id) WHERE attendee_id IS NOT NULL;
