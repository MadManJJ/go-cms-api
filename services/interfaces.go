package services

import (
	"github.com/MadManJJ/cms-api/dto"
)

// AppServiceInterface defines methods for the App service
type AppServiceInterface interface {
	GetTestData() (*dto.AppData, error)
	GetAdditionalData() (*dto.AppData, error)
}

// CMSServiceInterface defines methods for the CMS service
type CMSServiceInterface interface {
	GetTestData() (*dto.CMSData, error)
	GetAdditionalData() (*dto.CMSData, error)
}
