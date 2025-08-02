-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enums
CREATE TYPE page_language AS ENUM ('th', 'en');
CREATE TYPE page_mode AS ENUM ('Published', 'Preview', 'Histories', 'Draft');
CREATE TYPE publish_status AS ENUM ('UnPublished', 'Published');
CREATE TYPE workflow_status AS ENUM ('Draft', 'Approval_Pending', 'Waiting_Design_Approved', 'Schedule', 'Published', 'UnPublished', 'Waiting_Deletion', 'Delete');
CREATE TYPE file_type AS ENUM ('CSS', 'JS');
CREATE TYPE component_type AS ENUM (
	'Margin',
	'LeadComponent',
	'StickyHeader',
	'SectionContent',
	'Divider',
	'DynamicClassTextSection',
	'VideoWithEditor',
	'TC0101',
	'TC0102',
	'CoachProfileList',
	'X1004',
	'H21',
	'H22',
	'H23',
	'H24',
	'H31',
	'H32',
	'LargeGreenLinkButton',
	'ThreeLargeGreenLinkButton',
	'LargeWhiteLinkButton',
	'MidsizeWhiteLinkButtonLeftAligned',
	'MidsizeWhiteLinkButtonCentered',
	'MidsizeWhiteLinkButtonRightAligned',
	'List',
	'Box',
	'Quotation',
	'RelatedLinks',
	'RelatedArticles',
	'BL0501',
	'OneColumnImage',
	'TwoColumnImage',
	'Video',
	'VideoExternalLink',
	'TabContent',
	'L0201',
	'L0301',
	'L0401',
	'L0501',
	'L0601',
	'NormalText',
	'NormalTextRed',
	'Bold',
	'TextCentered',
	'Chatter',
	'TextList',
	'TextListNumber',
	'Notes',
	'Links',
	'LinksSeparateWindow',
	'AnchorLink',
	'Pdf',
	'U0201',
	'QrCode',
	'X0201',
	'X0301',
	'X0302',
	'X0401List',
	'X0501',
	'X0601',
	'X0701',
	'X0801',
	'X1001',
	'X1002',
	'X1101',
	'X1201',
	'X1301',
	'ImageOneColText',
	'ImageTwoColText',
	'ImageThreeColText',
	'ImageFourColText',
	'LeftImageRightTextWrapped',
	'Y0201',
	'Y0202',
	'Y0301',
	'Y0302',
	'ImageVideoCoverOneColText',
	'ImageVideoCoverTwoColText',
	'CampaignLfc',
	'GridContents',
	'GridContentsLfcFilter'
);

-- Create tables

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS category_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_type_id UUID NOT NULL,
    language_code page_language NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    weight INTEGER DEFAULT 0,
    publish_status publish_status DEFAULT 'UnPublished',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_type_id) REFERENCES category_types(id)
);

CREATE TABLE IF NOT EXISTS meta_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    description TEXT,
    cover_image VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS revisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    landing_content_id UUID,
    partner_content_id UUID,
    faq_content_id UUID,
    publish_status publish_status,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    author VARCHAR(255),
    message TEXT,
    description TEXT,
    FOREIGN KEY (landing_content_id) REFERENCES landing_contents(id),
    FOREIGN KEY (partner_content_id) REFERENCES partner_contents(id),
    FOREIGN KEY (faq_content_id) REFERENCES faq_contents(id)
);

CREATE TABLE IF NOT EXISTS faq_pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS faq_contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL,
    title VARCHAR(255),
    language page_language,
    authored_at TIMESTAMP WITH TIME ZONE NOT NULL,
    html_input TEXT,
    mode page_mode,
    workflow_status workflow_status,
    publish_status publish_status,
    publish_on TIMESTAMP WITH TIME ZONE,
    unpublish_on TIMESTAMP WITH TIME ZONE,
    authored_on TIMESTAMP WITH TIME ZONE,
    url_alias VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    meta_tag_id UUID UNIQUE,
    expired_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (page_id) REFERENCES faq_pages(id),
    FOREIGN KEY (meta_tag_id) REFERENCES meta_tags(id)
);

CREATE TABLE IF NOT EXISTS landing_pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS landing_contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL,
    title VARCHAR(255),
    language page_language,
    authored_at TIMESTAMP WITH TIME ZONE NOT NULL,
    html_input TEXT,
    mode page_mode,
    workflow_status workflow_status,
    publish_status publish_status,
    url_alias VARCHAR(255) NOT NULL UNIQUE,
    meta_tag_id UUID UNIQUE,
    publish_on TIMESTAMP WITH TIME ZONE,
    unpublish_on TIMESTAMP WITH TIME ZONE,
    authored_on TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approval_email TEXT[],
    FOREIGN KEY (page_id) REFERENCES landing_pages(id),
    FOREIGN KEY (meta_tag_id) REFERENCES meta_tags(id)
);

CREATE TABLE IF NOT EXISTS landing_content_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    landing_content_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    download_url VARCHAR(255) NOT NULL,
    file_type file_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (landing_content_id) REFERENCES landing_contents(id)
);

CREATE TABLE IF NOT EXISTS partner_pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS partner_contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL,
    title VARCHAR(255),
    thumbnail_image VARCHAR(255),
    thumbnail_alt_text VARCHAR(255),
    company_logo VARCHAR(255),
    company_alt_text VARCHAR(255),
    language page_language,
    company_name VARCHAR(255),
    company_detail TEXT,
    lead_body TEXT,
    challenges TEXT,
    solutions TEXT,
    results TEXT,
    authored_at TIMESTAMP WITH TIME ZONE NOT NULL,
    html_input TEXT,
    mode page_mode,
    workflow_status workflow_status,
    publish_status publish_status,
    url_alias VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    meta_tag_id UUID UNIQUE,
    is_recommended BOOLEAN DEFAULT false,
    publish_on TIMESTAMP WITH TIME ZONE,
    unpublish_on TIMESTAMP WITH TIME ZONE,
    authored_on TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approval_email TEXT[],
    FOREIGN KEY (page_id) REFERENCES partner_pages(id),
    FOREIGN KEY (meta_tag_id) REFERENCES meta_tags(id)
);

CREATE TABLE IF NOT EXISTS components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    landing_content_id UUID,
    partner_content_id UUID,
    faq_content_id UUID,
    type component_type,
    props JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (landing_content_id) REFERENCES landing_contents(id),
    FOREIGN KEY (partner_content_id) REFERENCES partner_contents(id),
    FOREIGN KEY (faq_content_id) REFERENCES faq_contents(id)
);

-- Create junction tables for many-to-many relationships
CREATE TABLE IF NOT EXISTS faq_content_categories (
    faq_content_id UUID NOT NULL,
    category_id UUID NOT NULL,
    PRIMARY KEY (faq_content_id, category_id),
    FOREIGN KEY (faq_content_id) REFERENCES faq_contents(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS landing_content_categories (
    landing_content_id UUID NOT NULL,
    category_id UUID NOT NULL,
    PRIMARY KEY (landing_content_id, category_id),
    FOREIGN KEY (landing_content_id) REFERENCES landing_contents(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS partner_content_categories (
    partner_content_id UUID NOT NULL,
    category_id UUID NOT NULL,
    PRIMARY KEY (partner_content_id, category_id),
    FOREIGN KEY (partner_content_id) REFERENCES partner_contents(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);