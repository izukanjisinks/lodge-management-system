-- Restrict workflow_type to only allow LEAVE_REQUEST
-- Drop the existing broad check constraint and replace with a narrower one

ALTER TABLE workflows
DROP CONSTRAINT valid_workflow_type;

ALTER TABLE workflows
ADD CONSTRAINT valid_workflow_type CHECK (
    workflow_type = 'LEAVE_REQUEST'
);
