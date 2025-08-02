package tests

import (
	"log"
	"testing"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/services"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type MockCMSAuthRepo struct {
	registerUser    func(user *models.User) (*models.User, error)
	findUserByEmail func(email string) (*models.User, error)
	findUserById    func(id uuid.UUID) (*models.User, error)
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

func TestCMSService_GetLoginLink(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("No .env file found or failed to load")
	}	
	
	t.Run("successfully get login link", func(t *testing.T) {
		cfg := config.New()
		repo := &MockCMSAuthRepo{}

		service := services.NewLineLoginService(cfg, repo)
		
		loginLink, err := service.GetLoginLink()
		assert.NoError(t, err)
		assert.NotNil(t, loginLink)
	})

	t.Run("failed to get login link: no client id", func(t *testing.T) {
		cfg := config.New()
		repo := &MockCMSAuthRepo{}
		
		cfg.Line.ClientId = ""

		service := services.NewLineLoginService(cfg, repo)
		
		loginLink, err := service.GetLoginLink()
		assert.Error(t, err)
		assert.Equal(t, "", loginLink)
	})
}