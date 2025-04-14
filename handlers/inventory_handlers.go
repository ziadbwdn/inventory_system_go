// handlers/inventory_handlers.go
package handlers

import (
	"inventory_system/database"
	"inventory_system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryHandler struct {
	DB *gorm.DB
}

type StockAdjustment struct {
	Action string `json:"action" binding:"required,oneof=add remove"`
	Value  int    `json:"value" binding:"required,gt=0"`
}

// GetInventory retrieves inventory information with optional product filtering
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	var inventories []models.Inventory
	db := h.DB

	// Apply product filter if provided
	if productID := c.Query("product_id"); productID != "" {
		db = db.Where("product_id = ?", productID)
	}

	// Join with products to get more information
	result := db.Preload("Product").Find(&inventories)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventory"})
		return
	}

	c.JSON(http.StatusOK, inventories)
}

// AdjustStock updates the inventory quantity for a specific product
func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	productID := c.Param("product_id")
	location := c.DefaultQuery("location", "Warehouse A") // Default to Warehouse A if not specified
	
	var adjustment StockAdjustment
	if err := c.ShouldBindJSON(&adjustment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	var product models.Product
	if result := h.DB.First(&product, productID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Find or create inventory entry
	var inventory models.Inventory
	result := h.DB.Where("product_id = ? AND location = ?", productID, location).First(&inventory)
	
	if result.Error != nil {
		// Create new inventory entry if it doesn't exist
		inventory = models.Inventory{
			ProductID: product.ID,
			Location:  location,
			Quantity:  0,
		}
	}

	// Adjust stock based on action
	if adjustment.Action == "add" {
		inventory.Quantity += adjustment.Value
	} else { // "remove"
		if inventory.Quantity < adjustment.Value {
			c.JSON(http.StatusConflict, gin.H{"error": "Insufficient stock"})
			return
		}
		inventory.Quantity -= adjustment.Value
	}

	// Save changes
	if result.Error != nil {
		// Create new record
		if createResult := h.DB.Create(&inventory); createResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return
		}
	} else {
		// Update existing record
		if updateResult := h.DB.Save(&inventory); updateResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return
		}
	}

	c.JSON(http.StatusOK, inventory)
}

// GetInventoryByLocation groups inventory by warehouse/location
func (h *InventoryHandler) GetInventoryByLocation(c *gin.Context) {
	results, err := database.GetStockDistributionByLocation(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventory by location"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetLowStockProducts returns products with quantity below threshold
func (h *InventoryHandler) GetLowStockProducts(c *gin.Context) {
	threshold, err := strconv.Atoi(c.DefaultQuery("threshold", "20"))
	if err != nil {
		threshold = 20
	}

	results, err := database.GetLowStockProducts(h.DB, threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve low stock products"})
		return
	}

	c.JSON(http.StatusOK, results)
}