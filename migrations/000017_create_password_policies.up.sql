-- Create password_policies table
CREATE TABLE password_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID, -- Reserved for future multi-tenant support
    min_length INT NOT NULL DEFAULT 8 CHECK (min_length >= 6 AND min_length <= 128),
    require_uppercase BOOLEAN NOT NULL DEFAULT true,
    require_lowercase BOOLEAN NOT NULL DEFAULT true,
    require_numbers BOOLEAN NOT NULL DEFAULT true,
    require_special_chars BOOLEAN NOT NULL DEFAULT true,
    max_failed_attempts INT NOT NULL DEFAULT 5 CHECK (max_failed_attempts >= 1 AND max_failed_attempts <= 100),
    lockout_duration_mins INT NOT NULL DEFAULT 30 CHECK (lockout_duration_mins >= 1 AND lockout_duration_mins <= 10080),
    password_expiry_days INT CHECK (password_expiry_days IS NULL OR (password_expiry_days >= 1 AND password_expiry_days <= 365)),
    otp_length INT NOT NULL DEFAULT 6 CHECK (otp_length >= 4 AND otp_length <= 10),
    otp_expiry_mins INT NOT NULL DEFAULT 5 CHECK (otp_expiry_mins >= 1 AND otp_expiry_mins <= 60),
    session_timeout_mins INT NOT NULL DEFAULT 30 CHECK (session_timeout_mins >= 1 AND session_timeout_mins <= 10080),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create unique index to ensure only one policy per organization
CREATE UNIQUE INDEX idx_password_policies_organization
    ON password_policies(organization_id)
    WHERE organization_id IS NOT NULL;

-- Create unique index to ensure only one global default policy (where organization_id IS NULL)
CREATE UNIQUE INDEX idx_password_policies_global
    ON password_policies((1))
    WHERE organization_id IS NULL;

-- Insert default global password policy
INSERT INTO password_policies (
    id,
    organization_id,
    min_length,
    require_uppercase,
    require_lowercase,
    require_numbers,
    require_special_chars,
    max_failed_attempts,
    lockout_duration_mins,
    password_expiry_days,
    otp_length,
    otp_expiry_mins,
    session_timeout_mins,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    NULL, -- Global default
    8,
    true,
    true,
    true,
    true,
    5,
    30,
    90,
    6,
    5,
    30,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Add comment to table
COMMENT ON TABLE password_policies IS 'Stores password policy configurations. NULL organization_id indicates global default policy.';
