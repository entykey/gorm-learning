package models

// User has many Orders, Order.CustomerID is the foreign key
// type Customer struct {
// 	CustomerId int    `gorm:"primaryKey"`
// 	Name       string `gorm:"not null"`
// 	Email      string `gorm:"not null"`
// 	Orders     []Order
// }

// modify Customer model to make it more seeding-friendly:
type Customer struct {
	CustomerId int `gorm:"primary_key"`
	Name       string
	Email      string
	Orders     []Order `gorm:"foreignkey:CustomerId"` // Explicit foreign key
}
