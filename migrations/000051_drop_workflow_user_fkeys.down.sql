ALTER TABLE workflow_instances ADD CONSTRAINT workflow_instances_created_by_fkey FOREIGN KEY (created_by) REFERENCES users(user_id);
ALTER TABLE assigned_tasks ADD CONSTRAINT assigned_tasks_assigned_to_fkey FOREIGN KEY (assigned_to) REFERENCES users(user_id);
ALTER TABLE assigned_tasks ADD CONSTRAINT assigned_tasks_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES users(user_id);
