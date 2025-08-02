package tests

import (
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCMSFormSubmissionRepo struct {
	createFormSubmission func(formSubmission *models.FormSubmission) (*models.FormSubmission, error)
	getFormSubmissions   func(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error)
	getFormSubmission    func(submissionId uuid.UUID) (*models.FormSubmission, error)
	getEmailContentsFormFormId func(formId uuid.UUID) ([]*models.EmailContent, error)
}

func (m *MockCMSFormSubmissionRepo) CreateFormSubmission(formSubmission *models.FormSubmission) (*models.FormSubmission, error) {
	return m.createFormSubmission(formSubmission)
}

func (m *MockCMSFormSubmissionRepo) GetFormSubmissions(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error) {
	return m.getFormSubmissions(formId, sort, page, limit)
}

func (m *MockCMSFormSubmissionRepo) GetFormSubmission(submissionId uuid.UUID) (*models.FormSubmission, error) {
	return m.getFormSubmission(submissionId)
}

func (m *MockCMSFormSubmissionRepo) GetEmailContentsFormFormId(formId uuid.UUID) ([]*models.EmailContent, error) {
	return m.getEmailContentsFormFormId(formId)
}

type MockEmailSendingService struct {
	mock.Mock
}

func (m *MockEmailSendingService) SendEmail(req dto.SendEmailRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func TestCMSService_CreateFormSubmission(t *testing.T) {
	mockFormSubmission := helpers.InitializeMockFormSubmission()
	formId := uuid.New()
	emailSendingService := &MockEmailSendingService{}	

	t.Run("successfully create form submission", func(t *testing.T) {
		repo := &MockCMSFormSubmissionRepo{
			createFormSubmission: func(formSubmission *models.FormSubmission) (*models.FormSubmission, error) {
				return mockFormSubmission, nil
			},
			getEmailContentsFormFormId: func(formId uuid.UUID) ([]*models.EmailContent, error) {
				return []*models.EmailContent{}, nil
			},
		}	
		emailSendingService.On("SendEmail", mock.AnythingOfType("dto.SendEmailRequest")).Return(nil)
	
		service := services.NewCMSFormSubmissionService(repo, emailSendingService)
	
		actualFormSubmission, err := service.CreateFormSubmission(formId, mockFormSubmission)
		assert.NoError(t, err)
		assert.Equal(t, mockFormSubmission, actualFormSubmission)	
	})

	t.Run("failed to create form submission", func(t *testing.T) {
		repo := &MockCMSFormSubmissionRepo{
			createFormSubmission: func(formSubmission *models.FormSubmission) (*models.FormSubmission, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailSendingService.ExpectedCalls = nil
		emailSendingService.On("SendEmail", mock.AnythingOfType("dto.SendEmailRequest")).Return(nil)
	
		service := services.NewCMSFormSubmissionService(repo, emailSendingService)
	
		actualFormSubmission, err := service.CreateFormSubmission(formId, mockFormSubmission)
		assert.Error(t, err)	
		assert.Nil(t, actualFormSubmission)
	})
}

func TestCMSService_GetFormSubmissions(t *testing.T) {
	mockFormSubmission := helpers.InitializeMockFormSubmission()
	formId := uuid.New()
	sort := "mock sort"	
	emailSendingService := &MockEmailSendingService{}	

	t.Run("successfully get form submissions", func(t *testing.T) {
		repo := &MockCMSFormSubmissionRepo{
			getFormSubmissions: func(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error) {
				return []*models.FormSubmission{mockFormSubmission}, int64(1), nil
			},
		}	
	
		service := services.NewCMSFormSubmissionService(repo, emailSendingService)
	
		actualFormSubmissions, totalCount, err := service.GetFormSubmissions(formId, sort, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, []*models.FormSubmission{mockFormSubmission}, actualFormSubmissions)	
	})

	t.Run("failed to get form submissions", func(t *testing.T) {
		repo := &MockCMSFormSubmissionRepo{
			getFormSubmissions: func(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error) {
				return nil, int64(0), errs.ErrInternalServerError
			},
		}
		
		service := services.NewCMSFormSubmissionService(repo, emailSendingService)
	
		actualFormSubmissions, totalCount, err := service.GetFormSubmissions(formId, sort, 1, 10)
		assert.Error(t, err)	
		assert.Nil(t, actualFormSubmissions)
		assert.Equal(t, int64(0), totalCount)
	})
}

func TestCMSService_GetFormSubmission(t *testing.T) {
	mockFormSubmission := helpers.InitializeMockFormSubmission()
	submissionId := uuid.New()
	emailSendingService := &MockEmailSendingService{}	
	t.Run("successfully get form submission", func(t *testing.T) {
		repo := &MockCMSFormSubmissionRepo{
			getFormSubmission: func(submissionId uuid.UUID) (*models.FormSubmission, error) {
				return mockFormSubmission, nil
			},
		}
	
		service := services.NewCMSFormSubmissionService(repo, emailSendingService)
	
		actualFormSubmission, err := service.GetFormSubmission(submissionId)
		assert.NoError(t, err)
		assert.Equal(t, mockFormSubmission, actualFormSubmission)	
	})

	t.Run("failed to get form submission", func(t *testing.T) {
		repo := &MockCMSFormSubmissionRepo{
			getFormSubmission: func(submissionId uuid.UUID) (*models.FormSubmission, error) {
				return nil, errs.ErrInternalServerError
			},
		}
	
		service := services.NewCMSFormSubmissionService(repo, emailSendingService)
	
		actualFormSubmission, err := service.GetFormSubmission(submissionId)
		assert.Error(t, err)	
		assert.Nil(t, actualFormSubmission)
	})
}