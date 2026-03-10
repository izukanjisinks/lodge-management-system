-- Remove the reseeded Leave Request Approval workflow (cascades to steps and transitions)
DELETE FROM workflows WHERE name = 'Leave Request Approval';
