ALTER TABLE individual_profiles
    ADD CONSTRAINT uq_individual_profiles_id_passport_number UNIQUE (id_passport_number, org_id);
