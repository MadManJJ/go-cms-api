DROP TRIGGER IF EXISTS set_timestamp_forms ON forms;
DROP TRIGGER IF EXISTS set_timestamp_form_sections ON form_sections;
DROP TRIGGER IF EXISTS set_timestamp_form_fields ON form_fields;
DROP TRIGGER IF EXISTS set_timestamp_form_submissions ON form_submissions;

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS form_submissions CASCADE;
DROP TABLE IF EXISTS form_fields CASCADE;
DROP TABLE IF EXISTS form_sections CASCADE;
DROP TABLE IF EXISTS forms CASCADE;

-- (Optional) Drop ENUM types if they are no longer needed by any other table
DROP TYPE IF EXISTS form_field_type_enum;

DROP FUNCTION IF EXISTS trigger_set_timestamp();