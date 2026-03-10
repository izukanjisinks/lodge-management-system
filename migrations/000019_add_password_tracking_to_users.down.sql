-- Remove password tracking fields from users table
ALTER TABLE users
    DROP COLUMN IF EXISTS change_password,
    DROP COLUMN IF EXISTS password_changed_at,
    DROP COLUMN IF EXISTS password_expires_at,
    DROP COLUMN IF EXISTS failed_login_attempts,
    DROP COLUMN IF EXISTS is_locked,
    DROP COLUMN IF EXISTS locked_until;
