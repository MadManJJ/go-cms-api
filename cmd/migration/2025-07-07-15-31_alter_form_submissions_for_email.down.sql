ALTER TABLE form_submissions
DROP CONSTRAINT IF EXISTS form_submissions_submitter_user_id_fkey;

ALTER TABLE form_submissions
DROP COLUMN IF EXISTS submitter_user_id;

ALTER TABLE form_submissions
ADD COLUMN IF NOT EXISTS submitted_email VARCHAR(255);

DROP INDEX IF EXISTS idx_form_submissions_submitter_user_id;

CREATE INDEX IF NOT EXISTS idx_form_submissions_submitted_email ON form_submissions(submitted_email);