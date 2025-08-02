package repositories

import (
	"fmt"

	"github.com/MadManJJ/cms-api/dto" // ยังคง import dto สำหรับ filters ถ้าจำเป็น
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CMSCategoryRepositoryInterface defines methods for managing categories (which are now category details).
type CMSCategoryRepositoryInterface interface {
	CreateCategory(category *models.Category) (*models.Category, error)
	GetCategoryByID(categoryID uuid.UUID) (*models.Category, error)
	UpdateCategory(category *models.Category) (*models.Category, error)
	DeleteCategory(categoryID uuid.UUID) error
	CountCategoriesByTypeAndLanguage(categoryTypeID uuid.UUID) (map[string]int, error)
	ListCategoriesByFilter(filters dto.CategoryFilter) ([]models.Category, error)
}

type cmsCategoryRepository struct {
	db *gorm.DB
}

func NewCMSCategoryRepository(db *gorm.DB) CMSCategoryRepositoryInterface {
	return &cmsCategoryRepository{db: db}
}

// CreateCategory creates a new category (detail).
func (r *cmsCategoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	// category.ID (UUID) should be set by BeforeCreate hook or DB default
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	// Reload to get CategoryType properly
	return r.GetCategoryByID(category.ID)
}

func (r *cmsCategoryRepository) GetCategoryByID(categoryID uuid.UUID) (*models.Category, error) {
	var cat models.Category
	// Preload CategoryType as Category (Detail) belongs to a CategoryType
	err := r.db.Preload("CategoryType").First(&cat, "id = ?", categoryID).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// ListCategoriesByFilter lists categories based on provided filters.
func (r *cmsCategoryRepository) ListCategoriesByFilter(filters dto.CategoryFilter) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Model(&models.Category{}).Preload("CategoryType")

	if filters.CategoryTypeID != nil && *filters.CategoryTypeID != "" {
		catTypeUUID, err := uuid.Parse(*filters.CategoryTypeID)
		if err != nil {
			return nil, fmt.Errorf("repository: invalid category_type_id filter format: %w", err)
		}
		query = query.Where("category_type_id = ?", catTypeUUID)
	}

	if filters.LanguageCode != nil && *filters.LanguageCode != "" {
		query = query.Where("language_code = ?", *filters.LanguageCode)
	}

	if filters.Name != nil && *filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+*filters.Name+"%")
	}

	if filters.PublishStatus != nil && *filters.PublishStatus != "" {
		query = query.Where("publish_status = ?", *filters.PublishStatus)
	}

	err := query.Order("weight asc, created_at asc").Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("error finding categories with filter: %w", err)
	}
	return categories, nil
}

func (r *cmsCategoryRepository) DeleteCategory(categoryID uuid.UUID) error {
	// Directly delete the Category (Detail) itself
	result := r.db.Delete(&models.Category{}, categoryID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete category ID %s: %w", categoryID, result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateCategory updates a category (detail).
func (r *cmsCategoryRepository) UpdateCategory(category *models.Category) (*models.Category, error) {
	// Ensure the category ID is set
	if err := r.db.Omit("CategoryType").Save(category).Error; err != nil {
		return nil, fmt.Errorf("failed to update category (ID: %s): %w", category.ID, err)
	}
	// Reload to get CategoryType properly after update
	return r.GetCategoryByID(category.ID)
}

// CountCategoriesByTypeAndLanguage counts categories for a specific type, grouped by language.
func (r *cmsCategoryRepository) CountCategoriesByTypeAndLanguage(categoryTypeID uuid.UUID) (map[string]int, error) {
	result := make(map[string]int)
	type langCount struct {
		LanguageCode enums.PageLanguage
		Count        int
	}
	var counts []langCount

	err := r.db.Model(&models.Category{}).
		Select("language_code, COUNT(*) as count").
		Where("category_type_id = ?", categoryTypeID).
		Group("language_code").
		Scan(&counts).Error

	if err != nil {
		return nil, fmt.Errorf("error counting categories for type %s by language: %w", categoryTypeID, err)
	}

	for _, lc := range counts {
		result[string(lc.LanguageCode)] = lc.Count
	}

	if _, ok := result[string(enums.PageLanguageTH)]; !ok {
		result[string(enums.PageLanguageTH)] = 0
	}
	if _, ok := result[string(enums.PageLanguageEN)]; !ok {
		result[string(enums.PageLanguageEN)] = 0
	}
	return result, nil
}
