package repositories

import (
	"errors"
	"time"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppFaqPageRepositoryInterface interface {
	GetFaqPageBySlug(slug string, preloads []string, isAlias bool, language string) (*models.FaqPage, error)
	GetFaqContentPreview(id uuid.UUID) (*models.FaqContent, error)	
}

type AppFaqPageRepository struct {
	db *gorm.DB
}

func NewAppFaqPageRepository(db *gorm.DB) *AppFaqPageRepository {
	return &AppFaqPageRepository{db: db}
}

func (r *AppFaqPageRepository) GetFaqPageBySlug(slug string, preloads []string, isAlias bool, language string) (*models.FaqPage, error) {
	var faqPage models.FaqPage
	query := r.db

	query = query.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
				return db.
						Where("faq_contents.workflow_status = ? AND faq_contents.language = ? AND faq_contents.mode != ?", enums.WorkflowPublished, language, "Histories").
						Order("faq_contents.created_at DESC")
		})		

	if len(preloads) == 0 {
		query = query.
							Preload("Contents.Revision").
							Preload("Contents.Categories").
							Preload("Contents.Components").
							Preload("Contents.MetaTag")
	} else {
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	// Correctly query using joined faq_contents
	if isAlias {
		query = query.
			Joins("JOIN faq_contents ON faq_contents.page_id = faq_pages.id").
			Where("faq_contents.url_alias = ?", slug)
	} else {
		query = query.
			Joins("JOIN faq_contents ON faq_contents.page_id = faq_pages.id").
			Where("faq_contents.url = ?", slug)
	}

	result := query.First(&faqPage)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errs.ErrNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &faqPage, nil
}

func (r *AppFaqPageRepository) GetFaqContentPreview(id uuid.UUID) (*models.FaqContent, error) {
	var faqContent models.FaqContent
	// Exclude expired_at = null
	err := r.db.
				Preload("MetaTag").
				Preload("Components").
				Where("id = ? AND mode = ? AND (expired_at > ?)", id, "Preview", time.Now()).
				First(&faqContent).Error

	if err != nil {
		return nil, err
	}

	return &faqContent, nil
}