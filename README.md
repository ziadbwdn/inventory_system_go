# Inventory Management System - New Branch (for Assignment Update Purposes)

A RESTful API for inventory management built with Go, Gin framework, and MySQL.

## Features

- Product management (CRUD operations)
- Inventory tracking across multiple locations
- Order processing with automatic inventory updates
- File upload and serving for product images
- Comprehensive data reporting

## Prerequisites

- Go 1.19 or higher
- MySQL 8.0 or higher

## Project Structure

```
/inventory_system_go
├── database/            # Database configuration and queries
│   ├── db.go
│   ├── queries.go
│   └── scripts/
│       └── schema.sql
├── docs/
│   └── documentation.pdf
├── handlers/            # Gin route logic
│   ├── inventory_handlers.go
│   ├── order_handlers.go
│   ├── product_handlers.go
│   └── image_handlers.go
├── main.go
├── models/              # Structs (Product, Inventory, Order)
│   └── models.go
├── routes/              # Gin router groups
│   └── router.go
├── uploads/             # Product images
│   └── products/
├── utils/               # Utility functions
│   └── file_utils.go
├── go.mod               # Dependencies (Gin, GORM, MySQL driver)
├── go.sum
└── README.md            # This file
```

## Database Setup

### Method 1: Run the Schema SQL Script

1. Log in to MySQL:
   ```bash
   mysql -u root -p
   ```

2. Run the schema script:
   ```bash
   mysql -u root -p < database/scripts/schema.sql
   ```

### Method 2: Let the Application Handle It

1. Set environment variables for database connection:
   ```bash
   export DB_USER=root
   export DB_PASSWORD=your_password
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_NAME=inventory
   ```

2. Run the application, which will create the database schema automatically using GORM AutoMigrate.

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/inventory-system.git
   cd inventory-system
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   export DB_USER=root
   export DB_PASSWORD=your_password
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_NAME=inventory
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. The server will start at http://localhost:8080

## API Documentation

### Products

#### Get all products
```bash
curl -X GET http://localhost:8080/products
```

#### Get products by category
```bash
curl -X GET "http://localhost:8080/products?category=Electronics"
```

#### Get products by price range
```bash
curl -X GET "http://localhost:8080/products?min_price=50&max_price=200"
```

#### Get a specific product
```bash
curl -X GET http://localhost:8080/products/1
```

#### Create a new product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Wireless Mouse","description":"Ergonomic wireless mouse","price":29.99,"category":"Electronics"}'
```

#### Update a product
```bash
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Product Name","price":39.99}'
```

#### Upload a product image
```bash
curl -X POST http://localhost:8080/products/1/upload \
  -F "image=@/path/to/image.jpg"
```

#### Get a product image
```bash
curl -X GET http://localhost:8080/products/1/image
```

### Inventory

#### Get all inventory
```bash
curl -X GET http://localhost:8080/inventory
```

#### Get inventory by product
```bash
curl -X GET "http://localhost:8080/inventory?product_id=1"
```

#### Adjust stock
```bash
curl -X PATCH "http://localhost:8080/inventory/1?location=Warehouse%20A" \
  -H "Content-Type: application/json" \
  -d '{"action":"add","value":10}'
```

#### Get inventory by location
```bash
curl -X GET http://localhost:8080/inventory/locations
```

#### Get low stock products
```bash
curl -X GET "http://localhost:8080/inventory/low-stock?threshold=15"
```

### Orders

#### Get all orders
```bash
curl -X GET http://localhost:8080/orders
```

#### Get a specific order
```bash
curl -X GET http://localhost:8080/orders/1
```

#### Create a new order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'
```

#### Get revenue by category
```bash
curl -X GET http://localhost:8080/orders/revenue
```

## File Upload/Download Workflow

```
┌─────────────┐         ┌──────────────┐         ┌────────────────┐
│ Client      │         │ API Server   │         │ File System    │
│             ├─────────┤              ├─────────┤                │
│             │         │              │         │                │
└─────────────┘         └──────────────┘         └────────────────┘
      │                        │                         │
      │  1. Upload Request     │                         │
      │───────────────────────>│                         │
      │                        │  2. Validate File       │
      │                        │  & Product ID           │
      │                        │                         │
      │                        │  3. Save File           │
      │                        │────────────────────────>│
      │                        │                         │
      │                        │  4. Store Path in DB    │
      │                        │                         │
      │  5. Return Success     │                         │
      │<───────────────────────│                         │
      │                        │                         │
      │  6. Get Image Request  │                         │
      │───────────────────────>│                         │
      │                        │  7. Retrieve File Path  │
      │                        │                         │
      │                        │  8. Read File           │
      │                        │<────────────────────────│
      │  9. Return Image       │                         │
      │<───────────────────────│                         │
      │                        │                         │
```

## Example Requests

### Upload Product Image

```bash
curl -X POST -F "image=@/path/to/image.jpg" http://localhost:8080/products/123/upload
```

### Get Product Image

```bash
curl -X GET http://localhost:8080/products/123/image -o product_image.jpg
```

## Security Features

- MIME type validation (not just file extension)
- File size limitation (max 5MB)
- Path traversal prevention
- Secure file storage structure
- Automatic cleanup of orphaned images

## Error Handling

The API returns appropriate HTTP status codes and JSON error messages:

- 400 Bad Request: Invalid input or file type
- 404 Not Found: Product not found
- 413 Request Entity Too Large: File too large
- 500 Internal Server Error: Server-side issues

## Dependencies

- Gin Web Framework
- GORM ORM
- UUID Generator
- MySQL Driver (or your chosen database)

## Database Optimization

The project implements several MySQL-specific optimizations:

1. **Proper Indexing**: Indexes on frequently queried columns
2. **Efficient JOINs**: Using appropriate JOIN types and index hints
3. **Connection Pooling**: Configured for optimal performance
4. **Query Optimization**: Structured queries to take advantage of MySQL's query optimizer

## Scaling Considerations

For high-volume applications, consider:

1. **Read Replicas**: Set up MySQL read replicas to distribute query load
2. **Load Balancing**: Deploy multiple API instances behind a load balancer
3. **Caching**: Implement Redis caching for frequently accessed data
4. **Partitioning**: Consider table partitioning for very large datasets
