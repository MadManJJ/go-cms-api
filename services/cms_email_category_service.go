package services

import (
	"errors"
	"fmt"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailCategoryServiceInterface interface {
	CreateCategory(req dto.CreateEmailCategoryRequest) (*dto.EmailCategoryResponse, error)
	GetCategoryByID(idStr string) (*dto.EmailCategoryResponse, error)
	GetCategoryByTitle(title string) (*dto.EmailCategoryResponse, error)
	ListCategories() ([]dto.EmailCategoryResponse, error)
	UpdateCategory(idStr string, req dto.UpdateEmailCategoryRequest) (*dto.EmailCategoryResponse, error)
	DeleteCategory(idStr string) error
}

type emailCategoryService struct {
	repo        repositories.EmailCategoryRepositoryInterface
	contentRepo repositories.EmailContentRepositoryInterface // To delete content when category is deleted
}

func NewEmailCategoryService(repo repositories.EmailCategoryRepositoryInterface, contentRepo repositories.EmailContentRepositoryInterface) EmailCategoryServiceInterface {
	return &emailCategoryService{repo: repo, contentRepo: contentRepo}
}

func mapEmailCategoryToResponse(category *models.EmailCategory) *dto.EmailCategoryResponse {
	fmt.Println("mapEmailCategoryToResponse", category)
	if category == nil {
		return nil
	}
	return &dto.EmailCategoryResponse{
		ID:        category.ID.String(),
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func (s *emailCategoryService) CreateCategory(req dto.CreateEmailCategoryRequest) (*dto.EmailCategoryResponse, error) {
	isUnique, err := s.repo.IsTitleUnique(req.Title, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check title uniqueness: %w", err)
	}
	if !isUnique {
		return nil, errors.New("email category title already exists")
	}

	category := &models.EmailCategory{
		Title: req.Title,
	}

	helpers.SanitizeEmailCategory(category)

	createdCategory, err := s.repo.Create(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create email category: %w", err)
	}
	return mapEmailCategoryToResponse(createdCategory), nil
}

func (s *emailCategoryService) GetCategoryByID(idStr string) (*dto.EmailCategoryResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}
	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email category not found")
		}
		return nil, fmt.Errorf("failed to get email category: %w", err)
	}
	return mapEmailCategoryToResponse(category), nil
}

func (s *emailCategoryService) GetCategoryByTitle(title string) (*dto.EmailCategoryResponse, error) {
	category, err := s.repo.FindByTitle(title)
	if err != nil {
		return nil, fmt.Errorf("failed to get email category by title: %w", err)
	}
	if category == nil { // FindByTitle returns nil, nil if not found
		return nil, errors.New("email category not found")
	}
	return mapEmailCategoryToResponse(category), nil
}

func (s *emailCategoryService) ListCategories() ([]dto.EmailCategoryResponse, error) {
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list email categories: %w", err)
	}
	responses := make([]dto.EmailCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *mapEmailCategoryToResponse(&category)
	}
	return responses, nil
}

func (s *emailCategoryService) UpdateCategory(idStr string, req dto.UpdateEmailCategoryRequest) (*dto.EmailCategoryResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}

	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email category not found for update")
		}
		return nil, fmt.Errorf("failed to find email category for update: %w", err)
	}

	if req.Title != "" && req.Title != category.Title {
		isUnique, err := s.repo.IsTitleUnique(req.Title, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check title uniqueness: %w", err)
		}
		if !isUnique {
			return nil, errors.New("email category title already exists")
		}
		category.Title = req.Title
	}

	helpers.SanitizeEmailCategory(category)

	updatedCategory, err := s.repo.Update(category)
	if err != nil {
		return nil, fmt.Errorf("failed to update email category: %w", err)
	}
	return mapEmailCategoryToResponse(updatedCategory), nil
}

func (s *emailCategoryService) DeleteCategory(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("invalid category ID format")
	}
	// First, delete all associated email contents
	if err := s.contentRepo.DeleteByCategoryID(id); err != nil {
		// If it's a record not found error for contents, it's okay, proceed to delete category
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to delete associated email contents: %w", err)
		}
	}
	// Then, delete the category itself
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("email category not found for delete")
		}
		return fmt.Errorf("failed to delete email category: %w", err)
	}
	return nil
}
