-- workflow_type was globally UNIQUE, preventing multiple orgs from using the same
-- workflow type (e.g. BOOKING_APPROVAL). Change to unique per org instead.
ALTER TABLE workflows DROP CONSTRAINT IF EXISTS workflows_workflow_type_key;
ALTER TABLE workflows ADD CONSTRAINT workflows_workflow_type_org_unique UNIQUE (workflow_type, org_id);
