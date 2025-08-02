package repositories

import (
	"errors"

	"github.com/MadManJJ/cms-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailCategoryRepositoryInterface interface {
	Create(category *models.EmailCategory) (*models.EmailCategory, error)
	FindByID(id uuid.UUID) (*models.EmailCategory, error)
	FindByTitle(title string) (*models.EmailCategory, error)
	FindAll() ([]models.EmailCategory, error)
	Update(category *models.EmailCategory) (*models.EmailCategory, error)
	Delete(id uuid.UUID) error
	IsTitleUnique(title string, excludeID uuid.UUID) (bool, error)
}

type emailCategoryRepository struct {
	db *gorm.DB
}

func NewEmailCategoryRepository(db *gorm.DB) EmailCategoryRepositoryInterface {
	return &emailCategoryRepository{db: db}
}

func (r *emailCategoryRepository) Create(category *models.EmailCategory) (*models.EmailCategory, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *emailCategoryRepository) FindByID(id uuid.UUID) (*models.EmailCategory, error) {
	var category models.EmailCategory
	if err := r.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *emailCategoryRepository) FindByTitle(title string) (*models.EmailCategory, error) {
	var category models.EmailCategory
	if err := r.db.Where("title = ?", title).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error here, service layer decides
		}
		return nil, err
	}
	return &category, nil
}

func (r *emailCategoryRepository) FindAll() ([]models.EmailCategory, error) {
	var categories []models.EmailCategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *emailCategoryRepository) Update(category *models.EmailCategory) (*models.EmailCategory, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *emailCategoryRepository) Delete(id uuid.UUID) error {
	// Add check if category is in use by EmailContent before deleting
	var count int64
	if err := r.db.Model(&models.EmailContent{}).Where("email_category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email category is in use and cannot be deleted")
	}

	if err := r.db.Delete(&models.EmailCategory{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *emailCategoryRepository) IsTitleUnique(title string, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.EmailCategory{}).Where("title = ?", title)
	if excludeID != uuid.Nil {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
