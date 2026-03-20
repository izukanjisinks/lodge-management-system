-- Reverse workflow seed — removes all seeded workflow data in dependency order

DELETE FROM workflow_history
WHERE instance_id IN (
    SELECT id FROM workflow_instances
    WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Booking Approval Workflow')
);

DELETE FROM assigned_tasks
WHERE instance_id IN (
    SELECT id FROM workflow_instances
    WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Booking Approval Workflow')
);

DELETE FROM workflow_instances
WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Booking Approval Workflow');

DELETE FROM workflow_transitions
WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Booking Approval Workflow');

DELETE FROM workflow_steps
WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Booking Approval Workflow');

DELETE FROM workflows WHERE name = 'Booking Approval Workflow';
