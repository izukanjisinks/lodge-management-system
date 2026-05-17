DELETE FROM individual_profiles WHERE guest_id IS NOT NULL;

ALTER TABLE individual_profiles ALTER COLUMN id_passport_number SET NOT NULL;
