package integration

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"
	"github.com/MadManJJ/cms-api/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupFaqPageTest(t *testing.T) (
	sqlmock.Sqlmock,
	services.CMSFaqPageServiceInterface,
	*config.Config,
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	repo := repositories.NewCMSFaqPageRepository(gormDB)

	testCfg := &config.Config{
		App: config.AppConfig{
			FrontendURLS: "http://localhost:8000,http://localhost:3000",
			WebBaseURL:   "https://example-frontend.com",
			APIBaseURL:   "http://localhost:8080",
			UploadPath:   t.TempDir(),
		},
	}

	service := services.NewCMSFaqPageService(repo, testCfg)
	return mock, service, testCfg, cleanup
}

func TestFaqPageLifecycle(t *testing.T) {
	mock, service, _, cleanup := setupFaqPageTest(t)
	defer cleanup()

	var createdPageID uuid.UUID
	var createdContentID uuid.UUID
	var createdMetaTagID uuid.UUID
	var duplicatedContentID uuid.UUID
	t.Run("1_CreateFaqPage_Success", func(t *testing.T) {
		t.Log("===> Start: 1_CreateFaqPage_Success")
		mockFaqPage := helpers.InitializeMockFaqPage()
		faqContent := mockFaqPage.Contents[0]

		newFaqPageID := uuid.New()
		newFaqContentID := uuid.New()
		newMetaTagID := uuid.New()

		// Uniqueness checks
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url = $1`)).
			WithArgs(faqContent.URL).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url_alias = $1`)).
			WithArgs(faqContent.URLAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Transaction and inserts
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newFaqPageID))

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newMetaTagID))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newFaqContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).AddRow(uuid.New(), uuid.New()))
		mock.ExpectCommit()

		resp, err := service.CreateFaqPage(mockFaqPage)
		require.NoError(t, err)
		require.NotNil(t, resp)

		createdPageID = newFaqPageID
		createdContentID = newFaqContentID
		createdMetaTagID = newMetaTagID

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdateFaqContent_Success", func(t *testing.T) {
		t.Log("===> Start: 2_UpdateFaqContent_Success")
		updatedContent := helpers.InitializeMockFaqPage().Contents[0]
		updatedContent.Title = "Updated FAQ Title"
		updatedContent.URL += "-updated"

		// Find existing content

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE id = $1 ORDER BY "faq_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, "en"))

		// Uniqueness checks
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url = $1 AND page_id != $2`)).
			WithArgs(updatedContent.URL, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(updatedContent.URLAlias, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Transaction and updates
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE id = $1 ORDER BY "faq_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, enums.PageLanguageEN))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_contents" SET "page_id"=$1,"title"=$2,"language"=$3,"authored_at"=$4,"html_input"=$5,"mode"=$6,"workflow_status"=$7,"publish_status"=$8,"publish_on"=$9,"unpublish_on"=$10,"authored_on"=$11,"url_alias"=$12,"url"=$13,"meta_tag_id"=$14,"expired_at"=$15,"created_at"=$16,"updated_at"=$17 WHERE "id" = $18`)).
			WithArgs(
				createdPageID,    // page_id
				sqlmock.AnyArg(), // title
				"en",             // language
				sqlmock.AnyArg(), // authored_at
				sqlmock.AnyArg(), // html_input
				"Histories",      // mode
				sqlmock.AnyArg(), // workflow_status
				sqlmock.AnyArg(), // publish_status
				sqlmock.AnyArg(), // publish_on
				sqlmock.AnyArg(), // unpublish_on
				sqlmock.AnyArg(), // authored_on
				sqlmock.AnyArg(), // url_alias
				sqlmock.AnyArg(), // url
				sqlmock.AnyArg(), // meta_tag_id
				sqlmock.AnyArg(), // expired_at
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at (เป็น timestamp)
				createdContentID, // WHERE id =
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).AddRow(uuid.New(), uuid.New()))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WithArgs(sqlmock.AnyArg(), createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		resp, err := service.UpdateFaqContent(updatedContent, createdContentID)
		require.NoError(t, err)
		assert.Equal(t, "Updated FAQ Title", resp.Title)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("3_FindFaqPageById_Success", func(t *testing.T) {
		t.Log("===> Start: 3_FindFaqPageById_Success")
		mockTime := time.Now()

		// 1. Main FAQ page query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages" WHERE id = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(createdPageID, mockTime, mockTime))

		// 2. Associated content query (หลัก)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE "faq_contents"."page_id" = $1 AND (faq_contents.mode != $2 AND faq_contents.mode != $3) ORDER BY faq_contents.created_at DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).AddRow(createdContentID, createdPageID, createdMetaTagID))

		// 3. Preload: Categories
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories" WHERE "faq_content_categories"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}))

		// 4. Preload: Components
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}))

		// 5. Preload: MetaTag

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(createdMetaTagID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(createdMetaTagID, "mock title", "mock desc"))

		// 6. Preload: Revisions
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}))

		_, err := service.FindFaqPageById(createdPageID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("4_FindContentByFaqPageId_Success", func(t *testing.T) {
		t.Log("===> Start: 4_FindContentByFaqPageId_Success")
		lang := string(enums.PageLanguageEN)
		mode := string(enums.PageModeDraft)

		// 1. Query หลัก
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "faq_contents"."id" LIMIT $4`)).
			WithArgs(createdPageID, lang, mode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language", "mode"}).AddRow(createdContentID, createdPageID, lang, mode))

		// 2. Preload: Categories
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories" WHERE "faq_content_categories"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}))

		// 3. Preload: Components
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}))

		// 4. Preload: Revisions
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}))

		_, err := service.FindContentByFaqPageId(createdPageID, lang, mode)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("5_DuplicateFaqContentToAnotherLanguage_Success", func(t *testing.T) {
		t.Log("===> Start: 5_DuplicateFaqContentToAnotherLanguage_Success")
		newRev := helpers.InitializeMockRevision()
		duplicatedContentID = uuid.New()

		// 1. Find original content and its preloaded data
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE id = $1 ORDER BY "faq_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "language", "page_id", "meta_tag_id"}).AddRow(createdContentID, "en", createdPageID, createdMetaTagID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories" WHERE "faq_content_categories"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(createdMetaTagID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(createdMetaTagID, "mock title", "mock desc"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."faq_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}))

		// 2. Transaction for duplication
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(duplicatedContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectCommit()

		_, err := service.DuplicateFaqContentToAnotherLanguage(createdContentID, newRev)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("6_DuplicateFaqPage_Success", func(t *testing.T) {
		t.Log("===> Start: 6_DuplicateFaqPage_Success")
		require.NotEqual(t, uuid.Nil, createdPageID, "A page must have been created in a previous step")

		// --- Arrange ---

		mockOriginalPage := helpers.InitializeMockFaqPage()
		mockOriginalContent := mockOriginalPage.Contents[0]

		mockOriginalPage.ID = createdPageID
		mockOriginalContent.ID = createdContentID
		mockOriginalContent.PageID = createdPageID
		mockOriginalContent.MetaTag.ID = createdMetaTagID

		newDuplicatedPageID := uuid.New()
		newDuplicatedContentID := uuid.New()
		newDuplicatedMetaTagID := uuid.New()
		newDuplicatedRevisionID := uuid.New()
		newDuplicatedComponentID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages" WHERE id = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdPageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE "faq_contents"."page_id" = $1 AND faq_contents.mode != $2 ORDER BY faq_contents.created_at DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).AddRow(createdContentID, createdPageID, createdMetaTagID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"category_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).AddRow(uuid.New(), createdContentID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WithArgs(createdMetaTagID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(createdMetaTagID, "Original Meta Title"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).AddRow(uuid.New(), createdContentID))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedPageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedMetaTagID))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedRevisionID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedComponentID))

		mock.ExpectCommit()

		// --- Act ---
		duplicatedPage, err := service.DuplicateFaqPage(createdPageID)

		// --- Assert ---
		require.NoError(t, err)
		require.NotNil(t, duplicatedPage)
		assert.NotEqual(t, createdPageID, duplicatedPage.ID, "Duplicated page ID should be new")
		require.NotEmpty(t, duplicatedPage.Contents, "Duplicated page should have content")
		assert.NotEqual(t, createdContentID, duplicatedPage.Contents[0].ID, "Duplicated content ID should be new")

		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("7_DeleteContentByFaqPageId_Success", func(t *testing.T) {
		t.Log("===> Start: 7_DeleteContentByFaqPageId_Success")

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "faq_contents"."id" LIMIT $4`)).
			WithArgs(createdPageID, "th", "Draft", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(duplicatedContentID, createdPageID))

		mock.ExpectExec(`DELETE FROM "components" WHERE "?faq_content_id"? = \$1`).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`DELETE FROM "revisions" WHERE "?faq_content_id"? = \$1`).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`DELETE FROM "faq_content_categories" WHERE "?faq_content_id"? = \$1`).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "faq_contents" WHERE "faq_contents"."id" = $1`)).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := service.DeleteContentByFaqPageId(createdPageID, "th", "Draft")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("8_DeleteFaqPage_Success", func(t *testing.T) {
		t.Log("===> Start: 8_DeleteFaqPage_Success")

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages" WHERE id = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdPageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE page_id = $1`)).
			WithArgs(createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(createdContentID, createdPageID))

		mock.ExpectExec(`DELETE FROM "components" WHERE faq_content_id IN \(\$1\)`).
			WithArgs(createdContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`DELETE FROM "faq_content_categories" WHERE faq_content_id IN \(\$1\)`).
			WithArgs(createdContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`DELETE FROM "revisions" WHERE faq_content_id IN \(\$1\)`).
			WithArgs(createdContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "faq_contents" WHERE page_id = $1`)).
			WithArgs(createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`DELETE FROM "faq_pages" WHERE id = \$1`).
			WithArgs(createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := service.DeleteFaqPage(createdPageID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
func TestFaqPageRevertFlow(t *testing.T) {
	t.Log("===> Start: TestFaqPageRevertFlow")
	mock, service, _, cleanup := setupFaqPageTest(t)
	defer cleanup()

	var pageID, contentV1ID, contentV2ID, revisionV1ID uuid.UUID
	var categoryID, metaTagV1ID uuid.UUID

	t.Run("1_CreateInitialPage_V1", func(t *testing.T) {
		t.Log("===> Start: 1_CreateInitialPage_V1")

		// --- Arrange ---
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqPage.Contents[0].Title = "Version 1"
		faqContent := mockFaqPage.Contents[0]

		// กำหนด ID ที่จะใช้
		pageID = uuid.New()
		contentV1ID = uuid.New()
		revisionV1ID = uuid.New()
		metaTagV1ID = uuid.New()
		categoryID = uuid.New()
		newComponentID := uuid.New()
		newCategoryTypeID := uuid.New()

		// --- Mock ---

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url = $1`)).
			WithArgs(faqContent.URL).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url_alias = $1`)).
			WithArgs(faqContent.URLAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Transaction และ INSERT
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(metaTagV1ID))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentV1ID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(revisionV1ID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newComponentID))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newCategoryTypeID))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))
		mock.ExpectQuery(`INSERT INTO "faq_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).AddRow(contentV1ID, categoryID))
		mock.ExpectCommit()

		// --- Act ---
		_, err := service.CreateFaqPage(mockFaqPage)

		// --- Assert ---
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdateContent_To_V2", func(t *testing.T) {
		t.Log("===> Start: 2_UpdateContent_To_V2")
		updatedContent := helpers.InitializeMockFaqPage().Contents[0]
		updatedContent.Title = "Version 2"

		contentV2ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE id = $1 ORDER BY "faq_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(contentV1ID, pageID, "en"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url = $1 AND page_id != $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Mock Transaction
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE id = $1 ORDER BY "faq_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(contentV1ID, pageID, enums.PageLanguageEN))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentV2ID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"faq_content_id"}))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err := service.UpdateFaqContent(updatedContent, contentV1ID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("3_RevertTo_V1", func(t *testing.T) {
		t.Log("===> Start: 3_RevertTo_V1")
		require.NotEqual(t, uuid.Nil, revisionV1ID, "Revision V1 ID must be available")

		revertAuthor := &models.Revision{Author: "Reverter User", PublishStatus: enums.PublishStatusPublished}
		newContentV3ID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE id = $1 ORDER BY "revisions"."id" LIMIT $2`)).
			WithArgs(revisionV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).AddRow(revisionV1ID, contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE id = $1 ORDER BY "faq_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "meta_tag_id", "language"}).
				AddRow(contentV1ID, pageID, "Version 1", metaTagV1ID, enums.PageLanguageEN))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories" WHERE "faq_content_categories"."faq_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).AddRow(contentV1ID, categoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(categoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(categoryID, "Mock Category"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."faq_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).AddRow(uuid.New(), contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(metaTagV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(metaTagV1ID, "Meta Title V1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."faq_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).AddRow(revisionV1ID, contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE page_id = $1 AND language = $2 ORDER BY created_at DESC,"faq_contents"."id" LIMIT $3`)).
			WithArgs(pageID, enums.PageLanguageEN, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(contentV2ID, pageID))

		// ----- Transaction: Revert -----
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(metaTagV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newContentV3ID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Mock Category", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}))

		mock.ExpectCommit()

		// --- Act ---
		revertedContent, err := service.RevertFaqContent(revisionV1ID, revertAuthor)

		// --- Assert ---
		require.NoError(t, err)
		require.NotNil(t, revertedContent)
		assert.Equal(t, "Version 1", revertedContent.Title, "Content should be reverted to Version 1's title")
		assert.Equal(t, enums.PageModeDraft, revertedContent.Mode, "Reverted content should be in Draft mode")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("4_FindRevisions_Success", func(t *testing.T) {
		t.Log("===> Start: 4_FindRevisions_Success")
		require.NotEqual(t, uuid.Nil, pageID, "PageID must exist from previous steps")

		// --- Arrange ---
		language := string(enums.PageLanguageEN)
		revisionV2ID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages" WHERE id = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WithArgs(pageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE "faq_contents"."page_id" = $1 AND language = $2`)).
			WithArgs(pageID, language).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "mode"}).
				AddRow(contentV1ID, pageID, "Histories").
				AddRow(contentV2ID, pageID, "Histories"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."faq_content_id" IN ($1,$2)`)).
			WithArgs(contentV1ID, contentV2ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id", "author", "created_at"}).
				AddRow(revisionV1ID, contentV1ID, "Initial Author", time.Now().Add(-2*time.Hour)).
				AddRow(revisionV2ID, contentV2ID, "Updater User", time.Now().Add(-1*time.Hour)))

		// --- Act ---
		revisions, err := service.FindRevisions(pageID, language)

		// --- Assert ---
		require.NoError(t, err)
		require.Len(t, revisions, 2, "Should find two revisions for this page in this language")

		assert.Equal(t, revisionV2ID, revisions[0].ID)
		assert.Equal(t, "Updater User", revisions[0].Author)
		assert.Equal(t, revisionV1ID, revisions[1].ID)
		assert.Equal(t, "Initial Author", revisions[1].Author)

		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestFaqPagePreviewFlow(t *testing.T) {
	mock, service, cfg, cleanup := setupFaqPageTest(t)
	defer cleanup()

	var pageID uuid.UUID
	var previewContentID uuid.UUID

	t.Run("1_CreateBasePageForPreview", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		pageID = uuid.New()

		mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"faq_content_id"}))
		mock.ExpectCommit()

		_, err := service.CreateFaqPage(mockFaqPage)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_PreviewContent_FirstTime_CreatesNew", func(t *testing.T) {
		t.Log("===> Start: 2_PreviewContent_FirstTime_CreatesNew")
		previewContent := helpers.InitializeMockFaqPage().Contents[0]
		previewContent.Language = enums.PageLanguageEN

		// --- Mock ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE page_id = $1 AND language = $2 AND mode = $3`)).WillReturnError(gorm.ErrRecordNotFound)

		previewContentID = uuid.New()
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(previewContentID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		// --- Act ---
		previewURL, err := service.PreviewFaqContent(pageID, previewContent)

		// --- Assert ---
		require.NoError(t, err)

		frontendURLs := strings.Split(cfg.App.FrontendURLS, ",")
		appURL := frontendURLs[1]
		expectedURL := fmt.Sprintf("%s/preview/%s/%s?id=%s", appURL, string(previewContent.Language), "faq", previewContentID.String())

		assert.Equal(t, expectedURL, previewURL)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("3_PreviewContent_SecondTime_UpdatesExisting", func(t *testing.T) {
		t.Log("===> Start: 3_PreviewContent_SecondTime_UpdatesExisting")
		updatedPreviewContent := helpers.InitializeMockFaqPage().Contents[0]
		updatedPreviewContent.Language = enums.PageLanguageEN
		updatedPreviewContent.Title = "Updated Preview Title"

		// --- Mock ---

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		existingPreview := *updatedPreviewContent
		existingPreview.ID = previewContentID
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "faq_contents"."id" LIMIT $4`)).
			WithArgs(pageID, "en", "Preview", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(existingPreview.ID))

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE faq_content_id = $1`)).
			WithArgs(previewContentID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectCommit()

		// --- Act ---
		previewURL, err := service.PreviewFaqContent(pageID, updatedPreviewContent)

		// --- Assert ---
		require.NoError(t, err)

		frontendURLs := strings.Split(cfg.App.FrontendURLS, ",")
		appURL := frontendURLs[1]
		expectedURL := fmt.Sprintf("%s/preview/%s/%s?id=%s", appURL, "en", "faq", previewContentID.String())

		assert.Equal(t, expectedURL, previewURL)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFaqPageCategoryFlow(t *testing.T) {
	mock, service, _, cleanup := setupFaqPageTest(t)
	defer cleanup()

	var pageID, contentID, faqCategoryID, keywordCategoryID uuid.UUID
	const faqCategoryTypeCode = "faq"
	const keywordCategoryTypeCode = "category_keywords"

	t.Run("1_CreatePageWithMixedCategories", func(t *testing.T) {
		t.Log("===> Start: 1_CreatePageWithMixedCategories")

		mockFaqPage := helpers.InitializeMockFaqPage()

		keywordCategoryType := &models.CategoryType{Name: "Keywords", TypeCode: keywordCategoryTypeCode, IsActive: true}
		keywordCategory := &models.Category{Name: "General Keywords", CategoryType: keywordCategoryType, LanguageCode: enums.PageLanguageEN}

		mockFaqPage.Contents[0].Categories = append(mockFaqPage.Contents[0].Categories, keywordCategory)

		pageID = uuid.New()
		contentID = uuid.New()
		faqCategoryID = uuid.New()
		keywordCategoryID = uuid.New()
		faqCategoryTypeID := uuid.New()
		keywordCategoryTypeID := uuid.New()

		mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(faqCategoryTypeID).AddRow(keywordCategoryTypeID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(faqCategoryID).AddRow(keywordCategoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(contentID, faqCategoryID).AddRow(contentID, keywordCategoryID))

		mock.ExpectCommit()

		// --- Act ---
		_, err := service.CreateFaqPage(mockFaqPage)

		// --- Assert ---
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_FindCategoriesByFaqTypeCode_Success", func(t *testing.T) {
		t.Log("===> Start: 2_FindCategoriesByFaqTypeCode_Success")
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModeDraft)

		// --- Mock ---

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT "id" FROM "faq_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 LIMIT $4`,
		)).WithArgs(pageID, language, mode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "category_id" FROM "faq_content_categories"`)).
			WithArgs(contentID).
			WillReturnRows(sqlmock.NewRows([]string{"category_id"}).AddRow(faqCategoryID).AddRow(keywordCategoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "categories"."id","categories"."category_type_id","categories"."language_code","categories"."name","categories"."description","categories"."weight","categories"."publish_status","categories"."created_at","categories"."updated_at" FROM "categories" JOIN category_types ON category_types.id = categories.category_type_id WHERE categories.id IN ($1,$2) AND category_types.type_code = $3`)).
			WithArgs(faqCategoryID, keywordCategoryID, faqCategoryTypeCode).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(faqCategoryID, "FAQ Category Name"))

		// --- Act ---
		categories, err := service.FindCategories(pageID, faqCategoryTypeCode, language, mode)

		// --- Assert ---
		require.NoError(t, err)
		require.Len(t, categories, 1, "Should only find one category with the specified type code")
		assert.Equal(t, faqCategoryID, categories[0].ID)
		assert.Equal(t, "FAQ Category Name", categories[0].Name)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
