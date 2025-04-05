package database

import (
	"fmt"
	. "gorm-learning/models"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

// only seee Customer to see if it works or not, no transaction
// STATUS: workeddocker compose up -d --build

/*
func SeedData(db *gorm.DB) error {
	// Check if table exists (GORM v1 way)
	if !db.HasTable(&Customer{}) {
		return fmt.Errorf("customers table does not exist")
	}

	// Check for existing data
	var count int
	if err := db.Model(&Customer{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count customers: %v", err)
	}

	if count > 0 {
		log.Println("Customers table already has data, skipping seeding")
		return nil
	}

	log.Println("Seeding customer data...")

	// Create customers (without orders to avoid relationship issues)
	customers := []Customer{
		{
			CustomerId: 1,
			Name:       "John Smith",
			Email:      "john.smith@example.com",
			Orders:     []Order{}, // Empty slice to avoid nil
		},
		{
			CustomerId: 2,
			Name:       "Jane Doe",
			Email:      "jane.doe@example.com",
			Orders:     []Order{},
		},
		{
			CustomerId: 3,
			Name:       "Bob Johnson",
			Email:      "bob.johnson@example.com",
			Orders:     []Order{},
		},
	}

	// Create customers
	if err := db.Create(&customers).Error; err != nil {
		return fmt.Errorf("failed to create customers: %v", err)
	}

	log.Printf("Successfully seeded %d customers", len(customers))
	return nil
}
*/

// STATUS: worked
/*
func SeedData(db *gorm.DB) error {
	// Skip if customers exist
	if db.HasTable("customers") && db.First(&Customer{}).Error == nil {
		log.Println("Customers already exist, skipping seeding")
		return nil
	}

	log.Println("Seeding customer data...")

	// Initialize with exported fields only
	customers := []*Customer{
		&Customer{
			CustomerId: 1,
			Name:       "John Smith",
			Email:      "john.smith@example.com",
		},
		&Customer{
			CustomerId: 2,
			Name:       "Jane Doe",
			Email:      "jane.doe@example.com",
		},
		&Customer{
			CustomerId: 3,
			Name:       "Bob Johnson",
			Email:      "bob.johnson@example.com",
		},
	}

	// Create customers one by one with error handling
	for _, customer := range customers {
		if err := db.Create(customer).Error; err != nil {
			return fmt.Errorf("failed to create customer %s: %v", customer.Name, err)
		}
		log.Printf("Created customer: %s", customer.Name)
	}

	return nil
}
*/

// func SeedData(db *gorm.DB) error {
// 	// Start transaction
// 	tx := db.Begin()
// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	// Defer rollback in case of panic
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// Check if data exists
// 	var count int64
// 	if err := tx.Model(&Customer{}).Count(&count).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if count > 0 {
// 		log.Println("Database already seeded, skipping")
// 		tx.Commit() // Explicitly commit read-only transaction
// 		return nil
// 	}

// 	log.Println("Not seeded. Seeding database...")

// 	// Seed Customers
// 	customers := []Customer{
// 		{CustomerId: 1, Name: "John Smith", Email: "john.smith@example.com"},
// 		{CustomerId: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
// 		{CustomerId: 3, Name: "Bob Johnson", Email: "bob.johnson@example.com"},
// 	}
// 	if err := tx.Create(&customers).Error; err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("Rolled back, failed to seed customers, error: %v", err)
// 	}
// 	fmt.Println("Successfully inserted Customers")

// 	// Seed Products
// 	products := []Product{
// 		{ProductId: 1, Name: "Product 1", Price: 9.99},
// 		{ProductId: 2, Name: "Product 2", Price: 14.99},
// 		{ProductId: 3, Name: "Product 3", Price: 19.99},
// 	}
// 	if err := tx.Create(&products).Error; err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("Rolled back, failed to seed products: %v", err)
// 	}
// 	fmt.Println("Successfully inserted Products")

// 	// Seed Orders
// 	orders := []Order{
// 		{OrderId: 1, OrderDate: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), CustomerId: customers[0].CustomerId},
// 		{OrderId: 2, OrderDate: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC), CustomerId: customers[1].CustomerId},
// 		{OrderId: 3, OrderDate: time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC), CustomerId: customers[2].CustomerId},
// 	}
// 	if err := tx.Create(&orders).Error; err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("Rolled back, failed to seed orders: %v", err)
// 	}
// 	fmt.Println("Successfully inserted Orders")

// 	// Seed OrderItems
// 	orderItems := []OrderItem{
// 		{OrderItemId: 1, OrderId: orders[0].OrderId, ProductId: products[0].ProductId, Quantity: 2},
// 		{OrderItemId: 2, OrderId: orders[0].OrderId, ProductId: products[1].ProductId, Quantity: 1},
// 		{OrderItemId: 3, OrderId: orders[1].OrderId, ProductId: products[0].ProductId, Quantity: 3},
// 		{OrderItemId: 4, OrderId: orders[1].OrderId, ProductId: products[2].ProductId, Quantity: 1},
// 		{OrderItemId: 5, OrderId: orders[2].OrderId, ProductId: products[1].ProductId, Quantity: 2},
// 		{OrderItemId: 6, OrderId: orders[2].OrderId, ProductId: products[2].ProductId, Quantity: 2},
// 	}
// 	if err := tx.Create(&orderItems).Error; err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("Rolled back, failed to seed order items: %v", err)
// 	}
// 	fmt.Println("Successfully inserted OrderItems")

// 	// Commit transaction
// 	if err := tx.Commit().Error; err != nil {
// 		return fmt.Errorf("failed to commit transaction: %v", err)
// 	}

// 	log.Println("✅ Data seeded successfully!")
// 	return nil
// }

func SeedData(db *gorm.DB) error {
	infoLogger := log.New(color.Output, color.GreenString("INFO: "), log.LstdFlags)
	errorLogger := log.New(color.Output, color.RedString("ERROR: "), log.LstdFlags)

	// Enable SQL Logging:
	db.LogMode(true) // Add this before seeding

	// Track which tables need seeding
	needsSeeding := map[string]bool{
		"customers":   true,
		"products":    true,
		"orders":      true,
		"order_items": true,
	}

	// Check each table individually
	infoLogger.Println("Checking database tables for existing data...")

	var count int64

	// Check customers
	if err := db.Model(&Customer{}).Count(&count).Error; err != nil {
		errorLogger.Printf("Failed to check customers table: %v", err)
		return err
	}
	needsSeeding["customers"] = count == 0
	infoLogger.Printf("Customers table: %d records", count)

	// Check products
	if err := db.Model(&Product{}).Count(&count).Error; err != nil {
		errorLogger.Printf("Failed to check products table: %v", err)
		return err
	}
	needsSeeding["products"] = count == 0
	infoLogger.Printf("Products table: %d records", count)

	// Check orders
	if err := db.Model(&Order{}).Count(&count).Error; err != nil {
		errorLogger.Printf("Failed to check orders table: %v", err)
		return err
	}
	needsSeeding["orders"] = count == 0
	infoLogger.Printf("Orders table: %d records", count)

	// Check order_items
	if err := db.Model(&OrderItem{}).Count(&count).Error; err != nil {
		errorLogger.Printf("Failed to check order_items table: %v", err)
		return err
	}
	needsSeeding["order_items"] = count == 0
	infoLogger.Printf("OrderItems table: %d records", count)

	// Seed only what's needed
	if needsSeeding["customers"] {
		infoLogger.Println("Seeding customers...")
		customers := []Customer{
			{CustomerId: 1, Name: "John Smith", Email: "john.smith@example.com"},
			{CustomerId: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
			{CustomerId: 3, Name: "Bob Johnson", Email: "bob.johnson@example.com"},
		}
		if err := db.Create(&customers).Error; err != nil {
			errorLogger.Printf("Failed to seed customers: %v", err)
			return fmt.Errorf("customer seeding failed: %v", err)
		}
		infoLogger.Printf("Seeded %d customers", len(customers))
	} else {
		infoLogger.Println("Skipping customers (already seeded)")
	}

	if needsSeeding["products"] {
		infoLogger.Println("Seeding products...")
		products := []Product{
			{ProductId: 1, Name: "Product 1", Price: 9.99},
			{ProductId: 2, Name: "Product 2", Price: 14.99},
			{ProductId: 3, Name: "Product 3", Price: 19.99},
		}
		// if err := db.Create(&products).Error; err != nil {
		// 	errorLogger.Printf("Failed to seed products: %v", err)
		// 	return fmt.Errorf("product seeding failed: %v", err)
		// }
		// infoLogger.Printf("Seeded %d products", len(products))

		// Create customers one by one with error handling
		for _, product := range products {
			if err := db.Create(product).Error; err != nil {
				return fmt.Errorf("failed to create customer %s: %v", product.Name, err)
			}
			log.Printf("Created product: %s", product.Name)
		}
	} else {
		infoLogger.Println("Skipping products (already seeded)")
	}

	if needsSeeding["orders"] {
		infoLogger.Println("Seeding orders...")
		orders := []Order{
			{OrderId: 1, OrderDate: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), CustomerId: 1},
			{OrderId: 2, OrderDate: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC), CustomerId: 2},
			{OrderId: 3, OrderDate: time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC), CustomerId: 3},
		}
		if err := db.Create(&orders).Error; err != nil {
			errorLogger.Printf("Failed to seed orders: %v", err)
			return fmt.Errorf("order seeding failed: %v", err)
		}
		infoLogger.Printf("Seeded %d orders", len(orders))
	} else {
		infoLogger.Println("Skipping orders (already seeded)")
	}

	if needsSeeding["order_items"] {
		infoLogger.Println("Seeding order items...")
		orderItems := []OrderItem{
			{OrderItemId: 1, OrderId: 1, ProductId: 1, Quantity: 2},
			{OrderItemId: 2, OrderId: 1, ProductId: 2, Quantity: 1},
			{OrderItemId: 3, OrderId: 2, ProductId: 1, Quantity: 3},
			{OrderItemId: 4, OrderId: 2, ProductId: 3, Quantity: 1},
			{OrderItemId: 5, OrderId: 3, ProductId: 2, Quantity: 2},
			{OrderItemId: 6, OrderId: 3, ProductId: 3, Quantity: 2},
		}
		if err := db.Create(&orderItems).Error; err != nil {
			errorLogger.Printf("Failed to seed order items: %v", err)
			return fmt.Errorf("order item seeding failed: %v", err)
		}
		infoLogger.Printf("Seeded %d order items", len(orderItems))
	} else {
		infoLogger.Println("Skipping order items (already seeded)")
	}

	infoLogger.Println("✅ Database seeding check completed")
	return nil
}
