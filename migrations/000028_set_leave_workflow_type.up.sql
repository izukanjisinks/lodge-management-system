-- Set workflow_type on the Leave Request Approval workflow (was missing from seed)
UPDATE workflows
SET workflow_type = 'LEAVE_REQUEST'
WHERE name = 'Leave Request Approval'
  AND workflow_type IS NULL;
