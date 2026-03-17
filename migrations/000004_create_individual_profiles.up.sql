CREATE TABLE IF NOT EXISTS individual_profiles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    full_name   VARCHAR(255) NOT NULL,
    phone       VARCHAR(50),
    id_number   VARCHAR(100),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_individual_profiles_user_id ON individual_profiles(user_id);
