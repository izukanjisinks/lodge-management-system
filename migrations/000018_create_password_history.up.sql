-- Create password_history table
CREATE TABLE password_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id and created_at for efficient lookups
CREATE INDEX idx_password_history_user_created
    ON password_history(user_id, created_at DESC);

-- Add comment to table
COMMENT ON TABLE password_history IS 'Tracks password history to prevent password reuse. Stores only hashed passwords.';
