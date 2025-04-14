// database/queries.go
package database

import (

	"gorm.io/gorm"
	"log"
)

// GetLowStockProducts returns products with quantity less than threshold
func GetLowStockProducts(db *gorm.DB, threshold int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// Using efficient JOIN with INDEX hints for MySQL
	err := db.Table("inventories").
		Select("products.id, products.name, products.category, inventories.location, inventories.quantity").
		Joins("JOIN products USE INDEX (PRIMARY) ON inventories.product_id = products.id").
		Where("inventories.quantity < ?", threshold).
		Find(&results).Error
	
	return results, err
}

// GetRevenueByCategory calculates total revenue per product category
func GetRevenueByCategory(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// Using JOIN with INDEX hints and efficient aggregation
	err := db.Table("orders").
        Select("products.category, SUM(orders.total_price) as total_revenue, COUNT(orders.order_id) as order_count").
        // Joins("JOIN products USE INDEX (PRIMARY) ON orders.product_id = products.id"). // Potential issue
        Joins("JOIN products ON orders.product_id = products.id"). // Use standard join first
        Group("products.category").
        Find(&results).Error

    // --- TEMPORARY LOGGING ---
    if err != nil {
        log.Printf("DEBUG: GetRevenueByCategory DB Error: %v\n", err)
    } else {
        log.Printf("DEBUG: GetRevenueByCategory Success. Rows found: %d\n", len(results))
    }
    // --- END TEMPORARY LOGGING ---

    return results, err
}

// GetStockDistributionByLocation shows total stock across all locations
func GetStockDistributionByLocation(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// Using index on location for faster grouping
	err := db.Table("inventories").
		Select("location, SUM(quantity) as total_stock, COUNT(DISTINCT product_id) as product_count").
		Group("location").
		Find(&results).Error
	
	return results, err
}

// GetTopSellingProducts returns the top selling products
func GetTopSellingProducts(db *gorm.DB, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := db.Table("orders").
		Select("products.id, products.name, SUM(orders.quantity) as total_sold, SUM(orders.total_price) as total_revenue").
		Joins("JOIN products ON orders.product_id = products.id").
		Group("products.id").
		Order("total_sold DESC").
		Limit(limit).
		Find(&results).Error
	
	return results, err
}

// GetInventoryValueByCategory calculates the total inventory value per category
func GetInventoryValueByCategory(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := db.Table("inventories").
		Select("products.category, SUM(inventories.quantity * products.price) as total_value").
		Joins("JOIN products ON inventories.product_id = products.id").
		Group("products.category").
		Find(&results).Error
	
	return results, err
}