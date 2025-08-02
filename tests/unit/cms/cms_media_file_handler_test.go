package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	"github.com/MadManJJ/cms-api/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockMediaFileService struct {
	mock.Mock
}

func (m *MockMediaFileService) UploadMediaFile(originalFilename string, mimeType string, fileData []byte, customPath *string, replace *bool, userID uuid.UUID) (*dto.MediaFileResponse, error) {
	args := m.Called(originalFilename, mimeType, fileData, customPath, replace, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaFileResponse), args.Error(1)
}

func (m *MockMediaFileService) GetMediaFileByID(idStr string) (*dto.MediaFileResponse, error) {
	args := m.Called(idStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaFileResponse), args.Error(1)
}

func (m *MockMediaFileService) ListMediaFiles(filter dto.MediaFileListFilter) (*dto.MediaFilesListResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaFilesListResponse), args.Error(1)
}

func (m *MockMediaFileService) DeleteMediaFile(idStr string, userID uuid.UUID) error {
	args := m.Called(idStr, userID)
	return args.Error(0)
}

func TestCMSMediaFileHandler(t *testing.T) {
	mockService := &MockMediaFileService{}
	handler := cmsHandler.NewMediaFileHandler(mockService)

	userId := uuid.New()

	app := fiber.New()
	app.Post("/cms/mediafiles", func (c *fiber.Ctx) error {
		c.Locals("userID", userId)
		return handler.HandleUploadMediaFile(c)
	})
	app.Get("/cms/mediafiles", handler.HandleListMediaFiles)
	app.Get("/cms/mediafiles/:id", handler.HandleGetMediaFileByID)
	app.Delete("/cms/mediafiles/:id", handler.HandleDeleteMediaFile)

	t.Run("POST /cms/mediafiles HandleUploadMediaFile", func(t *testing.T) {

		mediaResponse := dto.MediaFileResponse{
			ID:           uuid.New().String(),
			Name:         "test.png",
			DownloadURL:  "http://localhost:8080/mediafiles/" + uuid.New().String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}		

		t.Run("successfully update media file", func(t *testing.T) {
			body, contentType := helpers.CreateMultipartRequest(t)

			mockService.On("UploadMediaFile", "test.png", "image/png", []byte("dummy image data"), mock.AnythingOfType("*string"), mock.AnythingOfType("*bool"), userId).Return(&mediaResponse, nil)
	
			// --- Step 2: Send the request ---
			req := httptest.NewRequest("POST", "/cms/mediafiles", body)
			req.Header.Set("Content-Type", contentType)
	
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	
			var response dto.MediaFileResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, mediaResponse.ID, response.ID)
			assert.Equal(t, mediaResponse.Name, response.Name)
			assert.Equal(t, mediaResponse.DownloadURL, response.DownloadURL)
			assert.WithinDuration(t, mediaResponse.CreatedAt, response.CreatedAt, time.Second)
			assert.WithinDuration(t, mediaResponse.UpdatedAt, response.UpdatedAt, time.Second)
		})

		t.Run("failed to update media file", func(t *testing.T) {
			mockService.ExpectedCalls = nil

			body, contentType := helpers.CreateMultipartRequest(t)

			mockService.On("UploadMediaFile", "test.png", "image/png", []byte("dummy image data"), mock.AnythingOfType("*string"), mock.AnythingOfType("*bool"), userId).Return(nil, errs.ErrInternalServerError)
		
			// --- Step 2: Send the request ---
			req := httptest.NewRequest("POST", "/cms/mediafiles", body)
			req.Header.Set("Content-Type", contentType)
		
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})

	t.Run("GET /cms/mediafiles HandleListMediaFiles", func(t *testing.T) {
		search := "test"
		sort := "name"
		page := 1
		pageSize := 10
		order := "asc"
		filter := dto.MediaFileListFilter{
			Search: &search,
			Page: page,
			PageSize: pageSize,
			SortBy: &sort,
			Order: &order,
		}

		mediaFile := helpers.InitializeMockMediaFile()
		mediaFileItemResponse := dto.MediaFileListItemResponse{
			ID: mediaFile.ID.String(),
			Name: mediaFile.Name,
			DownloadURL: "http://localhost:8080/mediafiles/" + mediaFile.ID.String(),
			CreatedAt: mediaFile.CreatedAt,
		}
		mediaFileListResponse := dto.MediaFilesListResponse{
			Data: []dto.MediaFileListItemResponse{mediaFileItemResponse},
			Total: 1,
			Page: page,
			PageSize: pageSize,
		}

		body, err := json.Marshal(filter)
		require.NoError(t, err)				

		t.Run("successfully list media files", func(t *testing.T) {
			mockService.On("ListMediaFiles", mock.AnythingOfType("dto.MediaFileListFilter")).Return(&mediaFileListResponse, nil)
			
			req := httptest.NewRequest("GET", "/cms/mediafiles", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")	

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to list media files", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("ListMediaFiles", mock.AnythingOfType("dto.MediaFileListFilter")).Return(nil, errs.ErrInternalServerError)
			
			req := httptest.NewRequest("GET", "/cms/mediafiles", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")	
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})

	t.Run("GET /cms/mediafiles/:id HandleGetMediaFileByID", func(t *testing.T) {
		idStr := "test"
		mediaFile := helpers.InitializeMockMediaFile()
		mediaFileResponse := dto.MediaFileResponse{
			ID: mediaFile.ID.String(),
			Name: mediaFile.Name,
			DownloadURL: "http://localhost:8080/mediafiles/" + mediaFile.ID.String(),
			CreatedAt: mediaFile.CreatedAt,
			UpdatedAt: mediaFile.UpdatedAt,
		}

		t.Run("successfully get media file by ID", func(t *testing.T) {
			mockService.On("GetMediaFileByID", idStr).Return(&mediaFileResponse, nil)
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/mediafiles/%s", idStr), nil)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

		t.Run("failed to get media file by ID", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetMediaFileByID", idStr).Return(nil, errs.ErrInternalServerError)
			
			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/mediafiles/%s", idStr), nil)
			req.Header.Set("Content-Type", "application/json")	
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})

	t.Run("DELETE /cms/mediafiles/:id HandleDeleteMediaFile", func(t *testing.T) {
		idStr := "test"
		
		t.Run("successfully delete media file by ID", func(t *testing.T) {
			mockService.On("DeleteMediaFile", idStr, mock.AnythingOfType("uuid.UUID")).Return(nil)
			
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/mediafiles/%s", idStr), nil)
			req.Header.Set("Content-Type", "application/json")	
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		})

		t.Run("failed to delete media file by ID", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("DeleteMediaFile", idStr, mock.AnythingOfType("uuid.UUID")).Return(errs.ErrInternalServerError)
			
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/mediafiles/%s", idStr), nil)
			req.Header.Set("Content-Type", "application/json")	
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	})
}
