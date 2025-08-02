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

func TestCMSRepo_CreateFaqPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)

	t.Run("successfully create faq page", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()

		mock.ExpectBegin()

		// Expect FaqPage insert
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)

		// Expect MetaTag insert
		mock.ExpectQuery(`INSERT INTO "meta_tags"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)					

		// Expect FaqContent insert
		mock.ExpectQuery(`INSERT INTO "faq_contents"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)		

		// Expect Revision insert
		mock.ExpectQuery(`INSERT INTO "revisions"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)		

		// Expect Components insert
		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)	

		// Expect CategoryType insert
		mock.ExpectQuery(`INSERT INTO "category_types"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)			

			
		// Expect Category insert
		mock.ExpectQuery(`INSERT INTO "categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), 
			)		


		// Expect Content Category insert
		mock.ExpectQuery(`INSERT INTO "faq_content_categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
					AddRow(uuid.New(), uuid.New()), // match what's RETURNED
			)

		mock.ExpectCommit()

		faqPage, err := cmsFaqPageRepo.CreateFaqPage(mockFaqPage)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, mockFaqPage, faqPage)
	})

	t.Run("failed to create faq page", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mock.ExpectBegin()

		// Expect FaqPage insert
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		faqPage, err := cmsFaqPageRepo.CreateFaqPage(mockFaqPage)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Error(t, err)
		assert.Nil(t, faqPage)
	})	
}

func TestCMSRepo_FindAllFaqPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
		
	mockQuery := dto.FaqPageQuery{
		Title:            "Reset Password",
		CategoryFaq:      "Account",
		CategoryKeywords: "security",
		Status:           "Published",
	}
	sort := "created_at:DESC"
	page := 1
	limit := 10
	language := "en"

	t.Run("Successfully find all faq pages", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(DISTINCT("faq_pages"."id")) FROM "faq_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))

		faqPageId := uuid.New()
		// Mock the INTERSECT subquery
		mock.ExpectQuery(`SELECT DISTINCT faq_pages.* FROM "faq_pages"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(faqPageId, time.Now(), time.Now()))

		faqContentId := uuid.New()
		metaTagId := uuid.New()
		// Mock preload queries for Contents
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id","meta_tag_id"}).
				AddRow(faqContentId, faqPageId, metaTagId))
	
		categoryId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(faqContentId, categoryId))			
				
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
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(componentId1, faqContentId).
				AddRow(componentId2, faqContentId))		
				
		// Mock preload query for MetaTag
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))				

		revisionId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, faqContentId))				

		faqPages, totalCount, err := cmsFaqPageRepo.FindAllFaqPage(mockQuery, sort, page, limit, language)
		
		assert.NoError(t, err)
		assert.NotNil(t, faqPages)
		assert.Equal(t, int64(1), totalCount)
		assert.Len(t, faqPages, 1)
		assert.NotNil(t, faqPages[0].Contents)
		assert.Len(t, faqPages[0].Contents, 1)
		assert.Len(t, faqPages[0].Contents[0].Components, 2)
		assert.NotNil(t, faqPages[0].Contents[0].MetaTag)
		
		// Check the loaded content
		content := faqPages[0].Contents[0]
		assert.Equal(t, faqContentId, content.ID)
		
		// Check nested relationships
		assert.NotNil(t, content.Revision)
		assert.NotNil(t, content.Categories)
		assert.NotNil(t, content.Components)
		assert.NotNil(t, content.MetaTag)		

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find all faq pages", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(DISTINCT("faq_pages"."id")) FROM "faq_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))	

		// Mock the INTERSECT subquery
		mock.ExpectQuery(`SELECT DISTINCT faq_pages.* FROM "faq_pages"`).
			WillReturnError(errs.ErrInternalServerError)

		faqPages, totalCount, err := cmsFaqPageRepo.FindAllFaqPage(mockQuery, sort, page, limit, language)
		
		assert.Error(t, err)
		assert.Nil(t, faqPages)
		assert.Equal(t, int64(0), totalCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindFaqPageById(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	

	t.Run("successfully find faq page by id", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		faqPageId := uuid.New()
		mockFaqPage.ID = faqPageId
		
		// Mock the main FaqPage query - remove "contents" column since it's not stored in the table
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(mockFaqPage.ID, mockFaqPage.CreatedAt, mockFaqPage.UpdatedAt))

		// Generate test IDs
		faqContentId := uuid.New()
		categoryId := uuid.New()
		categoryTypeId := uuid.New()
		componentId1 := uuid.New()
		componentId2 := uuid.New()
		revisionId := uuid.New()
		metaTagId := uuid.New()

		// Mock preload query for Contents (ordered by created_at DESC)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents" WHERE "faq_contents"."page_id" = $1 AND (faq_contents.mode != $2 AND faq_contents.mode != $3) ORDER BY faq_contents.created_at DESC`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "page_id", "meta_tag_id",
			}).AddRow(
				faqContentId, mockFaqPage.ID, metaTagId,
			))

		// Mock preload query for Categories (many-to-many relationship)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories" WHERE "faq_content_categories"."faq_content_id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(faqContentId, categoryId))	
				
		// Mock preload query for the actual Categories
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(categoryId, categoryTypeId))	
				
		// Mock preload query for CategoryType
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types" WHERE "category_types"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryTypeId))				

		// Mock preload query for Components
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components" WHERE "components"."faq_content_id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(componentId1, faqContentId).
				AddRow(componentId2, faqContentId))

		// Mock preload query for MetaTag
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags" WHERE "meta_tags"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		// Mock preload query for Revision
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions" WHERE "revisions"."faq_content_id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, faqContentId))

		// Execute the function
		faqPage, err := cmsFaqPageRepo.FindFaqPageById(mockFaqPage.ID)
		
		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, faqPage)
		assert.Equal(t, mockFaqPage.ID, faqPage.ID)
		assert.Equal(t, mockFaqPage.CreatedAt, faqPage.CreatedAt)
		assert.Equal(t, mockFaqPage.UpdatedAt, faqPage.UpdatedAt)

		// Check that Contents were loaded
		assert.NotNil(t, faqPage.Contents)
		assert.Len(t, faqPage.Contents, 1)
		
		// Check the loaded content
		content := faqPage.Contents[0]
		assert.Equal(t, faqContentId, content.ID)
		assert.Equal(t, mockFaqPage.ID, content.PageID)
		
		// Check nested relationships
		assert.NotNil(t, content.Revision)
		assert.NotNil(t, content.Categories)
		assert.NotNil(t, content.Components)
		assert.NotNil(t, content.MetaTag)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed find faq page by id", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		faqPageId := uuid.New()
		mockFaqPage.ID = faqPageId
		
		// Mock the main FaqPage query - remove "contents" column since it's not stored in the table
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnError(errs.ErrInternalServerError)

		// Execute the function
		faqPage, err := cmsFaqPageRepo.FindFaqPageById(mockFaqPage.ID)
		
		// Assertions
		assert.Error(t, err)
		assert.Nil(t, faqPage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_UpdateFaqContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)		

	t.Run("successfully update faq content", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqContent := mockFaqPage.Contents[0]

		pageId := uuid.New()
		mockFaqContent.PageID = pageId
		prevContentId := uuid.New()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT for the previous content
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "mode"}).
				AddRow(prevContentId, mockFaqContent.PageID, enums.PageModePublished))

		// Expect UPDATE (setting previous mode to history)
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect INSERT for the metatag
		metatagId := uuid.New()
		mockFaqContent.MetaTagID = metatagId
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metatagId))			

		// Expect INSERT for the new content
		faqContentId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(faqContentId))		

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
					AddRow(componentId), 
			)					
				
		// Expect INSERT for the new category_types
		categoryTypeId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryTypeId))						

		// Expect INSERT for the new categories
		category := mockFaqContent.Categories[0]
		category.CategoryTypeID = categoryTypeId
		categoryId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(categoryId))				
				
		// Expect INSERT for the new faq_content_categories
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(faqContentId, categoryId))			
				
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_pages" SET "updated_at"=$1 WHERE id = $2`)).
			WillReturnResult(sqlmock.NewResult(1, 1))						

		// Expect commit
		mock.ExpectCommit()		

		updatedFaqContent, err := cmsFaqPageRepo.UpdateFaqContent(mockFaqContent, prevContentId)
		assert.Equal(t, updatedFaqContent, mockFaqContent)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully update faq content", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockFaqContent := mockFaqPage.Contents[0]

		pageId := uuid.New()
		mockFaqContent.PageID = pageId
		prevContentId := uuid.New()

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT for the previous content
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()		

		updatedFaqContent, err := cmsFaqPageRepo.UpdateFaqContent(mockFaqContent, prevContentId)
		assert.Nil(t, updatedFaqContent)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_DeleteFaqPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)		

	t.Run("successfully delete faq content", func(t *testing.T) {
		now := time.Now()
		pageId := uuid.New()
		contentId := uuid.New()
		contentTitle := "some title"

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(pageId, now, now))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title"}).
				AddRow(contentId, pageId, contentTitle))		
				
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "faq_content_categories"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))		

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "faq_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			
			
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "faq_pages"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))					

		mock.ExpectCommit()

		err := cmsFaqPageRepo.DeleteFaqPage(pageId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to delete faq content", func(t *testing.T) {
		pageId := uuid.New()

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsFaqPageRepo.DeleteFaqPage(pageId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindContentByFaqPageId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)		

	t.Run("successfully find faq content by page Id", func(t *testing.T) {
		pageId := uuid.New()
		contentId := uuid.New()
		categoryId := uuid.New()
		categoryTypeId := uuid.New()
		componentId1 := uuid.New()
		componentId2 := uuid.New()
		revisionId := uuid.New()

		mockFaqPage := helpers.InitializeMockFaqPage()
		component := helpers.InitializeMockComponent()

		mockContent := mockFaqPage.Contents[0]
		category := mockContent.Categories[0]
		categoryType := category.CategoryType
		revision := mockContent.Revision

		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "language", "mode"}).
				AddRow(contentId, pageId, mockContent.Title, language, mode))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(contentId, categoryId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id", "language_code", "name"}).
				AddRow(categoryId, categoryTypeId, category.LanguageCode, category.Name))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category_types"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type_code", "name", "is_active"}).
				AddRow(categoryTypeId, categoryType.TypeCode, categoryType.Name, categoryType.IsActive))	

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id", "type", "props"}).
				AddRow(componentId1, contentId, component.Type, component.Props).		
				AddRow(componentId2, contentId, component.Type, component.Props))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id", "publish_status", "author", "message", "description"}).
				AddRow(revisionId, contentId, revision.PublishStatus, revision.Author, revision.Message, revision.Description))					
										
		faqContent, err := cmsFaqPageRepo.FindContentByFaqPageId(pageId, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, faqContent.Title, mockContent.Title)
		assert.Equal(t, len(faqContent.Components), 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find faq content by page Id", func(t *testing.T) {
		pageId := uuid.New()

		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		faqContent, err := cmsFaqPageRepo.FindContentByFaqPageId(pageId, language, mode)
		assert.Error(t, err)		
		assert.Nil(t, faqContent)		
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSRepo_FindLatestFaqContentByPageId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)		

	t.Run("successfully find latest content by page id", func(t *testing.T) {
		pageId := uuid.New()
		contentId := uuid.New()
		contentTitle := "some title"
		contentLanguage := string(enums.PageLanguageEN)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title", "language"}).
			AddRow(contentId, pageId, contentTitle, contentLanguage))

		faqContent, err := cmsFaqPageRepo.FindLatestContentByPageId(pageId, contentLanguage)
		assert.NoError(t, err)
		assert.Equal(t, faqContent.ID , contentId)
		assert.Equal(t, faqContent.PageID , pageId)
		assert.Equal(t, faqContent.Title , contentTitle)
		assert.Equal(t, faqContent.Language , enums.PageLanguageEN)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find latest content by page id", func(t *testing.T) {
		pageId := uuid.New()
		contentLanguage := string(enums.PageLanguageEN)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		faqContent, err := cmsFaqPageRepo.FindLatestContentByPageId(pageId, contentLanguage)
		assert.Error(t, err)
		assert.Nil(t, faqContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_DeleteFaqContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	

	t.Run("successfully delete faq content by id", func(t *testing.T) {
		pageId := uuid.New()
		contentId := uuid.New()
		contentTitle := "some title"
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "title"}).
				AddRow(contentId, pageId, contentTitle))	

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "revisions"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "faq_content_categories"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "faq_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			

		mock.ExpectCommit()

		err := cmsFaqPageRepo.DeleteFaqContent(pageId, language, mode)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully delete faq content by id", func(t *testing.T) {
		pageId := uuid.New()
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		err := cmsFaqPageRepo.DeleteFaqContent(pageId, language, mode)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_DuplicateFaqPage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
	
	t.Run("successfully duplicate faq page", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		mockContent := mockFaqPage.Contents[0]
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		// Generate test IDs
		faqContentId := uuid.New()
		categoryId := uuid.New()
		categoryTypeId := uuid.New()
		componentId1 := uuid.New()
		componentId2 := uuid.New()
		revisionId := uuid.New()
		metaTagId := uuid.New()

		// Mock preload query for Contents (ordered by created_at DESC)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "page_id", "title","meta_tag_id",
			}).AddRow(
				faqContentId, pageId, mockContent.Title, metaTagId,
			))

		// Mock preload query for Categories (many-to-many relationship)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(faqContentId, categoryId))	
				
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
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(componentId1, faqContentId).
				AddRow(componentId2, faqContentId))

		// Mock preload query for MetaTag
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		// Mock preload query for Revision
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, faqContentId))		

		newContentId := uuid.New()
		newPageId := uuid.New()
		newMetaTagId := uuid.New()
		newRevisionId := uuid.New()
		newComponentId1 := uuid.New()
		newComponentId2 := uuid.New()

		mock.ExpectBegin()
		// Expect FaqPage insert
		mock.ExpectQuery(`INSERT INTO "faq_pages"`).
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

		mock.ExpectQuery(`INSERT INTO "faq_contents"`).
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
			
		mock.ExpectQuery(`INSERT INTO "faq_content_categories"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
					AddRow(newContentId, categoryId),
			)				
			
		mock.ExpectCommit()

		faqPage, err := cmsFaqPageRepo.DuplicateFaqPage(pageId)
		assert.NoError(t, err)
		assert.Equal(t, len(faqPage.Contents), 1)
		assert.Equal(t, len(faqPage.Contents[0].Components), 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to duplicate faq page", func(t *testing.T) {
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnError(errs.ErrInternalServerError)

		faqPage, err := cmsFaqPageRepo.DuplicateFaqPage(pageId)
		assert.Error(t, err)
		assert.Nil(t, faqPage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_DuplicateFaqContentToAnotherLanguage(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
	
	t.Run("successfully duplicate faq content to another language", func(t *testing.T) {
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

		mockFaqPage := helpers.InitializeMockFaqPage()

		content := mockFaqPage.Contents[0]
		newRevision := helpers.InitializeMockRevision()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "pageId", "title"}).
				AddRow(oldContentId, oldPageId, content.Title))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(oldContentId, oldCategoryId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(oldCategoryId, oldCategoryTypeId))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(oldComponentId1, oldContentId).
				AddRow(oldComponentId2, oldContentId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(oldRevisionId, oldContentId))		
				
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_contents"`)).
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
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(newContentId, newCategoryId))						

		mock.ExpectCommit()

		faqContent, err := cmsFaqPageRepo.DuplicateFaqContentToAnotherLanguage(oldContentId, newRevision)
		assert.NoError(t, err)
		assert.Equal(t, len(faqContent.Components), 2)
		assert.Equal(t, faqContent.ID, newContentId)
		assert.Equal(t, faqContent.Revision.ID, newRevisionId)
		assert.NotEqual(t, faqContent.ID, oldContentId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully duplicate faq content to another language", func(t *testing.T) {
		oldContentId := uuid.New()

		newRevision := helpers.InitializeMockRevision()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)				

		faqContent, err := cmsFaqPageRepo.DuplicateFaqContentToAnotherLanguage(oldContentId, newRevision)
		assert.Error(t, err)
		assert.Nil(t, faqContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_RevertFaqContent(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
	
	t.Run("successfully revert faq content", func(t *testing.T) {
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
		mockFaqPage := helpers.InitializeMockFaqPage()
		content := mockFaqPage.Contents[0]

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(oldRevisionId, oldContentId))		
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "pageId", "title"}).
				AddRow(oldContentId, oldPageId, content.Title))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(oldContentId, oldCategoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_type_id"}).
				AddRow(oldCategoryId, oldCategoryTypeId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(oldComponentId1, oldContentId).
				AddRow(oldComponentId2, oldContentId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(oldRevisionId, oldContentId))				
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).
			AddRow(oldContentId, oldPageId))					
				
		mock.ExpectBegin()	

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "faq_contents"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))			

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_contents"`)).
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
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "category_id"}).
				AddRow(newContentId, newCategoryId))					

		mock.ExpectCommit()

		faqContent, err := cmsFaqPageRepo.RevertFaqContent(oldRevisionId, newRevision)
		assert.NoError(t, err)
		assert.Equal(t, len(faqContent.Components), 2)
		assert.Equal(t, faqContent.ID, newContentId)
		assert.Equal(t, faqContent.Revision.ID, newRevisionId)
		assert.NotEqual(t, faqContent.ID, oldContentId)		
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully revert faq content", func(t *testing.T) {
		oldRevisionId := uuid.New()

		newRevision := helpers.InitializeMockRevision()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnError(errs.ErrInternalServerError)		

		faqContent, err := cmsFaqPageRepo.RevertFaqContent(oldRevisionId, newRevision)
		assert.Error(t, err)
		assert.Nil(t, faqContent)	
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_GetFaqPageCategory(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
	
	t.Run("successfully get categorise", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		content := mockFaqPage.Contents[0]
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "category_id" FROM "faq_content_categories"`)).
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

		categories, err := cmsFaqPageRepo.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.NoError(t, err)
		assert.Equal(t, len(categories), 4)
		assert.Equal(t, categories[0].CategoryTypeID, categoryTypeId1)
		assert.Equal(t, categories[1].CategoryTypeID, categoryTypeId1)
		assert.Equal(t, categories[2].CategoryTypeID, categoryTypeId2)
		assert.Equal(t, categories[3].CategoryTypeID, categoryTypeId2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get categories", func(t *testing.T) {
		mockFaqPage := helpers.InitializeMockFaqPage()
		content := mockFaqPage.Contents[0]
		category := content.Categories[0]
		categoryType := category.CategoryType

		pageId := uuid.New()
		categoryTypeCode := categoryType.TypeCode
		language := string(enums.PageLanguageEN)
		mode := string(enums.PageModePublished)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		categories, err := cmsFaqPageRepo.GetCategory(pageId, categoryTypeCode, language, mode)
		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_GetRevisionByFaqPageId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
	
	t.Run("successfully get revisions by page id", func(t *testing.T) {
		pageId := uuid.New()
		contentId1 := uuid.New()
		revisionId1 := uuid.New()
		language := string(enums.PageLanguageEN)


		mockFaqPage := helpers.InitializeMockFaqPage()
		content := mockFaqPage.Contents[0]
		revision := content.Revision

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(pageId, time.Now(), time.Now()))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "page_id", "title", "meta_tag_id",
			}).
				AddRow(contentId1, pageId, content.Title, content.MetaTagID))		
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "faq_content_id", "author", "message", "description",
			}).
				AddRow(revisionId1, contentId1, revision.Author, revision.Message, revision.Description))

		revisions, err := cmsFaqPageRepo.GetRevisionByFaqPageId(pageId, language)
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_pages"`)).
			WillReturnError(errs.ErrInternalServerError)

		revisions, err := cmsFaqPageRepo.GetRevisionByFaqPageId(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, revisions)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})	
}

func TestCMSRepo_IsFaqUrlDuplicate(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
		
	t.Run("url is duplicated", func(t *testing.T) {
		url := "random-url"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))		

		isUrlDuplicate, err := cmsFaqPageRepo.IsUrlDuplicate(url, pageId)
		assert.NoError(t, err)
		assert.True(t, isUrlDuplicate)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("url is not duplicated", func(t *testing.T) {
		url := "random-url"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(0))		

		isUrlDuplicate, err := cmsFaqPageRepo.IsUrlAliasDuplicate(url, pageId)
		assert.NoError(t, err)
		assert.False(t, isUrlDuplicate)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	

	t.Run("internal server error", func(t *testing.T) {
		url := "random-url"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)	

		_, err := cmsFaqPageRepo.IsUrlAliasDuplicate(url, pageId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	
}

func TestCMSRepo_IsFaqUrlAliasDuplicate(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
		
	t.Run("url alias is duplicated", func(t *testing.T) {
		urlAlias := "random-url-alias"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))		

		isUrlDuplicate, err := cmsFaqPageRepo.IsUrlAliasDuplicate(urlAlias, pageId)
		assert.NoError(t, err)
		assert.True(t, isUrlDuplicate)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})

	t.Run("url alias is not duplicated", func(t *testing.T) {
		urlAlias := "random-url-alias"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(0))		

		isUrlDuplicate, err := cmsFaqPageRepo.IsUrlAliasDuplicate(urlAlias, pageId)
		assert.NoError(t, err)
		assert.False(t, isUrlDuplicate)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	

	t.Run("internal server error", func(t *testing.T) {
		urlAlias := "random-url-alias"
		pageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)	

		_, err := cmsFaqPageRepo.IsUrlAliasDuplicate(urlAlias, pageId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())				
	})	
}

func TestCMSRepo_GetFaqPageIdByContentId(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	
		
	t.Run("successfully get page id by content id", func(t *testing.T) {
		contentId := uuid.New()
		expectedPageId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id"}).
				AddRow(contentId, expectedPageId))

		actualPageId, err := cmsFaqPageRepo.GetPageIdByContentId(contentId)
		assert.NoError(t, err)
		assert.Equal(t, expectedPageId, actualPageId)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("successfully get page id by content id", func(t *testing.T) {
		contentId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		actualPageId, err := cmsFaqPageRepo.GetPageIdByContentId(contentId)
		assert.Error(t, err)
		assert.Equal(t, actualPageId, uuid.Nil)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})
}

func TestCMSRepo_CreateFaqContentPreview(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	

	mockFaqPage := helpers.InitializeMockFaqPage()
	mockFaqContent := mockFaqPage.Contents[0]
	metaTagId := uuid.New()
	contentId := uuid.New()
	componentId := uuid.New()
	
	t.Run("successfully create faq content preview", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(metaTagId), 
			)		

		mock.ExpectQuery(`INSERT INTO "faq_contents"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(contentId), 
			)			

		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(componentId), 
			)				

		mock.ExpectCommit()

		actualFaqContent, err := cmsFaqPageRepo.CreateFaqContentPreview(mockFaqContent)
		assert.NoError(t, err)
		assert.Nil(t, actualFaqContent.Revision)
		assert.Nil(t, actualFaqContent.Categories)
		assert.NotNil(t, actualFaqContent.Components)
		assert.NotNil(t, actualFaqContent.MetaTag)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("failed to create faq content preview", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "meta_tags"`).
			WillReturnError(errs.ErrInternalServerError)			

		mock.ExpectRollback()

		actualFaqContent, err := cmsFaqPageRepo.CreateFaqContentPreview(mockFaqContent)
		assert.Error(t, err)
		assert.Nil(t, actualFaqContent)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})	
}

func TestCMSRepo_UpdateFaqContentPreview(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)		

	mockFaqPage := helpers.InitializeMockFaqPage()
	mockFaqContent := mockFaqPage.Contents[0]
	metaTagId := uuid.New()
	contentId := uuid.New()
	componentId := uuid.New()
	
	t.Run("successfully update faq content preview", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE faq_content_id = $1`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meta_tags"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(metaTagId),
			)

		mock.ExpectQuery(`INSERT INTO "faq_contents"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(contentId), 
			)			

		mock.ExpectQuery(`INSERT INTO "components"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(componentId), 
			)			

		mock.ExpectCommit()

		actualFaqContent, err := cmsFaqPageRepo.UpdateFaqContentPreview(mockFaqContent)
		assert.NoError(t, err)
		assert.Nil(t, actualFaqContent.Revision)
		assert.Nil(t, actualFaqContent.Categories)
		assert.NotNil(t, actualFaqContent.Components)
		assert.NotNil(t, actualFaqContent.MetaTag)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})

	t.Run("failed to update faq content preview", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "components" WHERE faq_content_id = $1`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		actualFaqContent, err := cmsFaqPageRepo.UpdateFaqContentPreview(mockFaqContent)
		assert.Error(t, err)
		assert.Nil(t, actualFaqContent)
		assert.NoError(t, mock.ExpectationsWereMet())		
	})	
}

func TestCMSRepo_FindFaqContentPreviewById(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsFaqPageRepo := repo.NewCMSFaqPageRepository(gormDB)	

	pageId := uuid.New()
	language := string(enums.PageLanguageEN)
	
	t.Run("successfully get content preview by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "language", "mode"}).
				AddRow(uuid.New(), pageId, language, "Preview"))

		actualFaqContent, err := cmsFaqPageRepo.FindFaqContentPreviewById(pageId, language)
		assert.NoError(t, err)
		assert.NotNil(t, actualFaqContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get content preview by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		actualFaqContent, err := cmsFaqPageRepo.FindFaqContentPreviewById(pageId, language)
		assert.Error(t, err)
		assert.Nil(t, actualFaqContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}