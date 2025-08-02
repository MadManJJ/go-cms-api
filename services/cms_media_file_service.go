package services

import (
	"errors"
	"fmt"
	"log"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MadManJJ/cms-api/config"
	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/models"
	"github.com/MadManJJ/cms-api/repositories"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type MediaFileServiceInterface interface {
	UploadMediaFile(originalFilename string, mimeType string, fileData []byte, customPath *string, replace *bool, userID uuid.UUID) (*dto.MediaFileResponse, error)
	GetMediaFileByID(idStr string) (*dto.MediaFileResponse, error)
	ListMediaFiles(filter dto.MediaFileListFilter) (*dto.MediaFilesListResponse, error)
	DeleteMediaFile(idStr string, userID uuid.UUID) error
}

type mediaFileService struct {
	cfg  *config.Config
	repo repositories.MediaFileRepositoryInterface
}

func NewMediaFileService(cfg *config.Config, repo repositories.MediaFileRepositoryInterface) MediaFileServiceInterface {

	return &mediaFileService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *mediaFileService) mapModelToResponse(mf *models.MediaFile) *dto.MediaFileResponse {
	if mf == nil {
		return nil
	}
	return &dto.MediaFileResponse{
		ID:          mf.ID.String(),
		Name:        mf.Name,
		DownloadURL: mf.DownloadURL,
		CreatedAt:   mf.CreatedAt,
		UpdatedAt:   mf.UpdatedAt,
	}
}

func (s *mediaFileService) mapModelToListItemResponse(mf *models.MediaFile) dto.MediaFileListItemResponse {
	return dto.MediaFileListItemResponse{
		ID:          mf.ID.String(),
		Name:        mf.Name,
		DownloadURL: mf.DownloadURL,
		CreatedAt:   mf.CreatedAt,
	}
}
func (s *mediaFileService) UploadMediaFile(
	originalFilename string,
	mimeType string, // นี่คือ MIME type ที่อาจจะมีพารามิเตอร์ เช่น "text/css; charset=utf-8"
	fileData []byte,
	customPathOpt *string,
	replaceOpt *bool,
	userID uuid.UUID,
) (*dto.MediaFileResponse, error) {

	// --- File Type Validation ---
	allowedMimeTypes := map[string]bool{
		"image/jpeg":             true,
		"image/png":              true,
		"image/webp":             true,
		"image/svg+xml":          true,
		"application/pdf":        true,
		"application/json":       true,
		"text/css":               true,
		"application/javascript": true,
		"text/javascript":        true,
	}

	
	mediaType, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		// Handle error if mimeType string is malformed, though unlikely if coming from standard detection
		log.Printf("Error parsing MIME type '%s': %v", mimeType, err)
		// Construct allowedTypesStr for error message (เหมือนเดิม)
		var allowedTypesStr []string
		for mt := range allowedMimeTypes {
			switch mt {
			case "image/jpeg":
				allowedTypesStr = append(allowedTypesStr, "JPEG")
			case "image/png":
				allowedTypesStr = append(allowedTypesStr, "PNG")
			case "image/webp":
				allowedTypesStr = append(allowedTypesStr, "WebP")
			case "image/svg+xml":
				allowedTypesStr = append(allowedTypesStr, "SVG")
			case "application/pdf":
				allowedTypesStr = append(allowedTypesStr, "PDF")
			case "application/json":
				allowedTypesStr = append(allowedTypesStr, "JSON")
			case "text/css":
				allowedTypesStr = append(allowedTypesStr, "CSS")
			case "application/javascript", "text/javascript":
				foundJS := false
				for _, ats := range allowedTypesStr {
					if ats == "JavaScript" {
						foundJS = true
						break
					}
				}
				if !foundJS {
					allowedTypesStr = append(allowedTypesStr, "JavaScript")
				}
			default:
				allowedTypesStr = append(allowedTypesStr, mt)
			}
		}
		return nil, fmt.Errorf("invalid or unparseable MIME type: %s. Allowed base types are %s", mimeType, strings.Join(allowedTypesStr, ", "))
	}

	// Construct the allowed types string for the error message dynamically (เหมือนเดิม)
	var allowedTypesStr []string
	for mt := range allowedMimeTypes {
		switch mt {
		case "image/jpeg":
			allowedTypesStr = append(allowedTypesStr, "JPEG")
		case "image/png":
			allowedTypesStr = append(allowedTypesStr, "PNG")
		case "image/webp":
			allowedTypesStr = append(allowedTypesStr, "WebP")
		case "image/svg+xml":
			allowedTypesStr = append(allowedTypesStr, "SVG")
		case "application/pdf":
			allowedTypesStr = append(allowedTypesStr, "PDF")
		case "application/json":
			allowedTypesStr = append(allowedTypesStr, "JSON")
		case "text/css":
			allowedTypesStr = append(allowedTypesStr, "CSS")
		case "application/javascript", "text/javascript":
			foundJS := false
			for _, ats := range allowedTypesStr {
				if ats == "JavaScript" {
					foundJS = true
					break
				}
			}
			if !foundJS {
				allowedTypesStr = append(allowedTypesStr, "JavaScript")
			}
		default:
			allowedTypesStr = append(allowedTypesStr, mt)
		}
	}

	// Now check the parsed mediaType (base type) against the allowedMimeTypes
	if !allowedMimeTypes[mediaType] { // ใช้ mediaType ที่ parse แล้ว
		return nil, fmt.Errorf("invalid file type: base type '%s' (from '%s') is not allowed. Allowed base types are %s", mediaType, mimeType, strings.Join(allowedTypesStr, ", "))
	}

	// --- SVG Sanitization (Example with basic string replace, consider a robust library for production) ---
	if mimeType == "image/svg+xml" {
		// IMPORTANT: This is a VERY basic sanitizer. For production, use a library like bluemonday
		// configured to allow only safe SVG elements and attributes.
		svgString := string(fileData)
		// Remove script tags (very basic)
		reScript := regexp.MustCompile(`(?i)<script.*?>.*?</script>`)
		svgString = reScript.ReplaceAllString(svgString, "")
		// Remove on* event handlers (very basic)
		reOnEvent := regexp.MustCompile(`(?i)\s(on\w+)=["'].*?["']`)
		svgString = reOnEvent.ReplaceAllString(svgString, "")
		fileData = []byte(svgString)
		// log.Printf("Sanitized SVG: %s", string(fileData)) // For debugging
	}

	// --- Path and Filename Logic ---
	uploadRoot := s.cfg.App.UploadPath // e.g., "./uploads"
	actualSubPath := ""
	if customPathOpt != nil {
		actualSubPath = filepath.Clean(*customPathOpt) // Clean the path
		if strings.HasPrefix(actualSubPath, "..") {    // Prevent directory traversal
			return nil, errors.New("invalid custom path")
		}
	}

	targetDir := filepath.Join(uploadRoot, actualSubPath)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory '%s': %w", targetDir, err)
	}

	filenameOnly := filepath.Base(originalFilename)
	ext := filepath.Ext(filenameOnly)
	base := strings.TrimSuffix(filenameOnly, ext)
	finalFilename := filenameOnly
	shouldReplace := false
	if replaceOpt != nil {
		shouldReplace = *replaceOpt
	}

	// Check for existing file by Name (and Path if your model supports it)
	// For now, assuming name must be unique globally or we generate unique names
	// If your model adds a 'Path' field, adjust repo.FindByNameAndPath
	existingFile, _ := s.repo.FindByNameAndPath(finalFilename, actualSubPath)

	if existingFile != nil {
		if shouldReplace {
			// Delete old file from disk and DB before saving new one
			oldFilePath := filepath.Join(uploadRoot, actualSubPath, existingFile.Name)
			if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
				log.Printf("Warning: failed to remove old file from disk '%s': %v", oldFilePath, err)
				// Decide if this is a critical error or if we can proceed
			}
			if err := s.repo.Delete(existingFile.ID); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Warning: failed to delete old file record from DB '%s': %v", existingFile.ID, err)
			}
		} else {
			// Generate a new unique name
			counter := 1
			for {
				finalFilename = fmt.Sprintf("%s___%d%s", base, counter, ext)
				checkExisting, _ := s.repo.FindByNameAndPath(finalFilename, actualSubPath)
				if checkExisting == nil {
					break
				}
				counter++
				if counter > 100 { // Safety break
					return nil, errors.New("could not generate a unique filename after 100 attempts")
				}
			}
		}
	}

	// --- Save File to Disk ---
	fullDiskPath := filepath.Join(targetDir, finalFilename)
	if err := os.WriteFile(fullDiskPath, fileData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file to disk '%s': %w", fullDiskPath, err)
	}

	// --- Create Database Record ---
	// Construct DownloadURL based on your API's serving mechanism
	// Example: http://localhost:8080/api/v1/files/subpath/image.jpg
	// Ensure cfg.Server.APIBaseURL ends WITHOUT a trailing slash
	// Ensure cfg.Server.StaticFilePrefix starts WITH a slash
	var urlPathBuilder strings.Builder
	urlPathBuilder.WriteString(s.cfg.App.StaticFilePrefix) // StaticFilePrefix ควรมี / นำหน้า เช่น "/files"

	if actualSubPath != "" { // actualSubPath คือ sub-directory ที่ clean แล้ว
		urlPathBuilder.WriteString("/")
		// Ensure actualSubPath uses forward slashes for URL
		urlPathBuilder.WriteString(strings.ReplaceAll(filepath.ToSlash(actualSubPath), "\\", "/"))
	}
	urlPathBuilder.WriteString("/")
	urlPathBuilder.WriteString(url.PathEscape(finalFilename)) // Escape filename for URL

	downloadURL := s.cfg.App.APIBaseURL + urlPathBuilder.String()

	mediaFileModel := &models.MediaFile{
		Name:        finalFilename,
		DownloadURL: downloadURL,
	}

	createdFile, err := s.repo.Create(mediaFileModel)
	if err != nil {
		// Attempt to remove the physically saved file if DB record creation fails
		if removeErr := os.Remove(fullDiskPath); removeErr != nil {
			log.Printf("CRITICAL: Failed to save file record to DB and also failed to remove orphaned file from disk '%s': %v", fullDiskPath, removeErr)
		}
		return nil, fmt.Errorf("failed to save media file record: %w", err)
	}

	return s.mapModelToResponse(createdFile), nil
}

func (s *mediaFileService) GetMediaFileByID(idStr string) (*dto.MediaFileResponse, error) {
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}
	file, err := s.repo.FindByID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("media file not found")
		}
		return nil, fmt.Errorf("failed to get media file by ID: %w", err)
	}
	return s.mapModelToResponse(file), nil
}

func (s *mediaFileService) ListMediaFiles(filter dto.MediaFileListFilter) (*dto.MediaFilesListResponse, error) {
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20 // Default page size
	}
	if filter.PageSize > 100 { // Max page size
		filter.PageSize = 100
	}

	files, total, err := s.repo.List(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list media files: %w", err)
	}

	responses := make([]dto.MediaFileListItemResponse, len(files))
	for i, file := range files {
		responses[i] = s.mapModelToListItemResponse(&file)
	}

	return &dto.MediaFilesListResponse{
		Data:     responses,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (s *mediaFileService) DeleteMediaFile(idStr string, userID uuid.UUID) error {
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("invalid ID format")
	}

	file, err := s.repo.FindByID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("media file not found")
		}
		return fmt.Errorf("failed to find media file for deletion: %w", err)
	}

	// Let's assume `file.DownloadURL` can be parsed to get the relative path
	parsedURL, err := url.Parse(file.DownloadURL)
	if err != nil {
		log.Printf("Warning: Could not parse DownloadURL to determine file path for deletion: %s. Skipping disk deletion.", file.DownloadURL)
	} else {
		// The path part of the URL, e.g. /files/images/my.jpg
		// We need to strip the static file prefix, e.g. /files
		relativePath := strings.TrimPrefix(parsedURL.Path, s.cfg.App.StaticFilePrefix)
		fullDiskPath := filepath.Join(s.cfg.App.UploadPath, filepath.FromSlash(relativePath))

		// Delete file from disk
		if err := os.Remove(fullDiskPath); err != nil {
			// Log error but don't necessarily fail the DB delete if file is already gone
			log.Printf("Warning: failed to delete file from disk '%s': %v", fullDiskPath, err)
			if !os.IsNotExist(err) {
				// If it's not a "file not exist" error, it might be more serious (e.g., permissions)
				// Depending on policy, you might want to return an error here.
			}
		}
	}

	// Delete record from DB
	if err := s.repo.Delete(uid); err != nil {
		return fmt.Errorf("failed to delete media file record: %w", err)
	}

	log.Printf("User %s deleted media file %s (ID: %s)", userID, file.Name, idStr)
	return nil
}
