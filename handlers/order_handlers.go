// handlers/order_handlers.go
package handlers

import (
	"inventory_system/database"
	"inventory_system/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

type CreateOrderInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// GetOrders retrieves all orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	var orders []models.Order
	
	result := h.DB.Preload("Product").Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder retrieves a single order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	result := h.DB.Preload("Product").First(&order, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateOrder creates a new order and updates inventory
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	var product models.Product
	if result := h.DB.First(&product, input.ProductID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if there's enough inventory (using first warehouse as default)
	var inventory models.Inventory
	result := h.DB.Where("product_id = ? AND location = ?", input.ProductID, "Warehouse A").First(&inventory)
	if result.Error != nil || inventory.Quantity < input.Quantity {
		c.JSON(http.StatusConflict, gin.H{"error": "Insufficient stock"})
		return
	}

	// Begin transaction
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create order
	order := models.Order{
		ProductID:  input.ProductID,
		Quantity:   input.Quantity,
		OrderDate:  time.Now(),
		TotalPrice: product.Price * float64(input.Quantity),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Update inventory
	inventory.Quantity -= input.Quantity
	if err := tx.Save(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Fetch the complete order with product details
	h.DB.Preload("Product").First(&order, order.OrderID)

	c.JSON(http.StatusCreated, order)
}

// GetRevenueByCategory gets revenue statistics grouped by product category
func (h *OrderHandler) GetRevenueByCategory(c *gin.Context) {
	results, err := database.GetRevenueByCategory(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve revenue data"})
		return
	}

	c.JSON(http.StatusOK, results)
}