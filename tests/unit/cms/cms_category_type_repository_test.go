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

func TestCMSRepo_CreateCategoryType(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryTypeRepo := repo.NewCMSCategoryTypeRepository(gormDB)

	mockCategoryType, _ := helpers.InitializeMockCategory()
	
	t.Run("successfully create category type", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()),
			)

		mock.ExpectCommit()
		
		actualCategoryTpye, err := cmsCategoryTypeRepo.Create(mockCategoryType)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategoryTpye)
	})

	t.Run("failed to create category type", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()
		
		actualCategoryTpye, err := cmsCategoryTypeRepo.Create(mockCategoryType)
		assert.Error(t, err)
		assert.Nil(t, actualCategoryTpye)
	})	
}

func TestCMSRepo_FindCategoryTypeByID(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryTypeRepo := repo.NewCMSCategoryTypeRepository(gormDB)

	categoryTypeId := uuid.New()
	
	t.Run("successfully get category", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)

		actualCategory, err := cmsCategoryTypeRepo.FindByID(categoryTypeId)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategory)
	})	

	t.Run("failed to get category", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnError(errs.ErrInternalServerError)

		actualCategory, err := cmsCategoryTypeRepo.FindByID(categoryTypeId)
		assert.Error(t, err)
		assert.Nil(t, actualCategory)
	})
}

func TestCMSRepo_FindCategoryTypeByCode(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryTypeRepo := repo.NewCMSCategoryTypeRepository(gormDB)

	categoryTypeId := uuid.New()
	code := "some random code"
	
	t.Run("successfully get category", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE type_code = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(categoryTypeId),
			)

		actualCategory, err := cmsCategoryTypeRepo.FindByCode(code)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategory)
	})	

	t.Run("failed to get category", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE type_code = $1`)).
			WillReturnError(errs.ErrInternalServerError)

		actualCategory, err := cmsCategoryTypeRepo.FindByCode(code)
		assert.Error(t, err)
		assert.Nil(t, actualCategory)
	})
}

func TestCMSRepo_FindAllCategoryTypes(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryTypeRepo := repo.NewCMSCategoryTypeRepository(gormDB)

	categoryTypeId := uuid.New()
	isActive := true
	
	t.Run("successfully get all active category types", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE is_active = $1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
					AddRow(categoryTypeId, "CODE1", "Category 1", true).
					AddRow(uuid.New(), "CODE2", "Category 2", true),
			)

		actualCategories, err := cmsCategoryTypeRepo.FindAll(&isActive)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategories)
		assert.Equal(t, 2, len(actualCategories))
	})

	t.Run("successfully get all category types without filter", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
					AddRow(categoryTypeId, "CODE1", "Category 1", true).
					AddRow(uuid.New(), "CODE2", "Category 2", false),
			)

		actualCategories, err := cmsCategoryTypeRepo.FindAll(nil)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategories)
		assert.Equal(t, 2, len(actualCategories))
	})

	t.Run("failed to get category types", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE is_active = $1`)).
			WillReturnError(errs.ErrInternalServerError)

		actualCategories, err := cmsCategoryTypeRepo.FindAll(&isActive)
		assert.Error(t, err)
		assert.Nil(t, actualCategories)
	})
}

func TestCMSRepo_UpdateCategoryType(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryTypeRepo := repo.NewCMSCategoryTypeRepository(gormDB)

	categoryTypeId := uuid.New()
	
	t.Run("successfully update category type", func(t *testing.T) {
		// First mock the SELECT query to find the category type
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
				AddRow(categoryTypeId, "CODE1", "Old Name", true))

		mock.ExpectBegin()

		// Then mock the UPDATE query
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "category_types" SET `)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		updates := map[string]interface{}{
			"name": "Updated Name",
			"is_active": false,
		}

		actualCategoryType, err := cmsCategoryTypeRepo.Update(categoryTypeId, updates)
		assert.NoError(t, err)
		assert.NotNil(t, actualCategoryType)
		assert.Equal(t, "Updated Name", actualCategoryType.Name)
		assert.Equal(t, false, actualCategoryType.IsActive)
	})

	t.Run("failed to update category type: not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnError(errs.ErrNotFound)

		updates := map[string]interface{}{
			"name": "Updated Name",
		}

		actualCategoryType, err := cmsCategoryTypeRepo.Update(categoryTypeId, updates)
		assert.Error(t, err)
		assert.Nil(t, actualCategoryType)
	})

	t.Run("failed to update category type: update error", func(t *testing.T) {
		// First mock the SELECT query to find the category type
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
				AddRow(categoryTypeId, "CODE1", "Old Name", true))

		mock.ExpectBegin()		

		// Then mock the UPDATE query with error
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "category_types" SET `)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		updates := map[string]interface{}{
			"name": "Updated Name",
		}

		actualCategoryType, err := cmsCategoryTypeRepo.Update(categoryTypeId, updates)
		assert.Error(t, err)
		assert.Nil(t, actualCategoryType)
	})
}

func TestCMSRepo_DeleteCategoryType(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsCategoryTypeRepo := repo.NewCMSCategoryTypeRepository(gormDB)

	categoryTypeId := uuid.New()
	
	t.Run("successfully delete category type", func(t *testing.T) {
		// First mock the SELECT query to find the category type
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
				AddRow(categoryTypeId, "CODE1", "Category 1", true))

		mock.ExpectBegin()

		// Then mock the DELETE query
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "category_types" WHERE id = $1`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := cmsCategoryTypeRepo.Delete(categoryTypeId)
		assert.NoError(t, err)
	})

	t.Run("failed to delete category type: not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnError(errs.ErrNotFound)

		err := cmsCategoryTypeRepo.Delete(categoryTypeId)
		assert.Error(t, err)
	})

	t.Run("failed to delete category type: delete error", func(t *testing.T) {
		// First mock the SELECT query to find the category type
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
				AddRow(categoryTypeId, "CODE1", "Category 1", true))

		mock.ExpectBegin()

		// Then mock the DELETE query with error
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "category_types" WHERE id = $1`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsCategoryTypeRepo.Delete(categoryTypeId)
		assert.Error(t, err)
	})
}

