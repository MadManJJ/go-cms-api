ALTER TABLE forms ADD COLUMN email_category_id UUID NULL;
ALTER TABLE forms ADD CONSTRAINT fk_forms_email_category FOREIGN KEY (email_category_id) REFERENCES email_categories(id) ON DELETE SET NULL;