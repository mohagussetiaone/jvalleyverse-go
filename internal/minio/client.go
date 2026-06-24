package minio

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"jvalleyverse/pkg/config"

	"github.com/google/uuid"
	minioLib "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps the minio client with app-specific utilities.
type Client struct {
	mc     *minioLib.Client
	bucket string
	cdnURL string
}

// DefaultClient is the package-level MinIO client.
var DefaultClient *Client

// ConnectMinio initializes the MinIO client and ensures the bucket exists.
func ConnectMinio() {
	cfg := config.AppConfig

	if cfg.MinioAccessKey == "" || cfg.MinioSecretKey == "" {
		log.Println("⚠️  MINIO_ACCESS_KEY or MINIO_SECRET_KEY not set — file upload disabled")
		DefaultClient = nil
		return
	}

	mc, err := minioLib.New(cfg.MinioEndpoint, &minioLib.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		log.Printf("⚠️  MinIO connection failed (upload disabled): %v\n", err)
		DefaultClient = nil
		return
	}

	ctx := context.Background()
	exists, err := mc.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		log.Printf("⚠️  MinIO bucket check failed: %v\n", err)
		DefaultClient = nil
		return
	}
	if !exists {
		log.Printf("⚠️  MinIO bucket '%s' does not exist — creating it\n", cfg.MinioBucket)
		if err := mc.MakeBucket(ctx, cfg.MinioBucket, minioLib.MakeBucketOptions{}); err != nil {
			log.Printf("⚠️  Failed to create MinIO bucket: %v\n", err)
			DefaultClient = nil
			return
		}
	}

	DefaultClient = &Client{
		mc:     mc,
		bucket: cfg.MinioBucket,
		cdnURL: strings.TrimRight(cfg.MinioCDNURL, "/"),
	}

	log.Printf("✅ MinIO connected — bucket: %s, CDN: %s\n", cfg.MinioBucket, cfg.MinioCDNURL)
}

// IsAvailable returns true when the MinIO client is ready.
func IsAvailable() bool {
	return DefaultClient != nil && DefaultClient.mc != nil
}

// UploadStatus describes the result of a file upload.
type UploadResult struct {
	ObjectName string `json:"object_name"`
	URL        string `json:"url"`
	Size       int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// UploadFile uploads a file from a multipart file header to MinIO.
// Returns the object name and public CDN URL.
func (c *Client) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*UploadResult, error) {
	// Generate unique file name
	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Open the file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = detectContentType(ext)
	}

	// Upload to MinIO
	_, err = c.mc.PutObject(ctx, c.bucket, objectName, src, fileHeader.Size,
		minioLib.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	// Build public CDN URL
	url := fmt.Sprintf("%s/%s", c.cdnURL, objectName)

	return &UploadResult{
		ObjectName:  objectName,
		URL:         url,
		Size:        fileHeader.Size,
		ContentType: contentType,
	}, nil
}

// DeleteFile removes an object from MinIO by its object name.
func (c *Client) DeleteFile(ctx context.Context, objectName string) error {
	return c.mc.RemoveObject(ctx, c.bucket, objectName, minioLib.RemoveObjectOptions{})
}

// GeneratePresignedUploadURL creates a presigned URL for direct browser upload.
// Expiration defaults to 15 minutes.
func (c *Client) GeneratePresignedUploadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	if expiry == 0 {
		expiry = 15 * time.Minute
	}
	u, err := c.mc.PresignedPutObject(ctx, c.bucket, objectName, expiry)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// detectContentType returns a MIME type based on file extension.
func detectContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".mp4":
		return "video/mp4"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}
