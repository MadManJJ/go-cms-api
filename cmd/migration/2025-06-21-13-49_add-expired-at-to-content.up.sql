ALTER TABLE faq_contents
ADD COLUMN expired_at TIMESTAMP;

ALTER TABLE partner_contents
ADD COLUMN expired_at TIMESTAMP;

ALTER TABLE landing_contents
ADD COLUMN expired_at TIMESTAMP;
