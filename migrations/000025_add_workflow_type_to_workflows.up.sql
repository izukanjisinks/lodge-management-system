ALTER TABLE workflows ADD COLUMN IF NOT EXISTS workflow_type VARCHAR(100) UNIQUE;

-- Backfill the seeded Booking Approval Workflow
UPDATE workflows SET workflow_type = 'BOOKING_APPROVAL' WHERE name = 'Booking Approval Workflow';
