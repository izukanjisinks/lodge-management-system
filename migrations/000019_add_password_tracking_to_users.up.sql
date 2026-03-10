-- Add password tracking fields to users table
ALTER TABLE users
    ADD COLUMN change_password BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN password_changed_at TIMESTAMP,
    ADD COLUMN password_expires_at TIMESTAMP,
    ADD COLUMN failed_login_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN is_locked BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN locked_until TIMESTAMP;

-- Add index for locked account queries
CREATE INDEX idx_users_locked_until ON users(locked_until) WHERE locked_until IS NOT NULL;

-- Add comment to new columns
COMMENT ON COLUMN users.change_password IS 'Flag to force user to change password on next login';
COMMENT ON COLUMN users.password_changed_at IS 'Timestamp when password was last changed';
COMMENT ON COLUMN users.password_expires_at IS 'Timestamp when current password expires';
COMMENT ON COLUMN users.failed_login_attempts IS 'Counter for tracking failed login attempts';
COMMENT ON COLUMN users.is_locked IS 'Permanent account lock flag (admin action)';
COMMENT ON COLUMN users.locked_until IS 'Timestamp until which account is temporarily locked due to failed attempts';
