package enums

import "database/sql/driver"

type ProviderType string

const (
	ProviderNormal ProviderType = "normal"
	ProviderLine   ProviderType = "line"
)

// PageLanguage represents the language options.
type PageLanguage string

const (
	PageLanguageTH PageLanguage = "th"
	PageLanguageEN PageLanguage = "en"
)

// PageMode represents the publication mode of a page.
type PageMode string

const (
	PageModePublished PageMode = "Published"
	PageModePreview   PageMode = "Preview"
	PageModeHistories PageMode = "Histories"
	PageModeDraft     PageMode = "Draft"
)

type PublishStatus string

const (
	PublishStatusNotPublished PublishStatus = "UnPublished"
	PublishStatusPublished    PublishStatus = "Published"
)

// WorkflowStatus represents the workflow states a page can be in.
type WorkflowStatus string

const (
	WorkflowDraft           WorkflowStatus = "Draft"
	WorkflowApprovalPending WorkflowStatus = "Approval_Pending"
	WorkflowWaitingDesign   WorkflowStatus = "Waiting_Design_Approved"
	WorkflowSchedule        WorkflowStatus = "Schedule"
	WorkflowPublished       WorkflowStatus = "Published"
	WorkflowUnPublished     WorkflowStatus = "UnPublished"
	WorkflowWaitingDeletion WorkflowStatus = "Waiting_Deletion"
	WorkflowDelete          WorkflowStatus = "Delete"
)

// FileType represents the types of files.
type FileType string

const (
	FileTypeCSS FileType = "CSS"
	FileTypeJS  FileType = "JS"
)

// ComponentType represents all possible UI component types.
type ComponentType string

const (
	ComponentMargin                         ComponentType = "Margin"
	ComponentLeadComponent                  ComponentType = "LeadComponent"
	ComponentStickyHeader                   ComponentType = "StickyHeader"
	ComponentSectionContent                 ComponentType = "SectionContent"
	ComponentDivider                        ComponentType = "Divider"
	ComponentDynamicClassTextSection        ComponentType = "DynamicClassTextSection"
	ComponentVideoWithEditor                ComponentType = "VideoWithEditor"
	ComponentTC0101                         ComponentType = "TC0101"
	ComponentTC0102                         ComponentType = "TC0102"
	ComponentCoachProfileList               ComponentType = "CoachProfileList"
	ComponentX1004                          ComponentType = "X1004"
	ComponentH21                            ComponentType = "H21"
	ComponentH22                            ComponentType = "H22"
	ComponentH23                            ComponentType = "H23"
	ComponentH24                            ComponentType = "H24"
	ComponentH31                            ComponentType = "H31"
	ComponentH32                            ComponentType = "H32"
	ComponentLargeGreenLinkButton           ComponentType = "LargeGreenLinkButton"
	ComponentThreeLargeGreenLinkButton      ComponentType = "ThreeLargeGreenLinkButton"
	ComponentLargeWhiteLinkButton           ComponentType = "LargeWhiteLinkButton"
	ComponentMidsizeWhiteLinkButtonLeft     ComponentType = "MidsizeWhiteLinkButtonLeftAligned"
	ComponentMidsizeWhiteLinkButtonCentered ComponentType = "MidsizeWhiteLinkButtonCentered"
	ComponentMidsizeWhiteLinkButtonRight    ComponentType = "MidsizeWhiteLinkButtonRightAligned"
	ComponentList                           ComponentType = "List"
	ComponentBox                            ComponentType = "Box"
	ComponentQuotation                      ComponentType = "Quotation"
	ComponentRelatedLinks                   ComponentType = "RelatedLinks"
	ComponentRelatedArticles                ComponentType = "RelatedArticles"
	ComponentBL0501                         ComponentType = "BL0501"
	ComponentOneColumnImage                 ComponentType = "OneColumnImage"
	ComponentTwoColumnImage                 ComponentType = "TwoColumnImage"
	ComponentVideo                          ComponentType = "Video"
	ComponentVideoExternalLink              ComponentType = "VideoExternalLink"
	ComponentTabContent                     ComponentType = "TabContent"
	ComponentL0201                          ComponentType = "L0201"
	ComponentL0301                          ComponentType = "L0301"
	ComponentL0401                          ComponentType = "L0401"
	ComponentL0501                          ComponentType = "L0501"
	ComponentL0601                          ComponentType = "L0601"
	ComponentNormalText                     ComponentType = "NormalText"
	ComponentNormalTextRed                  ComponentType = "NormalTextRed"
	ComponentBold                           ComponentType = "Bold"
	ComponentTextCentered                   ComponentType = "TextCentered"
	ComponentChatter                        ComponentType = "Chatter"
	ComponentTextList                       ComponentType = "TextList"
	ComponentTextListNumber                 ComponentType = "TextListNumber"
	ComponentNotes                          ComponentType = "Notes"
	ComponentLinks                          ComponentType = "Links"
	ComponentLinksSeparateWindow            ComponentType = "LinksSeparateWindow"
	ComponentAnchorLink                     ComponentType = "AnchorLink"
	ComponentPdf                            ComponentType = "Pdf"
	ComponentU0201                          ComponentType = "U0201"
	ComponentQrCode                         ComponentType = "QrCode"
	ComponentX0201                          ComponentType = "X0201"
	ComponentX0301                          ComponentType = "X0301"
	ComponentX0302                          ComponentType = "X0302"
	ComponentX0401List                      ComponentType = "X0401List"
	ComponentX0501                          ComponentType = "X0501"
	ComponentX0601                          ComponentType = "X0601"
	ComponentX0701                          ComponentType = "X0701"
	ComponentX0801                          ComponentType = "X0801"
	ComponentX1001                          ComponentType = "X1001"
	ComponentX1002                          ComponentType = "X1002"
	ComponentX1101                          ComponentType = "X1101"
	ComponentX1201                          ComponentType = "X1201"
	ComponentX1301                          ComponentType = "X1301"
	ComponentImageOneColText                ComponentType = "ImageOneColText"
	ComponentImageTwoColText                ComponentType = "ImageTwoColText"
	ComponentImageThreeColText              ComponentType = "ImageThreeColText"
	ComponentImageFourColText               ComponentType = "ImageFourColText"
	ComponentLeftImageRightTextWrapped      ComponentType = "LeftImageRightTextWrapped"
	ComponentY0201                          ComponentType = "Y0201"
	ComponentY0202                          ComponentType = "Y0202"
	ComponentY0301                          ComponentType = "Y0301"
	ComponentY0302                          ComponentType = "Y0302"
	ComponentImageVideoCoverOneColText      ComponentType = "ImageVideoCoverOneColText"
	ComponentImageVideoCoverTwoColText      ComponentType = "ImageVideoCoverTwoColText"
	ComponentCampaignLfc                    ComponentType = "CampaignLfc"
	ComponentGridContents                   ComponentType = "GridContents"
	ComponentGridContentsLfcFilter          ComponentType = "GridContentsLfcFilter"
)

type FormFieldType string

const (
	FieldTypeText          FormFieldType = "text"
	FieldTypeEmail         FormFieldType = "email"
	FieldTypeNumber        FormFieldType = "number"
	FieldTypePassword      FormFieldType = "password"
	FieldTypeDate          FormFieldType = "date"
	FieldTypeCheckbox      FormFieldType = "checkbox"
	FieldTypeDropdown      FormFieldType = "dropdown"
	FieldTypeCheckboxGroup FormFieldType = "checkboxgroup"
	FieldTypeRadio         FormFieldType = "radio"
	FieldTypeRadioGroup    FormFieldType = "radiogroup"
	FieldTypeTextArea      FormFieldType = "textarea"
	FieldTypeTextList      FormFieldType = "textlist"
	FieldTypeFile          FormFieldType = "file"
)

// GORM Scanner and Valuer for FormFieldType (เพื่อให้ GORM รู้จัก custom type นี้)
func (fft *FormFieldType) Scan(value interface{}) error {
	*fft = FormFieldType(value.(string))
	return nil
}

func (fft FormFieldType) Value() (driver.Value, error) {
	return string(fft), nil
}
