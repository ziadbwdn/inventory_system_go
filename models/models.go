// models/models.go
package models

import (
	"time"
)

// Product represents an item that can be sold
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey;type:int unsigned"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null;check:price >= 0"`
	Category    string    `json:"category" gorm:"size:50;not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ImagePath   string    `json:"image_path,omitempty" gorm:"size:255"`
}

// Inventory represents the stock of a product at a specific location
type Inventory struct {
	ProductID uint    `json:"product_id" gorm:"primaryKey;type:int unsigned"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null;default:0"`
	Location  string  `json:"location" gorm:"size:100;not null;primaryKey"`
}

// Order represents a customer order for a specific product
type Order struct {
	OrderID    uint      `json:"order_id" gorm:"primaryKey;type:int unsigned"`
	ProductID  uint      `json:"product_id" gorm:"type:int unsigned;not null"`
	Product    Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity   int       `json:"quantity" gorm:"not null;check:quantity > 0"`
	OrderDate  time.Time `json:"order_date" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"type:decimal(10,2);not null"`
}