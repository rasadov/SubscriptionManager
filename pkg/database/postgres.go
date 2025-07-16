package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB - returns connection to the postgres database
func NewPostgresDB(
	host string,
	user string,
	password string,
	dbName string,
	port int,
	sslMode string,
) (*gorm.DB, error) {
	postgresUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		host, user, password, dbName, port, sslMode)
	db, err := gorm.Open(postgres.Open(postgresUrl))

	if err != nil {
		return nil, err
	}

	// Check whether db works correctly
	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
