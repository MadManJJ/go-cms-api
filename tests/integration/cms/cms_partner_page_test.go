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

type manualMockPartnerEmailService struct {
	done        chan struct{}
	lastRequest dto.SendEmailRequest
	errToReturn error
}

// newManualMockPartnerEmailService creates a new instance of our manual mock.
func newManualMockPartnerEmailService() *manualMockPartnerEmailService {
	return &manualMockPartnerEmailService{
		done: make(chan struct{}),
	}
}

// SendEmail implements the EmailSendingServiceInterface.
// This is the key change: the method name must match the interface.
func (m *manualMockPartnerEmailService) SendEmail(req dto.SendEmailRequest) error {
	m.lastRequest = req
	defer close(m.done)
	return m.errToReturn
}

// setupPartnerPageTest provides all dependencies for the partner page service tests.
func setupPartnerPageTest(t *testing.T) (
	sqlmock.Sqlmock,
	services.CMSPartnerPageServiceInterface,
	*manualMockPartnerEmailService,
	*config.Config,
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)

	// Instantiate repositories
	repo := repositories.NewCMSPartnerPageRepository(gormDB)
	emailContentRepo := repositories.NewEmailContentRepository(gormDB)
	emailCategoryRepo := repositories.NewEmailCategoryRepository(gormDB)

	// Create an instance of our manual mock
	mockEmailService := newManualMockPartnerEmailService()

	testCfg := &config.Config{
		App: config.AppConfig{
			FrontendURLS: "http://localhost:8000,http://localhost:3000",
			WebBaseURL:   "https://example-frontend.com",
			CMSBaseURL:   "https://cms.example.com",
			APIBaseURL:   "http://localhost:8080",
			UploadPath:   t.TempDir(),
		},
	}

	// Create the service, injecting the manual mock which satisfies the interface
	service := services.NewCMSPartnerPageService(
		repo,
		mockEmailService, // Pass the manual mock here
		emailContentRepo,
		emailCategoryRepo,
		testCfg,
	)

	return mock, service, mockEmailService, testCfg, cleanup
}

func TestPartnerPageLifecycle(t *testing.T) {
	mock, service, _, _, cleanup := setupPartnerPageTest(t)
	defer cleanup()

	var createdPageID, createdContentID, createdMetaTagID, duplicatedContentID uuid.UUID

	t.Run("1_CreatePartnerPage_Success", func(t *testing.T) {
		t.Log("===> Start: 1_CreatePartnerPage_Success")
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := mockPartnerPage.Contents[0]

		createdPageID = uuid.New()
		createdContentID = uuid.New()
		createdMetaTagID = uuid.New()

		// Uniqueness checks
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1`)).
			WithArgs(partnerContent.URL).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1`)).
			WithArgs(partnerContent.URLAlias).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Transaction and inserts
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "partner_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdPageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdMetaTagID))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}).AddRow(uuid.New(), uuid.New()))
		mock.ExpectCommit()

		_, err := service.CreatePartnerPage(mockPartnerPage)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdatePartnerContent_Success", func(t *testing.T) {
		t.Log("===> Start: 2_UpdatePartnerContent_Success")
		updatedContent := helpers.InitializeMockPartnerPage().Contents[0]
		updatedContent.Title = "Updated Common Title"
		updatedContent.URLAlias += "-updated"
		updatedContent.URL += "-updated"

		newContentID := uuid.New()

		// --- Mock DB interactions ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(createdContentID, createdPageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1 AND page_id != $2`)).
			WithArgs(updatedContent.URL, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(updatedContent.URLAlias, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, enums.PageLanguageEN))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}).AddRow(newContentID, uuid.New()))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WithArgs(sqlmock.AnyArg(), createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 AND "partner_contents"."id" = $2 ORDER BY "partner_contents"."id" LIMIT $3`)).
			WithArgs(newContentID, newContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title"}).AddRow(newContentID, createdPageID, updatedContent.Title))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories" WHERE "partner_content_categories"."partner_content_id" = $1`)).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."partner_content_id" = $1`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."partner_content_id" = $1`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		// --- Act ---
		resp, err := service.UpdatePartnerContent(updatedContent, createdContentID)

		// --- Assert ---
		require.NoError(t, err)
		assert.Equal(t, "Updated Common Title", resp.Title)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("3_FindPartnerPageById_Success", func(t *testing.T) {
		t.Log("===> Start: 3_FindPartnerPageById_Success")
		mockTime := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_pages" WHERE id = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(createdPageID, mockTime, mockTime))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE "partner_contents"."page_id" = $1 AND (partner_contents.mode != $2 AND partner_contents.mode != $3) ORDER BY partner_contents.created_at DESC`)).
			WithArgs(createdPageID, "Histories", "Preview").
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).AddRow(createdContentID, createdPageID, createdMetaTagID))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories" WHERE "partner_content_categories"."partner_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."partner_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(createdMetaTagID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(createdMetaTagID, "mock title", "mock desc"))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."partner_content_id" = $1`)).
			WithArgs(createdContentID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}))

		_, err := service.FindPartnerPageById(createdPageID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("4_FindContentByPartnerPageId_Success", func(t *testing.T) {
		t.Log("===> Start: 4_FindContentByPartnerPageId_Success")
		lang := string(enums.PageLanguageEN)
		mode := string(enums.PageModeDraft)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "partner_contents"."id" LIMIT $4`)).
			WithArgs(createdPageID, lang, mode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language", "mode"}).AddRow(createdContentID, createdPageID, lang, mode))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}))

		_, err := service.FindContentByPartnerPageId(createdPageID, lang, mode)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("5_DuplicatePartnerContentToAnotherLanguage_Success", func(t *testing.T) {
		t.Log("===> Start: 5_DuplicatePartnerContentToAnotherLanguage_Success")
		newRev := helpers.InitializeMockRevision()
		duplicatedContentID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "language", "page_id", "meta_tag_id"}).AddRow(createdContentID, "en", createdPageID, createdMetaTagID))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).WithArgs(createdMetaTagID).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(createdMetaTagID, "mock title", "mock desc"))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(duplicatedContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		_, err := service.DuplicatePartnerContentToAnotherLanguage(createdContentID, newRev)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("6_DuplicatePartnerPage_Success", func(t *testing.T) {
		t.Log("===> Start: 6_DuplicatePartnerPage_Success")
		require.NotEqual(t, uuid.Nil, createdPageID, "A page must have been created in a previous step")

		newDuplicatedPageID := uuid.New()
		newDuplicatedContentID := uuid.New()
		newDuplicatedMetaTagID := uuid.New()
		newDuplicatedRevisionID := uuid.New()
		newDuplicatedComponentID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_pages" WHERE id = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdPageID))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE "partner_contents"."page_id" = $1 AND partner_contents.mode != $2 ORDER BY partner_contents.created_at DESC`)).
			WithArgs(createdPageID, "Histories").WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).AddRow(createdContentID, createdPageID, createdMetaTagID))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"category_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).AddRow(uuid.New(), createdContentID))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).WithArgs(createdMetaTagID).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(createdMetaTagID, "Original Meta Title"))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).WithArgs(createdContentID).WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).AddRow(uuid.New(), createdContentID))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "partner_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedPageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedMetaTagID))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedRevisionID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newDuplicatedComponentID))
		mock.ExpectCommit()

		_, err := service.DuplicatePartnerPage(createdPageID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("7_DeleteContentByPartnerPageId_Success", func(t *testing.T) {
		t.Log("===> Start: 7_DeleteContentByPartnerPageId_Success")

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "partner_contents"."id" LIMIT $4`)).
			WithArgs(createdPageID, "th", "Draft", 1).WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(duplicatedContentID, createdPageID))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components"`)).WithArgs(duplicatedContentID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions"`)).WithArgs(duplicatedContentID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "partner_content_categories"`)).WithArgs(duplicatedContentID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "partner_contents" WHERE "partner_contents"."id" = $1`)).WithArgs(duplicatedContentID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.DeleteContentByPartnerPageId(createdPageID, "th", "Draft")
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("8_DeletePartnerPage_Success", func(t *testing.T) {
		t.Log("===> Start: 8_DeletePartnerPage_Success")

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_pages" WHERE id = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WithArgs(createdPageID, 1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdPageID))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE page_id = $1`)).
			WithArgs(createdPageID).WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(createdContentID, createdPageID))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE partner_content_id IN ($1)`)).
			WithArgs(createdContentID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "partner_content_categories" WHERE partner_content_id IN ($1)`)).
			WithArgs(createdContentID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions" WHERE partner_content_id IN ($1)`)).
			WithArgs(createdContentID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "partner_contents" WHERE page_id = $1`)).
			WithArgs(createdPageID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "partner_pages" WHERE id = $1`)).
			WithArgs(createdPageID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.DeletePartnerPage(createdPageID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPartnerPageRevertFlow(t *testing.T) {
	mock, service, _, _, cleanup := setupPartnerPageTest(t)
	defer cleanup()

	var pageID, contentV1ID, contentV2ID, revisionV1ID uuid.UUID
	var categoryID, metaTagV1ID uuid.UUID

	t.Run("1_CreateInitialPage_V1", func(t *testing.T) {
		t.Log("===> Start: 1_CreateInitialPage_V1")

		// --- Arrange ---
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		mockPartnerPage.Contents[0].Title = "Version 1"
		partnerContent := mockPartnerPage.Contents[0]

		pageID = uuid.New()
		contentV1ID = uuid.New()
		revisionV1ID = uuid.New()
		metaTagV1ID = uuid.New()
		categoryID = uuid.New()
		newComponentID := uuid.New()
		newCategoryTypeID := uuid.New()

		// --- Mock ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1`)).
			WithArgs(partnerContent.URL).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1`)).
			WithArgs(partnerContent.URLAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "partner_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(metaTagV1ID))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentV1ID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(revisionV1ID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newComponentID))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newCategoryTypeID))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))
		mock.ExpectQuery(`INSERT INTO "partner_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}).AddRow(contentV1ID, categoryID))
		mock.ExpectCommit()

		// --- Act ---
		_, err := service.CreatePartnerPage(mockPartnerPage)

		// --- Assert ---
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdateContent_To_V2", func(t *testing.T) {
		t.Log("===> Start: 2_UpdateContent_To_V2")
		updatedContent := helpers.InitializeMockPartnerPage().Contents[0]
		updatedContent.Title = "Version 2"

		contentV2ID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(contentV1ID, pageID, "en"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1 AND page_id != $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Mock Transaction
		mock.ExpectBegin()
		// หา Content V1 อีกครั้งใน Transaction
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(contentV1ID, pageID, enums.PageLanguageEN))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentV2ID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id"}))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		_, err := service.UpdatePartnerContent(updatedContent, contentV1ID)
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
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).AddRow(revisionV1ID, contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(contentV1ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "meta_tag_id", "language"}).
				AddRow(contentV1ID, pageID, "Version 1", metaTagV1ID, enums.PageLanguageEN))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories" WHERE "partner_content_categories"."partner_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}).AddRow(contentV1ID, categoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(categoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(categoryID, "Mock Category"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."partner_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).AddRow(uuid.New(), contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WithArgs(metaTagV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(metaTagV1ID, "Meta Title V1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."partner_content_id" = $1`)).
			WithArgs(contentV1ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).AddRow(revisionV1ID, contentV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE page_id = $1 AND language = $2 ORDER BY created_at DESC,"partner_contents"."id" LIMIT $3`)).
			WithArgs(pageID, enums.PageLanguageEN, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).AddRow(contentV2ID, pageID))

		// ----- Transaction: Revert -----
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(metaTagV1ID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "partner_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newContentV3ID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Mock Category", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "partner_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}))

		mock.ExpectCommit()

		// --- Act ---
		revertedContent, err := service.RevertPartnerContent(revisionV1ID, revertAuthor)

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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_pages" WHERE id = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WithArgs(pageID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE "partner_contents"."page_id" = $1 AND language = $2`)).
			WithArgs(pageID, language).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "mode"}).
				AddRow(contentV1ID, pageID, "Histories").
				AddRow(contentV2ID, pageID, "Histories"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."partner_content_id" IN ($1,$2)`)).
			WithArgs(contentV1ID, contentV2ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id", "author", "created_at"}).
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

func TestPartnerPagePreviewFlow(t *testing.T) {
	mock, service, _, testCfg, cleanup := setupPartnerPageTest(t)
	defer cleanup()

	var pageID uuid.UUID
	var previewContentID uuid.UUID

	t.Run("1_CreateBasePageForPreview", func(t *testing.T) {
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := mockPartnerPage.Contents[0]
		pageID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1`)).
			WithArgs(partnerContent.URL).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1`)).
			WithArgs(partnerContent.URLAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "partner_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id"}))

		mock.ExpectCommit()

		_, err := service.CreatePartnerPage(mockPartnerPage)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("2_PreviewContent_FirstTime_CreatesNew", func(t *testing.T) {
		t.Log("===> Start: 2_PreviewContent_FirstTime_CreatesNew")
		previewContent := helpers.InitializeMockPartnerPage().Contents[0]
		previewContent.Language = enums.PageLanguageEN

		// --- Mock ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1 AND page_id != $2`)).
			WithArgs(previewContent.URL, pageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(previewContent.URLAlias, pageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "partner_contents"."id" LIMIT $4`)).
			WithArgs(pageID, "en", "Preview", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		previewContentID = uuid.New()
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(previewContentID))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		// --- Act ---
		previewURL, err := service.PreviewPartnerContent(pageID, previewContent)

		// --- Assert ---
		require.NoError(t, err)
		frontendURLs := strings.Split(testCfg.App.FrontendURLS, ",")
		appURL := frontendURLs[1]
		expectedURL, err := helpers.BuildPreviewURL(
			appURL,
			string(previewContent.Language), // หรือ "en"
			"partner",
			previewContentID,
		)
		assert.Equal(t, expectedURL, previewURL)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("3_PreviewContent_SecondTime_UpdatesExisting", func(t *testing.T) {
		t.Log("===> Start: 3_PreviewContent_SecondTime_UpdatesExisting")
		updatedPreviewContent := helpers.InitializeMockPartnerPage().Contents[0]
		updatedPreviewContent.Language = enums.PageLanguageEN
		updatedPreviewContent.Title = "Updated Preview Title"

		// --- Mock ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1 AND page_id != $2`)).
			WithArgs(updatedPreviewContent.URL, pageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(updatedPreviewContent.URLAlias, pageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		existingPreview := *updatedPreviewContent
		existingPreview.ID = previewContentID
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 ORDER BY "partner_contents"."id" LIMIT $4`)).
			WithArgs(pageID, "en", "Preview", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(existingPreview.ID))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE partner_content_id = $1`)).
			WithArgs(previewContentID).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()
		// --- Act ---
		previewURL, err := service.PreviewPartnerContent(pageID, updatedPreviewContent)

		// --- Assert ---
		require.NoError(t, err)
		frontendURLs := strings.Split(testCfg.App.FrontendURLS, ",")
		appURL := frontendURLs[1]
		expectedURL, err := helpers.BuildPreviewURL(
			appURL,
			"en",
			"partner",
			previewContentID,
		)
		assert.Equal(t, expectedURL, previewURL)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPartnerPageCategoryFlow(t *testing.T) {
	mock, service, _, _, cleanup := setupPartnerPageTest(t)
	defer cleanup()

	var pageID, contentID, partnerCategoryID, keywordCategoryID uuid.UUID

	const keywordCategoryTypeCode = "category_keywords"

	t.Run("1_CreatePageWithMixedCategories", func(t *testing.T) {
		t.Log("===> Start: 1_CreatePageWithMixedCategories")

		mockPartnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := mockPartnerPage.Contents[0]

		keywordCategoryType := &models.CategoryType{Name: "Keywords", TypeCode: keywordCategoryTypeCode, IsActive: true}
		keywordCategory := &models.Category{Name: "General Keywords", CategoryType: keywordCategoryType, LanguageCode: enums.PageLanguageEN}
		mockPartnerPage.Contents[0].Categories = append(mockPartnerPage.Contents[0].Categories, keywordCategory)

		pageID = uuid.New()
		contentID = uuid.New()
		partnerCategoryID = uuid.New()
		keywordCategoryID = uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1`)).
			WithArgs(partnerContent.URL).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1`)).
			WithArgs(partnerContent.URLAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "partner_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(pageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New())) // for 'category_keywords'
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(keywordCategoryID))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "partner_content_categories"`)).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}))

		mock.ExpectCommit()

		_, err := service.CreatePartnerPage(mockPartnerPage)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_FindCategoriesByPartnerTypeCode_Success", func(t *testing.T) {
		t.Log("===> Start: 2_FindCategoriesByPartnerTypeCode_Success")
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModeDraft)

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT "id" FROM "partner_contents" WHERE page_id = $1 AND language = $2 AND mode = $3 LIMIT $4`,
		)).WithArgs(pageID, language, mode, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contentID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "category_id" FROM "partner_content_categories" WHERE Partner_content_id = $1`)).
			WithArgs(contentID).
			WillReturnRows(sqlmock.NewRows([]string{"category_id"}).AddRow(partnerCategoryID).AddRow(keywordCategoryID))

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

func TestPartnerPagEmailLifecycle(t *testing.T) {
	t.Log("===> Start: TestPartnerPagEmailLifecycle")
	mock, service, mockEmailService, _, cleanup := setupPartnerPageTest(t)
	defer cleanup()

	var createdPageID uuid.UUID
	var createdContentID uuid.UUID

	t.Run("1_CreatePartnerPage_Success", func(t *testing.T) {
		t.Log("===> Start: 1_CreatePartnerPage_Success")
		mockPartnerPage := helpers.InitializeMockPartnerPage()
		partnerContent := mockPartnerPage.Contents[0]

		newPartnerPageID := uuid.New()
		newPartnerContentID := uuid.New()
		newMetaTagID := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1`)).
			WithArgs(partnerContent.URL).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1`)).
			WithArgs(partnerContent.URLAlias).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "partner_pages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newPartnerPageID))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newMetaTagID))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newPartnerContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "category_id"}).AddRow(uuid.New(), uuid.New()))
		mock.ExpectCommit()

		resp, err := service.CreatePartnerPage(mockPartnerPage)
		require.NoError(t, err)
		require.NotNil(t, resp)

		createdPageID = newPartnerPageID
		createdContentID = newPartnerContentID

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_UpdatePartnerContent_AndTriggerEmail_Success", func(t *testing.T) {
		t.Log("===> Start: 2_UpdatePartnerContent_AndTriggerEmail_Success")

		updatedContent := helpers.InitializeMockPartnerPage().Contents[0]
		updatedContent.Title = "Updated Common Title"
		updatedContent.URLAlias += "-updated"
		updatedContent.WorkflowStatus = enums.WorkflowWaitingDesign
		updatedContent.ApprovalEmail = []string{"admin1@example.com"}

		newContentID := uuid.New()

		// --- Mock DB interactions ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, enums.PageLanguageEN))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url = $1 AND page_id != $2`)).
			WithArgs(updatedContent.URL, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "partner_contents" WHERE url_alias = $1 AND page_id != $2`)).
			WithArgs(updatedContent.URLAlias, createdPageID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 ORDER BY "partner_contents"."id" LIMIT $2`)).
			WithArgs(createdContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language"}).AddRow(createdContentID, createdPageID, enums.PageLanguageEN))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_contents" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_contents"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newContentID))
		mock.ExpectQuery(`INSERT INTO "revisions"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "components"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "category_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "partner_content_categories"`).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id"}))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "partner_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WithArgs(sqlmock.AnyArg(), createdPageID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// --- Mock a final preload ---
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents" WHERE id = $1 AND "partner_contents"."id" = $2 ORDER BY "partner_contents"."id" LIMIT $3`)).
			WithArgs(newContentID, newContentID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "language", "workflow_status", "approval_email", "url_alias"}).
				AddRow(newContentID, createdPageID, updatedContent.Title, updatedContent.Language, updatedContent.WorkflowStatus, updatedContent.ApprovalEmail, updatedContent.URLAlias))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).WillReturnRows(sqlmock.NewRows([]string{"partner_content_id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
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
		resp, err := service.UpdatePartnerContent(updatedContent, createdContentID)

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
