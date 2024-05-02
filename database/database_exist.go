package database

import "github.com/jinzhu/gorm"

// DatabaseExists checks if the specified database exists
func DatabaseExists(db *gorm.DB, databaseName string) bool {
	var exists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?)", databaseName).Row().Scan(&exists)
	return exists
}
