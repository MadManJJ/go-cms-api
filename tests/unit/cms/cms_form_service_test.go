package tests

import (
	"testing"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type MockFormRepository struct {
	createForm                func(tx *gorm.DB, form *models.Form) (*models.Form, error)
	getFormByID               func(formID uuid.UUID) (*models.Form, error)
	getFormStructure          func(formID uuid.UUID) (*models.Form, error)
	listForms                 func(filter dto.FormListFilter) ([]models.Form, int64, error)
	updateForm                func(tx *gorm.DB, form *models.Form) (*models.Form, error)
	deleteForm                func(tx *gorm.DB, formID uuid.UUID) error
	checkFieldKeyExistsInForm func(formID uuid.UUID, fieldKey string, excludeFieldID *uuid.UUID) (bool, error)
	getFormWithFields         func(formID uuid.UUID) (*models.Form, error)
}

func (m *MockFormRepository) CreateForm(tx *gorm.DB, form *models.Form) (*models.Form, error) {
	return m.createForm(tx, form)
}

func (m *MockFormRepository) GetFormByID(formID uuid.UUID) (*models.Form, error) {
	return m.getFormByID(formID)
}

func (m *MockFormRepository) GetFormStructure(formID uuid.UUID) (*models.Form, error) {
	return m.getFormStructure(formID)
}

func (m *MockFormRepository) ListForms(filter dto.FormListFilter) ([]models.Form, int64, error) {
	return m.listForms(filter)
}

func (m *MockFormRepository) UpdateForm(tx *gorm.DB, form *models.Form) (*models.Form, error) {
	return m.updateForm(tx, form)
}

func (m *MockFormRepository) DeleteForm(tx *gorm.DB, formID uuid.UUID) error {
	return m.deleteForm(tx, formID)
}

func (m *MockFormRepository) CheckFieldKeyExistsInForm(formID uuid.UUID, fieldKey string, excludeFieldID *uuid.UUID) (bool, error) {
	return m.checkFieldKeyExistsInForm(formID, fieldKey, excludeFieldID)
}

func (m *MockFormRepository) GetFormWithFields(formID uuid.UUID) (*models.Form, error) {
	return m.getFormWithFields(formID)
}

func TestCMSService_CreateForm(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cfg := config.New()

	mockForm := helpers.InitializeMockForm()
	formSection := mockForm.Sections[0]
	formField := formSection.Fields[0]

	emailCategoryId := uuid.New().String()
	language := string(enums.PageLanguageEN)

	formFieldRequest := dto.FormFieldRequest{
		Label:        formField.Label,
		FieldKey:     formField.FieldKey,
		FieldType:    string(formField.FieldType),
		IsRequired:   formField.IsRequired,
		Placeholder:  formField.Placeholder,
		DefaultValue: formField.DefaultValue,
		OrderIndex:   formField.OrderIndex,
		Properties:   formField.Properties,
		Display:      formField.Display,
	}

	formSectionRequest := dto.FormSectionRequest{
		Title:       formSection.Title,
		Description: formSection.Description,
		Fields:      []dto.FormFieldRequest{formFieldRequest},
	}

	createFormReq := dto.CreateFormRequest{
		Name:            mockForm.Name,
		Description:     mockForm.Description,
		EmailCategoryID: &emailCategoryId,
		Language:        &language,
		Sections:        []dto.FormSectionRequest{formSectionRequest},
	}

	t.Run("successfully create form", func(t *testing.T) {
		mock.ExpectBegin()
		// CreateForm
		mock.ExpectCommit()

		formRepo := &MockFormRepository{
			createForm: func(tx *gorm.DB, form *models.Form) (*models.Form, error) {
				return mockForm, nil
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return nil, nil
			},
		}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.CreateNewForm(createFormReq)
		assert.NoError(t, err)
		assert.NotNil(t, actualForm)
		assert.Equal(t, mockForm.ID, actualForm.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to create form", func(t *testing.T) {
		mock.ExpectBegin()
		// CreateForm
		mock.ExpectRollback()

		formRepo := &MockFormRepository{
			createForm: func(tx *gorm.DB, form *models.Form) (*models.Form, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return nil, nil
			},
		}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.CreateNewForm(createFormReq)
		assert.Error(t, err)
		assert.Nil(t, actualForm)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSService_GetFormDetails(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cfg := config.New()

	mockForm := helpers.InitializeMockForm()

	t.Run("successfully get form details", func(t *testing.T) {
		formRepo := &MockFormRepository{
			getFormByID: func(formID uuid.UUID) (*models.Form, error) {
				return mockForm, nil
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.GetFormDetails(mockForm.ID)
		assert.NoError(t, err)
		assert.NotNil(t, actualForm)
		assert.Equal(t, mockForm.ID, actualForm.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get form details", func(t *testing.T) {
		formRepo := &MockFormRepository{
			getFormByID: func(formID uuid.UUID) (*models.Form, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.GetFormDetails(mockForm.ID)
		assert.Error(t, err)
		assert.Nil(t, actualForm)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSService_GetFormStructure(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cfg := config.New()

	mockForm := helpers.InitializeMockForm()

	t.Run("successfully get form structure", func(t *testing.T) {
		formRepo := &MockFormRepository{
			getFormStructure: func(formID uuid.UUID) (*models.Form, error) {
				return mockForm, nil
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.GetFormStructure(mockForm.ID)
		assert.NoError(t, err)
		assert.NotNil(t, actualForm)
		assert.Equal(t, mockForm.ID, actualForm.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get form structure", func(t *testing.T) {
		formRepo := &MockFormRepository{
			getFormStructure: func(formID uuid.UUID) (*models.Form, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.GetFormStructure(mockForm.ID)
		assert.Error(t, err)
		assert.Nil(t, actualForm)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSService_GetAllForms(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cfg := config.New()

	mockForm := helpers.InitializeMockForm()

	sort := "name_asc"
	page := 1
	itemsPerPage := 10
	filter := dto.FormListFilter{
		Name:         &mockForm.Name,
		Sort:         &sort,
		Page:         &page,
		ItemsPerPage: &itemsPerPage,
	}

	t.Run("successfully get all forms", func(t *testing.T) {
		formRepo := &MockFormRepository{
			listForms: func(filter dto.FormListFilter) ([]models.Form, int64, error) {
				return []models.Form{*mockForm}, int64(1), nil
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForms, err := service.GetAllForms(filter)
		assert.NoError(t, err)
		assert.NotNil(t, actualForms)
		assert.Equal(t, mockForm.ID, actualForms.Data[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to get all forms", func(t *testing.T) {
		formRepo := &MockFormRepository{
			listForms: func(filter dto.FormListFilter) ([]models.Form, int64, error) {
				return nil, int64(0), errs.ErrInternalServerError
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForms, err := service.GetAllForms(filter)
		assert.Error(t, err)
		assert.Nil(t, actualForms)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSService_UpdateExistingForm(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cfg := config.New()

	mockForm := helpers.InitializeMockForm()

	updatedForm := *mockForm
	updatedForm.Name = "Updated Form"

	updatedFormSection := updatedForm.Sections[0]
	updatedTitle := "Updated Section"
	updatedFormSection.Title = &updatedTitle

	updatedFormField := updatedFormSection.Fields[0]
	updatedFormField.Label = "Updated Field"

	emailCategoryId := uuid.New().String()
	language := string(enums.PageLanguageEN)

	updatedFormFieldRequest := dto.UpdateFormFieldRequest{
		Label:        updatedFormField.Label,
		FieldKey:     updatedFormField.FieldKey,
		FieldType:    string(updatedFormField.FieldType),
		IsRequired:   updatedFormField.IsRequired,
		Placeholder:  updatedFormField.Placeholder,
		DefaultValue: updatedFormField.DefaultValue,
		OrderIndex:   updatedFormField.OrderIndex,
		Properties:   updatedFormField.Properties,
		Display:      updatedFormField.Display,
	}

	updatedFormSectionRequest := dto.UpdateFormSectionRequest{
		Title:       updatedFormSection.Title,
		Description: updatedFormSection.Description,
		Fields:      []dto.UpdateFormFieldRequest{updatedFormFieldRequest},
	}

	updatedFormReq := dto.UpdateFormRequest{
		Name:            updatedForm.Name,
		Description:     updatedForm.Description,
		EmailCategoryID: &emailCategoryId,
		Language:        &language,
		Sections:        []dto.UpdateFormSectionRequest{updatedFormSectionRequest},
	}

	t.Run("successfully update existing form", func(t *testing.T) {
		mock.ExpectBegin()
		// UpdateForm
		mock.ExpectCommit()

		formRepo := &MockFormRepository{
			getFormByID: func(formID uuid.UUID) (*models.Form, error) {
				return mockForm, nil
			},
			updateForm: func(tx *gorm.DB, form *models.Form) (*models.Form, error) {
				return &updatedForm, nil
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return nil, nil
			},
		}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.UpdateExistingForm(mockForm.ID, updatedFormReq)
		assert.NoError(t, err)
		assert.NotNil(t, actualForm)
		assert.Equal(t, updatedForm.ID, actualForm.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to update existing form", func(t *testing.T) {
		mock.ExpectBegin()
		// UpdateForm
		mock.ExpectRollback()

		formRepo := &MockFormRepository{
			getFormByID: func(formID uuid.UUID) (*models.Form, error) {
				return mockForm, nil
			},
			updateForm: func(tx *gorm.DB, form *models.Form) (*models.Form, error) {
				return nil, errs.ErrInternalServerError
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{
			findByID: func(id uuid.UUID) (*models.EmailCategory, error) {
				return nil, nil
			},
		}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		actualForm, err := service.UpdateExistingForm(mockForm.ID, updatedFormReq)
		assert.Error(t, err)
		assert.Nil(t, actualForm)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCMSService_DeleteExistingForm(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cfg := config.New()

	formId := uuid.New()

	t.Run("successfully delete and existing form", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		formRepo := &MockFormRepository{
			deleteForm: func(tx *gorm.DB, formID uuid.UUID) error {
				return nil
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		err := service.DeleteExistingForm(formId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to delete and existing form", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		formRepo := &MockFormRepository{
			deleteForm: func(tx *gorm.DB, formID uuid.UUID) error {
				return errs.ErrInternalServerError
			},
		}
		emailCategoryRepo := &MockCMSEmailCategoryRepo{}

		service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, cfg)

		err := service.DeleteExistingForm(formId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
