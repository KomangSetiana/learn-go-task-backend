package config

import (
	"fmt"
	"log"
	"os"
	"tusk/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseConnection() *gorm.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return database
}

func CreateOwnerAccount(db *gorm.DB) {
	hasedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	owner := models.User{
		Role:     "Admin",
		Name:     "Owner",
		Password: string(hasedPasswordBytes),
		Email:    "owner@go.id",
	}

	if db.Where("email=?", owner.Email).First(&owner).RowsAffected == 0 {
		db.Create(&owner)
	} else {
		fmt.Println("Owner exists")
	}
}
