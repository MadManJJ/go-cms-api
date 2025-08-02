package repositories

import (
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailContentRepositoryInterface interface {
	Create(content *models.EmailContent) (*models.EmailContent, error)
	FindByID(id uuid.UUID) (*models.EmailContent, error)
	FindByCategoryIDAndLanguageAndLabel(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error)
	ListByFilters(filters dto.EmailContentFilter) ([]models.EmailContent, error)
	Update(content *models.EmailContent) (*models.EmailContent, error)
	Delete(id uuid.UUID) error
	DeleteByCategoryID(categoryID uuid.UUID) error
	FindEmailContentByCategoryIDAndLanguage(categoryID uuid.UUID, language enums.PageLanguage) ([]models.EmailContent, error)
}

type emailContentRepository struct {
	db *gorm.DB
}

func NewEmailContentRepository(db *gorm.DB) EmailContentRepositoryInterface {
	return &emailContentRepository{db: db}
}

func (r *emailContentRepository) Create(content *models.EmailContent) (*models.EmailContent, error) {
	if err := r.db.Create(content).Error; err != nil {
		return nil, err
	}
	// Preload category for response mapping
	return r.FindByID(content.ID)
}

func (r *emailContentRepository) FindByID(id uuid.UUID) (*models.EmailContent, error) {
	var content models.EmailContent
	if err := r.db.Preload("EmailCategory").First(&content, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *emailContentRepository) FindByCategoryIDAndLanguageAndLabel(categoryID uuid.UUID, language enums.PageLanguage, label string) (*models.EmailContent, error) {
	var content models.EmailContent
	err := r.db.Preload("EmailCategory").
		Where("email_category_id = ? AND language = ? AND label = ?", categoryID, language, label).
		First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *emailContentRepository) ListByFilters(filters dto.EmailContentFilter) ([]models.EmailContent, error) {
	var contents []models.EmailContent
	query := r.db.Model(&models.EmailContent{}).Preload("EmailCategory")

	if filters.EmailCategoryID != nil && *filters.EmailCategoryID != "" {
		catID, err := uuid.Parse(*filters.EmailCategoryID)
		if err == nil {
			query = query.Where("email_category_id = ?", catID)
		}
	}
	if filters.Language != nil && *filters.Language != "" {
		query = query.Where("language = ?", *filters.Language)
	}
	if filters.Label != nil && *filters.Label != "" {
		query = query.Where("label LIKE ?", "%"+*filters.Label+"%")
	}

	if err := query.Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}

func (r *emailContentRepository) Update(content *models.EmailContent) (*models.EmailContent, error) {
	if err := r.db.Save(content).Error; err != nil {
		return nil, err
	}
	// Preload category for response mapping
	return r.FindByID(content.ID)
}

func (r *emailContentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.EmailContent{}, "id = ?", id).Error
}

func (r *emailContentRepository) DeleteByCategoryID(categoryID uuid.UUID) error {
	return r.db.Where("email_category_id = ?", categoryID).Delete(&models.EmailContent{}).Error
}

func (r *emailContentRepository) FindEmailContentByCategoryIDAndLanguage(categoryID uuid.UUID, language enums.PageLanguage) ([]models.EmailContent, error) {
	var contents []models.EmailContent
	if err := r.db.Where("email_category_id = ? AND language = ?", categoryID, language).
		Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}
