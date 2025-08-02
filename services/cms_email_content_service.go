package services

import (
	"errors"
	"fmt"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailContentServiceInterface interface {
	CreateContent(req dto.CreateEmailContentRequest) (*dto.EmailContentResponse, error)
	GetContentByID(idStr string) (*dto.EmailContentResponse, error)
	GetContentByCategoryAndLangAndLabel(categoryIDStr string, language enums.PageLanguage, label string) (*dto.EmailContentResponse, error)
	ListContents(filters dto.EmailContentFilter) ([]dto.EmailContentResponse, error)
	UpdateContent(idStr string, req dto.UpdateEmailContentRequest) (*dto.EmailContentResponse, error)
	DeleteContent(idStr string) error
	GetEmailContentByCategoryIDAndLanguage(categoryIDStr string, language enums.PageLanguage) ([]*dto.EmailContentResponse, error)
}

type emailContentService struct {
	repo         repositories.EmailContentRepositoryInterface
	categoryRepo repositories.EmailCategoryRepositoryInterface
}

func NewEmailContentService(repo repositories.EmailContentRepositoryInterface, categoryRepo repositories.EmailCategoryRepositoryInterface) EmailContentServiceInterface {
	return &emailContentService{repo: repo, categoryRepo: categoryRepo}
}

func mapEmailContentToResponse(content *models.EmailContent) *dto.EmailContentResponse {
	fmt.Println("mapEmailContentToResponse", content)
	if content == nil {
		return nil
	}
	resp := &dto.EmailContentResponse{
		ID:              content.ID.String(),
		EmailCategoryID: content.EmailCategoryID.String(),
		Language:        content.Language,
		Label:           content.Label,

		SendTo:          content.SendTo,
		CcEmail:         content.CcEmail,
		BccEmail:        content.BccEmail,
		SendFromEmail:   content.SendFromEmail,
		SendFromName:    content.SendFromName,
		Subject:         content.Subject,
		TopImgLink:      content.TopImgLink,
		Header:          content.Header,
		Paragraph:       content.Paragraph,
		Footer:          content.Footer,
		FooterImageLink: content.FooterImageLink,

		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}

	if content.EmailCategory != nil && content.EmailCategory.ID != uuid.Nil { // Check if EmailCategory was preloaded
		fmt.Println("check ???")
		resp.EmailCategory = mapEmailCategoryToResponse(content.EmailCategory)
		fmt.Println("email category inside", resp.EmailCategory)
	}
	return resp
}

func (s *emailContentService) CreateContent(req dto.CreateEmailContentRequest) (*dto.EmailContentResponse, error) {
	categoryID, err := uuid.Parse(req.EmailCategoryID)
	if err != nil {
		return nil, errors.New("invalid email_category_id format")
	}

	_, err = s.categoryRepo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email category not found")
		}
		return nil, fmt.Errorf("failed to verify email category: %w", err)
	}

	// Check for uniqueness of (categoryID, language, label)
	existing, _ := s.repo.FindByCategoryIDAndLanguageAndLabel(categoryID, req.Language, req.Label)
	if existing != nil {
		return nil, errors.New("email content with this category, language, and label already exists")
	}

	contentModel := &models.EmailContent{
		EmailCategoryID: categoryID,
		Language:        req.Language,
		Label:           req.Label,
		SendTo:          req.SendTo,
		CcEmail:         req.CcEmail,
		BccEmail:        req.BccEmail,
		SendFromEmail:   req.SendFromEmail,
		SendFromName:    req.SendFromName,
		Subject:         req.Subject,
		TopImgLink:      req.TopImgLink,
		Header:          req.Header,
		Paragraph:       req.Paragraph,
		Footer:          req.Footer,
		FooterImageLink: req.FooterImageLink,
	}

	helpers.SanitizeEmailContent(contentModel)

	createdContent, err := s.repo.Create(contentModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create email content: %w", err)
	}
	return mapEmailContentToResponse(createdContent), nil
}

func (s *emailContentService) GetContentByID(idStr string) (*dto.EmailContentResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid content ID format")
	}
	content, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email content not found")
		}
		return nil, fmt.Errorf("failed to get email content: %w", err)
	}
	return mapEmailContentToResponse(content), nil
}

func (s *emailContentService) GetContentByCategoryAndLangAndLabel(categoryIDStr string, language enums.PageLanguage, label string) (*dto.EmailContentResponse, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}
	content, err := s.repo.FindByCategoryIDAndLanguageAndLabel(categoryID, language, label)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email content not found for the given category, language, and label")
		}
		return nil, fmt.Errorf("failed to get email content: %w", err)
	}
	return mapEmailContentToResponse(content), nil
}

func (s *emailContentService) ListContents(filters dto.EmailContentFilter) ([]dto.EmailContentResponse, error) {
	contents, err := s.repo.ListByFilters(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list email contents: %w", err)
	}
	responses := make([]dto.EmailContentResponse, len(contents))
	for i, content := range contents {
		responses[i] = *mapEmailContentToResponse(&content)
	}
	return responses, nil
}

func (s *emailContentService) UpdateContent(idStr string, req dto.UpdateEmailContentRequest) (*dto.EmailContentResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid content ID format")
	}

	content, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email content not found for update")
		}
		return nil, fmt.Errorf("failed to find email content for update: %w", err)
	}

	// Check for uniqueness if label or language is changing
	currentLang := content.Language
	if req.Language != nil {
		currentLang = *req.Language
	}
	currentLabel := content.Label
	if req.Label != nil {
		currentLabel = *req.Label
	}

	if (req.Language != nil && *req.Language != content.Language) || (req.Label != nil && *req.Label != content.Label) {
		existing, _ := s.repo.FindByCategoryIDAndLanguageAndLabel(content.EmailCategoryID, currentLang, currentLabel)
		if existing != nil && existing.ID != id { // If it's a different record
			return nil, errors.New("email content with this category, language, and label already exists")
		}
	}

	// Apply updates
	if req.Language != nil {
		content.Language = *req.Language
	}
	if req.Label != nil {
		content.Label = *req.Label
	}
	if req.SendTo != nil {
		content.SendTo = *req.SendTo
	}
	if req.CcEmail != nil {
		content.CcEmail = *req.CcEmail
	}
	if req.BccEmail != nil {
		content.BccEmail = *req.BccEmail
	}
	if req.SendFromEmail != nil {
		content.SendFromEmail = *req.SendFromEmail
	}
	if req.SendFromName != nil {
		content.SendFromName = *req.SendFromName
	}
	if req.Subject != nil {
		content.Subject = *req.Subject
	}
	if req.TopImgLink != nil {
		content.TopImgLink = *req.TopImgLink
	}
	if req.Header != nil {
		content.Header = *req.Header
	}
	if req.Paragraph != nil {
		content.Paragraph = *req.Paragraph
	}
	if req.Footer != nil {
		content.Footer = *req.Footer
	}
	if req.FooterImageLink != nil {
		content.FooterImageLink = *req.FooterImageLink
	}

	helpers.SanitizeEmailContent(content)

	updatedContent, err := s.repo.Update(content)
	if err != nil {
		return nil, fmt.Errorf("failed to update email content: %w", err)
	}
	return mapEmailContentToResponse(updatedContent), nil
}

func (s *emailContentService) DeleteContent(idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("invalid content ID format")
	}
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("email content not found for delete")
		}
		return fmt.Errorf("failed to delete email content: %w", err)
	}
	return nil
}

func (s *emailContentService) GetEmailContentByCategoryIDAndLanguage(categoryIDStr string, language enums.PageLanguage) ([]*dto.EmailContentResponse, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}

	contents, err := s.repo.FindEmailContentByCategoryIDAndLanguage(categoryID, language)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no email content found for the given category and language")
		}
		return nil, fmt.Errorf("failed to get email content: %w", err)
	}

	fmt.Println("Found contents:", len(contents))

	responses := make([]*dto.EmailContentResponse, len(contents))
	for i, content := range contents {
		responses[i] = mapEmailContentToResponse(&content)

		fmt.Println("mapped อันแรกเสร็จ")
	}
	fmt.Println("Mapped responses:", len(responses))
	return responses, nil
}
