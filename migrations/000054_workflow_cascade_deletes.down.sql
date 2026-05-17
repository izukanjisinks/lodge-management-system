-- Revert CASCADE back to default RESTRICT behaviour.

ALTER TABLE workflow_instances
    DROP CONSTRAINT IF EXISTS workflow_instances_workflow_id_fkey,
    ADD CONSTRAINT workflow_instances_workflow_id_fkey
        FOREIGN KEY (workflow_id) REFERENCES workflows(id);

ALTER TABLE workflow_instances
    DROP CONSTRAINT IF EXISTS workflow_instances_current_step_id_fkey,
    ADD CONSTRAINT workflow_instances_current_step_id_fkey
        FOREIGN KEY (current_step_id) REFERENCES workflow_steps(id);

ALTER TABLE assigned_tasks
    DROP CONSTRAINT IF EXISTS assigned_tasks_step_id_fkey,
    ADD CONSTRAINT assigned_tasks_step_id_fkey
        FOREIGN KEY (step_id) REFERENCES workflow_steps(id);
