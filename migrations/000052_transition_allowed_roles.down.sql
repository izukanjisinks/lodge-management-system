ALTER TABLE workflow_transitions DROP COLUMN IF EXISTS allowed_roles;

ALTER TABLE workflow_steps
    ADD COLUMN allowed_roles         JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN requires_all_approvers BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN min_approvals          INT NOT NULL DEFAULT 0;
