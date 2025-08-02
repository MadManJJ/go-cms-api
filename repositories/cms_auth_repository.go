package repositories

import (
	"github.com/MadManJJ/cms-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CMSAuthRepositoryInterface interface {
	RegisterUser(user *models.User) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserById(id uuid.UUID) (*models.User, error)
}

// CMSAuthRepo is an implementation of CMSAuthRepository
type CMSAuthRepository struct{
	db *gorm.DB
}

// NewCMSAuthRepo creates a new instance of CMSAuthRepo
func NewCMSAuthCMSAuthRepository(db *gorm.DB) *CMSAuthRepository {
	return &CMSAuthRepository{db: db}
}

// RegisterUser implements the CMSAuthRepository interface
func (r *CMSAuthRepository) RegisterUser(user *models.User) (*models.User, error) {

	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// FindUserByEmail implements the CMSAuthRepository interface
func (r *CMSAuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	
	return &user, nil
}

func (r *CMSAuthRepository) FindUserById(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	
	return &user, nil
}
