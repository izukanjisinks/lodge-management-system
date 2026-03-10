-- Rollback: Remove workflow_type column and constraints

DROP INDEX IF EXISTS idx_workflows_type_unique;

ALTER TABLE workflows
DROP CONSTRAINT IF EXISTS valid_workflow_type;

ALTER TABLE workflows
DROP COLUMN IF EXISTS workflow_type;
