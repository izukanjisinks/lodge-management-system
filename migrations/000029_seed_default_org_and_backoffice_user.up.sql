-- Seed default organization and backoffice user
-- Backoffice@123  (bcrypt cost 10)
--
-- backoffice@lodge.dev  → Backoffice@123

DO $$
DECLARE
    v_org_id UUID;
BEGIN

INSERT INTO organizations (name, email)
VALUES ('The Sanctuary Lodge', 'admin@lodge.dev')
ON CONFLICT DO NOTHING;

SELECT id INTO v_org_id FROM organizations WHERE email = 'admin@lodge.dev' LIMIT 1;

-- Backfill org_id on all existing seeded staff data
UPDATE users               SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE roles               SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE rooms               SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE bookings            SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE individual_profiles SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE corporate_profiles  SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE meal_plans          SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE invoices            SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE workflows           SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE workflow_instances  SET org_id = v_org_id WHERE org_id IS NULL;
UPDATE assigned_tasks      SET org_id = v_org_id WHERE org_id IS NULL;

-- Seed default backoffice user
INSERT INTO backoffice_users (full_name, email, password)
VALUES (
    'Platform Admin',
    'backoffice@lodge.dev',
    '$2a$10$j.VRwVUG44TQbGALiLJWSOOHaXh.rv4BxTOA2amgzoRBT9tK9qdI2'
)
ON CONFLICT (email) DO NOTHING;

RAISE NOTICE 'Default organization and backoffice user seeded successfully';

END $$;
