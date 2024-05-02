package controllers

import (
	"database/sql"
	"fmt"
	"time"

	// . "gorm-learning/models"

	. "gorm-learning/viewmodels"

	"log"
	// "log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// OrderController handles HTTP requests related to orders
type OrderController struct {
	ormdb *gorm.DB
	rawdb *sql.DB
}

// NewOrderController creates a new instance of OrderController
func NewOrderController(ormdb *gorm.DB, rawdb *sql.DB) *OrderController {
	return &OrderController{ormdb: ormdb, rawdb: rawdb}
}

// GetOrdersSimple retrieves orders with specific fields from the database
func (oc *OrderController) GetOrdersSimpleAnonymousObj(c fiber.Ctx) error {
	// Fetch data from the database and map it to the custom anonymous object
	rows, err := oc.ormdb.Table("orders").
		Select("orders.order_id, orders.order_date, orders.customer_id").
		Rows()

	if err != nil {
		errorMessage := fmt.Sprintf("Error querying database: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}
	defer rows.Close()

	// Define a custom anonymous object to collect the fetched columns
	type OrderInfo struct {
		OrderId    int       `json:"orderId"`
		OrderDate  time.Time `json:"orderDate"`
		CustomerId int       `json:"customerId"`
	}

	var orders []OrderInfo
	for rows.Next() {
		var order OrderInfo

		err := rows.Scan(&order.OrderId, &order.OrderDate, &order.CustomerId)
		if err != nil {
			errorMessage := fmt.Sprintf("Error scanning row: %s", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		errorMessage := fmt.Sprintf("Error iterating over rows: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Return the result as JSON
	return c.JSON(orders)
}

/*
func (oc *OrderController) GetOrdersSimple(c fiber.Ctx) error {
	// Fetch data from the database and map it to the custom anonymous object
	rows, err := oc.ormdb.Table("orders").
		Select("orders.order_id, orders.order_date, orders.customer_id").
		Rows()

	if err != nil {
		errorMessage := fmt.Sprintf("Error querying database: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order

		err := rows.Scan(&order.OrderId, &order.OrderDate, &order.CustomerId)
		if err != nil {
			errorMessage := fmt.Sprintf("Error scanning row: %s", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		errorMessage := fmt.Sprintf("Error iterating over rows: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Return the result as JSON
	return c.JSON(orders)
}
*/

/*
	func (oc *OrderController) RetrieveOrdersIncludeAllPropsGorm(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v", r)
				c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			}
		}()

		// Ensure db is not nil
		if oc.ormdb == nil {
			log.Println("Database connection is nil")
			return c.Status(http.StatusInternalServerError).SendString("Database connection is nil")
		}

		// Fetch data from the database and map it to the view model
		rows, err := oc.ormdb.Table("orders").
			Select("orders.order_id, orders.order_date, orders.customer_id, customers.customer_id, customers.name, customers.email, order_items.order_item_id, order_items.product_id, order_items.quantity, products.product_id, products.name, products.price").
			Joins("left join customers on orders.customer_id = customers.customer_id").
			Joins("left join order_items on orders.order_id = order_items.order_id").
			Joins("left join products on order_items.product_id = products.product_id").
			Rows()

		if err != nil {
			errorMessage := fmt.Sprintf("Error querying database: %s", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}
		defer rows.Close()

		var orders []OrderViewModel
		for rows.Next() {
			var order OrderViewModel
			// var customer Customer
			var customer struct {
				CustomerId int
				Name       string
				Email      string
			}
			var orderItem struct {
				OrderItemId int
				ProductId   int
				Quantity    int
				Product     struct {
					ProductId int
					Name      string
					Price     float64
				}
			}
			var product struct {
				ProductId int
				Name      string
				Price     float64
			}
			// var product Product

			err := rows.Scan(&order.OrderId, &order.OrderDate, &order.CustomerId, &customer.CustomerId, &customer.Name, &customer.Email, &orderItem.OrderItemId, &orderItem.ProductId, &orderItem.Quantity, &product.ProductId, &product.Name, &product.Price)
			if err != nil {
				errorMessage := fmt.Sprintf("Error scanning row: %s", err)
				log.Println(errorMessage)
				return c.Status(http.StatusInternalServerError).SendString(errorMessage)
			}

			order.Customer = customer
			orderItem.Product = product
			order.OrderItems = append(order.OrderItems, orderItem)
			orders = append(orders, order)
		}

		if err := rows.Err(); err != nil {
			errorMessage := fmt.Sprintf("Error iterating over rows: %s", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		// Calculate total for each order
		for i := range orders {
			var total float64
			for _, item := range orders[i].OrderItems {
				total += float64(item.Quantity) * item.Product.Price
			}
			orders[i].Total = total
		}

		// Return the result
		return c.JSON(orders)
	}
*/

// status: MissingProductMapping
func (oc *OrderController) MissingProductMappingRetrieveOrdersIncludeAllPropsGorm(c fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v", r)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}()

	// Ensure db is not nil
	if oc.ormdb == nil {
		log.Println("Database connection is nil")
		return c.Status(http.StatusInternalServerError).SendString("Database connection is nil")
	}

	// Fetch data from the database and map it to the view model
	rows, err := oc.ormdb.Table("orders").
		Select("orders.order_id, orders.order_date, orders.customer_id, customers.customer_id, customers.name, customers.email, order_items.order_item_id, order_items.product_id, order_items.quantity, products.product_id, products.name, products.price").
		Joins("left join customers on orders.customer_id = customers.customer_id").
		Joins("left join order_items on orders.order_id = order_items.order_id").
		Joins("left join products on order_items.product_id = products.product_id").
		Rows()

	if err != nil {
		errorMessage := fmt.Sprintf("Error querying database: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}
	defer rows.Close()

	var ordersMap = make(map[int]*OrderViewModel)
	for rows.Next() {
		var order OrderViewModel
		var customer struct {
			CustomerId int
			Name       string
			Email      string
		}
		var orderItem struct {
			OrderItemId int
			ProductId   int
			Quantity    int
			Product     struct {
				ProductId int
				Name      string
				Price     float64
			}
		}
		var product struct {
			ProductId int
			Name      string
			Price     float64
		}

		err := rows.Scan(&order.OrderId, &order.OrderDate, &order.CustomerId, &customer.CustomerId, &customer.Name, &customer.Email, &orderItem.OrderItemId, &orderItem.ProductId, &orderItem.Quantity, &product.ProductId, &product.Name, &product.Price)
		if err != nil {
			errorMessage := fmt.Sprintf("Error scanning row: %s", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		if _, ok := ordersMap[order.OrderId]; !ok {
			order.Customer = customer
			ordersMap[order.OrderId] = &order
		}
		ordersMap[order.OrderId].OrderItems = append(ordersMap[order.OrderId].OrderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		errorMessage := fmt.Sprintf("Error iterating over rows: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Convert map to slice
	var orders []OrderViewModel
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	// Calculate total for each order
	for i := range orders {
		var total float64
		for _, item := range orders[i].OrderItems {
			total += float64(item.Quantity) * item.Product.Price
		}
		orders[i].Total = total
	}

	// Return the result
	return c.JSON(orders)
}

func (oc *OrderController) WorkedRetrieveOrdersIncludeAllPropsGorm(c fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v", r)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}()

	// Ensure db is not nil
	if oc.ormdb == nil {
		log.Println("Database connection is nil")
		return c.Status(http.StatusInternalServerError).SendString("Database connection is nil")
	}

	// Fetch data from the database and map it to the view model
	rows, err := oc.ormdb.Table("orders").
		Select("orders.order_id, orders.order_date, orders.customer_id, customers.customer_id, customers.name, customers.email, order_items.order_item_id, order_items.product_id, order_items.quantity, products.product_id, products.name, products.price").
		Joins("left join customers on orders.customer_id = customers.customer_id").
		Joins("left join order_items on orders.order_id = order_items.order_id").
		Joins("left join products on order_items.product_id = products.product_id").
		Rows()

	if err != nil {
		errorMessage := fmt.Sprintf("Error querying database: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}
	defer rows.Close()

	var ordersMap = make(map[int]*OrderViewModel)
	for rows.Next() {
		var order OrderViewModel
		var customer struct {
			CustomerId int
			Name       string
			Email      string
		}
		var orderItem struct {
			OrderItemId int
			ProductId   int
			Quantity    int
			Product     struct {
				ProductId int
				Name      string
				Price     float64
			}
		}

		err := rows.Scan(&order.OrderId, &order.OrderDate, &order.CustomerId, &customer.CustomerId, &customer.Name, &customer.Email, &orderItem.OrderItemId, &orderItem.ProductId, &orderItem.Quantity, &orderItem.Product.ProductId, &orderItem.Product.Name, &orderItem.Product.Price)
		if err != nil {
			errorMessage := fmt.Sprintf("Error scanning row: %s", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		if _, ok := ordersMap[order.OrderId]; !ok {
			order.Customer = customer
			ordersMap[order.OrderId] = &order
		}
		ordersMap[order.OrderId].OrderItems = append(ordersMap[order.OrderId].OrderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		errorMessage := fmt.Sprintf("Error iterating over rows: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Convert map to slice
	var orders []OrderViewModel
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	// Calculate total for each order
	for i := range orders {
		var total float64
		for _, item := range orders[i].OrderItems {
			total += float64(item.Quantity) * item.Product.Price
		}
		orders[i].Total = total
	}

	// Return the result
	return c.JSON(orders)
}

// mapped correctly but wong thing is: each all order only has exactly 1 orderItem => wrong calculation
func (oc *OrderController) WrongRetrieveOrdersIncludeAllPropsGorm(c fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v", r)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}()

	// Ensure db is not nil
	if oc.ormdb == nil {
		log.Println("Database connection is nil")
		return c.Status(http.StatusInternalServerError).SendString("Database connection is nil")
	}

	// Fetch data from the database and map it to the view model
	rows, err := oc.ormdb.Table("orders").
		Select("orders.order_id, orders.order_date, orders.customer_id, customers.customer_id, customers.name, customers.email, order_items.order_item_id, order_items.product_id, order_items.quantity, products.product_id, products.name, products.price").
		Joins("left join customers on orders.customer_id = customers.customer_id").
		Joins("left join order_items on orders.order_id = order_items.order_id").
		Joins("left join products on order_items.product_id = products.product_id").
		Rows()

	if err != nil {
		errorMessage := fmt.Sprintf("Error querying database: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}
	defer rows.Close()

	var orders []OrderViewModel
	for rows.Next() {
		var order OrderViewModel
		// var customer Customer
		var customer struct {
			CustomerId int
			Name       string
			Email      string
		}
		var orderItem struct {
			OrderItemId int
			ProductId   int
			Quantity    int
			Product     struct {
				ProductId int
				Name      string
				Price     float64
			}
		}
		var product struct {
			ProductId int
			Name      string
			Price     float64
		}
		// var product Product

		err := rows.Scan(&order.OrderId, &order.OrderDate, &order.CustomerId, &customer.CustomerId, &customer.Name, &customer.Email, &orderItem.OrderItemId, &orderItem.ProductId, &orderItem.Quantity, &product.ProductId, &product.Name, &product.Price)
		if err != nil {
			errorMessage := fmt.Sprintf("Error scanning row: %s", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		order.Customer = customer
		orderItem.Product = product
		order.OrderItems = append(order.OrderItems, orderItem)
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		errorMessage := fmt.Sprintf("Error iterating over rows: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Calculate total for each order
	for i := range orders {
		var total float64
		for _, item := range orders[i].OrderItems {
			total += float64(item.Quantity) * item.Product.Price
		}
		orders[i].Total = total
	}

	// Return the result
	return c.JSON(orders)
}

// Handler to Get orders (joined props) using raw SQL
func (oc *OrderController) GetOrdersWithRawSQL(c fiber.Ctx) error {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v", r)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}()

	// Ensure db is not nil
	if oc.rawdb == nil {
		log.Println("Database connection is nil")
		return c.Status(http.StatusInternalServerError).SendString("Database connection is nil")
	}

	// Defer function to catch panics during database query
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic during database query: %v", r)
			c.Status(http.StatusInternalServerError).SendString("Error querying database")
		}
	}()

	// Define the SQL query
	query := `
		SELECT
			o.order_id AS OrderId,
			o.order_date AS OrderDate,
			o.customer_id AS CustomerId,
			c.customer_id AS Customer__CustomerId,
			c.name AS Customer__Name,
			c.email AS Customer__Email,
			oi.order_item_id AS OrderItems__OrderItemId,
			oi.product_id AS OrderItems__ProductId,
			oi.quantity AS OrderItems__Quantity,
			p.product_id AS OrderItems__Product__ProductId,
			p.name AS OrderItems__Product__Name,
			p.price AS OrderItems__Product__Price
		FROM
			orders o
		LEFT JOIN
			customers c ON o.customer_id = c.customer_id
		LEFT JOIN
			order_items oi ON o.order_id = oi.order_id
		LEFT JOIN
			products p ON oi.product_id = p.product_id
	`

	// Execute the query
	rows, err := oc.rawdb.Query(query)
	if err != nil {
		errorMessage := fmt.Sprintf("error executing query: %v", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}
	defer rows.Close()

	// Initialize a slice to store the results
	var orders []OrderViewModel

	// Iterate over the rows
	for rows.Next() {
		// Define variables to store the row values
		var order OrderViewModel
		// var customer Customer
		// var orderItem OrderItem
		// var product Product
		var customer struct {
			CustomerId int
			Name       string
			Email      string
		}
		var orderItem struct {
			OrderItemId int
			ProductId   int
			Quantity    int
			Product     struct {
				ProductId int
				Name      string
				Price     float64
			}
		}
		var product struct {
			ProductId int
			Name      string
			Price     float64
		}

		// Scan the row into variables
		var orderDate []uint8
		err := rows.Scan(
			&order.OrderId,
			&orderDate,
			&order.CustomerId,
			&customer.CustomerId,
			&customer.Name,
			&customer.Email,
			&orderItem.OrderItemId,
			&orderItem.ProductId,
			&orderItem.Quantity,
			&product.ProductId,
			&product.Name,
			&product.Price,
		)
		if err != nil {
			errorMessage := fmt.Sprintf("error scanning row: %v", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		// Convert orderDate to time.Time
		// order.OrderDate, err = time.Parse("2006-01-02", string(orderDate))
		order.OrderDate, err = time.Parse("2006-01-02 15:04:05.000", string(orderDate))
		if err != nil {
			errorMessage := fmt.Sprintf("error parsing order date: %v", err)
			log.Println(errorMessage)
			return c.Status(http.StatusInternalServerError).SendString(errorMessage)
		}

		// Assign the customer and product to the order
		order.Customer = customer
		orderItem.Product = product
		order.OrderItems = append(order.OrderItems, orderItem)

		// Append the order to the slice
		orders = append(orders, order)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		errorMessage := fmt.Sprintf("error iterating over rows: %v", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Process the results to calculate the total for each order
	for i := range orders {
		var total float64
		for _, item := range orders[i].OrderItems {
			total += float64(item.Quantity) * item.Product.Price
		}
		orders[i].Total = total
	}

	// Return the orders
	// return orders, nil
	return c.JSON(orders)
}
