package tests

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models/enums"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/MadManJJ/cms-api/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCMSRepo_CreateCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryRepo := repo.NewCMSCategoryRepository(gormDB)

	_, mockCategory := helpers.InitializeMockCategory() 
	categoryId := uuid.New()
	categoryTypeId := uuid.New()

	t.Run("successfully create category", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "category_types"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)				

		mock.ExpectQuery(`INSERT INTO "categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryId),
			)	

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "category_type_id"}).
					AddRow(categoryId, categoryTypeId),
			)		
			
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)				

		actualCategory, err := cmsCategoryRepo.CreateCategory(mockCategory)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategory)
		assert.NotNil(t, actualCategory.CategoryType)
	})

	t.Run("failed to create category", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "category_types"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)				

		mock.ExpectQuery(`INSERT INTO "categories"`).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()			

		actualCategory, err := cmsCategoryRepo.CreateCategory(mockCategory)
		assert.Error(t, err)
		assert.Nil(t, actualCategory)
	})	
}

func TestCMSRepo_GetCategoryByID(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryRepo := repo.NewCMSCategoryRepository(gormDB)
	
	categoryId := uuid.New()
	categoryTypeId := uuid.New()

	t.Run("successfully get category", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "category_type_id"}).
					AddRow(categoryId, categoryTypeId),
			)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)

		actualCategory, err := cmsCategoryRepo.GetCategoryByID(categoryId)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategory)
		assert.NotNil(t, actualCategory.CategoryType)
	})

	t.Run("failed to get category - not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		actualCategory, err := cmsCategoryRepo.GetCategoryByID(categoryId)
		assert.Error(t, err)
		assert.Nil(t, actualCategory)
	})
}

func TestCMSRepo_ListCategoriesByFilter(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryRepo := repo.NewCMSCategoryRepository(gormDB)
	
	categoryId := uuid.New()
	categoryTypeId := uuid.New()

	t.Run("successfully list categories with filter", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "category_type_id"}).
					AddRow(categoryId, categoryTypeId),
			)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)

		categoryTypeIdStr := categoryTypeId.String()
		language := enums.PageLanguageEN
		categoryName := "random name"
		publishStatus := enums.PublishStatusPublished
		filters := dto.CategoryFilter{
			CategoryTypeID: &categoryTypeIdStr,
			LanguageCode:   &language,
			Name:           &categoryName,
			PublishStatus:  &publishStatus,
		}

		categories, err := cmsCategoryRepo.ListCategoriesByFilter(filters)
		assert.NoError(t, err)
		assert.NotNil(t, categories)
	})

	t.Run("failed to list categories", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnError(errs.ErrInternalServerError)

		filters := dto.CategoryFilter{}
		categories, err := cmsCategoryRepo.ListCategoriesByFilter(filters)
		assert.Error(t, err)
		assert.Nil(t, categories)
	})
}

func TestCMSRepo_DeleteCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryRepo := repo.NewCMSCategoryRepository(gormDB)
	
	categoryId := uuid.New()

	t.Run("successfully delete category", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := cmsCategoryRepo.DeleteCategory(categoryId)
		assert.NoError(t, err)
	})

	t.Run("failed to delete category - not found", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories"`)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		err := cmsCategoryRepo.DeleteCategory(categoryId)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("failed to delete category: internal server error", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsCategoryRepo.DeleteCategory(categoryId)
		assert.Error(t, err)
		assert.NotNil(t, err)
	})	
}

func TestCMSRepo_UpdateCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryRepo := repo.NewCMSCategoryRepository(gormDB)
	
	_, mockCategory := helpers.InitializeMockCategory() 
	categoryId := uuid.New()
	categoryTypeId := uuid.New()

	t.Run("successfully update category", func(t *testing.T) {
		mockCategory.ID = categoryId
		mockCategory.CategoryTypeID = categoryTypeId

		mock.ExpectBegin()
		
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "category_type_id"}).
					AddRow(categoryId, categoryTypeId),
			)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)		
		
		actualCategory, err := cmsCategoryRepo.UpdateCategory(mockCategory)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategory)
		assert.NotNil(t, actualCategory.CategoryType)
	})

	t.Run("failed to update category", func(t *testing.T) {
		mockCategory.ID = categoryId
		mockCategory.CategoryTypeID = categoryTypeId
		
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories"`)).
			WillReturnError(errs.ErrInternalServerError)

		actualCategory, err := cmsCategoryRepo.UpdateCategory(mockCategory)
		assert.Error(t, err)
		assert.Nil(t, actualCategory)
	})
}

func TestCMSRepo_CountCategoriesByTypeAndLanguage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryRepo := repo.NewCMSCategoryRepository(gormDB)
	
	categoryTypeId := uuid.New()

	t.Run("successfully count categories", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT language_code, COUNT(*)`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"language_code", "count"}).
					AddRow("en", 2).
					AddRow("th", 3),
			)

		counts, err := cmsCategoryRepo.CountCategoriesByTypeAndLanguage(categoryTypeId)
		assert.NoError(t, err)
		assert.NotNil(t, counts)
		assert.Equal(t, 2, counts["en"])
		assert.Equal(t, 3, counts["th"])
	})

	t.Run("failed to count categories", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT language_code, COUNT(*)`)).
			WillReturnError(errs.ErrInternalServerError)

		counts, err := cmsCategoryRepo.CountCategoriesByTypeAndLanguage(categoryTypeId)
		assert.Error(t, err)
		assert.Nil(t, counts)
	})
}