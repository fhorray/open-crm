package db

import (
	"fmt"
	"log"
	"open-crm/core/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=secret dbname=open-crm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error trying to connect with database: ", err)
	}

	// Migrate DB
	db.AutoMigrate((&models.Lead{}))
	db.AutoMigrate((&models.User{}))

	DB = db
	fmt.Println("✅ DB Connected")
}
