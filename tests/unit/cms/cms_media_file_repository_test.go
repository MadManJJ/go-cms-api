package tests

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCMSRepo_CreateMediaFile(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()	

	cmsMediaFileRepo := repo.NewMediaFileRepository(gormDB)

	t.Run("successfully create media file", func(t *testing.T) {
		mockMediaFile := helpers.InitializeMockMediaFile()
		mock.ExpectBegin()
		
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "media_files"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()))

		mock.ExpectCommit()

		createdMediaFile, err := cmsMediaFileRepo.Create(mockMediaFile)
		assert.NoError(t, err)
		assert.Equal(t, mockMediaFile, createdMediaFile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to create media file", func(t *testing.T) {
		mockMediaFile := helpers.InitializeMockMediaFile()
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "media_files"`)).
			WillReturnError(errs.ErrInternalServerError)

		mock.ExpectRollback()

		createdMediaFile, err := cmsMediaFileRepo.Create(mockMediaFile)
		assert.Error(t, err)
		assert.Nil(t, createdMediaFile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSRepo_FindById(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()	

	cmsMediaFileRepo := repo.NewMediaFileRepository(gormDB)

	t.Run("successfully find media file by id", func(t *testing.T) {
		mediaFileId := uuid.New()
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "media_files" WHERE id = $1 ORDER BY "media_files"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(mediaFileId))

		mediaFile, err := cmsMediaFileRepo.FindByID(mediaFileId)
		assert.NoError(t, err)
		assert.Equal(t, mediaFileId, mediaFile.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find media file by id", func(t *testing.T) {
		mediaFileId := uuid.New()
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "media_files" WHERE id = $1 ORDER BY "media_files"."id" LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)

		mediaFile, err := cmsMediaFileRepo.FindByID(mediaFileId)
		assert.Error(t, err)
		assert.Nil(t, mediaFile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// ! Check this function name again later
func TestCMSRepo_FindByNameAndPath(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()	

	cmsMediaFileRepo := repo.NewMediaFileRepository(gormDB)
	
	name := "fake name"
	path := "fake path"

	t.Run("successfully find media file by name and path", func(t *testing.T) {
		mediaFileId := uuid.New()
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "media_files" WHERE name = $1 ORDER BY "media_files"."id" LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(mediaFileId))

		mediaFile, err := cmsMediaFileRepo.FindByNameAndPath(name, path)
		assert.NoError(t, err)
		assert.Equal(t, mediaFileId, mediaFile.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to find media file by name and path", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "media_files" WHERE name = $1 ORDER BY "media_files"."id" LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)

		mediaFile, err := cmsMediaFileRepo.FindByNameAndPath(name, path)
		assert.Error(t, err)
		assert.Nil(t, mediaFile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSRepo_ListMediaFile(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()	

	cmsMediaFileRepo := repo.NewMediaFileRepository(gormDB)

	search := "fake search"
	sortBy := "fake sort by"
	order := "fake order"
	page := 1
	pageSize := 10
	mediaFileFilter := dto.MediaFileListFilter{
		Search: &search,
		SortBy: &sortBy,
		Order: &order,
		Page: page,
		PageSize: pageSize,
	}

	t.Run("successfully list media file", func(t *testing.T) {
		mediaFileId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "media_files" WHERE name ILIKE $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(2))		
		
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "media_files" WHERE name ILIKE $1 ORDER BY created_at DESC LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(mediaFileId).
				AddRow(mediaFileId))

		mediaFiles, total, err := cmsMediaFileRepo.List(mediaFileFilter)
		assert.NoError(t, err)
		assert.Equal(t, mediaFileId, mediaFiles[0].ID)
		assert.Equal(t, int64(2), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to list media file", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "media_files" WHERE name ILIKE $1`)).
			WillReturnError(errs.ErrInternalServerError)

		mediaFiles, total, err := cmsMediaFileRepo.List(mediaFileFilter)
		assert.Error(t, err)
		assert.Nil(t, mediaFiles)
		assert.Equal(t, int64(0), total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSRepo_DeleteMediaFile(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()	

	cmsMediaFileRepo := repo.NewMediaFileRepository(gormDB)

	mediaFileId := uuid.New()

	t.Run("successfully delete media file", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "media_files" WHERE id = $1`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := cmsMediaFileRepo.Delete(mediaFileId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to delete media file", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "media_files" WHERE id = $1`)).
			WillReturnError(errs.ErrInternalServerError)
		mock.ExpectRollback()

		err := cmsMediaFileRepo.Delete(mediaFileId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}