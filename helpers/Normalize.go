package helpers

import (
	"strings"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/models/enums"
)

func NormalizeLanguage(language string) (string, error) {
	switch {
	case strings.EqualFold(language, string(enums.PageLanguageEN)):
		return string(enums.PageLanguageEN), nil
	case strings.EqualFold(language, string(enums.PageLanguageTH)):
		return string(enums.PageLanguageTH), nil
	default:
		return "", errs.ErrInvalidLanguageCode
	}
}

func NormalizeMode(mode string) (string, error) {
	switch {
	case strings.EqualFold(mode, string(enums.PageModePublished)):
		return string(enums.PageModePublished), nil
	case strings.EqualFold(mode, string(enums.PageModePreview)):
		return string(enums.PageModePreview), nil
	case strings.EqualFold(mode, string(enums.PageModeHistories)):
		return string(enums.PageModeHistories), nil
	case strings.EqualFold(mode, string(enums.PageModeDraft)):
		return string(enums.PageModeDraft), nil
	default:
		return "", errs.ErrInvalidMode
	}
}

func NormalizeWorkflowStatus(workFlowStatus string) (string, error) {
	switch {
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowDraft)):
		return string(enums.WorkflowDraft), nil
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowApprovalPending)):
		return string(enums.WorkflowApprovalPending), nil
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowWaitingDesign)):
		return string(enums.WorkflowWaitingDesign), nil
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowSchedule)):
		return string(enums.WorkflowSchedule), nil
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowPublished)):
		return string(enums.WorkflowPublished), nil
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowUnPublished)):
		return string(enums.WorkflowUnPublished), nil
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowWaitingDeletion)):
		return string(enums.WorkflowWaitingDeletion), nil
	case strings.EqualFold(workFlowStatus, string(enums.WorkflowDelete)):
		return string(enums.WorkflowDelete), nil
	default:
		return "", errs.ErrInvalidWorkflowStatus
	}
}

func NormalizePublishStatus(publishStatus string) (string, error) {
	switch {
	case strings.EqualFold(publishStatus, string(enums.PublishStatusPublished)):
		return string(enums.PublishStatusPublished), nil
	case strings.EqualFold(publishStatus, string(enums.PublishStatusNotPublished)):
		return string(enums.PublishStatusNotPublished), nil
	default:
		return "", errs.ErrInvalidPublishStatus
	}
}

func NormalizeFaqContent(faqContent *models.FaqContent) error {
	language := string(faqContent.Language)
	mode := string(faqContent.Mode)
	workFlowStatus := string(faqContent.WorkflowStatus)
	publishStatus := string(faqContent.PublishStatus)

	language, err := NormalizeLanguage(language)
	if err != nil {
		return err
	}

	mode, err = NormalizeMode(mode)
	if err != nil {
		return err
	}

	workFlowStatus, err = NormalizeWorkflowStatus(workFlowStatus)
	if err != nil {
		return err
	}

	publishStatus, err = NormalizePublishStatus(publishStatus)
	if err != nil {
		return err
	}

	faqContent.Language = enums.PageLanguage(language)
	faqContent.Mode = enums.PageMode(mode)
	faqContent.WorkflowStatus = enums.WorkflowStatus(workFlowStatus)
	faqContent.PublishStatus = enums.PublishStatus(publishStatus)

	return nil
}

func NormalizeRevision(revision *models.Revision) error {
	publishStatus := string(revision.PublishStatus)

	publishStatus, err := NormalizePublishStatus(publishStatus)
	if err != nil {
		return err
	}

	revision.PublishStatus = enums.PublishStatus(publishStatus)

	return nil
}

func NormalizeCategory(category *models.Category) error {
	publishStatus := string(category.PublishStatus)
	language := string(category.LanguageCode)

	publishStatus, err := NormalizePublishStatus(publishStatus)
	if err != nil {
		return err
	}

	language, err = NormalizeLanguage(language)
	if err != nil {
		return err
	}

	category.PublishStatus = enums.PublishStatus(publishStatus)
	category.LanguageCode = enums.PageLanguage(language)

	return nil
}

func NormalizeLandingContent(landingContent *models.LandingContent) error {

	language := string(landingContent.Language)
	mode := string(landingContent.Mode)
	workFlowStatus := string(landingContent.WorkflowStatus)
	publishStatus := string(landingContent.PublishStatus)

	// Normalize each field
	var err error
	language, err = NormalizeLanguage(language)
	if err != nil {
		return err
	}

	mode, err = NormalizeMode(mode)
	if err != nil {
		return err
	}

	workFlowStatus, err = NormalizeWorkflowStatus(workFlowStatus)
	if err != nil {
		return err
	}

	publishStatus, err = NormalizePublishStatus(publishStatus)
	if err != nil {
		return err
	}

	// Convert normalized strings back to enums
	landingContent.Language = enums.PageLanguage(language)
	landingContent.Mode = enums.PageMode(mode)
	landingContent.WorkflowStatus = enums.WorkflowStatus(workFlowStatus)
	landingContent.PublishStatus = enums.PublishStatus(publishStatus)

	return nil
}

func NormalizePartnerContent(partnerContent *models.PartnerContent) error {

	language := string(partnerContent.Language)
	mode := string(partnerContent.Mode)
	workFlowStatus := string(partnerContent.WorkflowStatus)
	publishStatus := string(partnerContent.PublishStatus)

	// Normalize each field
	var err error
	language, err = NormalizeLanguage(language)
	if err != nil {
		return err
	}

	mode, err = NormalizeMode(mode)
	if err != nil {
		return err
	}

	workFlowStatus, err = NormalizeWorkflowStatus(workFlowStatus)
	if err != nil {
		return err
	}

	publishStatus, err = NormalizePublishStatus(publishStatus)
	if err != nil {
		return err
	}

	// Convert normalized strings back to enums
	partnerContent.Language = enums.PageLanguage(language)
	partnerContent.Mode = enums.PageMode(mode)
	partnerContent.WorkflowStatus = enums.WorkflowStatus(workFlowStatus)
	partnerContent.PublishStatus = enums.PublishStatus(publishStatus)

	return nil
}
