-- Move role-based access control from steps to transitions.
-- A transition defines who can trigger it, not the step itself.
ALTER TABLE workflow_transitions
    ADD COLUMN allowed_roles JSONB NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE workflow_steps
    DROP COLUMN IF EXISTS allowed_roles,
    DROP COLUMN IF EXISTS requires_all_approvers,
    DROP COLUMN IF EXISTS min_approvals;
