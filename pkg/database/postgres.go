package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetPostgresDB - returns connection to the postgres database
func GetPostgresDB(postgresUrl string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(postgresUrl))

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Check whether db works correctly
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	err = sqlDB.Ping()

	if err != nil {
		log.Fatal("Failed to ping:", err)
	}

	return db
}
