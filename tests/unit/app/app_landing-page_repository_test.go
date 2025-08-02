package tests

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models/enums"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAppRepo_GetLandingPageBySlug(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	appLandingPageRepo := repo.NewAppLandingPageRepository(gormDB)

	slug := "about/us"
	language := string(enums.PageLanguageEN)

	t.Run("successfully get landing page by slug: no preload (url alias)", func(t *testing.T) {
		emptyPreload := []string{}

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()
		contentFileId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "landing_pages"."id","landing_pages"."created_at","landing_pages"."updated_at" FROM "landing_pages" JOIN landing_contents ON landing_contents.page_id = landing_pages.id WHERE landing_contents.url_alias = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "landing_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(componentId, contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_files"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(contentFileId, contentId))				

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(revisionId, contentId))

		landingPage, err := appLandingPageRepo.GetLandingPageByUrlAlias(slug, emptyPreload, language)
		assert.NoError(t, err)
		assert.Equal(t, landingPage.ID, pageId)
		assert.Equal(t, landingPage.Contents[0].ID, contentId)
		assert.NotNil(t, landingPage.Contents[0])
		assert.NotNil(t, landingPage.Contents[0].MetaTag)
		assert.NotNil(t, landingPage.Contents[0].Revision)
		assert.NotNil(t, landingPage.Contents[0].Components)
		assert.NotNil(t, landingPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get landing page by slug: preloads all (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}
		
		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "landing_pages"."id","landing_pages"."created_at","landing_pages"."updated_at" FROM "landing_pages" JOIN landing_contents ON landing_contents.page_id = landing_pages.id WHERE landing_contents.url_alias = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "landing_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(componentId, contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(revisionId, contentId))

		landingPage, err := appLandingPageRepo.GetLandingPageByUrlAlias(slug, preloads, language)
		assert.NoError(t, err)
		assert.Equal(t, landingPage.ID, pageId)
		assert.Equal(t, landingPage.Contents[0].ID, contentId)
		assert.NotNil(t, landingPage.Contents[0])
		assert.NotNil(t, landingPage.Contents[0].MetaTag)
		assert.NotNil(t, landingPage.Contents[0].Revision)
		assert.NotNil(t, landingPage.Contents[0].Components)
		assert.NotNil(t, landingPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get landing page by slug: partial preloads (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories"}

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "landing_pages"."id","landing_pages"."created_at","landing_pages"."updated_at" FROM "landing_pages" JOIN landing_contents ON landing_contents.page_id = landing_pages.id WHERE landing_contents.url_alias = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "landing_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"landing_content_id", "landing_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "landing_content_id"}).
				AddRow(revisionId, contentId))

		landingPage, err := appLandingPageRepo.GetLandingPageByUrlAlias(slug, preloads, language)
		assert.NoError(t, err)
		assert.Equal(t, landingPage.ID, pageId)
		assert.Equal(t, landingPage.Contents[0].ID, contentId)
		assert.NotNil(t, landingPage.Contents[0])
		assert.Nil(t, landingPage.Contents[0].MetaTag)
		assert.NotNil(t, landingPage.Contents[0].Revision)
		assert.Nil(t, landingPage.Contents[0].Components)
		assert.NotNil(t, landingPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get landing page by slug (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "landing_pages"."id","landing_pages"."created_at","landing_pages"."updated_at" FROM "landing_pages" JOIN landing_contents ON landing_contents.page_id = landing_pages.id WHERE landing_contents.url_alias = $1 ORDER BY "landing_pages"."id" LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)

		landingPage, err := appLandingPageRepo.GetLandingPageByUrlAlias(slug, preloads, language)
		assert.Error(t, err)
		assert.Nil(t, landingPage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}