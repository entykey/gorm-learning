package models

import "time"

// Order represents an order model
type Order struct {
	OrderId    int         `gorm:"primary_key"`
	OrderDate  time.Time   `gorm:"not null"`
	CustomerId int         `gorm:"not null"`
	Customer   Customer    `gorm:"foreignKey:CustomerId"`
	OrderItems []OrderItem `gorm:"foreignkey:OrderItemId"` // Explicit foreign key
	// Below are fields for ViewModel processing purpose:
	Total float64 `gorm:"-"` // Add Total field without mapping to database column
}
