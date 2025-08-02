package services

import (
	"fmt"
	"os"
	"time"

	"github.com/MadManJJ/cms-api/errs"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type CMSAuthServiceInterface interface {
	RegisterUser(user *models.User) (*models.User, error)
	LoginUser(user *models.User) (*models.User, string, error)
}

// CMSService handles business logic for CMS domain
type CMSAuthService struct {
	repo repositories.CMSAuthRepositoryInterface
}

// NewCMSService creates a new instance of CMSService
func NewCMSAuthService(repo repositories.CMSAuthRepositoryInterface) *CMSAuthService {
	return &CMSAuthService{
		repo: repo,
	}
}

// RegisterUser: registers a new user to the postgres database
func (s *CMSAuthService) RegisterUser(user *models.User) (*models.User, error) {
	// Validate the user data
	var validate = validator.New()
	if err := validate.Struct(user); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Store the hashed password
	hashedPasswordString := string(hashedPassword)
	user.Password = &hashedPasswordString

	// Set provider to normal
	user.Provider = "normal"

	return s.repo.RegisterUser(user)
}

// LoginUser: find user then validate the password and returns a JWT token
func (s *CMSAuthService) LoginUser(user *models.User) (*models.User, string, error) {

	if user.Email == nil || user.Password == nil {
		return nil, "", errs.ErrBadRequest
	}

	// Find user by email
	selectedUser, err := s.repo.FindUserByEmail(*user.Email)
	if err != nil {
		return nil, "", errs.ErrInvalidCredentials
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(*selectedUser.Password), []byte(*user.Password))
	if err != nil {
		return nil, "", errs.ErrInvalidCredentials
	}	

	// Create JWT token
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = selectedUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	t, err := token.SignedString([]byte(jwtSecretKey))
  if err != nil {
    return nil, "", err
  }

	return  selectedUser, t, nil
}