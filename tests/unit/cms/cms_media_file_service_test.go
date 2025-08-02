package tests

import (
	"testing"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockMediaFileRepository struct {
	create   func(file *models.MediaFile) (*models.MediaFile, error)
	findByID func(id uuid.UUID) (*models.MediaFile, error)
	findByNameAndPath func(name string, path string) (*models.MediaFile, error)
	list     func(filter dto.MediaFileListFilter) ([]models.MediaFile, int64, error)
	delete   func(id uuid.UUID) error
}

func (m *MockMediaFileRepository) Create(file *models.MediaFile) (*models.MediaFile, error) {
	return m.create(file)
}

func (m *MockMediaFileRepository) FindByID(id uuid.UUID) (*models.MediaFile, error) {
	return m.findByID(id)
}

func (m *MockMediaFileRepository) FindByNameAndPath(name string, path string) (*models.MediaFile, error) {
	return m.findByNameAndPath(name, path)
}

func (m *MockMediaFileRepository) List(filter dto.MediaFileListFilter) ([]models.MediaFile, int64, error) {
	return m.list(filter)
}

func (m *MockMediaFileRepository) Delete(id uuid.UUID) error {
	return m.delete(id)
}

func TestCMSService_UploadMediaFile(t *testing.T) {
	cfg := config.New()

	originalFilename := "test.png"
	mimeType := "image/png"
	fileData := []byte("test")
	customPathOpt := "some-path"
	replaceOpt := false
	userID := uuid.New()

	t.Run("successfully upload media file", func(t *testing.T) {
		mockMediaFile := helpers.InitializeMockMediaFile()
		mockMediaFile.ID = uuid.New()

		mediaFileResponse := &dto.MediaFileResponse{
			ID:          mockMediaFile.ID.String(),
			Name:        mockMediaFile.Name,
			DownloadURL: mockMediaFile.DownloadURL,
			CreatedAt:   mockMediaFile.CreatedAt,
			UpdatedAt:   mockMediaFile.UpdatedAt,
		}

		repo := &MockMediaFileRepository{
			create: func(file *models.MediaFile) (*models.MediaFile, error) {
				return mockMediaFile, nil
			},
			findByNameAndPath: func(name string, path string) (*models.MediaFile, error) {
				return nil, nil
			},
			delete: func(id uuid.UUID) error {
				return nil
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		actualMediaFile, err := service.UploadMediaFile(originalFilename, mimeType, fileData, &customPathOpt, &replaceOpt, userID)
		assert.NoError(t, err)
		assert.Equal(t, mediaFileResponse, actualMediaFile)
	})

	t.Run("failed to upload media file", func(t *testing.T) {
		repo := &MockMediaFileRepository{
			create: func(file *models.MediaFile) (*models.MediaFile, error) {
				return nil, errs.ErrInternalServerError
			},
			findByNameAndPath: func(name string, path string) (*models.MediaFile, error) {
				return nil, nil
			},
			delete: func(id uuid.UUID) error {
				return nil
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		actualMediaFile, err := service.UploadMediaFile(originalFilename, mimeType, fileData, &customPathOpt, &replaceOpt, userID)
		assert.Error(t, err)
		assert.Nil(t, actualMediaFile)
	})
}

func TestCMSService_GetMediaFileByID(t *testing.T) {
	cfg := config.New()

	mediaFileId := uuid.New()
	mockMediaFile := helpers.InitializeMockMediaFile()
	mockMediaFile.ID = mediaFileId

	t.Run("successfully get media file by ID", func(t *testing.T) {
		repo := &MockMediaFileRepository{
			findByID: func(id uuid.UUID) (*models.MediaFile, error) {
				return mockMediaFile, nil
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		actualMediaFile, err := service.GetMediaFileByID(mediaFileId.String())
		assert.NoError(t, err)
		assert.NotNil(t, actualMediaFile)
		assert.Equal(t, mockMediaFile.ID.String(), actualMediaFile.ID)
		assert.Equal(t, mockMediaFile.Name, actualMediaFile.Name)
	})

	t.Run("failed to get media file by ID", func(t *testing.T) {
		repo := &MockMediaFileRepository{
			findByID: func(id uuid.UUID) (*models.MediaFile, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		actualMediaFile, err := service.GetMediaFileByID(mediaFileId.String())
		assert.Error(t, err)
		assert.Nil(t, actualMediaFile)
	})
}

func TestCMSService_ListMediaFiles(t *testing.T) {
	cfg := config.New()

	mediaFileId := uuid.New()
	mockMediaFile := helpers.InitializeMockMediaFile()
	mockMediaFile.ID = mediaFileId
	
	t.Run("successfully list media files", func(t *testing.T) {
		total := int64(1)

		repo := &MockMediaFileRepository{
			list: func(filter dto.MediaFileListFilter) ([]models.MediaFile, int64, error) {
				return []models.MediaFile{*mockMediaFile}, total, nil
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		actualMediaFiles, err := service.ListMediaFiles(dto.MediaFileListFilter{})
		assert.NoError(t, err)
		assert.NotNil(t, actualMediaFiles)
		assert.Equal(t, mockMediaFile.ID.String(), actualMediaFiles.Data[0].ID)
		assert.Equal(t, mockMediaFile.Name, actualMediaFiles.Data[0].Name)
		assert.Equal(t, total, actualMediaFiles.Total)
	})

	t.Run("failed to list media files", func(t *testing.T) {
		repo := &MockMediaFileRepository{
			list: func(filter dto.MediaFileListFilter) ([]models.MediaFile, int64, error) {
				return nil, 0, errs.ErrInternalServerError
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		actualMediaFiles, err := service.ListMediaFiles(dto.MediaFileListFilter{})
		assert.Error(t, err)
		assert.Nil(t, actualMediaFiles)
	})
}

func TestCMSService_DeleteMediaFile(t *testing.T) {
	cfg := config.New()

	mediaFileId := uuid.New()
	mockMediaFile := helpers.InitializeMockMediaFile()
	mockMediaFile.ID = mediaFileId
	userId := uuid.New()	

	t.Run("successfully delete media file", func(t *testing.T) {
		repo := &MockMediaFileRepository{
			findByID: func(id uuid.UUID) (*models.MediaFile, error) {
				return mockMediaFile, nil
			},
			delete: func(id uuid.UUID) error {
				return nil
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		err := service.DeleteMediaFile(mediaFileId.String(), userId)
		assert.NoError(t, err)
	})

	t.Run("failed to delete media file", func(t *testing.T) {
		repo := &MockMediaFileRepository{
			findByID: func(id uuid.UUID) (*models.MediaFile, error) {
				return mockMediaFile, nil
			},
			delete: func(id uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}

		service := services.NewMediaFileService(cfg, repo)

		err := service.DeleteMediaFile(mediaFileId.String(), userId)
		assert.Error(t, err)
	})
}