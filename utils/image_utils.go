// image utils

package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// Maximum file size (5MB)
	MaxFileSize = 5 << 20
	// Base upload directory
	UploadDir = "./uploads/products"
)

// ValidateImage checks if the file is a valid image and within size limits
func ValidateImage(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > MaxFileSize {
		return fmt.Errorf("file exceeds 5MB limit")
	}

	// Open the file to check its content type
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read first 512 bytes to determine content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return err
	}

	// Reset the file reader
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// Check content type
	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/jpg" && contentType != "image/png" {
		return fmt.Errorf("only PNG/JPG/JPEG allowed")
	}

	return nil
}

// saveProductImage saves the uploaded file to the appropriate directory
func SaveProductImage(c *gin.Context, file *multipart.FileHeader, productID uint) (string, error) {
	// Generate unique filename with UUID
	fileExt := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + "-ProductImage" + fileExt

	// Ensure the directory exists
	productDir := filepath.Join(UploadDir, fmt.Sprintf("%d", productID))
	if err := os.MkdirAll(productDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Full path for the file
	dst := filepath.Join(productDir, newFilename)

	// Save the file using Gin's utility function
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return dst, nil
}

// GetProductImagePath retrieves the latest image for a product
func GetProductImagePath(productID uint) (string, error) {
	productDir := filepath.Join(UploadDir, fmt.Sprintf("%d", productID))
	
	// Check if directory exists
	if _, err := os.Stat(productDir); os.IsNotExist(err) {
		return "", fmt.Errorf("no image found for product")
	}

	// Read directory contents
	files, err := os.ReadDir(productDir)
	if err != nil {
		return "", err
	}

	// No files found
	if len(files) == 0 {
		return "", fmt.Errorf("no image found for product")
	}

	// Find the most recent file (based on name for simplicity)
	// In a real app, you might want to check file creation dates
	var latestFile string
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".jpg") || 
		   strings.HasSuffix(file.Name(), ".jpeg") || 
		   strings.HasSuffix(file.Name(), ".png")) {
			if latestFile == "" || file.Name() > latestFile {
				latestFile = file.Name()
			}
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("no image found for product")
	}

	return filepath.Join(productDir, latestFile), nil
}

// DeleteProductImages removes all images associated with a product
func DeleteProductImages(productID uint) error {
	productDir := filepath.Join(UploadDir, fmt.Sprintf("%d", productID))
	
	// Check if directory exists
	if _, err := os.Stat(productDir); os.IsNotExist(err) {
		// Directory doesn't exist, nothing to delete
		return nil
	}

	// Remove the directory and all its contents
	return os.RemoveAll(productDir)
}

// PathTraversalMiddleware prevents path traversal attacks
func PathTraversalMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        paramValue := c.Param("id")
        log.Printf("PathTraversalMiddleware checking param: %s", paramValue)
        
        if strings.Contains(paramValue, "..") {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
            c.Abort()
            return
        }
        log.Printf("PathTraversalMiddleware passed for param: %s", paramValue)
        c.Next()
    }
}