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

func TestCMSRepo_CreateEmailCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	emailCategoryRepo := repo.NewEmailCategoryRepository(gormDB)

	
	t.Run("successfully create email category", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "email_categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)			
		mock.ExpectCommit()

		actualEmailCategory, err := emailCategoryRepo.Create(mockEmailCategory)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("failed to create email category", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "email_categories"`).
			WillReturnError(errs.ErrInternalServerError)			
		mock.ExpectRollback()

		actualEmailCategory, err := emailCategoryRepo.Create(mockEmailCategory)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_FindEmailCategoryById(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	emailCategoryRepo := repo.NewEmailCategoryRepository(gormDB)

	emailCategoryId := uuid.New()
		
	t.Run("successfully find email category by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE id = $1 ORDER BY "email_categories"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(emailCategoryId))

		actualEmailCategory, err := emailCategoryRepo.FindByID(emailCategoryId)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("failed to find email category by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE id = $1 ORDER BY "email_categories"."id" LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)

		actualEmailCategory, err := emailCategoryRepo.FindByID(emailCategoryId)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_FindEmailCategoryByTitle(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	emailCategoryRepo := repo.NewEmailCategoryRepository(gormDB)

	emailCategoryTitle := "Welcome Emails"
		
	t.Run("successfully find email category by title", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE title = $1 ORDER BY "email_categories"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()))

		actualEmailCategory, err := emailCategoryRepo.FindByTitle(emailCategoryTitle)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("failed to find email category by title", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE title = $1 ORDER BY "email_categories"."id" LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)

		actualEmailCategory, err := emailCategoryRepo.FindByTitle(emailCategoryTitle)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_FindAllEmailCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	emailCategoryRepo := repo.NewEmailCategoryRepository(gormDB)
		
	t.Run("successfully find all email categories", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()).
				AddRow(uuid.New()))

		actualEmailCategories, err := emailCategoryRepo.FindAll()
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategories)
		assert.Len(t, actualEmailCategories, 2)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("failed to find all email categories", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories"`)).
			WillReturnError(errs.ErrInternalServerError)

		actualEmailCategories, err := emailCategoryRepo.FindAll()
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategories)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_UpdateEmailCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	emailCategoryRepo := repo.NewEmailCategoryRepository(gormDB)

	t.Run("successfully update email category", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "email_categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)						
		mock.ExpectCommit()

		actualEmailCategory, err := emailCategoryRepo.Update(mockEmailCategory)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("failed to update email category", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "email_categories"`)).
			WillReturnError(errs.ErrInternalServerError)						
		mock.ExpectRollback()

		actualEmailCategory, err := emailCategoryRepo.Update(mockEmailCategory)
		assert.Error(t, err)
		assert.Nil(t, actualEmailCategory)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})
}

func TestCMSRepo_DeleteEmailCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	emailCategoryRepo := repo.NewEmailCategoryRepository(gormDB)
	
	t.Run("successfully delete email category", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		mockEmailCategory.ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_contents" WHERE email_category_id = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).
					AddRow(0), 
			)								
		
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_categories" WHERE id = $1`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			

		mock.ExpectCommit()

		err := emailCategoryRepo.Delete(mockEmailCategory.ID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("failed to delete email category: email category is in use", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		mockEmailCategory.ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_contents" WHERE email_category_id = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).
					AddRow(1), 
			)					

		err := emailCategoryRepo.Delete(mockEmailCategory.ID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("failed to delete email category: internal server error", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		mockEmailCategory.ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_contents" WHERE email_category_id = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).
					AddRow(0), 
			)								
		
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_categories" WHERE id = $1`)).
			WillReturnError(errs.ErrInternalServerError)			

		mock.ExpectRollback()

		err := emailCategoryRepo.Delete(mockEmailCategory.ID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	
}

func TestCMSRepo_IsEmailCategoryTitleUnique(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	emailCategoryRepo := repo.NewEmailCategoryRepository(gormDB)	

	t.Run("successfully determine title is unique (with exclude id)", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		mockEmailCategory.ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_categories" WHERE title = $1 AND id != $2`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).
					AddRow(0), 
			)								

		actualIsUnique, err := emailCategoryRepo.IsTitleUnique(mockEmailCategory.Title, mockEmailCategory.ID)
		assert.NoError(t, err)
		assert.True(t, actualIsUnique)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("successfully determine title is unique (without exclude id)", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		mockEmailCategory.ID = uuid.Nil

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_categories" WHERE title = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).
					AddRow(0), 
			)								

		actualIsUnique, err := emailCategoryRepo.IsTitleUnique(mockEmailCategory.Title, mockEmailCategory.ID)
		assert.NoError(t, err)
		assert.True(t, actualIsUnique)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	

	t.Run("successfully determine title isn't unique", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		mockEmailCategory.ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_categories" WHERE title = $1 AND id != $2`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"count"}).
					AddRow(1), 
			)								

		actualIsUnique, err := emailCategoryRepo.IsTitleUnique(mockEmailCategory.Title, mockEmailCategory.ID)
		assert.NoError(t, err)
		assert.False(t, actualIsUnique)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("failed to determine title is unique", func(t *testing.T) {
		mockEmailCategory := helpers.InitializeMockEmailCategory()
		mockEmailCategory.ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_categories" WHERE title = $1 AND id != $2`)).
			WillReturnError(errs.ErrInternalServerError)				

		actualIsUnique, err := emailCategoryRepo.IsTitleUnique(mockEmailCategory.Title, mockEmailCategory.ID)
		assert.Error(t, err)
		assert.False(t, actualIsUnique)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	
}