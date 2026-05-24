-- Remove branch assignment from users, bookings, and rooms,
-- then delete the seeded default branches.
UPDATE users SET branch_id = NULL
WHERE branch_id IN (SELECT id FROM branches WHERE name LIKE '% (Main Branch)');

UPDATE bookings SET branch_id = NULL
WHERE branch_id IN (SELECT id FROM branches WHERE name LIKE '% (Main Branch)');

UPDATE rooms SET branch_id = NULL
WHERE branch_id IN (SELECT id FROM branches WHERE name LIKE '% (Main Branch)');

DELETE FROM branches WHERE name LIKE '% (Main Branch)';
