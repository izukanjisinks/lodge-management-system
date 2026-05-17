-- Make id_passport_number nullable so self-registered guests (who have no NRC/passport
-- at sign-up time) can still have a profile row. NULL values never collide in unique indexes.
ALTER TABLE individual_profiles ALTER COLUMN id_passport_number DROP NOT NULL;

-- Backfill individual_profiles for guests that registered before
-- the registration flow was updated to create a profile automatically.
INSERT INTO individual_profiles (id, guest_id, full_name, email, phone, status, created_at, updated_at)
SELECT
    gen_random_uuid(),
    g.id,
    g.full_name,
    g.email,
    g.phone,
    'active',
    NOW(),
    NOW()
FROM guests g
WHERE NOT EXISTS (
    SELECT 1 FROM individual_profiles ip WHERE ip.guest_id = g.id
);
