ALTER TABLE users
    ADD COLUMN IF NOT EXISTS full_name VARCHAR(255) NOT NULL DEFAULT '';

-- Backfill names for dev seed users (no-op if they don't exist)
UPDATE users SET full_name = 'System Admin'      WHERE email = 'admin@lodge.dev'        AND full_name = '';
UPDATE users SET full_name = 'Lodge Manager'     WHERE email = 'manager@lodge.dev'      AND full_name = '';
UPDATE users SET full_name = 'Front Desk'        WHERE email = 'receptionist@lodge.dev' AND full_name = '';
UPDATE users SET full_name = 'House Keeping'     WHERE email = 'cleaner@lodge.dev'      AND full_name = '';
