# Inventory Management System

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
/inventory-system  
  ├── database/          # Database configuration and queries
  │   └── scripts/       # SQL scripts for database setup
  ├── docs/              # documentation
  ├── handlers/          # Gin route logic  
  ├── models/            # Structs (Product, Inventory, Order)  
  ├── routes/            # Gin router groups  
  ├── uploads/           # Product images  
  ├── go.mod             # Dependencies (Gin, GORM, MySQL driver)  
  └── README.md          # This file  
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
