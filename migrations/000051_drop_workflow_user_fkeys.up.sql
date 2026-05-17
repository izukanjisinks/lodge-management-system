ALTER TABLE workflow_instances DROP CONSTRAINT IF EXISTS workflow_instances_created_by_fkey;
ALTER TABLE assigned_tasks DROP CONSTRAINT IF EXISTS assigned_tasks_created_by_fkey;
ALTER TABLE assigned_tasks DROP CONSTRAINT IF EXISTS assigned_tasks_assigned_to_fkey;
ALTER TABLE assigned_tasks DROP CONSTRAINT IF EXISTS assigned_tasks_assigned_by_fkey;
