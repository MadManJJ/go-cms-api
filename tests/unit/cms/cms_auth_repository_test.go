package tests

import (
	"regexp"
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	repo "github.com/MadManJJ/cms-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCMSRepo_RegisterUser(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsAuthRepo := repo.NewCMSAuthCMSAuthRepository(gormDB)

	mockUser := helpers.InitializeMockUserWithHashedPassword()

	t.Run("successfully register user", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "users"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(uuid.New()), // or any UUID
			)		

		mock.ExpectCommit()

		actualUser, err := cmsAuthRepo.RegisterUser(mockUser)
		assert.NoError(t, err)
		assert.Equal(t, mockUser, actualUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed to register user", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "users"`).
			WillReturnError(errs.ErrInternalServerError)	

		mock.ExpectRollback()

		actualUser, err := cmsAuthRepo.RegisterUser(mockUser)
		assert.Error(t, err)
		assert.Nil(t, actualUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindUserByEmail(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsAuthRepo := repo.NewCMSAuthCMSAuthRepository(gormDB)
	
	email := "user@example.com"
	mockUser := helpers.InitializeMockUserWithHashedPassword()
	userId := uuid.New()

	t.Run("successfully find user by email", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(userId, mockUser.Email, mockUser.Password),
			)		

		actualUser, err := cmsAuthRepo.FindUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, actualUser.Email)
		assert.Equal(t, mockUser.Password, actualUser.Password)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to find user by email", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnError(errs.ErrInternalServerError)		

		actualUser, err := cmsAuthRepo.FindUserByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, actualUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	
}

func TestCMSRepo_FindUserById(t *testing.T) {
	gormDB, mock, cleanup := helpers.SetupTestDB(t)
	defer cleanup()

	cmsAuthRepo := repo.NewCMSAuthCMSAuthRepository(gormDB)
	
	userId := uuid.New()
	mockUser := helpers.InitializeMockUserWithHashedPassword()

	t.Run("successfully find user by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(userId, mockUser.Email, mockUser.Password),
			) 		

		actualUser, err := cmsAuthRepo.FindUserById(userId)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, actualUser.Email)
		assert.Equal(t, mockUser.Password, actualUser.Password)
		assert.NoError(t, mock.ExpectationsWereMet())
	})	

	t.Run("failed to find user by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnError(errs.ErrInternalServerError) 		

		actualUser, err := cmsAuthRepo.FindUserById(userId)
		assert.Error(t, err)
		assert.Nil(t, actualUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}