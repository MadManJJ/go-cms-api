-- Allow NULL email and password (for LINE users)
ALTER TABLE users
ALTER COLUMN email DROP NOT NULL,
ALTER COLUMN password DROP NOT NULL;

-- Add provider column with default value
ALTER TABLE users
ADD COLUMN provider VARCHAR(20) DEFAULT 'normal';

-- Add a check constraint to ensure only allowed provider values
ALTER TABLE users
ADD CONSTRAINT provider_check CHECK (provider IN ('normal', 'line'));
