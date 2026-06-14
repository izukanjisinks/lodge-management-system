CREATE TABLE corporate_guests (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    corporate_profile_id UUID NOT NULL REFERENCES cor_profiles(id) ON DELETE CASCADE,
    first_name           VARCHAR(100) NOT NULL,
    last_name            VARCHAR(100) NOT NULL,
    phone                VARCHAR(50),
    email                VARCHAR(255),
    identification_card  VARCHAR(100) NOT NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_corporate_guests_corporate_profile_id ON corporate_guests(corporate_profile_id);
