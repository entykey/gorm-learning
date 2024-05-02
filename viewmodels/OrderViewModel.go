package viewmodels

import "time"

// Define a view model struct for representing orders with additional fields
type OrderViewModel struct {
	OrderId    int
	OrderDate  time.Time
	CustomerId int
	Customer   struct {
		CustomerId int
		Name       string
		Email      string
	}
	OrderItems []struct {
		OrderItemId int
		ProductId   int
		Quantity    int
		Product     struct {
			ProductId int
			Name      string
			Price     float64
		}
	}
	Total float64
}
