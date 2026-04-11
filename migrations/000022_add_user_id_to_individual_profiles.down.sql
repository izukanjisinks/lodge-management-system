DROP INDEX IF EXISTS idx_individual_profiles_user_id;
ALTER TABLE individual_profiles DROP COLUMN IF EXISTS user_id;
