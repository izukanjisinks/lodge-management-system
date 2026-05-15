-- Backfill allowed_roles on existing workflow transitions that were seeded
-- before allowed_roles was moved from steps to transitions.

-- "Submit for Approval" transition: receptionist, manager, or admin can submit
UPDATE workflow_transitions
SET allowed_roles = '["receptionist","admin","manager"]'::jsonb
WHERE action_name = 'Submit for Approval'
  AND allowed_roles = '[]'::jsonb;

-- "approve" / "Approve" transitions: manager or admin only
UPDATE workflow_transitions
SET allowed_roles = '["manager","admin"]'::jsonb
WHERE action_name ILIKE 'approve'
  AND allowed_roles = '[]'::jsonb;

-- "reject" / "Reject" transitions: manager or admin only
UPDATE workflow_transitions
SET allowed_roles = '["manager","admin"]'::jsonb
WHERE action_name ILIKE 'reject'
  AND allowed_roles = '[]'::jsonb;
