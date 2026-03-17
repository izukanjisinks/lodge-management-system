-- Workflow system — generic engine reused from HR system foundation

CREATE TABLE IF NOT EXISTS workflows (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    is_active   BOOLEAN DEFAULT TRUE,
    created_by  UUID NOT NULL REFERENCES users(user_id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workflow_steps (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id           UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    step_name             VARCHAR(255) NOT NULL,
    step_order            INT NOT NULL,
    initial               BOOLEAN DEFAULT FALSE,
    final                 BOOLEAN DEFAULT FALSE,
    allowed_roles         JSONB DEFAULT '[]'::jsonb,
    requires_all_approvers BOOLEAN DEFAULT FALSE,
    min_approvals         INT DEFAULT 0,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_step_order UNIQUE (workflow_id, step_order)
);

CREATE TABLE IF NOT EXISTS workflow_transitions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_step_id    UUID NOT NULL REFERENCES workflow_steps(id) ON DELETE CASCADE,
    to_step_id      UUID NOT NULL REFERENCES workflow_steps(id) ON DELETE CASCADE,
    action_name     VARCHAR(100) NOT NULL,
    condition_type  VARCHAR(100),
    condition_value TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_transition CHECK (from_step_id != to_step_id)
);

CREATE TABLE IF NOT EXISTS workflow_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id),
    current_step_id UUID NOT NULL REFERENCES workflow_steps(id),
    status          VARCHAR(50) NOT NULL DEFAULT 'pending',
    task_details    JSONB NOT NULL,
    created_by      UUID NOT NULL REFERENCES users(user_id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    due_date        TIMESTAMPTZ,
    priority        VARCHAR(20) DEFAULT 'medium',
    CONSTRAINT valid_instance_status CHECK (status IN ('pending', 'in_progress', 'completed', 'rejected', 'cancelled')),
    CONSTRAINT valid_priority CHECK (priority IN ('low', 'medium', 'high', 'urgent'))
);

CREATE TABLE IF NOT EXISTS assigned_tasks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id  UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    step_id      UUID NOT NULL REFERENCES workflow_steps(id),
    step_name    VARCHAR(255) NOT NULL,
    assigned_to  UUID NOT NULL REFERENCES users(user_id),
    assigned_by  UUID NOT NULL REFERENCES users(user_id),
    status       VARCHAR(50) NOT NULL DEFAULT 'pending',
    due_date     TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_task_status CHECK (status IN ('pending', 'in_progress', 'completed', 'skipped'))
);

CREATE TABLE IF NOT EXISTS workflow_history (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id       UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    from_step_id      UUID REFERENCES workflow_steps(id),
    to_step_id        UUID NOT NULL REFERENCES workflow_steps(id),
    action_taken      VARCHAR(100) NOT NULL,
    performed_by      UUID NOT NULL REFERENCES users(user_id),
    performed_by_name VARCHAR(255) NOT NULL,
    comments          TEXT,
    metadata          JSONB,
    timestamp         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
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

-- updated_at triggers
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_workflows_updated_at
    BEFORE UPDATE ON workflows
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_workflow_instances_updated_at
    BEFORE UPDATE ON workflow_instances
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_assigned_tasks_updated_at
    BEFORE UPDATE ON assigned_tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
