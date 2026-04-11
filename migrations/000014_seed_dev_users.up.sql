-- Seed development/test users for the Lodge Management System
-- Passwords shown in comments. DO NOT use these credentials in production.
--
-- admin@lodge.dev        → Admin@123
-- manager@lodge.dev      → Manager@123
-- receptionist@lodge.dev → Reception@123
-- cleaner@lodge.dev      → Cleaner@123

DO $$
DECLARE
    v_role_admin        UUID;
    v_role_manager      UUID;
    v_role_receptionist UUID;
    v_role_cleaner      UUID;
BEGIN

SELECT role_id INTO v_role_admin        FROM roles WHERE name = 'admin'        LIMIT 1;
SELECT role_id INTO v_role_manager      FROM roles WHERE name = 'manager'      LIMIT 1;
SELECT role_id INTO v_role_receptionist FROM roles WHERE name = 'receptionist' LIMIT 1;
SELECT role_id INTO v_role_cleaner      FROM roles WHERE name = 'cleaner'      LIMIT 1;

IF v_role_admin IS NULL THEN
    RAISE EXCEPTION 'Roles not seeded — ensure migration 000002 has run first';
END IF;

INSERT INTO users (email, password, role_id, is_active) VALUES
    -- Admin@123  (bcrypt cost 10)
    ('admin@lodge.dev',
     '$2a$10$CQCThSSuQ1J9o84X3C/1juRXCdfHXniFg8Pcj8.CLSk9UURxIVcEG',
     v_role_admin, TRUE),

    -- Manager@123  (bcrypt cost 10)
    ('manager@lodge.dev',
     '$2a$10$FwbvszWKQiWWGeoM1rQu3uqJbsougE9wWd4pTMJynzpxK1Sav4apy',
     v_role_manager, TRUE),

    -- Reception@123  (bcrypt cost 10)
    ('receptionist@lodge.dev',
     '$2a$10$sFxr22ih4Hh/u0wSX840PehTpGC8zzM8TDxgLMoa3ooIirnlMN07G',
     v_role_receptionist, TRUE),

    -- Cleaner@123  (bcrypt cost 10)
    ('cleaner@lodge.dev',
     '$2a$10$TBaW30/B6S.2rCS5VYcnrO/zO/T12tbwIk4hu.5dvPDhrYAnSji5a',
     v_role_cleaner, TRUE)

ON CONFLICT (email) DO NOTHING;

RAISE NOTICE 'Dev users seeded successfully';

END $$;
