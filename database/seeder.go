package database

import (
	"fmt"
	. "gorm-learning/models"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

/*
Output log:
2025-04-06 06:41:30 INFO: 2025/04/05 23:41:30 Creating database gorm_learning...
2025-04-06 06:41:30 INFO: 2025/04/05 23:41:30 Applying automatic migrations...
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Migrations applied successfully!
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeding customers...
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeded 3 customers
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeding products...
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeded 3 products
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeding orders...
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeded 3 orders
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeding order items...
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 Seeded 6 order items
2025-04-06 06:41:31 INFO: 2025/04/05 23:41:31 ✅ Database seeding completed
*/
// STATUS: worked flawlessly (jinzhu/gorm + mysql 8)
// Seeder handles all database seeding operations
type Seeder struct {
	db          *gorm.DB
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// NewSeeder creates a new Seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db:          db,
		infoLogger:  log.New(color.Output, color.GreenString("INFO: "), log.LstdFlags),
		errorLogger: log.New(color.Output, color.RedString("ERROR: "), log.LstdFlags),
	}
}

// SeedAll runs all seeding operations
func (s *Seeder) SeedAll() error {

	// Enable SQL Logging:
	s.db.LogMode(true) // Add this before seeding

	if err := s.SeedCustomers(); err != nil {
		return err
	}
	if err := s.SeedProducts(); err != nil {
		return err
	}
	if err := s.SeedOrders(); err != nil {
		return err
	}
	return s.SeedOrderItems()
}

// SeedCustomers handles customer data seeding
func (s *Seeder) SeedCustomers() error {
	if !s.db.HasTable(&Customer{}) {
		return fmt.Errorf("customers table does not exist")
	}

	var count int64
	s.db.Model(&Customer{}).Count(&count)

	if count > 0 {
		s.infoLogger.Printf("Skipping customers - %d records exist", count)
		return nil
	}

	s.infoLogger.Println("Seeding customers...")
	customers := []*Customer{
		{CustomerId: 1, Name: "John Smith", Email: "john.smith@example.com"},
		{CustomerId: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
		{CustomerId: 3, Name: "Bob Johnson", Email: "bob.johnson@example.com"},
	}

	for _, c := range customers {
		if err := s.db.Create(c).Error; err != nil {
			s.errorLogger.Printf("Failed to create customer %s: %v", c.Name, err)
			return fmt.Errorf("customer creation failed: %v", err)
		}
	}
	s.infoLogger.Printf("Seeded %d customers", len(customers))
	return nil
}

// SeedProducts handles product data seeding
func (s *Seeder) SeedProducts() error {
	if !s.db.HasTable(&Product{}) {
		return fmt.Errorf("products table does not exist")
	}

	var count int64
	s.db.Model(&Product{}).Count(&count)

	if count > 0 {
		s.infoLogger.Printf("Skipping products - %d records exist", count)
		return nil
	}

	s.infoLogger.Println("Seeding products...")
	products := []*Product{
		{ProductId: 1, Name: "Product 1", Price: 9.99},
		{ProductId: 2, Name: "Product 2", Price: 14.99},
		{ProductId: 3, Name: "Product 3", Price: 19.99},
	}

	for _, p := range products {
		if err := s.db.Create(p).Error; err != nil {
			s.errorLogger.Printf("Failed to create product %s: %v", p.Name, err)
			return fmt.Errorf("product creation failed: %v", err)
		}
	}
	s.infoLogger.Printf("Seeded %d products", len(products))
	return nil
}

// SeedOrders handles order data seeding
func (s *Seeder) SeedOrders() error {
	if !s.db.HasTable(&Order{}) {
		return fmt.Errorf("orders table does not exist")
	}

	var count int64
	s.db.Model(&Order{}).Count(&count)

	if count > 0 {
		s.infoLogger.Printf("Skipping orders - %d records exist", count)
		return nil
	}

	// Verify prerequisites
	var customerCount, productCount int64
	s.db.Model(&Customer{}).Count(&customerCount)
	s.db.Model(&Product{}).Count(&productCount)

	if customerCount == 0 || productCount == 0 {
		return fmt.Errorf("cannot seed orders - need customers and products first")
	}

	s.infoLogger.Println("Seeding orders...")
	orders := []*Order{
		{OrderId: 1, OrderDate: time.Now(), CustomerId: 1},
		{OrderId: 2, OrderDate: time.Now(), CustomerId: 2},
		{OrderId: 3, OrderDate: time.Now(), CustomerId: 3},
	}

	for _, o := range orders {
		if err := s.db.Create(o).Error; err != nil {
			s.errorLogger.Printf("Failed to create order %d: %v", o.OrderId, err)
			return fmt.Errorf("order creation failed: %v", err)
		}
	}
	s.infoLogger.Printf("Seeded %d orders", len(orders))
	return nil
}

// SeedOrderItems handles order item data seeding
func (s *Seeder) SeedOrderItems() error {
	if !s.db.HasTable(&OrderItem{}) {
		return fmt.Errorf("order_items table does not exist")
	}

	var count int64
	s.db.Model(&OrderItem{}).Count(&count)

	if count > 0 {
		s.infoLogger.Printf("Skipping order items - %d records exist", count)
		return nil
	}

	// Verify prerequisites
	var orderCount, productCount int64
	s.db.Model(&Order{}).Count(&orderCount)
	s.db.Model(&Product{}).Count(&productCount)

	if orderCount == 0 || productCount == 0 {
		return fmt.Errorf("cannot seed order items - need orders and products first")
	}

	s.infoLogger.Println("Seeding order items...")
	orderItems := []*OrderItem{
		{OrderItemId: 1, OrderId: 1, ProductId: 1, Quantity: 2},
		{OrderItemId: 2, OrderId: 1, ProductId: 2, Quantity: 1},
		{OrderItemId: 3, OrderId: 2, ProductId: 1, Quantity: 3},
		{OrderItemId: 4, OrderId: 2, ProductId: 3, Quantity: 1},
		{OrderItemId: 5, OrderId: 3, ProductId: 2, Quantity: 2},
		{OrderItemId: 6, OrderId: 3, ProductId: 3, Quantity: 2},
	}

	for _, oi := range orderItems {
		if err := s.db.Create(oi).Error; err != nil {
			s.errorLogger.Printf("Failed to create order item %d: %v", oi.OrderItemId, err)
			return fmt.Errorf("order item creation failed: %v", err)
		}
	}
	s.infoLogger.Printf("Seeded %d order items", len(orderItems))
	return nil
}

// only seee Customer to see if it works or not, no transaction
// STATUS: faield
/*
2025-04-06 01:43:08 INFO: 2025/04/05 18:43:08 Creating database gorm_learning...
2025-04-06 01:43:08 INFO: 2025/04/05 18:43:08 Applying automatic migrations...
2025-04-06 01:43:08 INFO: 2025/04/05 18:43:08 Migrations applied successfully!
2025-04-06 01:43:08 INFO: 2025/04/05 18:43:08 Starting database seeding...
2025-04-06 01:44:23 INFO: 2025/04/05 18:44:23 Creating database gorm_learning...
2025-04-06 01:44:23 INFO: 2025/04/05 18:44:23 Applying automatic migrations...
2025-04-06 01:43:08 2025/04/05 18:43:08 Seeding customer data...
2025-04-06 01:43:08 panic: reflect: call of reflect.Value.Interface on zero Value [recovered]
2025-04-06 01:43:08     panic: reflect: call of reflect.Value.Interface on zero Value
2025-04-06 01:43:08
2025-04-06 01:43:08 goroutine 1 [running]:
2025-04-06 01:43:08 github.com/jinzhu/gorm.(*Scope).callCallbacks.func1()
2025-04-06 01:43:08     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:865 +0x6b
2025-04-06 01:43:08 panic({0x89bfa0?, 0xc000233158?})
2025-04-06 01:43:08     /usr/local/go/src/runtime/panic.go:914 +0x21f
2025-04-06 01:43:08 reflect.valueInterface({0x0?, 0x0?, 0x4?}, 0xd8?)
2025-04-06 01:43:08     /usr/local/go/src/reflect/value.go:1495 +0xfb
2025-04-06 01:43:08 reflect.Value.Interface(...)
2025-04-06 01:43:08     /usr/local/go/src/reflect/value.go:1490
2025-04-06 01:43:08 github.com/jinzhu/gorm.createCallback(0xc000226900)
2025-04-06 01:43:08     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/callback_create.go:68 +0x46f
2025-04-06 01:43:08 github.com/jinzhu/gorm.(*Scope).callCallbacks(0xc000226900, {0xc000182400, 0x9, 0xc00020fa01?})
2025-04-06 01:43:08     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:869 +0x8d
2025-04-06 01:43:08 github.com/jinzhu/gorm.(*DB).Create(0xc000230a90, {0x87cac0?, 0xc0002330b0?})
2025-04-06 01:43:08     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/main.go:483 +0x49
2025-04-06 01:43:08 gorm-learning/database.SeedData(0xc0001951a0?)
2025-04-06 01:43:08     /app/database/seeder.go:57 +0x2b7
2025-04-06 01:43:08 main.initORM_DB()
2025-04-06 01:43:08     /app/main.go:93 +0xac7
2025-04-06 01:43:08 main.main()
2025-04-06 01:43:08     /app/main.go:211 +0x27
2025-04-06 01:44:23 INFO: 2025/04/05 18:44:23 Migrations applied successfully!
2025-04-06 01:44:23 INFO: 2025/04/05 18:44:23 Starting database seeding...
2025-04-06 01:44:23 2025/04/05 18:44:23 Seeding customer data...
2025-04-06 01:44:23 panic: reflect: call of reflect.Value.Interface on zero Value [recovered]
2025-04-06 01:44:23     panic: reflect: call of reflect.Value.Interface on zero Value
2025-04-06 01:44:23
2025-04-06 01:44:23 goroutine 1 [running]:
2025-04-06 01:44:23 github.com/jinzhu/gorm.(*Scope).callCallbacks.func1()
2025-04-06 01:44:23     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:865 +0x6b
2025-04-06 01:44:23 panic({0x89bfa0?, 0xc000276000?})
2025-04-06 01:44:23     /usr/local/go/src/runtime/panic.go:914 +0x21f
2025-04-06 01:44:23 reflect.valueInterface({0x0?, 0x0?, 0x4?}, 0x30?)
2025-04-06 01:44:23     /usr/local/go/src/reflect/value.go:1495 +0xfb
2025-04-06 01:44:23 reflect.Value.Interface(...)
2025-04-06 01:44:23     /usr/local/go/src/reflect/value.go:1490
2025-04-06 01:44:23 github.com/jinzhu/gorm.createCallback(0xc000203580)
2025-04-06 01:44:23     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/callback_create.go:68 +0x46f
2025-04-06 01:44:23 github.com/jinzhu/gorm.(*Scope).callCallbacks(0xc000203580, {0xc000202400, 0x9, 0xc00028ba01?})
2025-04-06 01:44:23     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:869 +0x8d
2025-04-06 01:44:23 github.com/jinzhu/gorm.(*DB).Create(0xc000101ad0, {0x87cac0?, 0xc000011f50?})
2025-04-06 01:44:23     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/main.go:483 +0x49
2025-04-06 01:44:23 gorm-learning/database.SeedData(0xc0002151a0?)
2025-04-06 01:44:23     /app/database/seeder.go:57 +0x2b7
2025-04-06 01:44:23 main.initORM_DB()
2025-04-06 01:44:23     /app/main.go:93 +0xac7
2025-04-06 01:44:23 main.main()
2025-04-06 01:44:23     /app/main.go:211 +0x27
*/
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

/*
2025-04-06 01:46:24 INFO: 2025/04/05 18:46:24 Applying automatic migrations...
2025-04-06 01:46:24 INFO: 2025/04/05 18:46:24 Migrations applied successfully!
2025-04-06 01:46:24 INFO: 2025/04/05 18:46:24 Starting database seeding...
2025-04-06 01:46:24
2025-04-06 01:46:24     _______ __
2025-04-06 01:46:24    / ____(_) /_  ___  _____
2025-04-06 01:46:24   / /_  / / __ \/ _ \/ ___/
2025-04-06 01:46:24  / __/ / / /_/ /  __/ /
2025-04-06 01:46:24 /_/   /_/_.___/\___/_/          v3.0.0-beta.2
2025-04-06 01:46:24 --------------------------------------------------
2025-04-06 01:46:24 INFO Server started on:     http://127.0.0.1:3000 (bound on host 0.0.0.0 and port 3000)
2025-04-06 01:46:24 INFO Total handlers count:  7
2025-04-06 01:46:24 INFO Prefork:                       Disabled
2025-04-06 01:46:24 INFO PID:                   1
2025-04-06 01:46:24 INFO Total process count:   1
2025-04-06 01:46:24
2025-04-06 01:46:24 2025/04/05 18:46:24 Seeding customer data...
2025-04-06 01:46:24 2025/04/05 18:46:24 Created customer: John Smith
2025-04-06 01:46:24 2025/04/05 18:46:24 Created customer: Jane Doe
2025-04-06 01:46:24 2025/04/05 18:46:24 Created customer: Bob Johnson
*/
// STATUS: worked
// /*
// func SeedData(db *gorm.DB) error {
// 	// Skip if customers exist
// 	if db.HasTable("customers") && db.First(&Customer{}).Error == nil {
// 		log.Println("Customers already exist, skipping seeding")
// 		return nil
// 	}

// 	log.Println("Seeding customer data...")

// 	// Initialize with exported fields only
// 	customers := []*Customer{
// 		&Customer{
// 			CustomerId: 1,
// 			Name:       "John Smith",
// 			Email:      "john.smith@example.com",
// 		},
// 		&Customer{
// 			CustomerId: 2,
// 			Name:       "Jane Doe",
// 			Email:      "jane.doe@example.com",
// 		},
// 		&Customer{
// 			CustomerId: 3,
// 			Name:       "Bob Johnson",
// 			Email:      "bob.johnson@example.com",
// 		},
// 	}

// 	// Create customers one by one with error handling
// 	for _, customer := range customers {
// 		if err := db.Create(customer).Error; err != nil {
// 			return fmt.Errorf("failed to create customer %s: %v", customer.Name, err)
// 		}
// 		log.Printf("Created customer: %s", customer.Name)
// 	}

// 	// Skip if products exist
// 	if db.HasTable("products") && db.First(&Product{}).Error == nil {
// 		log.Println("Products already exist, skipping seeding")
// 		return nil
// 	}

// 	log.Println("Seeding product data...")

// 	// Initialize with exported fields only
// 	products := []Product{
// 		{ProductId: 1, Name: "Product 1", Price: 9.99},
// 		{ProductId: 2, Name: "Product 2", Price: 14.99},
// 		{ProductId: 3, Name: "Product 3", Price: 19.99},
// 	}

// 	// Create customers one by one with error handling
// 	for _, product := range products {
// 		if err := db.Create(product).Error; err != nil {
// 			return fmt.Errorf("failed to create product %s: %v", product.Name, err)
// 		}
// 		log.Printf("Created product: %s", product.Name)
// 	}

// 	return nil
// }

/*
// STATUS: Customers & Products are seeded fine now
// fix to not skip Product just because Customer exist:
func SeedData(db *gorm.DB) error {
	infoLogger := log.New(color.Output, color.GreenString("INFO: "), log.LstdFlags)
	errorLogger := log.New(color.Output, color.RedString("ERROR: "), log.LstdFlags)

	// Seed Customers (independent check)
	if db.HasTable("customers") {
		var customerCount int64
		db.Model(&Customer{}).Count(&customerCount)

		if customerCount == 0 {
			infoLogger.Println("Seeding customers...")
			customers := []*Customer{
				{
					CustomerId: 1,
					Name:       "John Smith",
					Email:      "john.smith@example.com",
				},
				{
					CustomerId: 2,
					Name:       "Jane Doe",
					Email:      "jane.doe@example.com",
				},
				{
					CustomerId: 3,
					Name:       "Bob Johnson",
					Email:      "bob.johnson@example.com",
				},
			}

			for _, customer := range customers {
				if err := db.Create(customer).Error; err != nil {
					errorLogger.Printf("Failed to create customer %s: %v", customer.Name, err)
					return fmt.Errorf("customer creation failed: %v", err)
				}
				infoLogger.Printf("Created customer: %s", customer.Name)
			}
		} else {
			infoLogger.Printf("Skipping customer seeding - %d customers already exist", customerCount)
		}
	}

	// Seed Products (independent check)
	if db.HasTable("products") {
		var productCount int64
		db.Model(&Product{}).Count(&productCount)

		if productCount == 0 {
			infoLogger.Println("Seeding products...")
			products := []Product{
				{ProductId: 1, Name: "Product 1", Price: 9.99},
				{ProductId: 2, Name: "Product 2", Price: 14.99},
				{ProductId: 3, Name: "Product 3", Price: 19.99},
			}

			for _, product := range products {
				if err := db.Create(&product).Error; err != nil {
					errorLogger.Printf("Failed to create product %s: %v", product.Name, err)
					return fmt.Errorf("product creation failed: %v", err)
				}
				infoLogger.Printf("Created product: %s", product.Name)
			}
		} else {
			infoLogger.Printf("Skipping product seeding - %d products already exist", productCount)
		}
	}

	// Seed Orders and OrderItems (if needed)
	if db.HasTable("orders") {
		var orderCount int64
		db.Model(&Order{}).Count(&orderCount)

		if orderCount == 0 {
			infoLogger.Println("Seeding orders...")
			// Verify we have customers and products first
			var customerCount, productCount int64
			db.Model(&Customer{}).Count(&customerCount)
			db.Model(&Product{}).Count(&productCount)

			if customerCount == 0 || productCount == 0 {
				errorLogger.Println("Cannot seed orders - missing customers or products")
				return nil
			}

			orders := []Order{
				{OrderId: 1, OrderDate: time.Now(), CustomerId: 1},
				{OrderId: 2, OrderDate: time.Now(), CustomerId: 2},
				{OrderId: 3, OrderDate: time.Now(), CustomerId: 3},
			}

			for _, order := range orders {
				if err := db.Create(&order).Error; err != nil {
					errorLogger.Printf("Failed to create order %d: %v", order.OrderId, err)
					continue // Skip failed orders but continue trying others
				}
			}

			// Seed order items
			infoLogger.Println("Seeding order items...")
			orderItems := []OrderItem{
				{OrderItemId: 1, OrderId: 1, ProductId: 1, Quantity: 2},
				{OrderItemId: 2, OrderId: 1, ProductId: 2, Quantity: 1},
				{OrderItemId: 3, OrderId: 2, ProductId: 1, Quantity: 3},
				{OrderItemId: 4, OrderId: 2, ProductId: 3, Quantity: 1},
				{OrderItemId: 5, OrderId: 3, ProductId: 2, Quantity: 2},
				{OrderItemId: 6, OrderId: 3, ProductId: 3, Quantity: 2},
			}

			for _, item := range orderItems {
				if err := db.Create(&item).Error; err != nil {
					errorLogger.Printf("Failed to create order item %d: %v", item.OrderItemId, err)
				}
			}
		}
	}

	infoLogger.Println("✅ Seeding process completed")
	return nil
}
*/

// */

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

/*
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
*/

/*
Output:
2025-04-06 01:38:50 INFO: 2025/04/05 18:38:50 Applying automatic migrations...
2025-04-06 01:38:50 INFO: 2025/04/05 18:38:50 Migrations applied successfully!
2025-04-06 01:38:50 INFO: 2025/04/05 18:38:50 Starting database seeding...
2025-04-06 01:40:01 INFO: 2025/04/05 18:40:01 Creating database gorm_learning...
2025-04-06 01:40:01 INFO: 2025/04/05 18:40:01 Applying automatic migrations...
2025-04-06 01:38:50 2025/04/05 18:38:50 Orders table is empty, seeding data.
2025-04-06 01:38:50 panic: reflect: call of reflect.Value.Interface on zero Value [recovered]
2025-04-06 01:38:50     panic: reflect: call of reflect.Value.Interface on zero Value
2025-04-06 01:38:50
2025-04-06 01:38:50 goroutine 1 [running]:
2025-04-06 01:38:50 github.com/jinzhu/gorm.(*Scope).callCallbacks.func1()
2025-04-06 01:38:50     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:865 +0x6b
2025-04-06 01:38:50 panic({0x89d1a0?, 0xc00009fae8?})
2025-04-06 01:38:50     /usr/local/go/src/runtime/panic.go:914 +0x21f
2025-04-06 01:38:50 reflect.valueInterface({0x0?, 0x0?, 0x4?}, 0x8?)
2025-04-06 01:38:50     /usr/local/go/src/reflect/value.go:1495 +0xfb
2025-04-06 01:38:50 reflect.Value.Interface(...)
2025-04-06 01:38:50     /usr/local/go/src/reflect/value.go:1490
2025-04-06 01:38:50 github.com/jinzhu/gorm.createCallback(0xc0000a2680)
2025-04-06 01:38:50     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/callback_create.go:68 +0x46f
2025-04-06 01:38:50 github.com/jinzhu/gorm.(*Scope).callCallbacks(0xc0000a2680, {0xc000202400, 0x9, 0xc0000c5a01?})
2025-04-06 01:38:50     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:869 +0x8d
2025-04-06 01:38:50 github.com/jinzhu/gorm.(*DB).Create(0xc0002a88f0, {0x87dba0?, 0xc00009fa40?})
2025-04-06 01:38:50     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/main.go:483 +0x49
2025-04-06 01:38:50 gorm-learning/database.SeedData(0xc0002151a0?)
2025-04-06 01:38:50     /app/database/seeder.go:347 +0x246
2025-04-06 01:38:50 main.initORM_DB()
2025-04-06 01:38:50     /app/main.go:93 +0xac7
2025-04-06 01:38:50 main.main()
2025-04-06 01:38:50     /app/main.go:211 +0x27
2025-04-06 01:40:01 INFO: 2025/04/05 18:40:01 Migrations applied successfully!
2025-04-06 01:40:01 INFO: 2025/04/05 18:40:01 Starting database seeding...
2025-04-06 01:40:01 2025/04/05 18:40:01 Orders table is empty, seeding data.
2025-04-06 01:40:01 panic: reflect: call of reflect.Value.Interface on zero Value [recovered]
2025-04-06 01:40:01     panic: reflect: call of reflect.Value.Interface on zero Value
2025-04-06 01:40:01
2025-04-06 01:40:01 goroutine 1 [running]:
2025-04-06 01:40:01 github.com/jinzhu/gorm.(*Scope).callCallbacks.func1()
2025-04-06 01:40:01     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:865 +0x6b
2025-04-06 01:40:01 panic({0x89d1a0?, 0xc000010e88?})
2025-04-06 01:40:01     /usr/local/go/src/runtime/panic.go:914 +0x21f
2025-04-06 01:40:01 reflect.valueInterface({0x0?, 0x0?, 0x4?}, 0x30?)
2025-04-06 01:40:01     /usr/local/go/src/reflect/value.go:1495 +0xfb
2025-04-06 01:40:01 reflect.Value.Interface(...)
2025-04-06 01:40:01     /usr/local/go/src/reflect/value.go:1490
2025-04-06 01:40:01 github.com/jinzhu/gorm.createCallback(0xc000024700)
2025-04-06 01:40:01     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/callback_create.go:68 +0x46f
2025-04-06 01:40:01 github.com/jinzhu/gorm.(*Scope).callCallbacks(0xc000024700, {0xc0001b6380, 0x9, 0xc000317a01?})
2025-04-06 01:40:01     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/scope.go:869 +0x8d
2025-04-06 01:40:01 github.com/jinzhu/gorm.(*DB).Create(0xc00009dad0, {0x87dba0?, 0xc000010de0?})
2025-04-06 01:40:01     /go/pkg/mod/github.com/jinzhu/gorm@v1.9.16/main.go:483 +0x49
2025-04-06 01:40:01 gorm-learning/database.SeedData(0xc0001c91a0?)
2025-04-06 01:40:01     /app/database/seeder.go:347 +0x246
2025-04-06 01:40:01 main.initORM_DB()
2025-04-06 01:40:01     /app/main.go:93 +0xac7
2025-04-06 01:40:01 main.main()
2025-04-06 01:40:01     /app/main.go:211 +0x27
*/
/*
func SeedData(db *gorm.DB) error {
	var count int64
	db.Model(&Order{}).Count(&count)
	if count > 0 {
		// If the table is not empty, do not seed
		log.Printf("Orders table is not empty, skipping seeding.")
		return nil
	}
	log.Println("Orders table is empty, seeding data.")

	// Seed Customers
	customers := []Customer{
		{CustomerId: 1, Name: "John Smith", Email: "john.smith@example.com"},
		{CustomerId: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
		{CustomerId: 3, Name: "Bob Johnson", Email: "bob.johnson@example.com"},
	}
	db.Create(&customers)

	// Seed Products
	products := []Product{
		{ProductId: 1, Name: "Product 1", Price: 9.99},
		{ProductId: 2, Name: "Product 2", Price: 14.99},
		{ProductId: 3, Name: "Product 3", Price: 19.99},
	}
	db.Create(&products)

	// Seed Orders
	orders := []Order{
		{OrderId: 1, OrderDate: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), CustomerId: 1},
		{OrderId: 2, OrderDate: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC), CustomerId: 2},
		{OrderId: 3, OrderDate: time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC), CustomerId: 3},
	}
	db.Create(&orders)

	// Seed OrderItems
	orderItems := []OrderItem{
		{OrderItemId: 1, OrderId: 1, ProductId: 1, Quantity: 2},
		{OrderItemId: 2, OrderId: 1, ProductId: 2, Quantity: 1},
		{OrderItemId: 3, OrderId: 2, ProductId: 1, Quantity: 3},
		{OrderItemId: 4, OrderId: 2, ProductId: 3, Quantity: 1},
		{OrderItemId: 5, OrderId: 3, ProductId: 2, Quantity: 2},
		{OrderItemId: 6, OrderId: 3, ProductId: 3, Quantity: 2},
	}
	db.Create(&orderItems)

	return nil
}
*/
