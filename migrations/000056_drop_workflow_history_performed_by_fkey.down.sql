ALTER TABLE workflow_history
    ADD CONSTRAINT workflow_history_performed_by_fkey
        FOREIGN KEY (performed_by) REFERENCES users(user_id);
