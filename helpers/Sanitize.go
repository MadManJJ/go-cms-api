package helpers

import (
	"github.com/MadManJJ/cms-api/models"

	"github.com/microcosm-cc/bluemonday"
)

func SanitizeRevision(revision *models.Revision) {
	// Plain text
	textPolicy := bluemonday.StrictPolicy()

	revision.Author = textPolicy.Sanitize(revision.Author)
	revision.Message = textPolicy.Sanitize(revision.Message)
	revision.Description = textPolicy.Sanitize(revision.Description)
}

func SanitizeCategories(categories []*models.Category) {
	// Plain text
	textPolicy := bluemonday.StrictPolicy()

	for _, category := range categories {
		category.Name = textPolicy.Sanitize(category.Name)
		if category.Description != nil {
			*category.Description = textPolicy.Sanitize(*category.Description)
		}
	}
}

func SanitizeMetaTag(metaTag *models.MetaTag) {
	// Plain text
	textPolicy := bluemonday.StrictPolicy()

	metaTag.Title = textPolicy.Sanitize(metaTag.Title)
	metaTag.Description = textPolicy.Sanitize(metaTag.Description)
}

func SanitizeFaqContent(content *models.FaqContent) {
	// HTML content
	htmlPolicy := bluemonday.UGCPolicy()
	htmlPolicy.AllowAttrs("class").Globally()
	htmlPolicy.AllowAttrs("style").Globally()
	htmlPolicy.AllowAttrs("target").Globally()

	content.HTMLInput = htmlPolicy.Sanitize(content.HTMLInput)

	// Plain text
	textPolicy := bluemonday.StrictPolicy()

	content.Title = textPolicy.Sanitize(content.Title)
	content.URLAlias = textPolicy.Sanitize(content.URLAlias)
	content.URL = textPolicy.Sanitize(content.URL)

	if content.MetaTag != nil {
		SanitizeMetaTag(content.MetaTag)
	}

	if content.Revision != nil {
		SanitizeRevision(content.Revision)
	}

	if content.Categories != nil {
		SanitizeCategories(content.Categories)
	}
}

func SanitizeFaqPage(page *models.FaqPage) {
	if page.Contents != nil {
		for _, content := range page.Contents {
			SanitizeFaqContent(content)
		}
	}
}

func SanitizePartnerContent(content *models.PartnerContent) {
	// HTML content
	htmlPolicy := bluemonday.UGCPolicy()
	htmlPolicy.AllowAttrs("class").Globally()
	htmlPolicy.AllowAttrs("style").Globally()
	htmlPolicy.AllowAttrs("target").Globally()	

	content.HTMLInput = htmlPolicy.Sanitize(content.HTMLInput)
	content.CompanyDetail = htmlPolicy.Sanitize(content.CompanyDetail)
	content.LeadBody = htmlPolicy.Sanitize(content.LeadBody)
	content.Challenges = htmlPolicy.Sanitize(content.Challenges)
	content.Solutions = htmlPolicy.Sanitize(content.Solutions)
	content.Results = htmlPolicy.Sanitize(content.Results)

	// Plain text
	textPolicy := bluemonday.StrictPolicy()
	
	content.Title = textPolicy.Sanitize(content.Title)
	content.ThumbnailImage = textPolicy.Sanitize(content.ThumbnailImage)
	content.ThumbnailAltText = textPolicy.Sanitize(content.ThumbnailAltText)
	content.CompanyLogo = textPolicy.Sanitize(content.CompanyLogo)
	content.CompanyAltText = textPolicy.Sanitize(content.CompanyAltText)
	content.CompanyName = textPolicy.Sanitize(content.CompanyName)
	content.URLAlias = textPolicy.Sanitize(content.URLAlias)
	content.URL = textPolicy.Sanitize(content.URL)

	// Sanitize related models
	if content.MetaTag != nil {
		SanitizeMetaTag(content.MetaTag)
	}
	if content.Revision != nil {
		SanitizeRevision(content.Revision)
	}
	if content.Categories != nil {
		SanitizeCategories(content.Categories)
	}
}

func SanitizePartnerPage(page *models.PartnerPage) {
	if page.Contents != nil {
		for _, content := range page.Contents {
			SanitizePartnerContent(content)
		}
	}
}

func SanitizeLandingContent(content *models.LandingContent) {
	// HTML content
	htmlPolicy := bluemonday.UGCPolicy()
	htmlPolicy.AllowAttrs("class").Globally()
	htmlPolicy.AllowAttrs("style").Globally()
	htmlPolicy.AllowAttrs("target").Globally()	
	
	content.HTMLInput = htmlPolicy.Sanitize(content.HTMLInput)

	// Plain text
	textPolicy := bluemonday.StrictPolicy()
	
	content.Title = textPolicy.Sanitize(content.Title)
	content.UrlAlias = textPolicy.Sanitize(content.UrlAlias)

	// Sanitize related models
	if content.MetaTag != nil {
		SanitizeMetaTag(content.MetaTag)
	}
	if content.Revision != nil {
		SanitizeRevision(content.Revision)
	}
	if content.Categories != nil {
		SanitizeCategories(content.Categories)
	}
}

func SanitizeLandingPage(page *models.LandingPage) {
	if page.Contents != nil {
		for _, content := range page.Contents {
			SanitizeLandingContent(content)
		}
	}
}

func SanitizeEmailCategory(category *models.EmailCategory) {
	// Plain text
	textPolicy := bluemonday.StrictPolicy()

	category.Title = textPolicy.Sanitize(category.Title)
}

func SanitizeEmailContent(content *models.EmailContent) {
	// Plain text
	textPolicy := bluemonday.StrictPolicy()

	content.Label = textPolicy.Sanitize(content.Label)
	content.SendTo = textPolicy.Sanitize(content.SendTo)
	content.CcEmail = textPolicy.Sanitize(content.CcEmail)
	content.BccEmail = textPolicy.Sanitize(content.BccEmail)
	content.SendFromEmail = textPolicy.Sanitize(content.SendFromEmail)
	content.SendFromName = textPolicy.Sanitize(content.SendFromName)
	content.Subject = textPolicy.Sanitize(content.Subject)
	content.TopImgLink = textPolicy.Sanitize(content.TopImgLink)
	content.Header = textPolicy.Sanitize(content.Header)
	content.Paragraph = textPolicy.Sanitize(content.Paragraph)
	content.Footer = textPolicy.Sanitize(content.Footer)
	content.FooterImageLink = textPolicy.Sanitize(content.FooterImageLink)
}
