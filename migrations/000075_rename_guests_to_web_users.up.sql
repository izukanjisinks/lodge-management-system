-- Rename guests table to web_users and drop individual_profiles.
-- web_users is the primary identity table for website users.
-- individual_profiles is no longer needed — its guest_id FK pointed to guests(id)
-- which is now web_users, and the profile data it held is superseded by the
-- new corporate profile redesign (cor_profiles, corporate_guests).

ALTER TABLE guests RENAME TO web_users;

ALTER INDEX idx_guests_email RENAME TO idx_web_users_email;

-- reviews.guest_id FK points to individual_profiles(id); drop it before the table goes away.
-- Reviewer identity is now derived from the booking's web_user_id.
ALTER TABLE reviews DROP COLUMN IF EXISTS guest_id;
DROP INDEX IF EXISTS idx_reviews_guest_id;

DROP TABLE IF EXISTS individual_profiles;
