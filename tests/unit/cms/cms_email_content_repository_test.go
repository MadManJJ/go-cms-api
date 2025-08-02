package tests

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models/enums"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCMSRepo_CreateEmailContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)

	emailCategoryId := uuid.New()
	emailContentId := uuid.New()

	t.Run("successfully create email content", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()
		
		mock.ExpectBegin()	

		mock.ExpectQuery(`INSERT INTO "email_categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)			
			
		mock.ExpectQuery(`INSERT INTO "email_contents"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailContentId),
			)		
			
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label"}).
					AddRow(emailContentId, emailCategoryId, mockEmailContent.Label),
			)				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)					

		actualEmailContent, err := cmsEmailContentRepo.Create(mockEmailContent)
		assert.NoError(t, err)
		assert.Equal(t, mockEmailContent.Label, actualEmailContent.Label)
		assert.NotNil(t, actualEmailContent.EmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to create email content", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()
		
		mock.ExpectBegin()	

		mock.ExpectQuery(`INSERT INTO "email_categories"`).
			WillReturnError(errs.ErrInternalServerError)		
			
		mock.ExpectRollback()

		actualEmailContent, err := cmsEmailContentRepo.Create(mockEmailContent)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindEmailContentById(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)

	emailCategoryId := uuid.New()
	emailContentId := uuid.New()	
	
	t.Run("successfully find email content by id", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label"}).
					AddRow(emailContentId, emailCategoryId, mockEmailContent.Label),
			)				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)					

		actualEmailContent, err := cmsEmailContentRepo.FindByID(emailContentId)
		assert.NoError(t, err)
		assert.Equal(t, mockEmailContent.Label, actualEmailContent.Label)
		assert.NotNil(t, actualEmailContent.EmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to find email content by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		actualEmailContent, err := cmsEmailContentRepo.FindByID(emailContentId)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindByCategoryIDAndLanguageAndLabel(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)	

	language := enums.PageLanguageEN
	label := "some label"
	emailCategoryId := uuid.New()
	emailContentId := uuid.New()		

	t.Run("successfully find email content by category id, language, and label", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label", "language"}).
					AddRow(emailContentId, emailCategoryId, mockEmailContent.Label, language),
			)	
			
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)	
			
		actualEmailContent, err := cmsEmailContentRepo.FindByCategoryIDAndLanguageAndLabel(emailCategoryId, language, label)
		assert.NoError(t, err)
		assert.Equal(t, mockEmailContent.Label, actualEmailContent.Label)
		assert.NotNil(t, actualEmailContent.EmailCategory)		
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find email content by category id, language, and label", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents"`)).
			WillReturnError(errs.ErrInternalServerError)	
			
		actualEmailContent, err := cmsEmailContentRepo.FindByCategoryIDAndLanguageAndLabel(emailCategoryId, language, label)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContent)	
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_ListEmailContentByFilters(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)	
	
	emailContentId := uuid.New()
	emailContentId2 := uuid.New()
	emailCategoryId := uuid.New().String()
	language := enums.PageLanguageEN
	label := "some label"


	t.Run("successfully list email content with all filters", func(t *testing.T) {
		mockEmailContentFilter := dto.EmailContentFilter{
			EmailCategoryID: &emailCategoryId,
			Language: &language,
			Label: &label,
		}		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1 AND language = $2 AND label LIKE $3`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label", "language"}).
					AddRow(emailContentId, emailCategoryId, label, language).
					AddRow(emailContentId2, emailCategoryId, label, language),
			)	
			
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)			

		emailContents, err := cmsEmailContentRepo.ListByFilters(mockEmailContentFilter)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(emailContents))
		assert.NotNil(t, emailContents[0].EmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully list email content with category and language filters", func(t *testing.T) {
		mockEmailContentFilter := dto.EmailContentFilter{
			EmailCategoryID: &emailCategoryId,
			Language: &language,
		}		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1 AND language = $2`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label", "language"}).
					AddRow(emailContentId, emailCategoryId, label, language),
			)	
			
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)			

		emailContents, err := cmsEmailContentRepo.ListByFilters(mockEmailContentFilter)
		assert.NoError(t, err)
		assert.NotNil(t, emailContents[0].EmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("successfully list email content with category filter only", func(t *testing.T) {
		mockEmailContentFilter := dto.EmailContentFilter{
			EmailCategoryID: &emailCategoryId,
			Language: &language,
		}		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label", "language"}).
					AddRow(emailContentId, emailCategoryId, label, language),
			)	
			
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)			

		emailContents, err := cmsEmailContentRepo.ListByFilters(mockEmailContentFilter)
		assert.NoError(t, err)
		assert.NotNil(t, emailContents[0].EmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
	t.Run("successfully list email content with no filters", func(t *testing.T) {
		mockEmailContentFilter := dto.EmailContentFilter{}		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label", "language"}).
					AddRow(emailContentId, emailCategoryId, label, language),
			)	
			
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)			

		emailContents, err := cmsEmailContentRepo.ListByFilters(mockEmailContentFilter)
		assert.NoError(t, err)
		assert.NotNil(t, emailContents[0].EmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to list list email content with all filters", func(t *testing.T) {
		mockEmailContentFilter := dto.EmailContentFilter{
			EmailCategoryID: &emailCategoryId,
			Language: &language,
			Label: &label,
		}		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1 AND language = $2 AND label LIKE $3`)).
			WillReturnError(errs.ErrInternalServerError)	

		emailContents, err := cmsEmailContentRepo.ListByFilters(mockEmailContentFilter)
		assert.Error(t, err)
		assert.Nil(t, emailContents)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_UpdateEmailContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)		

	mockEmailContent := helpers.InitializeMockEmailContent()
	emailCategoryId := uuid.New()
	emailContentId := uuid.New()
	mockEmailContent.ID = emailContentId

	t.Run("successfully update email content", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "email_categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)

		mock.ExpectExec(`UPDATE "email_contents"`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "label"}).
					AddRow(emailContentId, emailCategoryId, mockEmailContent.Label),
			)				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(emailCategoryId),
			)			

		emailContent, err := cmsEmailContentRepo.Update(mockEmailContent)
		assert.NoError(t, err)
		assert.Equal(t, emailContent.EmailCategoryID, mockEmailContent.EmailCategoryID)
		assert.NotNil(t, emailContent.EmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to update email content", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "email_categories"`).
			WillReturnError(errs.ErrInternalServerError)
		
		mock.ExpectRollback()

		emailContent, err := cmsEmailContentRepo.Update(mockEmailContent)
		assert.Error(t, err)
		assert.Nil(t, emailContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_DeleteEmailContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)
	
	emailContentId := uuid.New()

	t.Run("successfully delete email content by email content id", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))	

		mock.ExpectCommit()

		err := cmsEmailContentRepo.Delete(emailContentId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to delete email content by email content id", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsEmailContentRepo.Delete(emailContentId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_DeleteEmailContentByCategoryID(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)
	
	categoryId := uuid.New()

	t.Run("successfully delete email content by category id", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))	

		mock.ExpectCommit()

		err := cmsEmailContentRepo.DeleteByCategoryID(categoryId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to delete email content by category id", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsEmailContentRepo.DeleteByCategoryID(categoryId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_FindEmailContentByCategoryIDAndLanguage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsEmailContentRepo := repo.NewEmailContentRepository(gormDB)
	categoryId := uuid.New()
	emailContentId := uuid.New()
	language := enums.PageLanguageEN

	t.Run("successfully find email content by category id and language", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1 AND language = $2`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email_category_id", "language"}).
					AddRow(emailContentId, emailContentId, language),
			)	

		_, err := cmsEmailContentRepo.FindEmailContentByCategoryIDAndLanguage(categoryId, language)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}