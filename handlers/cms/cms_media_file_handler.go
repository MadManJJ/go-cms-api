package cms

import (
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/MadManJJ/cms-api/dto"
	"github.com/MadManJJ/cms-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaFileHandler struct {
	Service  services.MediaFileServiceInterface
	validate *validator.Validate
}

func NewMediaFileHandler(service services.MediaFileServiceInterface) *MediaFileHandler {
	return &MediaFileHandler{
		Service:  service,
		validate: validator.New(),
	}
}

// HandleUploadMediaFile handles file uploads.
// @Summary      Upload Media File
// @Description  Uploads a new media file (image, pdf, svg).
// @Tags         CMS - Media Files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "The file to upload (jpg, png, webp, svg, pdf)"
// @Param        path formData string false "Optional sub-directory path for the file"
// @Param        replace formData boolean false "Optional. If true, replaces existing file with the same name and path. Default: false."
// @Success      201  {object} dto.MediaFileResponse "File uploaded successfully"
// @Failure      400  {object} dto.ErrorResponse "Bad Request (e.g., no file, invalid file type, invalid path)"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/media-files [post]
func (h *MediaFileHandler) HandleUploadMediaFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "File is required", Message: "No file uploaded in 'file' field."})
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "File upload error", Message: err.Error()})
	}

	var req dto.UploadMediaFileRequest
	if err := c.BodyParser(&req); err != nil {
		// This might not parse form fields correctly if mixed with file upload.
		// It's better to get form values directly.
	}

	customPathStr := c.FormValue("path")
	var customPath *string
	if customPathStr != "" {
		customPath = &customPathStr
	}

	replaceStr := c.FormValue("replace", "false") // default to "false" if not provided
	var replace bool
	if strings.ToLower(replaceStr) == "true" {
		replace = true
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to open uploaded file", Message: err.Error()})
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to read uploaded file", Message: err.Error()})
	}

	// Detect MIME type from content if header is not reliable enough or for extra validation
	// For robust SVG sanitization, you'd pass fileData (as string) to a sanitizer.
	// Here we'll use the header's MIME type for simplicity, but handler could validate it.
	// In a real app, you might want to use a library like `github.com/gabriel-vasile/mimetype`
	detectedMimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	if detectedMimeType == "" { // Fallback if extension is missing or unknown
		detectedMimeType = http.DetectContentType(fileData)
	}
	// If you have specific requirements for fileHeader.Header.Get("Content-Type"), use that or combine.
	// var userID uuid.UUID
	// userID := c.Locals("userID").(uuid.UUID) // Assuming userID is set by auth middleware and is uuid.UUID
	var userID uuid.UUID
	userIDLocal := c.Locals("userID")
	if userIDLocal != nil {
		var ok bool
		userID, ok = userIDLocal.(uuid.UUID)
		if !ok {
			// Log error หรือ return bad request ถ้า type ไม่ถูกต้อง แต่ในกรณีนี้ถ้ามีควรจะเป็น UUID
			log.Println("Error: userID in context is not of type uuid.UUID")
			// For now, let's assign a default if type assertion fails but value is not nil
			// This case should ideally not happen if auth middleware sets it correctly.
			userID = uuid.Nil // หรือ uuid.New() ถ้าต้องการค่าที่ไม่ nil สำหรับ test
		}
	} else {
		// *** สำหรับการทดสอบเท่านั้น: ถ้ายังไม่มี Auth Middleware ***
		log.Println("Warning: No userID found in context, using a default/mock userID for testing.")
		userID = uuid.New() // สร้าง UUID ใหม่สำหรับทดสอบ หรือใช้ค่าคงที่
		// หรือถ้า service/repo อนุญาต Nil UUID สำหรับบางกรณี (ซึ่งไม่ควรสำหรับ created_by)
		// userID = uuid.Nil
	}
	response, err := h.Service.UploadMediaFile(fileHeader.Filename, detectedMimeType, fileData, customPath, &replace, userID)
	if err != nil {
		if strings.Contains(err.Error(), "invalid file type") || strings.Contains(err.Error(), "invalid custom path") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Bad Request", Message: err.Error()})
		}
		// Log the full error for internal tracking
		log.Printf("Error uploading file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to upload media file", Message: "An internal error occurred."})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// HandleListMediaFiles lists all media files with pagination, search, and sort.
// @Summary      List Media Files
// @Description  Retrieves a paginated list of media files, with optional search by name and sorting.
// @Tags         CMS - Media Files
// @Produce      json
// @Param        search query string false "Search term for file name"
// @Param        page query int false "Page number (default: 1)"
// @Param        pageSize query int false "Items per page (default: 20, max: 100)"
// @Param        sortBy query string false "Sort by 'name' or 'created_at' (default: 'created_at')" Enums(name, created_at)
// @Param        order query string false "Sort order 'asc' or 'desc' (default: 'desc')" Enums(asc, desc)
// @Success      200  {object} dto.MediaFilesListResponse
// @Failure      400  {object} dto.ErrorResponse "Invalid query parameters"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/media-files [get]
func (h *MediaFileHandler) HandleListMediaFiles(c *fiber.Ctx) error {
	var filter dto.MediaFileListFilter
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid query parameters", Message: err.Error()})
	}

	// Validate filter (if any validation tags are added to MediaFileListFilter)
	if err := h.validate.Struct(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Validation failed for query parameters", Message: err.Error()})
	}

	response, err := h.Service.ListMediaFiles(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to list media files", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// HandleGetMediaFileByID retrieves a specific media file by its ID.
// @Summary      Get Media File by ID
// @Description  Retrieves details for a specific media file by its UUID.
// @Tags         CMS - Media Files
// @Produce      json
// @Param        id path string true "Media File ID (UUID)"
// @Success      200  {object} dto.MediaFileResponse
// @Failure      400  {object} dto.ErrorResponse "Invalid ID format"
// @Failure      404  {object} dto.ErrorResponse "Media file not found"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error"
// @Router       /cms/media-files/{id} [get]
func (h *MediaFileHandler) HandleGetMediaFileByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	response, err := h.Service.GetMediaFileByID(idStr)
	if err != nil {
		if errors.Is(err, errors.New("invalid UUID format")) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid ID format", Message: err.Error()})
		}
		if errors.Is(err, errors.New("media file not found")) || errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Media file not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to get media file", Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// HandleDeleteMediaFile deletes a media file.
// @Summary      Delete Media File
// @Description  Deletes a media file by its UUID from both the database and the disk.
// @Tags         CMS - Media Files
// @Param        id path string true "Media File ID (UUID)"
// @Success      204  "No Content - File deleted successfully"
// @Failure      400  {object} dto.ErrorResponse "Invalid ID format"
// @Failure      404  {object} dto.ErrorResponse "Media file not found"
// @Failure      500  {object} dto.ErrorResponse "Internal Server Error (e.g., failed to delete file from disk or DB)"
// @Router       /cms/media-files/{id} [delete]
func (h *MediaFileHandler) HandleDeleteMediaFile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var userID uuid.UUID
	userIDLocal := c.Locals("userID")
	if userIDLocal != nil {
		var ok bool
		userID, ok = userIDLocal.(uuid.UUID)
		if !ok {
			log.Println("Error: userID in context is not of type uuid.UUID")
			userID = uuid.Nil
		}
	} else {
		// *** สำหรับการทดสอบเท่านั้น ***
		log.Println("Warning: No userID found in context for delete, using a default/mock userID.")
		userID = uuid.New()
	}

	err := h.Service.DeleteMediaFile(idStr, userID)
	if err != nil {
		if errors.Is(err, errors.New("invalid ID format")) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid ID format", Message: err.Error()})
		}
		if errors.Is(err, errors.New("media file not found")) || errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "Not Found", Message: "Media file not found for deletion"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to delete media file", Message: err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
