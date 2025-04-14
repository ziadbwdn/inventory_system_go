// database/db.go
package database

import (
	// "errors"
	"fmt"
	"inventory_system/models"
	"log"
	"os"
	"time"
	"math/rand"
	"strings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize sets up the database connection and migrations
func Initialize() *gorm.DB {
	// Get database credentials from environment variables or use defaults
	dbUser := getEnv("DB_USER", "root")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "inventory_system")

	var dbPassword string
	fmt.Print("Enter database password: ") // Prompt the user
	_, errScan := fmt.Scan(&dbPassword)     // Read input into dbPassword
	if errScan != nil {
		log.Fatalf("Failed to read password: %v", errScan)
	}
	// --- End Modification ---

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open connection to database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool parameters
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)
	
	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)
	
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Migrate the schema
	err = db.AutoMigrate(&models.Product{}, &models.Inventory{}, &models.Order{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Check if we need to seed data
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count == 0 {
		seedData(db)
	}

	log.Println("Database connection established successfully")
	return db
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// seedData inserts sample records into the database
func seedData(db *gorm.DB) {
	// Sample products
	products := []models.Product{
		{Name: "Laptop", Description: "High-performance laptop with 20GB RAM", Price: 999.99, Category: "Electronics"},
		{Name: "Smartphone", Description: "Latest model with 128GB storage", Price: 699.99, Category: "Electronics"},
		{Name: "Headphones", Description: "Wireless noise-cancelling headphones", Price: 199.99, Category: "Electronics"},
		{Name: "T-shirt", Description: "Cotton t-shirt, size M", Price: 19.99, Category: "Apparel"},
		{Name: "Jeans", Description: "Blue denim jeans, slim fit", Price: 49.99, Category: "Apparel"},
		{Name: "Sneakers", Description: "Running shoes, size 11", Price: 89.99, Category: "Footwear"},
		{Name: "Coffee Table", Description: "Wooden coffee table", Price: 149.99, Category: "Furniture"},
		{Name: "Desk Chair", Description: "Ergonomic office chair", Price: 199.99, Category: "Furniture"},
		{Name: "Blender", Description: "High-speed blender for smoothies", Price: 79.99, Category: "Appliances"},
		{Name: "Toaster", Description: "6-slice toaster", Price: 39.99, Category: "Appliances"},
	}

	for _, product := range products {
		db.Create(&product)
	}

	// Sample inventory
	locations := []string{"Warehouse A", "Warehouse B", "Store 1", "Store 2"}
	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			log.Fatalf("Failed to create product: %v", err)
		}
	}
	
	// Update inventory creation to use the correct product IDs
	for _, product := range products {
		for _, location := range locations {
			inventory := models.Inventory{
				ProductID: product.ID, // Now correctly set
				Quantity:  generateRandomQuantity(product.Category, location),
				Location:  location,
			}
			if err := db.Create(&inventory).Error; err != nil {
				log.Printf("Failed to create inventory: %v", err)
			}
		}
	}

	// Sample orders
	for i := 0; i < 15; i++ {
		productID := uint(generateRandomInt(1, len(products)))
		quantity := generateRandomInt(1, 5)
		
		var product models.Product
		if err := db.First(&product, productID).Error; err != nil {
    		log.Printf("Failed to find product: %v", err)
    		continue // Skip this iteration
}
		
		order := models.Order{
			ProductID:  productID,
			Quantity:   quantity,
			OrderDate:  time.Now().AddDate(0, 0, -generateRandomInt(0, 30)), // Random date within last 30 days
			TotalPrice: product.Price * float64(quantity),
		}
		if err := db.Create(&order).Error; err != nil {
			log.Printf("Failed to create order: %v", err)
		}
		db.Create(&order)

		// Update inventory for the first location
		var inventory models.Inventory
		db.Where("product_id = ? AND location = ?", productID, "Warehouse A").First(&inventory)
		inventory.Quantity -= quantity
		if inventory.Quantity < 0 {
			inventory.Quantity = 0
		}
		db.Save(&inventory)
	}

	log.Println("Database seeded successfully")
}

// Helper function to generate random integer between min and max (inclusive)
// Uses math/rand - Ensure rand is seeded once at application startup!
func generateRandomInt(min, max int) int {
    if min > max {
        min, max = max, min // Swap if min > max
    }
	// rand.Intn(n) returns a value in [0, n). We want [min, max].
	// Range size is max - min + 1. Result is rand.Intn(rangeSize) + min.
	return rand.Intn(max-min+1) + min
}

// Helper function to generate random quantity based on product category and location
func generateRandomQuantity(category, location string) int {
	base := 30

	// Electronics have lower stock
	if category == "Electronics" {
		base = 15
	}

	// Warehouses have more stock than stores (Safer check)
	if strings.HasPrefix(location, "Warehouse") {
		base *= 2
	}

    // Add some randomness, ensuring quantity is not negative
    offset := generateRandomInt(-10, 20)
    quantity := base + offset
    if quantity < 0 {
        quantity = 0
    }
	return quantity
}