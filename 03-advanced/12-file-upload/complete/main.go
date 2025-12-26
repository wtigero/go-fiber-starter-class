package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

// ============ Constants ============

const (
	UploadDir   = "./uploads"
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"application/pdf": true,
}

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// ============ Types ============

type FileInfo struct {
	OriginalName string    `json:"original_name"`
	SavedName    string    `json:"saved_name"`
	Size         int64     `json:"size"`
	Type         string    `json:"type"`
	URL          string    `json:"url"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

func main() {
	// สร้าง upload directory
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    MaxFileSize,
		ErrorHandler: errorHandler,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Serve uploaded files
	app.Static("/files", UploadDir)

	// Serve static UI
	app.Static("/", "./static")

	// API routes
	app.Get("/api/info", infoHandler)
	app.Post("/upload", uploadSingleHandler)
	app.Post("/upload/multiple", uploadMultipleHandler)
	app.Post("/upload/image", uploadImageHandler)
	app.Get("/api/files", listFilesHandler)
	app.Delete("/api/files/:filename", deleteFileHandler)

	log.Println("🚀 File Upload server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

// ============ Handlers ============

func infoHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message":       "File Upload API",
		"max_file_size": fmt.Sprintf("%d MB", MaxFileSize/1024/1024),
		"allowed_types": getAllowedTypes(),
		"endpoints": fiber.Map{
			"upload_single":   "POST /upload",
			"upload_multiple": "POST /upload/multiple",
			"upload_image":    "POST /upload/image",
			"list_files":      "GET /api/files",
			"delete_file":     "DELETE /api/files/:filename",
		},
	})
}

func uploadSingleHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบไฟล์ กรุณาส่งไฟล์ในฟิลด์ 'file'",
		})
	}

	// Validate file
	if err := validateFile(file); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Save file
	fileInfo, err := saveFile(file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถบันทึกไฟล์ได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "อัพโหลดไฟล์สำเร็จ",
		"file":    fileInfo,
	})
}

func uploadMultipleHandler(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถอ่านฟอร์มได้",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบไฟล์ กรุณาส่งไฟล์ในฟิลด์ 'files'",
		})
	}

	var uploadedFiles []FileInfo
	var errors []string

	for _, file := range files {
		// Validate
		if err := validateFile(file); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %s", file.Filename, err.Error()))
			continue
		}

		// Save
		fileInfo, err := saveFile(file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: บันทึกไม่สำเร็จ", file.Filename))
			continue
		}

		uploadedFiles = append(uploadedFiles, *fileInfo)
	}

	return c.JSON(fiber.Map{
		"success":        len(errors) == 0,
		"message":        fmt.Sprintf("อัพโหลดสำเร็จ %d/%d ไฟล์", len(uploadedFiles), len(files)),
		"uploaded_files": uploadedFiles,
		"errors":         errors,
	})
}

func uploadImageHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบไฟล์ กรุณาส่งไฟล์ในฟิลด์ 'image'",
		})
	}

	// Validate image type
	if !allowedImageTypes[file.Header.Get("Content-Type")] {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "ไฟล์ต้องเป็นรูปภาพ (jpeg, png, gif, webp)",
		})
	}

	// Get resize dimensions
	width, _ := strconv.Atoi(c.FormValue("width", "0"))
	height, _ := strconv.Atoi(c.FormValue("height", "0"))
	quality, _ := strconv.Atoi(c.FormValue("quality", "85"))

	if quality < 1 || quality > 100 {
		quality = 85
	}

	// Save and process
	fileInfo, err := saveAndResizeImage(file, width, height, quality)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("ไม่สามารถประมวลผลรูปภาพได้: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "อัพโหลดและประมวลผลรูปภาพสำเร็จ",
		"file":    fileInfo,
		"options": fiber.Map{
			"width":   width,
			"height":  height,
			"quality": quality,
		},
	})
}

func listFilesHandler(c *fiber.Ctx) error {
	files, err := os.ReadDir(UploadDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถอ่านไดเรกทอรีได้",
		})
	}

	var fileList []fiber.Map
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, _ := file.Info()
		fileList = append(fileList, fiber.Map{
			"name":         file.Name(),
			"size":         info.Size(),
			"url":          "/files/" + file.Name(),
			"modified_at":  info.ModTime(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"files":   fileList,
		"count":   len(fileList),
	})
}

func deleteFileHandler(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Prevent path traversal
	filename = filepath.Base(filename)
	filePath := filepath.Join(UploadDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบไฟล์",
		})
	}

	if err := os.Remove(filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่สามารถลบไฟล์ได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบไฟล์สำเร็จ",
	})
}

// ============ Helpers ============

func validateFile(file *multipart.FileHeader) error {
	// Check size
	if file.Size > MaxFileSize {
		return fmt.Errorf("ไฟล์ใหญ่เกิน %d MB", MaxFileSize/1024/1024)
	}

	// Check MIME type
	contentType := file.Header.Get("Content-Type")
	if !allowedMimeTypes[contentType] {
		return fmt.Errorf("ประเภทไฟล์ไม่ได้รับอนุญาต: %s", contentType)
	}

	return nil
}

func saveFile(file *multipart.FileHeader) (*FileInfo, error) {
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	newName := uuid.New().String() + ext
	savePath := filepath.Join(UploadDir, newName)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	return &FileInfo{
		OriginalName: file.Filename,
		SavedName:    newName,
		Size:         file.Size,
		Type:         file.Header.Get("Content-Type"),
		URL:          "/files/" + newName,
		UploadedAt:   time.Now(),
	}, nil
}

func saveAndResizeImage(file *multipart.FileHeader, width, height, quality int) (*FileInfo, error) {
	// Open source file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Decode image
	img, err := imaging.Decode(src)
	if err != nil {
		return nil, err
	}

	// Resize if dimensions specified
	if width > 0 || height > 0 {
		if width > 0 && height > 0 {
			img = imaging.Fit(img, width, height, imaging.Lanczos)
		} else if width > 0 {
			img = imaging.Resize(img, width, 0, imaging.Lanczos)
		} else {
			img = imaging.Resize(img, 0, height, imaging.Lanczos)
		}
	}

	// Generate filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	newName := uuid.New().String() + ext
	savePath := filepath.Join(UploadDir, newName)

	// Save based on format
	var saveErr error
	switch ext {
	case ".png":
		saveErr = imaging.Save(img, savePath)
	case ".gif":
		saveErr = imaging.Save(img, savePath)
	default:
		saveErr = imaging.Save(img, savePath, imaging.JPEGQuality(quality))
	}

	if saveErr != nil {
		return nil, saveErr
	}

	// Get file info
	info, _ := os.Stat(savePath)

	return &FileInfo{
		OriginalName: file.Filename,
		SavedName:    newName,
		Size:         info.Size(),
		Type:         file.Header.Get("Content-Type"),
		URL:          "/files/" + newName,
		UploadedAt:   time.Now(),
	}, nil
}

func getAllowedTypes() []string {
	types := make([]string, 0, len(allowedMimeTypes))
	for t := range allowedMimeTypes {
		types = append(types, t)
	}
	return types
}

func errorHandler(c *fiber.Ctx, err error) error {
	return c.Status(500).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
