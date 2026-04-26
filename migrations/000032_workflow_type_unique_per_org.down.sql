ALTER TABLE workflows DROP CONSTRAINT IF EXISTS workflows_workflow_type_org_unique;
ALTER TABLE workflows ADD CONSTRAINT workflows_workflow_type_key UNIQUE (workflow_type);
