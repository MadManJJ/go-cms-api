package tests

import (
	"testing"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockCMSEmailContentRepo struct {
	create func(content *models.EmailContent) (*models.EmailContent, error)
	findByID func(id uuid.UUID) (*models.EmailContent, error)
	findByCategoryIDAndLanguageAndLabel func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error)
	listByFilters func(filters dto.EmailContentFilter) ([]models.EmailContent, error)
	update func(content *models.EmailContent) (*models.EmailContent, error)
	delete func(id uuid.UUID) error
	deleteByCategoryID func(categoryID uuid.UUID) error	
	findEmailContentByCategoryIDAndLanguage func(categoryID uuid.UUID, language enums.PageLanguage) ([]models.EmailContent, error)
}

func (m *MockCMSEmailContentRepo) Create(content *models.EmailContent) (*models.EmailContent, error) {
	return m.create(content)
}

func (m *MockCMSEmailContentRepo) FindByID(id uuid.UUID) (*models.EmailContent, error) {
	return m.findByID(id)
}

func (m *MockCMSEmailContentRepo) FindByCategoryIDAndLanguageAndLabel(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
	return m.findByCategoryIDAndLanguageAndLabel(categoryID, language, label)
}

func (m *MockCMSEmailContentRepo) ListByFilters(filters dto.EmailContentFilter) ([]models.EmailContent, error) {
	return m.listByFilters(filters)
}

func (m *MockCMSEmailContentRepo) Update(content *models.EmailContent) (*models.EmailContent, error) {
	return m.update(content)
}

func (m *MockCMSEmailContentRepo) Delete(id uuid.UUID) error {
	return m.delete(id)
}

func (m *MockCMSEmailContentRepo) DeleteByCategoryID(categoryID uuid.UUID) error {
	return m.deleteByCategoryID(categoryID)
}

func (m *MockCMSEmailContentRepo) FindEmailContentByCategoryIDAndLanguage(categoryID uuid.UUID, language enums.PageLanguage) ([]models.EmailContent, error) {
	return m.findEmailContentByCategoryIDAndLanguage(categoryID, language)
}

func TestCMSService_CreateEmailContent(t *testing.T) {
	mockEmailContent := helpers.InitializeMockEmailContent()
	emailContentId := uuid.New()
	emailCategoryId := uuid.New()
	mockEmailContent.ID = emailContentId
	mockEmailContent.EmailCategoryID = emailCategoryId

	emailContentReq := dto.CreateEmailContentRequest{
		EmailCategoryID: emailCategoryId.String(),
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

	t.Run("successfully create email content", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByCategoryIDAndLanguageAndLabel: func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
				return nil, nil
			},
			create: func(content *models.EmailContent) (*models.EmailContent, error) {
				return mockEmailContent, nil
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return mockEmailContent.EmailCategory, nil
			},
		}
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		actualEmailContentResponse, err := service.CreateContent(emailContentReq)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailContentResponse)
		assert.Equal(t, mockEmailContent.EmailCategoryID.String(), actualEmailContentResponse.EmailCategoryID)
		assert.Equal(t, mockEmailContent.Language, actualEmailContentResponse.Language)
		assert.Equal(t, mockEmailContent.Label, actualEmailContentResponse.Label)
		assert.Equal(t, mockEmailContent.SendTo, actualEmailContentResponse.SendTo)
		assert.Equal(t, mockEmailContent.CcEmail, actualEmailContentResponse.CcEmail)
		assert.Equal(t, mockEmailContent.BccEmail, actualEmailContentResponse.BccEmail)
		assert.Equal(t, mockEmailContent.SendFromEmail, actualEmailContentResponse.SendFromEmail)
		assert.Equal(t, mockEmailContent.SendFromName, actualEmailContentResponse.SendFromName)
		assert.Equal(t, mockEmailContent.Subject, actualEmailContentResponse.Subject)
		assert.Equal(t, mockEmailContent.TopImgLink, actualEmailContentResponse.TopImgLink)
		assert.Equal(t, mockEmailContent.Header, actualEmailContentResponse.Header)
		assert.Equal(t, mockEmailContent.Paragraph, actualEmailContentResponse.Paragraph)
		assert.Equal(t, mockEmailContent.Footer, actualEmailContentResponse.Footer)
		assert.Equal(t, mockEmailContent.FooterImageLink, actualEmailContentResponse.FooterImageLink)
	})

	t.Run("failed to create email content: internal server error", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByCategoryIDAndLanguageAndLabel: func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
				return nil, nil
			},
			create: func(content *models.EmailContent) (*models.EmailContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return mockEmailContent.EmailCategory, nil
			},
		}
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		actualEmailContentResponse, err := service.CreateContent(emailContentReq)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContentResponse)
	})	
}

func TestCMSService_GetEmailContentByID(t *testing.T) {
	mockEmailContent := helpers.InitializeMockEmailContent()
	emailContentId := uuid.New()
	mockEmailContent.ID = emailContentId
	
	emailContentResponse := &dto.EmailContentResponse{
		ID: mockEmailContent.ID.String(),
		EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
		Language: mockEmailContent.Language,
		Label: mockEmailContent.Label,
		SendTo: mockEmailContent.SendTo,
		CcEmail: mockEmailContent.CcEmail,
		BccEmail: mockEmailContent.BccEmail,
		SendFromEmail: mockEmailContent.SendFromEmail,
		SendFromName: mockEmailContent.SendFromName,
		Subject: mockEmailContent.Subject,
		TopImgLink: mockEmailContent.TopImgLink,
		Header: mockEmailContent.Header,
		Paragraph: mockEmailContent.Paragraph,
		Footer: mockEmailContent.Footer,
		FooterImageLink: mockEmailContent.FooterImageLink,
	}

	t.Run("successfully get email content by id", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByID: func(id uuid.UUID) (*models.EmailContent, error) {
				return mockEmailContent, nil
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}		

		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)

		actualEmailContentResponse, err := service.GetContentByID(emailContentId.String())
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailContentResponse)
		assert.Equal(t, emailContentResponse.ID, actualEmailContentResponse.ID)
		assert.Equal(t, emailContentResponse.Language, actualEmailContentResponse.Language)
		assert.Equal(t, emailContentResponse.Label, actualEmailContentResponse.Label)
		assert.Equal(t, emailContentResponse.SendTo, actualEmailContentResponse.SendTo)
		assert.Equal(t, emailContentResponse.CcEmail, actualEmailContentResponse.CcEmail)
		assert.Equal(t, emailContentResponse.BccEmail, actualEmailContentResponse.BccEmail)
		assert.Equal(t, emailContentResponse.SendFromEmail, actualEmailContentResponse.SendFromEmail)
		assert.Equal(t, emailContentResponse.SendFromName, actualEmailContentResponse.SendFromName)
		assert.Equal(t, emailContentResponse.Subject, actualEmailContentResponse.Subject)
		assert.Equal(t, emailContentResponse.TopImgLink, actualEmailContentResponse.TopImgLink)
		assert.Equal(t, emailContentResponse.Header, actualEmailContentResponse.Header)
		assert.Equal(t, emailContentResponse.Paragraph, actualEmailContentResponse.Paragraph)
		assert.Equal(t, emailContentResponse.Footer, actualEmailContentResponse.Footer)
		assert.Equal(t, emailContentResponse.FooterImageLink, actualEmailContentResponse.FooterImageLink)
	})

	t.Run("failed to get email content by id", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByID: func(id uuid.UUID) (*models.EmailContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}		

		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)

		actualEmailContentResponse, err := service.GetContentByID(emailContentId.String())
		assert.Error(t, err)
		assert.Nil(t, actualEmailContentResponse)
	})	
}

func TestCMSService_GetEmailContentByCategoryAndLangAndLabel(t *testing.T) {
	mockEmailContent := helpers.InitializeMockEmailContent()
	emailContentId := uuid.New()
	emailCategoryId := uuid.New()
	mockEmailContent.ID = emailContentId
	
	emailContentResponse := &dto.EmailContentResponse{
		ID: mockEmailContent.ID.String(),
		EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
		Language: mockEmailContent.Language,
		Label: mockEmailContent.Label,
		SendTo: mockEmailContent.SendTo,
		CcEmail: mockEmailContent.CcEmail,
		BccEmail: mockEmailContent.BccEmail,
		SendFromEmail: mockEmailContent.SendFromEmail,
		SendFromName: mockEmailContent.SendFromName,
		Subject: mockEmailContent.Subject,
		TopImgLink: mockEmailContent.TopImgLink,
		Header: mockEmailContent.Header,
		Paragraph: mockEmailContent.Paragraph,
		Footer: mockEmailContent.Footer,
		FooterImageLink: mockEmailContent.FooterImageLink,
	}

	t.Run("successfully get email content by category and lang and label", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByCategoryIDAndLanguageAndLabel: func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
				return mockEmailContent, nil
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}	
		
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)

		actualEmailContentResponse, err := service.GetContentByCategoryAndLangAndLabel(emailCategoryId.String(), mockEmailContent.Language, mockEmailContent.Label)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailContentResponse)
		assert.Equal(t, emailContentResponse.ID, actualEmailContentResponse.ID)
		assert.Equal(t, emailContentResponse.Language, actualEmailContentResponse.Language)
		assert.Equal(t, emailContentResponse.Label, actualEmailContentResponse.Label)
		assert.Equal(t, emailContentResponse.SendTo, actualEmailContentResponse.SendTo)
		assert.Equal(t, emailContentResponse.CcEmail, actualEmailContentResponse.CcEmail)
		assert.Equal(t, emailContentResponse.BccEmail, actualEmailContentResponse.BccEmail)
		assert.Equal(t, emailContentResponse.SendFromEmail, actualEmailContentResponse.SendFromEmail)
		assert.Equal(t, emailContentResponse.SendFromName, actualEmailContentResponse.SendFromName)
		assert.Equal(t, emailContentResponse.Subject, actualEmailContentResponse.Subject)
		assert.Equal(t, emailContentResponse.TopImgLink, actualEmailContentResponse.TopImgLink)
		assert.Equal(t, emailContentResponse.Header, actualEmailContentResponse.Header)
		assert.Equal(t, emailContentResponse.Paragraph, actualEmailContentResponse.Paragraph)
		assert.Equal(t, emailContentResponse.Footer, actualEmailContentResponse.Footer)
		assert.Equal(t, emailContentResponse.FooterImageLink, actualEmailContentResponse.FooterImageLink)
	})

	t.Run("failed to get email content by category and lang and label", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByCategoryIDAndLanguageAndLabel: func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}			
		
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)

		actualEmailContentResponse, err := service.GetContentByCategoryAndLangAndLabel(emailCategoryId.String(), mockEmailContent.Language, mockEmailContent.Label)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContentResponse)
	})
}

func TestCMSService_ListContents(t *testing.T) {
	mockEmailContent := helpers.InitializeMockEmailContent()
	emailContentId := uuid.New()
	emailCategoryId := uuid.New()
	mockEmailContent.ID = emailContentId
	now := time.Now()

	mockEmailContent.CreatedAt = now
	mockEmailContent.UpdatedAt = now
	
	emailContentResponses := []dto.EmailContentResponse{
		{
			ID: mockEmailContent.ID.String(),
			EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
			Language: mockEmailContent.Language,
			Label: mockEmailContent.Label,
			SendTo: mockEmailContent.SendTo,
			CcEmail: mockEmailContent.CcEmail,
			BccEmail: mockEmailContent.BccEmail,
			SendFromEmail: mockEmailContent.SendFromEmail,
			SendFromName: mockEmailContent.SendFromName,
			Subject: mockEmailContent.Subject,
			TopImgLink: mockEmailContent.TopImgLink,
			Header: mockEmailContent.Header,
			Paragraph: mockEmailContent.Paragraph,
			Footer: mockEmailContent.Footer,
			FooterImageLink: mockEmailContent.FooterImageLink,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}	

	emailCategoryIdStr := emailCategoryId.String()
	language := mockEmailContent.Language
	label := mockEmailContent.Label

	filter := dto.EmailContentFilter{
		EmailCategoryID: &emailCategoryIdStr,
		Language:        &language,
		Label:           &label,
	}

	t.Run("successfully list email contents", func(t *testing.T) { 
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			listByFilters: func(filters dto.EmailContentFilter) ([]models.EmailContent, error) {
				return []models.EmailContent{*mockEmailContent}, nil
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}		
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		actualEmailContentResponses, err := service.ListContents(filter)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailContentResponses)
		assert.Equal(t, emailContentResponses, actualEmailContentResponses)
	})

	t.Run("failed to list email contents", func(t *testing.T) { 
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			listByFilters: func(filters dto.EmailContentFilter) ([]models.EmailContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}			
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		actualEmailContentResponses, err := service.ListContents(filter)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContentResponses)
	})
}

func TestCMSService_UpdateContent(t *testing.T) {
	mockEmailContent := helpers.InitializeMockEmailContent()
	emailContentId := uuid.New()
	emailCategoryId := uuid.New()
	mockEmailContent.ID = emailContentId
	mockEmailContent.EmailCategoryID = emailCategoryId
	mockEmailContent.EmailCategory.ID = emailCategoryId

	updatedEmailContent := *mockEmailContent
	updatedEmailContent.ID = uuid.New()
	updatedEmailContent.EmailCategoryID = emailCategoryId
	updatedEmailContent.EmailCategory = mockEmailContent.EmailCategory
	updatedEmailContent.Language = enums.PageLanguage("en")
	updatedEmailContent.Label = "Updated Label"
	updatedEmailContent.SendTo = "updatedSendTo"
	updatedEmailContent.CcEmail = "updatedCcEmail"
	updatedEmailContent.BccEmail = "updatedBccEmail"
	
	updateEmailContentRequest := dto.UpdateEmailContentRequest{
		Language: &updatedEmailContent.Language,
		Label: &updatedEmailContent.Label,
		SendTo: &updatedEmailContent.SendTo,
		CcEmail: &updatedEmailContent.CcEmail,
		BccEmail: &updatedEmailContent.BccEmail,
		SendFromEmail: &updatedEmailContent.SendFromEmail,
		SendFromName: &mockEmailContent.SendFromName,
		Subject: &mockEmailContent.Subject,
		TopImgLink: &mockEmailContent.TopImgLink,
		Header: &mockEmailContent.Header,
		Paragraph: &mockEmailContent.Paragraph,
		Footer: &mockEmailContent.Footer,
		FooterImageLink: &mockEmailContent.FooterImageLink,
	}

	t.Run("successfully update email content", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByID: func(id uuid.UUID) (*models.EmailContent, error) {
				return mockEmailContent, nil
			},
			findByCategoryIDAndLanguageAndLabel: func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
				return mockEmailContent, nil
			},
			update: func(emailContent *models.EmailContent) (*models.EmailContent, error) {
				return &updatedEmailContent, nil
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}	

		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)			

		actualEmailContentResponse, err := service.UpdateContent(emailContentId.String(), updateEmailContentRequest)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailContentResponse)
		assert.Equal(t, updatedEmailContent.ID.String(), actualEmailContentResponse.ID)
		assert.Equal(t, updatedEmailContent.Language, actualEmailContentResponse.Language)
		assert.Equal(t, updatedEmailContent.Label, actualEmailContentResponse.Label)
		assert.Equal(t, updatedEmailContent.SendTo, actualEmailContentResponse.SendTo)
		assert.Equal(t, updatedEmailContent.CcEmail, actualEmailContentResponse.CcEmail)
		assert.Equal(t, updatedEmailContent.BccEmail, actualEmailContentResponse.BccEmail)
	})

	t.Run("failed to update email content: internal server error", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findByID: func(id uuid.UUID) (*models.EmailContent, error) {
				return mockEmailContent, nil
			},
			findByCategoryIDAndLanguageAndLabel: func(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
				return mockEmailContent, nil
			},
			update: func(emailContent *models.EmailContent) (*models.EmailContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}			
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)			

		actualEmailContentResponse, err := service.UpdateContent(emailContentId.String(), updateEmailContentRequest)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContentResponse)
	})	
}

func TestCMSService_DeleteContent(t *testing.T) {
	emailContentId := uuid.New()

	t.Run("successfully delete email content", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			delete: func(id uuid.UUID) error {
				return nil
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}		
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		err := service.DeleteContent(emailContentId.String())
		assert.NoError(t, err)
	})

	t.Run("failed to delete email content: internal server error", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			delete: func(id uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}		
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		err := service.DeleteContent(emailContentId.String())
		assert.Error(t, err)
	})
}

func TestCMSService_GetEmailContentByCategoryIDAndLanguage(t *testing.T) {
	mockEmailContent := helpers.InitializeMockEmailContent()
	emailContentId := uuid.New()
	emailCategoryId := uuid.New()
	mockEmailContent.ID = emailContentId
	language := mockEmailContent.Language
	now := time.Now()

	mockEmailContent.CreatedAt = now
	mockEmailContent.UpdatedAt = now
	
	// Change from normal value to pointer for some reason
	emailContentResponses := []*dto.EmailContentResponse{
		{
			ID: mockEmailContent.ID.String(),
			EmailCategoryID: mockEmailContent.EmailCategoryID.String(),
			Language: mockEmailContent.Language,
			Label: mockEmailContent.Label,
			SendTo: mockEmailContent.SendTo,
			CcEmail: mockEmailContent.CcEmail,
			BccEmail: mockEmailContent.BccEmail,
			SendFromEmail: mockEmailContent.SendFromEmail,
			SendFromName: mockEmailContent.SendFromName,
			Subject: mockEmailContent.Subject,
			TopImgLink: mockEmailContent.TopImgLink,
			Header: mockEmailContent.Header,
			Paragraph: mockEmailContent.Paragraph,
			Footer: mockEmailContent.Footer,
			FooterImageLink: mockEmailContent.FooterImageLink,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}	

	t.Run("successfully get email content by category ID and language", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findEmailContentByCategoryIDAndLanguage: func(categoryID uuid.UUID, language enums.PageLanguage) ([]models.EmailContent, error) {
				return []models.EmailContent{*mockEmailContent}, nil
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}			
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		actualEmailContentResponses, err := service.GetEmailContentByCategoryIDAndLanguage(emailCategoryId.String(), language)
		assert.NoError(t, err)
		assert.NotNil(t, actualEmailContentResponses)
		assert.Equal(t, emailContentResponses, actualEmailContentResponses)
	})

	t.Run("failed to get email content by category ID and language", func(t *testing.T) {
		mockEmailContentRepo := &MockCMSEmailContentRepo{
			findEmailContentByCategoryIDAndLanguage: func(categoryID uuid.UUID, language enums.PageLanguage) ([]models.EmailContent, error) {
				return nil, errs.ErrInternalServerError
			},
		}
	
		mockEmailCategoryRepo := &MockCMSEmailCategoryRepo{}			
	
		service := services.NewEmailContentService(mockEmailContentRepo, mockEmailCategoryRepo)
	
		actualEmailContentResponses, err := service.GetEmailContentByCategoryIDAndLanguage(emailCategoryId.String(), language)
		assert.Error(t, err)
		assert.Nil(t, actualEmailContentResponses)
	})	
}