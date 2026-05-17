CREATE TABLE IF NOT EXISTS guests (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                 VARCHAR(255) UNIQUE NOT NULL,
    password              VARCHAR(255) NOT NULL,
    full_name             VARCHAR(255) NOT NULL DEFAULT '',
    phone                 VARCHAR(50),
    is_active             BOOLEAN NOT NULL DEFAULT TRUE,
    change_password       BOOLEAN NOT NULL DEFAULT FALSE,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    is_locked             BOOLEAN NOT NULL DEFAULT FALSE,
    locked_until          TIMESTAMPTZ,
    last_login_at         TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_guests_email ON guests(email);

-- Replace user_id link on individual_profiles with guest_id
ALTER TABLE individual_profiles
    ADD COLUMN IF NOT EXISTS guest_id UUID REFERENCES guests(id) ON DELETE SET NULL;

ALTER TABLE individual_profiles DROP COLUMN IF EXISTS user_id;
