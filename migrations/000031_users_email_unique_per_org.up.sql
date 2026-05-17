-- Drop the global unique constraint on email
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

-- Add a per-org unique constraint — same email allowed across different orgs
ALTER TABLE users ADD CONSTRAINT users_email_org_unique UNIQUE (email, org_id);
