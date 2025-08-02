package integration

import (
	"regexp"
	"testing"
	"time"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"
	"github.com/MadManJJ/cms-api/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type manualMockFormEmailService struct {
	done        chan struct{}
	lastRequest dto.SendEmailRequest
	errToReturn error
}

func newManualMockEmailService() *manualMockFormEmailService {
	return &manualMockFormEmailService{
		done: make(chan struct{}),
	}
}

func (m *manualMockFormEmailService) SendEmail(req dto.SendEmailRequest) error {
	defer close(m.done)
	m.lastRequest = req
	return m.errToReturn
}

func setupSubmissionTest(t *testing.T) (
	*gorm.DB,
	sqlmock.Sqlmock,
	services.CMSFormSubmissionServiceInterface,
	*manualMockFormEmailService, // return manual mock
	func(),
) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	repo := repositories.NewFormSubmissionRepository(gormDB)
	mockEmailService := newManualMockEmailService()

	service := services.NewCMSFormSubmissionService(repo, mockEmailService)

	return gormDB, mock, service, mockEmailService, cleanup
}

func TestFormSubmissionLifecycle(t *testing.T) {
	_, mock, submissionService, mockEmailService, cleanup := setupSubmissionTest(t)
	defer cleanup()

	testFormID := uuid.New()
	var createdSubmissionID uuid.UUID

	t.Run("1_CreateFormSubmission_AndSendEmail_Success", func(t *testing.T) {
		t.Log("===> Start: 1_CreateFormSubmission_AndSendEmail_Success")
		submissionDataJSON := `{"feedback": "This is a test submission."}`
		submittedEmail := "test.user@example.com"
		submissionModel := &models.FormSubmission{
			SubmittedData:  datatypes.JSON(submissionDataJSON),
			SubmittedEmail: &submittedEmail,
		}

		tempSubmissionID := uuid.New()
		emailCatID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "form_submissions"`)).
			WithArgs(
				testFormID,
				submissionDataJSON,
				&submittedEmail,
				sqlmock.AnyArg(), // submitted_at
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempSubmissionID))
		mock.ExpectCommit()

		// Mock Preload("Form")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_submissions" WHERE id = $1 AND "form_submissions"."id" = $2 ORDER BY "form_submissions"."id" LIMIT $3`)).
			WithArgs(tempSubmissionID, tempSubmissionID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tempSubmissionID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE "forms"."id" = $1`)).
			WithArgs(testFormID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email_category_id"}).AddRow(testFormID, emailCatID))

		// Mock GetEmailContentsFormFormId
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "email_category_id" FROM "forms" WHERE id = $1 AND "forms"."deleted_at" IS NULL ORDER BY "forms"."id" LIMIT $2`)).
			WithArgs(testFormID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"email_category_id"}).AddRow(emailCatID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_contents" WHERE email_category_id = $1`)).
			WithArgs(emailCatID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label", "language", "email_category_id"}).AddRow(uuid.New(), "user-confirmation", "en", emailCatID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "email_categories" WHERE "email_categories"."id" = $1`)).
			WithArgs(emailCatID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(emailCatID, "Test Category"))

		// --- Act ---
		createdSubmission, err := submissionService.CreateFormSubmission(testFormID, submissionModel)

		// --- Assert ---
		require.NoError(t, err)
		require.NotNil(t, createdSubmission)
		createdSubmissionID = createdSubmission.ID

		select {
		case <-mockEmailService.done:
		case <-time.After(2 * time.Second):
			t.Fatal("Test timed out waiting for SendEmail to be called")
		}
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("2_GetFormSubmission_Success", func(t *testing.T) {
		t.Log("===> Start: 2_GetFormSubmission_Success")
		require.NotEqual(t, uuid.Nil, createdSubmissionID, "Submission must be created first")

		mock.ExpectQuery(`SELECT \* FROM "form_submissions" WHERE id = \$1 ORDER BY "form_submissions"\."id" LIMIT \$2`).
			WithArgs(createdSubmissionID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
				AddRow(createdSubmissionID, testFormID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE "forms"."id" = $1`)).
			WithArgs(testFormID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(testFormID, "Submission Test Form"))

		submission, err := submissionService.GetFormSubmission(createdSubmissionID)

		require.NoError(t, err)
		require.NotNil(t, submission)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("3_ListFormSubmissions_Success", func(t *testing.T) {
		t.Log("===> Start: 3_ListFormSubmissions_Success")
		require.NotEqual(t, uuid.Nil, testFormID, "Form ID must exist")
		require.NotEqual(t, uuid.Nil, createdSubmissionID, "Submission must exist to be listed")

		// Mock Count
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "form_submissions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "form_submissions" WHERE form_id = $1 ORDER BY created_at DESC LIMIT $2`)).
			WithArgs(testFormID, 10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "form_id"}).
				AddRow(createdSubmissionID, testFormID))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "forms" WHERE "forms"."id" = $1 AND "forms"."deleted_at" IS NULL`)).
			WithArgs(testFormID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(testFormID, "Submission Test Form"))

		submissions, totalCount, err := submissionService.GetFormSubmissions(testFormID, "created_at:desc", 1, 10)

		require.NoError(t, err)
		require.Len(t, submissions, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, createdSubmissionID, submissions[0].ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
