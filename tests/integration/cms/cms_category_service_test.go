package integration

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"
	"github.com/MadManJJ/cms-api/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupCategoryTest(t *testing.T) (
	*gorm.DB,
	sqlmock.Sqlmock,
	services.CMSCategoryTypeServiceInterface,
	services.CMSCategoryServiceInterface,
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	categoryRepo := repositories.NewCMSCategoryRepository(gormDB)
	categoryTypeRepo := repositories.NewCMSCategoryTypeRepository(gormDB)
	categoryService := services.NewCMSCategoryService(categoryRepo, categoryTypeRepo)
	categoryTypeService := services.NewCMSCategoryTypeService(categoryTypeRepo, categoryRepo, categoryService)
	return gormDB, mock, categoryTypeService, categoryService, cleanup
}

func TestCategoryLifecycle(t *testing.T) {
	_, mock, categoryTypeService, categoryService, cleanup := setupCategoryTest(t)
	defer cleanup()

	var createdCategoryTypeID uuid.UUID
	var createdCategoryID uuid.UUID
	typeName := "FAQ Topics"
	typeCode := "FAQ_TOPICS"

	// --- 1. Create Category Type ---
	t.Run("1_CreateCategoryType", func(t *testing.T) {
		req := dto.CreateCategoryTypeRequest{TypeCode: typeCode, Name: &typeName}
		tempID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "category_types" WHERE type_code = $1`)).
			WithArgs(typeCode).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempID))
		mock.ExpectCommit()

		resp, err := categoryTypeService.CreateCategoryType(req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		createdCategoryTypeID, _ = uuid.Parse(resp.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 2. Create Category ---
	t.Run("2_CreateCategory", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryTypeID, "CategoryType must be created first")

		categoryName := "Technical Issues"
		req := dto.CategoryCreateRequest{
			CategoryTypeID: createdCategoryTypeID.String(),
			LanguageCode:   enums.PageLanguageEN,
			Name:           categoryName,
			PublishStatus:  enums.PublishStatusPublished,
		}

		mock.ExpectQuery(`SELECT \* FROM "category_types" WHERE id = \$1 ORDER BY "category_types"\."id" LIMIT \$2`).
			WithArgs(createdCategoryTypeID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
				AddRow(createdCategoryTypeID, typeCode, typeName, true))

		tempCreatedCategoryID := uuid.New()
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempCreatedCategoryID))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT \* FROM "categories" WHERE id = \$1 ORDER BY "categories"\."id" LIMIT \$2`).
			WithArgs(tempCreatedCategoryID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "category_type_id"}).
				AddRow(tempCreatedCategoryID, categoryName, createdCategoryTypeID))

		mock.ExpectQuery(`SELECT \* FROM "category_types" WHERE "category_types"\."id" = \$1`).
			WithArgs(createdCategoryTypeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type_code"}).
				AddRow(createdCategoryTypeID, typeCode))

		resp, err := categoryService.CreateCategory(req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		createdCategoryID, _ = uuid.Parse(resp.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 3. Get and Verify Category Count on Type ---
	t.Run("3_VerifyCategoryCountOnType", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryTypeID, "CategoryType must exist")

		mock.ExpectQuery(`SELECT \* FROM "category_types" WHERE id = \$1 ORDER BY "category_types"\."id" LIMIT \$2`).
			WithArgs(createdCategoryTypeID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdCategoryTypeID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT language_code, COUNT(*) as count FROM "categories" WHERE category_type_id = $1 GROUP BY "language_code"`)).
			WithArgs(createdCategoryTypeID).
			WillReturnRows(sqlmock.NewRows([]string{"language_code", "count"}).AddRow("en", 1))

		resp, err := categoryTypeService.GetCategoryTypeByID(createdCategoryTypeID.String())

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 1, (*resp.ChildrenCount)["en"])
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 4. Update Category ---
	t.Run("4_UpdateCategory", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryID, "Category must exist")

		updatedName := "Technical Problems"
		req := dto.CategoryUpdateRequest{Name: &updatedName}

		mock.ExpectQuery(`SELECT \* FROM "categories" WHERE id = \$1 ORDER BY "categories"\."id" LIMIT \$2`).WithArgs(createdCategoryID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).AddRow(createdCategoryID, createdCategoryTypeID))
		mock.ExpectQuery(`SELECT \* FROM "category_types"`).WithArgs(createdCategoryTypeID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdCategoryTypeID))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "categories"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT \* FROM "categories" WHERE id = \$1 ORDER BY "categories"\."id" LIMIT \$2`).WithArgs(createdCategoryID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "category_type_id"}).AddRow(createdCategoryID, updatedName, createdCategoryTypeID))
		mock.ExpectQuery(`SELECT \* FROM "category_types"`).WithArgs(createdCategoryTypeID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdCategoryTypeID))

		resp, err := categoryService.UpdateCategoryByUUID(createdCategoryID.String(), req)

		require.NoError(t, err)
		assert.Equal(t, updatedName, resp.Name)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 5. Delete Category Type (Fail) ---
	t.Run("5_DeleteCategoryType_FailsWhenInUse", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryTypeID, "CategoryType must exist")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT language_code, COUNT(*)`)).
			WithArgs(createdCategoryTypeID).
			WillReturnRows(sqlmock.NewRows([]string{"language_code", "count"}).AddRow("en", 1))

		err := categoryTypeService.DeleteCategoryType(createdCategoryTypeID.String())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "in use")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 6. Delete Category (Success) ---
	t.Run("6_DeleteCategory", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryID, "Category must exist")

		mock.ExpectQuery(`SELECT \* FROM "categories" WHERE id = \$1 ORDER BY "categories"\."id" LIMIT \$2`).WithArgs(createdCategoryID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).AddRow(createdCategoryID, createdCategoryTypeID))
		mock.ExpectQuery(`SELECT \* FROM "category_types"`).WithArgs(createdCategoryTypeID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdCategoryTypeID))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(createdCategoryID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := categoryService.DeleteCategoryByUUID(createdCategoryID.String())

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 7. Delete Category Type (Success) ---
	t.Run("7_DeleteCategoryType_Success", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryTypeID, "CategoryType must exist")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT language_code, COUNT(*)`)).
			WithArgs(createdCategoryTypeID).
			WillReturnRows(sqlmock.NewRows([]string{"language_code", "count"}))

		mock.ExpectQuery(`SELECT \* FROM "category_types" WHERE id = \$1 ORDER BY "category_types"\."id" LIMIT \$2`).
			WithArgs(createdCategoryTypeID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdCategoryTypeID))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "category_types" WHERE id = $1`)).
			WithArgs(createdCategoryTypeID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := categoryTypeService.DeleteCategoryType(createdCategoryTypeID.String())

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
