package repositories

import (
	"github.com/MadManJJ/cms-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSCategoryTypeRepositoryInterface interface {
	Create(categoryType *models.CategoryType) (*models.CategoryType, error)
	FindByID(id uuid.UUID) (*models.CategoryType, error)
	FindByCode(code string) (*models.CategoryType, error)
	FindAll(isActive *bool) ([]models.CategoryType, error)
	Update(id uuid.UUID, updates map[string]interface{}) (*models.CategoryType, error)
	Delete(id uuid.UUID) error
	IsTypeCodeUnique(typeCode string, excludeID uuid.UUID) (bool, error)
}

type cmsCategoryTypeRepository struct {
	db *gorm.DB
}

func NewCMSCategoryTypeRepository(db *gorm.DB) CMSCategoryTypeRepositoryInterface {
	return &cmsCategoryTypeRepository{db: db}
}

func (r *cmsCategoryTypeRepository) Create(categoryType *models.CategoryType) (*models.CategoryType, error) {
	err := r.db.Create(categoryType).Error
	if err != nil {
		return nil, err
	}
	return categoryType, nil
}

func (r *cmsCategoryTypeRepository) FindByID(id uuid.UUID) (*models.CategoryType, error) {
	var ct models.CategoryType
	err := r.db.First(&ct, "id = ?", id).Error
	if err != nil {
		return nil, err
	}	
	return &ct, nil
}

func (r *cmsCategoryTypeRepository) FindByCode(code string) (*models.CategoryType, error) {
	var ct models.CategoryType
	err := r.db.Where("type_code = ?", code).First(&ct).Error
	if err != nil {
		return nil, err
	}		
	return &ct, nil
}

func (r *cmsCategoryTypeRepository) FindAll(isActive *bool) ([]models.CategoryType, error) {
	var cts []models.CategoryType
	query := r.db
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	err := query.Find(&cts).Error
	if err != nil {
		return nil, err
	}	
	return cts, nil
}

func (r *cmsCategoryTypeRepository) Update(id uuid.UUID, updates map[string]interface{}) (*models.CategoryType, error) {
	var ct models.CategoryType

	if err := r.db.First(&ct, "id = ?", id).Error; err != nil {
		return nil, err
	}

	err := r.db.Model(&ct).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return &ct, nil
}

func (r *cmsCategoryTypeRepository) Delete(id uuid.UUID) error {
	var ct models.CategoryType
	if err := r.db.First(&ct, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return r.db.Delete(&models.CategoryType{}, "id = ?", id).Error
}

func (r *cmsCategoryTypeRepository) IsTypeCodeUnique(typeCode string, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.CategoryType{}).Where("type_code = ?", typeCode)

	if excludeID != uuid.Nil {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
