package controllers

import (
	"database/sql"
	"fmt"
	"gorm-learning/models"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// CustomerController handles HTTP requests related to orders
type ProductController struct {
	ormdb *gorm.DB
	rawdb *sql.DB
}

/*
Docker Container logs:
2025-04-06 00:47:59 2025/04/05 17:47:59 No products found
2025-04-06 00:48:02 2025/04/05 17:48:02 No products found
*/

// NewOrderController creates a new instance of OrderController
func NewProductController(ormdb *gorm.DB, rawdb *sql.DB) *ProductController {
	return &ProductController{ormdb: ormdb, rawdb: rawdb}
}

func (uc *ProductController) GetAllProducts(c fiber.Ctx) error {
	var products []models.Product // must specify package "models.xxx" for VSCode Go auto import to act

	// Query all users from the database
	if err := uc.ormdb.Find(&products).Error; err != nil {
		errorMessage := fmt.Sprintf("Error retrieving products: %s", err)
		log.Println(errorMessage)
		return c.Status(http.StatusInternalServerError).SendString(errorMessage)
	}

	// Check if result is empty
	if len(products) == 0 {
		log.Println("No products found")
	}

	// Return the products as JSON
	return c.JSON(products)
}
