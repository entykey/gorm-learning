package main

import (
	"database/sql"
	"fmt"
	"gorm-learning/database"
	. "gorm-learning/database"
	. "gorm-learning/models"
	"gorm-learning/routes"
	"os"
	"time"

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

	// Get database credentials from environment
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// First connect to MySQL server (without specifying a database)
	dbString := fmt.Sprintf("%s:%s@tcp(%s:3306)/?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost)

	var defaultdb *gorm.DB
	var err error
	maxRetries := 20
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		defaultdb, err = gorm.Open("mysql", dbString)
		if err == nil {
			break
		}

		errorLogger.Printf("Attempt %d: Failed to connect to MySQL server: %v", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		errorLogger.Fatal("Failed to connect to MySQL server after retries:", err)
	}

	// Check if database exists
	if !DatabaseExists(defaultdb, dbName) {
		infoLogger.Printf("Creating database %s...", dbName)
		if err := defaultdb.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
			errorLogger.Fatal("Failed to create database:", err)
		}
	}

	// Now connect to the specific database
	dbString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbName)

	db, err = gorm.Open("mysql", dbString)
	if err != nil {
		errorLogger.Fatal("Failed to connect to database:", err)
	}

	// Create application user with proper privileges if using root initially
	if dbUser == "root" {
		db.Exec("CREATE USER IF NOT EXISTS 'user'@'%' IDENTIFIED BY 'password'")
		db.Exec("GRANT ALL PRIVILEGES ON gorm_learning.* TO 'user'@'%'")
		db.Exec("FLUSH PRIVILEGES")
	}

	// Auto-migrate the models
	infoLogger.Println("Applying automatic migrations...")
	db.AutoMigrate(&Customer{}, &Order{}, &Product{}, &OrderItem{})
	infoLogger.Println("Migrations applied successfully!")

	// Initialize seeder
	seeder := database.NewSeeder(db)

	// Seed data
	if err := seeder.SeedAll(); err != nil {
		errorLogger.Fatalf("Seeding failed: %v", err)
	}
	infoLogger.Println("✅ Database seeding completed")
}

// func initORM_DB() {
// 	// Create colorized loggers
// 	infoLogger := log.New(color.Output, color.GreenString("INFO: "), log.LstdFlags)
// 	errorLogger := log.New(color.Output, color.RedString("ERROR: "), log.LstdFlags)

// 	// (1) Database connection string (hardcoded, local machine's installed phpmyadmin's mysql)
// 	// dbName := "gorm_learning"
// 	// dbString := fmt.Sprintf("user:password@tcp(localhost:3306)/mysql?charset=utf8&parseTime=True&loc=Local")

// 	// (2) Docker based mysql service
// 	// Get database credentials from environment
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbName := os.Getenv("DB_NAME")

// 	// Database connection string
// 	dbString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
// 		dbUser, dbPassword, dbHost, dbName)

// 	// NOTE: use a seperate "db" instance to read first
// 	// Connect to the default "mysql" database
// 	defaultdb, err := gorm.Open("mysql", dbString)

// 	// for gorm.io's gorm
// 	// defaultdb, err := gorm.Open(mysql.Open(dbString), &gorm.Config{})
// 	if err != nil {
// 		errorLogger.Fatal("Failed to connect to default database:", err)
// 	}

// 	// defer defaultdb.Close()	// removed in GORM v2
// 	infoLogger.Printf("Successfully connected to the default database")

// 	// Check if the database exists
// 	if !DatabaseExists(defaultdb, dbName) {
// 		// Create the database if it doesn't exist
// 		infoLogger.Printf("Database not found for name %s, creating database...", dbName)
// 		if err := defaultdb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)).Error; err != nil {
// 			errorLogger.Fatal("Failed to create database:", err)
// 		}
// 		infoLogger.Printf("Database with name %s created.", dbName)
// 	} else {
// 		infoLogger.Printf("Database with name %s already exists.", dbName)
// 	}

// 	// NOTE: ormdb instance only start from here
// 	// Reconnect to the newly created database
// 	// dsn := fmt.Sprintf("user:password@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbName)

// 	// ERROR: This cause `db` to be nil
// 	// because this line creates a new local variable instead of assigning a value to the package-level variable db.
// 	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	// FIX: assign the database connection to the package-level variable db
// 	// db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	db, err = gorm.Open("mysql", fmt.Sprintf("user:password@tcp(localhost:3306)/%s?charset=utf8&parseTime=True&loc=Local", dbName))
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	// Auto-migrate the models
// 	infoLogger.Println("Applying automatic migrations...")
// 	db.AutoMigrate(&Customer{}, &Order{}, &Product{}, &OrderItem{})
// 	infoLogger.Println("Migrations applied successfully!")

// 	// Seed the database
// 	if err := SeedData(db); err != nil {
// 		log.Fatal("Failed to seed data:", err)
// 	}
// 	infoLogger.Println("Data seeded successfully!")
// }

// Function to initialize the raw DB connection
// func initRawDB() {
// 	var err error
// 	rawdb, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/gorm_learning")
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	// Set maximum number of connections in idle connection pool.
// 	rawdb.SetMaxIdleConns(10)

// 	// Set maximum number of open connections to the database.
// 	rawdb.SetMaxOpenConns(10)
// }

func initRawDB() {
	// Get database credentials from environment
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	var err error
	rawdb, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Verify connection
	err = rawdb.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	rawdb.SetMaxIdleConns(10)
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

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // default port
	}

	// Start Fiber server on all interfaces
	err := app.Listen("0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("Error starting Fiber server: %v", err)
	}
}
