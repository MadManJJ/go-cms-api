package helpers

import (
	"fmt"
	"net/url"
	"path"

	"github.com/google/uuid"
)

func BuildPreviewURL(base, language, urlPath string, id uuid.UUID) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
	}

	// Join the path segments safely
	baseURL.Path = path.Join(baseURL.Path, "preview", language, urlPath)

	// Add query parameter ?id=UUID
	query := baseURL.Query()
	query.Set("id", id.String())
	baseURL.RawQuery = query.Encode()

	return baseURL.String(), nil
}