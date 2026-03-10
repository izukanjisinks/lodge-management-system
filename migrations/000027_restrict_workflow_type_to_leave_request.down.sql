-- Restore the original broad check constraint
ALTER TABLE workflows
DROP CONSTRAINT valid_workflow_type;

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
