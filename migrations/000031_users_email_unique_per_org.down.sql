ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_org_unique;
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
