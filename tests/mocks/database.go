package mocks

import (
	"log"

	"github.com/rasadov/subscription-manager/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect to test database:", err)
	}

	err = db.AutoMigrate(&models.Subscription{})
	if err != nil {
		log.Fatal("failed to migrate test database:", err)
	}

	return db
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func ResetDB(db *gorm.DB) error {
	log.Println("Resetting test database...")
	if db == nil {
		return nil
	}
	err := db.Migrator().DropTable(&models.Subscription{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Subscription{})
	return err
}
