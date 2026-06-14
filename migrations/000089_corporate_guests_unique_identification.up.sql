-- Remove duplicate roster rows (keep the earliest per profile+ID card) so the
-- unique constraint can be applied. Re-point any booking_attendees that
-- referenced a now-deleted duplicate to the surviving row.
WITH ranked AS (
    SELECT id,
           corporate_profile_id,
           identification_card,
           ROW_NUMBER() OVER (
               PARTITION BY corporate_profile_id, identification_card
               ORDER BY created_at ASC, id ASC
           ) AS rn,
           FIRST_VALUE(id) OVER (
               PARTITION BY corporate_profile_id, identification_card
               ORDER BY created_at ASC, id ASC
           ) AS keep_id
    FROM corporate_guests
)
UPDATE booking_attendees ba
SET corporate_guest_id = ranked.keep_id
FROM ranked
WHERE ba.corporate_guest_id = ranked.id
  AND ranked.rn > 1;

WITH ranked AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY corporate_profile_id, identification_card
               ORDER BY created_at ASC, id ASC
           ) AS rn
    FROM corporate_guests
)
DELETE FROM corporate_guests
WHERE id IN (SELECT id FROM ranked WHERE rn > 1);

ALTER TABLE corporate_guests
    ADD CONSTRAINT uq_corporate_guests_profile_id_card
    UNIQUE (corporate_profile_id, identification_card);
