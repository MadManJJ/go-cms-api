// In helpers/slugify.go (หรือไฟล์ที่เหมาะสม)
package helpers

import (
	"regexp"
	"strings"
	// "github.com/gosimple/unidecode" // Optional: for better non-ASCII character handling
)

var (
	// Regex to find non-alphanumeric characters (excluding hyphens if we want to keep them)
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9ก-ฮเ-ไ\-]+`)
	// Regex to replace multiple hyphens with a single one
	multipleHyphensRegex = regexp.MustCompile(`-{2,}`)
)

// GenerateSlug creates a URL-friendly slug from a string.
func GenerateSlug(text string) string {
	if text == "" {
		return ""
	}

	// Optional: Transliterate non-ASCII characters to ASCII (e.g., "你好" -> "Ni Hao")
	// text = unidecode.Unidecode(text) // Uncomment if using github.com/gosimple/unidecode

	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces and common separators with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove all non-alphanumeric characters except hyphens (if you allow thai chars, adjust regex)
	// For Thai, you might want a more sophisticated slug generation that preserves Thai words or transliterates.
	// This example is basic.
	slug = nonAlphanumericRegex.ReplaceAllString(slug, "") // This will remove Thai chars if regex is as above

	// Replace multiple hyphens with a single hyphen
	slug = multipleHyphensRegex.ReplaceAllString(slug, "-")

	// Trim leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
