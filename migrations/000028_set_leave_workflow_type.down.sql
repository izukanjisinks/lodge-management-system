UPDATE workflows
SET workflow_type = NULL
WHERE name = 'Leave Request Approval'
  AND workflow_type = 'LEAVE_REQUEST';
