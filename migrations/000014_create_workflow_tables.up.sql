-- Workflow system tables
-- This migration creates the tables for the generic workflow system

-- Table: workflows (the templates)
CREATE TABLE IF NOT EXISTS workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_by UUID NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: workflow_steps (stages in a workflow template)
CREATE TABLE IF NOT EXISTS workflow_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    step_name VARCHAR(255) NOT NULL,
    step_order INT NOT NULL,
    initial BOOLEAN DEFAULT false,
    final BOOLEAN DEFAULT false,
    allowed_roles JSONB DEFAULT '[]'::jsonb, -- Array of role codes
    requires_all_approvers BOOLEAN DEFAULT false,
    min_approvals INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_step_order UNIQUE (workflow_id, step_order)
);

-- Table: workflow_transitions (allowed movements between steps)
CREATE TABLE IF NOT EXISTS workflow_transitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_step_id UUID NOT NULL REFERENCES workflow_steps(id) ON DELETE CASCADE,
    to_step_id UUID NOT NULL REFERENCES workflow_steps(id) ON DELETE CASCADE,
    action_name VARCHAR(100) NOT NULL, -- 'submit', 'approve', 'reject', 'reassign'
    condition_type VARCHAR(100), -- 'user_role', 'assigned_user_only', etc.
    condition_value TEXT, -- JSON for complex conditions
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_transition CHECK (from_step_id != to_step_id)
);

-- Table: workflow_instances (actual executions of workflows)
CREATE TABLE IF NOT EXISTS workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id),
    current_step_id UUID NOT NULL REFERENCES workflow_steps(id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'in_progress', 'completed', 'rejected', 'cancelled'
    task_details JSONB NOT NULL, -- Contains TaskDetails structure
    created_by UUID NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    due_date TIMESTAMP,
    priority VARCHAR(20) DEFAULT 'medium', -- 'low', 'medium', 'high', 'urgent'
    CONSTRAINT valid_status CHECK (status IN ('pending', 'in_progress', 'completed', 'rejected', 'cancelled')),
    CONSTRAINT valid_priority CHECK (priority IN ('low', 'medium', 'high', 'urgent'))
);

-- Table: assigned_tasks (action items for specific users)
CREATE TABLE IF NOT EXISTS assigned_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    step_id UUID NOT NULL REFERENCES workflow_steps(id),
    step_name VARCHAR(255) NOT NULL, -- Denormalized for performance
    assigned_to UUID NOT NULL REFERENCES users(user_id),
    assigned_by UUID NOT NULL REFERENCES users(user_id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'in_progress', 'completed', 'skipped'
    due_date TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_task_status CHECK (status IN ('pending', 'in_progress', 'completed', 'skipped'))
);

-- Table: workflow_history (audit trail)
CREATE TABLE IF NOT EXISTS workflow_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    from_step_id UUID REFERENCES workflow_steps(id), -- Nullable for initial creation
    to_step_id UUID NOT NULL REFERENCES workflow_steps(id),
    action_taken VARCHAR(100) NOT NULL, -- 'submit', 'approve', 'reject', 'reassign'
    performed_by UUID NOT NULL REFERENCES users(user_id),
    performed_by_name VARCHAR(255) NOT NULL, -- Denormalized for history preservation
    comments TEXT,
    metadata JSONB, -- Additional context
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_workflow_steps_workflow ON workflow_steps(workflow_id);
CREATE INDEX idx_workflow_transitions_workflow ON workflow_transitions(workflow_id);
CREATE INDEX idx_workflow_transitions_from_step ON workflow_transitions(from_step_id);
CREATE INDEX idx_workflow_instances_workflow ON workflow_instances(workflow_id);
CREATE INDEX idx_workflow_instances_status ON workflow_instances(status);
CREATE INDEX idx_workflow_instances_created_by ON workflow_instances(created_by);
CREATE INDEX idx_assigned_tasks_instance ON assigned_tasks(instance_id);
CREATE INDEX idx_assigned_tasks_assigned_to ON assigned_tasks(assigned_to);
CREATE INDEX idx_assigned_tasks_status ON assigned_tasks(status);
CREATE INDEX idx_workflow_history_instance ON workflow_history(instance_id);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_workflow_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_workflows_updated_at
    BEFORE UPDATE ON workflows
    FOR EACH ROW
    EXECUTE FUNCTION update_workflow_updated_at();

CREATE TRIGGER trigger_workflow_instances_updated_at
    BEFORE UPDATE ON workflow_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_workflow_updated_at();

CREATE TRIGGER trigger_assigned_tasks_updated_at
    BEFORE UPDATE ON assigned_tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_workflow_updated_at();
