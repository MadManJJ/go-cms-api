package tests

import (
	"regexp"
	"testing"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models/enums"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCMSRepo_CreateLandingPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)

	t.Run("successfully create landing page", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()

		mock.ExpectBegin()

		// Expect LandingPage insert
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)

		// Expect MetaTag insert
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)

		// Expect LandingContent insert
		mock.ExpectQuery(`INSERT INTO "landing_contents"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)

		// Expect Revision insert
		mock.ExpectQuery(`INSERT INTO "revisions"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)

		// Expect Components insert
		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)

		// Expect CategoryType insert
		mock.ExpectQuery(`INSERT INTO "category_types"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)

		// Expect Category insert
		mock.ExpectQuery(`INSERT INTO "categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)

		// Expect Content Category insert
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
					AddRow(uuid.New(), uuid.New()), // match what's RETURNED
			)

		mock.ExpectCommit()

		landingPage, err := cmsLandingPageRepo.CreateLandingPage(mockLandingPage)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, mockLandingPage, landingPage)
	})

	t.Run("failed to create landing page", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		mock.ExpectBegin()

		// Expect LandingPage insert
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		landingPage, err := cmsLandingPageRepo.CreateLandingPage(mockLandingPage)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.Nil(t, landingPage)
	})
}

func TestCMSRepo_FindAllLandingPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	
		
	mockQuery := dto.LandingPageQuery{
		Title:            "Reset Password",
		CategoryKeywords: "security",
		Status:           "Published",
	}
	sort := "created_at:DESC"
	page := 1
	limit := 10
	language := "en"

	t.Run("Successfully find all landing pages", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(DISTINCT("landing_pages"."id")) FROM "landing_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))
						
		landingPageId := uuid.New()
		// Mock the INTERSECT subquery
		mock.ExpectQuery(`SELECT DISTINCT landing_pages.* FROM "landing_pages"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(landingPageId, time.Now(), time.Now()))

		landingContentId := uuid.New()
		metaTagId := uuid.New()
		// Mock preload queries for Contents
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id","meta_tag_id"}).
				AddRow(landingContentId, landingPageId, metaTagId))
	
		categoryId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(landingContentId, categoryId))			
				
		categoryTypeId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(categoryId, categoryTypeId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryTypeId))	
		
		componentId1 := uuid.New()
		componentId2 := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(componentId1, landingContentId).
				AddRow(componentId2, landingContentId))		

		landingContentFileId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_files"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(landingContentFileId, landingContentId))				
				
		// Mock preload query for MetaTag
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))				

		revisionId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(revisionId, landingContentId))				

		landingPages, totalCount, err := cmsLandingPageRepo.FindAllLandingPage(mockQuery, sort, page, limit, language)
		
		assert.NoError(t, err)
		assert.NotNil(t, landingPages)
		assert.Equal(t, int64(1), totalCount)
		assert.Len(t, landingPages, 1)
		assert.NotNil(t, landingPages[0].Contents)
		assert.Len(t, landingPages[0].Contents, 1)
		assert.Len(t, landingPages[0].Contents[0].Components, 2)
		assert.NotNil(t, landingPages[0].Contents[0].MetaTag)
		
		// Check the loaded content
		content := landingPages[0].Contents[0]
		assert.Equal(t, landingContentId, content.ID)
		
		// Check nested relationships
		assert.NotNil(t, content.Revision)
		assert.NotNil(t, content.Categories)
		assert.NotNil(t, content.Components)
		assert.NotNil(t, content.MetaTag)		

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find all landing pages", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(DISTINCT("landing_pages"."id")) FROM "landing_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))
						
		// Mock the INTERSECT subquery
		mock.ExpectQuery(`SELECT DISTINCT landing_pages.* FROM "landing_pages"`).
			WillReturnError(errs.ErrInternalServerError)

		landingPages, totalCount, err := cmsLandingPageRepo.FindAllLandingPage(mockQuery, sort, page, limit, language)
		
		assert.Error(t, err)
		assert.Nil(t, landingPages)
		assert.Equal(t, int64(0), totalCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindLandingPageById(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)

	pageId := uuid.New()
	contentId := uuid.New()
	categoryId := uuid.New()
	categoryTypeId := uuid.New()
	componentId1 := uuid.New()
	componentId2 := uuid.New()
	metaTagId := uuid.New()
	revisionId := uuid.New()

	contentTitle := "some title"
	now := time.Now()

	t.Run("successfully find landing page by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(pageId, now, now))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "meta_tag_id"}).
				AddRow(contentId, pageId, contentTitle, metaTagId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(contentId, categoryId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(categoryId, categoryTypeId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryTypeId))			
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(componentId1, contentId).
				AddRow(componentId2, contentId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))							
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(revisionId, contentId))					
		
		landingPage, err := cmsLandingPageRepo.FindLandingPageById(pageId)
		assert.NoError(t, err)
		assert.Equal(t, landingPage.Contents[0].PageID, pageId)
		assert.Equal(t, len(landingPage.Contents), 1)
		assert.Equal(t, len(landingPage.Contents[0].Components), 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find landing page by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnError(errs.ErrInternalServerError)		
		
		landingPage, err := cmsLandingPageRepo.FindLandingPageById(pageId)
		assert.Error(t, err)
		assert.Nil(t, landingPage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_UpdateLandingContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)

	t.Run("successfully update landing content", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockLandingContent := mockLandingPage.Contents[0]

		pageId := uuid.New()
		mockLandingContent.PageID = pageId
		prevContentId := uuid.New()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT for the previous content
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "mode"}).
				AddRow(prevContentId, mockLandingContent.PageID, enums.PageModePublished))

		// Expect UPDATE (setting previous mode to history)
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect INSERT for the metatag
		metatagId := uuid.New()
		mockLandingContent.MetaTagID = metatagId
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metatagId))			

		// Expect INSERT for the new content
		landingContentId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(landingContentId))		

    // Expect INSERT for the new revision
		revisionId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(revisionId))		
				

		// Expect Components insert
		componentId := uuid.New()
		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(componentId), // or any UUID
			)					
				
		// Expect INSERT for the new category_types
		categoryTypeId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryTypeId))						

		// Expect INSERT for the new categories
		category := mockLandingContent.Categories[0]
		category.CategoryTypeID = categoryTypeId
		categoryId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryId))				
				
		// Expect INSERT for the new landing_content_categories
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(landingContentId, categoryId))				
				
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WillReturnResult(sqlmock.NewResult(1, 1))						

		// Expect commit
		mock.ExpectCommit()		

		updatedLandingContent, err := cmsLandingPageRepo.UpdateLandingContent(mockLandingContent, prevContentId)
		assert.Equal(t, updatedLandingContent, mockLandingContent)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("successfully update landing content", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockLandingContent := mockLandingPage.Contents[0]

		pageId := uuid.New()
		mockLandingContent.PageID = pageId
		prevContentId := uuid.New()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT for the previous content
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()		

		updatedLandingContent, err := cmsLandingPageRepo.UpdateLandingContent(mockLandingContent, prevContentId)
		assert.Nil(t, updatedLandingContent)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_DeleteLandingPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	

	t.Run("successfully delete landing content", func(t *testing.T) {
		now := time.Now()
		pageId := uuid.New()
		contentId := uuid.New()
		contentTitle := "some title"

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(pageId, now, now))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title"}).
				AddRow(contentId, pageId, contentTitle))		
				
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_content_categories"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
			
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_content_files"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))					

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			
			
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_pages"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))					

		mock.ExpectCommit()

		err := cmsLandingPageRepo.DeleteLandingPage(pageId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to delete landing content", func(t *testing.T) {
		pageId := uuid.New()

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsLandingPageRepo.DeleteLandingPage(pageId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindContentByLandingPageId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	

	t.Run("successfully find landing content by page Id", func(t *testing.T) {
		pageId := uuid.New()
		contentId := uuid.New()
		categoryId := uuid.New()
		categoryTypeId := uuid.New()
		componentId1 := uuid.New()
		componentId2 := uuid.New()
		revisionId := uuid.New()
		metaTagId := uuid.New()

		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]

		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "language", "mode", "meta_tag_id"}).
				AddRow(contentId, pageId, mockContent.Title, language, mode, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(contentId, categoryId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(categoryId, categoryTypeId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryTypeId))	

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(componentId1, contentId).		
				AddRow(componentId2, contentId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(revisionId, contentId))					
										
		landingContent, err := cmsLandingPageRepo.FindContentByLandingPageId(pageId, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, landingContent.Title, mockContent.Title)
		assert.Equal(t, len(landingContent.Components), 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to find landing content by page Id", func(t *testing.T) {
		pageId := uuid.New()

		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		landingContent, err := cmsLandingPageRepo.FindContentByLandingPageId(pageId, language, mode)
		assert.Error(t, err)		
		assert.Nil(t, landingContent)		
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindLatestLandingContentByPageId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	

	t.Run("successfully find latest content by page id", func(t *testing.T) {
		pageId := uuid.New()
		contentId := uuid.New()
		contentTitle := "some title"
		contentLanguage := string(enums.PageLanguageEN)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "language"}).
			AddRow(contentId, pageId, contentTitle, contentLanguage))

		landingContent, err := cmsLandingPageRepo.FindLatestContentByPageId(pageId, contentLanguage)
		assert.NoError(t, err)
		assert.Equal(t, landingContent.ID , contentId)
		assert.Equal(t, landingContent.PageID , pageId)
		assert.Equal(t, landingContent.Title , contentTitle)
		assert.Equal(t, landingContent.Language , enums.PageLanguageEN)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find latest content by page id", func(t *testing.T) {
		pageId := uuid.New()
		contentLanguage := string(enums.PageLanguageEN)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		landingContent, err := cmsLandingPageRepo.FindLatestContentByPageId(pageId, contentLanguage)
		assert.Error(t, err)
		assert.Nil(t, landingContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		
}

func TestCMSRepo_DeleteLandingContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)

	t.Run("successfully delete landing content by id", func(t *testing.T) {
		pageId := uuid.New()
		contentId := uuid.New()
		contentTitle := "some title"
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title"}).
				AddRow(contentId, pageId, contentTitle))	

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_content_categories"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "landing_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			

		mock.ExpectCommit()

		err := cmsLandingPageRepo.DeleteLandingContent(pageId, language, mode)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
	t.Run("successfully delete landing content by id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsLandingPageRepo.DeleteLandingContent(pageId, language, mode)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		
}

func TestCMSRepo_DuplicateLandingPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)

	t.Run("successfully duplicate landing page", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		mockContent := mockLandingPage.Contents[0]
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		// Generate test IDs
		landingContentId := uuid.New()
		categoryId := uuid.New()
		categoryTypeId := uuid.New()
		componentId1 := uuid.New()
		componentId2 := uuid.New()
		revisionId := uuid.New()
		metaTagId := uuid.New()

		// Mock preload query for Contents (ordered by created_at DESC)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "page_id", "title", "meta_tag_id",
			}).AddRow(
				landingContentId, pageId, mockContent.Title, metaTagId,
			))

		// Mock preload query for Categories (many-to-many relationship)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(landingContentId, categoryId))	
				
		// Mock preload query for the actual Categories
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(categoryId, categoryTypeId))	
				
		// Mock preload query for CategoryType
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryTypeId))				

		// Mock preload query for Components
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(componentId1, landingContentId).
				AddRow(componentId2, landingContentId))

		// Mock preload query for MetaTag
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		// Mock preload query for Revision
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(revisionId, landingContentId))		

		newContentId := uuid.New()
		newPageId := uuid.New()
		newMetaTagId := uuid.New()
		newRevisionId := uuid.New()
		newComponentId1 := uuid.New()
		newComponentId2 := uuid.New()

		mock.ExpectBegin()
		// Expect LandingPage insert
		mock.ExpectQuery(`INSERT INTO "landing_pages"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(newPageId),
			)		

		// Expect Metatag insert
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(newMetaTagId),
			)				

		mock.ExpectQuery(`INSERT INTO "landing_contents"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(newContentId),
			)		
			
		mock.ExpectQuery(`INSERT INTO "revisions"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(newRevisionId),
			)			

		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(newComponentId1).
					AddRow(newComponentId2),
			)		
			
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
			
		mock.ExpectQuery(`INSERT INTO "landing_content_categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
					AddRow(newContentId, categoryId),
			)				
			
		mock.ExpectCommit()

		landingPage, err := cmsLandingPageRepo.DuplicateLandingPage(pageId)
		assert.NoError(t, err)
		assert.Equal(t, len(landingPage.Contents), 1)
		assert.Equal(t, len(landingPage.Contents[0].Components), 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to duplicate landing page", func(t *testing.T) {
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnError(errs.ErrInternalServerError)

		landingPage, err := cmsLandingPageRepo.DuplicateLandingPage(pageId)
		assert.Error(t, err)
		assert.Nil(t, landingPage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_DuplicateLandingContentToAnotherLanguage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)
	
	t.Run("successfully duplicate landing content to another language", func(t *testing.T) {
		oldPageId := uuid.New()
		oldContentId := uuid.New()
		oldCategoryId := uuid.New()
		oldCategoryTypeId := uuid.New()
		oldComponentId1 := uuid.New()
		oldComponentId2 := uuid.New()
		oldRevisionId := uuid.New()

		newContentId := uuid.New()
		newRevisionId := uuid.New()
		newComponentId1 := uuid.New()
		newComponentId2 := uuid.New()
		newCategoryId := uuid.New()

		mockLandingPage := helpers.InitializeMockLandingPage()

		content := mockLandingPage.Contents[0]
		newRevision := helpers.InitializeMockRevision()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "pageId", "title"}).
				AddRow(oldContentId, oldPageId, content.Title))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(oldContentId, oldCategoryId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(oldCategoryId, oldCategoryTypeId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(oldComponentId1, oldContentId).
				AddRow(oldComponentId2, oldContentId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(oldRevisionId, oldContentId))						
				
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newContentId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newRevisionId))			
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newComponentId1).
				AddRow(newComponentId2))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newCategoryId))				
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(newContentId, newCategoryId))						

		mock.ExpectCommit()

		landingContent, err := cmsLandingPageRepo.DuplicateLandingContentToAnotherLanguage(oldContentId, newRevision)
		assert.NoError(t, err)
		assert.Equal(t, len(landingContent.Components), 2)
		assert.Equal(t, landingContent.ID, newContentId)
		assert.Equal(t, landingContent.Revision.ID, newRevisionId)
		assert.NotEqual(t, landingContent.ID, oldContentId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("successfully duplicate landing content to another language", func(t *testing.T) {
		oldContentId := uuid.New()

		newRevision := helpers.InitializeMockRevision()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)				

		landingContent, err := cmsLandingPageRepo.DuplicateLandingContentToAnotherLanguage(oldContentId, newRevision)
		assert.Error(t, err)
		assert.Nil(t, landingContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		
}

func TestCMSRepo_RevertLandingContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	
	t.Run("successfully revert landing content", func(t *testing.T) {
		oldRevisionId := uuid.New()
		oldContentId := uuid.New()
		oldPageId := uuid.New()
		oldCategoryId := uuid.New()
		oldCategoryTypeId := uuid.New()
		oldComponentId1 := uuid.New()
		oldComponentId2 := uuid.New()

		newContentId := uuid.New()
		newRevisionId := uuid.New()
		newComponentId1 := uuid.New()
		newComponentId2 := uuid.New()
		newCategoryId := uuid.New()

		newRevision := helpers.InitializeMockRevision()
		mockLandingPage := helpers.InitializeMockLandingPage()
		content := mockLandingPage.Contents[0]

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(oldRevisionId, oldContentId))		
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "pageId", "title"}).
				AddRow(oldContentId, oldPageId, content.Title))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(oldContentId, oldCategoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(oldCategoryId, oldCategoryTypeId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(oldComponentId1, oldContentId).
				AddRow(oldComponentId2, oldContentId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(oldRevisionId, oldContentId))				
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).
			AddRow(oldContentId, oldPageId))				
				
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "landing_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newContentId))		

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newRevisionId))			
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newComponentId1).
				AddRow(newComponentId2))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(newCategoryId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "category_id"}).
				AddRow(newContentId, newCategoryId))					

		mock.ExpectCommit()

		landingContent, err := cmsLandingPageRepo.RevertLandingContent(oldRevisionId, newRevision)
		assert.NoError(t, err)
		assert.Equal(t, len(landingContent.Components), 2)
		assert.Equal(t, landingContent.ID, newContentId)
		assert.Equal(t, landingContent.Revision.ID, newRevisionId)
		assert.NotEqual(t, landingContent.ID, oldContentId)		
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("successfully revert landing content", func(t *testing.T) {
		oldRevisionId := uuid.New()

		newRevision := helpers.InitializeMockRevision()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnError(errs.ErrInternalServerError)		

		landingContent, err := cmsLandingPageRepo.RevertLandingContent(oldRevisionId, newRevision)
		assert.Error(t, err)
		assert.Nil(t, landingContent)	
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_GetLandingPageCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	

	t.Run("successfully get categorise", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		content := mockLandingPage.Contents[0]
		category := content.Categories[0]
		categoryType := category.CategoryType

		pageId := uuid.New()
		contentId := uuid.New()
		categoryId1 := uuid.New()
		categoryId2 := uuid.New()
		categoryId3 := uuid.New()
		categoryId4 := uuid.New()
		categoryTypeId1 := uuid.New()
		categoryTypeId2 := uuid.New()
		categoryTypeCode := categoryType.TypeCode
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "category_id" FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"category_id"}).
				AddRow(categoryId1).
				AddRow(categoryId2).			
				AddRow(categoryId3).				
				AddRow(categoryId4))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "categories"."id","categories"."category_type_id","categories"."language_code","categories"."name","categories"."description","categories"."weight","categories"."publish_status","categories"."created_at","categories"."updated_at" FROM "categories" JOIN category_types ON category_types.id = categories.category_type_id`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "category_type_id",
			}).
				AddRow(categoryId1, categoryTypeId1).
				AddRow(categoryId2, categoryTypeId1).
				AddRow(categoryId3, categoryTypeId2).
				AddRow(categoryId4, categoryTypeId2))			

		categories, err := cmsLandingPageRepo.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, len(categories), 4)
		assert.Equal(t, categories[0].CategoryTypeID, categoryTypeId1)
		assert.Equal(t, categories[1].CategoryTypeID, categoryTypeId1)
		assert.Equal(t, categories[2].CategoryTypeID, categoryTypeId2)
		assert.Equal(t, categories[3].CategoryTypeID, categoryTypeId2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to get categories", func(t *testing.T) {
		mockLandingPage := helpers.InitializeMockLandingPage()
		content := mockLandingPage.Contents[0]
		category := content.Categories[0]
		categoryType := category.CategoryType

		pageId := uuid.New()
		categoryTypeCode := categoryType.TypeCode
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		categories, err := cmsLandingPageRepo.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		
}

func TestCMSRepo_GetRevisionByLandingPageId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	
	
	t.Run("successfully get revisions by page id", func(t *testing.T) {
		pageId := uuid.New()
		contentId1 := uuid.New()
		revisionId1 := uuid.New()
		language := string(enums.PageLanguageEN)

		mockLandingPage := helpers.InitializeMockLandingPage()
		content := mockLandingPage.Contents[0]
		revision := content.Revision

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(pageId, time.Now(), time.Now()))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "page_id", "title", "meta_tag_id",
			}).
				AddRow(contentId1, pageId, content.Title, content.MetaTagID))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "landing_content_id", "author", "message", "description",
			}).
				AddRow(revisionId1, contentId1, revision.Author, revision.Message, revision.Description))

		revisions, err := cmsLandingPageRepo.GetRevisionByLandingPageId(pageId, language)
		assert.NoError(t, err)
		assert.Equal(t, len(revisions), 1)
		assert.Equal(t, revisions[0].Author, revision.Author)
		assert.Equal(t, revisions[0].Message, revision.Message)
		assert.Equal(t, revisions[0].Description, revision.Description)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})	

	t.Run("successfully get revisions by page id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_pages"`)).
			WillReturnError(errs.ErrInternalServerError)

		revisions, err := cmsLandingPageRepo.GetRevisionByLandingPageId(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, revisions)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})	
}

func TestCMSRepo_IsLandingUrlAliasDuplicate(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	

	t.Run("url alias is duplicated", func(t *testing.T) {
		urlAlias := "random-url-alias"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))		

		isUrlDuplicate, err := cmsLandingPageRepo.IsUrlAliasDuplicate(urlAlias, pageId)
		assert.NoError(t, err)
		assert.True(t, isUrlDuplicate)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("url alias is not duplicated", func(t *testing.T) {
		urlAlias := "random-url-alias"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(0))		

		isUrlDuplicate, err := cmsLandingPageRepo.IsUrlAliasDuplicate(urlAlias, pageId)
		assert.NoError(t, err)
		assert.False(t, isUrlDuplicate)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	

	t.Run("internal server error", func(t *testing.T) {
		urlAlias := "random-url-alias"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)	

		_, err := cmsLandingPageRepo.IsUrlAliasDuplicate(urlAlias, pageId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	
}

func TestCMSRepo_GetLandingPageIdByContentId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsLandingPageRepo := repo.NewCMSLandingPageRepository(gormDB)	

	t.Run("successfully get page id by content id", func(t *testing.T) {
		contentId := uuid.New()
		expectedPageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).
				AddRow(contentId, expectedPageId))

		actualPageId, err := cmsLandingPageRepo.GetPageIdByContentId(contentId)
		assert.NoError(t, err)
		assert.Equal(t, expectedPageId, actualPageId)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("successfully get page id by content id", func(t *testing.T) {
		contentId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		actualPageId, err := cmsLandingPageRepo.GetPageIdByContentId(contentId)
		assert.Error(t, err)
		assert.Equal(t, actualPageId, uuid.Nil)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})	
}