// main.go
package main

import (
	"inventory_system/database"
	"inventory_system/routes"
	"log"
	"os"
)

func main() {
	// Create uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// Initialize database
	db := database.Initialize()

	// Setup router
	r := routes.SetupRouter(db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}