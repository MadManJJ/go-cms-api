package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCMSEmailContentService struct {
	mock.Mock
}

func (m *MockCMSEmailContentService) CreateContent(req dto.CreateEmailContentRequest) (*dto.EmailContentResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}	
	return args.Get(0).(*dto.EmailContentResponse), args.Error(1)
}

func (m *MockCMSEmailContentService) GetContentByID(idStr string) (*dto.EmailContentResponse, error) {
	args := m.Called(idStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}	
	return args.Get(0).(*dto.EmailContentResponse), args.Error(1)
}

func (m *MockCMSEmailContentService) GetContentByCategoryAndLangAndLabel(categoryIDStr string, language enums.PageLanguage, label string) (*dto.EmailContentResponse, error) {
	args := m.Called(categoryIDStr, language, label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}	
	return args.Get(0).(*dto.EmailContentResponse), args.Error(1)
}

func (m *MockCMSEmailContentService) ListContents(filters dto.EmailContentFilter) ([]dto.EmailContentResponse, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} 	
	return args.Get(0).([]dto.EmailContentResponse), args.Error(1)
}

func (m *MockCMSEmailContentService) UpdateContent(idStr string, req dto.UpdateEmailContentRequest) (*dto.EmailContentResponse, error) {
	args := m.Called(idStr, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} 	
	return args.Get(0).(*dto.EmailContentResponse), args.Error(1)
}

func (m *MockCMSEmailContentService) DeleteContent(idStr string) error {
	args := m.Called(idStr)
	return args.Error(0)
}

func (m *MockCMSEmailContentService) GetEmailContentByCategoryIDAndLanguage(categoryIDStr string, language enums.PageLanguage) ([]*dto.EmailContentResponse, error) {
	args := m.Called(categoryIDStr, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} 	
	return args.Get(0).([]*dto.EmailContentResponse), args.Error(1)
}

func TestCMSEmailContentHandler(t *testing.T) {
	mockService := &MockCMSEmailContentService{}
	handler := cmsHandler.NewEmailContentHandler(mockService)

	app := fiber.New()
	app.Post("/cms/email-contents", handler.HandleCreateEmailContent)
	app.Get("/cms/email-contents", handler.HandleListEmailContents)
	app.Get("/cms/email-contents/:id", handler.HandleGetEmailContent)
	app.Patch("/cms/email-contents/:id", handler.HandleUpdateEmailContent)
	app.Delete("/cms/email-contents/:id", handler.HandleDeleteEmailContent)
	app.Get("/cms/email-contents/category/:email_category_id/language/:language", handler.HandleGetEmailContentByCategoryAndLanguage)

	t.Run("POST /cms/email-contents HandleCreateEmailContent", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()
		mockEmailContent.EmailCategoryID = uuid.New()

		mockCreateEmailContentRequest := dto.CreateEmailContentRequest{
			EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
			Language:        mockEmailContent.Language,
			Label:           mockEmailContent.Label,
			EmailContentDetailBase: dto.EmailContentDetailBase{
				SendTo:          mockEmailContent.SendTo,
				CcEmail:         mockEmailContent.CcEmail,
				BccEmail:        mockEmailContent.BccEmail,
				SendFromEmail:   mockEmailContent.SendFromEmail,
				SendFromName:    mockEmailContent.SendFromName,
				Subject:         mockEmailContent.Subject,
				TopImgLink:      mockEmailContent.TopImgLink,
				Header:          mockEmailContent.Header,
				Paragraph:       mockEmailContent.Paragraph,
				Footer:          mockEmailContent.Footer,
				FooterImageLink: mockEmailContent.FooterImageLink,
			},
		}

		body, err := json.Marshal(mockCreateEmailContentRequest)
		require.NoError(t, err)					

		t.Run("successfully create email content", func(t *testing.T) {
			mockService.On("CreateContent", mock.AnythingOfType("dto.CreateEmailContentRequest")).Return(&dto.EmailContentResponse{
				ID:              mockEmailContent.ID.String(),
				EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
				Language:        mockEmailContent.Language,
				Label:           mockEmailContent.Label,
			}, nil)

			req := httptest.NewRequest("POST", "/cms/email-contents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to create email content", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateContent", mock.AnythingOfType("dto.CreateEmailContentRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("POST", "/cms/email-contents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})

	t.Run("GET /cms/email-contents HandleListEmailContents", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()
		emailCategoryId := uuid.New()
		emailCategoryIdStr := emailCategoryId.String()
		mockEmailContent.ID = uuid.New()
		mockEmailContent.EmailCategoryID = emailCategoryId

		emailContentFilter := dto.EmailContentFilter{
			EmailCategoryID: &emailCategoryIdStr,
			Language:        &mockEmailContent.Language,
			Label:           &mockEmailContent.Label,
		}

		body, err := json.Marshal(emailContentFilter)
		require.NoError(t, err)					

		t.Run("successfully list email contents", func(t *testing.T) {
			mockService.On("ListContents", mock.AnythingOfType("dto.EmailContentFilter")).Return([]dto.EmailContentResponse{
				{
					ID:              mockEmailContent.ID.String(),
					EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
					Language:        mockEmailContent.Language,
					Label:           mockEmailContent.Label,
				},
			}, nil)

			req := httptest.NewRequest("GET", "/cms/email-contents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to list email contents", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("ListContents", mock.AnythingOfType("dto.EmailContentFilter")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", "/cms/email-contents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})	

	t.Run("GET /cms/email-contents/:id HandleGetEmailContent", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()
		emailContentId := uuid.New()				
		mockEmailContent.ID = emailContentId

		t.Run("successfully get email content", func(t *testing.T) {
			mockService.On("GetContentByID", mock.AnythingOfType("string")).Return(&dto.EmailContentResponse{
				ID:              mockEmailContent.ID.String(),
				EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
				Language:        mockEmailContent.Language,
				Label:           mockEmailContent.Label,
			}, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/email-contents/%s", emailContentId.String()), nil)
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to get email content", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetContentByID", mock.AnythingOfType("string")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/email-contents/%s", emailContentId.String()), nil)
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})		

	t.Run("PATCH /cms/email-contents/:id HandleUpdateEmailContent", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()
		emailContentId := uuid.New()				
		mockEmailContent.ID = emailContentId

		updateEmailContent := *mockEmailContent
		updateEmailContent.Header = "Updated Header"
		updateEmailContent.Label = "Updated Label"
		updateEmailContent.SendTo = "Updated SendTo"
		updateEmailContent.CcEmail = "Updated CcEmail"
		updateEmailContent.BccEmail = "Updated BccEmail"

		updateEmailContentRequest := dto.UpdateEmailContentRequest{
			Language:        &updateEmailContent.Language,
			Label:           &updateEmailContent.Label,
			SendTo:          &updateEmailContent.SendTo,
			CcEmail:         &updateEmailContent.CcEmail,
			BccEmail:        &updateEmailContent.BccEmail,
			SendFromEmail:   &updateEmailContent.SendFromEmail,
			SendFromName:    &updateEmailContent.SendFromName,
			Subject:         &updateEmailContent.Subject,
			TopImgLink:      &updateEmailContent.TopImgLink,
			Header:          &updateEmailContent.Header,
			Paragraph:       &updateEmailContent.Paragraph,
			Footer:          &updateEmailContent.Footer,
			FooterImageLink: &updateEmailContent.FooterImageLink,
		}

		body, err := json.Marshal(updateEmailContentRequest)
		require.NoError(t, err)				
		
		t.Run("successfully update email content", func(t *testing.T) {
			mockService.On("UpdateContent", mock.AnythingOfType("string"), mock.AnythingOfType("dto.UpdateEmailContentRequest")).Return(&dto.EmailContentResponse{
				ID:              updateEmailContent.ID.String(),
				EmailCategoryID: updateEmailContent.EmailCategoryID.String(),
				Language:        updateEmailContent.Language,
				Label:           updateEmailContent.Label,
				SendTo:          updateEmailContent.SendTo,
				CcEmail:         updateEmailContent.CcEmail,
				BccEmail:        updateEmailContent.BccEmail,
				SendFromEmail:   updateEmailContent.SendFromEmail,
				SendFromName:    updateEmailContent.SendFromName,
				Subject:         updateEmailContent.Subject,
				TopImgLink:      updateEmailContent.TopImgLink,
				Header:          updateEmailContent.Header,
				Paragraph:       updateEmailContent.Paragraph,
				Footer:          updateEmailContent.Footer,
				FooterImageLink: updateEmailContent.FooterImageLink,
			}, nil)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/cms/email-contents/%s", emailContentId.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to update email content", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("UpdateContent", mock.AnythingOfType("string"), mock.AnythingOfType("dto.UpdateEmailContentRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/cms/email-contents/%s", emailContentId.String()), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")		
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})		

	t.Run("DELETE /cms/email-contents/:id HandleDeleteEmailContent", func(t *testing.T) {
		emailContentId := uuid.New()

		t.Run("successfully delete email content", func(t *testing.T) {
			mockService.On("DeleteContent", mock.AnythingOfType("string")).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/email-contents/%s", emailContentId.String()), nil)
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to delete email content", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("DeleteContent", mock.AnythingOfType("string")).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/cms/email-contents/%s", emailContentId.String()), nil)
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
	})

	t.Run("GET /cms/email-contents/category/:email_category_id/language/:language HandleGetEmailContentByCategoryAndLanguage", func(t *testing.T) {
		mockEmailContent := helpers.InitializeMockEmailContent()
		emailCategoryId := uuid.New()
		emailCategoryIdStr := emailCategoryId.String()
		mockEmailContent.ID = uuid.New()
		mockEmailContent.EmailCategoryID = emailCategoryId
		language := enums.PageLanguageEN
		
		t.Run("successfully get email content by category and language", func(t *testing.T) {
			mockService.On("GetEmailContentByCategoryIDAndLanguage", mock.AnythingOfType("string"), mock.AnythingOfType("enums.PageLanguage")).Return([]*dto.EmailContentResponse{
				{
					ID:              mockEmailContent.ID.String(),
					EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
					Language:        mockEmailContent.Language,
					Label:           mockEmailContent.Label,
				},
			}, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/email-contents/category/%s/language/%s", emailCategoryIdStr, language), nil)
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to get email content by category and language", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetEmailContentByCategoryIDAndLanguage", mock.AnythingOfType("string"), mock.AnythingOfType("enums.PageLanguage")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/cms/email-contents/category/%s/language/%s", emailCategoryIdStr, language), nil)
				
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})
}