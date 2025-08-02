package repositories

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormRepositoryInterface interface {
	CreateForm(tx *gorm.DB, form *models.Form) (*models.Form, error)
	GetFormByID(formID uuid.UUID) (*models.Form, error)
	GetFormStructure(formID uuid.UUID) (*models.Form, error)
	ListForms(filter dto.FormListFilter) ([]models.Form, int64, error)
	UpdateForm(tx *gorm.DB, form *models.Form) (*models.Form, error)
	DeleteForm(tx *gorm.DB, formID uuid.UUID) error
	CheckFieldKeyExistsInForm(formID uuid.UUID, fieldKey string, excludeFieldID *uuid.UUID) (bool, error)
	GetFormWithFields(formID uuid.UUID) (*models.Form, error)
}

type formRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) FormRepositoryInterface {
	return &formRepository{db: db}
}

// Constants for better maintainability
const (
	DefaultLimit = 10
	MaxLimit     = 100
)

func (r *formRepository) CreateForm(tx *gorm.DB, form *models.Form) (*models.Form, error) {
	if err := tx.Create(form).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") && strings.Contains(err.Error(), "forms_slug_key") {
			return nil, fmt.Errorf("%w: slug '%s' already exists: %v", errs.ErrBadRequest, form.Slug, err)
		}
		return nil, fmt.Errorf("failed to create form with associations: %w", err)
	}

	return r.getFormWithFullAssociations(tx, form.ID)
}

func (r *formRepository) GetFormByID(formID uuid.UUID) (*models.Form, error) {

	return r.getFormWithFullAssociations(r.db, formID)
}

func (r *formRepository) GetFormStructure(formID uuid.UUID) (*models.Form, error) {
	var form models.Form
	err := r.db.
		Preload("Sections", r.orderSections).
		Preload("Sections.Fields", r.orderFields).
		First(&form, "id = ?", formID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get form structure by ID %s: %w", formID, err)
	}
	return &form, nil
}

func (r *formRepository) ListForms(filter dto.FormListFilter) ([]models.Form, int64, error) {
	var forms []models.Form
	var totalItems int64

	query := r.buildFilterQuery(filter)

	if err := query.Model(&models.Form{}).Count(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count forms: %w", err)
	}

	query = r.applySortingAndPagination(query, filter)

	if err := query.Find(&forms).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list forms: %w", err)
	}

	return forms, totalItems, nil
}

func (r *formRepository) UpdateForm(tx *gorm.DB, form *models.Form) (*models.Form, error) {
	if form.ID == uuid.Nil {
		return nil, errors.New("form ID is required for update")
	}

	// Delete existing form sections first
	if err := tx.Where("form_id = ?", form.ID).Delete(&models.FormSection{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete existing form sections: %w", err)
	}

	// Update form metadata including the new slug
	if err := r.updateFormMetadata(tx, form); err != nil {
		return nil, err
	}

	// Recreate form sections
	if err := r.recreateFormSections(tx, form); err != nil {
		return nil, err
	}

	return r.getFormWithFullAssociations(tx, form.ID)
}

func (r *formRepository) updateFormMetadata(tx *gorm.DB, form *models.Form) error {
	updates := map[string]interface{}{
		"name":              form.Name,
		"slug":              form.Slug,
		"description":       form.Description,
		"updated_at":        gorm.Expr("NOW()"),
		"email_category_id": form.EmailCategoryID,
		"language":          form.Language,
	}

	if err := tx.Model(&models.Form{}).Where("id = ?", form.ID).Updates(updates).Error; err != nil {
		// Check for unique constraint violation on slug
		if strings.Contains(err.Error(), "unique constraint") && strings.Contains(err.Error(), "forms_slug_key") {
			return fmt.Errorf("%w: slug '%s' already exists: %v", errs.ErrBadRequest, form.Slug, err)
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return fmt.Errorf("%w: invalid reference, e.g., email_category_id not found: %v", errs.ErrBadRequest, err)
		}
		return fmt.Errorf("failed to update form metadata: %w", err)
	}
	return nil
}
func (r *formRepository) recreateFormSections(tx *gorm.DB, form *models.Form) error {

	if len(form.Sections) > 0 {

		for i := range form.Sections {
			form.Sections[i].FormID = form.ID
		}
		if err := tx.Create(&form.Sections).Error; err != nil {
			return fmt.Errorf("failed to recreate sections: %w", err)
		}
	}
	return nil
}

func (r *formRepository) DeleteForm(tx *gorm.DB, formID uuid.UUID) error {
	// Use hard delete to completely remove the record and free up the slug
	result := tx.Unscoped().Delete(&models.Form{}, formID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete form %s: %w", formID, result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// --- Private helper methods ---

func (r *formRepository) getFormWithFullAssociations(db *gorm.DB, formID uuid.UUID) (*models.Form, error) {
	var form models.Form
	err := db.
		Preload("Sections", r.orderSections).
		Preload("Sections.Fields", r.orderFields).
		Preload("EmailCategory").
		First(&form, "id = ?", formID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get form by ID %s with full associations: %w", formID, err)
	}
	return &form, nil
}

func (r *formRepository) orderSections(db *gorm.DB) *gorm.DB {
	return db.Order("form_sections.order_index ASC")
}

func (r *formRepository) orderFields(db *gorm.DB) *gorm.DB {
	return db.Order("form_fields.order_index ASC")
}

func (r *formRepository) CheckFieldKeyExistsInForm(formID uuid.UUID, fieldKey string, excludeFieldID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.FormField{}).
		Joins("JOIN form_sections ON form_sections.id = form_fields.section_id").
		Where("form_sections.form_id = ? AND form_fields.field_key = ?", formID, fieldKey)

	if excludeFieldID != nil && *excludeFieldID != uuid.Nil {
		query = query.Where("form_fields.id != ?", *excludeFieldID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("error checking field_key existence: %w", err)
	}

	return count > 0, nil
}

func (r *formRepository) GetFormWithFields(formID uuid.UUID) (*models.Form, error) {
	var form models.Form
	err := r.db.Preload("Sections.Fields").First(&form, "id = ?", formID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get form with fields: %w", err)
	}
	return &form, nil
}

func (r *formRepository) buildFilterQuery(filter dto.FormListFilter) *gorm.DB {
	query := r.db.Model(&models.Form{})

	if filter.Name != nil && *filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}

	if filter.CreatedAt != nil {
		startOfDay := filter.CreatedAt.Truncate(24 * time.Hour)
		endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)
		query = query.Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay)
	}

	return query
}

func (r *formRepository) applySortingAndPagination(query *gorm.DB, filter dto.FormListFilter) *gorm.DB {
	sortOrder := "updated_at DESC" // Default sort order (Newest)

	if filter.Sort != nil && *filter.Sort != "" {
		switch *filter.Sort {
		case "updated_at_asc":
			sortOrder = "updated_at ASC" // Oldest
		case "updated_at_desc":
			sortOrder = "updated_at DESC" // Newest
		case "name_asc":
			sortOrder = "name ASC" // A-Z
		case "name_desc":
			sortOrder = "name DESC" // Z-A
		}
	}
	query = query.Order(sortOrder)

	page := r.getPage(filter.Page)
	limit := r.getLimit(filter.ItemsPerPage)
	offset := (page - 1) * limit

	return query.Offset(offset).Limit(limit)
}

func (r *formRepository) getPage(page *int) int {
	if page != nil && *page > 0 {
		return *page
	}
	return 1
}

func (r *formRepository) getLimit(itemsPerPage *int) int {
	if itemsPerPage != nil && *itemsPerPage > 0 {
		if *itemsPerPage > MaxLimit {
			return MaxLimit
		}
		return *itemsPerPage
	}
	return DefaultLimit
}
