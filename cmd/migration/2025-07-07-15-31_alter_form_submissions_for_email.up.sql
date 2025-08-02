DROP INDEX IF EXISTS idx_form_submissions_submitted_email;

-- Step 2: Drop the new submitted_email column
ALTER TABLE form_submissions
DROP COLUMN IF EXISTS submitted_email;
