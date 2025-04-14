// handlers/product_handlers.go
package handlers

import (
	"inventory_system/models"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

type CreateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"required"`
}

type UpdateProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gte=0"`
	Category    string  `json:"category"`
}

// GetProducts retrieves all products with optional filtering
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var products []models.Product
	db := h.DB

	// Apply filters if provided
	if category := c.Query("category"); category != "" {
		db = db.Where("category = ?", category)
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			db = db.Where("price >= ?", price)
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			db = db.Where("price <= ?", price)
		}
	}

	result := db.Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct retrieves a single product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	result := h.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct adds a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate category
	if !isValidCategory(input.Category) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
	}

	result := h.DB.Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	
	if result := h.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply updates only for fields that were provided
	updates := make(map[string]interface{})
	
	if input.Name != "" {
		updates["name"] = input.Name
	}
	
	if input.Description != "" {
		updates["description"] = input.Description
	}
	
	if input.Price != 0 {
		updates["price"] = input.Price
	}
	
	if input.Category != "" {
		if !isValidCategory(input.Category) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
			return
		}
		updates["category"] = input.Category
	}

	// Apply updates
	if len(updates) > 0 {
		result := h.DB.Model(&product).Updates(updates)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}
	}

	// Get updated product
	h.DB.First(&product, id)
	c.JSON(http.StatusOK, product)
}

// UploadProductImage handles file uploads for product images
func (h *ProductHandler) UploadProductImage(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	
	if result := h.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG files are allowed"})
		return
	}

	// Create unique filename
	filename := "product_" + id + ext
	filepath := filepath.Join("uploads", filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update product with image path
	h.DB.Model(&product).Update("image_path", filepath)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"filepath": filepath,
		"product_id": id,
	})
}

// GetProductImage serves the product image
func (h *ProductHandler) GetProductImage(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	
	if result := h.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.ImagePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No image available for this product"})
		return
	}

	c.File(product.ImagePath)
}

// Utility function to validate product category
func isValidCategory(category string) bool {
	validCategories := []string{"Electronics", "Apparel", "Footwear", "Furniture", "Appliances"}
	for _, c := range validCategories {
		if category == c {
			return true
		}
	}
	return false
}