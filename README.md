# CMS API

A Go REST API built using Clean Architecture principles with Fiber framework.

## Project Structure

The project follows Clean Architecture principles with clear separation of concerns across multiple layers. Here's a detailed breakdown of the structure:

### Core Layers

- **DTO (Data Transfer Objects)**: Defines the data structures used for API requests/responses
- **Repository Layer**: Handles data access and implements database operations
- **Service Layer**: Contains business logic and use cases
- **Handler Layer**: Manages HTTP requests and responses
- **Domain Layer**: Contains core business models and interfaces

### Directory Structure

```
.
├── cmd/                  # Command-line applications
│   └── migration/        # Database migration scripts
│
├── config/              # Configuration files
├── dto/                 # Data Transfer Objects
│   ├── app/             # DTOs for app domain
│   └── cms/             # DTOs for CMS domain
│
├── handlers/            # HTTP request handlers
│   ├── app/             # Handlers for app domain
│   ├── cms/             # Handlers for CMS domain
│   └── common/          # Common handler utilities
│
├── helpers/             # Helper functions and utilities
├── middleware/          # HTTP middleware components
├── models/              # Database models and enums
│   └── enums/           # Enumerations used across the application
│
├── repositories/        # Data access layer
│   ├── app/             # Repositories for app domain
│   └── cms/             # Repositories for CMS domain (implements CMSAuthRepository)
│
├── services/            # Business logic layer
│   ├── app/             # Services for app domain
│   └── cms/             # Services for CMS domain
│
├── tests/               # Test suites
│   ├── integration/     # Integration tests
│   │   └── cms/         # CMS integration tests
│   └── unit/            # Unit tests
│       ├── app/         # App domain unit tests
│       ├── cms/         # CMS domain unit tests
│       └── common/      # Common utilities unit tests
│
├── .env.example         # Environment variables example
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

### Key Components

- **cmd/migration/**: Contains database migration scripts
- **config/**: Centralized configuration management
- **dto/**: Data Transfer Objects for both app and CMS domains
- **handlers/**: HTTP request handlers organized by domain
- **models/**: Database models and enumerations
- **repositories/**: Data access layer implementing the repository pattern
  - **cms/**: Implements `CMSAuthRepository` for CMS authentication
- **services/**: Business logic layer
- **tests/**: Comprehensive test suites including unit and integration tests

## Environment Variables

The application uses environment variables for configuration. A `.env.example` file is provided as a template. To get started:

1. Copy `.env.example` to a new file called `.env`
2. Update the values according to your environment

### Key Environment Variables

#### Server Configuration

- `PORT` - Port number the server will run on (default: 8080)
- `ENVIRONMENT` - Application environment (development, production, etc.)
- `APP_NAME` - Name of the application

#### Database

- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USERNAME` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name

#### URLs

- `FRONTEND_URLS` - Comma-separated list of allowed frontend URLs for CORS
- `API_BASE_URL` - Base URL of the API
- `CMS_BASE_URL` - Base URL for the CMS
- `WEB_BASE_URL` - Base URL for the web frontend

#### File Storage

- `UPLOAD_FILE_PATH` - Path where uploaded files will be stored
- `STATIC_FILE_PREFIX` - URL prefix for serving static files

#### Authentication

- `JWT_SECRET_KEY` - Secret key for JWT token generation and validation

#### Email Service (SendGrid)

- `SENDGRID_API_KEY` - API key for SendGrid email service

#### LINE Login

- `OAUTH_CLIENT_ID` - LINE Login client ID
- `OAUTH_REDIRECT_URI` - Callback URL for LINE Login
- `OAUTH_CLIENT_SECRET` - Client secret for LINE Login
- `TOKEN_URL` - LINE OAuth token URL
- `AUTHORIZE_URL` - LINE OAuth authorization URL

#### Development Tools

- `PGADMIN_DEFAULT_EMAIL` - Email for pgAdmin (development only)
- `PGADMIN_DEFAULT_PASSWORD` - Password for pgAdmin (development only)

> **Note**: Never commit your `.env` file to version control. It's included in `.gitignore` by default.

## API Endpoints

### Health Check

- GET `/health` - Health check endpoint

### App Domain

#### Test Endpoints

- GET `/api/v1/app/test` - Test endpoint
- GET `/api/v1/app/additional` - Additional test endpoint

#### Landing Pages

- GET `/api/v1/app/landingpages/:languageCode/by-alias` - Get landing page by URL alias
- GET `/api/v1/app/landingpages/previews/:id` - Get landing page preview

#### Partner Pages

- GET `/api/v1/app/partnerpages/:languageCode/by-alias` - Get partner page by alias
- GET `/api/v1/app/partnerpages/:languageCode/by-url` - Get partner page by URL
- GET `/api/v1/app/partnerpages/previews/:id` - Get partner page preview

#### FAQ Pages

- GET `/api/v1/app/faqpages/:languageCode/by-alias` - Get FAQ page by alias
- GET `/api/v1/app/faqpages/:languageCode/by-url` - Get FAQ page by URL
- GET `/api/v1/app/faqpages/previews/:id` - Get FAQ page preview

#### Forms

- GET `/api/v1/app/forms/:formId/structure` - Get form structure

### CMS Domain

#### Authentication

- POST `/api/v1/cms/auth/register` - Register a new user
- POST `/api/v1/cms/auth/login` - User login

#### FAQ Pages Management

- POST `/api/v1/cms/faqpages` - Create new FAQ page
- GET `/api/v1/cms/faqpages` - List all FAQ pages
- GET `/api/v1/cms/faqpages/:pageId` - Get FAQ page by ID
- GET `/api/v1/cms/faqpages/:pageId/latestcontents/:languageCode` - Get latest FAQ content
- POST `/api/v1/cms/faqpages/:revisionId/revisions` - Revert FAQ content
- PUT `/api/v1/cms/faqpages/:contentId/contents` - Update FAQ content
- DELETE `/api/v1/cms/faqpages/:pageId` - Delete FAQ page
- GET `/api/v1/cms/faqpages/:pageId/contents/:languageCode` - Get FAQ content by language
- DELETE `/api/v1/cms/faqpages/:pageId/contents/:languageCode` - Delete FAQ content
- POST `/api/v1/cms/faqpages/duplicate/:pageId/pages` - Duplicate FAQ page
- POST `/api/v1/cms/faqpages/duplicate/:contentId/contents` - Duplicate FAQ content to another language
- GET `/api/v1/cms/faqpages/category/:categoryTypeCode/:pageId/:languageCode` - Get FAQ category
- GET `/api/v1/cms/faqpages/revisions/:languageCode/:pageId` - Get FAQ revisions
- POST `/api/v1/cms/faqpages/previews/:pageId` - Preview FAQ content

#### Landing Pages Management

- POST `/api/v1/cms/landingpages` - Create new landing page
- GET `/api/v1/cms/landingpages` - List all landing pages
- GET `/api/v1/cms/landingpages/:pageId` - Get landing page by ID
- GET `/api/v1/cms/landingpages/:pageId/latestcontents/:languageCode` - Get latest landing page content
- POST `/api/v1/cms/landingpages/:revisionId/revisions` - Revert landing page content
- PUT `/api/v1/cms/landingpages/:contentId/contents` - Update landing page content
- DELETE `/api/v1/cms/landingpages/:pageId` - Delete landing page
- GET `/api/v1/cms/landingpages/:pageId/contents/:languageCode` - Get landing page content by language
- DELETE `/api/v1/cms/landingpages/:pageId/contents/:languageCode` - Delete landing page content
- POST `/api/v1/cms/landingpages/duplicate/:pageId/pages` - Duplicate landing page
- POST `/api/v1/cms/landingpages/duplicate/:contentId/contents` - Duplicate landing page content to another language
- GET `/api/v1/cms/landingpages/category/:categoryTypeCode/:pageId/:languageCode` - Get landing page category
- GET `/api/v1/cms/landingpages/revisions/:languageCode/:pageId` - Get landing page revisions
- POST `/api/v1/cms/landingpages/previews/:pageId` - Preview landing page content

#### Partner Pages Management

- POST `/api/v1/cms/partnerpages` - Create new partner page
- GET `/api/v1/cms/partnerpages` - List all partner pages
- GET `/api/v1/cms/partnerpages/:pageId` - Get partner page by ID
- GET `/api/v1/cms/partnerpages/:pageId/latestcontents/:languageCode` - Get latest partner page content
- POST `/api/v1/cms/partnerpages/:revisionId/revisions` - Revert partner page content
- PUT `/api/v1/cms/partnerpages/:contentId/contents` - Update partner page content
- DELETE `/api/v1/cms/partnerpages/:pageId` - Delete partner page
- GET `/api/v1/cms/partnerpages/:pageId/contents/:languageCode` - Get partner page content by language
- DELETE `/api/v1/cms/partnerpages/:pageId/contents/:languageCode` - Delete partner page content
- POST `/api/v1/cms/partnerpages/duplicate/:pageId/pages` - Duplicate partner page
- POST `/api/v1/cms/partnerpages/duplicate/:contentId/contents` - Duplicate partner page content to another language
- GET `/api/v1/cms/partnerpages/category/:categoryTypeCode/:pageId/:languageCode` - Get partner page category
- GET `/api/v1/cms/partnerpages/revisions/:languageCode/:pageId` - Get partner page revisions
- POST `/api/v1/cms/partnerpages/previews/:pageId` - Preview partner page content

#### Category Types Management

- POST `/api/v1/cms/category-types` - Create category type
- GET `/api/v1/cms/category-types` - List all category types
- GET `/api/v1/cms/category-types/:id` - Get category type by ID
- PATCH `/api/v1/cms/category-types/:id` - Update category type
- DELETE `/api/v1/cms/category-types/:id` - Delete category type
- GET `/api/v1/cms/category-types/:categoryTypeId/categories` - List categories for type

#### Categories Management

- POST `/api/v1/cms/categories` - Create category
- GET `/api/v1/cms/categories` - List all categories
- GET `/api/v1/cms/categories/:categoryUuid` - Get category by UUID
- PATCH `/api/v1/cms/categories/:categoryUuid` - Update category
- DELETE `/api/v1/cms/categories/:categoryUuid` - Delete category

#### Email Management

- POST `/api/v1/cms/email-categories` - Create email category
- GET `/api/v1/cms/email-categories` - List email categories
- GET `/api/v1/cms/email-categories/:id` - Get email category
- PATCH `/api/v1/cms/email-categories/:id` - Update email category
- DELETE `/api/v1/cms/email-categories/:id` - Delete email category

#### Email Contents

- POST `/api/v1/cms/email-contents` - Create email content
- GET `/api/v1/cms/email-contents` - List email contents
- GET `/api/v1/cms/email-contents/category/:email_category_id/language/:language` - Get content by category and language
- PATCH `/api/v1/cms/email-contents/:id` - Update email content
- DELETE `/api/v1/cms/email-contents/:id` - Delete email content

#### Media Files

- POST `/api/v1/cms/media-files` - Upload media file
- GET `/api/v1/cms/media-files` - List media files
- GET `/api/v1/cms/media-files/:id` - Get media file by ID
- DELETE `/api/v1/cms/media-files/:id` - Delete media file

#### Form Builder

- POST `/api/v1/cms/forms` - Create form
- GET `/api/v1/cms/forms` - List all forms
- GET `/api/v1/cms/forms/:formId` - Get form by ID
- PUT `/api/v1/cms/forms/:formId` - Update form
- DELETE `/api/v1/cms/forms/:formId` - Delete form
- POST `/api/v1/cms/forms/:formId/submissions` - Submit form data
- GET `/api/v1/cms/forms/:formId/submissions` - List form submissions
- GET `/api/v1/cms/forms/submissions/:submissionId` - Get form submission

### Common Endpoints

#### Email Sending

- POST `/api/v1/emails/send` - Send email

#### LINE Login

- GET `/api/v1/common/login-link` - Get LINE login link
- POST `/api/v1/common/authenticate` - Authenticate with LINE
- POST `/api/v1/common/refresh-token` - Refresh authentication token

#### Middleware Testing

- GET `/api/v1/middleware/test` - Test middleware (requires authentication)

## Getting started

### Prerequisites

- Go 1.18 or higher
- Git

### Installation

1. Clone the repository

```bash
git clone https://github.com/MadManJJ/go-cms-api.git
cd api
```

2. Install dependencies

```bash
go mod tidy
```

### Run the server

```bash
go run cmd/api/main.go
```

The server will start on port 8080. You can test the health check endpoint with:

```bash
curl http://localhost:8080/health
```

Expected output: `OK`

## Add your files

- [ ] [Create](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#create-a-file) or [upload](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#upload-a-file) files
- [ ] [Add files using the command line](https://docs.gitlab.com/ee/gitlab-basics/add-file.html#add-a-file-using-the-command-line) or push an existing Git repository with the following command:

```
cd existing_repo
git remote add origin https://github.com/MadManJJ/go-cms-api.git
git branch -M main
git push -uf origin main
```

## Test and Deploy

Command for run test

```shell
go test ./...
```
