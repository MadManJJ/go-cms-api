package repositories

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSLandingPageRepositoryInterface interface {
	CreateLandingPage(LandingPage *models.LandingPage) (*models.LandingPage, error)
	FindAllLandingPage(query dto.LandingPageQuery, sort string, page, limit int, language string) ([]models.LandingPage, int64, error)
	FindLandingPageById(id uuid.UUID) (*models.LandingPage, error)
	UpdateLandingContent(updateLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error)
	DeleteLandingPage(id uuid.UUID) error
	FindContentByLandingPageId(pageId uuid.UUID, language string, mode string) (*models.LandingContent, error)
	FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.LandingContent, error)
	// Deprecate
	CreateContentForLandingPage(LandingContent *models.LandingContent, lang string, mode string) (*models.LandingContent, error)
	DeleteLandingContent(pageId uuid.UUID, lang, mode string) error
	DuplicateLandingPage(pageId uuid.UUID) (*models.LandingPage, error)
	DuplicateLandingContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error)
	RevertLandingContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error)
	GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	GetRevisionByLandingPageId(pageId uuid.UUID, language string) ([]models.Revision, error)
	IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error)
	GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error)
	CreateLandingContentPreview(landingContentPreview *models.LandingContent) (*models.LandingContent, error)
	UpdateLandingContentPreview(landingContentPreview *models.LandingContent) (*models.LandingContent, error)
	FindLandingContentPreviewById(pageId uuid.UUID, language string) (*models.LandingContent, error)
}

type CMSLandingPageRepository struct {
	db *gorm.DB
}

func NewCMSLandingPageRepository(db *gorm.DB) *CMSLandingPageRepository {
	return &CMSLandingPageRepository{db: db}
}

func (r *CMSLandingPageRepository) CreateLandingPage(LandingPage *models.LandingPage) (*models.LandingPage, error) {
	//Create the LandingPage
	if err := r.db.Create(LandingPage).Error; err != nil {
		return nil, err
	}

	return LandingPage, nil
}

// Original postgres
// SELECT Landing_pages.* FROM Landing_pages
// JOIN Landing_contents ON Landing_contents.page_id = Landing_pages.id
// LEFT JOIN Landing_content_categories ON Landing_content_categories.Landing_content_id = Landing_contents.id
// JOIN categories ON Landing_content_categories.category_id = categories.id
// JOIN category_types ON category_types.id = categories.category_type_id
// WHERE category_types.type_code = 'Landing' AND categories.name ILIKE '%filter.CategoryLanding%'
// INTERSECT
// SELECT Landing_pages.* FROM Landing_pages
// JOIN Landing_contents ON Landing_contents.page_id = Landing_pages.id
// LEFT JOIN Landing_content_categories ON Landing_content_categories.Landing_content_id = Landing_contents.id
// JOIN categories ON Landing_content_categories.category_id = categories.id
// JOIN category_types ON category_types.id = categories.category_type_id
// WHERE category_types.type_code = 'category-keywords' AND categories.name ILIKE '%filter.CategoryLanding%'
func (r *CMSLandingPageRepository) FindAllLandingPage(query dto.LandingPageQuery, sort string, page, limit int, language string) ([]models.LandingPage, int64, error) {
	var landingPages []models.LandingPage
	var totalCount int64

	// Build base query with proper joins and filters
	baseQuery := r.db.Model(&models.LandingPage{}).
		Joins("JOIN landing_contents ON landing_contents.page_id = landing_pages.id").
		Where("landing_contents.mode != ? AND landing_contents.mode != ?", "Histories", "Preview")

	// Apply filters
	if query.Title != "" {
		baseQuery = baseQuery.Where("landing_contents.title ILIKE ?", "%"+query.Title+"%")
	}
	if query.UrlAlias != "" {
		baseQuery = baseQuery.Where("landing_contents.url_alias ILIKE ?", "%"+query.UrlAlias+"%")
	}
	if query.Status != "" {
		baseQuery = baseQuery.Where("landing_contents.workflow_status = ?", query.Status)
	}
	if language != "" {
		baseQuery = baseQuery.Where("landing_contents.language = ?", language)
	}
	if query.CategoryKeywords != "" {
		categorySubQuery := r.db.Table("landing_content_categories").
			Select("landing_content_categories.landing_content_id").
			Joins("JOIN categories ON landing_content_categories.category_id = categories.id").
			Joins("JOIN category_types ON category_types.id = categories.category_type_id").
			Where("category_types.type_code = ? AND categories.name ILIKE ?", "category-keywords", "%"+query.CategoryKeywords+"%")
		baseQuery = baseQuery.Where("landing_contents.id IN (?)", categorySubQuery)
	}

	// Clone query for counting to avoid modification
	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Model(&models.LandingPage{}).Distinct("landing_pages.id").Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Create a new variable for the final query to avoid modifying the baseQuery
	finalQuery := baseQuery

	// Apply sorting
	sortColumn := "landing_pages.created_at DESC" // default sort
	sortableColumns := map[string]string{
		"title":      "landing_contents.title",
		"url_alias":  "landing_contents.url_alias",
		"status":     "landing_contents.workflow_status",
		"created_at": "landing_pages.created_at",
		"updated_at": "landing_pages.updated_at",
	}

	// Check if sorting is on a joined table column which requires GROUP BY
	isSortOnJoinedTable := false
	if sort != "" {
		sortParts := strings.Split(sort, ":")
		if len(sortParts) == 2 {
			column := sortParts[0]
			direction := strings.ToUpper(sortParts[1])
			if direction != "ASC" && direction != "DESC" {
				direction = "DESC"
			}

			if col, ok := sortableColumns[column]; ok {
				// Identify if the column is from landing_contents
				if strings.HasPrefix(col, "landing_contents.") {
					isSortOnJoinedTable = true
					// Use an aggregate function for sorting to resolve ambiguity
					// MIN() is a safe choice to get a deterministic sort order
					sortColumn = fmt.Sprintf("MIN(%s) %s", col, direction)
				} else {
					sortColumn = fmt.Sprintf("%s %s", col, direction)
				}
			}
		}
	}

	if isSortOnJoinedTable {
		// If sorting is on a joined table, we must GROUP BY the primary table's ID
		// and select the primary table's columns.
		finalQuery = finalQuery.Select("landing_pages.*").Group("landing_pages.id")
	} else {
		// If sorting is on the primary table, DISTINCT is safe and efficient.
		finalQuery = finalQuery.Select("DISTINCT landing_pages.*")
	}

	// Get paginated and sorted landing pages with all necessary preloads
	offset := (page - 1) * limit
	err := finalQuery.
		Order(sortColumn).
		Offset(offset).
		Limit(limit).
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("landing_contents.mode != ? AND landing_contents.mode != ?", "Histories", "Preview").
				Order("landing_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		Preload("Contents.Files").
		Find(&landingPages).Error

	if err != nil {
		return nil, 0, err
	}

	return landingPages, totalCount, nil
}
func (r *CMSLandingPageRepository) FindLandingPageById(id uuid.UUID) (*models.LandingPage, error) {
	var LandingPage models.LandingPage

	err := r.db.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("landing_contents.mode != ? AND landing_contents.mode != ?", "Histories", "Preview").
				Order("landing_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		First(&LandingPage, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &LandingPage, nil
}

func (r *CMSLandingPageRepository) UpdateLandingContent(updateLandingContent *models.LandingContent, prevContentId uuid.UUID) (*models.LandingContent, error) {

	err := r.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("[REPO-BEFORE-NORMALIZE] Language: '%s'", updateLandingContent.Language)

		var oldContent models.LandingContent
		if err := tx.First(&oldContent, "id = ?", prevContentId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("previous content with ID %s not found", prevContentId)
			}
			return err
		}

		oldContent.Mode = enums.PageModeHistories
		if err := tx.Save(&oldContent).Error; err != nil {
			return fmt.Errorf("failed to archive old content: %w", err)
		}

		updateLandingContent.PageID = oldContent.PageID

		if err := helpers.NormalizeLandingContent(updateLandingContent); err != nil {
			return err
		}
		log.Printf("[REPO-AFTER-NORMALIZE] Language: '%s'", updateLandingContent.Language)
		if updateLandingContent.Revision != nil {
			if err := helpers.NormalizeRevision(updateLandingContent.Revision); err != nil {
				return err
			}
		}

		updateLandingContent.ID = uuid.Nil
		updateLandingContent.CreatedAt = time.Time{} // Reset เพื่อให้ DB generate
		updateLandingContent.UpdatedAt = time.Time{}

		if updateLandingContent.MetaTag != nil {
			updateLandingContent.MetaTag.ID = uuid.Nil
			updateLandingContent.MetaTagID = uuid.Nil // Clear FK
		}
		if updateLandingContent.Revision != nil {
			updateLandingContent.Revision.ID = uuid.Nil
			updateLandingContent.Revision.LandingContentID = nil // Clear FK
		}
		for _, component := range updateLandingContent.Components {
			if component != nil {
				component.ID = uuid.Nil
				component.LandingContentID = nil // Clear FK
			}
		}
		for _, file := range updateLandingContent.Files {
			if file != nil {
				file.ID = uuid.Nil
				file.LandingContentID = uuid.Nil // Clear FK
			}
		}

		if err := tx.Create(updateLandingContent).Error; err != nil {
			return fmt.Errorf("failed to create new content version: %w", err)
		}

		if err := tx.Model(&models.LandingPage{}).Where("id = ?", updateLandingContent.PageID).Update("updated_at", time.Now()).Error; err != nil {
			return fmt.Errorf("failed to update page timestamp: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 7. Preload ข้อมูลทั้งหมดกลับมาเพื่อความสมบูรณ์ของ Response
	if err := r.db.
		Preload("MetaTag").
		Preload("Revision").
		Preload("Components").
		Preload("Files").
		Preload("Categories.CategoryType").
		First(&updateLandingContent, "id = ?", updateLandingContent.ID).Error; err != nil {
		log.Printf("Warning: failed to fully preload associations for new content %s: %v", updateLandingContent.ID, err)

	}

	return updateLandingContent, nil
}

func (r *CMSLandingPageRepository) DeleteLandingPage(id uuid.UUID) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		var contents []models.LandingContent

		// Step 0: Ensure the LandingPage exists
		var page models.LandingPage
		if err := tx.First(&page, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.ErrNotFound
			}
			return err
		}

		// Step 1: Get all LandingContent entries for this page
		if err := tx.Where("page_id = ?", id).Find(&contents).Error; err != nil {
			return err
		}

		// Step 2: Collect all LandingContent IDs
		var contentIDs []uuid.UUID
		for _, content := range contents {
			contentIDs = append(contentIDs, content.ID)
		}

		if len(contentIDs) > 0 {
			// Step 3: Delete Components
			if err := tx.Where("Landing_content_id IN ?", contentIDs).Delete(&models.Component{}).Error; err != nil {
				return err
			}

			// Step 4: Delete LandingContentCategories
			if err := tx.Where("Landing_content_id IN ?", contentIDs).Delete(&models.LandingContentCategory{}).Error; err != nil {
				return err
			}

			// Step 5: Delete Revisions
			if err := tx.Where("Landing_content_id IN ?", contentIDs).Delete(&models.Revision{}).Error; err != nil {
				return err
			}

			// Step 6: Delete LandingContentFiles
			if err := tx.Where("Landing_content_id IN ?", contentIDs).Delete(&models.LandingContentFile{}).Error; err != nil {
				return err
			}
		}

		// Step 7: Delete LandingContents
		if err := tx.Where("page_id = ?", id).Delete(&models.LandingContent{}).Error; err != nil {
			return err
		}

		// Step 8: Delete LandingPage
		if err := tx.Where("id = ?", id).Delete(&models.LandingPage{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CMSLandingPageRepository) FindContentByLandingPageId(pageId uuid.UUID, language string, mode string) (*models.LandingContent, error) {
	var LandingContent models.LandingContent

	err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Categories.CategoryType").
		Preload("Components").
		Where("page_id = ? AND language = ? AND mode = ?", pageId, language, mode).
		First(&LandingContent).Error

	if err != nil {
		return nil, err
	}

	return &LandingContent, nil
}

func (r *CMSLandingPageRepository) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.LandingContent, error) {
	var LandingContent models.LandingContent
	if err := r.db.
		Where("page_id = ? AND language = ?", pageId, language).
		Order("created_at DESC").
		First(&LandingContent).Error; err != nil {
		return nil, err
	}

	return &LandingContent, nil
}

// Might be deprecate
func (r *CMSLandingPageRepository) CreateContentForLandingPage(LandingContent *models.LandingContent, lang string, mode string) (*models.LandingContent, error) {
	if err := r.db.Create(LandingContent).Error; err != nil {
		return nil, err
	}

	var createdLandingContent models.LandingContent
	if err := r.db.First(&createdLandingContent, "id = ?", LandingContent.ID).Error; err != nil {
		return nil, err
	}

	return &createdLandingContent, nil
}

func (r *CMSLandingPageRepository) DeleteLandingContent(pageId uuid.UUID, lang, mode string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		var LandingContent models.LandingContent
		if err := tx.First(&LandingContent, "page_id = ? AND language = ? AND mode = ?", pageId, lang, mode).Error; err != nil {
			return err
		}

		if err := tx.Where("Landing_content_id = ?", LandingContent.ID).Delete(&models.Component{}).Error; err != nil {
			return err
		}

		if err := tx.Where("Landing_content_id = ?", LandingContent.ID).Delete(&models.Revision{}).Error; err != nil {
			return err
		}

		if err := tx.Where("Landing_content_id = ?", LandingContent.ID).Delete(&models.LandingContentCategory{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&LandingContent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CMSLandingPageRepository) DuplicateLandingPage(pageId uuid.UUID) (*models.LandingPage, error) {
	// Already a pointer
	landingPage := &models.LandingPage{}

	err := r.db.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("landing_contents.mode != ?", "Histories").
				Order("landing_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		First(&landingPage, "id = ?", pageId).Error

	if err != nil {
		return nil, err
	}

	if len(landingPage.Contents) == 0 {
		return nil, errs.ErrNoNewContentToDuplicate
	}

	// For loop to duplicate two of the main language (en, th)
	landingContents := landingPage.Contents
	newContents := []*models.LandingContent{}
	for _, originalLandingContent := range landingContents {

		// Create a deep copy
		copyLandingContent := *originalLandingContent

		// Metatag
		if originalLandingContent.MetaTag != nil {
			metaTagCopy := *originalLandingContent.MetaTag
			metaTagCopy.ID = uuid.Nil
			metaTagCopy.CreatedAt = time.Time{}
			metaTagCopy.UpdatedAt = time.Time{}
			copyLandingContent.MetaTag = &metaTagCopy
		}

		// Revision
		if originalLandingContent.Revision != nil {
			revisionCopy := *originalLandingContent.Revision
			revisionCopy.ID = uuid.Nil
			revisionCopy.LandingContentID = nil
			revisionCopy.CreatedAt = time.Time{}
			revisionCopy.UpdatedAt = time.Time{}
			copyLandingContent.Revision = &revisionCopy
		}

		// Components
		if originalLandingContent.Components != nil {
			var componentCopies []*models.Component
			for _, component := range originalLandingContent.Components {
				comp := *component
				comp.ID = uuid.Nil
				comp.LandingContentID = nil
				comp.CreatedAt = time.Time{}
				comp.UpdatedAt = time.Time{}
				componentCopies = append(componentCopies, &comp)
			}
			copyLandingContent.Components = componentCopies
		}

		// UrlAlias: Random a string of length 3 and concat it to the back
		copyLandingContent.UrlAlias += "-" + helpers.RandomString(3)

		// Content
		copyLandingContent.ID = uuid.Nil
		copyLandingContent.PageID = uuid.Nil
		copyLandingContent.CreatedAt = time.Time{}
		copyLandingContent.UpdatedAt = time.Time{}

		// Attach the new content
		newContents = append(newContents, &copyLandingContent)
	}

	// Landing Page
	copyLandingPage := *landingPage
	copyLandingPage.ID = uuid.Nil
	copyLandingPage.CreatedAt = time.Time{}
	copyLandingPage.UpdatedAt = time.Time{}

	// Use the new contents
	copyLandingPage.Contents = newContents

	// Create the new landingPage
	if err := r.db.Create(&copyLandingPage).Error; err != nil {
		return nil, err
	}

	return &copyLandingPage, nil
}

func (r *CMSLandingPageRepository) DuplicateLandingContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	var landingContent models.LandingContent
	if err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Components").
		Preload("MetaTag").
		First(&landingContent, "id = ?", contentId).Error; err != nil {
		return nil, err
	}

	// Change language
	if landingContent.Language == enums.PageLanguageTH {
		landingContent.Language = enums.PageLanguageEN
	} else {
		landingContent.Language = enums.PageLanguageTH
	}

	landingContent.ID = uuid.Nil

	// Metatag
	landingContent.MetaTagID = uuid.Nil
	if landingContent.MetaTag != nil {
		landingContent.MetaTag.ID = uuid.Nil
		landingContent.MetaTag.CreatedAt = time.Time{}
		landingContent.MetaTag.UpdatedAt = time.Time{}
	}

	// Content
	landingContent.ID = uuid.Nil
	landingContent.CreatedAt = time.Time{}
	landingContent.UpdatedAt = time.Time{}

	// Attach revision, always needs a revision for keeping the history (but this is a new content, so it's a start of the history)
	if newRevision == nil {
		return nil, errs.ErrNoRevisionFound
	}
	landingContent.Revision = newRevision

	// Components
	if landingContent.Components != nil {
		for _, component := range landingContent.Components {
			component.CreatedAt = time.Time{}
			component.UpdatedAt = time.Time{}
			component.ID = uuid.Nil
			component.LandingContentID = nil
		}
	}

	if err := r.db.Create(&landingContent).Error; err != nil {
		return nil, err
	}

	return &landingContent, nil
}

func (r *CMSLandingPageRepository) RevertLandingContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.LandingContent, error) {
	var revision models.Revision
	// Get the revision
	if err := r.db.First(&revision, "id = ?", revisionId).Error; err != nil {
		return nil, err
	}

	LandingContent := &models.LandingContent{}
	// Get the content of that revision
	if err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Components").
		Preload("MetaTag").
		First(LandingContent, "id = ?", revision.LandingContentID).Error; err != nil {
		return nil, err
	}

	// Get the old content
	oldContent, err := r.FindLatestContentByPageId(LandingContent.PageID, string(LandingContent.Language))
	if err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Update the old content

		oldContent.Mode = enums.PageModeHistories
		if err := tx.Save(oldContent).Error; err != nil {
			return err
		}

		LandingContent.Mode = enums.PageModeDraft
		LandingContent.ID = uuid.Nil

		// Metatag
		LandingContent.MetaTagID = uuid.Nil
		if LandingContent.MetaTag != nil {
			LandingContent.MetaTag.ID = uuid.Nil
			LandingContent.MetaTag.CreatedAt = time.Time{}
			LandingContent.MetaTag.UpdatedAt = time.Time{}
		}

		// Content
		LandingContent.ID = uuid.Nil
		LandingContent.CreatedAt = time.Time{}
		LandingContent.UpdatedAt = time.Time{}

		// Attach revision, always needs a revision for keeping the history (but this is a new content, so it's a start of the history)
		if newRevision == nil {
			return errs.ErrNoRevisionFound
		}
		LandingContent.Revision = newRevision

		// Components
		if LandingContent.Components != nil {
			for _, component := range LandingContent.Components {
				component.CreatedAt = time.Time{}
				component.UpdatedAt = time.Time{}
				component.ID = uuid.Nil
				component.LandingContentID = nil
			}
		}

		if err := tx.Create(LandingContent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return LandingContent, nil
}

func (r *CMSLandingPageRepository) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	var result struct {
		ID uuid.UUID
	}
	if err := r.db.
		Model(&models.LandingContent{}).
		Select("id").
		Where("page_id = ? AND language = ? AND mode = ?", pageId, language, mode).
		Take(&result).Error; err != nil {
		return nil, err
	}

	LandingContentID := result.ID

	var categoryIDs []uuid.UUID
	if err := r.db.
		Model(&models.LandingContentCategory{}).
		Where("Landing_content_id = ?", LandingContentID).
		Pluck("category_id", &categoryIDs).Error; err != nil {
		return nil, err
	}

	var filteredCategories []models.Category
	err := r.db.Model(&models.Category{}).
		Joins("JOIN category_types ON category_types.id = categories.category_type_id").
		Where("categories.id IN ?", categoryIDs).
		Where("category_types.type_code = ?", categoryTypeCode).
		Find(&filteredCategories).Error

	if err != nil {
		return nil, err
	}

	return filteredCategories, nil
}

func (r *CMSLandingPageRepository) GetRevisionByLandingPageId(pageId uuid.UUID, language string) ([]models.Revision, error) {
	var landingPage models.LandingPage
	err := r.db.
		Preload("Contents", "language = ?", language).
		Preload("Contents.Revision").
		Where("id = ?", pageId).
		First(&landingPage).Error

	if err != nil {
		return nil, err
	}

	var revisions []models.Revision
	for _, content := range landingPage.Contents {
		if content.Revision != nil {
			revisions = append(revisions, *content.Revision)
		}
	}

	// Sort by created_at
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].CreatedAt.After(revisions[j].CreatedAt)
	})

	return revisions, nil
}

func (r *CMSLandingPageRepository) IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error) {
	var count int64
	if pageId == uuid.Nil { // No page yet
		err := r.db.Model(&models.LandingContent{}).
			Where("url_alias = ?", urlAlias).
			Count(&count).Error

		if err != nil {
			return false, err
		}

		return count > 0, nil
	}

	err := r.db.Model(&models.LandingContent{}).
		Where("url_alias = ? AND page_id != ?", urlAlias, pageId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CMSLandingPageRepository) GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error) {
	var landingContent models.LandingContent

	err := r.db.
		First(&landingContent, "id = ?", contentId).Error

	if err != nil {
		return uuid.Nil, err
	}

	return landingContent.PageID, nil
}

func (r *CMSLandingPageRepository) CreateLandingContentPreview(landingContentPreview *models.LandingContent) (*models.LandingContent, error) {
	landingContentPreview.Revision = nil
	landingContentPreview.Categories = nil
	if err := r.db.Create(landingContentPreview).Error; err != nil {
		return nil, err
	}

	return landingContentPreview, nil
}

func (r *CMSLandingPageRepository) UpdateLandingContentPreview(landingContentPreview *models.LandingContent) (*models.LandingContent, error) {
	// Overwrite full record (All fields)
	landingContentPreview.Revision = nil
	landingContentPreview.Categories = nil
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Delete the old components
		if err := tx.Where("landing_content_id = ?", landingContentPreview.ID).Delete(&models.Component{}).Error; err != nil {
			return err
		}

		if err := tx.Save(landingContentPreview).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return landingContentPreview, nil
}

func (r *CMSLandingPageRepository) FindLandingContentPreviewById(pageId uuid.UUID, language string) (*models.LandingContent, error) {
	var landingContentPreview models.LandingContent
	err := r.db.
		Where("page_id = ? AND language = ? AND mode = ?", pageId, string(language), "Preview").
		First(&landingContentPreview).Error

	if err != nil {
		return nil, err
	}

	return &landingContentPreview, nil
}
