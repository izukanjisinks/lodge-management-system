ALTER TABLE individual_profiles
    ADD COLUMN user_id UUID NULL REFERENCES users(user_id) ON DELETE SET NULL;

CREATE UNIQUE INDEX idx_individual_profiles_user_id ON individual_profiles(user_id)
    WHERE user_id IS NOT NULL;
