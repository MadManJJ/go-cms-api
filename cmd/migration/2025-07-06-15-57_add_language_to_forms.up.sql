-- Add the notification language column to the forms table
-- This assumes the 'page_language' ENUM type ('th', 'en') already exists.
-- CREATE TYPE page_language AS ENUM ('th', 'en');
ALTER TABLE forms
ADD COLUMN IF NOT EXISTS language page_language;


