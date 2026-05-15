ALTER TABLE workflow_history
    DROP CONSTRAINT IF EXISTS workflow_history_from_step_id_fkey,
    ADD CONSTRAINT workflow_history_from_step_id_fkey
        FOREIGN KEY (from_step_id) REFERENCES workflow_steps(id);

ALTER TABLE workflow_history
    DROP CONSTRAINT IF EXISTS workflow_history_to_step_id_fkey,
    ADD CONSTRAINT workflow_history_to_step_id_fkey
        FOREIGN KEY (to_step_id) REFERENCES workflow_steps(id);
