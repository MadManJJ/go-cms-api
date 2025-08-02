package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSCategoryServiceInterface interface {
	CreateCategory(req dto.CategoryCreateRequest) (*dto.CategoryResponse, error)
	GetCategoryByUUID(uuidStr string) (*dto.CategoryResponse, error)
	ListAllCategories(filter dto.CategoryFilter) ([]dto.CategoryResponse, error)
	UpdateCategoryByUUID(uuidStr string, req dto.CategoryUpdateRequest) (*dto.CategoryResponse, error)
	DeleteCategoryByUUID(uuidStr string) error
	MapCategoryModelToResponse(cat *models.Category) (*dto.CategoryResponse, error) // << เพิ่ม method นี้ใน Interface
}

type cmsCategoryService struct {
	categoryRepo     repositories.CMSCategoryRepositoryInterface
	categoryTypeRepo repositories.CMSCategoryTypeRepositoryInterface
}

func NewCMSCategoryService(
	categoryRepo repositories.CMSCategoryRepositoryInterface,
	categoryTypeRepo repositories.CMSCategoryTypeRepositoryInterface,
) CMSCategoryServiceInterface {
	return &cmsCategoryService{
		categoryRepo:     categoryRepo,
		categoryTypeRepo: categoryTypeRepo,
	}
}

// MapCategoryModelToResponse is now a public method of cmsCategoryService
func (s *cmsCategoryService) MapCategoryModelToResponse(cat *models.Category) (*dto.CategoryResponse, error) {
	if cat == nil {
		return nil, nil
	}

	response := &dto.CategoryResponse{
		ID:             cat.ID.String(),
		CategoryTypeID: cat.CategoryTypeID.String(),
		LanguageCode:   cat.LanguageCode,
		Name:           cat.Name,
		Description:    cat.Description,
		Weight:         cat.Weight,
		PublishStatus:  cat.PublishStatus,
		CreatedAt:      cat.CreatedAt,
		UpdatedAt:      cat.UpdatedAt,
	}

	if cat.CategoryType.ID != uuid.Nil {
		response.CategoryType = &dto.CategoryTypeResponse{
			ID:        cat.CategoryType.ID.String(),
			TypeCode:  cat.CategoryType.TypeCode,
			Name:      &cat.CategoryType.Name,
			IsActive:  cat.CategoryType.IsActive,
			CreatedAt: cat.CategoryType.CreatedAt,
			UpdatedAt: cat.CategoryType.UpdatedAt,
		}
		if cat.CategoryType.Name == "" {
			response.CategoryType.Name = nil
		}
	} else if cat.CategoryTypeID != uuid.Nil {
		// s.categoryTypeRepo is available here because MapCategoryModelToResponse is a method of cmsCategoryService
		ctModel, err := s.categoryTypeRepo.FindByID(cat.CategoryTypeID)
		if err == nil && ctModel != nil {
			response.CategoryType = &dto.CategoryTypeResponse{
				ID:        ctModel.ID.String(),
				TypeCode:  ctModel.TypeCode,
				Name:      &ctModel.Name,
				IsActive:  ctModel.IsActive,
				CreatedAt: ctModel.CreatedAt,
				UpdatedAt: ctModel.UpdatedAt,
			}
			if ctModel.Name == "" {
				response.CategoryType.Name = nil
			}
		} else {
			fmt.Printf("Warning (MapCategoryModelToResponse): Fallback fetch for CategoryType ID %s failed for Category %s. Error: %v\n", cat.CategoryTypeID, cat.ID, err)
		}
	}

	return response, nil
}

// CreateCategory, GetCategoryByUUID, ListAllCategories, UpdateCategoryByUUID, DeleteCategoryByUUID
// will now call s.MapCategoryModelToResponse(...) instead of a separate helper
func (s *cmsCategoryService) CreateCategory(req dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	// ... (logic สร้าง catModel) ...
	categoryTypeUUID, err := uuid.Parse(req.CategoryTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid category_type_id format: %s", req.CategoryTypeID)
	}
	categoryType, err := s.categoryTypeRepo.FindByID(categoryTypeUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category_type_id '%s' not found", req.CategoryTypeID)
		}
		return nil, fmt.Errorf("error fetching category type: %w", err)
	}
	if !categoryType.IsActive {
		return nil, fmt.Errorf("category_type with ID '%s' is not active", req.CategoryTypeID)
	}

	catModel := &models.Category{
		CategoryTypeID: categoryTypeUUID,
		LanguageCode:   req.LanguageCode,
		Name:           req.Name,
		Description:    req.Description,
		PublishStatus:  req.PublishStatus,
	}
	if req.Weight != nil {
		catModel.Weight = *req.Weight
	} else {
		catModel.Weight = 0
	}

	createdCat, err := s.categoryRepo.CreateCategory(catModel)
	if err != nil {
		// ... (error handling) ...
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("a category with name '%s' for language '%s' under this type already exists", req.Name, req.LanguageCode)
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return s.MapCategoryModelToResponse(createdCat) // เรียกใช้ method ของตัวเอง
}

func (s *cmsCategoryService) GetCategoryByUUID(uuidStr string) (*dto.CategoryResponse, error) {
	categoryID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, errors.New("invalid category UUID format")
	}
	cat, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category by UUID: %w", err)
	}
	return s.MapCategoryModelToResponse(cat)
}

func (s *cmsCategoryService) ListAllCategories(filter dto.CategoryFilter) ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.ListCategoriesByFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories from repository: %w", err)
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))
	for _, cat := range categories {
		resp, mapErr := s.MapCategoryModelToResponse(&cat)
		if mapErr != nil {
			fmt.Printf("Warning: Error mapping category %s: %v\n", cat.ID, mapErr)
			continue
		}
		if resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses, nil
}

func (s *cmsCategoryService) UpdateCategoryByUUID(uuidStr string, req dto.CategoryUpdateRequest) (*dto.CategoryResponse, error) {
	categoryID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category UUID format for update: %s", uuidStr)
	}

	existingCatModel, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category with ID %s not found for update", uuidStr)
		}
		return nil, fmt.Errorf("failed to fetch category %s for update: %w", uuidStr, err)
	}

	if req.Name != nil {
		existingCatModel.Name = *req.Name
	}
	if req.Description != nil {
		existingCatModel.Description = req.Description
	}
	if req.Weight != nil {
		existingCatModel.Weight = *req.Weight
	}
	if req.PublishStatus != nil {
		existingCatModel.PublishStatus = *req.PublishStatus
	}

	updatedModel, err := s.categoryRepo.UpdateCategory(existingCatModel)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, fmt.Errorf("a category with the new name for language '%s' under this type might already exist", existingCatModel.LanguageCode)
		}
		return nil, fmt.Errorf("failed to update category in repository: %w", err)
	}

	return s.MapCategoryModelToResponse(updatedModel)
}

func (s *cmsCategoryService) DeleteCategoryByUUID(uuidStr string) error {
	categoryID, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid category UUID format for delete: %s", uuidStr)
	}
	_, getErr := s.categoryRepo.GetCategoryByID(categoryID)
	if getErr != nil {
		if errors.Is(getErr, gorm.ErrRecordNotFound) {
			return fmt.Errorf("category with ID %s not found for delete", uuidStr)
		}
		return fmt.Errorf("failed to fetch category %s before delete: %w", uuidStr, getErr)
	}

	err = s.categoryRepo.DeleteCategory(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("category with ID %s not found during delete operation by repository", uuidStr)
		}
		return fmt.Errorf("failed to delete category %s from repository: %w", uuidStr, err)
	}
	return nil
}
