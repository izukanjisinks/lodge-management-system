-- Replace RESTRICT foreign keys on workflow tables with CASCADE so that
-- deleting a workflow or step also removes its instances, tasks, and history.

-- workflow_instances.workflow_id
ALTER TABLE workflow_instances
    DROP CONSTRAINT IF EXISTS workflow_instances_workflow_id_fkey,
    ADD CONSTRAINT workflow_instances_workflow_id_fkey
        FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE;

-- workflow_instances.current_step_id
ALTER TABLE workflow_instances
    DROP CONSTRAINT IF EXISTS workflow_instances_current_step_id_fkey,
    ADD CONSTRAINT workflow_instances_current_step_id_fkey
        FOREIGN KEY (current_step_id) REFERENCES workflow_steps(id) ON DELETE CASCADE;

-- assigned_tasks.step_id
ALTER TABLE assigned_tasks
    DROP CONSTRAINT IF EXISTS assigned_tasks_step_id_fkey,
    ADD CONSTRAINT assigned_tasks_step_id_fkey
        FOREIGN KEY (step_id) REFERENCES workflow_steps(id) ON DELETE CASCADE;
