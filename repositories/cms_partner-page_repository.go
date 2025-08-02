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

type CMSPartnerPageRepositoryInterface interface {
	CreatePartnerPage(PartnerPage *models.PartnerPage) (*models.PartnerPage, error)
	FindAllPartnerPage(query dto.PartnerPageQuery, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error)
	FindPartnerPageById(id uuid.UUID) (*models.PartnerPage, error)
	UpdatePartnerContent(updatePartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error)
	DeletePartnerPage(id uuid.UUID) error
	FindContentByPartnerPageId(pageId uuid.UUID, language string, mode string) (*models.PartnerContent, error)
	FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.PartnerContent, error)
	// Deprecate
	CreateContentForPartnerPage(PartnerContent *models.PartnerContent, lang string, mode string) (*models.PartnerContent, error)
	DeletePartnerContent(pageId uuid.UUID, lang, mode string) error
	DuplicatePartnerPage(pageId uuid.UUID) (*models.PartnerPage, error)
	DuplicatePartnerContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error)
	RevertPartnerContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error)
	GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	GetRevisionByPartnerPageId(pageId uuid.UUID, language string) ([]models.Revision, error)
	IsUrlDuplicate(url string, pageId uuid.UUID) (bool, error)
	IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error)
	GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error)
	CreatePartnerContentPreview(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error)
	UpdatePartnerContentPreview(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error)
	FindPartnerContentPreviewById(pageId uuid.UUID, language string) (*models.PartnerContent, error)
}

type CMSPartnerPageRepository struct {
	db *gorm.DB
}

func NewCMSPartnerPageRepository(db *gorm.DB) *CMSPartnerPageRepository {
	return &CMSPartnerPageRepository{db: db}
}

func (r *CMSPartnerPageRepository) CreatePartnerPage(PartnerPage *models.PartnerPage) (*models.PartnerPage, error) {
	err := r.db.Transaction(func(r *gorm.DB) error {
		// Create the PartnerPage (and its content)
		if err := r.Create(PartnerPage).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return PartnerPage, nil
}

// Original postgres
// SELECT Partner_pages.* FROM Partner_pages
// JOIN Partner_contents ON Partner_contents.page_id = Partner_pages.id
// LEFT JOIN Partner_content_categories ON Partner_content_categories.Partner_content_id = Partner_contents.id
// JOIN categories ON Partner_content_categories.category_id = categories.id
// JOIN category_types ON category_types.id = categories.category_type_id
// WHERE category_types.type_code = 'Partner' AND categories.name ILIKE '%filter.CategoryPartner%'
// INTERSECT
// SELECT Partner_pages.* FROM Partner_pages
// JOIN Partner_contents ON Partner_contents.page_id = Partner_pages.id
// LEFT JOIN Partner_content_categories ON Partner_content_categories.Partner_content_id = Partner_contents.id
// JOIN categories ON Partner_content_categories.category_id = categories.id
// JOIN category_types ON category_types.id = categories.category_type_id
// WHERE category_types.type_code = 'category-keywords' AND categories.name ILIKE '%filter.CategoryPartner%'
func (r *CMSPartnerPageRepository) FindAllPartnerPage(query dto.PartnerPageQuery, sort string, page, limit int, language string) ([]models.PartnerPage, int64, error) {
	var partnerPages []models.PartnerPage
	var totalCount int64

	// Build base query with proper joins and filters
	baseQuery := r.db.Model(&models.PartnerPage{}).
		Joins("JOIN partner_contents ON partner_contents.page_id = partner_pages.id").
		Where("partner_contents.mode != ? AND partner_contents.mode != ?", "Histories", "Preview")

	// Content filters
	if query.Title != "" {
		baseQuery = baseQuery.Where("partner_contents.title ILIKE ?", "%"+query.Title+"%")
	}
	if query.UrlAlias != "" {
		baseQuery = baseQuery.Where("partner_contents.url_alias ILIKE ?", "%"+query.UrlAlias+"%")
	}
	if query.URL != "" {
		baseQuery = baseQuery.Where("partner_contents.url ILIKE ?", "%"+query.URL+"%")
	}
	if query.Status != "" {
		baseQuery = baseQuery.Where("partner_contents.workflow_status = ?", query.Status)
	}
	if language != "" {
		baseQuery = baseQuery.Where("partner_contents.language = ?", language)
	}

	// Category filters
	categoryFilters := map[string]string{
		"partner":            query.CategoryPartner,
		"category-keywords":  query.CategoryKeywords,
		"category-scale":     query.CategoryScale,
		"category-industry":  query.CategoryIndustry,
		"category-goal":      query.CategoryGoal,
		"category-functions": query.CategoryFunctions,
	}

	for typeCode, filterValue := range categoryFilters {
		if filterValue != "" {
			subQuery := r.db.Table("partner_content_categories").
				Select("partner_content_categories.partner_content_id").
				Joins("JOIN categories ON partner_content_categories.category_id = categories.id").
				Joins("JOIN category_types ON category_types.id = categories.category_type_id").
				Where("category_types.type_code = ? AND categories.name ILIKE ?", typeCode, "%"+filterValue+"%")
			baseQuery = baseQuery.Where("partner_contents.id IN (?)", subQuery)
		}
	}

	// Clone query for counting to avoid modification
	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Model(&models.PartnerPage{}).Distinct("partner_pages.id").Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	finalQuery := baseQuery

	// Apply sorting
	sortColumn := "partner_pages.created_at DESC" // default sort
	sortableColumns := map[string]string{
		"title":      "partner_contents.title",
		"url_alias":  "partner_contents.url_alias",
		"url":        "partner_contents.url",
		"status":     "partner_contents.workflow_status",
		"created_at": "partner_pages.created_at",
		"updated_at": "partner_pages.updated_at",
	}

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
				if strings.HasPrefix(col, "partner_contents.") {
					isSortOnJoinedTable = true
					// Use MIN() for deterministic sorting on joined column
					sortColumn = fmt.Sprintf("MIN(%s) %s", col, direction)
				} else {
					sortColumn = fmt.Sprintf("%s %s", col, direction)
				}
			}
		}
	}

	if isSortOnJoinedTable {
		// Use GROUP BY when sorting on joined table columns
		finalQuery = finalQuery.Select("partner_pages.*").Group("partner_pages.id")
	} else {
		// Use DISTINCT for primary table columns
		finalQuery = finalQuery.Select("DISTINCT partner_pages.*")
	}

	// Get paginated and sorted pages with all necessary preloads
	offset := (page - 1) * limit
	err := finalQuery.
		Order(sortColumn).
		Offset(offset).
		Limit(limit).
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("partner_contents.mode != ? AND partner_contents.mode != ?", "Histories", "Preview").
				Order("partner_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		Find(&partnerPages).Error

	if err != nil {
		return nil, 0, err
	}

	return partnerPages, totalCount, nil
}
func (r *CMSPartnerPageRepository) FindPartnerPageById(id uuid.UUID) (*models.PartnerPage, error) {
	var PartnerPage models.PartnerPage

	err := r.db.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("partner_contents.mode != ? AND partner_contents.mode != ?", "Histories", "Preview").
				Order("partner_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		First(&PartnerPage, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &PartnerPage, nil
}

func (r *CMSPartnerPageRepository) UpdatePartnerContent(updatePartnerContent *models.PartnerContent, prevContentId uuid.UUID) (*models.PartnerContent, error) {

	err := r.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("[REPO-BEFORE-NORMALIZE] Language: '%s'", updatePartnerContent.Language)

		var oldContent models.PartnerContent
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

		updatePartnerContent.PageID = oldContent.PageID

		if err := helpers.NormalizePartnerContent(updatePartnerContent); err != nil {
			return err
		}
		log.Printf("[REPO-AFTER-NORMALIZE] Language: '%s'", updatePartnerContent.Language)
		if updatePartnerContent.Revision != nil {
			if err := helpers.NormalizeRevision(updatePartnerContent.Revision); err != nil {
				return err
			}
		}

		updatePartnerContent.ID = uuid.Nil
		updatePartnerContent.CreatedAt = time.Time{} // Reset เพื่อให้ DB generate
		updatePartnerContent.UpdatedAt = time.Time{}

		if updatePartnerContent.MetaTag != nil {
			updatePartnerContent.MetaTag.ID = uuid.Nil
			updatePartnerContent.MetaTagID = uuid.Nil // Clear FK
		}
		if updatePartnerContent.Revision != nil {
			updatePartnerContent.Revision.ID = uuid.Nil
			updatePartnerContent.Revision.PartnerContentID = nil // Clear FK
		}
		for _, component := range updatePartnerContent.Components {
			if component != nil {
				component.ID = uuid.Nil
				component.PartnerContentID = nil // Clear FK
			}
		}

		if err := tx.Create(updatePartnerContent).Error; err != nil {
			return fmt.Errorf("failed to create new content version: %w", err)
		}

		if err := tx.Model(&models.PartnerPage{}).Where("id = ?", updatePartnerContent.PageID).Update("updated_at", time.Now()).Error; err != nil {
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
		Preload("Categories.CategoryType").
		First(&updatePartnerContent, "id = ?", updatePartnerContent.ID).Error; err != nil {
		log.Printf("Warning: failed to fully preload associations for new content %s: %v", updatePartnerContent.ID, err)

	}

	return updatePartnerContent, nil
}
func (r *CMSPartnerPageRepository) DeletePartnerPage(id uuid.UUID) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		var contents []models.PartnerContent

		// Step 0: Ensure the PartnerPage exists
		var page models.PartnerPage
		if err := tx.First(&page, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.ErrNotFound
			}
			return err
		}

		// Step 1: Get all PartnerContent entries for this page
		if err := tx.Where("page_id = ?", id).Find(&contents).Error; err != nil {
			return err
		}

		// Step 2: Collect all PartnerContent IDs
		var contentIDs []uuid.UUID
		for _, content := range contents {
			contentIDs = append(contentIDs, content.ID)
		}

		if len(contentIDs) > 0 {
			// Step 3: Delete Components
			if err := tx.Where("partner_content_id IN ?", contentIDs).Delete(&models.Component{}).Error; err != nil {
				return err
			}

			// Step 4: Delete PartnerContentCategories
			if err := tx.Where("partner_content_id IN ?", contentIDs).Delete(&models.PartnerContentCategory{}).Error; err != nil {
				return err
			}

			// Step 5: Delete Revisions
			if err := tx.Where("partner_content_id IN ?", contentIDs).Delete(&models.Revision{}).Error; err != nil {
				return err
			}
		}

		// Step 6: Delete PartnerContents
		if err := tx.Where("page_id = ?", id).Delete(&models.PartnerContent{}).Error; err != nil {
			return err
		}

		// Step 7: Delete PartnerPage
		if err := tx.Where("id = ?", id).Delete(&models.PartnerPage{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CMSPartnerPageRepository) FindContentByPartnerPageId(pageId uuid.UUID, language string, mode string) (*models.PartnerContent, error) {
	var PartnerContent models.PartnerContent

	err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Categories.CategoryType").
		Preload("Components").
		Where("page_id = ? AND language = ? AND mode = ?", pageId, language, mode).
		First(&PartnerContent).Error

	if err != nil {
		return nil, err
	}

	return &PartnerContent, nil
}

func (r *CMSPartnerPageRepository) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
	var PartnerContent models.PartnerContent
	if err := r.db.
		Where("page_id = ? AND language = ?", pageId, language).
		Order("created_at DESC").
		First(&PartnerContent).Error; err != nil {
		return nil, err
	}

	return &PartnerContent, nil
}

// Might be deprecate
func (r *CMSPartnerPageRepository) CreateContentForPartnerPage(PartnerContent *models.PartnerContent, lang string, mode string) (*models.PartnerContent, error) {
	if err := r.db.Create(PartnerContent).Error; err != nil {
		return nil, err
	}

	var createdPartnerContent models.PartnerContent
	if err := r.db.First(&createdPartnerContent, "id = ?", PartnerContent.ID).Error; err != nil {
		return nil, err
	}

	return &createdPartnerContent, nil
}

func (r *CMSPartnerPageRepository) DeletePartnerContent(pageId uuid.UUID, lang, mode string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		var PartnerContent models.PartnerContent
		if err := tx.First(&PartnerContent, "page_id = ? AND language = ? AND mode = ?", pageId, lang, mode).Error; err != nil {
			return err
		}

		if err := tx.Where("partner_content_id = ?", PartnerContent.ID).Delete(&models.Component{}).Error; err != nil {
			return err
		}

		if err := tx.Where("partner_content_id = ?", PartnerContent.ID).Delete(&models.Revision{}).Error; err != nil {
			return err
		}

		if err := tx.Where("partner_content_id = ?", PartnerContent.ID).Delete(&models.PartnerContentCategory{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&PartnerContent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CMSPartnerPageRepository) DuplicatePartnerPage(pageId uuid.UUID) (*models.PartnerPage, error) {
	// Already a pointer
	partnerPage := &models.PartnerPage{}

	err := r.db.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("partner_contents.mode != ?", "Histories").
				Order("partner_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		First(&partnerPage, "id = ?", pageId).Error

	if err != nil {
		return nil, err
	}

	if len(partnerPage.Contents) == 0 {
		return nil, errs.ErrNoNewContentToDuplicate
	}

	// For loop to duplicate two of the main language (en, th)
	partnerContents := partnerPage.Contents
	newContents := []*models.PartnerContent{}
	for _, originalPartnerContent := range partnerContents {

		// Create a deep copy
		copyPartnerContent := *originalPartnerContent

		// Metatag
		if originalPartnerContent.MetaTag != nil {
			metaTagCopy := *originalPartnerContent.MetaTag
			metaTagCopy.ID = uuid.Nil
			metaTagCopy.CreatedAt = time.Time{}
			metaTagCopy.UpdatedAt = time.Time{}
			copyPartnerContent.MetaTag = &metaTagCopy
		}

		// Revision
		if originalPartnerContent.Revision != nil {
			revisionCopy := *originalPartnerContent.Revision
			revisionCopy.ID = uuid.Nil
			revisionCopy.PartnerContentID = nil
			revisionCopy.CreatedAt = time.Time{}
			revisionCopy.UpdatedAt = time.Time{}
			copyPartnerContent.Revision = &revisionCopy
		}

		// Components
		if originalPartnerContent.Components != nil {
			var componentCopies []*models.Component
			for _, component := range originalPartnerContent.Components {
				comp := *component
				comp.ID = uuid.Nil
				comp.PartnerContentID = nil
				comp.CreatedAt = time.Time{}
				comp.UpdatedAt = time.Time{}
				componentCopies = append(componentCopies, &comp)
			}
			copyPartnerContent.Components = componentCopies
		}

		// Url and UrlAlias: Random a string of length 3 and concat it to the back
		copyPartnerContent.URL += "-" + helpers.RandomString(3)

		// Only for partner and partner page, because both of these pages can have no alias
		if copyPartnerContent.URLAlias != "" {
			copyPartnerContent.URLAlias += "-" + helpers.RandomString(3)
		}

		// Content
		copyPartnerContent.ID = uuid.Nil
		copyPartnerContent.PageID = uuid.Nil
		copyPartnerContent.CreatedAt = time.Time{}
		copyPartnerContent.UpdatedAt = time.Time{}

		// Attach the new content
		newContents = append(newContents, &copyPartnerContent)
	}

	// Partner Page
	copyPartnerPage := *partnerPage
	copyPartnerPage.ID = uuid.Nil
	copyPartnerPage.CreatedAt = time.Time{}
	copyPartnerPage.UpdatedAt = time.Time{}

	// Use the new contents
	copyPartnerPage.Contents = newContents

	// Create the new partnerPage
	if err := r.db.Create(&copyPartnerPage).Error; err != nil {
		return nil, err
	}

	return &copyPartnerPage, nil
}

func (r *CMSPartnerPageRepository) DuplicatePartnerContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	var PartnerContent models.PartnerContent
	if err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Components").
		Preload("MetaTag").
		First(&PartnerContent, "id = ?", contentId).Error; err != nil {
		return nil, err
	}

	// Change language
	if PartnerContent.Language == enums.PageLanguageTH {
		PartnerContent.Language = enums.PageLanguageEN
	} else {
		PartnerContent.Language = enums.PageLanguageTH
	}

	PartnerContent.ID = uuid.Nil

	// Metatag
	PartnerContent.MetaTagID = uuid.Nil
	if PartnerContent.MetaTag != nil {
		PartnerContent.MetaTag.ID = uuid.Nil
		PartnerContent.MetaTag.CreatedAt = time.Time{}
		PartnerContent.MetaTag.UpdatedAt = time.Time{}
	}

	// Content
	PartnerContent.ID = uuid.Nil
	PartnerContent.CreatedAt = time.Time{}
	PartnerContent.UpdatedAt = time.Time{}

	// Attach revision, always needs a revision for keeping the history (but this is a new content, so it's a start of the history)
	if newRevision == nil {
		return nil, errs.ErrNoRevisionFound
	}
	PartnerContent.Revision = newRevision

	// Components
	if PartnerContent.Components != nil {
		for _, component := range PartnerContent.Components {
			component.CreatedAt = time.Time{}
			component.UpdatedAt = time.Time{}
			component.ID = uuid.Nil
			component.PartnerContentID = nil
		}
	}

	if err := r.db.Create(&PartnerContent).Error; err != nil {
		return nil, err
	}

	return &PartnerContent, nil
}

func (r *CMSPartnerPageRepository) RevertPartnerContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.PartnerContent, error) {
	var revision models.Revision
	// Get the revision
	if err := r.db.First(&revision, "id = ?", revisionId).Error; err != nil {
		return nil, err
	}

	PartnerContent := &models.PartnerContent{}
	// Get the content of that revision
	if err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Components").
		Preload("MetaTag").
		First(PartnerContent, "id = ?", revision.PartnerContentID).Error; err != nil {
		return nil, err
	}

	// Get the old content
	oldContent, err := r.FindLatestContentByPageId(PartnerContent.PageID, string(PartnerContent.Language))
	if err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Update the old content

		oldContent.Mode = enums.PageModeHistories
		if err := tx.Save(oldContent).Error; err != nil {
			return err
		}

		// Change the new one to draft
		PartnerContent.Mode = enums.PageModeDraft

		// Metatag
		PartnerContent.MetaTagID = uuid.Nil
		if PartnerContent.MetaTag != nil {
			PartnerContent.MetaTag.ID = uuid.Nil
			PartnerContent.MetaTag.CreatedAt = time.Time{}
			PartnerContent.MetaTag.UpdatedAt = time.Time{}
		}

		// Content
		PartnerContent.ID = uuid.Nil
		PartnerContent.CreatedAt = time.Time{}
		PartnerContent.UpdatedAt = time.Time{}

		// Attach revision, always needs a revision for keeping the history (but this is a new content, so it's a start of the history)
		if newRevision == nil {
			return errs.ErrNoRevisionFound
		}
		PartnerContent.Revision = newRevision

		// Components
		if PartnerContent.Components != nil {
			for _, component := range PartnerContent.Components {
				component.CreatedAt = time.Time{}
				component.UpdatedAt = time.Time{}
				component.ID = uuid.Nil
				component.PartnerContentID = nil
			}
		}

		if err := tx.Create(PartnerContent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return PartnerContent, nil
}

func (r *CMSPartnerPageRepository) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	var result struct {
		ID uuid.UUID
	}
	if err := r.db.
		Model(&models.PartnerContent{}).
		Select("id").
		Where("page_id = ? AND language = ? AND mode = ?", pageId, language, mode).
		Take(&result).Error; err != nil {
		return nil, err
	}

	PartnerContentID := result.ID

	var categoryIDs []uuid.UUID
	if err := r.db.
		Model(&models.PartnerContentCategory{}).
		Where("Partner_content_id = ?", PartnerContentID).
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

func (r *CMSPartnerPageRepository) GetRevisionByPartnerPageId(pageId uuid.UUID, language string) ([]models.Revision, error) {
	var PartnerPage models.PartnerPage
	err := r.db.
		Preload("Contents", "language = ?", language).
		Preload("Contents.Revision").
		Where("id = ?", pageId).
		First(&PartnerPage).Error

	if err != nil {
		return nil, err
	}

	var revisions []models.Revision
	for _, content := range PartnerPage.Contents {
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

func (r *CMSPartnerPageRepository) IsUrlDuplicate(url string, pageId uuid.UUID) (bool, error) {
	var count int64
	if pageId == uuid.Nil { // No page yet
		err := r.db.Model(&models.PartnerContent{}).
			Where("url = ?", url).
			Count(&count).Error

		if err != nil {
			return false, err
		}

		return count > 0, nil
	}

	err := r.db.Model(&models.PartnerContent{}).
		Where("url = ? AND page_id != ?", url, pageId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CMSPartnerPageRepository) IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error) {
	var count int64
	if pageId == uuid.Nil { // No page yet
		err := r.db.Model(&models.PartnerContent{}).
			Where("url_alias = ?", urlAlias).
			Count(&count).Error

		if err != nil {
			return false, err
		}

		return count > 0, nil
	}

	err := r.db.Model(&models.PartnerContent{}).
		Where("url_alias = ? AND page_id != ?", urlAlias, pageId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CMSPartnerPageRepository) GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error) {
	var partnerContent models.PartnerContent

	err := r.db.
		First(&partnerContent, "id = ?", contentId).Error

	if err != nil {
		return uuid.Nil, err
	}

	return partnerContent.PageID, nil
}

func (r *CMSPartnerPageRepository) CreatePartnerContentPreview(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error) {
	partnerContentPreview.Revision = nil
	partnerContentPreview.Categories = nil
	if err := r.db.Create(partnerContentPreview).Error; err != nil {
		return nil, err
	}

	return partnerContentPreview, nil
}

func (r *CMSPartnerPageRepository) UpdatePartnerContentPreview(partnerContentPreview *models.PartnerContent) (*models.PartnerContent, error) {
	// Overwrite full record (All fields)
	partnerContentPreview.Revision = nil
	partnerContentPreview.Categories = nil
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Delete the old components

		if err := tx.Where("partner_content_id = ?", partnerContentPreview.ID).Delete(&models.Component{}).Error; err != nil {
			return err
		}

		if err := tx.Save(partnerContentPreview).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return partnerContentPreview, nil
}

func (r *CMSPartnerPageRepository) FindPartnerContentPreviewById(pageId uuid.UUID, language string) (*models.PartnerContent, error) {
	var partnerContentPreview models.PartnerContent
	err := r.db.
		Where("page_id = ? AND language = ? AND mode = ?", pageId, string(language), "Preview").
		First(&partnerContentPreview).Error

	if err != nil {
		return nil, err
	}

	return &partnerContentPreview, nil
}
