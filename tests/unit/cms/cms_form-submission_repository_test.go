package tests

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCMSRepo_CreateFormSubmission(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormSubmissionRepo := repo.NewFormSubmissionRepository(gormDB)

	t.Run("successfully create form submission", func(t *testing.T) {
		mockSubmission := helpers.InitializeMockFormSubmission()
		formId := uuid.New()
		formSubmissionId := uuid.New()

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "form_submissions"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_submissions" WHERE id = $1 AND "form_submissions"."id" = $2 ORDER BY "form_submissions"."id" LIMIT $3`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
				AddRow(formSubmissionId, formId))			

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE "forms"."id" = $1 AND "forms"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(formId))					

		formSubmission, err := cmsFormSubmissionRepo.CreateFormSubmission(mockSubmission)

		assert.NoError(t, err)
		assert.Equal(t, mockSubmission, formSubmission)
		assert.NotNil(t, formSubmission.Form)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to create form submission", func(t *testing.T) {
		mockSubmission := helpers.InitializeMockFormSubmission()

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "form_submissions"`).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		formSubmission, err := cmsFormSubmissionRepo.CreateFormSubmission(mockSubmission)

		assert.Error(t, err)
		assert.Nil(t, formSubmission)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSRepo_GetFormSubmissions(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormSubmissionRepo := repo.NewFormSubmissionRepository(gormDB)
	
	t.Run("successfully get form submissions", func(t *testing.T) {
		formSubmissionId := uuid.New()
		formId := uuid.New()
		sort := "created_at:DESC"
		page := 1
		limit := 10

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "form_submissions" WHERE form_id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))		
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_submissions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
				AddRow(formSubmissionId, formId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(formId))
		
		formSubmissions, totalCount, err := cmsFormSubmissionRepo.GetFormSubmissions(formId, sort, page, limit)
		
		assert.NoError(t, err)
		assert.NotNil(t, formSubmissions)
		assert.Equal(t, int64(1), totalCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get form submissions", func(t *testing.T) {
		formId := uuid.New()
		sort := "created_at:DESC"
		page := 1
		limit := 10
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "form_submissions" WHERE form_id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))		
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_submissions"`)).
			WillReturnError(errs.ErrInternalServerError)

		formSubmissions, totalCount, err := cmsFormSubmissionRepo.GetFormSubmissions(formId, sort, page, limit)
		
		assert.Error(t, err)
		assert.Nil(t, formSubmissions)
		assert.Equal(t, int64(0), totalCount)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_GetFormSubmission(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormSubmissionRepo := repo.NewFormSubmissionRepository(gormDB)
	
	t.Run("successfully get form submission", func(t *testing.T) {
		formSubmissionId := uuid.New()
		formId := uuid.New()
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_submissions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
				AddRow(formSubmissionId, formId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(formId))		

		formSubmission, err := cmsFormSubmissionRepo.GetFormSubmission(formSubmissionId)
		
		assert.NoError(t, err)
		assert.NotNil(t, formSubmission)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
	t.Run("failed to get form submission", func(t *testing.T) {
		formSubmissionId := uuid.New()
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_submissions"`)).
			WillReturnError(errs.ErrInternalServerError)

		formSubmission, err := cmsFormSubmissionRepo.GetFormSubmission(formSubmissionId)
		
		assert.Error(t, err)
		assert.Nil(t, formSubmission)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_GetEmailContentsFormFormId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormSubmissionRepo := repo.NewFormSubmissionRepository(gormDB)
	
	formId := uuid.New()
	emailCategoryId := uuid.New()
	emailContentId1 := uuid.New()
	emailContentId2 := uuid.New()

	t.Run("successfully get email contents form form id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "email_category_id" FROM "forms"`)).
			WillReturnRows(sqlmock.NewRows([]string{"email_category_id"}).
				AddRow(emailCategoryId))		

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email_category_id"}).
				AddRow(emailContentId1, emailCategoryId).
				AddRow(emailContentId2, emailCategoryId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE "email_categories"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(emailCategoryId))						

		emailContents, err := cmsFormSubmissionRepo.GetEmailContentsFormFormId(formId)
		
		assert.NoError(t, err)
		assert.NotNil(t, emailContents)
		assert.Equal(t, 2, len(emailContents))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get email contents form form id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "email_category_id" FROM "forms"`)).
			WillReturnError(errs.ErrInternalServerError)

		emailContents, err := cmsFormSubmissionRepo.GetEmailContentsFormFormId(formId)
		
		assert.Error(t, err)
		assert.Nil(t, emailContents)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}