ALTER TABLE assigned_tasks ADD COLUMN IF NOT EXISTS branch_id UUID REFERENCES branches(id);

CREATE INDEX idx_assigned_tasks_branch_id ON assigned_tasks(branch_id);
