ALTER TABLE bookings
    ADD COLUMN corporate_client_id UUID REFERENCES corporate_profiles(id) ON DELETE SET NULL;

CREATE INDEX idx_bookings_corporate_client_id ON bookings(corporate_client_id) WHERE corporate_client_id IS NOT NULL;
