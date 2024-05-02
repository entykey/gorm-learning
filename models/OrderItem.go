package models

// OrderItem represents an order item model
type OrderItem struct {
	OrderItemId int     `gorm:"primaryKey"`
	OrderId     int     `gorm:"not null"`
	Order       Order   `gorm:"foreignKey:OrderId;references:OrderItemId"`
	ProductId   int     `gorm:"not null"`
	Product     Product `gorm:"foreignKey:ProductId;references:OrderItemId"`
	Quantity    int     `gorm:"not null"`
}

// type OrderItem struct {
// 	OrderItemId int `gorm:"primaryKey"`
// 	OrderId     int `gorm:"not null"`
// 	ProductId   int `gorm:"not null"`
// 	Quantity    int `gorm:"not null"`
// 	Product     Product
// }
