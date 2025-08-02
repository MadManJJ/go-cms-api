package repositories

import (
	"fmt"
	"log"
	"strings"

	"github.com/MadManJJ/cms-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormSubmissionRepositoryInterface interface {
	CreateFormSubmission(formSubmission *models.FormSubmission) (*models.FormSubmission, error)
	GetFormSubmissions(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error)
	GetFormSubmission(submissionId uuid.UUID) (*models.FormSubmission, error)
	GetEmailContentsFormFormId(formId uuid.UUID) ([]*models.EmailContent, error)
}

type FormSubmissionRepository struct {
	db *gorm.DB
}

func NewFormSubmissionRepository(db *gorm.DB) FormSubmissionRepositoryInterface {
	return &FormSubmissionRepository{db: db}
}

func (r *FormSubmissionRepository) CreateFormSubmission(formSubmission *models.FormSubmission) (*models.FormSubmission, error) {
	if err := r.db.Create(formSubmission).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Form").First(formSubmission, "id = ?", formSubmission.ID).Error; err != nil {

		log.Printf("ERROR: Failed to preload Form for submission ID %s: %v", formSubmission.ID, err)
		return nil, err
	}

	return formSubmission, nil
}

func (r *FormSubmissionRepository) GetFormSubmissions(formId uuid.UUID, sort string, page, limit int) ([]*models.FormSubmission, int64, error) {
	db := r.db.Model(&models.FormSubmission{}).Where("form_id = ?", formId)

	var totalCount int64
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, totalCount, err
	}

	db = db.
		Preload("Form")

	// Sort & Pagination
	if sort != "" {
		sortParts := strings.Split(sort, ":")
		if len(sortParts) == 2 {
			column := sortParts[0]
			direction := strings.ToUpper(sortParts[1])
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC" // fallback to ASC if direction is invalid
			}
			db = db.Order(fmt.Sprintf("%s %s", column, direction))
		} else {
			// fallback if the format is wrong
			db = db.Order("created_at DESC")
		}
	} else {
		db = db.Order("created_at DESC")
	}

	offset := (page - 1) * limit
	db = db.Offset(offset).Limit(limit)

	var formSubmissions []*models.FormSubmission
	if err := db.Find(&formSubmissions).Error; err != nil {
		return nil, 0, err
	}

	return formSubmissions, totalCount, nil
}

func (r *FormSubmissionRepository) GetFormSubmission(submissionId uuid.UUID) (*models.FormSubmission, error) {
	var formSubmission models.FormSubmission
	if err := r.db.Preload("Form").First(&formSubmission, "id = ?", submissionId).Error; err != nil {
		return nil, err
	}

	return &formSubmission, nil
}

func (r *FormSubmissionRepository) GetEmailContentsFormFormId(formId uuid.UUID) ([]*models.EmailContent, error) {
	var form models.Form
	if err := r.db.Select("email_category_id").Where("id = ?", formId).First(&form).Error; err != nil {
		return nil, err
	}

	var emailContents []*models.EmailContent
	if err := r.db.Preload("EmailCategory").Where("email_category_id = ?", form.EmailCategoryID).Find(&emailContents).Error; err != nil {
		return nil, err
	}

	return emailContents, nil
}
