package controllers

import (
	"database/sql"
	"fmt"
	"gorm-learning/models"
	"log"
	"net/http"

	// . "gorm-learning/models"

	// "log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// CustomerController handles HTTP requests related to orders
type CustomerController struct {
	ormdb *gorm.DB
	rawdb *sql.DB
}

/*
Docker Container logs:
2025-04-05 23:56:56 app-1  | 2025/04/05 16:56:56 No customers found
2025-04-05 23:57:09 app-1  | 2025/04/05 16:57:09 No customers found
*/

// NewOrderController creates a new instance of OrderController
func NewCustomerController(ormdb *gorm.DB, rawdb *sql.DB) *CustomerController {
	return &CustomerController{ormdb: ormdb, rawdb: rawdb}
}

func (uc *CustomerController) GetAllCustomers(c fiber.Ctx) error {
	var customers []models.Customer // must specify package "models.xxx" for VSCode Go auto import to act

	// Query all users from the database
	if err := uc.ormdb.Find(&customers).Error; err != nil {
		errorMessage := fmt.Sprintf("Error retrieving customers: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Check if result is empty
	if len(customers) == 0 {
		log.Println("No customers found")
	}

	// Return the customers as JSON
	return c.JSON(customers)
}
