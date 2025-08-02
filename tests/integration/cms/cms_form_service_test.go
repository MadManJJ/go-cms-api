package integration

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/repositories"
	"github.com/MadManJJ/cms-api/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// แก้ไข setupFormTest ให้ถูกต้องตามที่ service ต้องการ
func setupFormTest(t *testing.T) (
	sqlmock.Sqlmock,
	services.CMSFormServiceInterface,
	*config.Config,
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	formRepo := repositories.NewFormRepository(gormDB)
	emailCategoryRepo := repositories.NewEmailCategoryRepository(gormDB)
	testCfg := &config.Config{
		App: config.AppConfig{
			WebBaseURL: "https://example-frontend.com",
		},
	}
	service := services.NewCMSFormService(gormDB, formRepo, emailCategoryRepo, testCfg)
	return mock, service, testCfg, cleanup
}

func TestFormLifecycle(t *testing.T) {
	mock, formService, _, cleanup := setupFormTest(t)
	defer cleanup()

	var createdFormID uuid.UUID
	var emailCategoryID = uuid.New()

	t.Run("1_CreateForm_Success", func(t *testing.T) {
		t.Log("===> Start: 1_CreateForm_Success")
		description := "Please provide your contact information."
		emailCatIDStr := emailCategoryID.String()
		lang := "en"

		reqDTO := dto.CreateFormRequest{
			Name:            "Contact Us Form",
			Description:     &description,
			EmailCategoryID: &emailCatIDStr,
			Language:        &lang,
			Sections: []dto.FormSectionRequest{
				{
					Title:       helpers.Ptr("Personal Information"),
					Description: helpers.Ptr("Your personal details"),
					OrderIndex:  0,
					Fields: []dto.FormFieldRequest{
						{Label: "Full Name", FieldKey: "full_name", FieldType: "text", IsRequired: true, OrderIndex: 0},
						{Label: "Email Address", FieldKey: "email", FieldType: "email", IsRequired: true, OrderIndex: 1},
					},
				},
			},
		}

		tempFormID := uuid.New()
		tempSectionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE id = $1 ORDER BY "email_categories"."id" LIMIT $2`)).
			WithArgs(emailCategoryID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(emailCategoryID))

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "forms"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempFormID))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "form_sections"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempSectionID))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "form_fields"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1 AND "forms"."deleted_at" IS NULL ORDER BY "forms"."id" LIMIT $2`)).
			WithArgs(tempFormID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "email_category_id"}).AddRow(tempFormID, reqDTO.Name, helpers.GenerateSlug(reqDTO.Name), emailCategoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE "email_categories"."id" = $1`)).
			WithArgs(emailCategoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(emailCategoryID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1 ORDER BY form_sections.order_index ASC`)).
			WithArgs(tempFormID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).AddRow(tempSectionID, tempFormID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1 ORDER BY form_fields.order_index ASC`)).
			WithArgs(tempSectionID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).AddRow(uuid.New(), tempSectionID).AddRow(uuid.New(), tempSectionID))

		// 2.3 Commit Transaction
		mock.ExpectCommit()

		// --- Act ---
		resp, err := formService.CreateNewForm(reqDTO)

		// --- Assert ---
		require.NoError(t, err)
		require.NotNil(t, resp)
		createdFormID = resp.ID
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("2_GetFormDetails_Success", func(t *testing.T) {
		t.Log("===> Start: 2_GetFormDetails_Success")
		require.NotEqual(t, uuid.Nil, createdFormID, "Form must be created first")
		tempSectionID := uuid.New()

		// Mock getFormWithFullAssociations
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1 AND "forms"."deleted_at" IS NULL ORDER BY "forms"."id" LIMIT $2`)).
			WithArgs(createdFormID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email_category_id"}).AddRow(createdFormID, "Contact Us Form", nil))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1 ORDER BY form_sections.order_index ASC`)).
			WithArgs(createdFormID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).AddRow(tempSectionID, createdFormID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1 ORDER BY form_fields.order_index ASC`)).
			WithArgs(tempSectionID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).AddRow(uuid.New(), tempSectionID))

		resp, err := formService.GetFormDetails(createdFormID)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("3_UpdateForm_Success", func(t *testing.T) {
		t.Log("===> Start: 3_UpdateForm_Success")
		require.NotEqual(t, uuid.Nil, createdFormID, "Form must be created first")

		updatedDesc := "Please provide your contact information and message."
		updatedReqDTO := dto.UpdateFormRequest{
			Name:        "Updated Contact Form",
			Description: &updatedDesc,
			Sections: []dto.UpdateFormSectionRequest{
				{
					Title:      helpers.Ptr("Your Message Section"),
					OrderIndex: 1,
					Fields: []dto.UpdateFormFieldRequest{
						{Label: "Your Message", FieldKey: "message", FieldType: "textarea", IsRequired: true, OrderIndex: 1},
					},
				},
			},
		}

		tempExistingSectionID := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1 AND "forms"."deleted_at" IS NULL ORDER BY "forms"."id" LIMIT $2`)).
			WithArgs(createdFormID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "email_category_id"}).AddRow(createdFormID, "contact-us-form", nil))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1 ORDER BY form_sections.order_index ASC`)).
			WithArgs(createdFormID).WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).AddRow(tempExistingSectionID, createdFormID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1 ORDER BY form_fields.order_index ASC`)).
			WithArgs(tempExistingSectionID).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		// Mock Transaction
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "form_sections" WHERE form_id = $1`)).WithArgs(createdFormID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "forms" SET`)).WillReturnResult(sqlmock.NewResult(1, 1))
		tempNewSectionID := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "form_sections"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempNewSectionID))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "form_fields"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		// Mock การ Preload หลัง Transaction
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1 AND "forms"."deleted_at" IS NULL ORDER BY "forms"."id" LIMIT $2`)).
			WithArgs(createdFormID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email_category_id"}).AddRow(createdFormID, updatedReqDTO.Name, nil))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1 ORDER BY form_sections.order_index ASC`)).
			WithArgs(createdFormID).WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).AddRow(tempNewSectionID, createdFormID))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1 ORDER BY form_fields.order_index ASC`)).
			WithArgs(tempNewSectionID).WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).AddRow(uuid.New(), tempNewSectionID))

		mock.ExpectCommit()

		resp, err := formService.UpdateExistingForm(createdFormID, updatedReqDTO)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("4_DeleteForm_Success", func(t *testing.T) {
		t.Log("===> Start: 4_DeleteForm_Success")
		require.NotEqual(t, uuid.Nil, createdFormID, "Form must be created first")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "forms" SET "deleted_at"=$1 WHERE "forms"."id" = $2 AND "forms"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), createdFormID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := formService.DeleteExistingForm(createdFormID)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
