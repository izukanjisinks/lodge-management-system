CREATE TYPE leave_request_status AS ENUM ('pending', 'approved', 'rejected', 'cancelled');

CREATE TABLE IF NOT EXISTS leave_requests (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id    UUID NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
    leave_type_id  UUID NOT NULL REFERENCES leave_types(id) ON DELETE RESTRICT,
    start_date     DATE NOT NULL,
    end_date       DATE NOT NULL,
    total_days     INT NOT NULL DEFAULT 0,
    reason         TEXT NOT NULL DEFAULT '',
    status         leave_request_status NOT NULL DEFAULT 'pending',
    reviewed_by    UUID NULL REFERENCES employees(id) ON DELETE SET NULL,
    reviewed_at    TIMESTAMPTZ NULL,
    review_comment TEXT NOT NULL DEFAULT '',
    attachment_url TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leave_requests_employee ON leave_requests(employee_id);
CREATE INDEX idx_leave_requests_status ON leave_requests(status);
CREATE INDEX idx_leave_requests_dates ON leave_requests(start_date, end_date);
