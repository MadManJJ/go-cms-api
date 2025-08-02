-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Or for PostgreSQL 13+: CREATE EXTENSION IF NOT EXISTS "pgcrypto";




-- ENUM type for form_fields.field_type
CREATE TYPE form_field_type_enum AS ENUM (
    'text', 'email', 'number', 'password', 'date',
    'checkbox', 'dropdown', 'checkboxgroup', 'radio', 'radiogroup',
    'textarea', 'textlist', 'file'
);

-- Table: forms
CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE, 
    description TEXT,
    url_send VARCHAR(2048),
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indexes for forms
CREATE INDEX idx_forms_name ON forms(name);
CREATE INDEX idx_forms_slug ON forms(slug); 
CREATE INDEX idx_forms_created_by_user_id ON forms(created_by_user_id);
CREATE INDEX idx_forms_deleted_at ON forms(deleted_at) WHERE deleted_at IS NOT NULL;

-- Table: form_sections
CREATE TABLE form_sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    title VARCHAR(255),
    description TEXT,
    order_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for form_sections
CREATE INDEX idx_form_sections_form_id ON form_sections(form_id);
CREATE INDEX idx_form_sections_form_id_order_index ON form_sections(form_id, order_index);

-- Table: form_fields
CREATE TABLE form_fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    section_id UUID NOT NULL REFERENCES form_sections(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    field_key VARCHAR(100) NOT NULL,
    field_type form_field_type_enum NOT NULL,
    placeholder VARCHAR(255),
    is_required BOOLEAN NOT NULL DEFAULT FALSE,
    default_value TEXT,
    properties JSONB, 
    display JSONB,   
    order_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for form_fields
CREATE INDEX idx_form_fields_section_id ON form_fields(section_id);
CREATE INDEX idx_form_fields_field_key ON form_fields(field_key);
CREATE INDEX idx_form_fields_field_type ON form_fields(field_type);

-- Table: form_submissions
CREATE TABLE form_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE RESTRICT,
    submitted_data JSONB NOT NULL,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    submitter_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for form_submissions
CREATE INDEX idx_form_submissions_form_id ON form_submissions(form_id);
CREATE INDEX idx_form_submissions_submitted_at ON form_submissions(submitted_at DESC);
CREATE INDEX idx_gin_form_submissions_submitted_data ON form_submissions USING GIN (submitted_data);
CREATE INDEX idx_form_submissions_submitter_user_id ON form_submissions(submitter_user_id);

-- Create triggers to automatically update the updated_at timestamp
CREATE TRIGGER set_timestamp_forms
BEFORE UPDATE ON forms
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_form_sections
BEFORE UPDATE ON form_sections
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_form_fields
BEFORE UPDATE ON form_fields
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_form_submissions
BEFORE UPDATE ON form_submissions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();