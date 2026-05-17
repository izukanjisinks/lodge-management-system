-- workflow_history references workflow_steps via from_step_id and to_step_id
-- without CASCADE, deleting a step (or workflow) errors on these FKs.

ALTER TABLE workflow_history
    DROP CONSTRAINT IF EXISTS workflow_history_from_step_id_fkey,
    ADD CONSTRAINT workflow_history_from_step_id_fkey
        FOREIGN KEY (from_step_id) REFERENCES workflow_steps(id) ON DELETE CASCADE;

ALTER TABLE workflow_history
    DROP CONSTRAINT IF EXISTS workflow_history_to_step_id_fkey,
    ADD CONSTRAINT workflow_history_to_step_id_fkey
        FOREIGN KEY (to_step_id) REFERENCES workflow_steps(id) ON DELETE CASCADE;
