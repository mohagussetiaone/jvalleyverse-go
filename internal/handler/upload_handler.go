package handler

import (
	"fmt"
	"path/filepath"
	"strings"

	"jvalleyverse/internal/minio"

	"github.com/gofiber/fiber/v2"
)

// UploadHandler handles file uploads to MinIO.
type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// allowedExtensions defines which file types can be uploaded.
var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".svg": true,
	".mp4":  true, ".pdf": true, ".zip": true,
}

// maxUploadSize is 10 MB.
const maxUploadSize = 10 << 20

// Upload receives a multipart file and stores it in MinIO.
//
//	POST /api/upload
//	Content-Type: multipart/form-data
//	Fields:
//	  - file: the file to upload (required)
//	  - folder: subfolder path, e.g. "courses", "avatars", "showcases", "lessons" (required)
//
//	Response 201:
//	  { "url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/uuid.jpg",
//	    "object_name": "courses/uuid.jpg",
//	    "size": 12345,
//	    "content_type": "image/jpeg" }
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
	if !minio.IsAvailable() {
		return c.Status(503).JSON(fiber.Map{
			"error": "File upload is not available (MinIO not configured)",
		})
	}

	// Parse folder from form field
	folder := strings.TrimSpace(c.FormValue("folder"))
	if folder == "" {
		return c.Status(400).JSON(fiber.Map{"error": "folder is required (e.g. courses, avatars, showcases, lessons)"})
	}

	// Validate folder name (no path traversal)
	if strings.Contains(folder, "..") || strings.Contains(folder, "/") {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid folder name"})
	}

	// Get the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file is required"})
	}

	// Validate file size
	if fileHeader.Size > maxUploadSize {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("file too large, maximum size is %d MB", maxUploadSize/(1<<20)),
		})
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("file type '%s' is not allowed (allowed: jpg, jpeg, png, gif, webp, svg, mp4, pdf, zip)", ext),
		})
	}

	// Upload to MinIO
	result, err := minio.DefaultClient.UploadFile(c.UserContext(), fileHeader, folder)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "file upload failed"})
	}

	return c.Status(201).JSON(result)
}
