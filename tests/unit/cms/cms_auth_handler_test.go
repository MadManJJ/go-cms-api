package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	cmsHandler "github.com/MadManJJ/cms-api/handlers/cms"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCMSAuthService struct {
	mock.Mock
}

func (m *MockCMSAuthService) RegisterUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}	
	return args.Get(0).(*models.User), args.Error(1)	
}

func (m *MockCMSAuthService) LoginUser(user *models.User) (*models.User, string, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}	
	return args.Get(0).(*models.User), args.Get(1).(string), args.Error(2)
}

func TestCMSAuthHandler(t *testing.T) {
	mockService := &MockCMSAuthService{}
	handler := cmsHandler.NewAuthCMSHandler(mockService)

	app := fiber.New()
	app.Post("/cms/auth/register", handler.HandleRegister)
	app.Post("/cms/auth/login", handler.HandleLogin)

	t.Run("POST /cms/auth/register HandleRegister", func(t *testing.T) {
		mockUser := helpers.InitializeMockUser()
		mockUserWithHashedPassword := helpers.InitializeMockUserWithHashedPassword()
		body, err := json.Marshal(mockUser)
		require.NoError(t, err)				
		
		t.Run("successfully register", func(t *testing.T) {
			mockService.On("RegisterUser", mockUser).Return(mockUserWithHashedPassword, nil)						

			req := httptest.NewRequest("POST", "/cms/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed to register: invalid body", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("RegisterUser", mockUser).Return(mockUserWithHashedPassword, nil)	

			invalidBody, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", "/cms/auth/register", bytes.NewReader(invalidBody))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed to register: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("RegisterUser", mockUser).Return(nil, errs.ErrInternalServerError)			

			req := httptest.NewRequest("POST", "/cms/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})

	t.Run("POST /cms/auth/login HandleLogin", func(t *testing.T) {
		mockUser := helpers.InitializeMockUser()
		mockUserWithHashedPassword := helpers.InitializeMockUserWithHashedPassword()
		body, err := json.Marshal(mockUser)
		require.NoError(t, err)	
		token := "fake token"
		
		t.Run("successfully login", func(t *testing.T) {
			mockService.On("LoginUser", mockUser).Return(mockUserWithHashedPassword, token, nil)		

			req := httptest.NewRequest("POST", "/cms/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	
		
		t.Run("failed to login: invalid body", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("LoginUser", mockUser).Return(mockUserWithHashedPassword, token, nil)	

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", "/cms/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})	

		t.Run("failed to login: invalid creadentials", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("LoginUser", mockUser).Return(nil, "", errs.ErrInvalidCredentials)	

			body, err := json.Marshal("invalid body")
			require.NoError(t, err)	

			req := httptest.NewRequest("POST", "/cms/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			mockService.AssertExpectations(t)
		})			
		
		t.Run("failed to login: internal server error", func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.On("LoginUser", mockUser).Return(nil, "", errs.ErrInternalServerError)	

			req := httptest.NewRequest("POST", "/cms/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
			mockService.AssertExpectations(t)
		})		
	})
}