-- Reverse demo seed — removes all seeded demo data in dependency order
DELETE FROM invoice_line_items
WHERE invoice_id IN (
    SELECT id FROM invoices WHERE invoice_number LIKE 'INV-%-000%'
);

DELETE FROM invoices WHERE invoice_number LIKE 'INV-%-000%';

DELETE FROM booking_meal_plans
WHERE booking_id IN (
    SELECT b.id FROM bookings b
    JOIN individual_profiles ip ON b.client_id = ip.id
    WHERE ip.email IN ('john.mwewa@gmail.com','sarah.banda@gmail.com','peter.okafor@gmail.com')
    UNION
    SELECT b.id FROM bookings b
    JOIN corporate_profiles cp ON b.client_id = cp.id
    WHERE cp.email IN ('bookings@acme.co.zm','events@zesco.co.zm')
);

DELETE FROM bookings
WHERE client_id IN (
    SELECT id FROM individual_profiles WHERE email IN ('john.mwewa@gmail.com','sarah.banda@gmail.com','peter.okafor@gmail.com')
    UNION
    SELECT id FROM corporate_profiles  WHERE email IN ('bookings@acme.co.zm','events@zesco.co.zm')
);

DELETE FROM meal_plans WHERE name IN ('Bed & Breakfast','Half Board','Full Board');

DELETE FROM corporate_profiles  WHERE email IN ('bookings@acme.co.zm','events@zesco.co.zm');
DELETE FROM individual_profiles WHERE email IN ('john.mwewa@gmail.com','sarah.banda@gmail.com','peter.okafor@gmail.com');

-- Restore room availability before deleting
UPDATE rooms SET is_available = TRUE WHERE name IN ('Room 201','Suite 301');
DELETE FROM rooms WHERE name IN ('Room 101','Room 102','Room 201','Room 202','Suite 301','Cabin 01','Conference A');
