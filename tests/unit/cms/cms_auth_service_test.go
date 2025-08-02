package tests

import (
	"testing"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/helpers"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockCMSAuthRepo struct {
	registerUser func(user *models.User) (*models.User, error)
	findUserByEmail func(email string) (*models.User, error)
	findUserById func(id uuid.UUID) (*models.User, error)
}

func (m *MockCMSAuthRepo) RegisterUser(user *models.User) (*models.User, error) {
	return m.registerUser(user)
}

func (m *MockCMSAuthRepo) FindUserByEmail(email string) (*models.User, error) {
	return m.findUserByEmail(email)
}

func (m *MockCMSAuthRepo) FindUserById(id uuid.UUID) (*models.User, error) {
	return m.findUserById(id)
}

func TestCMSService_RegisterUser(t *testing.T) {
	mockUser := helpers.InitializeMockUser()
	mockUserWithHashedPassword := helpers.InitializeMockUserWithHashedPassword()

	t.Run("successfully register user", func(t *testing.T) {
		repo := &MockCMSAuthRepo{
			registerUser: func(user *models.User) (*models.User, error) {
				return mockUserWithHashedPassword, nil
			},
		}

		service := services.NewCMSAuthService(repo)

		actualUser, err := service.RegisterUser(mockUser)
		assert.NoError(t, err)
		assert.Equal(t, mockUserWithHashedPassword, actualUser)
		assert.NotEqual(t, actualUser.Password, mockUser.Password)
	})

	t.Run("failed to register user", func(t *testing.T) {
		repo := &MockCMSAuthRepo{
			registerUser: func(user *models.User) (*models.User, error) {
				return nil, errs.ErrInternalServerError
			},
		}

		service := services.NewCMSAuthService(repo)

		actualUser, err := service.RegisterUser(mockUser)
		assert.Error(t, err)
		assert.Nil(t, actualUser)
	})	
}

func TestCMSService_LoginUser(t *testing.T) {
	mockUser := helpers.InitializeMockUser()
	mockUserWithHashedPassword := helpers.InitializeMockUserWithHashedPassword()

	t.Run("successfully login", func(t *testing.T) {
		repo := &MockCMSAuthRepo{
			findUserByEmail: func(email string) (*models.User, error) {
				return mockUserWithHashedPassword, nil
			},
		}

		service := services.NewCMSAuthService(repo)

		actualUser, token, err := service.LoginUser(mockUser)
		assert.NoError(t, err)
		assert.Equal(t, mockUserWithHashedPassword, actualUser)
		assert.NotNil(t, token)
	})

	t.Run("failed to login: incorrect password", func(t *testing.T) {
		repo := &MockCMSAuthRepo{
			findUserByEmail: func(email string) (*models.User, error) {
				return mockUserWithHashedPassword, nil
			},
		}

		service := services.NewCMSAuthService(repo)

		password := "incorrect_password"
		mockUser.Password = &password
		actualUser, token, err := service.LoginUser(mockUser)
		assert.Error(t, err)
		assert.Nil(t, actualUser)
		assert.Equal(t, "", token)
	})	
	
	t.Run("failed to login: not found", func(t *testing.T) {
		repo := &MockCMSAuthRepo{
			findUserByEmail: func(email string) (*models.User, error) {
				return nil, errs.ErrNotFound
			},
		}

		service := services.NewCMSAuthService(repo)

		actualUser, token, err := service.LoginUser(mockUser)
		assert.Error(t, err)
		assert.Nil(t, actualUser)
		assert.Equal(t, "", token)
	})	
}