// handlers/image_handlers.go
package handlers

import (
	"net/http"
	"strconv"
	"log"
	"inventory_system/models"
	"inventory_system/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImageHandler handles product image operations
type ImageHandler struct {
	DB *gorm.DB
}

// NewImageHandler creates a new image handler
func NewImageHandler(db *gorm.DB) *ImageHandler {
	return &ImageHandler{DB: db}
}

// UploadProductImage handles uploading an image for a specific product
func (h *ImageHandler) UploadProductImage(c *gin.Context) {
	// Parse product ID
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        log.Printf("Error parsing product ID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }
    log.Printf("Product ID parsed: %d", productID)

	// Verify product exists
	var product models.Product
	if err := h.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product ID invalid"})
		return
	}

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// Limit request body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, utils.MaxFileSize)

	// Validate file
	if err := utils.ValidateImage(file); err != nil {
		if err.Error() == "file exceeds 5MB limit" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Save file
	filePath, err := utils.SaveProductImage(c, file, uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Update product database record
	product.ImagePath = filePath
	if err := h.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"image_path": filePath,
	})
}

// DownloadProductImage serves the image for a specific product
func (h *ImageHandler) DownloadProductImage(c *gin.Context) {
	// Parse product ID
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Verify product exists
	var product models.Product
	if err := h.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product ID invalid"})
		return
	}

	// Get image path
	imagePath, err := utils.GetProductImagePath(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No image found for product"})
		return
	}

	// Serve the file
	c.File(imagePath)
}