package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MadManJJ/cms-api/dto" // For filter DTO
	"github.com/MadManJJ/cms-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaFileRepositoryInterface interface {
	Create(file *models.MediaFile) (*models.MediaFile, error)
	FindByID(id uuid.UUID) (*models.MediaFile, error)
	FindByNameAndPath(name string, path string) (*models.MediaFile, error)
	List(filter dto.MediaFileListFilter) ([]models.MediaFile, int64, error)
	Delete(id uuid.UUID) error
}

type mediaFileRepository struct {
	db *gorm.DB
}

func NewMediaFileRepository(db *gorm.DB) MediaFileRepositoryInterface {
	return &mediaFileRepository{db: db}
}

func (r *mediaFileRepository) Create(file *models.MediaFile) (*models.MediaFile, error) {
	if err := r.db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *mediaFileRepository) FindByID(id uuid.UUID) (*models.MediaFile, error) {
	var file models.MediaFile
	if err := r.db.First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// FindByNameAndPath searches for a file by its name and path (if path is used).

func (r *mediaFileRepository) FindByNameAndPath(name string, path string) (*models.MediaFile, error) {
	var file models.MediaFile
	// Assuming 'Path' field exists in MediaFile model as discussed.
	// If not, the query needs adjustment or path parameter might be irrelevant at DB level.
	// For now, let's assume we're only checking by name if Path is not part of the model.
	// query := r.db.Where("name = ?", name)
	// if path != "" { // If your model has a Path field to store subdirectories
	//  query = query.Where("path = ?", path)
	// }
	// err := query.First(&file).Error
	// For the current model in api.zip:
	err := r.db.Where("name = ?", name).First(&file).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not an error if not found, service layer handles logic
		}
		return nil, err
	}
	return &file, nil
}

func (r *mediaFileRepository) List(filter dto.MediaFileListFilter) ([]models.MediaFile, int64, error) {
	var files []models.MediaFile
	var total int64

	query := r.db.Model(&models.MediaFile{})

	if filter.Search != nil && *filter.Search != "" {
		// PostgreSQL case-insensitive search
		query = query.Where("name ILIKE ?", "%"+strings.ToLower(*filter.Search)+"%")
	}

	// Count total records matching the filter (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting media files: %w", err)
	}

	// Apply sorting
	if filter.SortBy != nil && *filter.SortBy != "" {
		orderDirection := "ASC"
		if filter.Order != nil && strings.ToLower(*filter.Order) == "desc" {
			orderDirection = "DESC"
		}
		// Sanitize SortBy to prevent SQL injection if it were raw SQL. GORM handles this better.
		validSortColumns := map[string]string{"name": "name", "created_at": "created_at"}
		if dbSortColumn, ok := validSortColumns[*filter.SortBy]; ok {
			query = query.Order(fmt.Sprintf("%s %s", dbSortColumn, orderDirection))
		} else {
			query = query.Order("created_at DESC") // Default sort
		}
	} else {
		query = query.Order("created_at DESC") // Default sort
	}

	// Apply pagination
	page := 1
	if filter.Page > 0 {
		page = filter.Page
	}
	pageSize := 10 // Default page size
	if filter.PageSize > 0 {
		pageSize = filter.PageSize
	}
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	if err := query.Find(&files).Error; err != nil {
		return nil, total, fmt.Errorf("error listing media files: %w", err)
	}

	return files, total, nil
}

func (r *mediaFileRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.MediaFile{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Or a custom "not found" error
	}
	return nil
}
