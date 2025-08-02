package helpers

import (
	"encoding/json"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

func NewMockCategoryType(code, name string, active bool, createdAt, updatedAt time.Time) *models.CategoryType {
	return &models.CategoryType{
		TypeCode:  code,
		Name:      name,
		IsActive:  active,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewMockCreateCategoryTypeRequest(TypeCode, Name string, IsActive bool) *dto.CreateCategoryTypeRequest {
	return &dto.CreateCategoryTypeRequest{
		TypeCode: TypeCode,
		Name:     &Name,
		IsActive: &IsActive,
	}
}

func NewMockUpdateCategoryTypeRequest(Name string, IsActive bool) *dto.UpdateCategoryTypeRequest {
	return &dto.UpdateCategoryTypeRequest{
		Name:     &Name,
		IsActive: &IsActive,
	}
}

func NewMockCategoryTypeResponse(id uuid.UUID, code, name string, active bool, childCount *int, createdAt, updatedAt time.Time) *dto.CategoryTypeResponse {
	idStr := ""
	if id != uuid.Nil {
		idStr = id.String()
	}

	var childrenCount *map[string]int
	if childCount != nil {
		cc := map[string]int{
			"th": *childCount,
		}
		childrenCount = &cc
	} else {
		empty := make(map[string]int)
		childrenCount = &empty
	}

	return &dto.CategoryTypeResponse{
		ID:            idStr,
		TypeCode:      code,
		Name:          &name,
		IsActive:      active,
		ChildrenCount: childrenCount,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func NewCategory(
	categoryType models.CategoryType,
	languageCode enums.PageLanguage,
	name string,
	description *string,
	weight int,
	publishStatus enums.PublishStatus,
	createdAt, updatedAt time.Time,
) *models.Category {
	return &models.Category{
		CategoryType:  &categoryType,
		LanguageCode:  languageCode,
		Name:          name,
		Description:   description,
		Weight:        weight,
		PublishStatus: publishStatus,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func InitializeMockCategory() (*models.CategoryType, *models.Category) {
	now := time.Now()

	categoryTypeCode := "CTG001"
	categoryTypeName := "category type name"
	active := true
	var mockCategoryType = NewMockCategoryType(
		categoryTypeCode,
		categoryTypeName,
		active,
		now,
		now,
	)

	categoryName := "category name"
	categoryDescription := "some desc"
	categoryWeight := 15
	var mockCategory = NewCategory(
		*mockCategoryType,
		enums.PageLanguageEN,
		categoryName,
		&categoryDescription,
		categoryWeight,
		enums.PublishStatusPublished,
		now,
		now,
	)

	return mockCategoryType, mockCategory
}

func InitializeMockFaqPage() *models.FaqPage {
	now := time.Now()

	_, category := InitializeMockCategory()
	component := InitializeMockComponent()

	return &models.FaqPage{
		CreatedAt: now,
		UpdatedAt: now,
		Contents: []*models.FaqContent{
			{
				Title:          "Mock FAQ Title",
				Language:       enums.PageLanguageEN,
				AuthoredAt:     now,
				HTMLInput:      "<p>This is a mock FAQ content.</p>",
				Mode:           enums.PageModeDraft,
				WorkflowStatus: enums.WorkflowDraft,
				PublishStatus:  enums.PublishStatusPublished,
				PublishOn:      now.Add(24 * time.Hour),
				UnpublishOn:    now.Add(48 * time.Hour),
				AuthoredOn:     now,
				URLAlias:       "mock-faq-title",
				URL:            "/faq/mock-faq-title",
				MetaTag: &models.MetaTag{
					Title:       "mock-title",
					Description: "mock-desc",
					CoverImage:  "some_img",
				},
				Revision: &models.Revision{
					PublishStatus: enums.PublishStatusPublished,
					Author:        "some author name",
					Message:       "some message",
					Description:   "some desc",
				},
				Categories: []*models.Category{
					category,
				},
				Components: []*models.Component{
					component,
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
}

func InitializeMockComponent() *models.Component {
	props, _ := json.Marshal(map[string]interface{}{
		"text": "Sample component text",
		"size": "large",
	})

	return &models.Component{
		LandingContentID: nil,
		LandingContent:   nil,
		PartnerContentID: nil,
		PartnerContent:   nil,
		FaqContentID:     nil,
		FaqContent:       nil,
		Type:             "text", // replace with an actual value from your enum
		Props:            datatypes.JSON(props),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func InitializeMockRevision() *models.Revision {
	return &models.Revision{
		ID:               uuid.New(),
		LandingContentID: nil,
		PartnerContentID: nil,
		FaqContentID:     nil,
		PublishStatus:    enums.PublishStatusPublished,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Author:           "John Doe",
		Message:          "Initial revision",
		Description:      "Sample revision description",
		LandingContent:   nil,
		PartnerContent:   nil,
		FaqContent:       nil,
	}
}

func InitializeMockPartnerPage() *models.PartnerPage {
	now := time.Now()

	_, category := InitializeMockCategory()
	component := InitializeMockComponent()

	return &models.PartnerPage{
		CreatedAt: now,
		UpdatedAt: now,
		Contents: []*models.PartnerContent{
			{
				Title:            "Mock Partner Title",
				ThumbnailImage:   "mock_thumbnail.jpg",
				ThumbnailAltText: "Mock Thumbnail Alt",
				CompanyLogo:      "mock_logo.jpg",
				CompanyAltText:   "Mock Logo Alt",
				Language:         enums.PageLanguageEN,
				CompanyName:      "Mock Company",
				CompanyDetail:    "Mock Company Detail",
				LeadBody:         "Mock lead body content.",
				Challenges:       "Mock challenges description.",
				Solutions:        "Mock solutions description.",
				Results:          "Mock results content.",
				AuthoredAt:       now,
				HTMLInput:        "<p>This is mock partner content.</p>",
				Mode:             enums.PageModeDraft,
				WorkflowStatus:   enums.WorkflowDraft,
				PublishStatus:    enums.PublishStatusPublished,
				URLAlias:         "/partners/mock-partner-title",
				URL:              "/partners/mock-partner-title",
				MetaTag: &models.MetaTag{
					Title:       "mock-partner-meta-title",
					Description: "mock-partner-meta-description",
					CoverImage:  "mock-partner-cover.jpg",
				},
				IsRecommended: true,
				PublishOn:     now.Add(24 * time.Hour),
				UnpublishOn:   now.Add(48 * time.Hour),
				AuthoredOn:    now,
				CreatedAt:     now,
				UpdatedAt:     now,
				ApprovalEmail: pq.StringArray{"approver1@example.com", "approver2@example.com"},
				Revision: &models.Revision{
					PublishStatus: enums.PublishStatusPublished,
					Author:        "partner author name",
					Message:       "partner revision message",
					Description:   "partner revision description",
				},
				Categories: []*models.Category{
					category,
				},
				Components: []*models.Component{
					component,
				},
			},
		},
	}
}

func InitializeMockLandingPage() *models.LandingPage {
	now := time.Now()

	_, category := InitializeMockCategory()
	component := InitializeMockComponent()

	return &models.LandingPage{
		CreatedAt: now,
		UpdatedAt: now,
		Contents: []*models.LandingContent{
			{
				Title:          "Mock Landing Title",
				Language:       enums.PageLanguageEN,
				AuthoredAt:     now,
				HTMLInput:      "<p>This is mock landing page content.</p>",
				Mode:           enums.PageModeDraft,
				WorkflowStatus: enums.WorkflowDraft,
				PublishStatus:  enums.PublishStatusPublished,
				UrlAlias:       "mock-landing-title",
				MetaTag: &models.MetaTag{
					Title:       "mock-landing-meta-title",
					Description: "mock-landing-meta-description",
					CoverImage:  "mock-landing-cover.jpg",
				},
				PublishOn:     ptrTime(now.Add(24 * time.Hour)),
				UnpublishOn:   ptrTime(now.Add(48 * time.Hour)),
				AuthoredOn:    ptrTime(now),
				CreatedAt:     now,
				UpdatedAt:     now,
				ApprovalEmail: pq.StringArray{"approver1@example.com", "approver2@example.com"},
				Revision: &models.Revision{
					PublishStatus: enums.PublishStatusPublished,
					Author:        "landing author name",
					Message:       "landing revision message",
					Description:   "landing revision description",
				},
				Categories: []*models.Category{
					category,
				},
				Components: []*models.Component{
					component,
				},
			},
		},
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func InitializeMockEmailCategory() *models.EmailCategory {
	return &models.EmailCategory{
		Title:     "Welcome Emails",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func InitializeMockEmailContent() *models.EmailContent {
	emailCategory := InitializeMockEmailCategory()

	return &models.EmailContent{
		Language:        enums.PageLanguageEN,
		Label:           "welcome_new_user",
		SendTo:          "user@example.com",
		CcEmail:         "cc@example.com",
		BccEmail:        "bcc@example.com",
		SendFromEmail:   "noreply@example.com",
		SendFromName:    "Support Team",
		Subject:         "Welcome to Our Platform!",
		TopImgLink:      "https://example.com/images/welcome.png",
		Header:          "Hello!",
		Paragraph:       "Thank you for joining our platform. We’re excited to have you!",
		Footer:          "Best regards,\nThe Team",
		FooterImageLink: "https://example.com/images/footer.png",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		EmailCategory:   emailCategory,
	}
}

func InitializeMockUserWithHashedPassword() *models.User {
	user := InitializeMockUser()
	hashedPassword := string(GetHashed(*user.Password))
	email := "user@example.com"
	return &models.User{
		Email:    &email,
		Password: &hashedPassword,
		Provider: "normal",
	}
}

func InitializeMockUser() *models.User {
	email := "user@example.com"
	password := "12345678"
	return &models.User{
		Email:    &email,
		Password: &password,
	}
}

func InitializeMockFormSubmission() *models.FormSubmission {
	data := map[string]interface{}{
		"field1": "value1",
		"field2": 123,
	}
	jsonData, _ := json.Marshal(data)

	return &models.FormSubmission{
		SubmittedData: datatypes.JSON(jsonData),
	}
}

func InitializeMockForm() *models.Form {
	title := "Section 1"
	label := "Field 1"
	fieldKey := "field1"
	fieldType := enums.FieldTypeText // Assuming this is defined in your enums
	now := time.Now()

	return &models.Form{
		ID:        uuid.New(),
		Name:      "Mock Form",
		Slug:      "mock-form",
		CreatedAt: now,
		UpdatedAt: now,
		Sections: []models.FormSection{
			{
				ID:         uuid.New(),
				Title:      &title,
				OrderIndex: 0,
				CreatedAt:  now,
				UpdatedAt:  now,
				Fields: []models.FormField{
					{
						ID:         uuid.New(),
						Label:      label,
						FieldKey:   fieldKey,
						FieldType:  fieldType,
						IsRequired: true,
						OrderIndex: 0,
						CreatedAt:  now,
						UpdatedAt:  now,
						Properties: datatypes.JSON([]byte(`{"maxLength": 100}`)),
						Display:    datatypes.JSON([]byte(`{"width": "full"}`)),
					},
				},
			},
		},
	}
}

func InitializeMockMediaFile() *models.MediaFile {
	return &models.MediaFile{
		Name:        "Mock Media File",
		DownloadURL: "mock-media-file.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
