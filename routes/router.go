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

	// Define routes
	app.Get("/orders-miss-product", orderController.MissingProductMappingRetrieveOrdersIncludeAllPropsGorm)
	app.Get("/orders-worked", orderController.WorkedRetrieveOrdersIncludeAllPropsGorm)
	app.Get("/orders-wrong", orderController.WrongRetrieveOrdersIncludeAllPropsGorm)
	app.Get("/orders-raw", orderController.GetOrdersWithRawSQL)
	// app.Get("/orders-simple-orm", orderController.GetOrdersSimple)
	app.Get("/orders-simple-orm-anonymous", orderController.GetOrdersSimpleAnonymousObj)

	// Add more routes here as needed
}
