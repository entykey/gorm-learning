package models

// Product represents a product model
// type Product struct {
// 	ProductId  int         `gorm:"primaryKey"`
// 	Name       string      `gorm:"not null"`
// 	Price      float64     `gorm:"not null"`
// 	OrderItems []OrderItem `gorm:"many2many:order_items;"`
// 	// Orders []Order `gorm:"many2many:order_items;"`
// }

type Product struct {
	ProductId int     `gorm:"primary_key"`
	Name      string  `json:"name" gorm:"column:name"`   // Ensure correct column name
	Price     float64 `json:"price" gorm:"column:price"` // Ensure correct column name
}
