-- Remove seeded Leave Approval workflow data
DELETE FROM workflow_transitions WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Leave Request Approval');
DELETE FROM workflow_steps WHERE workflow_id = (SELECT id FROM workflows WHERE name = 'Leave Request Approval');
DELETE FROM workflows WHERE name = 'Leave Request Approval';
