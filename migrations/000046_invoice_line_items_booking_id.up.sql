ALTER TABLE invoice_line_items
    ADD COLUMN booking_id UUID REFERENCES bookings(id) ON DELETE SET NULL;

CREATE INDEX idx_invoice_line_items_booking_id ON invoice_line_items(booking_id) WHERE booking_id IS NOT NULL;
