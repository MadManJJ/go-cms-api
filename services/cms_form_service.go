package services

import (
	"errors"
	"fmt"
	"math"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Interface ---
type CMSFormServiceInterface interface {
	CreateNewForm(req dto.CreateFormRequest) (*dto.FormResponse, error)
	GetFormDetails(formID uuid.UUID) (*dto.FormResponse, error)
	GetFormStructure(formID uuid.UUID) (*dto.FormResponse, error)
	GetAllForms(filter dto.FormListFilter) (*dto.PaginatedFormListResponse, error)
	UpdateExistingForm(formID uuid.UUID, req dto.UpdateFormRequest) (*dto.FormResponse, error)
	DeleteExistingForm(formID uuid.UUID) error
}

// --- Implementation ---
type cmsFormService struct {
	db                *gorm.DB
	formRepo          repositories.FormRepositoryInterface
	validate          *validator.Validate
	emailCategoryRepo repositories.EmailCategoryRepositoryInterface
	cfg               *config.Config
}

func NewCMSFormService(
	db *gorm.DB,
	formRepo repositories.FormRepositoryInterface,
	emailCategoryRepo repositories.EmailCategoryRepositoryInterface,
	receivedCfg *config.Config,
) CMSFormServiceInterface {
	return &cmsFormService{
		db:                db,
		formRepo:          formRepo,
		validate:          validator.New(),
		emailCategoryRepo: emailCategoryRepo,
		cfg:               receivedCfg,
	}
}

func (s *cmsFormService) CreateNewForm(req dto.CreateFormRequest) (*dto.FormResponse, error) {
	// 1. Validate field key uniqueness in request
	fieldKeysInRequest := make(map[string]bool)
	for _, sectionReq := range req.Sections {
		for _, fieldReq := range sectionReq.Fields {
			if fieldKeysInRequest[fieldReq.FieldKey] {
				return nil, fmt.Errorf("%w: duplicate field_key '%s' in request", errs.ErrBadRequest, fieldReq.FieldKey)
			}
			fieldKeysInRequest[fieldReq.FieldKey] = true
		}
	}

	// 2. Build the GORM model from the DTO
	var emailCatIDPtr *uuid.UUID
	if req.EmailCategoryID != nil && *req.EmailCategoryID != "" {
		parsedUUID, err := uuid.Parse(*req.EmailCategoryID)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid email_category_id format: %v", errs.ErrBadRequest, err)
		}
		_, err = s.emailCategoryRepo.FindByID(parsedUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("%w: email_category_id '%s' not found", errs.ErrBadRequest, *req.EmailCategoryID)
			}
			return nil, fmt.Errorf("service: failed to verify email category: %w", err)
		}
		emailCatIDPtr = &parsedUUID
	}

	var languagePtr *enums.PageLanguage
	if req.Language != nil && *req.Language != "" {
		lang := enums.PageLanguage(*req.Language)
		languagePtr = &lang
	}

	formModel := &models.Form{
		Name:            req.Name,
		Slug:            helpers.GenerateSlug(req.Name),
		Description:     req.Description,
		EmailCategoryID: emailCatIDPtr,
		Language:        languagePtr,
		Sections:        s.buildFormSectionsFromDTO(req.Sections),
	}

	// 3. Execute creation in a transaction
	var createdForm *models.Form
	err := s.executeInTransaction(func(tx *gorm.DB) error {
		var txErr error
		createdForm, txErr = s.formRepo.CreateForm(tx, formModel)
		return txErr
	})
	if err != nil {
		return nil, fmt.Errorf("service: failed to create form: %w", err)
	}

	return s.mapFormModelToResponse(createdForm), nil
}

func (s *cmsFormService) buildFormSectionsFromDTO(sectionDTOs []dto.FormSectionRequest) []models.FormSection {
	if len(sectionDTOs) == 0 {
		return nil
	}
	formSections := make([]models.FormSection, len(sectionDTOs))
	for i, sectionReq := range sectionDTOs {
		formSections[i] = models.FormSection{
			Title:       sectionReq.Title,
			Description: sectionReq.Description,
			OrderIndex:  i + 1,
			Fields:      s.buildFormFieldsFromDTO(sectionReq.Fields),
		}
	}
	return formSections
}

func (s *cmsFormService) buildFormFieldsFromDTO(fieldDTOs []dto.FormFieldRequest) []models.FormField {
	if len(fieldDTOs) == 0 {
		return nil
	}
	formFields := make([]models.FormField, len(fieldDTOs))
	for i, fieldReq := range fieldDTOs {
		formFields[i] = models.FormField{
			Label:        fieldReq.Label,
			FieldKey:     fieldReq.FieldKey,
			FieldType:    enums.FormFieldType(fieldReq.FieldType),
			IsRequired:   fieldReq.IsRequired,
			Placeholder:  fieldReq.Placeholder,
			DefaultValue: fieldReq.DefaultValue,
			OrderIndex:   i + 1, // หรือใช้ OrderIndex จาก request
			Properties:   fieldReq.Properties,
			Display:      fieldReq.Display,
		}
	}
	return formFields
}
func (s *cmsFormService) GetFormDetails(formID uuid.UUID) (*dto.FormResponse, error) {
	form, err := s.formRepo.GetFormByID(formID)
	if err != nil {
		// Let the handler deal with specific error types (e.g., 404 vs 500).
		return nil, err
	}
	return s.mapFormModelToResponse(form), nil
}

func (s *cmsFormService) GetFormStructure(formID uuid.UUID) (*dto.FormResponse, error) {
	form, err := s.formRepo.GetFormStructure(formID)
	if err != nil {
		return nil, err
	}
	return s.mapFormModelToResponse(form), nil
}

func (s *cmsFormService) GetAllForms(filter dto.FormListFilter) (*dto.PaginatedFormListResponse, error) {
	// Validate filter DTO.
	if err := s.validate.Struct(filter); err != nil {
		return nil, fmt.Errorf("%w: invalid list filter parameters: %v", errs.ErrBadRequest, err)
	}

	// The repository handles pagination and sorting.
	forms, totalItems, err := s.formRepo.ListForms(filter)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get forms list from repository: %w", err)
	}

	// Map models to response DTOs.
	formListItems := s.mapFormModelsToListItemResponses(forms)

	// Create pagination metadata.
	meta := s.createPaginationMeta(totalItems, filter)

	return &dto.PaginatedFormListResponse{
		Data: formListItems,
		Meta: meta,
	}, nil
}

func (s *cmsFormService) UpdateExistingForm(formID uuid.UUID, req dto.UpdateFormRequest) (*dto.FormResponse, error) {
	// 1. Validate field key uniqueness in the request DTO
	fieldKeysInRequest := make(map[string]bool)
	for _, sectionReq := range req.Sections {
		for _, fieldReq := range sectionReq.Fields {
			if fieldKeysInRequest[fieldReq.FieldKey] {
				return nil, fmt.Errorf("%w: duplicate field_key '%s' in request", errs.ErrBadRequest, fieldReq.FieldKey)
			}
			fieldKeysInRequest[fieldReq.FieldKey] = true
		}
	}

	// 2. Fetch the existing form to ensure it exists
	existingForm, err := s.formRepo.GetFormByID(formID)
	if err != nil {
		return nil, err // Let handler manage not found error
	}

	// 3. Build the GORM model for update from the DTO
	var emailCatIDPtr *uuid.UUID
	if req.EmailCategoryID != nil && *req.EmailCategoryID != "" {
		parsedUUID, err := uuid.Parse(*req.EmailCategoryID)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid email_category_id format: %v", errs.ErrBadRequest, err)
		}
		_, err = s.emailCategoryRepo.FindByID(parsedUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("%w: email_category_id '%s' not found", errs.ErrBadRequest, *req.EmailCategoryID)
			}
			return nil, fmt.Errorf("service: failed to verify email category: %w", err)
		}
		emailCatIDPtr = &parsedUUID
	}

	// Build the base model for the update, preserving non-updatable fields
	formToUpdate := &models.Form{
		ID:              formID,
		Name:            req.Name,
		Description:     req.Description,
		Slug:            helpers.GenerateSlug(req.Name),
		CreatedAt:       existingForm.CreatedAt,
		EmailCategoryID: emailCatIDPtr,
		Sections:        s.buildFormSectionsFromUpdateDTO(req.Sections),
	}

	// Handle the language pointer separately
	if req.Language != nil && *req.Language != "" {
		// If a new language is provided, create a new pointer to the correct enum type
		lang := enums.PageLanguage(*req.Language)
		formToUpdate.Language = &lang // Assign the address of the new language value
	} else {
		// If no new language is provided, keep the pointer from the existing record
		formToUpdate.Language = existingForm.Language
	}
	// --- END OF CORRECTION ---

	// 4. Execute update in a transaction
	var updatedForm *models.Form
	err = s.executeInTransaction(func(tx *gorm.DB) error {
		var txErr error
		updatedForm, txErr = s.formRepo.UpdateForm(tx, formToUpdate)
		return txErr
	})

	if err != nil {
		return nil, fmt.Errorf("service: failed to update form: %w", err)
	}

	// 5. Map the updated model to a response DTO
	return s.mapFormModelToResponse(updatedForm), nil
}

// You might need a new helper for Update DTO as well
func (s *cmsFormService) buildFormSectionsFromUpdateDTO(sectionDTOs []dto.UpdateFormSectionRequest) []models.FormSection {
	if len(sectionDTOs) == 0 {
		return nil
	}
	formSections := make([]models.FormSection, len(sectionDTOs))
	for i, sectionReq := range sectionDTOs {
		formSections[i] = models.FormSection{
			Title:       sectionReq.Title,
			Description: sectionReq.Description,
			OrderIndex:  i + 1,
			Fields:      s.buildFormFieldsFromUpdateDTO(sectionReq.Fields),
		}
	}
	return formSections
}

func (s *cmsFormService) buildFormFieldsFromUpdateDTO(fieldDTOs []dto.UpdateFormFieldRequest) []models.FormField {
	if len(fieldDTOs) == 0 {
		return nil
	}
	formFields := make([]models.FormField, len(fieldDTOs))
	for i, fieldReq := range fieldDTOs {
		formFields[i] = models.FormField{
			Label:        fieldReq.Label,
			FieldKey:     fieldReq.FieldKey,
			FieldType:    enums.FormFieldType(fieldReq.FieldType),
			IsRequired:   fieldReq.IsRequired,
			Placeholder:  fieldReq.Placeholder,
			DefaultValue: fieldReq.DefaultValue,
			OrderIndex:   i + 1,
			Properties:   fieldReq.Properties,
			Display:      fieldReq.Display,
		}
	}
	return formFields
}
func (s *cmsFormService) DeleteExistingForm(formID uuid.UUID) error {
	err := s.executeInTransaction(func(tx *gorm.DB) error {
		return s.formRepo.DeleteForm(tx, formID)
	})
	if err != nil {
		return fmt.Errorf("service: failed to delete form: %w", err)
	}
	return nil
}

// --- Private Helper Methods ---

// executeInTransaction simplifies running operations in a GORM transaction.
func (s *cmsFormService) executeInTransaction(fn func(*gorm.DB) error) error {
	return s.db.Transaction(fn)
}

// validateRequestAndFieldKeys combines validation steps for cleaner calls.
func (s *cmsFormService) validateRequestAndFieldKeys(req models.Form) error {
	if err := s.validate.Struct(req); err != nil {
		return fmt.Errorf("%w: %v", errs.ErrBadRequest, err)
	}

	fieldKeysInRequest := make(map[string]bool)
	for _, section := range req.Sections {
		for _, field := range section.Fields {
			if fieldKeysInRequest[field.FieldKey] {
				return fmt.Errorf("%w: duplicate field_key '%s' in request", errs.ErrBadRequest, field.FieldKey)
			}
			fieldKeysInRequest[field.FieldKey] = true
		}
	}
	return nil
}

// buildFormSections and buildFormFields build the nested structure, setting order indices if needed.
func (s *cmsFormService) buildFormSections(sections []models.FormSection, formID uuid.UUID) []models.FormSection {
	if len(sections) == 0 {
		return nil
	}
	formSections := make([]models.FormSection, len(sections))
	for i, sectionReq := range sections {
		formSections[i] = models.FormSection{
			FormID:      formID,
			Title:       sectionReq.Title,
			Description: sectionReq.Description,
			OrderIndex:  i + 1, // Ensure consistent ordering
			Fields:      s.buildFormFields(sectionReq.Fields),
		}
	}
	return formSections
}

func (s *cmsFormService) buildFormFields(fields []models.FormField) []models.FormField {
	if len(fields) == 0 {
		return nil
	}
	formFields := make([]models.FormField, len(fields))
	for i, fieldReq := range fields {
		formFields[i] = models.FormField{
			Label:        fieldReq.Label,
			FieldKey:     fieldReq.FieldKey,
			FieldType:    fieldReq.FieldType,
			IsRequired:   fieldReq.IsRequired,
			Placeholder:  fieldReq.Placeholder,
			DefaultValue: fieldReq.DefaultValue,
			OrderIndex:   i + 1,
			Properties:   fieldReq.Properties,
			Display:      fieldReq.Display,
		}
	}
	return formFields
}

// --- DTO Mapping ---

func (s *cmsFormService) mapFormModelToResponse(form *models.Form) *dto.FormResponse {
	if form == nil {
		return nil
	}
	var emailCatIDStrPtr *string
	if form.EmailCategoryID != nil {
		idStr := form.EmailCategoryID.String()
		emailCatIDStrPtr = &idStr
	}
	var LanguageStrPtr *string
	if form.Language != nil {
		langStr := string(*form.Language)
		LanguageStrPtr = &langStr
	}
	return &dto.FormResponse{
		ID:              form.ID,
		Name:            form.Name,
		Slug:            form.Slug,
		Description:     form.Description,
		CreatedAt:       form.CreatedAt,
		UpdatedAt:       form.UpdatedAt,
		EmailCategoryID: emailCatIDStrPtr,
		Language:        LanguageStrPtr,
		Sections:        s.mapSectionsToResponse(form.Sections),
	}

}

func (s *cmsFormService) mapSectionsToResponse(sections []models.FormSection) []dto.FormSectionResponse {
	if len(sections) == 0 {
		return nil
	}
	respSections := make([]dto.FormSectionResponse, len(sections))
	for i, section := range sections {
		respSections[i] = dto.FormSectionResponse{
			ID:          section.ID,
			Title:       section.Title,
			Description: section.Description,
			OrderIndex:  section.OrderIndex,
			Fields:      s.mapFieldsToResponse(section.Fields),
		}
	}
	return respSections
}

func (s *cmsFormService) mapFieldsToResponse(fields []models.FormField) []dto.FormFieldResponse {
	if len(fields) == 0 {
		return []dto.FormFieldResponse{}
	}
	respFields := make([]dto.FormFieldResponse, len(fields))
	for i, field := range fields {
		respFields[i] = dto.FormFieldResponse{
			ID:           field.ID,
			Label:        field.Label,
			FieldKey:     field.FieldKey,
			FieldType:    string(field.FieldType),
			IsRequired:   field.IsRequired,
			Placeholder:  field.Placeholder,
			DefaultValue: field.DefaultValue,
			OrderIndex:   field.OrderIndex,
			Properties:   field.Properties,
			Display:      field.Display,
		}
	}
	return respFields
}

func (s *cmsFormService) mapFormModelsToListItemResponses(forms []models.Form) []dto.FormListItemResponse {
	formListItems := make([]dto.FormListItemResponse, len(forms))
	for i, form := range forms {
		formListItems[i] = dto.FormListItemResponse{
			ID:        form.ID,
			Name:      form.Name,
			Slug:      form.Slug,
			CreatedAt: form.CreatedAt,
			UpdatedAt: form.UpdatedAt,
		}
	}
	return formListItems
}

func (s *cmsFormService) createPaginationMeta(totalItems int64, filter dto.FormListFilter) dto.PaginationMeta {
	currentPage := 1
	if filter.Page != nil && *filter.Page > 0 {
		currentPage = *filter.Page
	}
	itemsPerPage := repositories.DefaultLimit
	if filter.ItemsPerPage != nil && *filter.ItemsPerPage > 0 {
		itemsPerPage = *filter.ItemsPerPage
		if itemsPerPage > repositories.MaxLimit {
			itemsPerPage = repositories.MaxLimit
		}
	}
	totalPages := 0
	if totalItems > 0 && itemsPerPage > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))
	}
	return dto.PaginationMeta{
		TotalItems:   totalItems,
		ItemsPerPage: itemsPerPage,
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
	}
}
