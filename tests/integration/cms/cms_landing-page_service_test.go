package integration

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
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

// MockEmailSendingService to simulate email sending behavior
type manualMockLandingEmailService struct {
	done        chan struct{}
	lastRequest dto.SendEmailRequest
	errToReturn error
}

// newManualMockEmailService creates a new instance of our manual mock.
func newManualMockLandingEmailService() *manualMockLandingEmailService {
	return &manualMockLandingEmailService{
		done: make(chan struct{}),
	}
}

// SendEmail implements the EmailSendingServiceInterface for the manual mock.
func (m *manualMockLandingEmailService) SendEmail(req dto.SendEmailRequest) error {
	// Store the request for potential assertions later
	m.lastRequest = req

	defer close(m.done)

	// Return a pre-configured error, if any.
	return m.errToReturn
}

// setupLandingPageTest is updated to provide all necessary dependencies for the service.
func setupLandingPageTest(t *testing.T) (
	sqlmock.Sqlmock,
	services.CMSLandingPageServiceInterface,
	*manualMockLandingEmailService,
	*config.Config,
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	repo := repositories.NewCMSLandingPageRepository(gormDB)
	emailContentRepo := repositories.NewEmailContentRepository(gormDB)
	emailCategoryRepo := repositories.NewEmailCategoryRepository(gormDB)

	mockEmailService := newManualMockLandingEmailService()

	testCfg := &config.Config{
		App: config.AppConfig{
			FrontendURLS: "http://localhost:8000,http://localhost:3000",
			WebBaseURL:   "https://example-frontend.com",
			CMSBaseURL:   "https://cms.example.com",
			APIBaseURL:   "http://localhost:8080",
			UploadPath:   t.TempDir(),
		},
	}

	service := services.NewCMSLandingPageService(
		repo,
		mockEmailService,
		emailContentRepo,
		emailCategoryRepo,
		testCfg,
	)

	return mock, service, mockEmailService, testCfg, cleanup // <--- return manual mock
}
func TestLandingPageLifecycle(t *testing.T) {
	mock, service, _, _, cleanup := setupLandingPageTest(t)

	defer cleanup()

	var createdPageID uuid.UUID
	var createdContentID uuid.UUID
	var createdMetaTagID uuid.UUID
	var duplicatedContentID uuid.UUID
	t.Run("1_CreateLandingPage_Success", func(t *testing.T) {
		t.Log("===> Start: 1_CreateLandingPage_Success")
		mockLandingPage := helpers.InitializeMockLandingPage()
		landingContent := mockLandingPage.Contents[0]

		newLandingPageID := uuid.New()
		newLandingContentID := uuid.New()
		newMetaTagID := uuid.New()

		// Uniqueness checks

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1`)).
			WithArgs(landingContent.UrlAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Transaction and inserts
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newLandingPageID))

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newMetaTagID))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newLandingContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).AddRow(uuid.New(), uuid.New()))
		mock.ExpectCommit()

		resp, err := service.CreateLandingPage(mockLandingPage)
		require.NoError(t, err)
		require.NotNil(t, resp)

		createdPageID = newLandingPageID
		createdContentID = newLandingContentID
		createdMetaTagID = newMetaTagID

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdateLandingContent_Success", func(t *testing.T) {
		t.Log("===> Start: 2_UpdateLandingContent_Success")
		updatedContent := helpers.InitializeMockLandingPage().Contents[0]
		updatedContent.Title = "Updated Common Title"
		updatedContent.UrlAlias += "-updated"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, "en"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(updatedContent.UrlAlias, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Transaction and updates
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, enums.PageLanguageEN))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).AddRow(uuid.New(), uuid.New()))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		resp, err := service.UpdateLandingContent(updatedContent, createdContentID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Common Title", resp.Title)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("3_FindLandingPageById_Success", func(t *testing.T) {
		t.Log("===> Start: 3_FindLandingPageById_Success")
		mockTime := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages" WHERE id = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(createdPageID, mockTime, mockTime))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE "landing_contents"."page_id" = $1 AND (landing_contents.mode != $2 AND landing_contents.mode != $3) ORDER BY landing_contents.created_at DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).AddRow(createdContentID, createdPageID, createdMetaTagID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories" WHERE "landing_content_categories"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(createdMetaTagID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(createdMetaTagID, "mock title", "mock desc"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}))

		_, err := service.FindLandingPageById(createdPageID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("4_FindContentByLandingPageId_Success", func(t *testing.T) {
		t.Log("===> Start: 4_FindContentByLandingPageId_Success")
		lang := string(enums.PageLanguageEN)
		mode := string(enums.PageModeDraft)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "landing_contents"."id" LIMIT $4`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language", "mode"}).AddRow(createdContentID, createdPageID, lang, mode))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories" WHERE "landing_content_categories"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}))

		_, err := service.FindContentByLandingPageId(createdPageID, lang, mode)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("5_DuplicateLandingContentToAnotherLanguage_Success", func(t *testing.T) {
		t.Log("===> Start: 5_DuplicateLandingContentToAnotherLanguage_Success")
		newRev := helpers.InitializeMockRevision()
		duplicatedContentID = uuid.New()

		// 1. Find original content and its preloaded data
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "language", "page_id", "meta_tag_id"}).AddRow(createdContentID, "en", createdPageID, createdMetaTagID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories" WHERE "landing_content_categories"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(createdMetaTagID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(createdMetaTagID, "mock title", "mock desc"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."landing_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}))

		// 2. Transaction for duplication
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(duplicatedContentID))

		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectCommit()

		_, err := service.DuplicateLandingContentToAnotherLanguage(createdContentID, newRev)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("6_DuplicateLandingPage_Success", func(t *testing.T) {
		t.Log("===> Start: 6_DuplicateLandingPage_Success")
		require.NotEqual(t, uuid.Nil, createdPageID, "A page must have been created in a previous step")

		// --- Arrange ---

		mockOriginalPage := helpers.InitializeMockLandingPage()
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages" WHERE id = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdPageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE "landing_contents"."page_id" = $1 AND landing_contents.mode != $2 ORDER BY landing_contents.created_at DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).AddRow(createdContentID, createdPageID, createdMetaTagID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"category_id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).AddRow(uuid.New(), createdContentID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WithArgs(createdMetaTagID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(createdMetaTagID, "Original Meta Title"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).AddRow(uuid.New(), createdContentID))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedPageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedMetaTagID))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedRevisionID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedComponentID))

		mock.ExpectCommit()

		// --- Act ---
		duplicatedPage, err := service.DuplicateLandingPage(createdPageID)

		// --- Assert ---
		require.NoError(t, err)
		require.NotNil(t, duplicatedPage)
		assert.NotEqual(t, createdPageID, duplicatedPage.ID, "Duplicated page ID should be new")
		require.NotEmpty(t, duplicatedPage.Contents, "Duplicated page should have content")
		assert.NotEqual(t, createdContentID, duplicatedPage.Contents[0].ID, "Duplicated content ID should be new")

		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("7_DeleteContentByLandingPageId_Success", func(t *testing.T) {
		t.Log("===> Start: 7_DeleteContentByLandingPageId_Success")

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "landing_contents"."id" LIMIT $4`)).
			WithArgs(createdPageID, "th", "Draft", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(duplicatedContentID, createdPageID))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE Landing_content_id = $1`)).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions" WHERE Landing_content_id = $1`)).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_content_categories" WHERE Landing_content_id = $1`)).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_contents" WHERE "landing_contents"."id" = $1`)).
			WithArgs(duplicatedContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := service.DeleteContentByLandingPageId(createdPageID, "th", "Draft")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("8_DeleteLandingPage_Success", func(t *testing.T) {
		t.Log("===> Start: 8_DeleteLandingPage_Success")

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages" WHERE id = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdPageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE page_id = $1`)).
			WithArgs(createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(createdContentID, createdPageID))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE Landing_content_id IN ($1)`)).
			WithArgs(createdContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_content_categories" WHERE Landing_content_id IN ($1)`)).
			WithArgs(createdContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions" WHERE Landing_content_id IN ($1)`)).
			WithArgs(createdContentID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_content_files" WHERE Landing_content_id IN ($1)`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_contents" WHERE page_id = $1`)).
			WithArgs(createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_pages" WHERE id = $1`)).
			WithArgs(createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := service.DeleteLandingPage(createdPageID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
func TestLandingPageRevertFlow(t *testing.T) {
	mock, service, _, _, cleanup := setupLandingPageTest(t)
	defer cleanup()

	var pageID, contentV1ID, contentV2ID, revisionV1ID uuid.UUID
	var categoryID, metaTagV1ID uuid.UUID

	t.Run("1_CreateInitialPage_V1", func(t *testing.T) {
		t.Log("===> Start: 1_CreateInitialPage_V1")

		// --- Arrange ---
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockLandingPage.Contents[0].Title = "Version 1"
		landingContent := mockLandingPage.Contents[0]

		pageID = uuid.New()
		contentV1ID = uuid.New()
		revisionV1ID = uuid.New()
		metaTagV1ID = uuid.New()
		categoryID = uuid.New()
		newComponentID := uuid.New()
		newCategoryTypeID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1`)).
			WithArgs(landingContent.UrlAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Transaction และ INSERT
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(metaTagV1ID))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentV1ID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(revisionV1ID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newComponentID))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newCategoryTypeID))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).AddRow(contentV1ID, categoryID))
		mock.ExpectCommit()

		// --- Act ---
		_, err := service.CreateLandingPage(mockLandingPage)

		// --- Assert ---
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdateContent_To_V2", func(t *testing.T) {
		t.Log("===> Start: 2_UpdateContent_To_V2")
		updatedContent := helpers.InitializeMockLandingPage().Contents[0]
		updatedContent.Title = "Version 2"

		contentV2ID = uuid.New()

		// Mock การหา Content V1
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(contentV1ID, pageID, "en"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Mock Transaction
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(contentV1ID, pageID, enums.PageLanguageEN))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentV2ID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id"}))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err := service.UpdateLandingContent(updatedContent, contentV1ID)
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
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).AddRow(revisionV1ID, contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "meta_tag_id", "language"}).
				AddRow(contentV1ID, pageID, "Version 1", metaTagV1ID, enums.PageLanguageEN))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories" WHERE "landing_content_categories"."landing_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).AddRow(contentV1ID, categoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(categoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(categoryID, "Mock Category"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."landing_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).AddRow(uuid.New(), contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(metaTagV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(metaTagV1ID, "Meta Title V1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."landing_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).AddRow(revisionV1ID, contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE page_id = $1 AND language = $2 ORDER BY created_at DESC,"landing_contents"."id" LIMIT $3`)).
			WithArgs(pageID, enums.PageLanguageEN, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(contentV2ID, pageID))

		// ----- Transaction: Revert -----
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(metaTagV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newContentV3ID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Mock Category", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}))

		mock.ExpectCommit()

		// --- Act ---
		revertedContent, err := service.RevertLandingContent(revisionV1ID, revertAuthor)

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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages" WHERE id = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WithArgs(pageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE "landing_contents"."page_id" = $1 AND language = $2`)).
			WithArgs(pageID, language).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "mode"}).
				AddRow(contentV1ID, pageID, "Histories").
				AddRow(contentV2ID, pageID, "Histories"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."landing_content_id" IN ($1,$2)`)).
			WithArgs(contentV1ID, contentV2ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id", "author", "created_at"}).
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

func TestLandingPagePreviewFlow(t *testing.T) {
	mock, service, _, testCfg, cleanup := setupLandingPageTest(t)
	defer cleanup()

	var pageID uuid.UUID
	var previewContentID uuid.UUID

	t.Run("1_CreateBasePageForPreview", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		landingContent := mockLandingPage.Contents[0]
		pageID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1`)).
			WithArgs(landingContent.UrlAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "landing_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id"}))
		mock.ExpectCommit()

		_, err := service.CreateLandingPage(mockLandingPage)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("2_PreviewContent_FirstTime_CreatesNew", func(t *testing.T) {
		t.Log("===> Start: 2_PreviewContent_FirstTime_CreatesNew")
		previewContent := helpers.InitializeMockLandingPage().Contents[0]
		previewContent.Language = enums.PageLanguageEN

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(previewContent.UrlAlias, pageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "landing_contents"."id" LIMIT $4`)).
			WithArgs(pageID, "en", "Preview", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		previewContentID = uuid.New()
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(previewContentID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		// --- Act ---
		previewURL, err := service.PreviewLandingContent(pageID, previewContent)

		// --- Assert ---
		require.NoError(t, err)
		frontendURLs := strings.Split(testCfg.App.FrontendURLS, ",")
		appURL := frontendURLs[1]
		expectedURL, err := helpers.BuildPreviewURL(
			appURL,
			string(previewContent.Language),
			"landing",
			previewContentID,
		)
		assert.Equal(t, expectedURL, previewURL)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("3_PreviewContent_SecondTime_UpdatesExisting", func(t *testing.T) {
		t.Log("===> Start: 3_PreviewContent_SecondTime_UpdatesExisting")
		updatedPreviewContent := helpers.InitializeMockLandingPage().Contents[0]
		updatedPreviewContent.Language = enums.PageLanguageEN
		updatedPreviewContent.Title = "Updated Preview Title"

		// --- Mock ---

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(updatedPreviewContent.UrlAlias, pageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		existingPreview := *updatedPreviewContent
		existingPreview.ID = previewContentID
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "landing_contents"."id" LIMIT $4`)).
			WithArgs(pageID, "en", "Preview", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(existingPreview.ID))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE landing_content_id = $1`)).
			WithArgs(previewContentID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()
		// --- Act ---
		previewURL, err := service.PreviewLandingContent(pageID, updatedPreviewContent)

		// --- Assert ---
		require.NoError(t, err)
		frontendURLs := strings.Split(testCfg.App.FrontendURLS, ",")
		appURL := frontendURLs[1]
		expectedURL, err := helpers.BuildPreviewURL(
			appURL,
			"en",
			"landing",
			previewContentID,
		)
		assert.Equal(t, expectedURL, previewURL)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLandingPageCategoryFlow(t *testing.T) {
	mock, service, _, _, cleanup := setupLandingPageTest(t)
	defer cleanup()

	var pageID, contentID, landingCategoryID, keywordCategoryID uuid.UUID

	const keywordCategoryTypeCode = "category_keywords"

	t.Run("1_CreatePageWithMixedCategories", func(t *testing.T) {
		t.Log("===> Start: 1_CreatePageWithMixedCategories")

		mockLandingPage := helpers.InitializeMockLandingPage()

		landingContent := mockLandingPage.Contents[0]

		keywordCategoryType := &models.CategoryType{Name: "Keywords", TypeCode: keywordCategoryTypeCode, IsActive: true}
		keywordCategory := &models.Category{Name: "General Keywords", CategoryType: keywordCategoryType, LanguageCode: enums.PageLanguageEN}
		mockLandingPage.Contents[0].Categories = append(mockLandingPage.Contents[0].Categories, keywordCategory)

		pageID = uuid.New()
		contentID = uuid.New()
		landingCategoryID = uuid.New()
		keywordCategoryID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1`)).
			WithArgs(landingContent.UrlAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(keywordCategoryID))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_content_categories"`)).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}))

		mock.ExpectCommit()

		_, err := service.CreateLandingPage(mockLandingPage)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_FindCategoriesByLandingTypeCode_Success", func(t *testing.T) {
		t.Log("===> Start: 2_FindCategoriesByLandingTypeCode_Success")
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModeDraft)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT "id" FROM "landing_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 LIMIT $4`,
		)).WithArgs(pageID, language, mode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "category_id" FROM "landing_content_categories" WHERE Landing_content_id = $1`)).
			WithArgs(contentID).
			WillReturnRows(sqlmock.NewRows([]string{"category_id"}).AddRow(landingCategoryID).AddRow(keywordCategoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "categories"."id","categories"."category_type_id","categories"."language_code","categories"."name","categories"."description","categories"."weight","categories"."publish_status","categories"."created_at","categories"."updated_at" FROM "categories" JOIN category_types ON category_types.id = categories.category_type_id WHERE categories.id IN ($1,$2) AND category_types.type_code = $3`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(keywordCategoryID, "General Keywords"))

		categories, err := service.GetCategory(pageID, keywordCategoryTypeCode, language, mode)
		require.NoError(t, err)
		require.Len(t, categories, 1, "Should only find one category with the specified type code")
		assert.Equal(t, keywordCategoryID, categories[0].ID)
		assert.Equal(t, "General Keywords", categories[0].Name)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLandingPagEmailLifecycle(t *testing.T) {
	t.Log("===> Start: TestLandingPagEmailLifecycle")
	mock, service, mockEmailService, _, cleanup := setupLandingPageTest(t)
	defer cleanup()

	var createdPageID uuid.UUID
	var createdContentID uuid.UUID

	t.Run("1_CreateLandingPage_Success", func(t *testing.T) {
		t.Log("===> Start: 1_CreateLandingPage_Success")
		mockLandingPage := helpers.InitializeMockLandingPage()
		landingContent := mockLandingPage.Contents[0]

		newLandingPageID := uuid.New()
		newLandingContentID := uuid.New()
		newMetaTagID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1`)).
			WithArgs(landingContent.UrlAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newLandingPageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newMetaTagID))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newLandingContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).AddRow(uuid.New(), uuid.New()))
		mock.ExpectCommit()

		resp, err := service.CreateLandingPage(mockLandingPage)
		require.NoError(t, err)
		require.NotNil(t, resp)

		createdPageID = newLandingPageID
		createdContentID = newLandingContentID

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdateLandingContent_AndTriggerEmail_Success", func(t *testing.T) {
		t.Log("===> Start: 2_UpdateLandingContent_AndTriggerEmail_Success")

		updatedContent := helpers.InitializeMockLandingPage().Contents[0]
		updatedContent.Title = "Updated Common Title"
		updatedContent.UrlAlias += "-updated"
		updatedContent.WorkflowStatus = enums.WorkflowWaitingDesign
		updatedContent.ApprovalEmail = []string{"admin1@example.com"}

		newContentID := uuid.New()

		// --- Mock DB interactions ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, enums.PageLanguageEN))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(updatedContent.UrlAlias, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 ORDER BY "landing_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, enums.PageLanguageEN))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id"}))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WithArgs(sqlmock.AnyArg(), createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// --- Mock a final preload ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents" WHERE id = $1 AND "landing_contents"."id" = $2 ORDER BY "landing_contents"."id" LIMIT $3`)).
			WithArgs(newContentID, newContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "language", "workflow_status", "approval_email", "url_alias"}).
				AddRow(newContentID, createdPageID, updatedContent.Title, updatedContent.Language, updatedContent.WorkflowStatus, updatedContent.ApprovalEmail, updatedContent.UrlAlias))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).WillReturnRows(sqlmock.NewRows([]string{"landing_content_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_files"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		// --- Mock DB interactions for Email ---
		approvalCatID := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE title = $1 ORDER BY "email_categories"."id" LIMIT $2`)).
			WithArgs("Approve", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(approvalCatID, "Approve"))
		templateID := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1 AND language = $2`)).
			WithArgs(approvalCatID.String(), updatedContent.Language).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label"}).AddRow(templateID, "email_to_admin"))

		// --- Act ---
		resp, err := service.UpdateLandingContent(updatedContent, createdContentID)

		// --- Assert ---
		require.NoError(t, err)
		assert.Equal(t, "Updated Common Title", resp.Title)

		select {
		case <-mockEmailService.done:
		case <-time.After(2 * time.Second): // Timeout
			t.Fatal("Test timed out waiting for SendEmail to be called")
		}

		require.NoError(t, mock.ExpectationsWereMet())

	})
}
