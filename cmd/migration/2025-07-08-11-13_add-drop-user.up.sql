ALTER TABLE forms
DROP COLUMN IF EXISTS created_by_user_id;

DROP INDEX IF EXISTS idx_forms_created_by_user_id;