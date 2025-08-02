package tests

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCMSRepo_CreateForm(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)

	t.Run("successfully create form", func(t *testing.T) {
		mockForm := helpers.InitializeMockForm()
		formId := uuid.New()
		formSectionId := uuid.New()
		formFieldId := uuid.New()

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "forms"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(formId), 
			)

		mock.ExpectQuery(`INSERT INTO "form_sections"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "form_id"}).
					AddRow(formSectionId, formId), 
			)			

		mock.ExpectQuery(`INSERT INTO "form_fields"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "section_id"}).
					AddRow(formFieldId, formSectionId), 
			)				

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(formId))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
					AddRow(formSectionId, formId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).
					AddRow(formFieldId, formSectionId))					

		form, err := cmsFormRepo.CreateForm(gormDB, mockForm)

		assert.NoError(t, err)
		assert.NotNil(t, form)
		assert.NotNil(t, form.Sections)
		assert.NotNil(t, form.Sections[0].Fields)
		assert.Equal(t, formId, form.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed create form: internal server error", func(t *testing.T) {
		mockForm := helpers.InitializeMockForm()

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "forms"`).
			WillReturnError(errs.ErrInternalServerError)		

		mock.ExpectRollback()

		form, err := cmsFormRepo.CreateForm(gormDB, mockForm)

		assert.Error(t, err)
		assert.Nil(t, form)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_GetFormByID(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)

	t.Run("successfully get form by ID", func(t *testing.T) {
		formId := uuid.New()
		formSectionId := uuid.New()
		formFieldId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(formId))		

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
					AddRow(formSectionId, formId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).
					AddRow(formFieldId, formSectionId))					

		form, err := cmsFormRepo.GetFormByID(formId)

		assert.NoError(t, err)
		assert.NotNil(t, form)
		assert.NotNil(t, form.Sections)
		assert.NotNil(t, form.Sections[0].Fields)
		assert.Equal(t, formId, form.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed get form by ID: not found", func(t *testing.T) {
		formId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms"`)).
			WillReturnError(errs.ErrInternalServerError)

		form, err := cmsFormRepo.GetFormByID(formId)

		assert.Error(t, err)
		assert.Nil(t, form)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_GetFormStructure(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)
	
	t.Run("successfully get form structure", func(t *testing.T) {
		formId := uuid.New()
		formSectionId := uuid.New()
		formFieldId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(formId))	

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
					AddRow(formSectionId, formId))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).
					AddRow(formFieldId, formSectionId))					

		form, err := cmsFormRepo.GetFormStructure(formId)

		assert.NoError(t, err)
		assert.NotNil(t, form)
		assert.NotNil(t, form.Sections)
		assert.NotNil(t, form.Sections[0].Fields)
		assert.Equal(t, formId, form.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed get form structure: not found", func(t *testing.T) {
		formId := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms"`)).
			WillReturnError(errs.ErrInternalServerError)

		form, err := cmsFormRepo.GetFormStructure(formId)

		assert.Error(t, err)
		assert.Nil(t, form)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_ListForms(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)

	name := "test"
	createdBy := uuid.New()
	sort := "name_asc"
	page := 1
	limit := 10
	formFilter := dto.FormListFilter{
		Name: &name,
		Sort: &sort,
		Page: &page,
		ItemsPerPage: &limit,
	}	

	totalItems := int64(25)
		
	t.Run("successfully list forms", func(t *testing.T) {
		formId := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "forms"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(totalItems))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE name ILIKE $1 AND "forms"."deleted_at" IS NULL ORDER BY name ASC LIMIT $2`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_by_user_id"}).
				AddRow(formId, createdBy))				

		forms, actualTotalItems, err := cmsFormRepo.ListForms(formFilter)

		assert.NoError(t, err)
		assert.NotNil(t, forms)
		assert.Equal(t, totalItems, actualTotalItems)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to list forms", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "forms"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(totalItems))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE name ILIKE $1 AND "forms"."deleted_at" IS NULL ORDER BY name ASC LIMIT $2`)).
			WillReturnError(errs.ErrInternalServerError)				

		forms, actualTotalItems, err := cmsFormRepo.ListForms(formFilter)

		assert.Error(t, err)
		assert.Nil(t, forms)
		assert.Equal(t, int64(0), actualTotalItems)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		
}

func TestCMSRepo_UpdateForm(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)

	formId := uuid.New()
	form := helpers.InitializeMockForm()
	form.ID = formId

	updatedForm := *form
	updatedForm.Name = "updated name"
	updatedForm.Slug = "updated-slug"

	updateTitle := "updated section name"
	updatedFormSection := &form.Sections[0]
	updatedFormSection.ID = uuid.New()
	updatedFormSection.FormID = formId
	updatedFormSection.Title = &updateTitle
	updatedFormSection.Fields[0].ID = uuid.New()
	updatedFormSection.Fields[0].SectionID = updatedFormSection.ID
	updatedForm.Sections[0] = *updatedFormSection

	updatedForm.Sections = []models.FormSection{
		*updatedFormSection,	
	}	

	t.Run("successfully update form", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "form_sections" WHERE form_id = $1`)).
			WillReturnResult(sqlmock.NewResult(1, 1))	
			
		mock.ExpectCommit()			
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "forms" SET "description"=$1,"email_category_id"=$2,"language"=$3,"name"=$4,"slug"=$5,"updated_at"=NOW() WHERE id = $6 AND "forms"."deleted_at" IS NULL`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "form_sections"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(updatedFormSection.ID))
				
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "form_fields"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(updatedFormSection.Fields[0].ID))					

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(updatedForm.ID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
					AddRow(updatedFormSection.ID, updatedForm.ID))					

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).
					AddRow(updatedFormSection.Fields[0].ID, updatedFormSection.ID))		

		actualForm, err := cmsFormRepo.UpdateForm(gormDB, &updatedForm)

		assert.NoError(t, err)
		assert.Equal(t, updatedForm.ID, actualForm.ID)
		assert.NotNil(t, actualForm.Sections)
		assert.NotNil(t, actualForm.Sections[0].Fields)		
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to update form", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "form_sections" WHERE form_id = $1`)).
			WillReturnError(errs.ErrInternalServerError)
			
		mock.ExpectRollback()		

		actualForm, err := cmsFormRepo.UpdateForm(gormDB, &updatedForm)

		assert.Error(t, err)
		assert.Nil(t, actualForm)
		assert.NoError(t, mock.ExpectationsWereMet())
	})		
} 

func TestCMSRepo_DeleteForm(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)
	
	formId := uuid.New()

	t.Run("successfully delete form", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "forms" WHERE "forms"."id" = $1`)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()
		
		err := cmsFormRepo.DeleteForm(gormDB, formId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to delete form", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "forms" WHERE "forms"."id" = $1`)).
			WillReturnError(errs.ErrInternalServerError)
			
		mock.ExpectRollback()			
		err := cmsFormRepo.DeleteForm(gormDB, formId)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_CheckFieldKeyExistsInForm(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)
	
	formId := uuid.New()
	fieldKey := "test"
	excludeFieldId := uuid.New()	

	t.Run("successfully check field key exists in form", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "form_fields" JOIN form_sections ON form_sections.id = form_fields.section_id WHERE (form_sections.form_id = $1 AND form_fields.field_key = $2) AND form_fields.id != $3`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))					

		exists, err := cmsFormRepo.CheckFieldKeyExistsInForm(formId, fieldKey, &excludeFieldId)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to check field key exists in form", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "form_fields" JOIN form_sections ON form_sections.id = form_fields.section_id WHERE (form_sections.form_id = $1 AND form_fields.field_key = $2) AND form_fields.id != $3`)).
			WillReturnError(errs.ErrInternalServerError)				

		exists, err := cmsFormRepo.CheckFieldKeyExistsInForm(formId, fieldKey, &excludeFieldId)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_GetFormWithFields(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()
	
	cmsFormRepo := repo.NewFormRepository(gormDB)
	
	formId := uuid.New()	
	formSectionId := uuid.New()
	formFieldId := uuid.New()

	t.Run("successfully get form with fields", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(formId))	
				
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_sections" WHERE "form_sections"."form_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
					AddRow(formSectionId, formId))						
						
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_fields" WHERE "form_fields"."section_id" = $1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "section_id"}).
					AddRow(formFieldId, formSectionId))						

		form, err := cmsFormRepo.GetFormWithFields(formId)

		assert.NoError(t, err)
		assert.NotNil(t, form)
		assert.NotNil(t, form.Sections)
		assert.NotNil(t, form.Sections[0].Fields)
		assert.Equal(t, formId, form.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
	t.Run("failed to get form with fields", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE id = $1`)).
			WillReturnError(errs.ErrInternalServerError)						

		form, err := cmsFormRepo.GetFormWithFields(formId)

		assert.Error(t, err)
		assert.Nil(t, form)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}