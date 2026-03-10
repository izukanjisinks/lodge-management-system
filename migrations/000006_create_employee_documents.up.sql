CREATE TYPE document_type AS ENUM ('contract', 'id_document', 'certification', 'offer_letter', 'warning_letter', 'other');

CREATE TABLE IF NOT EXISTS employee_documents (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id   UUID NOT NULL REFERENCES employees(id) ON DELETE RESTRICT,
    document_type document_type NOT NULL DEFAULT 'other',
    title         VARCHAR(255) NOT NULL,
    description   TEXT NOT NULL DEFAULT '',
    file_url      TEXT NOT NULL,
    file_name     VARCHAR(255) NOT NULL,
    file_size     BIGINT NOT NULL DEFAULT 0,
    mime_type     VARCHAR(100) NOT NULL DEFAULT '',
    uploaded_by   UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    expiry_date   DATE NULL,
    is_verified   BOOLEAN NOT NULL DEFAULT FALSE,
    verified_by   UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
    verified_at   TIMESTAMPTZ NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ NULL
);

CREATE INDEX idx_employee_documents_employee ON employee_documents(employee_id);
CREATE INDEX idx_employee_documents_type ON employee_documents(document_type);
CREATE INDEX idx_employee_documents_deleted ON employee_documents(deleted_at);
