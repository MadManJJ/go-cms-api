package tests

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelper_NormalizeLanguage(t *testing.T) {
	t.Run("successfully normalize language (EN)", func(t *testing.T) {
		language := "EN"
		normalizedLanguage, err := helpers.NormalizeLanguage(language)
		assert.NoError(t, err)
		assert.Equal(t, normalizedLanguage, "en")
	})

	t.Run("successfully normalize language (eN)", func(t *testing.T) {
		language := "eN"
		normalizedLanguage, err := helpers.NormalizeLanguage(language)
		assert.NoError(t, err)
		assert.Equal(t, normalizedLanguage, "en")
	})	

	t.Run("successfully normalize language (TH)", func(t *testing.T) {
		language := "TH"
		normalizedLanguage, err := helpers.NormalizeLanguage(language)
		assert.NoError(t, err)
		assert.Equal(t, normalizedLanguage, "th")
	})	

	t.Run("successfully normalize language (Th)", func(t *testing.T) {
		language := "Th"
		normalizedLanguage, err := helpers.NormalizeLanguage(language)
		assert.NoError(t, err)
		assert.Equal(t, normalizedLanguage, "th")
	})	

	t.Run("failed to normalize: unsupported language", func(t *testing.T) {
		language := "Fr"
		normalizedLanguage, err := helpers.NormalizeLanguage(language)
		assert.Error(t, err)
		assert.Equal(t, normalizedLanguage, "")
	})	
}

func TestHelper_NormalizeMode(t *testing.T) {
	t.Run("successfully normalize mode (pubLIshed)", func(t *testing.T) {
		mode := "pubLIshed"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Published")
	})	

	t.Run("successfully normalize mode (Published)", func(t *testing.T) {
		mode := "Published"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Published")
	})	

	t.Run("successfully normalize mode (preVIew)", func(t *testing.T) {
		mode := "preVIew"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Preview")
	})	

	t.Run("successfully normalize mode (Preview)", func(t *testing.T) {
		mode := "Preview"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Preview")
	})	
	
	t.Run("successfully normalize mode (hisToRies)", func(t *testing.T) {
		mode := "hisToRies"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Histories")
	})		

	t.Run("successfully normalize mode (Histories)", func(t *testing.T) {
		mode := "Histories"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Histories")
	})		

	t.Run("successfully normalize mode (DRAft)", func(t *testing.T) {
		mode := "DRAft"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Draft")
	})	

	t.Run("successfully normalize mode (Draft)", func(t *testing.T) {
		mode := "Draft"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.NoError(t, err)
		assert.Equal(t, normalizedMode, "Draft")
	})	

	t.Run("failed to normalize: unsupported mode", func(t *testing.T) {
		mode := "unknowMode"
		normalizedMode, err := helpers.NormalizeMode(mode)
		assert.Error(t, err)
		assert.Equal(t, normalizedMode, "")
	})
}

func TestHelper_NormalizeWorkflowStatus(t *testing.T) {
	t.Run("successfully normalize workflow status (draFT)", func(t *testing.T) {
		status := "draFT"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Draft")
	})

	t.Run("successfully normalize workflow status (Draft)", func(t *testing.T) {
		status := "Draft"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Draft")
	})

	t.Run("successfully normalize workflow status (approval_Pending)", func(t *testing.T) {
		status := "approval_Pending"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Approval_Pending")
	})

	t.Run("successfully normalize workflow status (Approval_Pending)", func(t *testing.T) {
		status := "Approval_Pending"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Approval_Pending")
	})

	t.Run("successfully normalize workflow status (waiting_Design_approved)", func(t *testing.T) {
		status := "waiting_Design_approved"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Waiting_Design_Approved")
	})

	t.Run("successfully normalize workflow status (Waiting_Design_Approved)", func(t *testing.T) {
		status := "Waiting_Design_Approved"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Waiting_Design_Approved")
	})

	t.Run("successfully normalize workflow status (schedule)", func(t *testing.T) {
		status := "schedule"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Schedule")
	})

	t.Run("successfully normalize workflow status (Schedule)", func(t *testing.T) {
		status := "Schedule"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Schedule")
	})

	t.Run("successfully normalize workflow status (published)", func(t *testing.T) {
		status := "published"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Published")
	})

	t.Run("successfully normalize workflow status (Published)", func(t *testing.T) {
		status := "Published"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Published")
	})

	t.Run("successfully normalize workflow status (unPublished)", func(t *testing.T) {
		status := "unPublished"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "UnPublished")
	})

	t.Run("successfully normalize workflow status (UnPublished)", func(t *testing.T) {
		status := "UnPublished"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "UnPublished")
	})

	t.Run("successfully normalize workflow status (waiting_Deletion)", func(t *testing.T) {
		status := "waiting_Deletion"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Waiting_Deletion")
	})

	t.Run("successfully normalize workflow status (Waiting_Deletion)", func(t *testing.T) {
		status := "Waiting_Deletion"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Waiting_Deletion")
	})

	t.Run("successfully normalize workflow status (delete)", func(t *testing.T) {
		status := "delete"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Delete")
	})

	t.Run("successfully normalize workflow status (Delete)", func(t *testing.T) {
		status := "Delete"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, normalizedStatus, "Delete")
	})

	t.Run("failed to normalize: unsupported workflow status", func(t *testing.T) {
		status := "unknownStatus"
		normalizedStatus, err := helpers.NormalizeWorkflowStatus(status)
		assert.Error(t, err)
		assert.Equal(t, normalizedStatus, "")
	})
}

func TestHelper_NormalizeFaqContent(t *testing.T) {
	t.Run("successfully normalize valid FaqContent fields", func(t *testing.T) {
		faqContent := &models.FaqContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("pubLIshed"),
			WorkflowStatus: enums.WorkflowStatus("draFT"),
			PublishStatus:  enums.PublishStatus("pubLIshed"),
		}

		err := helpers.NormalizeFaqContent(faqContent)
		assert.NoError(t, err)
		assert.Equal(t, enums.PageLanguageEN, faqContent.Language)
		assert.Equal(t, enums.PageModePublished, faqContent.Mode)
		assert.Equal(t, enums.WorkflowDraft, faqContent.WorkflowStatus)
		assert.Equal(t, enums.PublishStatusPublished, faqContent.PublishStatus)
	})

	t.Run("fails to normalize with invalid language", func(t *testing.T) {
		faqContent := &models.FaqContent{
			Language:       enums.PageLanguage("unknownLang"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizeFaqContent(faqContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid mode", func(t *testing.T) {
		faqContent := &models.FaqContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("invalidMode"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizeFaqContent(faqContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid workflow status", func(t *testing.T) {
		faqContent := &models.FaqContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("not_a_real_status"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizeFaqContent(faqContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid publish status", func(t *testing.T) {
		faqContent := &models.FaqContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("badStatus"),
		}

		err := helpers.NormalizeFaqContent(faqContent)
		assert.Error(t, err)
	})
}

func TestHelper_NormalizeRevision(t *testing.T) {
	t.Run("successfully normalize valid Revision publish status", func(t *testing.T) {
		revision := &models.Revision{
			PublishStatus: enums.PublishStatus("published"),
		}

		err := helpers.NormalizeRevision(revision)
		assert.NoError(t, err)
		assert.Equal(t, enums.PublishStatusPublished, revision.PublishStatus)
	})

	t.Run("fails to normalize with invalid publish status", func(t *testing.T) {
		revision := &models.Revision{
			PublishStatus: enums.PublishStatus("badStatus"),
		}

		err := helpers.NormalizeRevision(revision)
		assert.Error(t, err)
	})
}

func TestHelper_NormalizeCategory(t *testing.T) {
	t.Run("successfully normalize valid Category fields", func(t *testing.T) {
		category := &models.Category{
			PublishStatus: enums.PublishStatus("unPublished"),
			LanguageCode:  enums.PageLanguage("TH"),
		}

		err := helpers.NormalizeCategory(category)
		assert.NoError(t, err)
		assert.Equal(t, enums.PublishStatusNotPublished, category.PublishStatus)
		assert.Equal(t, enums.PageLanguageTH, category.LanguageCode)
	})

	t.Run("fails to normalize with invalid publish status", func(t *testing.T) {
		category := &models.Category{
			PublishStatus: enums.PublishStatus("invalid"),
			LanguageCode:  enums.PageLanguage("EN"),
		}

		err := helpers.NormalizeCategory(category)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid language code", func(t *testing.T) {
		category := &models.Category{
			PublishStatus: enums.PublishStatus("Published"),
			LanguageCode:  enums.PageLanguage("invalidLang"),
		}

		err := helpers.NormalizeCategory(category)
		assert.Error(t, err)
	})
}

func TestHelper_NormalizeLandingContent(t *testing.T) {
	t.Run("successfully normalize valid LandingContent fields", func(t *testing.T) {
		landingContent := &models.LandingContent{
			Language:       enums.PageLanguage("th"),
			Mode:           enums.PageMode("preView"),
			WorkflowStatus: enums.WorkflowStatus("approval_pending"),
			PublishStatus:  enums.PublishStatus("unPublished"),
		}

		err := helpers.NormalizeLandingContent(landingContent)
		assert.NoError(t, err)
		assert.Equal(t, enums.PageLanguageTH, landingContent.Language)
		assert.Equal(t, enums.PageModePreview, landingContent.Mode)
		assert.Equal(t, enums.WorkflowApprovalPending, landingContent.WorkflowStatus)
		assert.Equal(t, enums.PublishStatusNotPublished, landingContent.PublishStatus)
	})

	t.Run("fails to normalize with invalid language", func(t *testing.T) {
		landingContent := &models.LandingContent{
			Language:       enums.PageLanguage("xx"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizeLandingContent(landingContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid mode", func(t *testing.T) {
		landingContent := &models.LandingContent{
			Language:       enums.PageLanguage("en"),
			Mode:           enums.PageMode("invalidMode"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizeLandingContent(landingContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid workflow status", func(t *testing.T) {
		landingContent := &models.LandingContent{
			Language:       enums.PageLanguage("en"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("not_a_status"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizeLandingContent(landingContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid publish status", func(t *testing.T) {
		landingContent := &models.LandingContent{
			Language:       enums.PageLanguage("en"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("bad_status"),
		}

		err := helpers.NormalizeLandingContent(landingContent)
		assert.Error(t, err)
	})
}

func TestHelper_NormalizePartnerContent(t *testing.T) {
	t.Run("successfully normalize valid PartnerContent fields", func(t *testing.T) {
		partnerContent := &models.PartnerContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("preview"),
			WorkflowStatus: enums.WorkflowStatus("waiting_design_approved"),
			PublishStatus:  enums.PublishStatus("UnPublished"),
		}

		err := helpers.NormalizePartnerContent(partnerContent)
		assert.NoError(t, err)
		assert.Equal(t, enums.PageLanguageEN, partnerContent.Language)
		assert.Equal(t, enums.PageModePreview, partnerContent.Mode)
		assert.Equal(t, enums.WorkflowWaitingDesign, partnerContent.WorkflowStatus)
		assert.Equal(t, enums.PublishStatusNotPublished, partnerContent.PublishStatus)
	})

	t.Run("fails to normalize with invalid language", func(t *testing.T) {
		partnerContent := &models.PartnerContent{
			Language:       enums.PageLanguage("XX"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizePartnerContent(partnerContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid mode", func(t *testing.T) {
		partnerContent := &models.PartnerContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("InvalidMode"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizePartnerContent(partnerContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid workflow status", func(t *testing.T) {
		partnerContent := &models.PartnerContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("nope_status"),
			PublishStatus:  enums.PublishStatus("Published"),
		}

		err := helpers.NormalizePartnerContent(partnerContent)
		assert.Error(t, err)
	})

	t.Run("fails to normalize with invalid publish status", func(t *testing.T) {
		partnerContent := &models.PartnerContent{
			Language:       enums.PageLanguage("EN"),
			Mode:           enums.PageMode("Published"),
			WorkflowStatus: enums.WorkflowStatus("Draft"),
			PublishStatus:  enums.PublishStatus("bad_status"),
		}

		err := helpers.NormalizePartnerContent(partnerContent)
		assert.Error(t, err)
	})
}

func TestHelper_SanitizeFaqPage(t *testing.T) {
	faqPage := helpers.InitializeMockFaqPage()
	faqPageContent := faqPage.Contents[0]
	faqPageContent.Title = "<script>alert('xss')</script>"
	faqPageContent.MetaTag.Description = "<script>xss</script>"
	faqPageContent.URL = "<h1>URLLL</h1>"

	helpers.SanitizeFaqPage(faqPage)
	
	fmt.Println(faqPageContent.URLAlias)
	assert.Equal(t, "", faqPageContent.Title)
	assert.Equal(t, "", faqPageContent.MetaTag.Description)
	assert.Equal(t, "URLLL", faqPageContent.URL)
	assert.Equal(t, "<p>This is a mock FAQ content.</p>", faqPageContent.HTMLInput)
}

func TestHelper_TestUUIDFromSub(t *testing.T) {
	sub := "test-sub"

	uuid1 := helpers.UUIDFromSub(sub)
	uuid2 := helpers.UUIDFromSub(sub)
	
	assert.Equal(t, uuid1, uuid2)
}

func TestHelper_GetHashed(t *testing.T) {
	password := "test-password"

	hashed1 := helpers.GetHashed(password)
	hashed2 := helpers.GetHashed(password)
	
	assert.NotEqual(t, hashed1, password)
	assert.NotEqual(t, hashed2, password)
	assert.NotEqual(t, hashed1, hashed2)
}

func TestHelper_GetUserIDFromContext(t *testing.T) {
	app := fiber.New()
	userId := uuid.New()
	var actualUserId uuid.UUID
	var err error

	app.Get("/", func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"user_id": userId.String(),
		}
		c.Locals("user", claims)		

		actualUserId, err = helpers.GetUserIDFromContext(c)
		if err != nil {
			return err
		}
		return c.SendString("Hello, World!")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, userId, actualUserId)
}

func TestHelper_BuildPreviewURL(t *testing.T) {
	baseURL := "https://example.com"
	language := "en"
	urlPath := "test-path"
	id := uuid.New()

	expectedURL, err := url.Parse(baseURL)
	require.NoError(t, err)
	expectedURL.Path = path.Join(expectedURL.Path, "preview", language, urlPath)
	query := expectedURL.Query()
	query.Set("id", id.String())
	expectedURL.RawQuery = query.Encode()

	actualURL, err := helpers.BuildPreviewURL(baseURL, language, urlPath, id)
	require.NoError(t, err)
	assert.Equal(t, expectedURL.String(), actualURL)
}