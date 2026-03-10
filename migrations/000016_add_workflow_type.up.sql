-- Add workflow_type column to workflows table
-- This ensures each workflow has a unique type identifier (e.g., LEAVE_REQUEST, EMPLOYEE_ONBOARDING)
-- Only one workflow should exist per type

ALTER TABLE workflows
ADD COLUMN workflow_type VARCHAR(100);

-- Add unique constraint to ensure only one workflow per type
CREATE UNIQUE INDEX idx_workflows_type_unique ON workflows(workflow_type)
WHERE workflow_type IS NOT NULL AND is_active = true;

-- Add check constraint for valid workflow types
ALTER TABLE workflows
ADD CONSTRAINT valid_workflow_type CHECK (
    workflow_type IN (
        'LEAVE_REQUEST',
        'EMPLOYEE_ONBOARDING',
        'EMPLOYEE_OFFBOARDING',
        'PERFORMANCE_REVIEW',
        'EXPENSE_CLAIM',
        'TRAINING_REQUEST',
        'ASSET_REQUEST',
        'GRIEVANCE',
        'PROMOTION_REQUEST',
        'TRANSFER_REQUEST'
    )
);

COMMENT ON COLUMN workflows.workflow_type IS 'Unique identifier for workflow purpose. Only one active workflow allowed per type.';
