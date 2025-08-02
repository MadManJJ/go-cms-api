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

func setupEmailTest(t *testing.T) (
	*gorm.DB,
	sqlmock.Sqlmock,
	services.EmailCategoryServiceInterface,
	services.EmailContentServiceInterface,
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	categoryRepo := repositories.NewEmailCategoryRepository(gormDB)
	contentRepo := repositories.NewEmailContentRepository(gormDB)
	categoryService := services.NewEmailCategoryService(categoryRepo, contentRepo)
	contentService := services.NewEmailContentService(contentRepo, categoryRepo)
	return gormDB, mock, categoryService, contentService, cleanup
}

func TestEmailTemplateLifecycle(t *testing.T) {
	_, mock, emailCategoryService, emailContentService, cleanup := setupEmailTest(t)
	defer cleanup()

	var createdCategoryID uuid.UUID
	var createdContentID uuid.UUID
	categoryTitle := "Order Confirmation Emails"
	contentLabel := "order-success-customer"

	// --- 1. Create Email Category (Success) ---
	t.Run("1_CreateEmailCategory", func(t *testing.T) {
		req := dto.CreateEmailCategoryRequest{Title: categoryTitle}
		tempID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_categories" WHERE title = $1`)).
			WithArgs(categoryTitle).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "email_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempID))
		mock.ExpectCommit()

		resp, err := emailCategoryService.CreateCategory(req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		createdCategoryID, _ = uuid.Parse(resp.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 2. Create Email Content (Success) ---
	t.Run("2_CreateEmailContent", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryID, "CategoryType must be created first")

		req := dto.CreateEmailContentRequest{
			EmailCategoryID: createdCategoryID.String(),
			Language:        enums.PageLanguageEN,
			Label:           contentLabel,
			EmailContentDetailBase: dto.EmailContentDetailBase{
				SendFromEmail: "shop@example.com",
				Subject:       "Your Order #[Order.ID] is Confirmed!",
			},
		}

		mock.ExpectQuery(`SELECT \* FROM "email_categories" WHERE id = \$1 ORDER BY "email_categories"\."id" LIMIT \$2`).
			WithArgs(createdCategoryID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdCategoryID))

		mock.ExpectQuery(`SELECT \* FROM "email_contents" WHERE email_category_id = \$1 AND language = \$2 AND label = \$3 ORDER BY "email_contents"\."id" LIMIT \$4`).
			WithArgs(createdCategoryID, enums.PageLanguageEN, contentLabel, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		tempContentID := uuid.New()
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "email_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempContentID))
		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT \* FROM "email_contents" WHERE id = \$1 ORDER BY "email_contents"\."id" LIMIT \$2`).
			WithArgs(tempContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label", "email_category_id"}).AddRow(tempContentID, contentLabel, createdCategoryID))
		mock.ExpectQuery(`SELECT \* FROM "email_categories" WHERE "email_categories"\."id" = \$1`).
			WithArgs(createdCategoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(createdCategoryID, categoryTitle))

		resp, err := emailContentService.CreateContent(req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		createdContentID, _ = uuid.Parse(resp.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 3. Update Email Content (Success) ---
	t.Run("3_UpdateEmailContent", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdContentID, "EmailContent must exist")

		updatedSubject := "Updated Lifecycle Subject! " + uuid.NewString()[:6]
		req := dto.UpdateEmailContentRequest{Subject: &updatedSubject}

		// 1. Service จะเรียก FindByID เพื่อหา Content เดิม
		mock.ExpectQuery(`SELECT \* FROM "email_contents" WHERE id = \$1 ORDER BY "email_contents"\."id" LIMIT \$2`).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "subject", "email_category_id", "label", "language"}).
				AddRow(createdContentID, "Old Subject", createdCategoryID, contentLabel, enums.PageLanguageEN))

		// Preload EmailCategory
		mock.ExpectQuery(`SELECT \* FROM "email_categories" WHERE "email_categories"\."id" = \$1`).
			WithArgs(createdCategoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(createdCategoryID, categoryTitle))

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "email_categories"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), createdCategoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdCategoryID))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "email_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		mock.ExpectQuery(`SELECT \* FROM "email_contents" WHERE id = \$1 ORDER BY "email_contents"\."id" LIMIT \$2`).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "subject", "email_category_id"}).AddRow(createdContentID, updatedSubject, createdCategoryID))

		mock.ExpectQuery(`SELECT \* FROM "email_categories" WHERE "email_categories"\."id" = \$1`).
			WithArgs(createdCategoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(createdCategoryID, categoryTitle))

		resp, err := emailContentService.UpdateContent(createdContentID.String(), req)

		require.NoError(t, err)
		assert.Equal(t, updatedSubject, resp.Subject)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 4. Delete Email Category (Fails) ---
	t.Run("4_DeleteEmailCategory_FailsWhenInUse", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryID, "EmailCategory must exist")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_contents" WHERE email_category_id = $1`)).
			WithArgs(createdCategoryID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "email_contents" WHERE email_category_id = $1`)).
			WithArgs(createdCategoryID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		err := emailCategoryService.DeleteCategory(createdCategoryID.String())

		require.Error(t, err)
		assert.Contains(t, err.Error(), "in use")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 5. Delete Email Content (Success) ---
	t.Run("5_DeleteEmailContent", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdContentID, "EmailContent must exist")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_contents" WHERE id = $1`)).
			WithArgs(createdContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := emailContentService.DeleteContent(createdContentID.String())

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 6. Delete Email Category (Success) ---
	t.Run("6_DeleteEmailCategory_Success", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdCategoryID, "EmailCategory must exist")

		// 1. Mock DeleteByCategoryID ของ content repo
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "email_contents"`).
			WithArgs(createdCategoryID).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		// 2. Mock Delete ของ category repo
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
			WithArgs(createdCategoryID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "email_categories" WHERE id = $1`)).
			WithArgs(createdCategoryID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := emailCategoryService.DeleteCategory(createdCategoryID.String())

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
