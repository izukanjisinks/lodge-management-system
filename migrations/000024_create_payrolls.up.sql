CREATE TABLE IF NOT EXISTS payrolls (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'PROCESSING', 'COMPLETED', 'CANCELLED')),
    processed_by    UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
    processed_at    TIMESTAMPTZ NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payrolls_status ON payrolls(status);
CREATE INDEX idx_payrolls_period ON payrolls(start_date, end_date);
