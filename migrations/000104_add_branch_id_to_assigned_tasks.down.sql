DROP INDEX IF EXISTS idx_assigned_tasks_branch_id;

ALTER TABLE assigned_tasks DROP COLUMN IF EXISTS branch_id;
