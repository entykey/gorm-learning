package database

import (
	. "gorm-learning/models"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

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
