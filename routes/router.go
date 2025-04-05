package routes

import (
	"database/sql"
	"gorm-learning/controllers"

	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// RegisterRoutes registers all routes for the application
func RegisterRoutes(app *fiber.App, ormdb *gorm.DB, rawdb *sql.DB) {
	orderController := controllers.NewOrderController(ormdb, rawdb)
	customerController := controllers.NewCustomerController(ormdb, rawdb)
	productController := controllers.NewProductController(ormdb, rawdb)

	// Define routes
	app.Get("/Customer/", customerController.GetAllCustomers)
	app.Get("/Product/", productController.GetAllProducts)

	app.Get("/orders-miss-product", orderController.MissingProductMappingRetrieveOrdersIncludeAllPropsGorm)
	app.Get("/orders-worked", orderController.WorkedRetrieveOrdersIncludeAllPropsGorm)
	app.Get("/orders-not-grouped", orderController.NotGroupedRetrieveOrdersIncludeAllPropsGorm)
	app.Get("/orders-raw-not-grouped", orderController.NotGroupedGetOrdersWithRawSQL)
	// app.Get("/orders-simple-orm", orderController.GetOrdersSimple)
	app.Get("/orders-simple-orm-anonymous", orderController.GetOrdersSimpleAnonymousObj)

	// Add more routes here as needed
}
