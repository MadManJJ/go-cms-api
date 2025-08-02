package tests

import (
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models/enums"
	repo "github.com/MadManJJ/cms-api/repositories"

	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAppRepo_GetFaqPageBySlug(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()	

	appFaqPageRepo := repo.NewAppFaqPageRepository(gormDB)

	slug := "about/us"
	language := string(enums.PageLanguageEN)

	t.Run("successfully get faq page by slug: no preload (url alias)", func(t *testing.T) {
		emptyPreload := []string{}
		isAlias := true

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "faq_pages"."id","faq_pages"."created_at","faq_pages"."updated_at" FROM "faq_pages" JOIN faq_contents ON faq_contents.page_id = faq_pages.id WHERE faq_contents.url_alias = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "faq_content_categories"}).
				AddRow(contentId, categoryId))			
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(componentId, contentId))				
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, contentId))						

		faqPage, err := appFaqPageRepo.GetFaqPageBySlug(slug, emptyPreload, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, faqPage.ID, pageId)
		assert.Equal(t, faqPage.Contents[0].ID, contentId)
		assert.NotNil(t, faqPage.Contents[0])
		assert.NotNil(t, faqPage.Contents[0].MetaTag)
		assert.NotNil(t, faqPage.Contents[0].Revision)
		assert.NotNil(t, faqPage.Contents[0].Components)
		assert.NotNil(t, faqPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("successfully get faq page by slug: no preload (url)", func(t *testing.T) {
		emptyPreload := []string{}
		isAlias := false

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "faq_pages"."id","faq_pages"."created_at","faq_pages"."updated_at" FROM "faq_pages" JOIN faq_contents ON faq_contents.page_id = faq_pages.id WHERE faq_contents.url = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "faq_content_categories"}).
				AddRow(contentId, categoryId))			
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(componentId, contentId))				
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, contentId))						

		faqPage, err := appFaqPageRepo.GetFaqPageBySlug(slug, emptyPreload, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, faqPage.ID, pageId)
		assert.Equal(t, faqPage.Contents[0].ID, contentId)
		assert.NotNil(t, faqPage.Contents[0])
		assert.NotNil(t, faqPage.Contents[0].MetaTag)
		assert.NotNil(t, faqPage.Contents[0].Revision)
		assert.NotNil(t, faqPage.Contents[0].Components)
		assert.NotNil(t, faqPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		

	t.Run("successfully get faq page by slug: preloads all (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}
		isAlias := true
		
		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "faq_pages"."id","faq_pages"."created_at","faq_pages"."updated_at" FROM "faq_pages" JOIN faq_contents ON faq_contents.page_id = faq_pages.id WHERE faq_contents.url_alias = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "faq_content_categories"}).
				AddRow(contentId, categoryId))			
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(componentId, contentId))				
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, contentId))						

		faqPage, err := appFaqPageRepo.GetFaqPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, faqPage.ID, pageId)
		assert.Equal(t, faqPage.Contents[0].ID, contentId)
		assert.NotNil(t, faqPage.Contents[0])
		assert.NotNil(t, faqPage.Contents[0].MetaTag)
		assert.NotNil(t, faqPage.Contents[0].Revision)
		assert.NotNil(t, faqPage.Contents[0].Components)
		assert.NotNil(t, faqPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get faq page by slug: preloads all (url)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}
		isAlias := false
		
		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "faq_pages"."id","faq_pages"."created_at","faq_pages"."updated_at" FROM "faq_pages" JOIN faq_contents ON faq_contents.page_id = faq_pages.id WHERE faq_contents.url = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "faq_content_categories"}).
				AddRow(contentId, categoryId))			
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(componentId, contentId))				
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, contentId))						

		faqPage, err := appFaqPageRepo.GetFaqPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, faqPage.ID, pageId)
		assert.Equal(t, faqPage.Contents[0].ID, contentId)
		assert.NotNil(t, faqPage.Contents[0])
		assert.NotNil(t, faqPage.Contents[0].MetaTag)
		assert.NotNil(t, faqPage.Contents[0].Revision)
		assert.NotNil(t, faqPage.Contents[0].Components)
		assert.NotNil(t, faqPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("successfully get faq page by slug: partial preloads (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories"}
		isAlias := true		

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "faq_pages"."id","faq_pages"."created_at","faq_pages"."updated_at" FROM "faq_pages" JOIN faq_contents ON faq_contents.page_id = faq_pages.id WHERE faq_contents.url_alias = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "faq_content_categories"}).
				AddRow(contentId, categoryId))											

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, contentId))						

		faqPage, err := appFaqPageRepo.GetFaqPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, faqPage.ID, pageId)
		assert.Equal(t, faqPage.Contents[0].ID, contentId)
		assert.NotNil(t, faqPage.Contents[0])
		assert.Nil(t, faqPage.Contents[0].MetaTag)
		assert.NotNil(t, faqPage.Contents[0].Revision)
		assert.Nil(t, faqPage.Contents[0].Components)
		assert.NotNil(t, faqPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("successfully get faq page by slug: partial preloads (url)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories"}
		isAlias := false		

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "faq_pages"."id","faq_pages"."created_at","faq_pages"."updated_at" FROM "faq_pages" JOIN faq_contents ON faq_contents.page_id = faq_pages.id WHERE faq_contents.url = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"faq_content_id", "faq_content_categories"}).
				AddRow(contentId, categoryId))											

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
				AddRow(revisionId, contentId))						

		faqPage, err := appFaqPageRepo.GetFaqPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, faqPage.ID, pageId)
		assert.Equal(t, faqPage.Contents[0].ID, contentId)
		assert.NotNil(t, faqPage.Contents[0])
		assert.Nil(t, faqPage.Contents[0].MetaTag)
		assert.NotNil(t, faqPage.Contents[0].Revision)
		assert.Nil(t, faqPage.Contents[0].Components)
		assert.NotNil(t, faqPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		

	t.Run("failed to get faq page by slug (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}
		isAlias := true

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "faq_pages"."id","faq_pages"."created_at","faq_pages"."updated_at" FROM "faq_pages" JOIN faq_contents ON faq_contents.page_id = faq_pages.id WHERE faq_contents.url_alias = $1 ORDER BY "faq_pages"."id" LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)

		faqPage, err := appFaqPageRepo.GetFaqPageBySlug(slug, preloads, isAlias, language)
		assert.Error(t, err)
		assert.Nil(t, faqPage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestAppRepo_GetFaqContentPreview(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()	

	appFaqPageRepo := repo.NewAppFaqPageRepository(gormDB)
	
	pageId := uuid.New()
	contentId := uuid.New()
	metaTagId := uuid.New()
	componentId := uuid.New()

	t.Run("successfully get faq content preview", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))		

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "faq_content_id"}).
					AddRow(componentId, contentId))						

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))		
		
		actualFaqContent, err := appFaqPageRepo.GetFaqContentPreview(pageId)
		assert.NoError(t, err)
		assert.NotNil(t, actualFaqContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get faq content preview", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "faq_contents"`)).
			WillReturnError(errs.ErrInternalServerError)

		actualFaqContent, err := appFaqPageRepo.GetFaqContentPreview(contentId)
		assert.Error(t, err)
		assert.Nil(t, actualFaqContent)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}