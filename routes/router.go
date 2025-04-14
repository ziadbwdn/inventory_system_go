// routes/router.go
package routes

import (
	"inventory_system/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configures all the routes for our application
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Middleware for CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Initialize handlers
	productHandler := &handlers.ProductHandler{DB: db}
	inventoryHandler := &handlers.InventoryHandler{DB: db}
	orderHandler := &handlers.OrderHandler{DB: db}

	// Static file serving
	r.Static("/uploads", "./uploads")

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	// Product routes
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("", productHandler.GetProducts)
		productRoutes.GET("/:id", productHandler.GetProduct)
		productRoutes.POST("", productHandler.CreateProduct)
		productRoutes.PUT("/:id", productHandler.UpdateProduct)
		productRoutes.POST("/:id/upload", productHandler.UploadProductImage)
		productRoutes.GET("/:id/image", productHandler.GetProductImage)
	}

	// Inventory routes
	inventoryRoutes := r.Group("/inventory")
	{
		inventoryRoutes.GET("", inventoryHandler.GetInventory)
		inventoryRoutes.PATCH("/:product_id", inventoryHandler.AdjustStock)
		inventoryRoutes.GET("/locations", inventoryHandler.GetInventoryByLocation)
		inventoryRoutes.GET("/low-stock", inventoryHandler.GetLowStockProducts)
	}

	// Order routes
	orderRoutes := r.Group("/orders")
	{
		orderRoutes.GET("", orderHandler.GetOrders)
		orderRoutes.GET("/:id", orderHandler.GetOrder)
		orderRoutes.POST("", orderHandler.CreateOrder)
		orderRoutes.GET("/revenue", orderHandler.GetRevenueByCategory)
	}

	return r
}