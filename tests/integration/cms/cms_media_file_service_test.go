package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/repositories"
	"github.com/MadManJJ/cms-api/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupMediaFileTest(t *testing.T) (
	sqlmock.Sqlmock,
	services.MediaFileServiceInterface,
	*config.Config,
	string,
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	repo := repositories.NewMediaFileRepository(gormDB)

	tempDir := t.TempDir()

	testCfg := &config.Config{
		App: config.AppConfig{
			UploadPath:       tempDir,
			APIBaseURL:       "http://localhost:8080",
			StaticFilePrefix: "/files",
		},
	}

	service := services.NewMediaFileService(testCfg, repo)

	return mock, service, testCfg, tempDir, cleanup
}

func TestMediaFileLifecycle(t *testing.T) {
	mock, service, cfg, tempDir, cleanup := setupMediaFileTest(t)
	defer cleanup()

	var createdFileID uuid.UUID
	var createdFilename string
	var createdDownloadURL string
	testUserID := uuid.New()

	// --- 1. Upload a new file (Success) ---
	t.Run("1_UploadNewFile_Success", func(t *testing.T) {
		filename := "test_image.jpg"
		fileData := []byte("fake-jpeg-data")
		mimeType := "image/jpeg"
		var customPath *string = nil
		replace := false

		mock.ExpectQuery(`SELECT \* FROM "media_files" WHERE name = \$1 ORDER BY "media_files"\."id" LIMIT \$2`).
			WithArgs(filename, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "media_files"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		resp, err := service.UploadMediaFile(filename, mimeType, fileData, customPath, &replace, testUserID)

		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, filename, resp.Name)
		expectedURL := fmt.Sprintf("%s%s/%s", cfg.App.APIBaseURL, cfg.App.StaticFilePrefix, filename)
		assert.Equal(t, expectedURL, resp.DownloadURL)

		filePath := filepath.Join(tempDir, filename)
		_, fileErr := os.Stat(filePath)
		assert.NoError(t, fileErr, "File should exist on disk")

		createdFileID, _ = uuid.Parse(resp.ID)
		createdFilename = resp.Name
		createdDownloadURL = resp.DownloadURL
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 2. Upload a duplicate file (no-replace) ---
	t.Run("2_UploadDuplicateFile_NoReplace_GeneratesNewName", func(t *testing.T) {
		filename := "test_image.jpg"
		fileData := []byte("new-jpeg-data")
		mimeType := "image/jpeg"
		var customPath *string = nil
		replace := false

		mock.ExpectQuery(`SELECT \* FROM "media_files" WHERE name = \$1 ORDER BY "media_files"\."id" LIMIT \$2`).
			WithArgs(filename, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(createdFileID, createdFilename))

		expectedNewFilename := "test_image___1.jpg"
		mock.ExpectQuery(`SELECT \* FROM "media_files" WHERE name = \$1 ORDER BY "media_files"\."id" LIMIT \$2`).
			WithArgs(expectedNewFilename, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "media_files"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		resp, err := service.UploadMediaFile(filename, mimeType, fileData, customPath, &replace, testUserID)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, expectedNewFilename, resp.Name)

		_, errOld := os.Stat(filepath.Join(tempDir, filename))
		_, errNew := os.Stat(filepath.Join(tempDir, expectedNewFilename))
		assert.NoError(t, errOld, "Original file should still exist")
		assert.NoError(t, errNew, "New file with unique name should exist")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 3. Upload a duplicate file (with-replace) ---
	t.Run("3_UploadDuplicateFile_WithReplace_Success", func(t *testing.T) {
		filename := "test_image.jpg"
		fileData := []byte("replaced-jpeg-data")
		mimeType := "image/jpeg"
		var customPath *string = nil
		replace := true

		mock.ExpectQuery(`SELECT \* FROM "media_files" WHERE name = \$1 ORDER BY "media_files"\."id" LIMIT \$2`).
			WithArgs(filename, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "download_url"}).AddRow(createdFileID, createdFilename, createdDownloadURL))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "media_files" WHERE id = $1`)).
			WithArgs(createdFileID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "media_files"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(createdFileID))
		mock.ExpectCommit()

		resp, err := service.UploadMediaFile(filename, mimeType, fileData, customPath, &replace, testUserID)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, filename, resp.Name)

		filePath := filepath.Join(tempDir, filename)
		readData, _ := os.ReadFile(filePath)
		assert.Equal(t, fileData, readData, "File content should be replaced")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 4. List and Get File ---
	t.Run("4_ListAndGetFile_Success", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdFileID, "File must be created first")

		filter := dto.MediaFileListFilter{Search: &createdFilename}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "media_files" WHERE name ILIKE $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "media_files" WHERE name ILIKE $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(createdFileID, createdFilename))

		listResp, err := service.ListMediaFiles(filter)
		require.NoError(t, err)
		require.Len(t, listResp.Data, 1)

		mock.ExpectQuery(`SELECT \* FROM "media_files" WHERE id = \$1 ORDER BY "media_files"\."id" LIMIT \$2`).
			WithArgs(createdFileID.String(), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(createdFileID, createdFilename))

		getResp, err := service.GetMediaFileByID(createdFileID.String())
		require.NoError(t, err)
		assert.Equal(t, createdFilename, getResp.Name)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	// --- 5. Delete File ---
	t.Run("5_DeleteFile_Success", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdFileID, "File must be created first")

		mock.ExpectQuery(`SELECT \* FROM "media_files" WHERE id = \$1 ORDER BY "media_files"\."id" LIMIT \$2`).
			WithArgs(createdFileID.String(), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "download_url"}).AddRow(createdFileID, createdFilename, createdDownloadURL))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "media_files" WHERE id = $1`)).
			WithArgs(createdFileID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.DeleteMediaFile(createdFileID.String(), testUserID)

		require.NoError(t, err)

		filePath := filepath.Join(tempDir, createdFilename)
		_, fileErr := os.Stat(filePath)
		assert.True(t, os.IsNotExist(fileErr), "File should be deleted from disk")

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
