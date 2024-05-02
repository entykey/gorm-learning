package main

import (
	"database/sql"
	"fmt"
	. "gorm-learning/database"
	. "gorm-learning/models"
	"gorm-learning/routes"

	"log"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v3"

	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB   // orm
var rawdb *sql.DB // raw

func initORM_DB() {
	// Create colorized loggers
	infoLogger := log.New(color.Output, color.GreenString("INFO: "), log.LstdFlags)
	errorLogger := log.New(color.Output, color.RedString("ERROR: "), log.LstdFlags)

	// Database connection string
	databaseName := "gorm_learning"
	dbString := fmt.Sprintf("user:password@tcp(localhost:3306)/mysql?charset=utf8&parseTime=True&loc=Local")

	// NOTE: use a seperate "db" instance to read first
	// Connect to the default "mysql" database
	defaultdb, err := gorm.Open("mysql", dbString)

	// for gorm.io's gorm
	// defaultdb, err := gorm.Open(mysql.Open(dbString), &gorm.Config{})
	if err != nil {
		errorLogger.Fatal("Failed to connect to default database:", err)
	}

	// defer defaultdb.Close()	// removed in GORM v2
	infoLogger.Printf("Successfully connected to the default database")

	// Check if the database exists
	if !DatabaseExists(defaultdb, databaseName) {
		// Create the database if it doesn't exist
		infoLogger.Printf("Database not found for name %s, creating database...", databaseName)
		if err := defaultdb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", databaseName)).Error; err != nil {
			errorLogger.Fatal("Failed to create database:", err)
		}
		infoLogger.Printf("Database with name %s created.", databaseName)
	} else {
		infoLogger.Printf("Database with name %s already exists.", databaseName)
	}

	// NOTE: ormdb instance only start from here
	// Reconnect to the newly created database
	// dsn := fmt.Sprintf("user:password@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseName)

	// ERROR: This cause `db` to be nil
	// because this line creates a new local variable instead of assigning a value to the package-level variable db.
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// FIX: assign the database connection to the package-level variable db
	// db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err = gorm.Open("mysql", fmt.Sprintf("user:password@tcp(localhost:3306)/%s?charset=utf8&parseTime=True&loc=Local", databaseName))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the models
	infoLogger.Println("Applying automatic migrations...")
	db.AutoMigrate(&Customer{}, &Order{}, &Product{}, &OrderItem{})
	infoLogger.Println("Migrations applied successfully!")

	// Seed the database
	if err := SeedData(db); err != nil {
		log.Fatal("Failed to seed data:", err)
	}
	infoLogger.Println("Data seeded successfully!")
}

// Function to initialize the raw DB connection
func initRawDB() {
	var err error
	rawdb, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/gorm_learning")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set maximum number of connections in idle connection pool.
	rawdb.SetMaxIdleConns(10)

	// Set maximum number of open connections to the database.
	rawdb.SetMaxOpenConns(10)
}

func main() {
	// Initialize database connection
	initORM_DB()
	initRawDB()
	defer db.Close() // removed in gorm.io registry
	defer rawdb.Close()

	// Create Fiber app
	app := fiber.New()

	// Register routes
	routes.RegisterRoutes(app, db, rawdb)

	// Start Fiber server
	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf("Error starting Fiber server: %v", err)
	}
}
