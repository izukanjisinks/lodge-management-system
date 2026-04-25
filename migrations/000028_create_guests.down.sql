ALTER TABLE individual_profiles ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(user_id) ON DELETE SET NULL;
ALTER TABLE individual_profiles DROP COLUMN IF EXISTS guest_id;

DROP INDEX IF EXISTS idx_guests_email;
DROP TABLE IF EXISTS guests;
