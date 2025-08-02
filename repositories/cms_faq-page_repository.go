package repositories

import (
	"errors"
	"fmt"
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

type CMSFaqPageRepositoryInterface interface {
	CreateFaqPage(faqPage *models.FaqPage) (*models.FaqPage, error)
	FindAllFaqPage(query dto.FaqPageQuery, sort string, page, limit int, language string) ([]models.FaqPage, int64, error)
	FindFaqPageById(id uuid.UUID) (*models.FaqPage, error)
	UpdateFaqContent(updateFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error)
	DeleteFaqPage(id uuid.UUID) error
	FindContentByFaqPageId(pageId uuid.UUID, language string, mode string) (*models.FaqContent, error)
	FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.FaqContent, error)
	// Deprecate
	CreateContentForFaqPage(faqContent *models.FaqContent, lang string, mode string) (*models.FaqContent, error)
	DeleteFaqContent(pageId uuid.UUID, lang, mode string) error
	DuplicateFaqPage(pageId uuid.UUID) (*models.FaqPage, error)
	DuplicateFaqContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error)
	RevertFaqContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error)
	GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error)
	GetRevisionByFaqPageId(pageId uuid.UUID, language string) ([]models.Revision, error)
	IsUrlDuplicate(url string, pageId uuid.UUID) (bool, error)
	IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error)
	GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error)
	CreateFaqContentPreview(faqContentPreview *models.FaqContent) (*models.FaqContent, error)
	UpdateFaqContentPreview(faqContentPreview *models.FaqContent) (*models.FaqContent, error)
	FindFaqContentPreviewById(pageId uuid.UUID, language string) (*models.FaqContent, error)
}

type CMSFaqPageRepository struct {
	db *gorm.DB
}

func NewCMSFaqPageRepository(db *gorm.DB) *CMSFaqPageRepository {
	return &CMSFaqPageRepository{db: db}
}

func (r *CMSFaqPageRepository) CreateFaqPage(faqPage *models.FaqPage) (*models.FaqPage, error) {
	// Create the faqPage (and its content)
	if err := r.db.Create(faqPage).Error; err != nil {
		return nil, err
	}

	return faqPage, nil
}

// Original postgres
// SELECT faq_pages.* FROM faq_pages
// JOIN faq_contents ON faq_contents.page_id = faq_pages.id
// LEFT JOIN faq_content_categories ON faq_content_categories.faq_content_id = faq_contents.id
// JOIN categories ON faq_content_categories.category_id = categories.id
// JOIN category_types ON category_types.id = categories.category_type_id
// WHERE category_types.type_code = 'faq' AND categories.name ILIKE '%filter.CategoryFaq%'
// INTERSECT
// SELECT faq_pages.* FROM faq_pages
// JOIN faq_contents ON faq_contents.page_id = faq_pages.id
// LEFT JOIN faq_content_categories ON faq_content_categories.faq_content_id = faq_contents.id
// JOIN categories ON faq_content_categories.category_id = categories.id
// JOIN category_types ON category_types.id = categories.category_type_id
// WHERE category_types.type_code = 'category-keywords' AND categories.name ILIKE '%filter.CategoryFaq%'
func (r *CMSFaqPageRepository) FindAllFaqPage(query dto.FaqPageQuery, sort string, page, limit int, language string) ([]models.FaqPage, int64, error) {
	var faqPages []models.FaqPage
	var totalCount int64

	// Build base query with proper joins and filters
	baseQuery := r.db.Model(&models.FaqPage{}).
		Joins("JOIN faq_contents ON faq_contents.page_id = faq_pages.id").
		Where("faq_contents.mode != ? AND faq_contents.mode != ?", "Histories", "Preview")

	// Content filters
	if query.Title != "" {
		baseQuery = baseQuery.Where("faq_contents.title ILIKE ?", "%"+query.Title+"%")
	}
	if query.UrlAlias != "" {
		baseQuery = baseQuery.Where("faq_contents.url_alias ILIKE ?", "%"+query.UrlAlias+"%")
	}
	if query.URL != "" {
		baseQuery = baseQuery.Where("faq_contents.url ILIKE ?", "%"+query.URL+"%")
	}
	if query.Status != "" {
		baseQuery = baseQuery.Where("faq_contents.workflow_status = ?", query.Status)
	}
	if language != "" {
		baseQuery = baseQuery.Where("faq_contents.language = ?", language)
	}

	// Category filters
	categoryFilters := map[string]string{
		"faq":               query.CategoryFaq,
		"category-keywords": query.CategoryKeywords,
	}
	for typeCode, filterValue := range categoryFilters {
		if filterValue != "" {
			subQuery := r.db.Table("faq_content_categories").
				Select("faq_content_categories.faq_content_id").
				Joins("JOIN categories ON faq_content_categories.category_id = categories.id").
				Joins("JOIN category_types ON category_types.id = categories.category_type_id").
				Where("category_types.type_code = ? AND categories.name ILIKE ?", typeCode, "%"+filterValue+"%")
			baseQuery = baseQuery.Where("faq_contents.id IN (?)", subQuery)
		}
	}

	// Clone query for counting
	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Model(&models.FaqPage{}).Distinct("faq_pages.id").Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	finalQuery := baseQuery

	// Apply sorting
	sortColumn := "faq_pages.created_at DESC" // default sort
	sortableColumns := map[string]string{
		"title":           "faq_contents.title",
		"url_alias":       "faq_contents.url_alias",
		"url":             "faq_contents.url",
		"status":          "faq_contents.workflow_status",
		"workflow_status": "faq_contents.workflow_status",
		"language":        "faq_contents.language",
		"created_at":      "faq_pages.created_at",
		"updated_at":      "faq_pages.updated_at",
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
				if strings.HasPrefix(col, "faq_contents.") {
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
		finalQuery = finalQuery.Select("faq_pages.*").Group("faq_pages.id")
	} else {
		// Use DISTINCT for primary table columns
		finalQuery = finalQuery.Select("DISTINCT faq_pages.*")
	}

	// Get paginated and sorted pages with all necessary preloads
	offset := (page - 1) * limit
	err := finalQuery.
		Order(sortColumn).
		Offset(offset).
		Limit(limit).
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("faq_contents.mode != ? AND faq_contents.mode != ?", "Histories", "Preview").
				Order("faq_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		Find(&faqPages).Error

	if err != nil {
		return nil, 0, err
	}

	return faqPages, totalCount, nil
}
func (r *CMSFaqPageRepository) FindFaqPageById(id uuid.UUID) (*models.FaqPage, error) {
	var faqPage models.FaqPage

	err := r.db.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("faq_contents.mode != ? AND faq_contents.mode != ?", "Histories", "Preview").
				Order("faq_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		First(&faqPage, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &faqPage, nil
}

func (r *CMSFaqPageRepository) UpdateFaqContent(updateFaqContent *models.FaqContent, prevContentId uuid.UUID) (*models.FaqContent, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Update the previous content mode to history

		var faqContent models.FaqContent
		if err := tx.First(&faqContent, "id = ?", prevContentId).Error; err != nil {
			return err
		}
		faqContent.Mode = enums.PageModeHistories

		// Save the updated content
		if err := tx.Save(&faqContent).Error; err != nil {
			return err
		}

		// Attach the page id to the updated content
		updateFaqContent.PageID = faqContent.PageID

		// Can't update to another language, because we only allow that from duplicate to another language
		updateFaqContent.Language = faqContent.Language

		now := time.Now()
		updateFaqContent.CreatedAt = now
		updateFaqContent.UpdatedAt = now

		// Create new content for that page
		if err := tx.Create(updateFaqContent).Error; err != nil {
			return err
		}

		// Update the page's updated_at to the current time
		if err := tx.Model(&models.FaqPage{}).Where("id = ?", updateFaqContent.PageID).Update("updated_at", now).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return updateFaqContent, nil
}

func (r *CMSFaqPageRepository) DeleteFaqPage(id uuid.UUID) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		var contents []models.FaqContent

		// Step 0: Ensure the FaqPage exists
		var page models.FaqPage
		if err := tx.First(&page, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.ErrNotFound
			}
			return err
		}

		// Step 1: Get all FaqContent entries for this page
		if err := tx.Where("page_id = ?", id).Find(&contents).Error; err != nil {
			return err
		}

		// Step 2: Collect all FaqContent IDs
		var contentIDs []uuid.UUID
		for _, content := range contents {
			contentIDs = append(contentIDs, content.ID)
		}

		if len(contentIDs) > 0 {
			// Step 3: Delete Components
			if err := tx.Where("faq_content_id IN ?", contentIDs).Delete(&models.Component{}).Error; err != nil {
				return err
			}

			// Step 4: Delete FaqContentCategories
			if err := tx.Where("faq_content_id IN ?", contentIDs).Delete(&models.FaqContentCategory{}).Error; err != nil {
				return err
			}

			// Step 5: Delete Revisions
			if err := tx.Where("faq_content_id IN ?", contentIDs).Delete(&models.Revision{}).Error; err != nil {
				return err
			}
		}

		// Step 6: Delete FaqContents
		if err := tx.Where("page_id = ?", id).Delete(&models.FaqContent{}).Error; err != nil {
			return err
		}

		// Step 7: Delete FaqPage
		if err := tx.Where("id = ?", id).Delete(&models.FaqPage{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CMSFaqPageRepository) FindContentByFaqPageId(pageId uuid.UUID, language string, mode string) (*models.FaqContent, error) {
	var faqContent models.FaqContent

	err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Categories.CategoryType").
		Preload("Components").
		Where("page_id = ? AND language = ? AND mode = ?", pageId, language, mode).
		First(&faqContent).Error

	if err != nil {
		return nil, err
	}

	return &faqContent, nil
}

func (r *CMSFaqPageRepository) FindLatestContentByPageId(pageId uuid.UUID, language string) (*models.FaqContent, error) {
	var faqContent models.FaqContent
	if err := r.db.
		Where("page_id = ? AND language = ?", pageId, language).
		Order("created_at DESC").
		First(&faqContent).Error; err != nil {
		return nil, err
	}

	return &faqContent, nil
}

// Might be deprecate
func (r *CMSFaqPageRepository) CreateContentForFaqPage(faqContent *models.FaqContent, lang string, mode string) (*models.FaqContent, error) {
	if err := r.db.Create(faqContent).Error; err != nil {
		return nil, err
	}

	var createdFaqContent models.FaqContent
	if err := r.db.First(&createdFaqContent, "id = ?", faqContent.ID).Error; err != nil {
		return nil, err
	}

	return &createdFaqContent, nil
}

func (r *CMSFaqPageRepository) DeleteFaqContent(pageId uuid.UUID, lang, mode string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var faqContent models.FaqContent
		if err := tx.First(&faqContent, "page_id = ? AND language = ? AND mode = ?", pageId, lang, mode).Error; err != nil {
			return err
		}

		if err := tx.Where("faq_content_id = ?", faqContent.ID).Delete(&models.Component{}).Error; err != nil {
			return err
		}

		if err := tx.Where("faq_content_id = ?", faqContent.ID).Delete(&models.Revision{}).Error; err != nil {
			return err
		}

		if err := tx.Where("faq_content_id = ?", faqContent.ID).Delete(&models.FaqContentCategory{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&faqContent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CMSFaqPageRepository) DuplicateFaqPage(pageId uuid.UUID) (*models.FaqPage, error) {
	// Already a pointer
	faqPage := &models.FaqPage{}

	err := r.db.
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("faq_contents.mode != ?", "Histories").
				Order("faq_contents.created_at DESC")
		}).
		Preload("Contents.Revision").
		Preload("Contents.Categories").
		Preload("Contents.Categories.CategoryType").
		Preload("Contents.Components").
		Preload("Contents.MetaTag").
		First(&faqPage, "id = ?", pageId).Error

	if err != nil {
		return nil, err
	}

	if len(faqPage.Contents) == 0 {
		return nil, errs.ErrNoNewContentToDuplicate
	}

	// For loop to duplicate two of the main language (en, th)
	faqContents := faqPage.Contents
	newContents := []*models.FaqContent{}
	for _, originalFaqContent := range faqContents {

		// Create a deep copy
		copyFaqContent := *originalFaqContent

		// Metatag
		if originalFaqContent.MetaTag != nil {
			metaTagCopy := *originalFaqContent.MetaTag
			metaTagCopy.ID = uuid.Nil
			metaTagCopy.CreatedAt = time.Time{}
			metaTagCopy.UpdatedAt = time.Time{}
			copyFaqContent.MetaTag = &metaTagCopy
		}

		// Revision
		if originalFaqContent.Revision != nil {
			revisionCopy := *originalFaqContent.Revision
			revisionCopy.ID = uuid.Nil
			revisionCopy.FaqContentID = nil
			revisionCopy.CreatedAt = time.Time{}
			revisionCopy.UpdatedAt = time.Time{}
			copyFaqContent.Revision = &revisionCopy
		}

		// Components
		if originalFaqContent.Components != nil {
			var componentCopies []*models.Component
			for _, component := range originalFaqContent.Components {
				comp := *component
				comp.ID = uuid.Nil
				comp.FaqContentID = nil
				comp.CreatedAt = time.Time{}
				comp.UpdatedAt = time.Time{}
				componentCopies = append(componentCopies, &comp)
			}
			copyFaqContent.Components = componentCopies
		}

		// Url and UrlAlias: Random a string of length 3 and concat it to the back
		copyFaqContent.URL += "-" + helpers.RandomString(3)

		// Only for faq and faq page, because both of these pages can have no alias
		if copyFaqContent.URLAlias != "" {
			copyFaqContent.URLAlias += "-" + helpers.RandomString(3)
		}

		// Content
		copyFaqContent.ID = uuid.Nil
		copyFaqContent.PageID = uuid.Nil
		copyFaqContent.CreatedAt = time.Time{}
		copyFaqContent.UpdatedAt = time.Time{}

		// Attach the new content
		newContents = append(newContents, &copyFaqContent)
	}

	// Faq Page
	copyFaqPage := *faqPage
	copyFaqPage.ID = uuid.Nil
	copyFaqPage.CreatedAt = time.Time{}
	copyFaqPage.UpdatedAt = time.Time{}

	// Use the new contents
	copyFaqPage.Contents = newContents

	// Create the new faqPage
	if err := r.db.Create(&copyFaqPage).Error; err != nil {
		return nil, err
	}

	return &copyFaqPage, nil
}

func (r *CMSFaqPageRepository) DuplicateFaqContentToAnotherLanguage(contentId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	var faqContent models.FaqContent
	if err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Components").
		Preload("MetaTag").
		First(&faqContent, "id = ?", contentId).Error; err != nil {
		return nil, err
	}

	// Change language
	if faqContent.Language == enums.PageLanguageTH {
		faqContent.Language = enums.PageLanguageEN
	} else {
		faqContent.Language = enums.PageLanguageTH
	}

	faqContent.ID = uuid.Nil

	// Metatag
	faqContent.MetaTagID = uuid.Nil
	if faqContent.MetaTag != nil {
		faqContent.MetaTag.ID = uuid.Nil
		faqContent.MetaTag.CreatedAt = time.Time{}
		faqContent.MetaTag.UpdatedAt = time.Time{}
	}

	// Content
	faqContent.ID = uuid.Nil
	faqContent.CreatedAt = time.Time{}
	faqContent.UpdatedAt = time.Time{}

	// Attach revision, always needs a revision for keeping the history (but this is a new content, so it's a start of the history)
	if newRevision == nil {
		return nil, errs.ErrNoRevisionFound
	}
	faqContent.Revision = newRevision

	// Components
	if faqContent.Components != nil {
		for _, component := range faqContent.Components {
			component.CreatedAt = time.Time{}
			component.UpdatedAt = time.Time{}
			component.ID = uuid.Nil
			component.FaqContentID = nil
		}
	}

	if err := r.db.Create(&faqContent).Error; err != nil {
		return nil, err
	}

	return &faqContent, nil
}

func (r *CMSFaqPageRepository) RevertFaqContent(revisionId uuid.UUID, newRevision *models.Revision) (*models.FaqContent, error) {
	var revision models.Revision
	// Get the revision
	if err := r.db.First(&revision, "id = ?", revisionId).Error; err != nil {
		return nil, err
	}

	faqContent := &models.FaqContent{}
	// Get the content of that revision
	if err := r.db.
		Preload("Revision").
		Preload("Categories").
		Preload("Components").
		Preload("MetaTag").
		First(faqContent, "id = ?", revision.FaqContentID).Error; err != nil {
		return nil, err
	}

	// Get the old content
	oldContent, err := r.FindLatestContentByPageId(faqContent.PageID, string(faqContent.Language))
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
		faqContent.Mode = enums.PageModeDraft

		// Metatag
		faqContent.MetaTagID = uuid.Nil
		if faqContent.MetaTag != nil {
			faqContent.MetaTag.ID = uuid.Nil
			faqContent.MetaTag.CreatedAt = time.Time{}
			faqContent.MetaTag.UpdatedAt = time.Time{}
		}

		// Content
		faqContent.ID = uuid.Nil
		faqContent.CreatedAt = time.Time{}
		faqContent.UpdatedAt = time.Time{}

		// Attach revision, always needs a revision for keeping the history (but this is a new content, so it's a start of the history)
		if newRevision == nil {
			return errs.ErrNoRevisionFound
		}
		faqContent.Revision = newRevision

		// Components
		if faqContent.Components != nil {
			for _, component := range faqContent.Components {
				component.CreatedAt = time.Time{}
				component.UpdatedAt = time.Time{}
				component.ID = uuid.Nil
				component.FaqContentID = nil
			}
		}

		if err := tx.Create(faqContent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return faqContent, nil
}

func (r *CMSFaqPageRepository) GetCategory(pageId uuid.UUID, categoryTypeCode, language, mode string) ([]models.Category, error) {
	var result struct {
		ID uuid.UUID
	}
	if err := r.db.
		Model(&models.FaqContent{}).
		Select("id").
		Where("page_id = ? AND language = ? AND mode = ?", pageId, language, mode).
		Take(&result).Error; err != nil {
		return nil, err
	}

	faqContentID := result.ID

	var categoryIDs []uuid.UUID
	if err := r.db.
		Model(&models.FaqContentCategory{}).
		Where("faq_content_id = ?", faqContentID).
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

func (r *CMSFaqPageRepository) GetRevisionByFaqPageId(pageId uuid.UUID, language string) ([]models.Revision, error) {
	var faqPage models.FaqPage
	err := r.db.
		Preload("Contents", "language = ?", language).
		Preload("Contents.Revision").
		Where("id = ?", pageId).
		First(&faqPage).Error

	if err != nil {
		return nil, err
	}

	var revisions []models.Revision
	for _, content := range faqPage.Contents {
		if content.Revision != nil {
			revisions = append(revisions, *content.Revision)
		}
	}

	// Sort by created_at, first revision the newest
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].CreatedAt.After(revisions[j].CreatedAt)
	})

	return revisions, nil
}

func (r *CMSFaqPageRepository) IsUrlDuplicate(url string, pageId uuid.UUID) (bool, error) {
	var count int64
	if pageId == uuid.Nil { // No page yet
		err := r.db.Model(&models.FaqContent{}).
			Where("url = ?", url).
			Count(&count).Error

		if err != nil {
			return false, err
		}

		return count > 0, nil
	}

	err := r.db.Model(&models.FaqContent{}).
		Where("url = ? AND page_id != ?", url, pageId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CMSFaqPageRepository) IsUrlAliasDuplicate(urlAlias string, pageId uuid.UUID) (bool, error) {
	var count int64
	if pageId == uuid.Nil { // No page yet
		err := r.db.Model(&models.FaqContent{}).
			Where("url_alias = ?", urlAlias).
			Count(&count).Error

		if err != nil {
			return false, err
		}

		return count > 0, nil
	}

	err := r.db.Model(&models.FaqContent{}).
		Where("url_alias = ? AND page_id != ?", urlAlias, pageId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *CMSFaqPageRepository) GetPageIdByContentId(contentId uuid.UUID) (uuid.UUID, error) {
	var faqContent models.FaqContent

	err := r.db.
		First(&faqContent, "id = ?", contentId).Error

	if err != nil {
		return uuid.Nil, err
	}

	return faqContent.PageID, nil
}

func (r *CMSFaqPageRepository) CreateFaqContentPreview(faqContentPreview *models.FaqContent) (*models.FaqContent, error) {
	faqContentPreview.Revision = nil
	faqContentPreview.Categories = nil
	if err := r.db.Create(faqContentPreview).Error; err != nil {
		return nil, err
	}

	return faqContentPreview, nil
}

func (r *CMSFaqPageRepository) UpdateFaqContentPreview(faqContentPreview *models.FaqContent) (*models.FaqContent, error) {
	// Overwrite full record (All fields)
	faqContentPreview.Revision = nil
	faqContentPreview.Categories = nil
	err := r.db.Transaction(func(tx *gorm.DB) error {

		// Delete the old components
		if err := tx.Where("faq_content_id = ?", faqContentPreview.ID).Delete(&models.Component{}).Error; err != nil {
			return err
		}

		if err := tx.Save(faqContentPreview).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return faqContentPreview, nil
}

func (r *CMSFaqPageRepository) FindFaqContentPreviewById(pageId uuid.UUID, language string) (*models.FaqContent, error) {
	var faqContentPreview models.FaqContent
	err := r.db.
		Where("page_id = ? AND language = ? AND mode = ?", pageId, string(language), "Preview").
		First(&faqContentPreview).Error

	if err != nil {
		return nil, err
	}

	return &faqContentPreview, nil
}
