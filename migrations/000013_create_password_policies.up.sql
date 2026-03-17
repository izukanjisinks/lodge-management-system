CREATE TABLE IF NOT EXISTS password_policies (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id       UUID,
    min_length            INT NOT NULL DEFAULT 8 CHECK (min_length >= 6 AND min_length <= 128),
    require_uppercase     BOOLEAN NOT NULL DEFAULT TRUE,
    require_lowercase     BOOLEAN NOT NULL DEFAULT TRUE,
    require_numbers       BOOLEAN NOT NULL DEFAULT TRUE,
    require_special_chars BOOLEAN NOT NULL DEFAULT TRUE,
    max_failed_attempts   INT NOT NULL DEFAULT 5 CHECK (max_failed_attempts >= 1 AND max_failed_attempts <= 100),
    lockout_duration_mins INT NOT NULL DEFAULT 30 CHECK (lockout_duration_mins >= 1 AND lockout_duration_mins <= 10080),
    password_expiry_days  INT CHECK (password_expiry_days IS NULL OR (password_expiry_days >= 1 AND password_expiry_days <= 365)),
    otp_length            INT NOT NULL DEFAULT 6 CHECK (otp_length >= 4 AND otp_length <= 10),
    otp_expiry_mins       INT NOT NULL DEFAULT 5 CHECK (otp_expiry_mins >= 1 AND otp_expiry_mins <= 60),
    session_timeout_mins  INT NOT NULL DEFAULT 30 CHECK (session_timeout_mins >= 1 AND session_timeout_mins <= 10080),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_password_policies_organization
    ON password_policies(organization_id)
    WHERE organization_id IS NOT NULL;

CREATE UNIQUE INDEX idx_password_policies_global
    ON password_policies((1))
    WHERE organization_id IS NULL;

-- Default global policy
INSERT INTO password_policies (
    organization_id, min_length, require_uppercase, require_lowercase,
    require_numbers, require_special_chars, max_failed_attempts,
    lockout_duration_mins, password_expiry_days, otp_length, otp_expiry_mins,
    session_timeout_mins
) VALUES (
    NULL, 8, TRUE, TRUE, TRUE, TRUE, 5, 30, 90, 6, 5, 30
);

CREATE TABLE IF NOT EXISTS password_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_history_user_id ON password_history(user_id);

-- Track failed login attempts and password change requirement on users
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS failed_login_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_locked             BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS locked_until          TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_login_at         TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS password_changed_at   TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS password_expires_at   TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS change_password        BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON TABLE password_policies IS 'Password policy configuration. NULL organization_id = global default.';
