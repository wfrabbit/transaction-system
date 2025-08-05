package db

import (
	"log"

	"wangfeng/transaction-system/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn *string) error {
	if dsn == nil {
		log.Fatal("database URL is not provided")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate models
	if err := DB.AutoMigrate(&model.Account{}, &model.Transaction{}); err != nil {
		return err
	}

	log.Println("Database connected and migrated successfully")
	return nil
}
