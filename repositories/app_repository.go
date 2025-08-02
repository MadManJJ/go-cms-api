package repositories

import "github.com/MadManJJ/cms-api/dto"

// AppRepository defines the interface for App data access
type AppRepository interface {
	GetTestData() (*dto.AppData, error)
	GetAdditionalData() (*dto.AppData, error)
}
