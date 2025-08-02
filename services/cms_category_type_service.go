package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSCategoryTypeServiceInterface interface {
	CreateCategoryType(req dto.CreateCategoryTypeRequest) (*dto.CategoryTypeResponse, error)
	GetCategoryTypeByID(idStr string) (*dto.CategoryTypeResponse, error)
	GetCategoryTypeByCode(code string) (*dto.CategoryTypeResponse, error)
	ListCategoryTypes(isActive *bool) ([]dto.CategoryTypeResponse, error)
	UpdateCategoryType(idStr string, req dto.UpdateCategoryTypeRequest) (*dto.CategoryTypeResponse, error)
	DeleteCategoryType(idStr string) error
	GetCategoryTypeWithDetails(categoryTypeIDStr string, languageCodeStr string) (*dto.CategoryTypeWithDetailsResponse, error)
}

type cmsCategoryTypeService struct {
	categoryTypeRepo repositories.CMSCategoryTypeRepositoryInterface
	categoryRepo     repositories.CMSCategoryRepositoryInterface // สำหรับ ChildrenCount
	categoryService  CMSCategoryServiceInterface                 // << เพิ่ม field นี้
}

// NewCMSCategoryTypeService constructor
func NewCMSCategoryTypeService(
	categoryTypeRepo repositories.CMSCategoryTypeRepositoryInterface,
	categoryRepo repositories.CMSCategoryRepositoryInterface,
	categoryService CMSCategoryServiceInterface, // << เพิ่ม parameter นี้
) CMSCategoryTypeServiceInterface {
	return &cmsCategoryTypeService{
		categoryTypeRepo: categoryTypeRepo,
		categoryRepo:     categoryRepo,
		categoryService:  categoryService, // << กำหนดค่า
	}
}

func (s *cmsCategoryTypeService) mapModelToCategoryTypeResponse(ct *models.CategoryType, childrenCountForThisType map[string]int) *dto.CategoryTypeResponse {
	if ct == nil {
		return nil
	}
	idStr := ""
	if ct.ID != uuid.Nil {
		idStr = ct.ID.String()
	}

	// Ensure childrenCountForThisType is not nil if we want to always show it
	finalChildrenCount := childrenCountForThisType
	if finalChildrenCount == nil {
		finalChildrenCount = make(map[string]int) // Default to empty map
		// Optionally ensure all languages have a count, even if 0
		finalChildrenCount[string(enums.PageLanguageTH)] = 0
		finalChildrenCount[string(enums.PageLanguageEN)] = 0
	}

	namePtr := &ct.Name
	if ct.Name == "" {
		namePtr = nil
	}

	return &dto.CategoryTypeResponse{
		ID:            idStr,
		TypeCode:      ct.TypeCode,
		Name:          namePtr,
		IsActive:      ct.IsActive,
		ChildrenCount: &finalChildrenCount, // Return pointer to the map
		CreatedAt:     ct.CreatedAt,
		UpdatedAt:     ct.UpdatedAt,
	}
}

func (s *cmsCategoryTypeService) CreateCategoryType(req dto.CreateCategoryTypeRequest) (*dto.CategoryTypeResponse, error) {
	unique, err := s.categoryTypeRepo.IsTypeCodeUnique(req.TypeCode, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("error checking type_code uniqueness: %w", err)
	}
	if !unique {
		return nil, errors.New("type_code already exists")
	}

	model := models.CategoryType{
		TypeCode: req.TypeCode,
	}
	if req.Name != nil {
		model.Name = *req.Name
	}
	if req.IsActive != nil {
		model.IsActive = *req.IsActive
	} else {
		model.IsActive = true // Default
	}

	created, err := s.categoryTypeRepo.Create(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to create category type: %w", err)
	}

	// For a new category type, children count is initially 0 for all languages
	initialChildrenCount := make(map[string]int)
	initialChildrenCount[string(enums.PageLanguageTH)] = 0
	initialChildrenCount[string(enums.PageLanguageEN)] = 0
	return s.mapModelToCategoryTypeResponse(created, initialChildrenCount), nil
}

func (s *cmsCategoryTypeService) GetCategoryTypeByID(idStr string) (*dto.CategoryTypeResponse, error) {
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}
	ct, err := s.categoryTypeRepo.FindByID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category type not found")
		}
		return nil, fmt.Errorf("failed to get category type by id: %w", err)
	}

	childrenCount, countErr := s.categoryRepo.CountCategoriesByTypeAndLanguage(ct.ID)
	if countErr != nil {
		fmt.Printf("Warning: failed to count categories for type %s: %v. Counts will be empty.\n", uid, countErr)
		childrenCount = make(map[string]int)
		childrenCount[string(enums.PageLanguageTH)] = 0
		childrenCount[string(enums.PageLanguageEN)] = 0
	}

	return s.mapModelToCategoryTypeResponse(ct, childrenCount), nil
}

func (s *cmsCategoryTypeService) GetCategoryTypeByCode(code string) (*dto.CategoryTypeResponse, error) {
	ct, err := s.categoryTypeRepo.FindByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category type not found for code: %s", code)
		}
		return nil, fmt.Errorf("failed to get category type by code: %w", err)
	}

	childrenCount, countErr := s.categoryRepo.CountCategoriesByTypeAndLanguage(ct.ID)
	if countErr != nil {
		fmt.Printf("Warning: failed to count categories for type code %s (ID: %s): %v. Counts will be empty.\n", code, ct.ID, countErr)
		childrenCount = make(map[string]int)
		childrenCount[string(enums.PageLanguageTH)] = 0
		childrenCount[string(enums.PageLanguageEN)] = 0
	}
	return s.mapModelToCategoryTypeResponse(ct, childrenCount), nil
}

func (s *cmsCategoryTypeService) ListCategoryTypes(isActive *bool) ([]dto.CategoryTypeResponse, error) {
	categoryTypes, err := s.categoryTypeRepo.FindAll(isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to list category types: %w", err)
	}

	if len(categoryTypes) == 0 {
		return []dto.CategoryTypeResponse{}, nil
	}

	responses := make([]dto.CategoryTypeResponse, 0, len(categoryTypes))
	for _, ct := range categoryTypes {
		childrenCount, countErr := s.categoryRepo.CountCategoriesByTypeAndLanguage(ct.ID)
		if countErr != nil {
			fmt.Printf("Warning: failed to count categories for type %s: %v. Counts will be set to empty/zero.\n", ct.ID, countErr)
			childrenCount = make(map[string]int)
			childrenCount[string(enums.PageLanguageTH)] = 0
			childrenCount[string(enums.PageLanguageEN)] = 0
		}
		mapped := s.mapModelToCategoryTypeResponse(&ct, childrenCount)
		if mapped != nil {
			responses = append(responses, *mapped)
		}
	}
	return responses, nil
}

func (s *cmsCategoryTypeService) UpdateCategoryType(idStr string, req dto.UpdateCategoryTypeRequest) (*dto.CategoryTypeResponse, error) {
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid ID format for update")
	}

	_, err = s.categoryTypeRepo.FindByID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category type not found for update")
		}
		return nil, fmt.Errorf("error finding category type for update: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		return s.GetCategoryTypeByID(idStr)
	}

	updated, err := s.categoryTypeRepo.Update(uid, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update category type: %w", err)
	}

	childrenCount, countErr := s.categoryRepo.CountCategoriesByTypeAndLanguage(updated.ID)
	if countErr != nil {
		fmt.Printf("Warning: failed to count categories for updated type %s: %v. Counts will be empty.\n", uid, countErr)
		childrenCount = make(map[string]int)
		childrenCount[string(enums.PageLanguageTH)] = 0
		childrenCount[string(enums.PageLanguageEN)] = 0
	}

	return s.mapModelToCategoryTypeResponse(updated, childrenCount), nil
}

func (s *cmsCategoryTypeService) DeleteCategoryType(idStr string) error {
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("invalid ID format for delete")
	}

	childrenCountMap, countErr := s.categoryRepo.CountCategoriesByTypeAndLanguage(uid)
	if countErr != nil {
		fmt.Printf("Warning: Could not verify category usage for type %s during deletion: %v\n", uid, countErr)
	}

	totalCategories := 0
	for _, count := range childrenCountMap {
		totalCategories += count
	}

	if totalCategories > 0 {
		return fmt.Errorf("cannot delete category type: it is still in use by %d categories (details)", totalCategories)
	}

	err = s.categoryTypeRepo.Delete(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category type not found for delete")
		}
		return fmt.Errorf("failed to delete category type: %w", err)
	}
	return nil
}

func (s *cmsCategoryTypeService) GetCategoryTypeWithDetails(categoryTypeIDStr string, languageCodeStr string) (*dto.CategoryTypeWithDetailsResponse, error) {
	categoryTypeUUID, err := uuid.Parse(categoryTypeIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category_type_id format: %s", categoryTypeIDStr)
	}

	categoryType, err := s.categoryTypeRepo.FindByID(categoryTypeUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category_type with ID '%s' not found", categoryTypeIDStr)
		}
		return nil, fmt.Errorf("error fetching category_type: %w", err)
	}

	var langEnum enums.PageLanguage
	switch strings.ToLower(languageCodeStr) {
	case string(enums.PageLanguageTH):
		langEnum = enums.PageLanguageTH
	case string(enums.PageLanguageEN):
		langEnum = enums.PageLanguageEN
	default:
		return nil, fmt.Errorf("invalid language code: %s. Supported: th, en", languageCodeStr)
	}

	tempLang := langEnum
	filter := dto.CategoryFilter{
		CategoryTypeID: &categoryTypeIDStr,
		LanguageCode:   &tempLang,
	}

	categoryDetailResponses, err := s.categoryService.ListAllCategories(filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories for type %s and language %s: %w", categoryTypeIDStr, languageCodeStr, err)
	}

	response := &dto.CategoryTypeWithDetailsResponse{
		ID:         categoryType.ID.String(),
		TypeCode:   categoryType.TypeCode,
		Name:       &categoryType.Name,
		IsActive:   categoryType.IsActive,
		CreatedAt:  categoryType.CreatedAt,
		UpdatedAt:  categoryType.UpdatedAt,
		Categories: categoryDetailResponses,
	}
	if categoryType.Name == "" {
		response.Name = nil
	}

	return response, nil
}
