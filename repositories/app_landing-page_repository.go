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

type AppLandingPageRepositoryInterface interface {
	GetLandingPageByUrlAlias(urlAlias string, preloads []string, language string) (*models.LandingPage, error)
	GetLandingContentPreview(id uuid.UUID) (*models.LandingContent, error)
}

type AppLandingPageRepository struct {
	db *gorm.DB
}

func NewAppLandingPageRepository(db *gorm.DB) *AppLandingPageRepository {
	return &AppLandingPageRepository{db: db}
}

func (r *AppLandingPageRepository) GetLandingPageByUrlAlias(urlAlias string, preloads []string, language string) (*models.LandingPage, error) {
	var landingPage models.LandingPage
	query := r.db

	query = query.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
				return db.
						Where("landing_contents.workflow_status = ? AND landing_contents.language = ? AND landing_contents.mode != ?", enums.WorkflowPublished, language, "Histories").
						Order("landing_contents.created_at DESC")
		})

	if len(preloads) == 0 {
		query = query.
							Preload("Contents.Files").
							Preload("Contents.Revision").
							Preload("Contents.Categories").
							Preload("Contents.Components").
							Preload("Contents.MetaTag")
	} else {
		for _, preload := range preloads {
			query = query.Preload(preload)
		}
	}

	// find landing page by url_alias
	query = query.
		Joins("JOIN landing_contents ON landing_contents.page_id = landing_pages.id").
		Where("landing_contents.url_alias = ?", urlAlias)	
	result := query.First(&landingPage)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errs.ErrNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}
	
	return &landingPage, nil
}

func (r *AppLandingPageRepository) GetLandingContentPreview(id uuid.UUID) (*models.LandingContent, error) {
	var landingContent models.LandingContent
	// Exclude expired_at = null
	err := r.db.
				Preload("Files").
				Preload("MetaTag").
				Preload("Components").	
				Where("id = ? AND mode = ? AND (expired_at > ?)", id, "Preview", time.Now()).
				First(&landingContent).Error

	if err != nil {
		return nil, err
	}

	return &landingContent, nil
}