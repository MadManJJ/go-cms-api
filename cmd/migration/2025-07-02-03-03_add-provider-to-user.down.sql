-- Remove the provider check constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS provider_check;

-- Remove the provider column
ALTER TABLE users
DROP COLUMN IF EXISTS provider;

-- Make email and password NOT NULL again
ALTER TABLE users
ALTER COLUMN email SET NOT NULL,
ALTER COLUMN password SET NOT NULL;
