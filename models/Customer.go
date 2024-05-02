package models

// User has many Orders, Order.CustomerID is the foreign key
type Customer struct {
	CustomerId int    `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Email      string `gorm:"not null"`
	Orders     []Order
}
