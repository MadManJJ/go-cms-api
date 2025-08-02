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

func TestAppRepo_GetPartnerPageBySlug(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	appPartnerPageRepo := repo.NewAppPartnerPageRepository(gormDB)

	slug := "about/us"
	language := string(enums.PageLanguageEN)

	t.Run("successfully get partner page by slug: no preload (url alias)", func(t *testing.T) {
		emptyPreload := []string{}
		isAlias := true

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "partner_pages"."id","partner_pages"."created_at","partner_pages"."updated_at" FROM "partner_pages" JOIN partner_contents ON partner_contents.page_id = partner_pages.id WHERE partner_contents.url_alias = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "partner_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(componentId, contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(revisionId, contentId))

		partnerPage, err := appPartnerPageRepo.GetPartnerPageBySlug(slug, emptyPreload, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, partnerPage.ID, pageId)
		assert.Equal(t, partnerPage.Contents[0].ID, contentId)
		assert.NotNil(t, partnerPage.Contents[0])
		assert.NotNil(t, partnerPage.Contents[0].MetaTag)
		assert.NotNil(t, partnerPage.Contents[0].Revision)
		assert.NotNil(t, partnerPage.Contents[0].Components)
		assert.NotNil(t, partnerPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get partner page by slug: no preload (url)", func(t *testing.T) {
		emptyPreload := []string{}
		isAlias := false

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "partner_pages"."id","partner_pages"."created_at","partner_pages"."updated_at" FROM "partner_pages" JOIN partner_contents ON partner_contents.page_id = partner_pages.id WHERE partner_contents.url = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "partner_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(componentId, contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(revisionId, contentId))

		partnerPage, err := appPartnerPageRepo.GetPartnerPageBySlug(slug, emptyPreload, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, partnerPage.ID, pageId)
		assert.Equal(t, partnerPage.Contents[0].ID, contentId)
		assert.NotNil(t, partnerPage.Contents[0])
		assert.NotNil(t, partnerPage.Contents[0].MetaTag)
		assert.NotNil(t, partnerPage.Contents[0].Revision)
		assert.NotNil(t, partnerPage.Contents[0].Components)
		assert.NotNil(t, partnerPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get partner page by slug: preloads all (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}
		isAlias := true

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "partner_pages"."id","partner_pages"."created_at","partner_pages"."updated_at" FROM "partner_pages" JOIN partner_contents ON partner_contents.page_id = partner_pages.id WHERE partner_contents.url_alias = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "partner_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(componentId, contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(revisionId, contentId))

		partnerPage, err := appPartnerPageRepo.GetPartnerPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, partnerPage.ID, pageId)
		assert.Equal(t, partnerPage.Contents[0].ID, contentId)
		assert.NotNil(t, partnerPage.Contents[0])
		assert.NotNil(t, partnerPage.Contents[0].MetaTag)
		assert.NotNil(t, partnerPage.Contents[0].Revision)
		assert.NotNil(t, partnerPage.Contents[0].Components)
		assert.NotNil(t, partnerPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get partner page by slug: preloads all (url)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}
		isAlias := false

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		componentId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "partner_pages"."id","partner_pages"."created_at","partner_pages"."updated_at" FROM "partner_pages" JOIN partner_contents ON partner_contents.page_id = partner_pages.id WHERE partner_contents.url = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "partner_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "components"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(componentId, contentId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meta_tags"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(revisionId, contentId))

		partnerPage, err := appPartnerPageRepo.GetPartnerPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, partnerPage.ID, pageId)
		assert.Equal(t, partnerPage.Contents[0].ID, contentId)
		assert.NotNil(t, partnerPage.Contents[0])
		assert.NotNil(t, partnerPage.Contents[0].MetaTag)
		assert.NotNil(t, partnerPage.Contents[0].Revision)
		assert.NotNil(t, partnerPage.Contents[0].Components)
		assert.NotNil(t, partnerPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get partner page by slug: partial preloads (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories"}
		isAlias := true

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "partner_pages"."id","partner_pages"."created_at","partner_pages"."updated_at" FROM "partner_pages" JOIN partner_contents ON partner_contents.page_id = partner_pages.id WHERE partner_contents.url_alias = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "partner_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(revisionId, contentId))

		partnerPage, err := appPartnerPageRepo.GetPartnerPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, partnerPage.ID, pageId)
		assert.Equal(t, partnerPage.Contents[0].ID, contentId)
		assert.NotNil(t, partnerPage.Contents[0])
		assert.Nil(t, partnerPage.Contents[0].MetaTag)
		assert.NotNil(t, partnerPage.Contents[0].Revision)
		assert.Nil(t, partnerPage.Contents[0].Components)
		assert.NotNil(t, partnerPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully get partner page by slug: partial preloads (url)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories"}
		isAlias := false

		pageId := uuid.New()
		contentId := uuid.New()
		metaTagId := uuid.New()
		categoryId := uuid.New()
		revisionId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "partner_pages"."id","partner_pages"."created_at","partner_pages"."updated_at" FROM "partner_pages" JOIN partner_contents ON partner_contents.page_id = partner_pages.id WHERE partner_contents.url = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(pageId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_contents"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "page_id", "meta_tag_id"}).
				AddRow(contentId, pageId, metaTagId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "partner_content_categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"partner_content_id", "partner_content_categories"}).
				AddRow(contentId, categoryId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "revisions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "partner_content_id"}).
				AddRow(revisionId, contentId))

		partnerPage, err := appPartnerPageRepo.GetPartnerPageBySlug(slug, preloads, isAlias, language)
		assert.NoError(t, err)
		assert.Equal(t, partnerPage.ID, pageId)
		assert.Equal(t, partnerPage.Contents[0].ID, contentId)
		assert.NotNil(t, partnerPage.Contents[0])
		assert.Nil(t, partnerPage.Contents[0].MetaTag)
		assert.NotNil(t, partnerPage.Contents[0].Revision)
		assert.Nil(t, partnerPage.Contents[0].Components)
		assert.NotNil(t, partnerPage.Contents[0].Categories)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get partner page by slug (url alias)", func(t *testing.T) {
		preloads := []string{"Contents.Revision", "Contents.Categories", "Contents.Components", "Contents.MetaTag"}
		isAlias := true

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "partner_pages"."id","partner_pages"."created_at","partner_pages"."updated_at" FROM "partner_pages" JOIN partner_contents ON partner_contents.page_id = partner_pages.id WHERE partner_contents.url_alias = $1 ORDER BY "partner_pages"."id" LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)

		partnerPage, err := appPartnerPageRepo.GetPartnerPageBySlug(slug, preloads, isAlias, language)
		assert.Error(t, err)
		assert.Nil(t, partnerPage)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}