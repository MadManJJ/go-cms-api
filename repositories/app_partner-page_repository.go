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

type AppPartnerPageRepositoryInterface interface {
	GetPartnerPageBySlug(slug string, preloads []string, isAlias bool, language string) (*models.PartnerPage, error)
	GetPartnerContentPreview(id uuid.UUID) (*models.PartnerContent, error)
}

type AppPartnerPageRepository struct {
	db *gorm.DB
}

func NewAppPartnerPageRepository(db *gorm.DB) *AppPartnerPageRepository {
	return &AppPartnerPageRepository{db: db}
}

func (r *AppPartnerPageRepository) GetPartnerPageBySlug(slug string, preloads []string, isAlias bool, language string) (*models.PartnerPage, error) {
	var partnerPage models.PartnerPage
	query := r.db

	query = query.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
				return db.
						Where("partner_contents.workflow_status = ? AND partner_contents.language = ? AND partner_contents.mode != ?", enums.WorkflowPublished, language, "Histories").
						Order("partner_contents.created_at DESC")
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
	
	// can query for both url_alias and url
	if isAlias {
		query = query.
			Joins("JOIN partner_contents ON partner_contents.page_id = partner_pages.id").
			Where("partner_contents.url_alias = ?", slug)
	} else {
		query = query.
			Joins("JOIN partner_contents ON partner_contents.page_id = partner_pages.id").
			Where("partner_contents.url = ?", slug)
	}
	result := query.First(&partnerPage)
	
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errs.ErrNotFound
	}	

	if result.Error != nil {
		return nil, result.Error
	}
	
	return &partnerPage, nil	
}

func (r *AppPartnerPageRepository) GetPartnerContentPreview(id uuid.UUID) (*models.PartnerContent, error) {
	var partnerContent models.PartnerContent
	// Exclude expired_at = null
	err := r.db.
				Preload("MetaTag").
				Preload("Components").	
				Where("id = ? AND mode = ? AND (expired_at > ?)", id, "Preview", time.Now()).
				First(&partnerContent).Error

	if err != nil {
		return nil, err
	}

	return &partnerContent, nil
}