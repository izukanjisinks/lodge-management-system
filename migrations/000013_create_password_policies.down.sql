ALTER TABLE users
    DROP COLUMN IF EXISTS failed_login_attempts,
    DROP COLUMN IF EXISTS locked_until,
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS password_changed_at,
    DROP COLUMN IF EXISTS change_password;

DROP TABLE IF EXISTS password_history;
DROP TABLE IF EXISTS password_policies;
