ALTER TABLE corporate_booking_requests ADD COLUMN IF NOT EXISTS web_user_id UUID REFERENCES web_users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_corporate_booking_requests_web_user_id ON corporate_booking_requests(web_user_id);
