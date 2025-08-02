package repositories

import (
	"github.com/MadManJJ/cms-api/dto"
)

// CMSRepository defines the interface for CMS data access
type CMSRepository interface {
	GetTestData() (*dto.CMSData, error)
	GetAdditionalData() (*dto.CMSData, error)
}

