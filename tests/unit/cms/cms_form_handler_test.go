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

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCMSFormService struct {
	mock.Mock
}

func (m *MockCMSFormService) CreateNewForm(req dto.CreateFormRequest) (*dto.FormResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FormResponse), args.Error(1)
}

func (m *MockCMSFormService) GetFormDetails(formID uuid.UUID) (*dto.FormResponse, error) {
	args := m.Called(formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FormResponse), args.Error(1)
}

func (m *MockCMSFormService) GetFormStructure(formID uuid.UUID) (*dto.FormResponse, error) {
	args := m.Called(formID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FormResponse), args.Error(1)
}

func (m *MockCMSFormService) GetAllForms(filter dto.FormListFilter) (*dto.PaginatedFormListResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaginatedFormListResponse), args.Error(1)
}

func (m *MockCMSFormService) UpdateExistingForm(formID uuid.UUID, req dto.UpdateFormRequest) (*dto.FormResponse, error) {
	args := m.Called(formID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FormResponse), args.Error(1)
}

func (m *MockCMSFormService) DeleteExistingForm(formID uuid.UUID) error {
	args := m.Called(formID)
	return args.Error(0)
}

func TestFormHandler(t *testing.T) {
	mockService := &MockCMSFormService{}
	handler := cmsHandler.NewCMSFormHandler(mockService)
	userId := uuid.New()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"user_id": userId.String(),
		}
		c.Locals("user", claims)
		return c.Next()
	})

	app.Post("/api/v1/cms/forms", handler.HandleCreateForm)
	app.Get("/api/v1/cms/forms/:formId", handler.HandleGetForm)
	app.Get("/api/v1/cms/forms/:formId/structure", handler.HandleGetFormStructure)
	app.Get("/api/v1/cms/forms", handler.HandleListForms)
	app.Put("/api/v1/cms/forms/:formId", handler.HandleUpdateForm)
	app.Delete("/api/v1/cms/forms/:formId", handler.HandleDeleteForm)

	form := helpers.InitializeMockForm()
	body, err := json.Marshal(form)
	require.NoError(t, err)

	formFieldResponse := &dto.FormFieldResponse{
		ID:         form.Sections[0].Fields[0].ID,
		Label:      form.Sections[0].Fields[0].Label,
		FieldKey:   form.Sections[0].Fields[0].FieldKey,
		FieldType:  string(form.Sections[0].Fields[0].FieldType),
		IsRequired: form.Sections[0].Fields[0].IsRequired,
		OrderIndex: form.Sections[0].Fields[0].OrderIndex,
	}

	formSectionResponse := &dto.FormSectionResponse{
		ID:          form.Sections[0].ID,
		Title:       form.Sections[0].Title,
		Description: form.Sections[0].Description,
		OrderIndex:  form.Sections[0].OrderIndex,
		Fields:      []dto.FormFieldResponse{*formFieldResponse},
	}

	formResponseResponse := &dto.FormResponse{
		Name:        form.Name,
		Slug:        form.Slug,
		Description: form.Description,
		CreatedAt:   form.CreatedAt,
		UpdatedAt:   form.UpdatedAt,
		Sections:    []dto.FormSectionResponse{*formSectionResponse},
	}

	t.Run("POST /api/v1/cms/forms HandleCreateForm", func(t *testing.T) {

		t.Run("successfully create form", func(t *testing.T) {
			mockService.On("CreateNewForm", mock.AnythingOfType("dto.CreateFormRequest")).Return(formResponseResponse, nil)

			req := httptest.NewRequest("POST", "/api/v1/cms/forms", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to create form: invalid body", func(t *testing.T) {
			mockService.ExpectedCalls = nil

			// Invalid body
			body, err := json.Marshal("invalid body")
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/cms/forms", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to create form: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("CreateNewForm", mock.AnythingOfType("dto.CreateFormRequest")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("POST", "/api/v1/cms/forms", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("GET /api/v1/cms/forms/:formId HandleGetForm", func(t *testing.T) {
		formId := uuid.New()

		t.Run("successfully get form", func(t *testing.T) {
			mockService.On("GetFormDetails", mock.AnythingOfType("uuid.UUID")).Return(formResponseResponse, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/cms/forms/%s", formId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to get form: invalid form id", func(t *testing.T) {
			mockService.ExpectedCalls = nil

			// Invalid form ID
			invalidFormId := "invalid-form-id"
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/cms/forms/%s", invalidFormId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to get form: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFormDetails", mock.AnythingOfType("uuid.UUID")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/cms/forms/%s", formId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("GET /api/v1/cms/forms/:formId/structure HandleGetFormStructure", func(t *testing.T) {
		formId := uuid.New()

		t.Run("successfully get form structure", func(t *testing.T) {
			mockService.On("GetFormStructure", mock.AnythingOfType("uuid.UUID")).Return(formResponseResponse, nil)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/cms/forms/%s/structure", formId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to get form structure: invalid form id", func(t *testing.T) {
			mockService.ExpectedCalls = nil

			// Invalid form ID
			invalidFormId := "invalid-form-id"
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/cms/forms/%s/structure", invalidFormId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to get form structure: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetFormStructure", mock.AnythingOfType("uuid.UUID")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/cms/forms/%s/structure", formId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("GET /api/v1/cms/forms HandleListForms", func(t *testing.T) {
		formId := uuid.New()

		formListItemReponse := &dto.FormListItemResponse{
			ID:        formId,
			Name:      form.Name,
			Slug:      form.Slug,
			CreatedAt: form.CreatedAt,
			UpdatedAt: form.UpdatedAt,
		}

		paginationMeta := &dto.PaginationMeta{
			TotalItems:   1,
			ItemsPerPage: 10,
			CurrentPage:  1,
			TotalPages:   5,
		}

		paginatedResponse := &dto.PaginatedFormListResponse{
			Data: []dto.FormListItemResponse{*formListItemReponse},
			Meta: *paginationMeta,
		}

		t.Run("successfully list forms", func(t *testing.T) {
			mockService.On("GetAllForms", mock.AnythingOfType("dto.FormListFilter")).Return(paginatedResponse, nil)

			req := httptest.NewRequest("GET", "/api/v1/cms/forms", nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to list forms: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("GetAllForms", mock.AnythingOfType("dto.FormListFilter")).Return(nil, errs.ErrInternalServerError)

			req := httptest.NewRequest("GET", "/api/v1/cms/forms", nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("PUT /api/v1/cms/forms/:formId HandleUpdateForm", func(t *testing.T) {
		formId := uuid.New()
		form := helpers.InitializeMockForm()
		form.ID = formId

		updatedForm := *form
		updatedForm.Name = "Updated Form"

		formResponseResponse.Name = updatedForm.Name

		t.Run("successfully update form", func(t *testing.T) {
			mockService.On("UpdateExistingForm", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("dto.UpdateFormRequest")).Return(formResponseResponse, nil)

			reqBody, err := json.Marshal(&updatedForm)
			assert.NoError(t, err)

			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/cms/forms/%s", formId), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to update form: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("UpdateExistingForm", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("dto.UpdateFormRequest")).Return(nil, errs.ErrInternalServerError)

			reqBody, err := json.Marshal(&updatedForm)
			assert.NoError(t, err)

			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/cms/forms/%s", formId), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})

	t.Run("DELETE /api/v1/cms/forms/:formId HandleDeleteForm", func(t *testing.T) {
		formId := uuid.New()

		t.Run("successfully delete form", func(t *testing.T) {
			mockService.On("DeleteExistingForm", mock.AnythingOfType("uuid.UUID")).Return(nil)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/cms/forms/%s", formId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
			mockService.AssertExpectations(t)
		})

		t.Run("failed to delete form: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("DeleteExistingForm", mock.AnythingOfType("uuid.UUID")).Return(errs.ErrInternalServerError)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/cms/forms/%s", formId), nil)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	})
}
